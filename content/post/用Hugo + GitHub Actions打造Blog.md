---
date: 2020-06-28T12:32:14+08:00
title: "用Hugo + GitHub Actions打造Blog"
---

大概有三年沒寫Blog了, 最近覺得年紀大了, 變得越來越會忘東忘西的, 是有必要強迫自己寫一些東西了

每次荒廢了很久後重新執筆, 好像就會習慣換個系統, 並且做個紀錄, 像是之前 - [這篇](http://blog.jln.co/%E5%8F%88%E5%86%8D%E6%8A%8ABlog%E6%90%AC%E5%AE%B6%E4%BA%86/) XD 
這次...也不要例外好了, 雖然之前用 [Jekyll](https://jekyllrb.com/) + GitHub Page還算堪用, 但就是覺得它不是很快(是太慢了), 加上這麼久沒用了, 也忘了差不多了, 還是拋棄它吧

GitHub是不錯的免費空間, 所以還是沿用吧, 把 [Jekyll](https://jekyllrb.com/) 換掉應該就差不多了, 這樣的話找另外一套靜態網頁產生器就足夠了, 也不用特地架server, domain name也沿用, 那這次該換甚麼呢? 最近這幾年, 比較迷 Go, 所以沒多考慮, 打算就採用用Go寫的 [Hugo](https://gohugo.io/), 當然希望以前的內容還是可以承繼下來, 之前在[Jekyll](https://jekyllrb.com/)上的功能能延續自然就更好了, 廢話說太多了, 廢話說太多了, 先來看一下這次做的改變

## 安裝 [Hugo](https://gohugo.io/)

[Hugo](https://gohugo.io/)是一套用Go寫的Open source靜態網頁產生器, 特點就是快, 非常快, 使用方法也很簡單, 如果你在Mac底下, 也像我一樣用 [Homebrew](https://brew.sh/), 那只要執行下面的指令:

```shell
brew install hugo
```

不過最近我在家的工作環境比較常是Windows底下用[WSL(Windows subsystem for Linux)](https://docs.microsoft.com/en-us/windows/wsl/wsl2-index)下開發東西，這邊就不得不誇誇微軟了，有了ＷＳＬ後，我在Windows下工作，也跟我在Ｍａｃ上一樣順手，WSL說穿了就是一個Linux環境, 我用的是Ubuntu, 不過我是想在Linux下也是用 我是想在Linux下也是用 [Homebrew](https://brew.sh/) 來安裝, 不想用apt(習慣了), 這時候就可以用下面的指令來安裝 [Homebrew](https://brew.sh/)

```shell
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/master/install.sh)"
```

## 移植原本Jekyll上的內容

這是整個所有過程中最簡單的事了, 當你裝好Hugo後只要執行

```shell
hugo import jekyll jekyll_root_path target_path
```

這邊的jekyll_root_path是你原始jekyll內容的路徑, 它會自動幫你轉換所有的內容到新路徑, 又快又簡單

移植完內容後, 可以到新的目錄執行```hugo server```先來預覽結果...不過, 慢著~~看起來好像跟想像的不同!因為我們還沒套用主題跟做相關設定

## 套用原本的主題

雖然換了新系統, 不過我還是希望風格可以接近之前的, [Jekyll](https://jekyllrb.com/) 是可以套用主題的, 而我之前套用的主題是改自 [hpstr](https://github.com/mmistakes/jekyll-theme-hpstr) 這個主題的(也只是換一下標頭途而已)

還好Hugo上也是有人幫忙移植了 [Hugo版hpstr](https://themes.gohugo.io/hpstr-hugo-theme/)

Hugo的主題(themes)是被放在```themes```的目錄下, 所以新增一個新的theme很簡單, 只要把theme下載回來解開放在那目錄就好, 如果要試試套用的結果, 就執行

```shell
hugo server --theme hpstr
```

這邊hpstr是theme的名字, 不過到這邊為止, 都只是預覽而已, 真正產生網頁時都還不會套用theme

## 設定網站

剛剛有說過, 這時候如果去產生網頁, 並不會真的套用你所想要的theme, 必須先做好設定才可以

在剛剛從Jekyll移植過來的目錄裡面, 已經有產生了一個 ```config.yaml``` 的設定檔, 如果希望產生網頁時套用我們剛剛選用的hpstr, 那只要在 ```config.yaml``` 裡加上一行 ```theme: hpstr```就好

這個config.yaml長的就像這樣:

```yaml
baseURL: http://blog.jln.co
disablePathToLower: true
languageCode: en-us
title: Le murmure de Julian
theme: hpstr
```

這個設定檔格式是以yaml格式儲存的, 但它其實也支援了 TOML, JSON等其他的格式, 這是因為, 創作了Hugo的大神 [spf13](https://spf13.com/) 同時也是著名 Go 套件 [viper](https://github.com/spf13/viper)的作者, 相信很多寫Go的攻城獅們很多都使用過 [viper](https://github.com/spf13/viper) 來讀取設定檔吧

這邊因為某些原因, 我後來決定不用yaml格式來存我的設定, 而是改用TOML, 因此我把config.yaml, 轉換成 config.toml , 以下是我最近一個版本的設定檔:

```toml
baseURL= "https://blog.jln.co"
disablePathToLower= true
languageCode="zh-tw"
title= "Le murmure de Julian"
theme="hpstr"
googleAnalytics = "UA-79243751-1"
PygmentsCodeFences = true
Paginate = 5
hasCJKLanguage = true
enableEmoji = true

[params]
        subtitle = "朱隸安貓囈語錄"
        images = ["/images/avatar.png"]
        [params.author]
                name = "Julian Shen"
                avatar = "/images/avatar.png"
                bio = "Softward developer"
                github = "julianshen"
                email = "julianshen22@gmail.com"
                linkedin = "julianshen"
                instagram = "julianshen"
        [params.image]
                feature = "/images/bkg2.jpg"

[outputFormats.RSS]
    mediatype = "application/rss"
    baseName = "feed"

[permalinks]
    post = "/:slug/"
```

這邊有些設定的意義, 容後再說(我沒忘的話), 但基本上有這些:

* baseURL - 你網站的URL
* languageCode - 語系
* title - 標題
* theme - 主題, 前面說過了
* googleAnalytics - GA的Tracking ID
* hasCJKLanguage - 這個設定關係到算字數,閱讀時間的

(好像都講得差不多了)

## 維持原本的URL格式

我原本URL格式是長這樣:

``` http://blog.jln.co/筆記Vue.js-Slot的應用/ ```

但Hugo實際產生的格式是這樣

``` http://blog.jln.co/post/筆記Vue.js-Slot的應用/ ```

原本的網站其實已經有被搜尋引擎爬過了, 所以我並不想改變URL格式, 因此我在設定檔內加入

```toml
[permalinks]
    post = "/:slug/"
```

這邊就是為了設定文章URL的格式, 當然不只有這邊可以設定, 也不是只有 ```slug``` 這個變數可以用, 詳細方法可以參考 [這篇說明](https://gohugo.io/content-management/urls/)

## RSS的位置

Hugo產生的網站也是會包含RSS的連結, 但它預設是放在index.xml, 但在我舊有用Jekyll產生的網站, 其實是放在feed.xml, 所以我在 [IFTTT](https://ifttt.com/)的設定是feed.xml, 這邊我也是不想變動, 所以加入了以下設定

```toml
[outputFormats.RSS]
    mediatype = "application/rss"
    baseName = "feed"
```

這樣你就會發現RSS不再是放在```index.xml```了, 而是放在```feed.xml```

## Shortcodes

在Jekyll上有一個非常好用的東西, 比如說你在文章內加入 ```{% youtube FhoPTyMUgX0 %}``` 最後產生的網頁就會自動嵌入對應的Youtube影片, 如果是用 ``` {% gist julianshen/229f4ac32b3893816bd7636b96fe6f7d %} ``` , 那就會嵌入對應的gist程式碼

這如果在Hugo上沒有, 就頭痛了, 還好, 也是存在的, 它叫 [shortcodes](https://gohugo.io/content-management/shortcodes/)

雖然功能一樣, 但格式是不同的, 以剛剛兩個例子來說, 它就分別變成 (把%%改成{})

* ```{ {< youtube FhoPTyMUgX0 >} }```
* ```{ {< gist julianshen 229f4ac32b3893816bd7636b96fe6f7d >} }```

## 流程圖嵌入

之前我寫過一篇 "[替Jekyll的markdown加上簡易流程圖功能](https://blog.jln.co/Blog-%E6%9B%BFJekyll%E7%9A%84markdown%E5%8A%A0%E4%B8%8A%E7%B0%A1%E6%98%93%E6%B5%81%E7%A8%8B%E5%9C%96%E5%8A%9F%E8%83%BD/)", 採用的是[Jekyll-mermaid](https://github.com/jasonbellamy/jekyll-mermaid)

[Jekyll-mermaid](https://github.com/jasonbellamy/jekyll-mermaid)是透過[mermaid-js](https://github.com/mermaid-js/mermaid)讓我們可以很簡單的在文章內加入流程圖

但在Hugo,並沒有這樣的plugin, 所以必須要用別的方式達成

首先, 先到 themes底下你的主題(這邊是hpstr)的```layouts/partials```找看看有沒head.html, 在裡面加入一行

```html
<script async src="https://unpkg.com/mermaid@8.2.3/dist/mermaid.min.js"></script>
```

這是為了在每個網頁都可以載入mermaid.js, 再來就是在內容目錄下的```layouts/shortcodes```裡面建立一個 ```mermain.html``` (注意喔, 不是在theme底下那個```layouts```喔), 內容如下:

```html
<div class="mermaid">
  {{.Inner}}
</div>
```

這是建立一個新的shortcodes叫mermaid, 內容會轉化成一個div, 這div class是mermaid, mermaid.js會透過class找到這個div, 並將裡面內容轉成流程圖, 因此, 你就可以用這樣的shortcodes:

```golang
 { {<mermaid>}}
 graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
 { {</mermaid>}}
```
畫出一個這樣的圖

{{<mermaid>}}
graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
{{</mermaid>}}

把原本轉過來的文章都改成新的shortcodes就大功告成了

## [Open Graph](https://ogp.me/)

[Open Graph](https://ogp.me/)很重要, 它決定了你的文章被分享到社群網路上的樣子, 長得太不起眼, 沒人會注意, 一篇文章的OG, 差不多會長像這樣:

```html
<meta property="og:title" content="[Blog] 替Jekyll的markdown加上簡易流程圖功能">
<meta property="og:description" content="對一個developer的blog來說, 流程圖似乎是蠻需要的, 比較能夠清楚來解釋一些東西, 但每個東西都轉圖檔還蠻麻煩的, 下面介紹一個有用的J"><meta property="og:type" content="article">
<meta property="og:url" content="https://blog.jln.co/Blog-%E6%9B%BFJekyll%E7%9A%84markdown%E5%8A%A0%E4%B8%8A%E7%B0%A1%E6%98%93%E6%B5%81%E7%A8%8B%E5%9C%96%E5%8A%9F%E8%83%BD/">
<meta property="og:image" content="https://blog.jln.co/images/posts/2016-08-31-%5Bblog%5D-%E6%9B%BFjekyll%E7%9A%84markdown%E5%8A%A0%E4%B8%8A%E7%B0%A1%E6%98%93%E6%B5%81%E7%A8%8B%E5%9C%96%E5%8A%9F%E8%83%BD.md.jpg">
```

這邊, ```og:title```跟 ```og:image``` 很重要的, 沒有圖, 就非常不起眼了, 標題不吸引人也是很不起眼, 所以我們要讓每篇文章都有對應的圖

如果我們去看 Hugo 的 [Open graph template原始檔](https://github.com/gohugoio/hugo/blob/master/tpl/tplimpl/embedded/templates/opengraph.html)

```golang
{{ with $.Params.images }}{{ range first 6 . -}}
<meta property="og:image" content="{{ . | absURL }}" />
{{ end }}{{ else -}}
{{- $images := $.Resources.ByType "image" -}}
{{- $featured := $images.GetMatch "*feature*" -}}
{{- if not $featured }}{{ $featured = $images.GetMatch "{*cover*,*thumbnail*}" }}{{ end -}}
{{- with $featured -}}
<meta property="og:image" content="{{ $featured.Permalink }}"/>
{{ else -}}
{{- with $.Site.Params.images -}}
<meta property="og:image" content="{{ index . 0 | absURL }}"/>
{{ end }}{{ end }}{{ end }}
```

這邊就不解釋 Hugo的tempalte語法了, 直接講答案, 跟 ```og:image``` 相關的設定是:

* 文章Front matter內設定的圖片(images)
* 文章內名字含有"feature", "cover", "thumbnail" 的圖片
* 網站設定(config.toml)裡params.images的第一張 (等於就是用這張當預設圖了)

由此可知, 只要在config.toml 裡面放上這設定

```toml
[params]
        subtitle = "朱隸安貓囈語錄"
        images = ["/images/avatar.png"]
```

那在文章沒半張圖時, ```/images/avatar.png```就會是預設圖

如果我想自己在文章中指定圖片呢? 我們在新增文章後, 每篇文章都會有這樣的表頭, 叫做front matter, Hugo支援yaml, toml, json等各種front matter格式, 以下這個由 ```---```分隔表頭跟文章的, 用的就是yaml格式

```yaml
 ---
 date: "2017-01-21T00:22:49Z"
 images:
 - /images/posts/2017-01-21-在heroku上用apt-get安裝套件.md.jpg
 tags:
 - server
 - heroku
 title: 在Heroku上用apt-get安裝套件
 ---
```

這邊就用了```images```屬性來設定了這篇文章的 ```og:image```

## 優化社群分享

上面提到了設定Open graph相關內容的基本, 不過, 還有一個問題

Facebook會針對你提供的圖片大小來決定分享出去的版面設計, 可以[參考這邊](https://developers.facebook.com/docs/sharing/webmasters/images), 最大版面用的圖是1200x630, 不然分享出去的就小小一塊沒人看到了

另外一個問題是, 如果你沒在front matter設定圖片, 那OG就會使用網站預設圖, 或是文章中含有或是文章中含有"feature", "cover", "thumbnail", 而這邊front matter上放的圖片還是自己放的, 如果不小心忘了, 絕大多數是忘了, 那可能每篇文章用的都是一樣的圖, 比較單調, 更慘的是, 萬一設的預設圖太小(像我一樣), 或是根本沒設, 那分享出去的版面就會很不顯眼

雖說會去掃文章中的圖片, 找出含有"feature", "cover", "thumbnail" 的圖片, 跟 [Jekyll Auto image](https://github.com/merlos/jekyll-auto-image) 有點類似(不過Auto image不會限制檔名), 當然你也可以把內建的open graph template拿出來改, 放到自己主題layouts下

如果你要改用自己的, 該怎做呢? 首先就是把內建的opengraph.html和twittercard.html拷貝一份到主題的layout/partials目錄下改, 並且也是去改layouts/partials/head.html, 你會在head.html發現這樣的內容(以hpstr為例):

```golang
<!-- Open Graph and Twitter cards -->
{{ template "_internal/opengraph.html" . }}
{{ template "_internal/twitter_cards.html" . }}
```

把這兩個路徑改成自己的就好了

還不滿足!!因為還是有些問題:

* 文章內的圖片可能本來就不大, 小於1200x630
* 文章內的圖片是外連的, 有可能之後消失不見
* 文章內根本沒有圖片

所以我希望:

* 被用在 ```og:image```的圖片要被放大裁切到1200x630的比例
* 這個圖片要放本地不能外連
* 沒有圖的狀況.....產生一個給它用....用標題文字來當圖, 大小也要1200x630

在這需求之下, 我自己就寫了一個工具叫 ```ogp```, 原始碼如下:

{{<gist julianshen f129b6db74c1dc0a93647acd6f9e0be1>}}

這工具做的是

{{<mermaid>}}
graph LR
A["分離Front matter跟文章內容"] --> I{是否含images屬性}
I -->|有| H
I -->|無| B["掃描文章內含的圖片"]
B --> C{是否有圖}
C -->|有| D[將圖片縮放並裁減]
C -->|無| E[產生文字圖片]
D --> F[在Front matter插入images屬性]
E --> F
F --> G[寫回原文]
G --> H[結束]
{{</mermaid>}}

這邊套用了幾個套件來處理, 

* Hugo本身的 Page parser, 用來處理Front matter跟Content分離, 由於Hugo本身就是用Go寫的, 所以很輕易地可以用 ``` import "github.com/gohugoio/hugo/parser/pageparser" ```來使用
* "gopkg.in/yaml.v2", 用來輸出yaml的
* "github.com/yuin/goldmark", 這是Hugo用的Mark down parser, 我這邊用來找出所有的圖片, 如果是外站的就下載回來
* "github.com/h2non/bimg", 這是一個基於libvips的圖片處理套件, 安裝上會需要libvips, 這邊就不詳述, 以後如果有時間再另外寫一篇好了
* "github.com/Iwark/text2img", 這有點有趣, 我本來想自己寫一個把標題轉成圖片的來當作```og:image```的材料的, 沒想到真有人已經做了, 這邊我是自己fork回來做了些小修正

## 自動化處理文章的```og:img```

有了上面的```ogp```後, 就可以"手動"來產生文章的```og:image```了, 當然, 這完全不方便, 我每次寫完一篇新文章後, 我就得自己手動跑一次 ```ogp```, 身為一個懶惰的攻城獅(做了前面這麼多還懶惰?), 當然要想想怎麼來自動化囉!!

這邊想到的方式是利用 [git-hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks), 甚麼是[git-hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks)呢?簡單的說, 就是用來在你做git操作時, 可以讓你觸發某些動作, 例如我這邊要用到的pre-commit和post-commit, pre-commit是在你下commit命令後但還沒真正做commit動作前觸發, post-commit則是在動作發生之後

要使用[git-hooks](https://git-scm.com/book/en/v2/Customizing-Git-Git-Hooks)的話, 要把你要執行的寫成腳本(script)放到 ```.git/hooks/``` 目錄下, 例如, pre-commit的腳本的檔名就是 ```.git/hooks/pre-commit```, 記得要把這檔案用 ```chmod +x``` 把權限改成可執行, 下面就是我用的pre-commit

```bash
#!/bin/bash
echo "Run precommit"
newfiles=`git diff --cached --name-status | awk -vFPAT='([^]*)|("[^]+)' '$1!="D" { print $2 }'`

for n in $newfiles
do
   if [[ $n =~ content/post/.*\.md$ ]]; then
       f=`basename $n`
       docker run -v `pwd`:/blog julianshen/ogpp $f > /tmp/$f
       mv /tmp/$f content/post/$f
       rm /tmp/$f
       touch .commit
   fi
done
exit 0
```

這邊所做的動作是找到這次有變更的文章(副檔名.md), 然後針對每個檔去跑ogp產生 ```og:image``` (我這邊把ogp包成docker image方便使用)

前面我說到, 會利用到 ```pre-commit``` 跟 ```post-commit```, 那 ```post-commit``` 又是用到哪個地方呢? 在 ```pre-commit```這邊, 我用了 ```ogp``` 有產生了新檔案(圖檔), 但對於git來說, 這個檔案並不在stage中(```git add```後), 所以當做完commit後, 這個檔案並不會被放進去, 因此我們得找一個方式把它一起放進去

如果你有發現到在```pre-commit```這腳本中, 在跑完 ```ogp```後有加了一個 ```touch .commit```的動作, 這目的就是為了來解決前面所說的問題, 因為有產生新的檔案, 所以這邊建立一個 ```.commit```的檔案來標記一下有新檔案產生, 這檔案本身沒太大意義, 而是為了在後面 ```post-commit```使用

那, 我的 ```post-commit```的腳本就會長成這樣:

```bash
#!/bin/bash
echo "Post commit"
if [ -e .commit ]
    then
    echo "Add rest files"
    rm .commit
    git add .
    git commit --amend -C HEAD --no-verify
fi
exit 0
```

很簡單, 只要有 ```.commit```存在, 就會去多做一次```git add```和```git commit```把新的檔案放進去

要注意, 這兩個hooks都只在本地端作用, 如果要讓它其他電腦也會有作用, 要記得複製過去才會有用

## 產生網頁並佈署

寫完文章後, 由於原始文章是 mark down 語法寫的, 如果沒把它轉成網頁是沒辦法在瀏覽器上看的, 產生網頁很簡單, 執行

```bash
hugo --minify
```

這就可以了, 當然你如果要簡單的執行 ```hugo```不加任何參數也行, 產生的靜態網頁, 就會放到 ```public``` 這個目錄下

所以把 ```public```目錄裡面的所有內容, 都放到你github page的repository下就好了!

慢著!!!慢著!!! 都手動嗎???

懶惰的攻城獅會同意嗎??? 沒關係, 接下來我介紹一下怎把它自動化

## 使用[GitHub Actions](https://docs.github.com/en/actions)來自動化流程

[GitHub Actions](https://docs.github.com/en/actions)是GitHub的持續集成(Continue Integration, 簡稱CI)的服務, 可能有不少朋友已經知道CI是甚麼了, 這邊不多做介紹, 實際上使用過GitHub Actions後會發現, 它簡單而且強大

由於我這個blog的內容都是放在GitHub page, 而且原始檔也是放在GitHub, 所以使用[GitHub Actions](https://docs.github.com/en/actions)來自動化建置頁面跟佈署看來也很理所當然

[GitHub Actions](https://docs.github.com/en/actions)的工作流程(workflow)設定檔都放在```.github/workflows```底下, 所以只要新增一個就可以使用了, 這邊要特別注意一下, 如果你用GitHub tokens來存取這目錄, 需要特別有```workflow```權限, 通常如果你clone你的repository時是用http的話, 那可能就無法更動到這個檔(無法push), 建議是使用ssh

我發現, 寫到這邊, 篇幅已經非常的長了, 如果要仔細介紹一下這工具的話, 那可能還要很長的篇幅來說, 所以這邊就不仔細介紹了, 底下分享一下我用的設定:

```yaml
# This is a basic workflow to help you get started with Actions

name: CI

# Controls when the action will run. Triggers the workflow on push or pull request
# events but only for the master branch
on:
  push:
    paths: ["content/**", ".github/workflows/main.yml", "config.toml"]

# A workflow run is made up of one or more jobs that can run sequentially or in parallel
jobs:
  # This workflow contains a single job called "build"
  build:
    # The type of runner that the job will run on
    runs-on: ubuntu-latest

    # Steps represent a sequence of tasks that will be executed as part of the job
    steps:
    # Checks-out your repository under $GITHUB_WORKSPACE, so your job can access it
    - uses: actions/checkout@v2

    - name: Hugo setup
      uses: peaceiris/actions-hugo@v2.4.12
      with:
        # The Hugo version to download (if necessary) and use. Example: 0.58.2
        hugo-version: latest
        # Download (if necessary) and use Hugo extended version. Example: true
        extended: false
    
    - name: Build pages
      run: hugo --minify

    - name: Deploy
      uses: peaceiris/actions-gh-pages@v3
      with:
        personal_token: ${{ secrets.PERSONAL_TOKEN }}
        external_repository: julianshen/julianshen.github.com
        publish_branch: master
        publish_dir: ./public
        cname: blog.jln.co
```

如果你用過另一個CI工具 [drone](https://drone.io/) 的話, 你會發現, 對比起來, [GitHub Actions](https://docs.github.com/en/actions)的彈性更大, 更方便, 觸發建置的條件不一定要是某個git事件(一般是push, pr), 檔案的更動也可以觸發,檔案的更動也可以觸發, 以這段為例:

```yaml
on:
  push:
    paths: ["content/**", ".github/workflows/main.yml", "config.toml"]
```

這就是指 content目錄下, 以及workflow或Hugo設定有改變後, 就會觸發

至於Hugo的建置跟發布, 就在請大家自己看設定檔吧, 其實很簡單的

## 總結

好就沒寫blog, 也好久沒寫長文了, 我發現我廢話一樣多, 不過這邊也花了很長時間弄, 所以也多了一點

還有一些沒寫上, 例如說我現在是用vscode打這篇文章, 搭配hugoify這個plugin, 並且同時用瀏覽器預覽, 雖然整套工具很geek, 不過弄好後還蠻好玩的就是了 

