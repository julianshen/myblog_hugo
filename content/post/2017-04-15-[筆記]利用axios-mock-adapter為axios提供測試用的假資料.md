---
date: "2017-04-15T12:46:27Z"
images:
- /images/posts/2017-04-15-[筆記]利用axios-mock-adapter為axios提供測試用的假資料.md.jpg
tags:
- javascript
- axios
- rest
title: '[筆記]利用axios-mock-adapter為axios提供測試用的假資料'
---

[axios](https://github.com/mzabriskie/axios)是蠻好用的javascript http client, 不僅可以在browser上跑, 也可以在node.js上用, 而且Promise形態的API寫起來就比較好看, 如果搭配async/await的寫法, 看起來就更加漂亮了

```javascript
function loadUser(uid) {
  axios.get('/user?ID=12345')
    .then(response => {
      console.log(response)
    })
    .catch(error => {
      console.log(error)
    })
}
```

或是（async/await)

```javascript
async function loadUser(uid) {
  try {
    data = await axios.get('/user?ID=12345')
	console.log(data)
  } catch(e) {
    console.log(e)
  }
}
```

但如果開發時期或是要做Unit testing需要用假資料來取代server api直接回傳呢? 目前我看到兩套方案, 一個是[axios](https://github.com/mzabriskie/axios)作者做的[moxios](https://github.com/mzabriskie/moxios)另外一個是[axios-mock-adapter](https://github.com/ctimmerm/axios-mock-adapter), [moxios](https://github.com/mzabriskie/moxios)看起來好像比較適合在Unit testing時, 而我是想在開發過程中使用, 所以我選的是[axios-mock-adapter](https://github.com/ctimmerm/axios-mock-adapter)

使用[axios-mock-adapter](https://github.com/ctimmerm/axios-mock-adapter)還蠻簡單的:

```javascript
import axios from 'axios'
import MockAdapter from 'axios-mock-adapter'

let mock = new MockAdapter(axios)

mock.onGet('/users').reply(200, {
  users: [
    { id: 1, name: 'John Smith' },
    { id: 2, name: 'John Doe' }
  ]
})

axios.get('/users')
  .then(response => {
    console.log(response.data)
  })
```

創建uri跟假資料的對應很簡單, 基本上也就是'on''Method', 比如說`onGet`, `onPost`, 另外還有一個`onAny`可以處理所有的HTTP methods

做過mock後, axios呼叫這個uri所拿回來的資料通通就都會是假資料了, 這樣也不用為了塞假資料開發測試而去改動自己的程式