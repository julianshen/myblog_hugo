---
date: 2023-07-07T21:27:41+08:00
title: "Dead or Alive? Are You Ready?"
slug: "Dead-or-Alive-Are-You-Ready"
images: 
- "https://og.jln.co/jlns1/RGVhZCBvciBBbGl2ZT8gQXJlIFlvdSBSZWFkeT8"
draft: false
---
這篇主要是要來講K8S上的liveness probe/readiness probe/startup probe, 這應該是比較不會被注意到的題目, 大家可能會想, 不過就是health check嘛, 有啥難的? 

不過其實在K8S上面其實不是只有單單health check這麼簡單, 它上面有liveness probe, readiness probe, startup probe, 每一種都有它的不同作用, 如果沒特別注意各自特性, 其實也是有可能會碰到災難的, 首先來看看, 怎麼使用這三種不同的"探針"(翻成探針好像怪怪的, 但我想不到比較好的說法)

這個定義是針對Pod (Deployment ...), 所以你的YAML可能會像這樣

```yaml
apiVersion: v1
kind: Pod
metadata:
  labels:
    test: liveness
  name: liveness-http
spec:
  containers:
  - name: liveness
    image: registry.k8s.io/liveness
    args:
    - /server
    livenessProbe:
      httpGet:
        path: /healthz
        port: 8080
        httpHeaders:
        - name: Custom-Header
          value: Awesome
      periodSeconds: 3
    startupProbe:
      httpGet:
        path: /healthz
        port: liveness-port
      failureThreshold: 5
      initialDelaySeconds: 30
      periodSeconds: 10 
    readinessProbe:
      httpGet:
        path: /healthz
        port: liveness-port
      failureThreshold: 30
      periodSeconds: 10   
```

先說一下, 上面並不是一個好寫法, 三種Probe都用了同一個endpoint, 先不多做解釋, 看完應該就比較會清楚啥不好

不管哪一種, 都有三種檢查機制可以用:
- exec 執行一個特別的命令, status code是0則是成功
- grpc 呼叫grpc來確認狀態(參考 [https://grpc.github.io/grpc/core/md_doc_health-checking.html](https://grpc.github.io/grpc/core/md_doc_health-checking.html))
- httpGet 呼叫HTTP GET來確認狀態, HTTP status 200-400(以內)都算成功
- tcpSocket 可以建立連線就表示沒問題

就是選擇適合來確認你服務的健康狀態的來使用就好了(廢話)

那這三種有啥分別?

- livenessProbe 用來確認你的服務是不是還在"執行", 如果不是, pod就會被砍掉, 然後會依據restartPolicy的設定是不是要重啟你的pod
- readinessProbe 活著(liveness)和準備好(readiness)差別在哪? 這邊的readiness指的是你的容器(container), 是不是已經"準備好提供服務"了, 如果是, 才會把請求(request)轉送到你的容器
- startupProbe 特性其實接近livenessProbe, 但是是用在容器剛啟動時, 用來確認容器是否正常啟動, 如果檢查失敗(時間超過 initialDelaySeconds[*一開始等多久*] + failureThreshold[*失敗幾次*] × periodSeconds[*間隔多久重試*])一樣是砍掉pod, 然後會依據restartPolicy的設定是不是要重啟你的pod

詳細點來說,

## 啥時用livenessProbe?

俗語說的好, "__`重開機治百病`__"(哪來的俗語呀?), 簡單的說, 如果你的容器卡住不動了, 怎麼搖都懷疑人生, 沒反應了, 需要透過重開, 重新投胎, 才能(有機會)恢復正常, 那就是livenessProbe可以發揮的地方了

這邊說的是卡住不動沒回應這類的, 所以像是你的程式碰到dead lock, 無窮迴圈而無法正常回應都算, 但程式結束, 不正常離開(像Java的uncaught RuntimeException?) 其實不用等到livenessProbe打失敗, 就會照restartPolicy來處理, 所以可以知道, 當發現有Pod重啟的情形, 應該就不外乎是自然死, 意外死, 還有就是因為livenessProbe打施打失敗被謀殺了

那....

### 該不該在livenessProbe去檢查相關的服務或資料庫有沒活著?

不該!(張惠妹/周杰倫), 這還蠻常見的狀況, 就想說, 這health check 嘛, 我資料庫連不上, 當然就不健康囉, 所以就回傳了Failure了, Spring boot acuator的health endpoint也是會幫你檢查相關依賴的資源的狀況列入檢查(不過它有為Kubernetes有對應的做法啦, 後面再說)

這樣會導致啥狀況? 你明明服務還好好的, 然後只是資料庫連不上, 結果你的Pod就被砍了(好無辜), 然後只要資料庫還沒修好, 就一直復活一直死(好悲哀)

所以liveness的probe應該簡單到只是確認這個容器有沒被卡死, 資料庫連不上只是不能服務, 資料庫修好了還是可以繼續服務呀

### 老是因為打livenessProbe時timeout被砍, 那我是不是盡量把timeoutSeconds盡量設的越大越好?

其實這檢查的目的只是要確認容器有沒被卡死, 所以livenessProbe應該盡可能越簡單越好, 不太適合去做一些需要大量運算或是複雜的事, 因為那可能會因為你的pod或node的繁忙程度去影響到它執行時間長短差距很大, 那如果因為這樣去調高timeout時間也不太合理, 因為也很難確定要多久才能確定它"真的被卡住了", 再加上你可能會因為頻繁做這些複雜運算(因為每periodSecond會被探測一次)影響系統效能(像是不要在livenessProbe實作內call DB query)

依據你livenessProbe正常會回應的時間再多給點應該就足夠了, 設非常大的話, 搞不好容器真的卡住了, 但卻反應慢了

### 我可不可以不要設livenessProbe

為何不可? 前面一直說, 它只是用來偵測程式是不是被"卡住了", 如果沒被卡死的風險, 不需要靠重啟手段來回復的話, 那不用設是可以呀

## 那啥時用readinessProbe呢?

readinessProbe跟livenessProbe的差異在於會不會出人(Pod)命, readiness為success的話, 上游(Service)才會把請求(request)送來給我處理, 不然的話, 就會收不到(又是廢話), 所以從這邊就可以看出前面那題的答案了, liveness和readniess的檢查邏輯應該會不一樣的(所以不太適合同一個endpoint搞定)

### 那要不要去檢查相關的服務或資料庫有沒活著?

可以, 也建議, 因為如果後面的服務或資料庫死了, 表示請求送進來也會處理失敗, 那不如先把它擋在外面等到服務正常了

## startupProbe呢?

這個通常用在啟動會很久的容器, 為了怕太早打livenessProbe, readinessProbe導致高失敗率(因為啟動很久, 太早打一定都失敗的), 所以用這個probe來確認容器成功後才真的去實施那兩個探測

容器啟動很久其實不是一個很好的practice, 所以這個其實也是萬不得以才在用, 如果啟動時間不長的話, 為probe設定initialDelaySeconds 就已經很足夠了

## 誰去打這些Probe的?

一個比較錯誤的想像是, Kubernetes在某個地方, control plane或那裡有個服務去打所有的probe, 其實不是, 這樣它會累死

其實是由每個node的kubelet來負責, 當被加入一個Pod時, kubelet就會為這些probe每個都起一個go routine來根據設定的規則做檢查(感覺這設計沒太好, 會起不少go routine)

可以參考kubelet的[實做細節](https://github.com/kubernetes/kubernetes/blob/7581ae812327fc8218204f678143a6f116cad931/pkg/kubelet/prober/prober_manager.go#L169-L213)

這樣其實比較合理啦, 每個node的kubelet就顧自己家後院就好

## Spring boot acuator的Kubernetes Probes

Spring boot acuator預設的health endpoint是`/actuator/health`, 但這其實不好一體適用於liveness和readiness

Spring boot要用`/actuator/health/liveness`在livenessProbe而`/actuator/health/readiness`在readinessProbe, 可以參考[這篇....](https://docs.spring.io/spring-boot/docs/2.3.0.RELEASE/reference/html/production-ready-features.html#production-ready-kubernetes-probes)