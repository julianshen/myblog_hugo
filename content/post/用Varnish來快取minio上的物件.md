---
date: 2022-10-08T16:10:26+08:00
title: "用Varnish來快取minio上的物件"
images: 
- "https://og.jln.co/jlns1/55SoVmFybmlzaOS-huW_q-WPlm1pbmlv5LiK55qE54mp5Lu2"
draft: false
---

剛好想說要用Varish來做一下Minio(S3)的cache, 研究了一下順便做個紀錄

先在Ubuntu上裝來試試, 可以用`apt` 來安裝：

```
apt install varnish
```

在Ubuntu 20.04上(WSL用的版本)是6.2.1的版本, 最新版應該是7.2, 不過沒差, 做法都一樣

Varnish預設的設定在`/etc/varnish/default.vcl`, 打開這檔案你就可以看到像這樣的內容:

```vcl
vcl 4.0;

# Default backend definition. Set this to point to your content server.
backend default {
    .host = "127.0.0.1";
}
sub vcl_recv {
    # Happens before we check if we have this in cache already.
    #
    # Typically you clean up the request here, removing cookies you don't need,
    # rewriting the request, etc.
}

sub vcl_backend_response {
    # Happens after we have read the response headers from the backend.
    #
    # Here you clean the response headers, removing silly Set-Cookie headers
    # and other mistakes your backend does.
}

sub vcl_deliver {
    # Happens when we have all the pieces we need, and are about to send the
    # response to the client.
    #
    # You can do accounting or modifying the final object here.
}
```

Varnish的設定檔用的是一種叫做vcl的語言, 它會被Varnish先compile過後才會被使用, 所以改好這檔案後, 如果你跑 `sudo system start varnish` (這是WSL2上用的, 其他地方可能就是`systemctl`), 如果你寫錯了, 一開始跑就可以發現出錯了

以上面那個例子來說, 它會預設快取你local上的web server

但如果是要連接Minio (S3)是不夠的, 因為如果單純把backend設成 Minio server, 那client還是會需要access key和secret key才可以存取, 如果你希望讓它跟存取靜態網站一樣, 那你可以能會希望把這兩個設定放在後端

Vanish出場是沒支援可以call S3 API的, 這時候就要透過一個[VMOD - AWSRest](https://github.com/xcir/libvmod-awsrest), 這VMOD是可以在你去backend (Minio/S3) 拿資料前先幫你用你的access key, secret key算好簽章(signature), 所以我們要先安裝這個VMOD

安裝VMOD你會先需要`libvarnishapi-dev`, 可以用`apt install libvarnishapi-dev`來安裝, 另外AWSRest還會需要mhash, 你還會需要安裝`apt-get install libmhash-dev`

裝好後, 從 https://github.com/xcir/libvmod-awsrest 抓取最新的source code, 進入目錄後執行

```
./autogen.sh
./configure
make
sudo make install
```

沒意外的話就可以完成安裝, 要確認是不是已經安裝好了, 我們可以在default.vcl加上

```vcl
vcl 4.0;
import awsrest;

# ....
```

重啟varnish有成功, 表示應該是沒啥問題才對

我在我本地端電腦跑了個Minio, port為9000, 有一個bucket叫做`mmmbux`, 裡面有個檔案, key為`20220101/a.c`, access key為`TGhYs2FYBGMYueAz`, secrect key為`IM2SgF7LxIlZVbeo3Vv7OdQzA7pnZFB1`, Varnish則是跑在port 6081上

首先我們來看看怎讓client/browser在不用提供access key/secret key的狀況下可以存取物件

``` vcl
vcl 4.0;
import awsrest;

backend default {
    .host = "127.0.0.1";
    .port = "9000";
}

sub vcl_recv {
    set req.http.host = "127.0.0.1";
    awsrest.v4_generic(
        service           = "s3",
        region            = "ap-northeast-1",
        access_key        = "TGhYs2FYBGMYueAz",
        secret_key        = "IM2SgF7LxIlZVbeo3Vv7OdQzA7pnZFB1",
        signed_headers    = "host;",
        canonical_headers = "host:" + req.http.host + awsrest.lf()
    );
}
```

Ok, 其實就很簡單的在vcl_recv上加上那幾行就好, 這時候你就可以用 `http://localhost:6081/mmmbux/20220101/a.c` 來存取 `mmmbux` 這bucket上 `2022/0101/a.c` 這個檔案了

那, 如果我不想把bucket name當作url的一部分呢?

```vcl
sub vcl_recv {
    set req.http.host = "127.0.0.1";
    set req.url = "mmmbux/" + req.url
    awsrest.v4_generic(
        service           = "s3",
        region            = "ap-northeast-1",
        access_key        = "TGhYs2FYBGMYueAz",
        secret_key        = "IM2SgF7LxIlZVbeo3Vv7OdQzA7pnZFB1",
        signed_headers    = "host;",
        canonical_headers = "host:" + req.http.host + awsrest.lf()
    );
}
```

上面這段就是把你進來的url加上`mmmbux/`當新的url, 這樣做的話, 你的新url就會是 `http://localhost:6081/20220101/a.c`

那如果我想進一步, 把它變成 `http://localhost:6081/files/20220101/a.c` 呢?

```vcl
sub vcl_recv {
    set req.http.host = "127.0.0.1";
    
    if (req.url ~ "^/files/") {
        set req.url = regsub(req.url, "^/files/", "/mmmbux/");
        awsrest.v4_generic(
          service           = "s3",
          region            = "ap-northeast-1",
          access_key        = "TGhYs2FYBGMYueAz",
          secret_key        = "IM2SgF7LxIlZVbeo3Vv7OdQzA7pnZFB1",
          signed_headers    = "host;",
          canonical_headers = "host:" + req.http.host + awsrest.lf()
        );
    } else {
        return(synth(404));
    }
}
```

上面這段就是把`/files/`後面的都到 `mmmbux`這bucket去抓, 然後其他目錄都回傳 `404 Not found`