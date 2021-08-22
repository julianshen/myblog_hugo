---
date: 2021-08-22T12:42:50+08:00
title: "用Kubernetes ConfigMap實現配置的熱更新"
images: 
- "https://og.jln.co/jlns1/55SoS3ViZXJuZXRlcyBDb25maWdNYXDlr6bnj77phY3nva7nmoTnhrHmm7TmlrA"
---

程式配置(Configuration)的熱更新(hot reload)應該是建置服務會常碰到一個題目, 常會有狀況需要在不動用release去調整程式配置的狀況, 比較常見的做法應該是將這些配置集中管理, 因此就有相關的解決方案產生像是:

 * [Spring cloud config server](https://cloud.spring.io/spring-cloud-config/multi/multi__spring_cloud_config_server.html)
 * [Netflix Archaius](https://github.com/Netflix/archaius)
 * [LINE Central Dogma](https://line.github.io/centraldogma/)
 * [HashiCorp Consul](https://www.hashicorp.com/products/consul) (不過近來這個產品已經被延伸到Service Mash的領域去了)
 * [Azure App configuration](https://azure.microsoft.com/en-us/services/app-configuration/)

真要找, 應該還有, 這種中央管理的方式, 無非就是想要把分布在不同系統的所有的設定, 做一個集中管理, 隨時可以進行線上更新, 不過帶來的問題點就是除了要綁定選定系統用相關的API開發外, 這類的服務也是有可能是SPOF 

在Kubernetes原生(Kubernetes Native)的角度來看這件事, Kubernetes就有內建ConfigMap, Secret, 是否還有必要導入這類的解決方案? 利用ConfigMap是否可以達成線上做熱更新的目的? 我的想法是, 如果用ConfigMap做到熱更新, 那麼搭配 GitOps 的流程, 這樣就可以做到簡單又兼顧集中管理的特性了(更新紀錄在git都可以查到, 另外可以用PR確保更改配置的安全性, 避免誤更, 在多叢集配置下也可以分享同一個git repository) 

## 使用ConfigMap

這邊沒特別要說明怎麼去用ConfigMap, 那個 [官方文件](https://kubernetes.io/docs/concepts/configuration/configmap/) 寫得很清楚, 先來看看ConfigMap在配合Pod/Deployment的兩個常見用法

先拿下面這範例來看:

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: game-demo
data:
  # property-like keys; each key maps to a simple value
  player_initial_lives: "3"
  ui_properties_file_name: "user-interface.properties"

  # file-like keys
  game.properties: |
    enemy.types=aliens,monsters
    player.maximum-lives=5    
  user-interface.properties: |
    color.good=purple
    color.bad=yellow
    allow.textmode=true 
```

上面這個ConfigMap我們可以在Pod這樣使用它:

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: configmap-demo-pod
spec:
  containers:
    - name: demo
      image: alpine
      command: ["sleep", "3600"]
      env:
        # Define the environment variable
        - name: PLAYER_INITIAL_LIVES # Notice that the case is different here
                                     # from the key name in the ConfigMap.
          valueFrom:
            configMapKeyRef:
              name: game-demo           # The ConfigMap this value comes from.
              key: player_initial_lives # The key to fetch.
        - name: UI_PROPERTIES_FILE_NAME
          valueFrom:
            configMapKeyRef:
              name: game-demo
              key: ui_properties_file_name
      volumeMounts:
      - name: config
        mountPath: "/config"
        readOnly: true
  volumes:
    # You set volumes at the Pod level, then mount them into containers inside that Pod
    - name: config
      configMap:
        # Provide the name of the ConfigMap you want to mount.
        name: game-demo
        # An array of keys from the ConfigMap to create as files
        items:
        - key: "game.properties"
          path: "game.properties"
        - key: "user-interface.properties"
          path: "user-interface.properties"
```

一個是用`valueFrom`把ConfigMap裡面的設定拿來放在環境變數使用(參考上面範例)

另一個則是透過 `volumes` 把設定內容掛載成檔案

為了達成熱更新, 我們會有興趣的是, 當我們ConfigMap更新時, 相對應的內容會不會改變, 答案是只有第二種掛載成檔案的, 會隨之更新, 而第一種, 當ConfigMap更新時, 相關的環境變數是不會跟著變的

至於掛載成檔案的, 當ConfigMap內容做過更動時, 相對應的檔案內容也會更新, 但...不是即時的, 根據文件

```
The kubelet checks whether the mounted ConfigMap is fresh on every periodic sync. However, the kubelet uses its local cache for getting the current value of the ConfigMap. The type of the cache is configurable using the ConfigMapAndSecretChangeDetectionStrategy field in the KubeletConfiguration struct. A ConfigMap can be either propagated by watch (default), ttl-based, or by redirecting all requests directly to the API server. As a result, the total delay from the moment when the ConfigMap is updated to the moment when new keys are projected to the Pod can be as long as the kubelet sync period + cache propagation delay, where the cache propagation delay depends on the chosen cache type (it equals to watch propagation delay, ttl of cache, or zero correspondingly).
```

也就是說預期會有根據你設定是用watch, ttl-based, 全透過API取得更新跟cache時間造成的時間差, 也就是雖然ConfigMap也是一種集中式管理(放在etcd), 但實際上還是會有數秒到數十秒的更新時間差(我實測最多碰到一分鐘後才更新)

因此如果需要做到配置的熱更新, 那我們可以選擇是第二種掛載成檔案的作法, 藉由監控檔案內容的改變, 再由程式去做熱更新

## 觀測ConfigMap的異動狀況

既然是檔案, 那我們可不可以由Linux的inotify去監控檔案的異動狀況? inotify是Linux核心的一個系統呼叫, 現在主流伺服端的程式設計應該也比較少用C直接去呼叫這些System call了吧? 不過, 基本上還是可行的, 這邊有一篇"[用 Sidecar 应用 Configmap 更新](https://cloud.tencent.com/developer/article/1557278)", 這邊就用 `inotifywait` 這個指令放在sidecar中去監控config檔案, 在有變動時, 發送訊號重啟主程序

這方法的優點是, 程式可以不用自行監控ConfigMap的變化, 缺點就是, 重啟這件事是不可控的, 當你的服務有多個實體(instance)時, 也有可能這些全部會在同一時間被重啟, 造成你的服務被下線

另外一個就是在程式內自行監控, 現在Kubernetes大行其道, 已然是顯學, 如果已經採用它來管理配置系統的話, 在設計上配合它來做, 也是無可厚非, Dev要能針對Ops來設計, 才能真的有DevOps, 更何況這部分只需要監控檔案, 並不需要綁死Kubernetes API

監控檔案異動的作法, 各語言有自己包裝, golang有[fsnotify](https://fsnotify.org/), Java則有nio裡的[WatcherService](https://www.baeldung.com/java-nio2-watchservice)

這邊先以Java Nio做一個簡單的測試(實際是以Kotlin實作):

```kotlin
suspend fun watchConfig(configFileName: String) {
	val dir:Path = Paths.get(configFileName).parent
	val fileName = Paths.get(configFileName).fileName

	val watcher = FileSystems.getDefault().newWatchService()

	dir.register(watcher, StandardWatchEventKinds.ENTRY_CREATE, StandardWatchEventKinds.ENTRY_DELETE, StandardWatchEventKinds.ENTRY_MODIFY)
	while(true) {
		val key =watcher.take()
		key.pollEvents().forEach { it ->
			if(it.context() == fileName) {
				reloadConfig()
			}
		}

		if(!key.reset()) {
			key.cancel()
			watcher.close()
			break
		}
	}
}
```

以前面config map掛載的範例來看的話, 假設, 我們用 `watchConfig("/config/game.properties")` 來監控"/config/game.properties", 這邊的"/config/game.properties"是由ConfigMap裡的 `game.properties` 來的, 所以變更這邊的`game.properties`, "/config/game.properties"也會跟著改變

但, 上面這段程式是"完全沒用的", 即使 `game.properties`和`/config/game.properties`內容都被改變了, 這邊的 `reloadConfig()` 也完全不會被觸發!!!! 如果使用golang的fsnotify, 也會是一樣的狀況

為什麼呢? 難道是這樣掛載的檔案有啥特異? 先來 `ls -l`看一下:

```
ls -l /config/game.properties
lrwxrwxrwx 1 root root 24 Aug 21 16:55 /config/game.properties -> ..data/game.properties
```

這邊可以發現`/config/game.properties`是一個Symbolic link連到`..data/game.properties`去, 這樣就導致我們監控不到它嗎? 其實還不只, 再來`ls -l ..data`看看

```
ls -l ..data
lrwxrwxrwx 1 root root 31 Aug 21 16:58 ..data -> ..2021_08_21_16_56_33.873456784
```

Ok, `..data`也是一個Symbolic link, 所以實際上ConfigMap被變更過後, 真正檔案變更大guy會是這樣的:

```
CREATE ..2021_08_21_16_58_04.661956783
CREATE ..2021_08_21_16_58_04.661956783/game.properties
CREATE ..data_tmp (link to ..2021_08_21_16_58_04.661956783)
MOVE ..data_tmp ..data
DELETE ..2021_08_21_16_56_33.873456784
```

所以我們原本直覺應該會是認為它是會直接變更`/config/game.properties`內容, 但實際上`/config/game.properties`是一直沒被變動的, 它一直是一個連結到`/config/..data/game.properties`的Symbolic link, 所以觀測對象是不對的, 因此得這樣改:

```kotlin
suspend fun watchConfig(configFileName: String) {
	var path:Path = Paths.get(configFileName)
	val parent:Path = path.parent

	while (Files.isSymbolicLink(path)) {
		path = Files.readSymbolicLink(path)
	}

	val realParent:String = path.parent.name
	val watcher = FileSystems.getDefault().newWatchService()

	parent.register(watcher, StandardWatchEventKinds.ENTRY_CREATE, StandardWatchEventKinds.ENTRY_DELETE, StandardWatchEventKinds.ENTRY_MODIFY)
	while(true) {
		val key =watcher.take()
		key.pollEvents().forEach { it ->
			if(it.context() == realParent) {
				reloadConfig()
			}
		}

		if(!key.reset()) {
			key.cancel()
			watcher.close()
			break
		}
	}
}
```

這邊的`realParent`其實就是`..data`, 有變動的會是它, 所以監控它就好了

## 使用golang的spf13/viper

如果你是用golang並且是用sp13大神的[viper](https://github.com/spf13/viper), 來管理設定檔, 那你只需要透過`viper.WatchConfig()`來監控ConfigMap掛載下來的設定檔就好

```golang
viper.WatchConfig()
viper.OnConfigChange(func(e fsnotify.Event) {
	fmt.Println("Config file changed:", e.Name)
})
```

這是因為[viper](https://github.com/spf13/viper)有針對這一狀況修正過, 有興趣可以參考["WatchConfig and Kubernetes (#284)"](https://github.com/spf13/viper/commit/e0f7631cf3ac7e7530949c7e154855076b0a4c17)這段

## Reloader

如果程式不想配合著改, 或大部分都是透過環境變數的方式來使用ConfigMap的話, 又怕使用前面inotify sidecar的作法會造成問題, 希望有更好的方式去RollOut, 那可以參考一下[Reloader](https://github.com/stakater/Reloader)

[Reloader](https://github.com/stakater/Reloader)會去監控ConfigMap跟Secret的變動, 來重啟跟他們有相關的DeploymentConfigs, Deployments, Daemonsets Statefulsets 和 Rollouts, 由於它是以Kubernetes conrtroller的形式存在, 並且採用Kubernetes API去監控資源: https://github.com/stakater/Reloader/blob/99a38bff8ea1346191b6a96583d3fbad72573ea5/internal/pkg/controller/controller.go#L47

安裝方法很簡單, 只需要用:

```
kubectl apply -f https://raw.githubusercontent.com/stakater/Reloader/master/deployments/kubernetes/reloader.yaml
```

裝到你所需要的namespace即可, 然後在你的Deployment設定上加上一個annotation `reloader.stakater.com/auto: "true"`:

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: app
  annotations:
    reloader.stakater.com/auto: "true"
```

這樣reloader就會幫你監控這個Deployment用到相關的ConfigMap跟Secret, 不管是用環境變數的方式, 還是掛載檔案的方式, 都適用, 並且由於它是直接透過Kubernetes API, 因此ConfigMap或是Secret有變化都是即時會監測到, 然後它就會用rolling update的方式去重啟相關的instances, 相較之下會比用sidecar的方式保險