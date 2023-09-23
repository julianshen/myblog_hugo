---
date: 2023-09-19T01:38:36+08:00
title: "用NATS來實現分散式微服務"
slug: "Yong-Natslai-Shi-Xian-Fen-San-Shi-Wei-Fu-Wu"
images: 
- "https://og.jln.co/jlns1/55SoTkFUU-S-huWvpuePvuWIhuaVo-W8j-W-ruacjeWLmQ"
draft: false
---
在近幾年, 微服務(Micro service)架構大部分的人應該不陌生了, 不管是面試, 實戰, 應該都已經聽到快爛了, 不過, 這篇來講講一個基於NATS的做法

首先, 先來了解一下[NATS到底是啥東西?](https://nats.io/about/)簡單來說, 它是一個輕量(Container image只有小小的18MB), 高效, 且安全的訊息佇列(Message Queue), 就基本的Pub/Sub用法來說, 它也的確像是這樣, 很容易就會把它跟Kafka, RabbitMQ等等歸為同一類, 那, 如果要談用Message Queue做微服務的溝通核心, 那有啥好講的? 不就是像是發佈訂閱(Pub Sub), 做成非同步架構, 那有啥好講的?

在微服務架構下, 要完成一件事, 各微服務之間的溝通是非常吃重的, 一般來說比較直覺的方式就是制定介面(API)來當作各微服務間溝通的協議, 微服務之間透過呼叫API的方式來與另一個服務做溝通, 不管是透過REST API或是透過gRPC, 這都屬於同步(Synchronized)的溝通方式, 也就是任一次呼叫在一定時間內都會預期有回覆(或錯誤)

再另一種方式就是利用Message Queue做成非同步的做法, 也就是呼叫方把訊息發佈到Message Queue內, 再由另一方訂閱方把訊息收去處理, 因為每次呼叫並不會需要預期有回應的結果, 呼叫方把訊息發佈後, 就不理了, 所以也就不會造成程式的阻塞, 適合需要處理很久的操作, 缺點就是呼叫方不容易拿到執行結果

如果只是要講後者, 那這篇講到這邊差不多就可以下課了(那我還寫幹嘛), 其實NATS的目標應該不僅止于Message queue, 由網站上寫的[有關NATS的相關內容](https://nats.io/about/), 可以知道它目標是作分散式應用程式的中樞神經系統, 所以其實除了非同步的方式外, 也可以識做成同步架構

## Request-Reply
前面有說到, 微服務間的溝通方式, 其中一種就是一個微服務透過API呼叫另一個微服務, 而這個API可以預期的狀況是: 
1. 成功並取得結果
2. 失敗並取得錯誤相關訊息
3. 在等待一段時間後(timeout), 呼叫失敗

NATS也提供機制讓你達成這樣的結果, 雖然NATS的基本就是Pub Sub, 但還是提供了Request/Reply的做法

不多囉唆, 先看一下程式：

```golang
    nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}
	defer nc.Close()
	subj, payload := args[0], []byte(args[1])

	msg, err := nc.Request(subj, payload, 2*time.Second)
	if err != nil {
		if nc.LastError() != nil {
			log.Fatalf("%v for request", nc.LastError())
		}
		log.Fatalf("%v for request", err)
	}
```
上頭這隻程式是一個 **"requester"**, 他把請求送到一個NATS subject, 並且等待並接收回傳訊息, 其實看起來就跟一個publisher沒啥兩樣, 差別就是他會卡在那邊等待回應(或timeout)

```golang
    //Responder
    nc, err := nats.Connect(*urls, opts...)
	if err != nil {
		log.Fatal(err)
	}

	subj, reply, i := args[0], args[1], 0

	nc.QueueSubscribe(subj, *queueName, func(msg *nats.Msg) {
		i++
		printMsg(msg, i)
		msg.Respond([]byte(reply))
	})
	nc.Flush()
```

上面則是相對於 **"requester"** 的 **"responder"** , 其實跟個subscriber差不多, 就是把訊息接回來處理,多一個回傳的動作(`msg.Respond([]byte(reply))`)而已, 從抽象角度來看, 跟我們直接拿REST API實作有點類似:
![](/images/posts/restapi.png)

但實際上, 他的做法比較是這樣的:

![](/images/posts/reqresp.png)

好像不太意外, 但這樣有啥好處, 我不就直接寫rest不就好了? 我們先來看一下負載平衡的做法好了:

![](/images/posts/lbrrr.png)

在這做法下, NATS其實就擔當起load balancer這角色了, 其實, 不知道你有沒注意到, 他也兼顧了service discovery的角色, 傳統你呼叫一個API service, 你必須先知道他的endpoint, 但在這邊你只要知道subject就好了, 因為responder是在監聽著那個subject, 因此, 還可以變形成這樣:

![](/images/posts/crosszoneee.png)

就可以簡單的實現到跨區呼叫或故障轉移(failover)

## NATS Service API

這應該是一個美麗(?)的未來, 不久前看到這段影片, 其實也真的就不久, 三月放出來的影片, 離現在也沒多久

{{< youtube byHGNUqIONw >}}

剛開始看到覺得, 頗酷的呀, 感覺就是在原本request/reply機制上再加上更多像是monitor和tracing的機制, 並讓它變得更像RPC call

但為了寫這篇時, 做實驗後發現, 他講的東西目前也都還沒push到main trunk去的樣子, 像是schema, 說有支援typescript也還沒, 還有`nats service`相關的指令也都還沒有, main裡面還沒找到相關的source code

所以這篇就沒打算寫太多了, 免得未來差異太大, 相關細節還是可以去看那段影片

先簡單來看一下程式會長成怎樣:

```golang
// GreeterServer is the server API for Greeter service.
type GreeterServer interface {

	// Sends a greeting
	SayHello(in *HelloRequest) *HelloReply
}

func RegisterGreeterServer(conn *nats_go.Conn, subject string, greeter GreeterServer) error {
	srv, err := micro.AddService(conn, micro.Config{
		Name:    "greeter",
		Version: "1.0.0",
	})
	if err != nil {
		return err
	}
	grp := srv.AddGroup(subject)
	grp.AddEndpoint("sayhello", micro.HandlerFunc(func(r micro.Request) {
		req := &HelloRequest{}
		proto.Unmarshal(r.Data(), req)
		resp := greeter.SayHello(req)
		data, _ := proto.Marshal(resp)
		r.Respond(data)
	}))
	return nil
}

type GreeterClient struct {
	subject string
	timeout time.Duration
	conn    *nats_go.Conn
}

func NewGreeterClient(conn *nats_go.Conn, subject string, timeout time.Duration) *GreeterClient {
	return &GreeterClient{subject, timeout, conn}
}

func (c *GreeterClient) SayHello(in *HelloRequest) (*HelloReply, error) {
	data, _ := proto.Marshal(in)
	msg, err := c.conn.Request(c.subject+".sayhello", data, c.timeout)
	if err != nil {
		return nil, err
	}
	reply := new(HelloReply)
	proto.Unmarshal(msg.Data, reply)
	return reply, nil
}
```

這段client/server的程式就跟request/reply的感覺差不多, 只是多了一些東西

其實我也試著想結合grpc跟這機制, 因此寫個[小工具叫NPC](https://github.com/julianshen/npc/), 所以以上的程式其實是由底下這個定義產生的:

```proto
syntax = "proto3";
option go_package = "nnrpc/pb";

// The greeting service definition.
// - version: 1.0.0
service Greeter {
  // Sends a greeting
  rpc SayHello (HelloRequest) returns (HelloReply) {}
}

// The request message containing the user's name.
message HelloRequest {
  string name = 1;
}

// The response message containing the greetings
message HelloReply {
  string message = 1;
}
```

(這邊就不談怎寫protoc的plugin了)