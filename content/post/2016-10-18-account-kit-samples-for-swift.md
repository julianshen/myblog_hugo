---
date: "2016-10-18T21:01:31Z"
images:
- /images/posts/2016-10-18-account-kit-samples-for-swift.md.jpg
tags:
- Swift
- Facebook
- Account Kit
title: Account Kit samples for Swift
---

做一個網路服務, 使用者驗證是蠻麻煩的一件事, 我們是可以裝作沒看到, 不做驗證, 但這樣的下場就是有假使用者, 有殭屍, 最簡單的方式是信任第三方服務像是Google, Facebook,
現在的人大多數都有Google, Facebook帳號了, 這樣其實也沒多大的問題, 但還是還是有人沒有, 而且也不是每個人都放心把Facebook帳號交給我們, 因此退而求其次就是用E-mail,
用E-mail認證雖然也是一個好方式, 但還是要建置整套發信機制(或是花錢買mailgun來送信), 而且在手機上就麻煩了, 來回在App跟e-mail間切換很不方便,
因此就會想用簡訊認證, 至於簡訊認證, 除了一個"貴"字以外, 要搞定各個國家的也是一個麻煩(當然, 花錢可解, 有Twilio這種服務可以用)

所幸有Facebook提供的這個可以用[Account Kit](https://developers.facebook.com/docs/accountkit), 在初期使用者不太多的時候, 不收費的確很吸引人呀(雖然他不是唯一一個這樣的服務, 之後再介紹其他的),
但由於他[iOS的範例](https://github.com/fbsamples/account-kit-samples-for-ios)是用Objective C寫的, 我的Objective C實在不太行,
加上, 要了解一個東西, 寫一遍就知道了, 所以順手翻譯了一個Swift的版本, Source如下:

## [Account Kit samples for Swift](https://github.com/julianshen/account-kit-samples-for-swift) ##

原本的版本, 我是覺得寫的不是太好, 花了好一些功夫看, 自己翻過來的這個版本, 也還沒debug過, 基礎的應該堪用啦! 至於iOS的account kit文件可以參考: [iOS 專用 Account Kit](https://developers.facebook.com/docs/accountkit/ios)

註: plist裡面寫的app id不是我的喔!是原本Sample用的