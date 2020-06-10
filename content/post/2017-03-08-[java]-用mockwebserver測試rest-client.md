---
date: "2017-03-08T21:42:36Z"
tags:
- Java
- REST
- Http
title: '[Java] 用MockWebServer測試REST client'
---

應該來規定自己一週至少要寫一篇文章的, 這禮拜剛回歸工作的生活, 回歸了Java, 先從今天算起, 看多久能寫個一篇

這次來寫寫怎麼測試REST client, 測試最直覺的當然是讓Client直接連到Server, 但這樣變數比較多, 比如說網路斷了呀, Server掛掉了呀, 測試資料也不穩定(資料庫內的資料並不一定是固定的), 不太利於自動化測試, 如果只是要測試Client邏輯, 自然擺脫這些因素比較好, 餵假資料(設計好的資料)是比較好的選擇

但總不可能為了測試, 寫一個測試用的假server吧? 為了這樣的需求, [Okhttp](http://square.github.io/okhttp/)有提供一個叫做[MockWebServer](https://github.com/square/okhttp/tree/master/mockwebserver)的(Android當然也可以用), 這個就是為了這用途而出現的

要使用MockWebServer的話, 它並不直接包含在[Okhttp](http://square.github.io/okhttp/)的包裝內, 要另外含入:

Maven:

```xml
<dependency>
  <groupId>com.squareup.okhttp3</groupId>
  <artifactId>mockwebserver</artifactId>
  <version>(insert latest version)</version>
  <scope>test</scope>
</dependency>
```

Gradle:

```groovy
compile 'com.squareup.okhttp3:mockwebserver:(insert latest version)'
```

使用方法也很簡單, 以Retrofit來當例子:

```java
// 建立一個MockWebServer
MockWebServer server = new MockWebServer();
// 建立假的回應資料
server.enqueue(new MockResponse().setBody("{\"status\":\"ok\"}"));

Retrofit retrofit = new Retrofit.Builder()
        .baseUrl(server.url("/"))
        .addCallAdapterFactory(RxJavaCallAdapterFactory.create()) 
        .build(); 
service = retrofit.create(Service.class);

// 啟動server 
server.start();
service.subscribe(subscriber);
server.shutdown();
```

MockWebServer是一個貨真價實的http server, 所以有自己的Url, 藉由呼叫 `server.url("/")` 可以取得它的url, 使用MockResponse來回傳假資料(比如說是JSON), 另外可以藉由takeRequest來驗證client送的request是否正確

```java
RecordedRequest request1 = server.takeRequest();
assertEquals("/v1/chat/messages/", request1.getPath());
assertNotNull(request1.getHeader("Authorization"));
```

那如果要測試不止一個URL呢?那就可以利用Dispatcher, 如下:

```java
final Dispatcher dispatcher = new Dispatcher() {

    @Override
    public MockResponse dispatch(RecordedRequest request) throws InterruptedException {

        if (request.getPath().equals("/v1/login/auth/")){
            return new MockResponse().setResponseCode(200);
        } else if (request.getPath().equals("v1/check/version/")){
            return new MockResponse().setResponseCode(200).setBody("version=9");
        } else if (request.getPath().equals("/v1/profile/info")) {
            return new MockResponse().setResponseCode(200).setBody("{\\\"info\\\":{\\\"name\":\"Lucas Albuquerque\",\"age\":\"21\",\"gender\":\"male\"}}");
        }
        return new MockResponse().setResponseCode(404);
    }
};
server.setDispatcher(dispatcher);
```

在Unit test中要測試REST Client的話, MockWebServer應該算是蠻好用的一個工具