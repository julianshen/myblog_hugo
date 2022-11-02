---
date: 2022-11-02T22:01:38+08:00
title: "在與世隔絕的環境安裝Dapr開發環境"
images: 
- "https://og.jln.co/jlns1/5Zyo6IiH5LiW6ZqU57WV55qE55Kw5aKD5a6J6KOdRGFwcumWi-eZvOeSsOWigw"
draft: false
---

這邊的與世隔絕當然不是真的與世隔絕啦! 我指的是無法連上外面的docker registry的環境, 例如docker hub或是quay.io, 開發環境指的是你要在local可以用來開發測試的standalone模式

首先, 你要先裝好docker(或podman)

### 如果有private registry

有幾個images會需要放到registry裡面的, 像是

- dapr
- 3rdparty/redis
- 3rdparty/zipkin
- placement
- daprd

列表可以在 [這邊](https://github.com/orgs/dapr/packages?repo_name=dapr) 找到 (這邊的是放在github上的)

接著安裝dapr cli, 如果可以直接連上internet, 那就用官方文件的做法:

```sh
wget -q https://raw.githubusercontent.com/dapr/cli/master/install/install.sh -O - | DAPR_INSTALL_DIR="$HOME/dapr" /bin/bash
```

但如果不行呢? 那就想辦法到[github上下載](https://github.com/dapr/cli/releases/tag/v1.9.1)cli回來安裝

因為要從private history來安裝dapr, 所以在`dapr init`要多下一個參數, 像是:

```sh
dapr init --image-registry MY_REGISTRY_URL
```

這樣就會從你private registry抓回來裝了

但如果連private registry都沒辦法放上相關的image呢?

### 不依靠docker registry安裝

Dapr很貼心呀! 還有[install bundle](https://github.com/dapr/installer-bundle/releases), 可以想辦法先去前一個連結裝bundle回來

這個bundle裡面已經有個dapr cli了, 所以解開後, 把dapr複製到你要的目錄, 例如:

```sh
sudo cp ./dapr /usr/local/bin
```

現在我們就有cli可以用了, 但images怎辦? bundle裡面其實就包含有相關的image的tar檔了, 所以把init方式改成:

```sh
dapr init --from-dir . -s
```

這樣就會用local image安裝了, 這邊多加了個`-s`表示是slim mode, 沒有redis, 沒zipkin的, 因為local images沒有放, 但如果自己需要(比如說要寫state store需要redis), 那就要自己去加component

以redis當state store為例: 假設我們需要一個redis來當state store, 且這redis我們預先跑在本機, port為為6379, 此store名稱為mystore, 我們可以在`~/.dapr/components/`這目錄加上一個mystore.yaml, 

```yaml
apiVersion: dapr.io/v1alpha1
kind: Component
metadata:
  name: mystore
spec:
  type: state.redis
  version: v1
  metadata:
  - name: redisHost
    value: localhost:6379
  - name: redisPassword
    value: ""
  - name: actorStateStore
    value: "true"
```

這樣我們就有一個名叫mystore的state store了, 不只state store, 其他也可以如法炮製 