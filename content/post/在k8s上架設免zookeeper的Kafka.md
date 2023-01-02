---
date: 2023-01-02T21:45:11+08:00
title: "在k8s上架設免zookeeper的Kafka"
slug: "Zai-K8sshang-Jia-She-Mian-Zookeeperde-Kafka"
images: 
- "https://og.jln.co/jlns1/5Zyoazhz5LiK5p626Kit5YWNem9va2VlcGVy55qES2Fma2E"
draft: false
---
架設Kafka一個稍稍討厭的點就是先架好Zookeeper, 在Kafka 2.8 (2021/4發布)之後, 支援了以自家的KRaft實現的Quorum controller, 這就可以不用再依賴zookeeper了, Confluent[這篇文章](https://developer.confluent.io/learn/kraft/)有簡單的介紹一下Quorum controller是怎運作的

在我前面[這篇](https://blog.jln.co/dapr-raw-payload-pub-sub/)有提到, 如何用docker跑無zookeeper的Kafka:

```shell
docker run -it --name kafka-zkless -p 9092:9092 -e LOG_DIR=/tmp/logs quay.io/strimzi/kafka:latest-kafka-2.8.1-amd64 /bin/sh -c 'export CLUSTER_ID=$(bin/kafka-storage.sh random-uuid) && bin/kafka-storage.sh format -t $CLUSTER_ID -c config/kraft/server.properties && bin/kafka-server-start.sh config/kraft/server.properties'
```

那如果要架設在K8S上, 可以怎麼做呢? 原本的Kafka需要依賴zookeeper, 加上Kafka的eco system其實蠻多東西的, 一般也不會光只用Kafka本身而已, Kafka Bridge, Kafka connect, schema registry, 進階一點就Kafka stream, KSQL, 規模大一點還需要用上Cruise control, Mirror Maker, 所以用Operator來架設可能會比單純寫manifest, helm chart來的好用, 而比較常見(有名的?)Kafka Operator大致上有這三個(就我知道的啦):

1. [Confluent Operator](https://docs.confluent.io/5.5.1/installation/operator/index.html), 由Confluent這家公司發布的, 由於Confluent這家公司的背景([https://docs.confluent.io/5.5.1/installation/operator/index.html](https://docs.confluent.io/5.5.1/installation/operator/index.html)), Kafka雖是Open source但就是他們家的產品, 所以這個也算是官方出品的Operator, 但這個功能上比較起來稍弱, 而且並沒啥更新, 當然也就還沒看到KRaft相關的支援
1. [KOperator](https://banzaicloud.com/docs/supertubes/kafka-operator/), [萬歲雲(Bonzai Cloud)](https://banzaicloud.com/)出品, 由於Bonzai Cloud目前是Cisco的, 所以這個也可以算大公司出品(?), 我自己是還沒用過, 但看架構, 預設就會架起Cruise control跟Prometheus, 感覺架構上考量是比較完整的, 另外就是也考量到部屬到Istio mesh的部分, 用Envoy來做external LB, 以及用等等, 另外一個值得一提的是Kafka這種Stateful application, 它卻並不是採用Statefulset來部署(它的文件有提到`All Kafka on Kubernetes operators use StatefulSet to create a Kafka Cluster`, 但事實是後來Strimzi也採用一樣的策略了), 但一樣的, 也還沒有支援KRaft
1.  [Strimzi Operator](https://strimzi.io/), 這應該算蠻廣泛被利用的一個Operator, 支援豐富, 更新迅速, 也是可以支援Cruise control (不一定要開), 基本該支援的, 應該也都差不多了, 而且從0.29就支援了KRaft, 不過這個Operator基本消費的記憶體就需要到300MB了

總和以上, 看起來如果要在K8S上玩KRaft的話, Strimzi是一個比較適合的選擇

## 安裝Strimzi operator

用以下指令安裝:

```shell
kubectl create namespace kafka
kubectl create -f 'https://strimzi.io/install/latest?namespace=kafka' -n kafka
```

這除了會建立strimzi-cluster-operator這個Deployment, 也會建立相關的ClusterRoles, ClusterRoleBindings, 和相關的CRD, 所以要先確定你有權限建立這些(尤其是Cluster level的), 其實相當簡單

另外, 用以下的manifest就會幫你建立好一個Kafka Cluster

```yaml
apiVersion: kafka.strimzi.io/v1beta2
kind: Kafka
metadata:
  name: my-cluster
spec:
  kafka:
    version: 3.3.1
    replicas: 3
    listeners:
      - name: plain
        port: 9092
        type: internal
        tls: false
      - name: tls
        port: 9093
        type: internal
        tls: true
    config:
      offsets.topic.replication.factor: 3
      transaction.state.log.replication.factor: 3
      transaction.state.log.min.isr: 2
      default.replication.factor: 3
      min.insync.replicas: 2
      inter.broker.protocol.version: "3.3"
    storage:
      type: ephemeral
  zookeeper:
    replicas: 3
    storage:
      type: ephemeral
```

這會建立一個replica數量為3的Kafka cluster,以及對應的Zookeeper, 這邊的Storage type為ephemeral, 這表示它會用emptyDir當Volume, 如果你有相對應的PVC, 也可以把這替換掉

這邊還是會幫你建立出zookeeper, 那如何能擺脫zookeeper呢?

## 打開實驗性功能

截至這邊文章寫的時間的版本(0.32), KRaft還是一個實驗性功能, 要以環境變數打開, 如下:

```shell
kubectl set env deployments/strimzi-cluster-operator STRIMZI_FEATURE_GATES=+UseKRaft -n kafka
```

strimzi是靠STRIMZI_FEATURE_GATES來當作feature toggle, 在0.32只有一個實驗性功能的開關, 那就是`UseKRaft`, 上面那行指令就可以把這功能打開

用上面一模一樣的Manfest(Zookeeper那段要留著, 雖然沒用, 但在這版本還是必須), 就可以開出一個不依賴zookeeper的kafka cluster了, 以下是整個操作過程:

[![asciicast](https://asciinema.org/a/nibWve6U2E94pn1ljcNx6vscC.svg)](https://asciinema.org/a/nibWve6U2E94pn1ljcNx6vscC)

這邊你可能會發現, 我是用一個strimzipotset的資源來確認是否Kafka有沒正確被開成功, 你如果再去看Replica set, Stateful set, 你會發現找不到Kafka相關的, Strimzi其實就是靠自己的controller來管理Kafka的pods

你也可以用 `kubectl get kafkas -n kafka`來確認kafka這namespace下的kafka cluster的狀況

這個Manifest其實我拿掉了EntityOperator的部分, 是因為KRaft功能目前還沒支援TopicOperator, 沒拿掉會報錯

## 是該拋棄Zookeeper了嗎?
KRaft相對很新, 以三大有名的Kafka Operator來說, 目前也只有Strimzi有支援, 而且才剛開始, 實務上來說, 以功能, 穩定性或許應該還不是時候在production廣泛使用, 真要用, 還是多測試一下再說吧, 短時間還是跟zookeeper做做好朋友

