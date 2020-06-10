---
date: "2016-10-21T09:07:59Z"
tags:
- concept
- idea
- social network
title: 淺談Social Feed - 概念篇
---

如果要找出一個我過去幾年工作中比較具有代表性的東西, 想了一下, 應該就是social feed這東西了(寫這篇時, 想了一下該用啥名詞, 以往我會叫Timeline, 不過Social feed應該更為貼切一點),
趁現在才剛離職有些時間, 把這些東西整理一下, 主要還是以以前做過的東西的概念為主, 希望沒忘掉太多

這系列打算由三篇來構成, 除了這篇概念篇外, 另外還會有兩篇比較細節一點的內容, 分為client和server的部分(之後寫完會再更新鏈結):

1. [淺談Social Feed - 多服務彙整式的social feed (Client)](/淺談social-feed-多服務彙整式的social-feed-client)
1. [淺談Social Feed - 後端架構實作 (Server)](/淺談social-feed-後端架構實作-server/)

## 什麼是Social feed? ##

如果要簡單的來解釋, 應該可以說是一條依時間線排序的社群動態, 來看看下面這張Twitter的畫面:

![](/images/posts/p1021_twitter1.png)

這算是一個簡單的例子, 基本上就是把社群動態一條線的排下來(所以之前我也比較習慣的把它叫做Timeline, 不過這邊會改叫Social feed是因為考慮到也有不是依時間排序的),
我不太確定最早是誰採用這樣的設計, 找了一些早期社群網站(Frienster, Myspace), 最早都還沒有這樣的設計:

Friendster (2002 原來那時候有繁中版?!):

![](http://cbsnews2.cbsistatic.com/hub/i/r/2011/07/06/4a83463d-a643-11e2-a3f0-029118418759/resize/620x465/d16e61e54ecd40c300fa4a6d0e52bd37/friendster-AP070927023229.jpg)

Myspace (2003):

![](http://linapps.s3.amazonaws.com/linapps/photomojo/kxan.com/photos/2014/12/g15554-popular-releases-since-perry-took-office/287973-myspace-was-launched-in-july-2003-45b2f.jpg)

即使是Facebook也要到2008年才有這樣的雛形(看下面這段video還蠻有趣的, 我好像也是那段時間開始跟這東西結下了孽緣)

{{< youtube Nl7igMfeOvo >}}

再看看2006年的Twitter, 似乎就比較像是一個雛形的樣了, 不過那時似乎還像是一個粗劣的網站

![](http://robbiesblog.com/wp-content/uploads/2015/08/Twitter-2006-1024x350.png)

當然也不是所有的Social feed都是由上而下的, 另一個有名的變形就是創了河道的Plurk, Plurk在2008創立, 在台灣紅了一陣子, 但在智慧型手機的浪潮沒跟的很好, 到了手機上就很難發揮河道這樣的特色了

![](/images/posts/post1021_plurk1.png)

## 現今的設計 ##

現在這時代, 智慧型手機當道, mobile first曾經流行過一陣子, 大家放很多心力在手機上, 但Social feed這東西, 到了手機上, 各家變化就不太大了, 大致上都很類似, 來看幾個手機上的範例:

#### Facebook: ####

![](/images/posts/post1021_facebook1.png)

#### Twitter: ####

![](/images/posts/post1021_twitter2.png)

#### Google+: ####

![](/images/posts/post1021_googlep1.png)

#### Linkedin: ####

![](/images/posts/post1021_linkedin.png)

#### Instagram ####

![](/images/posts/post1021_instagram1.png)

#### 新浪微博 ####

![](/images/posts/post1021_weibao.png)

#### Pinterest ####

![](/images/posts/post1021_pinterest.png)

這邊可以見到的是Pinterest採取了一個跟其他人不同的呈現方式, 但, 大體上的構成還是跟大家都相似的

另外可以發覺的是, 從2006在PC Web到現在2016手機上, 內容變複雜了, 從純粹文字到多媒體內容, 我個人其實不是那麼愛這種轉變, 因為要接收的資訊變多了, 雖然畫面變豐富更多,
但另外帶來的一個缺點, 尤其是在手機上, 螢幕已經不夠付載一則動態的資訊量了

## 使用者介面構成與行為 ##

先以Facebook的介面當做例子來解釋(其他各家都大同小異啦):

![](/images/posts/post1021_facebook2.jpg)

大致上可以分做為三部分 - 作者資訊(藍色, 黃色區域), 內文資訊(綠色, 紅色, 咖啡色區域), 社群互動功能(紫色區域)

### 作者資訊 (藍色, 黃色區) ###

光字面上意思就已經表達完這部分了, 大致上都是放作者的圖像跟ID(或名字), 在Twitter因為還有轉推的動作, 所以還包含了轉貼人的資訊, 近年比較流行的設計是會用圓形的頭像(像上面例子的Instgram, 微博, Google+),
圓形的頭像大半在client app要顯示時處理掉就好, 只是一個圓形的mask, 就如同我以前寫過這篇:

[圓形大頭貼 - 使用Picasso的Transformation](/android-圓形大頭貼-使用picasso的transformation/)

### 發布時間 (綠色區域) ###

一般說來, 發布時間這部分, 不會直接用絕對時間(幾年幾月幾日幾時幾分), 而是用"3分鐘前", "4天前"這樣的絕對時間, 這樣的顯示方式似乎就已經是一個約定俗成的默契了

這種時間格式有個好處, 不用管時差問題, social feed的內容可能來自於世界各地不同的朋友, 每個人時區不同, 與其轉成當地時區的時間格式, 還不如以這種方式表示來得直接一點, 也不用管字串會不會有太長的問題

做這樣一個東西, 也是不用重新造輪子啦, 已經有了[moment.js](http://momentjs.com)這樣方便的東西可以用了, 當然他也是有被移植到Javascript以外的, 比如說[SwiftMoment](https://github.com/akosma/SwiftMoment)

### 內文以及相關資訊(紅, 橘, 咖啡區域) ###

早期的social feed, 內容大多只有文字, 就算有連結的轉貼, 也只是多一個hyper link, 整體上讀起來還是文字, 接下來圖片被帶入後, 就變成有圖文夾雜格式出現,
Facebook這種通用的social network服務, 內容種類較多, 因此就會夾雜不同格式的內容, 除了純文字內容外, 還有圖片, 影像, 甚至, 現在多了個直播, 而像是Instagram這類以影像為主的,
格式就較為統一, 不過, 基本上也只是不同內容的內文顯示格式略有不同外, 在後面的資料結構理應大同小異

後來可能人們(不見得是使用者吧)不再滿足於單調的內容, 尤其是在社群網路上分享文章鏈結(像是分享新聞的)越來越多, 一堆超連結看起來也醜,
後來就出現了Twitter Card 和 Open graph這類的東西:

1. [Open graph](https://developers.facebook.com/docs/sharing/webmasters)
1. [Twitter Card](https://dev.twitter.com/cards/overview)

這進一步讓你可以去定義你自己的網站, 而社群服務像是Facebook再把你的網站當作一個物件, 以物件的類別來決定怎麼去呈現這個鏈結, 在視覺上就再更加的豐富

不過不管內容有多少種, 差別真的就是呈現方式的多寡, 呈現方式也是有限的, 在Client顯示設計上是可以設計靈活點可以擴展, 不過倒也不用考慮到會有無窮無盡的形式

另外跟內文相關的資訊常見的還會有喜歡這個動態的數目(Facebook還有多種情緒表示), 回文數目, 分享數目(不見得每個服務都有), 使用者可以透過這些數字來了解到這篇貼文的熱門程度,
但這些數字, 其實在大部分的服務裡都只能單參考用, 數字未必準確, 這是因為一來很難及時地把某則貼文按讚的狀態更新給所有人(對Server的負擔大, Client實作也複雜),
或許在視覺上可以用一些比較相對的表示方式而非絕對數量的表示

另外有些服務會節錄幾則(通常最多三則)回文跟著文章下面一起顯示, 像Facebook網頁版, 但一樣, 它也是難於即時的更新

### 社群互動功能 (紫色區域) ###

一般常見就是“喜歡(like)", "回文(comment)", "分享(share)", 社群網路的精神主要還是在互動跟分享, 因此這幾個也差不多是最精簡也必備的了

### 更進階一點的內容 ###

通常還會有所謂的 hashtag和mention (這邊以Twitter用詞)

所謂的hashtag是由User自訂, 跟這則內文相關的關鍵詞, 以"#"開頭, 差不多也是個約定俗成, 最早應該早在IRC時代就有在使用了,
什麼?沒聽過IRC?沒關係, 知道從很古早時代就有了, 設計上通常會把 #hashtag 作成連結的形式, 點下去顯示相關的文章

而mention指的是"提到"某某人, 所以通常的形式都是"@"後面加User ID, 這也已經是一種約定俗成的方式了, 設計上也都會是一個點下去就到那個人的資訊頁面的連結

### 時間線回朔 ###

這部分我們以前都把它稱作"load more", 不過覺得這樣講好像很難知其所以然, 先看一下圖:

![](/images/posts/loadmore1.png)

一般來說, 從Server端抓回來的文章不會一次傳回所有的, 因為那會實在太多了, 尤其對重度使用者來說, 從開始使用以來到現在可能為數不少, 因此當我們把整social feed (或time line或說stream)往下一直拉時,
總是會見底的

在以往的UI設計上大多會放一條"touch to load more"之類的讓使用者再讀取舊一點的資料(所以以往我們都會把它叫做load more),
但這樣的缺點是使用者體驗不會太好, 通常看到這條後就跳掉不看了, 因此後來就流行做成上圖那樣, 快拉到最下面時就預先抓取, 來不及抓完, 使用者就會看到轉圈圈的進度

最好的體驗應該是讓使用者無縫接軌, 可以一直一直往下拉不用中斷, 但這邊就存在有調整的空間了, 太晚觸發的話, 使用者滑到最下面還是會有等待時間, 等待時間只要一長了, 常常就沒耐性跑了,
所以如果可以提早一點抓取, 是可以減少拉到最下面的等待時間, 但到底要提早多少? 太早也不是一件好事, 可能會導致client太過頻繁跟server索取資料, 但實際上又用不到那麼多,
以至於浪費了太多的網路傳輸量, 以及增加了server的負擔, 但使用者滑動的速率每個都不同, 所以這是一件不好拿捏的事, 可能要經過多次試驗才會有比較好的體驗

### 資料更新 ###

這比較會出現在有背景更新的場合, 如果每次使用者要看最新內容要自己觸發更新, 更新結束前他不能做任何動作, 那就沒這問題, 這問題主要出現在使用者在瀏覽過程中, 背景更新有了新資料進來,
輕微的話, 資料從最上頭插入導致他正在看的位置跑掉了, 嚴重的話, 可能整個刷新後, 內容都不同了, 這當然對使用者體驗很不好, 現在大部分的設計都會設計成非同步更新, 也就是就算是由使用者觸發,
更新時, 使用者還是有機會做動作, 更新時間如果太長了, 就容易發生這狀況

這部分的解法通常像下圖Twitter的做法:

![](/images/posts/IMG_9618.png)

不直接刷新頁面, 而是先提示使用者有新的貼文, 這樣的感覺就好多了

## 聚合式的social feed ##

所謂的聚合式的, 就是把一堆Social service的feeds全部串在一起, 因為多數的使用者擁有不止一個社群網路帳號, 把所有放在一起在瀏覽上就不用一個個網站跑或一個app跳過一個app

最早著名的有Friendfeed, 它早已經被Facebook買下不存在了, 不過它就是這樣一個概念:

![](http://d2chgkz0kdtxdm.cloudfront.net/wp-content/uploads/2009/04/friendfeed-after.png)

另外手機上還有HTC的[Friend stream](http://mobile.htc.com/learnmore/desires/eng/howtos/GUID-89A364D2-2043-46BD-9249-AA7BE00577A9.html),
不過這個血和汗做來的產品也不在了 orz 

![](http://mobile.htc.com/learnmore/chacha/cht/howtos/GUID-9619718F-A721-4BAF-A4B4-C202E2A0D2B1-web.png)
(好, 我拿chacha的畫面的確是故意的 :P)

另外還有一個叫做[HootSuite](https://hootsuite.com)的也是類似:

![](/images/posts/post1021_hootsuite.png)

不過HootSuite跟前兩者有所不同的地方是它並未把Social feed全部整合在一個時間線, 而只是並列顯示, 全部整合的難度稍高, 這就留待下一篇來說明了

## 小結 ##

整理了這麼多, 是後面兩篇的前置, 這邊的概念等於是設計一個social feeds的"需求", 後面兩篇會再用到這些概念
