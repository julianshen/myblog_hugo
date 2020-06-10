---
date: "2017-04-01T10:00:58Z"
tags:
- golang
title: '[Golang] Go的httptest'
---

在標準的go package中除了已經內建了http相關的實作外, 還有一個`net/http/httptest`, 這的package是用來給寫http相關測試用的, 可分為測試http server (http handler)和http client的(提供mock server給client)

如果要測試http handler, 所需要的是ResponseRecorder, 基本上相關的也只需要兩個方法`NewRequest`, `NewRecorder`, 參照下面範例:

```go
package main

import (
        "fmt"
        "github.com/stretchr/testify/assert"
        "io"
        "io/ioutil"
        "net/http"
        "net/http/httptest"
        "testing"
)

func TestHttpServer(t *testing.T) {
        assert := assert.New(t)
        handler := func(w http.ResponseWriter, r *http.Request) {
                io.WriteString(w, "<html><body>Hello World!</body></html>")
        }

        req := httptest.NewRequest("GET", "http://example.com/foo", nil)
        w := httptest.NewRecorder()
        handler(w, req)

        resp := w.Result()
        body, _ := ioutil.ReadAll(resp.Body)

        assert.Equal(200, resp.StatusCode)
        assert.Equal("text/html; charset=utf-8", resp.Header.Get("Content-Type"))
        fmt.Println(string(body))
}
```

範例中先用`NewRequest`假造出一個連接到`http://example.com/foo`的request, 而對於一個http handler來說的話, 所需要的參數一個是*Request, 另一個則是ResponseWriter了, ResponseWriter便可以由`NewRecorder`假冒, 再將這兩者傳給handler去處理, 並可以由`ResponseRecorder.Result`取得回傳內容來驗證

那如果是要測試http client呢?測client通常我們會需要mock server來偽造成真著server餵給client適當的假資料來測試, 以下是一個測試的範例:

```go
func TestServer(t *testing.T) {
        assert := assert.New(t)
        ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
                fmt.Fprintln(w, "Hello, client")
        }))
        defer ts.Close()

        res, err := http.Get(ts.URL)
        if err != nil {
                t.Fatal(err)
        }
        greeting, err := ioutil.ReadAll(res.Body)
        res.Body.Close()
        if err != nil {
                t.Fatal(err)
        }

        assert.Equal("Hello, client\n", string(greeting))
}
```

利用`httptest.NewServer`創建一個test server, 這邊的handler就看你需要利用什麼樣的假資料測試來做, 上面這例子只是用單純的http client來測, 回傳總是是"Hello, client", 但假設你是測試restful API, 那也可以準備一系列的JSON回傳, 你只要把`ts.URL`當作你的API endpoint給你的REST client的實作使用即可