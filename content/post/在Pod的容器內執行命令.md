---
title: "在Pod的容器內執行命令"
date: 2021-12-08T10:13:18+08:00
images: 
- "https://og.jln.co/jlns1/5ZyoUG9k55qE5a655Zmo5YWn5Z-36KGM5ZG95Luk"
---

一般來說, container通常會設計成只專注在它單一的任務上, 也就是通常不會把一個http server跟db server跑在同一個container內, Kubernetes 的Pod的設計, 讓我們可以在同一個Pod內放多個containers, 因此可以延伸出init container, sidecar container來輔助原本的container, 中間可以透過分享Volume或是直接透過loopback網路來共享資料, 但還是會有情境是, 你會希望可以在某個container空間內執行某個程式, docker的話, 你可以用 `docker exec` ,那在Kubernetes呢?

Kubernetes 也是可以用`kubectl` 達成同樣的目的, 像是:

```shell
 kubectl exec mycontainer ls
```

其實跟`docker exec` 是類似的, 如果是要執行shell進去執行其他的維護:

```shell
kubectl exec mycontainer -it /bin/sh
```

但如果是, 你要從其他的pod去執行其他的pod裡面的指令呢? 像是用Cron job定期去執行某特定pod裡面的程式?

我們也是可以透過呼叫Kubernetes API來達成這目的的, 如同下面這個範例:

```go
func runCommand(clientset *kubernetes.Clientset, pod string, ns string, cmd ...string) (string, error) {
	req := clientset.CoreV1().RESTClient().Post().Resource("pods").Name(pod).Namespace(ns).SubResource("exec")
	option := &v1.PodExecOptions{
		Command: cmd,
		Stdin:   false,
		Stdout:  true,
		Stderr:  true,
		TTY:     false,
	}

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	req.VersionedParams(option, scheme.ParameterCodec)

	exec, err := remotecommand.NewSPDYExecutor(config, "POST", req.URL())
	if err != nil {
		return "", err
	}

	sopt := remotecommand.StreamOptions{
		Stdout: &stdout,
		Stderr: &stderr,
	}

	err = exec.Stream(sopt)

	if err != nil {
		return "", err
	}

	if stderr.String() != "" {
		return stdout.String(), errors.New(stderr.String())
	}

	return stdout.String(), nil
}
```

由於Kubernetes API還沒包裝`exec`這個資源, 所以要用`.Resource("pods").Name(pod).Namespace(ns).SubResource("exec")`去取用, stdin, stdout, stderr都還是可以串接回來的

但如果你直接在你的pod內執行這段的話, 其實不會成功的, 因為你沒有權限可以去做, 你必須先設定好一個Role如下:

```yaml
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: executor
rules:
  - verbs:
      - get
      - list
    apiGroups:
      - ''
    resources:
      - pods
      - pods/log
  - verbs:
      - create
    apiGroups:
      - ''
    resources:
      - pods/exec

```

我的情境是會用list pods找出pod的名稱, 在用這個名稱去找到特定的pod執行, 所以會需要前半段get和list的權限, 後半段對`pods/exec`做`create`的權限才是執行這段程式真正需要的