---
date: "2016-09-28T02:17:12Z"
images:
- /images/posts/2016-09-28-使用aws-lambda和github來提供中華職棒賽程資料.md.jpg
tags:
- Golang
- Github
- server
- AWS
- lambda
title: 使用AWS lambda和Github來提供中華職棒賽程資料
---

不知不覺的突然就多出了兩天颱風假, 這颱風實在很威, 乒乒乓乓的, 不過, 也沒做什麼, 時間就快過完了, 現在才想到, 還是來寫點什麼, 嚴格說來這些東西並不完全是颱風假時弄的, 只是拖得有點久

起因是, 之前(很久...追朔到去年)想寫個App, 需要用到中華職棒賽程的資料, 拖了很久一直沒真的去做, 斷斷續續的, 最近才先把資料這部分補齊, 首先需求是:

1. 當月之後的賽程資料, 但中華職棒並沒有API, 只有(很爛)的網頁, 因此資料勢必得從網頁去解析
2. 由Client app直接去解析html, 會比較麻煩(如果網站更新了, 就要更新App), 不是那麼可行
3. 不想花錢(或不想花太多錢)弄一個server, 更不用說還要考慮Scaling

而賽程表這樣的資料的特性則是:

1. 球季是3~10月
2. 資料內容除週一(休賽)外, 幾乎每天都會變, 但不會一兩個小時或幾分鐘就變一次
3. 變動的內容可能是:
	4. 比賽結束, 比數有更新
	5. 延賽或停賽
6. 一個月才幾十場比賽而已, 基本上不太需要有search或query的功能, 依據月份分類也就足夠了

因此我採用的做法是:

1. 利用AWS lambda定時解析中華職棒網站的資料
2. 資料以json格式存在github (使用Github api)
3. Client透過CDN去要這些json的raw content

### 賽程解析 ###

這部分我是用Go + Goquery來寫的, source code在這邊: [cpblschedule](https://github.com/julianshen/cpblschedule), 這code沒啥整理過, 光解析這堆亂七八糟的html就夠頭痛囉, 就沒啥整理

我做成了一個package, 因此要使用可用下列指令先安裝:

``` go get -u github.com/julianshen/cpblschedule ```

裡面也很簡單就一個function而已, 因此要使用可以參考:

```go
import "github.com/julianshen/cpblschedule"

func main() {
	matches, err := cpblschedule.ParseCPBLSchedule(year, month)
    ....
}
```

### AWS Lambda ###

這邊就不介紹這東西是什麼了, 網路上文章一大堆, 基本上他是AWS一個severless的解決方案(這算廣告詞吧), 會使用這個的原因有二:

1. 依我的用量應該是免費(事實證明, 其實還是要花點錢, 我忘了算網路傳輸的費用了, 不過這不多)
2. 可以用Cloud watch排程觸發

不過有個小問題

	**他不支援Go!!!!!**

而我上面那個解析的東東是go寫的, 那不就寫心酸的

所幸還有別的辦法,就是把程式編譯成執行檔, 然後用nodejs去包裝它, 不過這有點煩瑣, 所幸還有工具

[apex](https://github.com/apex/apex), 這是讓你更簡單的去建立lambda function的工具, 而且他正好也可以幫你簡單做好上面所說的包裝

安裝及使用就看文件吧, 不特別說了, 但要如何用go寫一個lambda function handle呢?以下是範例:

```go
import (
	"encoding/json"
	"github.com/apex/go-apex"
)
func main() {
	apex.HandleFunc(func(event json.RawMessage, ctx *apex.Context) (interface{}, error) {
		dosomething()
		return nil, nil
	})
}
```

#### Github as API source ####

資料既然變動不頻繁, 就用lambda定期產生然後把結果放到Github上就可了

Github API的Go的實做是Google放出來的[go-github](https://github.com/google/go-github/), 文件還蠻眼花撩亂的, 不過在這應用需要的API不多:

1. client.Repositories.GetContents - 取得內容
2. client.Repositories.CreateFile - 創立一個新檔
3. client.Repositories.UpdateFile - 更新某個檔

之所以需要1的原因是要確認檔案是不是已經在repository裡面了, 如果沒有就用create, 如果有就拿SHA hash去更新內容

GetContents會把檔案內容一併給抓回來, 這可以用來在更新檔案前先比較, 如果不比較, 就算沒更動, API也會新增一個新的commit, 為了避免不要太誇張, 還是先比較一下好了

那之後client怎樣存取這些資料? 找到檔案, 選取raw就可以知道raw的url了, client每次就抓這個URL就好, 但為了避免過量地request湧到github, 因此透過一個CDN來存取可能會好一點

這時候就可以用[RawGit](http://rawgit.com), 這邊透過MaxCDN, 讓你可以去存取Github上的raw content, 而你的檔案的網址會是像這樣:  

	https://cdn.rawgit.com/user/repo/tag/file

大致上就這樣