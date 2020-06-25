---
date: "2017-02-04T16:17:19Z"
images:
- /images/posts/2017-02-04-用ifttt-+-heroku跑定時的tasks.md.jpg
tags:
- ifttt
- heroku
title: 用ifttt + heroku跑定時的tasks
---

這篇是延續"[使用AWS lambda和Github來提供中華職棒賽程資料](http://blog.jln.co/%E4%BD%BF%E7%94%A8aws-lambda%E5%92%8Cgithub%E4%BE%86%E6%8F%90%E4%BE%9B%E4%B8%AD%E8%8F%AF%E8%81%B7%E6%A3%92%E8%B3%BD%E7%A8%8B%E8%B3%87%E6%96%99/)",
之前的做法是用Cloud watch加上lambda來做這件事, 但我跑的東西並不是那麼的頻繁, 在AWS上還是會被收到流量的費用,
因此就打算用更經濟的方式, 利用heroku免費的額度來做這事(真是壞客戶 XD)

目的是定時(比如說每四小時)去爬一些網頁的資訊, 爬這些網頁其實也不需要花太久時間

用Cloud watch + lambda的好處是不用架一台server, 但用Heroku這種PAAS其實也不用太去管server這事

Heroku是可以設定[scheduled tasks](https://devcenter.heroku.com/articles/scheduled-jobs-custom-clock-processes)的, 但額外的work dyno是要另外付費的,
因此, 如果需求不是需要太頻繁, 也不需要執行太久的, 這時候就可以利用[ifttt](https://ifttt.com/)來定時觸發一個url的方式來做

![](/images/posts/ifttt1.png)

要定時觸發一個URL, ifttt applet該怎麼設定呢? 首先"this"要選用的是Date & Time, 如下:

![](/images/posts/ifttt_time.png)

設定上並沒有很多, 就像是每小時, 每天之類的, 沒辦法訂多個, 如果需要一次多個設定, 那就多新增幾個Applets吧

![](/images/posts/ifttt_time_sel.png)

這邊設定每小時, 就設定每小時的15分來觸發吧

![](/images/posts/ifttt_set_time.png)

那"that"的動作呢? 觸發URL的動作是利用"Maker", 這是設計給iot用的吧, 不過, 拿來做這用途也是沒問題的:

![](/images/posts/ifttt_maker.png)

Maker只有一個選項"Make a web request"

![](/images/posts/ifttt_choose_action.png)

設定很單純, 就給定URL, 使用的HTTP Method(GET, POST, PUT ...), Content-type, 跟Body

![](/images/posts/ifttt_make_req.png)

這邊我用的是POST + Json, Json裡面會帶一個TOKEN來辨識, 以免有心人士利用了這個URL, go的檢查範例如下:

```go
func checkID(body io.Reader) bool {
        data, err := ioutil.ReadAll(body)

        if err != nil {
                return false
        }

        var rbody struct {
                Id string
        }
        err = json.Unmarshal(data, &rbody)

        if err != nil {
                return false
        }

        return rbody.Token == os.Getenv("SEC_TOKEN") && rbody.Token != ""
}
```

接到request後, 其實是可以把執行的task丟給另一個go routine處理, 原本的就可以回傳給ifttt, 避免執行太久而timeout的問題, 不過對heroku來說, 這還是在同一個web dyno上就是了