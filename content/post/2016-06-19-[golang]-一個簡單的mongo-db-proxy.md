---
date: "2016-06-19T13:17:04Z"
images:
- /images/posts/2016-06-19-[golang]-一個簡單的mongo-db-proxy.md.jpg
tags:
- golang
- mongodb
title: '[Golang] 一個簡單的Mongo db proxy'
---
之前被Parse搞的半死, 一直很好奇它的API到Mongodb的request之間到底是怎樣的對應

要弄清楚這個其實也不難, 把Mongodb的profiler全打開去看log就好了(```db.setProfilingLevel(2)```), 但這也是有缺點, profiler會寫到```system.profile```這個collection去, 而它是固定大小, 不能無限制的放, 再加上它還要多寫入這段, 多多少少影響效能

我需要的是一個從外部來觀察的工具, 不會影響到DB本身, 並且也可以將網路本身所花費的時間也包含進去, 所以想到的是在中間插一個proxy server

在現成的工具找到一個叫[MonoDB Proxy](https://github.com/christkv/mongodb-proxy)的工具, 這是用nodejs寫的, 勉強可以, 也證明了這個方法是可行的, 但這工具雖然有做到代理這部份, 但在log部分, 由於它並未解析bson, 所以詳細的內容並不好看, 所以就自己來寫一個

#### 功能需求
1. 支援[mongodb wire protocol](https://docs.mongodb.com/manual/reference/mongodb-wire-protocol/), 而不是只是單純的轉送資料
1. 印出request跟response內JSON的內容
1. 要能夠知道每個request所需要的時間(含網路)

#### 成品
最後寫出的的成品在這: [https://github.com/julianshen/mongoproxy](https://github.com/julianshen/mongoproxy)

整個還蠻簡單的:

1. wire.go 實作wire protocol
1. proxy.go 實作從client收資料並轉寫到server端
1. cmd/mp/main.go command line主程式的部分

#### 使用方法
這個工具是用Go寫的, 所以使用之前需要先安裝go

##### 安裝
```go get julianshen/mongoproxy/mp```

這個步驟做完後, 就可以把mp這個指令裝好了, 確定 $GOPATH/bin是在你路徑內, mp這個檔也是在那邊

##### 使用
``` mp --port=6001 --remote=mydb:27017 --response```

其中:

1. port是你這個proxy server的服務點
1. remote是遠端的mongodb (host:port)
1. Reponse是需不需要log回傳的部分

#### 實作Wire protocol

本來覺得Wire protocol會蠻複雜的, 結果, 其實是蠻簡單的

所有的wire protocol request都會有一個標準的表頭:

```c
struct MsgHeader {
    int32   messageLength; // total message size, including this
    int32   requestID;     // identifier for this message
    int32   responseTo;    // requestID from the original request
                           //   (used in responses from db)
    int32   opCode;        // request type - see table below
}
```

對應golang, 我定義成這樣:

```go
type MsgHeader struct {
	MessageLength int32 // total message size, including this
	RequestID     int32 // identifier for this message
	ResponseTo    int32 // requestID from the original request
	//   (used in responses from db)
	Opcode // request type - see table below
}
```

因為一開始就可以讀到整個訊息長度的, 所以就蠻好解析的, wire protocol的實作我是有參考了[dvara](https://github.com/facebookgo/dvara), 本來是有想拿它的code來改, 但看了一下發現它也沒完整實作wire protocol, 秉著自己也來了解一下這部份的想法, 就重頭自己刻了

跟[dvara](https://github.com/facebookgo/dvara)不同的地方是, 我用go的binary package來讀header而非自己刻一個, binary.Read的確是一個蠻好用的工具, 用底下的code就可以讀出header這個資料結構:

```go
h := MsgHeader{}
err := binary.Read(r, binary.LittleEndian, &h)
```

另外, 除header外, 各request的所帶的欄位各自不同, 這部份的作法就是定好各個所需的資料結構, 用reflection的方式來讀取各相關資料:

```go
v := reflect.ValueOf(req)
v = v.Elem()

// 根據資料結構內定義的每個欄位用相關的方法讀取
for i := 0; i < v.NumField(); i++ {
    f := v.Field(i)
    t := f.Type()

    if bytesRead == int(h.MessageLength) {
        break
    } else if bytesRead > int(h.MessageLength) {
        return nil, ErrorWrongLen
    }

    switch {
    case t == reflect.TypeOf((bson.D)(nil)):
        d, n, e := readDoc(bufferReader)
```

#### 解析bson

這部份就沒再重新造輪子了, 直接用golang著名的mongodb driver mgo裡的bson lib: https://godoc.org/gopkg.in/mgo.v2/bson , 這個bson lib已經寫的很不錯了, 直接拿來用即可

在這個package內, 泛用的bson資料結構有兩種: [```bson.M```](https://godoc.org/gopkg.in/mgo.v2/bson#M) 和 [```bson.D```](https://godoc.org/gopkg.in/mgo.v2/bson#D)

這兩個是不同用途的,仔細看一下M跟D的定義:

```go
type M map[string]interface{}
```
和

```go
type D []DocElem
```

如果你是要把解析出的資料用map來操作, M是蠻方便的, 一開始我也是依著之前我寫相關的東西的習慣用M, 不過這邊卻是不可以用M的, 這也是我碰到bug的地方

由於M解析出的是Map, 所以每個field的順序它並沒記住, 但偏偏在wire protocol裡, 尤其是 $cmd, 順序是重要的, 所以Unmarshal出的M再Marshal回去, 順序可能不是原本的順序了, 而這在這個proxy應用上, client寫什麼東西過來就要寫什麼到server才不至於出錯

#### 測試和Debug

我是用[Wireshark](https://www.wireshark.org)來驗證我的實作有沒問題, Wireshark預設會把到27017 port的資料解析成wire protocol的內容供閱讀, 當然也可以自己手動請它解析

#### 其他應用
這樣的proxy應該可以不只應用在debug, 像是Parse open source的[dvara](https://github.com/facebookgo/dvara), 它也是利用了proxy來做connection pooling, 應該也可以用在request routing和caching的應用上
