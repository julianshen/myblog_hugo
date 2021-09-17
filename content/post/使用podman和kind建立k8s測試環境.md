---
date: 2021-09-17T14:48:19+08:00
title: "使用podman和kind建立k8s測試環境"
images: 
- "https://og.jln.co/jlns1/5L2_55SocG9kbWFu5ZKMa2luZOW7uueri2s4c-a4rOippueSsOWigw"
---

[Docker Desktop要收錢了](https://www.linuxadictos.com/zh-TW/Docker-%E6%A1%8C%E9%9D%A2%E5%B0%87%E4%B8%8D%E5%86%8D%E5%B0%8D%E4%BC%81%E6%A5%AD%E5%85%8D%E8%B2%BB%EF%BC%8C%E7%8F%BE%E5%9C%A8%E5%B0%87%E6%8C%89%E6%9C%88%E8%A8%82%E9%96%B1%E9%80%B2%E8%A1%8C%E7%AE%A1%E7%90%86.html), 雖然不是跟個人開發者收, 而且一家公司走向營利也合理, 但這作法說實在有點粗糙, 是時候改用不同的工具來玩了

現在在開發階段多多少少會需要在local有一個測試環境亂搞, 但一般的小電腦, 跑起Docker + K8S, 也沒辦法跑太多容器了, 能夠有盡量輕量化的東西當然最好, 在我的linux PC上, 我現在用的是podman + KIND (Windows下我還是用Docker desktop, 還沒切換過去)

## Podman

[Podman](https://podman.io/)是一個由Redhat開發的工具, 跟Docker最大的不同在於, 它沒需要跑一個daemon常駐在那邊, Docker desktop即使你沒跑任何container狀況下daemon還會在, [Podman](https://podman.io/)則完全沒這問題

安裝的話請參考: [Podman Installation Instructions](https://podman.io/getting-started/installation)

安裝好後, 指令幾乎跟docker 類似, 像是 `docker ps`可以用`podman ps`取代, `docker run`可以用`podman run`取代, 幾乎沒太大問題

有些狀況, 會需要有docker daemon, 像是如果使用[pack](https://buildpacks.io/docs/tools/pack/), 由於[pack](https://buildpacks.io/docs/tools/pack/)會需要呼叫docker daemon來建立image, 如果使用podman, 在這狀況就得跑一個service:

```
podman system service --time=0 tcp:localhost:1234
```

這邊可以是tcp或是unix socket, 然後把DOCKER HOST改指到這邊來就好了

## [kind](https://kind.sigs.k8s.io/)

[kind](https://kind.sigs.k8s.io/)是一個讓你建立local K8S cluster的工具, 其他類似的還有[MicroK8s](https://microk8s.io/)和[MiniKube](https://github.com/kubernetes/minikube)

為何選[kind](https://kind.sigs.k8s.io/)? 有時候會需要模擬多個nodes的cluster環境, 不管是MicoK8s還是MiniKube, 在單機建立出的k8s都只有一個node, 但[kind](https://kind.sigs.k8s.io/)卻可以在單機建立出多nodes的環境

在Mac下可以用Homebrew安裝:

```
brew install kind
```

建立一個cluster很簡單, 只要執行

```
kind create cluster
```

如果你是要用[Podman](https://podman.io/)也可以
```
KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster
```

但第一次使用podman建立會發生一個問題:

```
KIND_EXPERIMENTAL_PROVIDER=podman kind create cluster
using podman due to KIND_EXPERIMENTAL_PROVIDER
enabling experimental podman provider
Creating cluster "kind" ...
 ✗ Ensuring node image (kindest/node:v1.21.1) 🖼 
ERROR: failed to create cluster: failed to pull image "kindest/node@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6": command "podman pull kindest/node@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6" failed with error: exit status 125
Command Output: Error: short-name resolution enforced but cannot prompt without a TTY
```

這是由於podman在pull image時碰到short name(像是java, ubuntu, fedora這類的), 它會請你從`/etc/containers/registries.conf`裡面設的search host挑一個是可以找到這個image的host, 但在kind跳不出來給你挑, 這時候只要看哪個pull不下來(像這邊是`kindest/node@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6`)就自己手動執行一次

```
podman pull kindest/node@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6
```

一次就好, 之後它會記起來, 這時候再來跑kind就沒問題了

要毀掉一個cluster也很簡單
```
kind delete cluster
```

雖然文件上有說支援[Rootless](https://kind.sigs.k8s.io/docs/user/rootless/), 不過實際上試, 要排除的問題還很多, 不建議使用

如果你需要用`kubectl`或是[Lens](https://k8slens.dev/)去存取建立出來的cluster, 那可以輸出kubeconfig
```
kind export kubeconfig
```

那, 說好的多肉...喔...多nodes環境呢? 建立以下這樣的config(例如檔名叫config.yaml):

```yaml
kind: Cluster
apiVersion: kind.x-k8s.io/v1alpha4
nodes:
- role: control-plane
- role: worker
- role: worker
- role: worker
```

然後用```kind create cluster --config config.yaml```就可以建立出一個4 nodes (一個control plane, 三個worker)的環境了(在單機)
