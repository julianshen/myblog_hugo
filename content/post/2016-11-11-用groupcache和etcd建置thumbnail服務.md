---
date: "2016-11-11T14:53:03Z"
tags:
- Golang
title: 用Groupcache和etcd建置thumbnail服務
---
在之前工作的時候, 做了一個專門用來產生thumbnail(縮圖)的服務, 當時這東西主要的目的是為了因應[Zencircle](http://www.zencircle.com)會有不同尺寸的縮圖的需求,
而且每次client app改版又可能多新的尺寸, 因此當時寫了這個叫Minami的服務, 當時幾個簡單的需求是:

1. 要能夠被CDN所cache (因此URL設定上不採用query string,而是簡單的URL)
1. 能夠容易被deploy
1. 能夠的簡單的被擴展 (加一台新的instance就可以)
1. 不需要太多額外的dependencies

不過那時候寫的版本, 沒寫得很好, 這兩天花了點時間重寫了一個叫做[Minami_t](https://github.com/julianshen/minami_t)(本來Minami這名字就是來自於Minami Takahashi, 所以加個"t" XD),
新的這個重寫的版本採一樣的架構(使用了[groupcache](https://github.com/golang/groupcache)), 但多加了Peer discovery的功能(使用[etcd](https://github.com/coreos/etcd)), 但少了
臉部辨識跟色情照片偵測功能(原本在前公司的版本有, 新寫的這個我懶得加了)

我把這次重寫的版本放到github上: [Minami_t](https://github.com/julianshen/minami_t)

不過這算是一個sample project, 影像來源來自於Imgur, 如何使用或如何改成支援自己的Image host, 那就自行看source code吧, 這版本縮圖的部分用了我改過的[VIPS](https://github.com/julianshen/vips),
當然原來版本的VIPS也是可用, 這版本只是我當初為了支援Face crop所改出來的

### Groupcache ###

先來說說為什麼採用[groupcache](https://github.com/golang/groupcache)? 我不是很確定當時為何會看到[groupcache](https://github.com/golang/groupcache)這來, 但後來想想, 採用它的原因可能是看到[這份投影片](https://talks.golang.org/2013/oscon-dl.slide#43),
它是memchached的作者寫來用在dl.google.com上面的, 架構上剛好也適合thumbnail service, 可能剛好投影片又提到thumbnail(我腦波也太弱了吧), 所以當初採用它來實作這個service

架構上會像是這樣:

![](/images/posts/minami1.001.jpeg)

Groupcache有幾個特色

1. Embedded, 不像memcached, redis需要額外的server, 它是嵌入在你原本的程式內的
1. Shared, Cache是可以所有Peer共享的, 資料未必放在某特定的Peer上, 有可能在本機, 也可能在另一台, 當然如果剛好在本機時就會快一點
1. LRU, Cache總量有上限限制的, 過久沒使用的資料有可能會被移出記憶體
1. Immutable, key所對應的值不像memcached, redis可以修改, 而是當cache miss時, 他會再透過你實作的getter去抓真正的資料

要讓Groupcache可以在不同node間共享cache, 就必須開啟HTTPPool, 像下面

```golang
	ln, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return nil, err
	}

	if port == 0 {
		port = ln.Addr().(*net.TCPAddr).Port
	}

	_url := fmt.Sprintf("http://%s:%d", ip, port)
	pool := groupcache.NewHTTPPool(_url)

	go func() {
		log.Printf("Group cache served at port %s:%d\n", ip, port)
		if err := http.Serve(ln, http.HandlerFunc(pool.ServeHTTP)); err != nil {
			log.Printf("GROUPCACHE PORT %d, ERROR: %s\n", port, err.Error())
			os.Exit(-1)
		}
	}()
```

Groupcache 的getter範例:

```golang
func getter(ctx groupcache.Context, key string, dest groupcache.Sink) error {
	log.Println("Cache missed for " + key)

	params := strings.Split(key, ":")

	if len(params) != 3 {
		return ErrWrongKey
	}

	d := ctx.(*Downloader)
	fileName, err := d.Download("http://i.imgur.com/" + params[2])

	if err != nil {
		return err
	}

	//Should assume correct since it is checked at where it is from
	width, _ := strconv.Atoi(params[0])
	height, _ := strconv.Atoi(params[1])

	data, err := resize(fileName, width, height)

	if err != nil {
		return err
	}

	dest.SetBytes(data)
	return nil
}
```

### etcd ###

我之前寫的版本有個問題是, 沒有自動的peer discovery的功能, 所以必須手動加peer, 這版本把etcd導入, etcd已經是coreos的核心之一了, 簡單, 又蠻好用的,
不過選它也是它直接有Go的client了

Peer discovery的部分, 參考了[Go kit](https://github.com/go-kit/kit)的[etcd實作](https://github.com/go-kit/kit/tree/master/sd/etcd),
[Go kit](https://github.com/go-kit/kit)是一個蠻好的Go的微服務框架, 它裡面也有實作用etcd做service discovery, 這一部分正好是這邊需要的, 因此
參考並寫出了這邊這個版本

重點是要能夠在有新server加入後就新增到peer list去, 有server離開後要拿掉, 因此必須利用到etcd的watch功能

```golang
func (s *ServiceRegistry) Watch(watcher Watcher) {
	key := fmt.Sprintf("/%s/nodes", s.name)
	log.Println("watch " + key)
	w := s.etcd_client.Watcher(key, &etcd.WatcherOptions{AfterIndex: 0, Recursive: true})

	var retryInterval time.Duration = 1

	for {
		_, err := w.Next(s.ctx)

		if err != nil {
			log.Printf("Failed to connect to etcd. Will retry after %d sec \n", retryInterval)
			time.Sleep(retryInterval * time.Second)

			retryInterval = (retryInterval * 2) % 4096
		} else {
			if retryInterval > 1 {
				retryInterval = 1
			}

			list, err := s.GetNodes()
			if err == nil {
				watcher(list)
			} else {
				//skip first
			}
		}
	}
}
````

Watch可以用來監測某一個key有無改變, 因此我們只要一直監測server node的list就好(指定一個key來放), 因此流程是這樣的:

1. Server開啟後, 自己到etcd註冊自己, 並把etcd上找到的nodes全加到peer list中
1. 另一台由etcd發現有另一台出現後, 把它加到peer list中
1. Server下線後, 要移除自己的註冊, 其他機器要從peer list把它移除

問題點在最後一點, Server下線有可能是被kill的, 也有可能按ctrl-c中斷的, 這時候就要監聽os的signal,
在程式被結束前, 可以先去移除註冊, 像這樣:

```golang
//Listening to exit signals for graceful leave
go func() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-c
	log.Println("I'm leaving")
	cm.Leave()
	os.Exit(0)
}()
```

這只是一個sample而已, 還有一些待改進的