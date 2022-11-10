---
date: 2022-11-10T13:57:03+08:00
title: "[Hugo] 讓文章的URL不要是中文字"
slug: "Hugo-Rang-Wen-Zhang-De-Urlbu-Yao-Shi-Zhong-Wen-Zi"
images: 
- "https://og.jln.co/jlns1/W0h1Z29dIOiuk-aWh-eroOeahFVSTOS4jeimgeaYr-S4reaWh-Wtlw"
draft: false
---

在發表[前一篇](dapr-raw-payload-pub-sub)時, [Evan](https://www.evanlin.com/) 跟我說, "能不能用slug把url改好看點?", 從開始寫blog來, 其實我也沒在意有沒人看, 所以SEO相關的也沒太在意, 但既然有人講了, 我就來弄一下吧

### Prama links的設定

首先先來看看pramalinks設定的部分:

```toml
[permalinks]
    post = "/:slug/"
```

這邊就是用來設定你URL長相的地方, 我的原文都放在post目錄內, 所以我一開始的設定就是以slug當URL沒錯, 那怎還是有中文呢?

### Archetypes, 初始文章的設定

當你用`hugo new filename`, 他會拿`themes/[theme_name]/archetypes/default.md`當範本來建立初始文章, 像我原本的設定是:

```yaml
---
date: {{ .Date }}
title: "{{ replace .Name "-" " " | title }}"
images: 
- "https://og.jln.co/jlns1/{{ replace .Name "-" " " | title | base64Encode | replaceRE "=+$" "" | replaceRE "\\+" "-" | replaceRE "/" "_"}}"
draft: true
---
```

slug是沒設定的, 不過它似乎應該會用title去算slug, 所以在Paramlinks那邊沒啥影響, 但其實你也可以加一行變成:

```yaml
---
date: {{ .Date }}
title: "{{ replace .Name "-" " " | title }}"
slug: "{{ anchorize .Name | title }}"
images: 
- "https://og.jln.co/jlns1/{{ replace .Name "-" " " | title | base64Encode | replaceRE "=+$" "" | replaceRE "\\+" "-" | replaceRE "/" "_"}}"
draft: true
---
```

anchorize是可以把"This is a cat"轉成"this-is-a-cat", 其實這兩段效果差不多, 問題在, 不支援中文, 因此像是"這是中文 Chinese", 其實翻成的是"這是中文-chinese", 如果再把這段放到url, unicode url並不是很好看

找了半天, 並不是很好用, 本來想寫個wrapper script來做新建文章好了, 然後自動加入比較好看的slug, 但轉念一想, 何不直接改Hugo?何不直接改Hugo?何不直接改Hugo? (回音持續)

### 支援中文的slug轉換的go module

推薦這個[github.com/gosimple/slug](https://github.com/gosimple/slug), 這不只支援中文, 很多語言都有!用法也很簡單:

```golang
package main

import (
	"fmt"
	"github.com/gosimple/slug"
)

func main() {
	someText := slug.Make("影師")
	fmt.Println(someText) // Will print: "ying-shi"
}
```

(以上範例來自它github)它很貼心的幫你把中文字轉成拼音了, 這用來做url感覺還蠻適合的呀!

### 修改Hugo

那我要加在Hugo哪裡呢?前面說到`anchorize`, 那其實仿效它做一個`slugify`不就好了, 那也容易知道要改哪, 抄`anchorize`就對了

要改的只有兩個檔:

- tpl/urls/init.go
- tpl/urls/urls.go

在`tpl/urls/urls.go`加入:
```golang
func (ns *Namespace) Slugify(s any) (string, error) {
	ss, err := cast.ToStringE(s)
	if err != nil {
		return "", nil
	}
	return slug.Make(ss), nil
}
```

然後在`tpl/urls/init.go`的`init()`內加入:

```golang
ns.AddMethodMapping(ctx.Slugify,
    []string{"slugify"},
    [][2]string{
        {`{{ "This is a title" | slugify }}`, `this-is-a-title`},
    },
)
```

重新建置安裝好後, 就可以把`Archetypes`改成:

```yaml
---
date: {{ .Date }}
title: "{{ replace .Name "-" " " | title }}"
slug: "{{ slugify .Name | title }}"
images: 
- "https://og.jln.co/jlns1/{{ replace .Name "-" " " | title | base64Encode | replaceRE "=+$" "" | replaceRE "\\+" "-" | replaceRE "/" "_"}}"
draft: true
---
```

搞定! 結果寫這些code五分鐘, 但卻花了五十分鐘找方法呀, 果然Open source還是自己改來用最快