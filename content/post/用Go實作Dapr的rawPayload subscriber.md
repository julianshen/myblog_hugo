---
date: 2022-11-08T00:33:12+08:00
title: "用Go實作Dapr的rawPayload Subscriber"
images: 
- "https://og.jln.co/jlns1/55SoR2_lr6bkvZxEYXBy55qEcmF3UGF5bG9hZCBTdWJzY3JpYmVy"
draft: false
slug: dapr-raw-payload-pub-sub
aliases: 
- /用Go實作Dapr的rawPayload-Subscriber/
---

本來沒預計寫這篇的, 不過後來想想, 本來想寫的篇幅太大, 先寫這篇幫後面內容暖身, 後續相關內容會再更新到下面連結:

1. 待定

這篇並不是要寫怎用go實做Dapr的pubsub, 不完全是, 實做pubsub部分請參考[官方文件](https://docs.dapr.io/developing-applications/building-blocks/pubsub/howto-publish-subscribe/), 基本的Dapr的publisher跟subscriber是用所謂CloudEvent的格式在傳遞, 用CloudEvent的好處是, 由於CloudEvent會幫忙夾帶一些metadata, 因此也就可以實現分散式追蹤(Tracing)的功能, 但缺點就是無法支援一些原本寫好的legacy publisher或subscriber, 所幸Dapr的pubsub還是支援raw payload可以讓你自組你的訊息格式

在開始之前, 為了測試實做, 我這邊採用了Kafka, 但由於Dapr把實做封裝得不錯, 所以其實也不一定要用Kafka, 不過支援了Kraft之後的Kafka, 由於可以去掉對zoo keeper的依賴, 所以算蠻簡單裝的

### 安裝Kafka

使用docker跑Kafka, 應該是最簡單的方式, 只要執行

```shell
docker run -it --name kafka-zkless -p 9092:9092 -e LOG_DIR=/tmp/logs quay.io/strimzi/kafka:latest-kafka-2.8.1-amd64 /bin/sh -c 'export CLUSTER_ID=$(bin/kafka-storage.sh random-uuid) && bin/kafka-storage.sh format -t $CLUSTER_ID -c config/kraft/server.properties && bin/kafka-server-start.sh config/kraft/server.properties'
```

這樣Kafka就可以順利活起來了, 完全不需要跑zoo keeper...喔耶...

建議也可以順便跑一下[Kafka map](https://github.com/dushixiang/kafka-map), 這樣待會可以直接發event來測試

### 新增Kafka component

有了Kafka後, 我們需要在Dapr新增這個component 才可以讓Dapr應用程式使用, 在`~/.dape/components`底下加一個檔案(可叫做kafka.yaml),內容是:

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: kafka-pubsub
spec:
  type: pubsub.kafka
  version: v1
  metadata:
  - name: brokers # Required. Kafka broker connection setting
    value: "localhost:9092"
  - name: consumerGroup # Optional. Used for input bindings.
    value: "group1"
  - name: authType # Required.
    value: "none"
  - name: disableTls # Optional. Disable TLS. This is not safe for production!! You should read the `Mutual TLS` section for how to use TLS.
    value: "true"
```

內容就不多做解釋了, 官方文件會有更清楚的說明, 這邊先說明, 這新增上去後, 我們會多一個pubsub component叫做`kafka-pubsub`, 這名字寫程式會用到囉

### 寫個subscriber來接收事件(event)吧
[官方文件](https://docs.dapr.io/developing-applications/building-blocks/pubsub/pubsub-raw/)其實有寫如何寫一個接收raw payload的subscriber, 但不像其他文件一樣有多種語言範例, 只有Python, PHP兩種

![](/post/images/20221108224329.png)

但其實, 如上圖, Dapr用sidecar的作法, 簡化了寫pubsub的複雜度, 而且減低了對語言的依賴, 也不像是Istio是從系統的角度設計, 算是有點有趣的作法, 你不用理解pubsub, 也不用特別知道你是用Kafka, NATS, 或者是RabbitMQ, 寫法都一樣, 不囉嗦, 直接看code

```golang
package main

import (
	"encoding/base64"
	"encoding/json"
	"log"

	"github.com/gin-gonic/gin"
)

type Subscription struct {
	PubsubName string            `json:"pubsubname"`
	Topic      string            `json:"topic"`
	Route      string            `json:"route,omitempty"`
	Metadata   map[string]string `json:"metadata,omitempty"`
}

type Event struct {
	Topic string `json:"topic"`
	Data  string `json:"data_base64"`
}

var sub = &Subscription{
	PubsubName: "kafka-pubsub",
	Topic:      "myevent",
	Route:      "create",
	Metadata: map[string]string{
		"rawPayload": "true",
	},
}

func main() {
	r := gin.Default()
	r.GET("/dapr/subscribe", func(ctx *gin.Context) {
		ctx.JSON(200, []*Subscription{sub})
	})

	r.POST("/create", func(ctx *gin.Context) {
		var event Event
		err := json.NewDecoder(ctx.Request.Body).Decode(&event)

		if err == nil {
			decoded, _ := base64.RawStdEncoding.DecodeString(event.Data)
			log.Println(string(decoded))
		}

		ctx.JSON(200, map[string]bool{
			"success": true,
		})
	})

	r.Run(":6002")
}
```

咦, 這不像是在寫subscriber呀, 倒像是一個web service, 沒錯, 實際上的subscribe的部分被封裝在Dapr內了, Dapr等於收到Event後會打給我們的程式

那他怎知道要收到哪個queue哪個topic要打到哪個endpoint? 很簡單, 你只要有一個叫做`/dapr/subscribe`的endpoint, 在開始執行後, Dapr會自行打這endpoint了解你希望幫忙它收哪些event, 這邊我們希望收的 PubsubName (這邊是我們剛剛加的`kafka-pubsub`), 另外我們希望收`myevent`這個topic, 然後我們會希望收到event後打`/create`這個endpoint, 這有個好處, 你換成另一個完全不一樣的方案, 比如說Redis, 是不需要重新改code的

那我們在`/create`會收到甚麼呢?基本上就是包裝成CloudEvent的資料結構, 不對, 我們不是要收raw payload嗎?別急, 它只是收到後幫你包裝, 你的raw payload是被base64編碼好好地放在欄位`data_base64`中

這邊我特別沒用任何Dapr SDK, 然後也用gin來寫(Dapr的sdk裡用的是Gorilla),主要是為了展示, 這簡單到不用SDK呀(其實sdk也還沒支援raw payload subscriber相關的呀 XD)

### 執行

```sh
dapr run --app-id subs --app-port 6002 --dapr-http-port 3601 --dapr-grpc-port 60001 --log-level debug go run main.go
```

指令如上, 可以設定log level把debug訊息打開, 這邊有一點需要注意的, 這浪費我半天的青春, app port一定要設對, 我們程式內用`6002`那麼這邊的app port就要是`6002`, 不然Dapr不但會不知道要打事件給你, 連一開始的設定都拿不到(就是打`dapr/subscribe`)

### 測試

測試方式很簡單, 如果你剛剛有裝Kafka map, 去那個topic發送一個訊息(按Produce message), 看有沒收到一樣的訊息就可以了