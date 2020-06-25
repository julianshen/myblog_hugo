---
date: "2017-01-01T01:56:29Z"
images:
- /images/posts/2017-01-01-[go]-利用goquery-+-otto來分析網頁.md.jpg
tags:
- golang
- goquery
- otto
title: '[Go] 利用goquery + otto來分析網頁'
---

2017的第一篇blog就先又來回的go跟爬蟲了

[goquery](https://github.com/PuerkitoBio/goquery)在Go是像jquery一樣的存在, 可以讓你用jquery一樣的selector來解析html檔案內容, 像這樣:

```go
doc.Find(".sidebar-reviews article .content-block")
```

只要有一些jquery的知識大概就不難上手, 但一般的網頁內容並不只有html這麼簡單而已, 有非常多的網頁是根本把資料給放在javascript, 在browser端再重組而成,
碰到這種, 光解析網頁原始檔是不夠的, 因為很多tag的內容都是之後被javascript所產生的

因此解析這類的網頁, 便可能需要去解析javascript的內容, 甚至是執行javascript, 在Go有一個好用的套件叫[Otto](https://github.com/robertkrimen/otto)可以來解決這麻煩

這邊以[Line Today](https://today.line.me/TW)的網頁當作範例來說明, 如果你打開Line today的網頁原始檔來看, 絕大部分都是javascript, 而且資料似乎都在categoryJson內(如下)

```
<script>
    var categoryJson = {
		"categories":[
			{"id":100003,"name":"今日頭條","source":0,"title":null,"thumbnail":null,"url":{"type":"FILE","hash":
		...
		}
		...
</script>
```

一個土法煉鋼的方式是找到categoryJson前後的{},把它substring出來, 這內容當Json來分析, 但這樣其實還蠻容易出錯的

另一個方式是先用goquery取得這個script內容, 然後用Otto去執行它後把值取出來, 像:

```go
url := "https://today.line.me/TW"
doc, err := goquery.NewDocument(url)

if err != nil {
	log.Println(err)
	return
}

s := doc.Find("script")
script := s.Nodes[1].FirstChild.Data

vm := otto.New()
vm.Run(script)

if value, err := vm.Get("categoryJson"); err == nil {
	goData, _ := value.Export()
	mapData := goData.([]map[string]interface{})
	...
}
```

otto很方便, 可以直接執行這個javascript, 還可以取裡面的值, 取出來的值其實是一個叫Value的struct, 如果要方便用的話, 就要用Export()來把資料轉成`interface{}`

但`interface{}`對object來說還是很不好取用, 如果確定你要的值是個javascript object, 那轉成`[]map[string]interface{}`會稍稍方便點

不過針對很複雜也很深的物件來說, `[]map[string]interface{}`也不是那麼的好用, 當你要存取的值要往下好多層時, 光轉型就要做到瘋掉, 多層次的map很麻煩的呀!

因此我的做法是, 多加一段javascript來減低取出的結構的複雜度, 像這樣:

```go
vm := otto.New()
vm.Run(script)
_, err = vm.Run(`
	var news = categoryJson.categories[0].templates[0].sections[0].articles;

	var articles = [];

	for(var n in news) {
		var article = {
			'title': news[n].title,
			'link': news[n].url.url,
			'publisher': news[n].publisher
		};
		articles.push(article);
	}
`)

if err != nil {
	fmt.Println(err)
}

if value, err := vm.Get("articles"); err == nil {
	goData, _ := value.Export()

	mapData := goData.([]map[string]interface{})
	todayNews = make([]TodayNews, 0)
	for _, data := range mapData {
		news := TodayNews{}
		news.Title = data["title"].(string)
		news.Link = data["link"].(string)
		news.ImageUrl = data["imageUrl"].(string)
		todayNews = append(todayNews, news)
	}
}
```

這樣一來就稍微簡單多了, 不過目前碰到的缺點是, 似乎有些ES6的語法, 它不太認得呀