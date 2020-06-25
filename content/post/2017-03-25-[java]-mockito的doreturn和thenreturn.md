---
date: "2017-03-25T20:55:24Z"
images:
- /images/posts/2017-03-25-[java]-mockito的doreturn和thenreturn.md.jpg
tags:
- java
- mockito
- test
title: '[Java] Mockito的doReturn和thenReturn'
---

在做測試時, 利用假資料來做測試還算是一個蠻常被利用的技巧, 除了可以減少測試中的變動因子, 維持測試的scope的穩定度, 避免因為非程式本身造成的問題影響測試外, 還有就是在有跟別人API的對接的場合, 在還沒實際的API測試時, 一樣可以測試介面實作有沒問題

[Mockito](http://site.mockito.org/)在這方面(mocking)算是一個很有趣的東西, 前陣子在做公司的東西時, 拿了Mockito做了一些unit tests, 就想要寫這一篇, 不過又一直拖著沒寫了 :P

如果對[Mockito](http://site.mockito.org/)沒什麼接觸的話, 可以看一下他的[文件](http://site.mockito.org/), 並沒有很多很複雜的API, 初次看可能會覺得有點神奇, 不過, 這邊並不是要講它神奇的原理, 而是探討一下他的`thenReturn`和`doReturn`

在[Mockito](http://site.mockito.org/)中, 你要創立一個mock object是像這樣:

```java
// mock creation
List mockedList = mock(List.class);
```

上面這例子就會創造一個虛假的List object, 你可以說這個object擁有List的特性, 但他卻是個假貨!

針對mock object, 你可能會去做這樣一件事:

```java
when(mockedList.get(0)).thenReturn("first");
/// 或是:
doReturn("first").when(mockedList).get(0);
```

上面這兩句話(好吧, 是兩行程式, 但Mockito實在太口語化了)在這狀況下是代表同一件事, 照理說應該也會有一樣的結果, 那為何要有兩種寫法呢?

這邊要注意一點是, "Mock"這件事是針對object, 創建出來的object實例(instance), 並不是針對類別(Class), 因此當你使用`when(mockedList.get(0)).thenReturn("first")`指的是如果你呼叫了`mockedList.get(0)`會回傳"first", 但不代表你呼叫所有`List.get(0)`都是回傳"first", 利用了`mock`創建出來的實例, 真的是個假貨, 呼叫它任何的方法(method), 都不會真的呼叫到你真正在類別裡面定義的實作, 而是被導引到空殼去了

再回到兩種寫法的問題, 針對上述的mock object, 這兩種寫法都是沒問題的, 完全一模一樣, 但[Mockito](http://site.mockito.org/)裡面還有另一種形式的mock, 叫做`Spy`

```java
   List list = new LinkedList();
   List spy = spy(list);

   //optionally, you can stub out some methods:
   when(spy.size()).thenReturn(100);

   //using the spy calls real methods
   spy.add("one");
   spy.add("two");
```

mock是偽造了全部的東西, 或許就像是個天才雷普利, 但spy想做的只是偽造部分的內容, 以上面的例子來說, `spy.add`會呼叫到真正的方法(method) - `add`, 也就是這個list實際上會有"one"和"two"兩個東西, 但呼叫`size()`時回傳的會是100, 所以以下這例子會有問題

```java
   //optionally, you can stub out some methods:
   when(spy.size()).thenReturn(100);

   //using the spy calls real methods
   spy.add("one");
   spy.add("two");
   spy.get(spy.size() - 1);
```

如果一般正常狀況, 這個list的大小正好是加入的大小（這邊為2), 但因為我們偽造了`size()`讓他回傳了100, 這邊就會有問題了, 因為`get`實際上呼叫到的會是"real method"

其實回來看when的語法的話, 會發現本身就是蠻神奇的了, `when(obj.method1()).thenReturn(anyObject())`, 照一般Java語法來看, 為什麼把obj.method1()的執行結果帶給when當參數？為何when的回傳結果呼叫了thenReturn之後我們呼叫`obj.method1()`都會是回傳我們指定的回傳值?

想知道? 問香蕉....不是啦...這必須要知道Mockito的原理, 不過因為它是利用了很多很冷門的神奇技巧, 不在這邊探討範圍, 這邊只需要知道一件事, `when(obj.method1())`的確真的是呼叫了`obj.method1()`沒錯, 但對`obj = mock(MyClass.class)`來說, obj完全是沒作用的空殼, 所以呼叫了`obj.method1()`, 什麼事都不會發生的!

但對於Spy來說就不一樣了, Spy不是完全體的mock, 他是個代理, 最後還是會呼叫後面真正的method的, 因此如果碰到`when(obj.method1()).thenReturn(anyObject())`, `obj.method1()`絕對會起作用的, 但我們通常的目的就是想取代這個method的執行結果, 因此這樣會有副作用, 也就是當我們要假冒他前, 會先觸發一次, 就跟詐騙別人前先跑去通知警察一樣(什麼爛比喻), 這不是我們樂見的, 因此我們需要的就是`doReturn`

`doReturn("my result").when(obj).method1()`在這裡的when會再把obj mock一次, 以至於`method1()`暫時變成個空殼不會在這邊被觸發, 這樣就可以達到我們前面說的目的了, 但, 缺點是, when ... thenReturn , 因為when可以先知道方法的回傳值型態, 因此`thenReturn`裡面放的值只能是那型態, 所以寫錯了, comoile time會抓出來, 但反過來用的`doReturn`, 它接的是`Object`, 並無法在compile time做檢查, 所以如果有出錯, 要到執行時期才會發現