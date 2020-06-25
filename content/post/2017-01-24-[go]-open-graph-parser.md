---
date: "2017-01-24T01:38:19Z"
images:
- /images/posts/2017-01-24-[go]-open-graph-parser.md.jpg
tags:
- golang
- web
- opengraph
- html
title: '[Go] Open graph parser'
---

最近因為寫bot, 處理不少的HTML資料, 其中最常用的就是去取的[Open Graph](http://ogp.me/)的內容,
取這部分的資料是做啥用呢? 現今, 多數的網頁已會用[Open Graph](http://ogp.me/)和[Twitter Card](https://dev.twitter.com/cards/overview)
來描述網頁的一些屬性, 比如說標題, 相關圖片, 關聯影片等等, 而不管是[Open Graph](http://ogp.me/)還是[Twitter Card](https://dev.twitter.com/cards/overview)
都以HTML的meta tags存在的, 像這樣:

```html
<meta property="og:site_name" content="TechCrunch" />
<meta property="og:site" content="social.techcrunch.com" />
<meta property="og:title" content="Hugo Barra is leaving his position as head of international at Xiaomi after 3.5&nbsp;years" />
<meta property="og:description" content="Chinese smartphone maker Xiaomi is losing its head of international and primary English-language spokesperson Hugo Barra after he announced his exit from the.." />
<meta property="og:image" content="https://tctechcrunch2011.files.wordpress.com/2016/05/hugo-barra-2.jpg?w=764&amp;h=400&amp;crop=1" />
<meta property="og:url" content="http://social.techcrunch.com/2017/01/22/hugo-barra-is-leaving-his-position-as-head-of-international-at-xiaomi-after-3-5-years/" />
<meta property="og:type" content="article" />
<meta name='twitter:card' content='summary_large_image' />
<meta name='twitter:image:src' content='https://tctechcrunch2011.files.wordpress.com/2016/05/hugo-barra-2.jpg?w=764&#038;h=400&#038;crop=1' />
<meta name='twitter:site' content='@techcrunch' />
<meta name='twitter:url' content='https://techcrunch.com/2017/01/22/hugo-barra-is-leaving-his-position-as-head-of-international-at-xiaomi-after-3-5-years/' />
<meta name='twitter:description' content='Chinese smartphone maker Xiaomi is losing its head of international and primary English-language spokesperson Hugo Barra after he announced his exit from the company.' />
```

(範例剛好跟上時事 :P)

在做我的bot[新聞萬事通](https://line.me/R/ti/p/%40cur4648v)就是用這個來取得圖片跟標題的內容, 當然, 你在Facebook分享了一個連結, Facebook會自動帶上縮圖跟內容, 資料來源也是來自OG

## OG

雖然Go也可以找到一兩個Open graph的parser, 看起來好像還算堪用, 不過因為自己用得多了, 所以索性自己寫一個叫OG:

[GitHub/julianshen - OG](https://github.com/julianshen/og)

仿golang本身的json package, 也利用了Reflection, 因此還有點靈活度, 可以支援更多的資訊

### 安裝OG

```
go get -u github.com/julianshen/og
```

### 基本資料結構

```go
type PageInfo struct {
	Title    string `meta:"og:title"`
	Type     string `meta:"og:type"`
	Url      string `meta:"og:url"`
	Site     string `meta:"og:site"`
	SiteName string `meta:"og:site_name"`
	Images   []*OgImage
	Videos   []*OgVideo
	Audios   []*OgAudio
	Twitter  *TwitterCard
	Content  string
}
```

PageInfo是預設的資料結構, 前面有提到, 有彈性可以支援更多的資訊, 所以這個除了直接使用外也可當作一個參考用的定義,
主要就是採用了Go的struct field tags, 定義了一個叫做"`meta`"的自訂tag, 對照前面的html範例就可知道, 後面的值
是html meta tag裡面的property, 因此你可以自定義自己的資料結構, 然後仿這個規則加上tag即可

也可支援巢狀式, Arrays, Pointer

### GetPageInfo

GetPageInfo是以已經定義好的PageInfo這個struct為主, (希望)已經包含比較基本的og或twitter card tags了,
直接呼叫對應的API, 就可以取得這個url裡面的PageInfo的資料了

```go
urlStr := "https://techcrunch.com/2017/01/22/yahoo-hacking-sec/"
pageInfo, e := og.GetPageInfoFromUrl(urlStr)
```

另外還有一個`PageInfo.Content`, 這是網頁去掉廣告跟多餘的東西的純文字內容, 這就是`GetPageData`所沒有的了

### GetPageData

跟GetPageInfo不同, 這很適合用在自訂資料結構, 舉個例說, 如果我只想抓`og:image`的資料就像這樣:

```go
ogImage := og.OgImage{}
urlStr := "https://techcrunch.com/2017/01/22/yahoo-hacking-sec/"
og.GetPageDataFromUrl(urlStr, &ogImage)
```

## 簡單談談[Go reflection](https://blog.golang.org/laws-of-reflection)

這個package因為要仿json package使用struct field tags, 所以用了Go的[Reflection](https://blog.golang.org/laws-of-reflection)機制

老實說, Go的reflection沒有像Java的設計那麼好, 寫到後面還蠻容易昏頭的, 而主要有三種要先搞清楚

- [Type](https://golang.org/pkg/reflect/#Type)
- [Kind](https://golang.org/pkg/reflect/#Kind)
- [Value](https://golang.org/pkg/reflect/#Value)

以這個例子來說:

```go
ogImage := og.OgImage{}
pImage := &ogImage
```

ogImage的Type是`struct OgImage`, 而Value是`og.OgImage{}`, pImage的Type即是`*OgImage`,
而ogImage的Kind則是`struct`, pImage則是`ptr`(pointer)

比較難搞懂的就是什麼時候要用Type, 什麼時候要用Value, Kind又是什麼時候? (作業?!)

```go
type := reflect.TypeOf(ogImage)
value := reflect.Value(pImage)
```

用'TypeOf'可以取得變數的Type, 而'ValueOf‘則是值(value)的部分, 而value雖然代表變數的值, 但它是一個叫做Value的struct,並不是原本的資料型態
所以不要把它直接當參數呼叫函數用, 如果真有需要, 可以用`Value.Interface{}`轉成`interface{}`用

用`StructField.Tag.Lookup`則是可以查Field裡面的tag內容

所以這個og parser的原理就是走過所有fields, 只要有meta tags的話就拿去搜尋對應的html tags,
有的話再把html tag裡面content屬性填入值

最後小抱怨一下, [Open Graph](http://ogp.me/)和[Twitter Card](https://dev.twitter.com/cards/overview)
實在很不一致, 一個用property, 一個用name, 因此在[OG](https://github.com/julianshen/og)的做法是先用property, 如果再找不到用name去找