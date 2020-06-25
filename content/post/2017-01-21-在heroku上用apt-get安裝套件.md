---
date: "2017-01-21T00:22:49Z"
images:
- /images/posts/2017-01-21-在heroku上用apt-get安裝套件.md.jpg
tags:
- server
- heroku
title: 在Heroku上用apt-get安裝套件
---

[Heroku](https://www.heroku.com)蠻好用的, 也用了好幾年了, 拿來做prototype真是方便, 不過慚愧的是, 我還沒付過錢給他(真惡劣)
最近chatbot玩得比較多, 不想花錢租server, 所以就比較頻繁的用它, 說到這裡, 照例, 先來廣告一下:

新的[Line叭寇(Barcode)小幫手](https://line.me/R/ti/p/%40dlk1367a) (轉含有條碼的圖片給它, 它會幫你解讀):

[![](http://qr-official.line.me/L/zMCfmfxLHk.png)](https://line.me/R/ti/p/%40dlk1367a)

好了, 回歸正題, Heroku 雖然是一個Paas的服務, 但它彈性非常大, 透過不同的buildpack也可以支援不同的語言跟框架,
不像Google的GAE, 支援的平台就比較有限

Heroku本身也是跑在Linux上, 因此, 如果你需要額外的套件, 其實也是沒問題的, 舉個例子, 雖然我這個叭寇小幫手是用Go來寫的
但卻會用到[ZBar](http://zbar.sourceforge.net/)這個讀取條碼的C程式庫, Heroku上當然沒裝, 所以建置Go時, 會因為找不到
程式庫而失敗

在Linux上, 我們可以用`apt-get`去安裝套件, 以zbar這例子是`apt-get install libzbar-dev`, 但在Heroku上又該怎麼裝呢?

還是要透過buildpack, 在command line下執行下面指令:

```
heroku buildpacks:add --index 1 https://github.com/heroku/heroku-buildpack-apt
```

apt這個build pack是放在官方的Github上, 因為我們希望該需要的, 軟體在一開始就把它準備好了

除了加build pack外, 還是不夠的, 你的套件還是沒被安裝, 因此繼續加一個叫`Aptfile`的檔案, 內容是你需要安裝的套件

當你push你的程式到heroku上後, 它就會根據`Aptfile`去安裝相對應的套件了

