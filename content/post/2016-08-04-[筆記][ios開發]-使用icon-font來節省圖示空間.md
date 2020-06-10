---
date: "2016-08-04T21:48:57Z"
tags:
- iOS
- Swift
- mobiledev
title: '[iOS開發筆記] 使用Icon font來節省圖示空間'
---

寫App有一個讓人頭痛的是App大小的問題, 而這大小有部分是由App裡面所用的圖所貢獻, 為了減少這部份消耗掉的資源, 不管是用較大壓縮率的格式來壓縮圖檔, 還是其他, 大家都想盡辦法想解決這問題, 現在由於平面化的UI設計, 使得又有不錯的方法來解決這問題, 平面化UI設計的特色是大部分的圖檔都是單色而非五顏六色, 這使得用向量圖, 甚至用字型來解決這問題變得可行

#### 免費的圖標字型(icon font)

把所有的向量圖示變成字型檔可以節省不少空間, 以流行的[Font Awesome](http://fontawesome.io) 來說, 它包含了634個圖標, 卻只佔了153KB, 這在以往可能是不到十個圖標的檔案就會達到的大小, 相較之下節省了不少空間, 像這樣開放的圖示字型, 可以找到不少:

* [Font Awesome](http://fontawesome.io) : 蠻流行的一個開放icon組, 提供了ttf, woff等字型檔格式
* [Google material icons](https://design.google.com/icons/) : Google開放源碼的免費icon組, 它不只提供Android, iOS可使用的圖檔外, 也提供了字型檔的部分, 而且它的字型檔支援了Ligatures (後面會再提到它好用的地方),這也使得它比Font Awesome來的好用
* [Weather Icons](https://erikflowers.github.io/weather-icons/) : 顧名思義, 這提供了222個可以用於表示天氣的icons, 不過對於風向的表示的部分, 它是用同一個圖示只是在web上利用css旋轉來顯示不同方向的風, 這一點應用到App上的話, 我是還沒找到比較適合表達的方式
* [Octicons](https://octicons.github.com) : 由GitHub開源出來的圖標字型, 圖標不多, 但自己新增應該蠻方便的(自己增加svg檔用grunt去build) 

除了這些之外, 應該還可以找到不少免費的圖標字型(icon font)可以用, [IconFontKit](https://github.com/ElfSundae/IconFontKit)這邊就列了不少(它也整合了)可用的圖標字型

使用這些, 除了可以節省app的大小, 也可以省下不少設計圖標的時間, 但也不是沒缺點, 因為是字型的關係, 它每一個icon都是對應到一個unicode字元, 這字元大多數跟icon的形狀沒關係, 也就不是那麼好對應, 通常都要查一下對照表找出字元碼

#### 利用現成的framework整合
要在iOS上使用這些圖標字型(icon font)的方式好幾種,寫程式去load字型是一種, 當然就有不少大德, 寫好包裝可以讓你用cocoapods或是cathage直接引入, 這邊有幾個不錯的:

* [IconFontKit](https://github.com/ElfSundae/IconFontKit) : 這個整合應該是最完整的, 以Objective C寫的, 缺點是, 整合太多了, 反而變肥了
* [FontAwesome.swift](https://github.com/thii/FontAwesome.swift) : 這一個是針對FontAwesome, 還蠻輕量的, 以Swift時做的, 類似的有用Objective C寫的[FontAwesomeKit](https://github.com/PrideChung/FontAwesomeKit), 不過比較起來還是[FontAwesome.swift](https://github.com/thii/FontAwesome.swift)比較輕量, 如果用的只是FontAwesome, 那還是這個好
* [FontWeather.swift](https://github.com/julianshen/FontWeather.swift) : 老王賣瓜一下, 因為想寫的東西剛好有需要用到[Weather Icons](https://erikflowers.github.io/weather-icons/), 所以我就用了[FontAwesome.swift](https://github.com/thii/FontAwesome.swift) 的方式做了包裝

用這些現成的framework的好處是, 一來減去自己手動包裝字型進app的複雜度, 二來是, 這些已經幫你定義好一些對應圖標的常數, 讓你用比較方便的方式而不是記憶unicode字元來對應這些圖標

但它也是有缺點的, 大部分這些的作法都是runtime才去載入跟註冊字型, 因此你必須是在程式內設定你的UILabel, UIButton的字型, 無法事先就在Interface Builder做預覽, 所以個人比較喜歡的方式就是自己動手來

#### 手動在xcode上使用自訂字型

自己手動加的好處就是, Interface Builder上就可以套用, 直接就可以看到結果, 但就是稍微繁瑣了一點

![Interface Builder直接看結果](/images/posts/p1608041.png)

首先, 要把字型檔拖入你的Project裡面:

![ttf files in project](/images/posts/p1608042.png)

接著打開Info.plist, 加上一個新的項目叫做*"Fonts provided by application"*

這個是一個陣列(Array), 它的內容就是你要加入的字型檔檔名, 把你要加的每一個都列進去

![Info.plist](/images/posts/p1608043.png)

接著, 在Interface Builder裡你所要使用icon font的地方, 比如說UILabel設定你的字型, 原本的字型是設定為*"System"*, 把它改成*"Custom"*, 並選定你所需要的字型名稱, 例如FontAwesome, 要注意的是, 字型名稱不一定等同於你字型檔的名字:

![Interface builder](/images/posts/p1608044.png)

![Interface builder](/images/posts/p1608045.png)

接下來在**Text**的部分輸入這個圖示的代表的Unicode字元就好, 不是Unicode碼, 而是那個字元本身, 這挺不方便的, 可能用copy paste的才有辦法, 這是這個方法最大的缺點

這問題還是有方法克服, 這也就是前面為何提到會比較推薦使用[Google material icons](https://design.google.com/icons/)而不是[Font Awesome](http://fontawesome.io) , 這原因就是Ligatures

#### Ligatures

Ligatures是一個字型上蠻方便的特色的, 關於Ligatures可以先看一下這篇, 這是在[Google material icons](https://design.google.com/icons/)提到的一篇文章:

> [The Era of Symbol Fonts](http://alistapart.com/article/the-era-of-symbol-fonts)
    
剛剛提到的一個很大的缺點是, 你要知道圖示對應的Unicode碼才可以在你的UI上顯示你想要的圖示, 這相當不方便, 尤其那些Unicode碼可能根本完全不代表任何意義

比較人性點的作法是當你想要一個圖示代表藍芽, 用bluetooth就可以找到對應圖示, 而Ligatures就是一個這樣的存在

我們先來看看, 如果使用沒有而Ligatures的FontAwesome, 你在Text打上**"Contacts"**會是怎樣一個情形?

![Ligatures1](/images/posts/p1608046.png)

它會直接一字不漏的呈現**"Contacts"**,這還是因為FontAwesome有包含原本英數字字型在裡面, 有些其他的自行更慘, 根本就是一片白

讓我們再看看用[Google material icons](https://design.google.com/icons/)的字型,同樣的東西會有什麼結果

![Ligatures2](/images/posts/p1608047.png)

因為這個字型有支援Ligatures, 所以在這邊contacts就會被直接代換成它對應的圖示了, 我們就不用寄那種完全看不懂的unicode碼了

但大部分的字型其實也都沒有, 所以自訂字型該怎做?那就留待之後研究了