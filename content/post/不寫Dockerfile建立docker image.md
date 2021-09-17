---
date: 2021-09-17T14:21:09+08:00
title: "不寫Dockerfile建立docker Image"
images: 
- "https://og.jln.co/jlns1/5LiN5a-rRG9ja2VyZmlsZeW7uueri2RvY2tlciBJbWFnZQ"
---

人都是懶的, 尤其如果拿同樣的工具, 開發不同的service, 又要部屬到cloud native環境, 免不了要一直重複寫類似的Dockerfile來建立docker image, 有沒懶方法?

有, 就是[Cloud native buildpacks](https://buildpacks.io/), 有用過[Heroku](https://www.heroku.com/)的應該會有點耳熟, 對, 就是那個buildpacks, 只是把它變成一個標準

首先, 你需要的是[pack](https://buildpacks.io/docs/tools/pack/)這工具, 在mac底下可以用 brew安裝

```
brew install buildpacks/tap/pack
```

Windows下可以用scoop

```
scoop install pack
```

安裝好後, 可以執行:

``` 
pack builder suggest
```

來看看有甚麼buildpacks 可以用

```
Suggested builders:
    Google:                gcr.io/buildpacks/builder:v1      Ubuntu 18 base image with buildpacks for .NET, Go, Java, Node.js, and Python
    Heroku:                heroku/buildpacks:18              Base builder for Heroku-18 stack, based on ubuntu:18.04 base image
    Heroku:                heroku/buildpacks:20              Base builder for Heroku-20 stack, based on ubuntu:20.04 base image
    Paketo Buildpacks:     paketobuildpacks/builder:base     Ubuntu bionic base image with buildpacks for Java, .NET Core, NodeJS, Go, Python, Ruby, NGINX and Procfile
    Paketo Buildpacks:     paketobuildpacks/builder:full     Ubuntu bionic base image with buildpacks for Java, .NET Core, NodeJS, Go, Python, PHP, Ruby, Apache HTTPD, NGINX and Procfile
    Paketo Buildpacks:     paketobuildpacks/builder:tiny     Tiny base image (bionic build image, distroless-like run image) with buildpacks for Java Native Image and Go
```

看你使用的開發語言或框架, 選擇適合的buildpack, 比如說像是node.js或是go, 只要簡單的在你的project底下執行:

```
pack build image_name --builder gcr.io/buildpacks/builder:v1
```

就可以建立出一個名為`image_name`的docker image了, 而且建置過程都是自動偵測

不過當然沒辦法適用所有的狀況, 如果你有不同的語言不同的需求, 其實也可以自建buildpack喔