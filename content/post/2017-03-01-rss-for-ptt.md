---
date: "2017-03-01T17:08:32Z"
tags:
- golang
- Ptt
- bbs
- rss
title: RSS for Ptt
---

最近才發現, ptt的rss功能好像拿掉了, 這樣好像就不能拿feedly之類的來訂閱版面內容, 反正我自己有寫了一個[gopttcrawler](https://github.com/julianshen/gopttcrawler)
所幸自己來寫一個吧!

source code在: [pttrss](https://github.com/julianshen/pttrss)

可以自行deploy到heroku去, 如果不想這麼麻煩, 可以用:

https://ptt.cowbay.wtf/rss/版名

例如: 
 
- 表特版的rss url是 - https://ptt.cowbay.wtf/rss/Beauty
- 電影版的是 - https://ptt.cowbay.wtf/rss/movie

要找到英文版名才可以

資料30分鐘才會更新一次, 不會時時更新, 避免被灌