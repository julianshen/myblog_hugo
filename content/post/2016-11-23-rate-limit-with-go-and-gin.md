---
date: "2016-11-23T12:05:38Z"
tags:
- golang
- web
- cloud developing
- Gin
- Go
title: Rate limit with Go and Gin
---

昨天趁等著去面試前稍微把這之前想要寫一下的這題目打包成一個[Gin](https://github.com/gin-gonic/gin)的middleware :

* [Gin-limiter](https://github.com/julianshen/gin-limiter) 

Rate limiting 通常在很多開放API的服務內會常看到, 像是[Twitter](https://dev.twitter.com/rest/public/rate-limiting),
像是[Facebook](https://developers.facebook.com/docs/graph-api/advanced/rate-limiting)或是[新浪微博](http://open.weibo.com/wiki/Rate-limiting),
其目的就是希望API不要被特定節點頻繁存取以致於造成伺服器端的過載

### rate limiter ###

一般的Rate limiting的設計大致上來說就是限制某一個特定的節點(或使用者或API Key等等)在一段特定的時間內的存取次數,
比如說, 限制一分鐘最多60次存取這樣的規則, 最直覺的方式我們是可以起一個timer和一個counter, counter大於60就限制存取, timer則每60秒重置counter一次,
看似這樣就好了, 但其實這有漏洞, 假設我在第59秒時瞬間存取了60次, 第61秒又瞬間存取了60次, 在這設計上是合法的, 因為counter在第60秒時就被重置了,
但實質上卻違反了一分鐘最多60次這限制, 因為他在兩秒內就存取了120次, 遠大於我們設計的限制, 當然我們也可以用Sliding time window來解決,
但那個實作上就稍稍複雜點

目前兩個主流比較常見的做法是[Token Bucket](https://en.wikipedia.org/wiki/Token_bucket)和[Leaky Bucket](https://en.wikipedia.org/wiki/Leaky_bucket),
這兩個原理上大同小異

先來說說[Token Bucket](https://en.wikipedia.org/wiki/Token_bucket), 他的做法是, 假設你有個桶子, 裡面是拿來裝令牌(Token)的,
桶子不是Doraemon的四次元口袋, 所以他空間是有限的, 令牌(Token)的作用在於, 要執行命令的人, 如果沒從桶子內摸到令牌, 就不准執行,
然後我們一段時間內丟一些令牌進去, 如果桶子裡面已經裝滿就不丟, 以上個例子來說, 我們可以準備一個最多可以裝60個令牌的桶子, 每秒鐘丟一個進去,
如果消耗速度大於每秒一個, 自然桶子很快就乾了, 就沒牌子拿了

[Leaky Bucket](https://en.wikipedia.org/wiki/Leaky_bucket)跟[Token Bucket](https://en.wikipedia.org/wiki/Token_bucket)很像, 不過就是反過來, 
我們把每次的存取都當作一滴水滴入桶子中, 桶子滿了就會溢出(拒絕存取), 桶子底下打個洞, 讓水以固定速率流出去, 這樣一樣能達到類似的效果

![Leaky Bucket](https://upload.wikimedia.org/wikipedia/commons/c/c4/Leaky_bucket_analogy.JPG)

### Go的rate limiter實作 ###

Go官方的package內其實是有rate limiter的實作的: 

* [package rate](https://godoc.org/golang.org/x/time/rate)

照他的說法他是實作了[Token Bucket](https://en.wikipedia.org/wiki/Token_bucket), [創建一個Limiter](https://godoc.org/golang.org/x/time/rate#NewLimiter),
要給的參數是[Limit](https://godoc.org/golang.org/x/time/rate#Limit)和b, 這個Limit指的是每秒鐘丟多少Token進桶字(? 我不知道有沒理解錯),
而b是桶子的大小

實際上去用了之後發現好像也不是那麼好用, 可能我理解有問題, 出現的並不是我想像的結果, 因此我換用了[Juju's ratelimit](https://github.com/juju/ratelimit),
這個是在[gokit](https://github.com/go-kit/kit)這邊看到它有用, 所以應該不會差到哪去, 一樣也是Token Bucket,給的參數就是多久餵一次牌子, 跟桶子的大小, 這就簡單用了一點

### 包裝成Gin middleware ###

要套在web server上使用的話, 包裝成middleware是比較方便的, 因此我就花了點時間把[Juju's ratelimit](https://github.com/juju/ratelimit)包裝成這個:

* [Gin-limiter](https://github.com/julianshen/gin-limiter) 

使用範例如下:

```go
    //Allow only 10 requests per minute per API-Key
	lm := limiter.NewRateLimiter(time.Minute, 10, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-API-KEY")
		if key != "" {
			return key, nil
		}
		return "", errors.New("API key is missing")
	})
	//Apply only to /ping
	r.GET("/ping", lm.Middleware(), func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})

	//Allow only 5 requests per second per user
	lm2 := limiter.NewRateLimiter(time.Second, 5, func(ctx *gin.Context) (string, error) {
		key := ctx.Request.Header.Get("X-USER-TOKEN")
		if key != "" {
			return key, nil
		}
		return "", errors.New("User is not authorized")
	})

	//Apply to a group
	x := r.Group("/v2")
	x.Use(lm2.Middleware())
	{
		x.GET("/ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong",
			})
		})
		x.GET("/another_ping", func(c *gin.Context) {
			c.JSON(200, gin.H{
				"message": "pong pong",
			})
		})
	}
```

這邊本來想說為了簡單理解一點把參數設計成"每分鐘不能超過10次"這樣的描述, 然後後面再轉換成實際的fillInterval, 不過好像覺得怪怪的,
有點不太符合Token Bucket的特質, 寫成middleware後的彈性就較大一點, 可以全部都用一個limiter或是分不同的資源不同限制都可

這邊建構時要傳入一個用來產生key的函數, 這是考慮到每個人想限制的依據不同, 例如根據API key, 或是session, 或是不同來源之類的, 由這函數去
產生一個對應的key來找到limiter, 如果傳回error, 就是這個request不符合規則, 直接把他拒絕掉

### 跨server的rate limit ###

這方法只能針對單一server, 但現在通常都是多台server水平擴展, 因此也是會需要橫跨server的解決方案, 這部分的話, 用Redis來實作Token Bucket是可行的, 這等下次再來弄好了