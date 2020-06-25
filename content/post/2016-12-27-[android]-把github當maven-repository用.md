---
date: "2016-12-27T15:39:59Z"
images:
- /images/posts/2016-12-27-[android]-把github當maven-repository用.md.jpg
tags:
- Android
- Github
- jitpack.io
- maven
- gradle
- java
title: '[Android] 把Github當Maven repository用'
---

自從Android導入gradle之後, 使用開放的第三方的程式庫就越來越方便了, 雖然方便, 但也不免會碰到這類的問題:

1. 想要的功能在master branch上更新了, 但卻遲遲不release以至於想用新的功能無法用
1. 程式碼已經沒在維護了, maven repository上一直都還是有問題的舊版本, 明知道怎麼修卻無法代他release到maven repository上去, PR又遲遲沒人理
1. 想加上自己的私有功能, 又不想包整包source codes到app裡面去
1. 想要開放自己做的程式庫卻覺得release到maven很麻煩

還好有[Jitpack](https://jitpack.io/)這東西, 剛剛就是碰到一個東西有問題, 想把它修掉直接用, 研究了一下

用法很簡單, 首先要先把`maven { url 'https://jitpack.io' }`加入到repositories裡面去

```groovy
allprojects {
    repositories {
        maven { url 'https://jitpack.io' }
    }
}
```

然後在你的`dependencies`裡面加入:

```groovy
compile 'com.github.User:Repo:Tag'
```

User就是你Github的user name, Repo是Repository的名稱, Tag是Git的tag名稱, 舉個例：https://github.com/CowBayStudio/wire-gradle-plugin ,
像這樣的Url, User name就是CowBayStudio, repo就是wire-gradle-plugin , 如果你的Git repo並沒有任何的Tag, 可以用Commit hash或是branch-SNAPSHOT(例如master-SNAPSHOT), 當然
自己加上tag會是比較好的做法, 比較好控管

當你去build你的app時, 在抓這個dependency時, Jitpack就會自動幫你把code從github抓下來build好,
當然是developer就免不了有bug, 導致build fail, 去 Jitpack網站把你Github的URL貼上去就可以找到build log了, 像是[這個](https://jitpack.io/com/github/CowBayStudio/wire-gradle-plugin/ver11/build.log)

我在弄我的東西時, 就發生了build fail的狀況, 而問題的原因是其中一個依賴的jar是JAVA 8 build出來的, 但jitpack用Java 7去build我的程式庫(文件說default應該是Java 8呀, 騙我!),
這時候可以加一個jitpack.yml, 內容如下: 

```yml
jdk:
  - oraclejdk8
```

透過jitpack.yml可以有很多客制可設定, 詳情就看一下[文件](https://jitpack.io/docs/BUILDING/#build-customization)吧!

所以碰到自己想動的程式庫, 可以直接fork出來改, 改出來的就可以直接這樣用了

那private repo呢?付錢給[Jitpack](https://jitpack.io/private#subscribe)就有啦! XD