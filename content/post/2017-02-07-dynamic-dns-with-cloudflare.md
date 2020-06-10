---
date: "2017-02-07T15:40:05Z"
tags:
- DNS
- Cloudflare
title: Dynamic DNS with Cloudflare
---

這陣子都在寫line bot, 本來都host在heroku上面的, 簡單且方便, 後來突發奇想, 想用Rasberry pi 跑看看(跑得動喔)

Line的webhook有一個需求就是要有SSL連結, 走https, 但我不想申請一個certificate, 在raspberry pi上弄, 所幸[Cloudflare](https://www.cloudflare.com/)
有提供免費的SSL certificate, 利用他們的flexible SSL就可以了

![](https://support.cloudflare.com/hc/en-us/article_attachments/206124658/cfssl_flexible.png)

flexible SSL的方式是client到他們CDN server走的是SSL沒錯, 但他們server到你的server則是可以走一般的http connection,
再來的第二個問題是, 我家的網路是浮動IP的（後來才去申請固定IP）, 所以必須能動態更新[Cloudflare](https://www.cloudflare.com/)上的DNS紀錄

還好Cloudflare是有[API](https://api.cloudflare.com/)的

直接自己自幹一個也是可以啦, 但Cloudflare其實也有一個客製版的ddclient:

[Dynamic DNS Client: ddclient](https://www.cloudflare.com/technical-resources/#ddclient)

步驟可以照著上面文件的步驟來做, 可以從My settings -> Account -> Global API Key取得API key當作ddclient的密碼

接下來碰到的問題是, 我ddclient是跑在raspberry pi上, ddclient預設是用local IP, 這很明顯不對, 因為會用到內部的IP而不是對外那個, 而我家的ASUS無線分享器並沒支援Cloudflare, 我也不太想改firmware,
但這還是有解的, 把ddclient.conf裡加上這行:

```
use=web, web=checkip.dyndns.org/, web-skip='IP Address' # found after IP Address
```

這是告訴ddclient不要用local ip而是用web api去找出IP

但這一切....都還是太麻煩了....raspbeery pi總是會不小心碰掉電源, 總是會當機或跑不動, 更何況, 我都已經跑一個server了, ddclient不要再來搶記憶體了啦

最後我的解法是: [DNS-O-Matic](https://www.dnsomatic.com/)

這是一個Dynamic DNS的服務, 我的無線分享器也有支援, 它不是自身有DNS server, 而是可以代你去更新妳的DNS紀錄, 而且, 有支援Cloudflare!!! OK, 結案 (偷懶!)