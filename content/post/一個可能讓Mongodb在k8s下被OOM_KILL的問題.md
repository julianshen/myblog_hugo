---
date: 2022-12-10T02:15:17+08:00
title: "一個可能讓Mongodb在k8s下被OOM_KILL的問題"
slug: "Ge-Ke-Neng-Rang-Mongodbzai-K8sxia-Bei-Oom_killde-Wen-Ti"
images: 
- "https://og.jln.co/jlns1/5LiA5YCL5Y-v6IO96K6TTW9uZ29kYuWcqGs4c-S4i-iiq09PTV9LSUxM55qE5ZWP6aGM"
draft: false
---

做了些實驗, 紀錄一下, 剛好寫這篇時巴西也被幹掉了, 也記錄一下 XD

這個問題其實要滿足以下條件才可能發生:
- MongoDB版本在 4.4.14, 5.3.0, 5.0.7, 4.2.20 之前, 這算是一個Mongodb在2022一月才fix的一個bug, 所以比這幾版舊是有可能的
- MongoDB instance在K8S上有設memory limit (且這limit要小於host memory的一半?)
- K8S所在的Host OS 的cgroup版本為V1, 可以[參考這文件](https://github.com/opencontainers/runc/blob/main/docs/cgroup-v2.md), Ubuntu 21.10, Fedora 31之後都開啟V2了, 不過如果你用的是WSL2, 由於WSL2的Kenel還是V1, 是試不出這問題的 (我是找了台Fedora來試)

### 問題是甚麼?

查這問題的起因當然是碰到Mongodb被OOM Kill, 後來發現好像這也算蠻常踩到的坑, 只是好像沒人寫出完整可能性

碰到被OOM Kill第一個會思考的是, 他為何要那麼多記憶體?要給他多少才夠? 另外一個是, 由於是發生在container, 跑在K8S上, 一個疑問是, 那MongoDB是否會遵守設給他的resource limit? 還是他會當node所有記憶體都是他可用的?

### 有沒有哪裡可疑的?

有哪些東西會吃記憶體? 連線會, index會, 但其實其中一個比較可疑的是給Wired Tiger的cache, 根據[這份文件](https://www.mongodb.com/docs/manual/core/wiredtiger/), Wired Tiger的cache會用掉

- 50% 的(總記憶體 - 1GB) 或是
- 256MB

也就是至少256MB, 然後如果你有64G記憶體, 他就會用掉最多 (64-1)/2, 到這聽起來好像沒啥問題, 只用一半還不至於有撐爆的問題, 會不會是其他的地方?

但另一個問題是, 直接裝在單機沒問題, 如果是跑在K8S上的容器, Memory limit我們是給在K8S上, MongoDB到底會以memory limit當總記憶體大小還是以整個node全部可用的記憶體計算?

其實根據文件的補充說明, 它是有考慮到的, 它會以[hostInfo.system.memLimitMB](https://www.mongodb.com/docs/manual/reference/command/hostInfo/#mongodb-data-hostInfo.system.memLimitMB) 來計算

```
In some instances, such as when running in a container, the database can have memory constraints that are lower than the total system memory. In such instances, this memory limit, rather than the total system memory, is used as the maximum RAM available.
```

而它這資訊是透過cgroup去抓的, K8S也是用cgroup做資源管理的, 所以這值會等於你設定的limit

我第一次在WSL下測試, 把memory limit設為2Gi, [hostInfo.system.memLimitMB](https://www.mongodb.com/docs/manual/reference/command/hostInfo/#mongodb-data-hostInfo.system.memLimitMB) 也的確是這個值(用mongo client下`db.hostInfo()`即可查詢)

那看來應該沒問題呀, 問題在哪?

後來查到一個bug : [https://jira.mongodb.org/browse/SERVER-60412](https://jira.mongodb.org/browse/SERVER-60412), 原來cgroup v1, v2抓這些資訊的位置是不同的, 所以導致舊版的會有抓不正確的狀況

看到這就來做個實驗, 找了台有開啟v2的fedora (with podman), 跑了k3s, 在這k3s上分別跑了4.4.13, 4.4.15兩個版本去做測試, memory limit都設為2Gi, 用`db.hostInfo()`查詢memLimitMB得到下面結果:

#### 4.4.13
![](/images/posts/mongo4.13.mem.png)

#### 4.4.15
![](/images/posts/mongo4.15.mem.png)

Bingo! 4.4.13果然抓到的memLimitMB是整個node的記憶體大小而非limit, 這樣如果node的記憶體大小遠大於limit, Wired Tiger cache是有可能用超過limit的

當然, 這只是其中一種可能性, 不見得一定都是這情形, 但碰到這類狀況, 這的確是可以考慮查的一個方向