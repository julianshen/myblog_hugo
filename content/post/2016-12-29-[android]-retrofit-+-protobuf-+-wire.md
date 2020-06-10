---
date: "2016-12-29T15:33:25Z"
tags:
- Retrofit
- Android
- Protobuf
title: '[Android] Retrofit + Protobuf + Wire'
---

[Retrofit](https://square.github.io/retrofit/)目前已然成為最熱門的呼叫REST API的開源程式庫了,不過,
大部分的Retrofit多是用來處理Json類的REST API, 但Retrofit的能力卻不僅限於此

透過Converter, Retrofit可以處理的不只是Json, 還包含了XML, [Protobuf](https://developers.google.com/protocol-buffers/), 甚至你自己的自訂格式, 剛好最近看Json不太順眼,
想來試試[Protobuf](https://developers.google.com/protocol-buffers/), 就來試試這部分

[Retrofit](https://square.github.io/retrofit/)官方網頁上列的[Protobuf](https://developers.google.com/protocol-buffers/) converter有兩種,
一種是使用Google親生的[Protobuf](https://developers.google.com/protocol-buffers/), 這個只要引入```com.squareup.retrofit2:converter-protobuf```即可使用,
另一種是有點比較吸引我的是Square自己開發出來的[Wire](https://github.com/square/wire), Sqaure真的是蠻喜歡打造自己的東西的, 打造出來的又比別人威,
這點讓我對[Wire](https://github.com/square/wire)的興趣比較大, 因此這篇主要是以[Wire](https://github.com/square/wire)當範例來介紹

使用Wire converter其實相當簡單的:

在build.gradle內加入相關的dependencies, 包含了wire runtime跟retorfit的converter:

```groovy
compile 'com.squareup.retrofit2:converter-wire:2.1.0'
compile 'com.squareup.wire:wire-runtime:2.2.0'
```

在build service時, 用`addConverterFactory`將`WireConverterFactory`加入即可

```java
Retrofit retrofit = new Retrofit.Builder()
		.baseUrl(baseUrl)
		.addConverterFactory(WireConverterFactory.create())
		.addCallAdapterFactory(RxJava2CallAdapterFactory.create())
		.build();
DataService service = retrofit.create(DataService.class);
```

DataService的定義如下(這範例搭配了rxjava2):

```java
public interface DataService {
    @GET("data")
    Single<Posts> getPosts();
}
```

這邊有個問題, `Posts`這個class並不是直接用Java刻出來的, 而是由`.proto`檔產生的, 內容如下:

```protobuf
syntax = "proto3";
package wtf.cowbay.dp;

message FBPost {
	string id = 1;
	string title = 2;
	string message = 3;
	string imgsrc = 4;
	string target = 5;
	string createdtime = 6;
}

message Posts {
	repeated FBPost data = 1;
}
```

必須要用工具產生Java class才能使用, Google的Protobuf有自己的工具, 而Wire也有自己的, 如果手動自己執行工具產生檔案後再加入, 未免太鳥,

還好Wire有[wire-gradle-plugin](https://github.com/square/wire-gradle-plugin), 不過這個有個大問題, 雖然放在Square的repo下, 但似乎
不是官方版本, 而是有人貢獻的, 因此並沒有跟上最新的2.x的版本, [JakeWharton大神說將會有官方版本](https://github.com/square/wire-gradle-plugin/issues/11),
但似乎從六月到現在都沒出現, 所以只好[自力救濟自己改 - 我的版本](https://github.com/CowBayStudio/wire-gradle-plugin)

使用這個plugin很簡單, 先在第一層的build.gradle裡的buildscript加入:

```groovy
buildscript {
    repositories {
        jcenter()
        maven {
            url "https://jitpack.io"
        }
    }
    dependencies {
        classpath 'com.android.tools.build:gradle:2.2.3'

        classpath 'com.github.CowBayStudio:wire-gradle-plugin:ver12'
    }
}
```

主要是加入jitpack.io, 並把我的版本的plugin放到dependencies去

接下來在app的build.gradle裡加入:

```groovy
apply plugin: 'com.squareup.wire'
```

`.proto`檔案放置的位子在`app/src/main/proto`, 把檔案放在這邊, build的時候就會自動產生這些對應的Java classes了