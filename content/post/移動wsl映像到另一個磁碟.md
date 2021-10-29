---
date: 2021-10-30T00:35:14+08:00
title: "移動wsl映像到另一個磁碟"
images: 
- "https://og.jln.co/jlns1/56e75YuVd3Ns5pig5YOP5Yiw5Y-m5LiA5YCL56OB56Kf"
---

[WSL](https://docs.microsoft.com/zh-tw/windows/wsl/about) 當作開發環境固然方便好用的, 但也吃蠻大空間的, 最近在開發的東西, 需要跑一個postgresql, 裝不少資料, 吃掉我蠻多硬碟空間的, 偏偏我SSD就只有小小的512G (好吧, 的確寒酸到不像開發者的電腦), 一下子就吃滿滿了, 所以就必須要把這個給搬到我另一個比較大的磁碟救急(說是救急的意思是, 預期它會吃上1TB, 不過這又是會碰到一個問題, 之後再解了)

WSL新開的映像檔都是放在`C:`(畢竟不是謎片,不會自動住到D槽...咦?!), 要把它搬家的話, 需要先把它export出來, export的方法很簡單:

```
wsl --shutdown
wsl --export Ubuntu-20.04 d:\ubuntuback.tar
```

先shutdown是想保險一點, 所以也先把相關的視窗(像是Terminal, VSCode)都關一關, 如果映像檔很大, 這預期要做非常久, 像我這個有200G以上, 放下去, 基本上我就去看電視不管它了, 當然也要確保一下你目標硬碟夠大, 這邊`Ubuntu-20.04`是我要備份的目標, 如果不知道名字是啥可以用`wsl -l`查詢

export完之後, 接著就用:

```
wsl --import Ubuntu20dev e:\wsl\dev D:\ubuntuback.tar
```

這邊`Ubuntu20dev`是新的名字, 不要跟舊的重複了, import完後就可以用`wsl -d Ubuntu20dev` 登入進去玩了

不過登入後, 咦, 等一下, 怎麼會是用root? 之前舊的並不是呀! 要解決這個問題, 新增這個檔案`/etc/wsl.conf`, 裡面內容是

```
[user]
default=yourloginname
```

把wsl再shutdown之後再重新進來就不會是root了