---
date: 2021-08-02T00:52:04+08:00
images: 
- "https://og.jln.co/tt/jlns1/5Zyobm9kZWpz5L2/55SoIFR5cGVzY3JpcHQg5ZG85Y+rdGhyaWZ0IENsaWVudA"
title: "在nodejs使用typescript呼叫thrift client"
---

在[Apache Thrift](https://thrift.apache.org/)的官網上, 有提供了[如何在nodejs下呼叫Thrift client的範例](http://thrift.apache.org/tutorial/nodejs)

這邊這個範例其實針對的是javascript, 而其由Thrift idl產生javascript的指令是:
```sh
thrift -r --gen js:node tutorial.thrift
```

但如果我們想要用typescript來寫呢? 這邊產生的程式碼就沒有適用於typescript的封裝, 那這個thrift generator有沒支援typescript呢? 如果我們用 ``` thrift --help ``` 來看看它的說明:

```
  js (Javascript):
    jquery:          Generate jQuery compatible code.
    node:            Generate node.js compatible code.
    ts:              Generate TypeScript definition files.
    with_ns:         Create global namespace objects when using node.js
    es6:             Create ES6 code with Promises
    thrift_package_output_directory=<path>:
                     Generate episode file and use the <path> as prefix
    imports=<paths_to_modules>:
                     ':' separated list of paths of modules that has episode files in their root   
```

除了我們剛剛用的```js:node```外, 還有一個```js:ts```, 似乎好像是有支援, 但你如果直接用```thrift -r --gen js:ts tutorial.thrift``` ,它產生的typescript code是給browser用的, 並非給nodejs用的, 這怎回事?難道就不能兼顧嗎?其實可以, 如果你去看[這段程式碼](https://github.com/apache/thrift/blob/master/lib/nodets/Makefile.am), 就會發現答案是用```-gen js:node,ts```, 範例沒寫, help也沒寫清楚

假設我們有一個範例叫sample.thrift:

```thrift
service SampleService {
    string hello(1: i64 a, 2: i64 b)
    void hello2()
}
```
那我們用這個指令
```sh
thrift -s -gen js:node,ts sample.ts
```
那就會在```gen-nodejs```產生以下四個檔

* sample_types.d.ts
* sample_types.js
* SampleService.d.ts
* SampleService.js

那我們如何在我們程式裡面呼叫Thrift client呢?參考以下範例:

```typescript
import {createConnection, TFramedTransport, TBinaryProtocol, createClient, Connection} from "thrift";
import { Client } from "./gen-nodejs/SampleService";
import Int64 = require('node-int64');

const conn:Connection = createConnection("localhost", 8080, {
    transport : TFramedTransport,
    protocol : TBinaryProtocol
  });

const client:Client = createClient(Client, conn);
(async () => {
    console.log(await client.hello(new Int64(11), new Int64(34)));
    conn.end();
  })()
```

對照一下原本javascript版本:

```javascript
const thrift = require('thrift');

const SampleService = require('./gen-nodejs/SampleService');

var transport = thrift.TFramedTransport;
var protocol = thrift.TBinaryProtocol;

var connection = thrift.createConnection("localhost", 8080, {
    transport : transport,
    protocol : protocol
  });

var client = thrift.createClient(SampleService, connection);

client.hello(1, 2).then(resp => {
    console.log(resp);
}).fin(() => {
    connection.end();
});
```

相較之下, typescript的版本好像好讀一些