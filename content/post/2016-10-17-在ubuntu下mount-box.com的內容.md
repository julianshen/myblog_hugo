---
date: "2016-10-17T20:38:54Z"
images:
- /images/posts/2016-10-17-在ubuntu下mount-box.com的內容.md.jpg
tags:
- box.com
- Ubuntu
- Linux
title: 在Ubuntu下mount box.com的內容
---

雖然很久沒用box.com的服務了, 不過既然老婆大人問起, 就來寫一下這解法吧

box.com是一個像Dropbox一樣的網路磁碟, 不過它目標客戶跟Dropbox不同, 是比較傾向企業用戶, 可以讓用戶很簡單的分享檔案,
存取box.com除了一般使用Web介面的方式外, 還有其他的方式, 像是透過它的REST API, 另外還有一種就是透過WebDav, 
如果要寫程式去存取它, 一般可以用這兩種方式, 用REST稍微複雜一點, 還要搞定OAuth2的部分, 但透過WebDav的話就簡單多了, 可以掛載成為你作業系統底下的目錄, 當本地檔案來處理

#### 安裝davfs2 ####

首先, 你會需要的是davfs2, 在Ubuntu下用apt-get安裝:

```sh
sudo apt-get install davfs2
```

#### 設定帳號密碼 ####

修改`/etc/davfs2/secrets`, 加入

```
https://dav.box.com/dav box.com帳號 密碼
```

#### 掛載 ####

執行底下指令掛載

```
sudo mkdir /mnt/box.com
sudo -t davfs https://dav.box.com/dav /mnt/box.com
```

成功之後,就會在/mnt/box.com底下看到你的檔案了(包含人家分享給你協作的檔案根目錄), 之後當本地端檔案存取即可, 設定好auto mount即可在開機後掛載