---
date: 2022-11-13T13:11:27+08:00
title: "使用Dapr Input/Output Binding連接Kafka"
slug: "Shi-Yong-Dapr-Input-Output-Bindinglian-Jie-Kafka"
images: 
- "https://og.jln.co/jlns1/5L2_55SoRGFwciBJbnB1dO-8j091dHB1dCBCaW5kaW5n6YCj5o6lS2Fma2E"
draft: false
---
在Dapr元件有一種叫做Binding的元件(component)讓你的app跟外部系統做一個連結的, 這元件可分為兩類:

1. Input BIndings: 用來接受外部事件的觸發,像是Webhook, 從Queue來的events, 甚至是人家新發的Tweets, 應該都可以歸為這一類
2. Ouput Bindings: 呼叫外部系統的動作,命名為Outpiut其實會讓人誤以為是資料的輸出,但其實,他不只可以用在資料輸出,呼叫外部系統的動作都可以包含在內, 舉個例子, [GraphQL 的Output binding](https://github.com/dapr/components-contrib/tree/master/bindings/graphql)定義了兩個操作(Operations), 一個是QueryOperation, 一個是MutationOperation, 熟悉GraphQL的應該知道,MutationOperation一般才是應用在資料操作,而Query感覺就跟輸出比較無關了

一開始我也有點搞不清楚這個模式目的在做啥的, 要接受事件觸發,我們有pub sub了,而state store本身就用在資料輸出, 感覺的確有點重複,但由上述兩點來看,其實Binding定義的範圍廣泛多了,它並不特定限制在Queue或是資料庫

但有個東西同時支援了Pub sub, input binding, output binding, 一開始我是看Kafka這應用,才讓我覺得有點錯亂,[前一篇](dapr-raw-payload-pub-sub)有講過了怎實作Subscriber, 這邊來比較一下, 利用Input binding的話,會有什麼不一樣?

### 建立Binding元件
![](/post/images/55875EA9-6C5C-48BB-8289-FF95B5CB5409.jpeg)

要用Kafka來觸發我們服務(如上圖),跟寫Subscriber一樣,我們需要在`~/.dapr/components`裡先建立好元件, 假設我們新增一個`kafka-binding.yaml`, 內容如下:

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: kafka-binding
spec:
  type: bindings.kafka
  version: v1
  metadata:
  - name: brokers
    value: localhost:9092
  - name: topics
    value: mytopic
  - name: consumerGroup
    value: group1
  - name: publishTopic
    value: mytopic
  - name: authRequired
    value: "false"
```

上面這個其實已經是定義好input和output binding了,topics定義的是input binding要聽取事件的topic, 而publishTopic定義的則是Output binding要輸出資料的目標

### 實作Input binding

跟實做subscriber差不多, Input binding也是實作一個webhook讓Dapr打進來而已, 這邊假設Kafka會收到的資料會是一個數字

```golang
package main

import (
    "log"
    "net/http"
    "github.com/gin-gonic/gin"
)

func dataBinding(ctx *gin.Context) {
    var data int
    if err := ctx.Bind(&data); err != nil {
        ctx.AbortWithStatus(500)
        return
    }
  
    log.Println(data)
    ctx.Status(200)
}

func main() {
    r := gin.Default()
    r.POST("/kafka-binding", dataBinding)
    r.OPTIONS("/kafka-binding", func(ctx *gin.Context) {
        ctx.Status(200)
    })

    http.ListenAndServe(":6003", r)
}
```

這邊幾個重點:
1. endpoint path跟你的binding名稱一樣, 當然這可以在元件設定那邊改
2. OPTIONS有點像是讓Dapr確認你有沒支援這個binding的health endpoint, 在程式一開始跑就會被call, 這邊只要回OK, 其實都好
3. 跟pub sub不一樣的是, 這邊會收到的格式不一定會是cloudevent, 除非publisher那邊過來的就是cloudevent, 因此, tracing應該是追蹤不到才是

### 實做Output binding
那如何實作跟剛剛的Input binding匹配的Output binding呢? 範例如下:

```golang
package main

import (
    "context"
    "log"
    "math/rand"
    "strconv"
    "time"
    
    dapr "github.com/dapr/go-sdk/client"
)

func main() {
    BINDING_NAME := "kafka-binding"
    BINDING_OPERATION := "create"

    for i := 0; i < 10; i++ {
        time.Sleep(5000)
        rand.Seed(time.Now().UnixMicro())
        dataId := rand.Intn(1000-1) + 1
        client, err := dapr.NewClient()
        if err != nil {
            panic(err)
        }

        defer client.Close()
        ctx := context.Background()

        in := &dapr.InvokeBindingRequest{Name: BINDING_NAME, Operation: BINDING_OPERATION, Data: []byte(strconv.Itoa(dataId))}

        client.InvokeOutputBinding(ctx, in)
        log.Println("Sending message: " + strconv.Itoa(dataId))
    }
}
```

這邊重點在於:
```golang
in := &dapr.InvokeBindingRequest{Name: BINDING_NAME, Operation: BINDING_OPERATION, Data: []byte(strconv.Itoa(dataId))}

client.InvokeOutputBinding(ctx, in)
```
雖說是"Output" binding, 但這邊用的名字是"Invoke", 跟Output沒啥相關, Operation則是元件訂的, Kafka binding只定義一個"create", 就是讓你送訊息用的, Data則是要傳送的資料, 以Byte array表示

### 用subscriber接收output binding來的事件

這邊就不用多作解釋了, 從前面不難發現, 它接的就是raw payload, 這部分可以參考 [前一篇](dapr-raw-payload-pub-sub) 

那啥時該用哪一種呢? 以Kafka這範例來說, 我是認為如果publisher跟subscriber都是自己實做的話, 應該是要選用pub sub, 用cloud events的話, 可以享受到distributed tracing帶來的好處, 如果不是, 差異應該不大, 都蠻簡單實作的