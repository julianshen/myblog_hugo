---
date: "2016-11-04T15:11:58Z"
tags:
- javascript
- ES6
title: '[ES6] Generators'
---

這個語法蠻有趣的, 早上一直都在看這個想能拿來幹嘛? 結果早上有個phone interview就有點腦袋小小的轉不過來,
不過這不是重點, 先簡單的來講一下[Generators](https://developer.mozilla.org/zh-TW/docs/Web/JavaScript/Reference/Statements/function*)

產生器? 顧名思義或許是這樣, 可以先看一下MDN上的[文件](https://developer.mozilla.org/zh-TW/docs/Web/JavaScript/Reference/Statements/function*),
這是一個從ES6開始才有的功能, 因此要在比較新的瀏覽器, 或是nodejs 6之後才有支援, 先來看看MDN上那個簡單的例子:

```javascript
function* idMaker(){
  var index = 0;
  while(index < 3)
    yield index++;
}

var gen = idMaker();

console.log(gen.next().value); // 0
console.log(gen.next().value); // 1
console.log(gen.next().value); // 2
console.log(gen.next().value); // undefined
``` 

一個generator是以"function*"宣告的, 可以說它是一種特別的function, 跟一般的function一次跑到結束不同,
它是可以中途暫停的, 以上面這例子來說, 如果習慣傳統的寫法, 可能會有點小暈, while loop在idMaker裡面不是一次做到完嗎?
剛剛說generator是可被中途暫停的, 因此, 在第一輪的時候會先停在"yield"處, 當你呼叫next()才會到下一個yield位置停下來並回傳,
我的感覺是比較像一種有帶state的function之類的

有什麼用途? 從上面的例子, 當然最直覺想到的是ID產生器或是計數器之類的東西, 網路上應該可以找到一些不同的用法, 比如說搭配Promise, 有興趣可以自己找找,
是不只可以用在產生器, 拿我早上interview被問到的實作strstr, 不是很難的東西, 我原本拿go寫,出了點小槌, 而且也只能找第一個發生的字串, 後來用generator改了這版本:

{{< gist julianshen 6c06ccfa0942829ea24973778a96ab64 >}}

以這個來說, 第一次呼叫會找出第一個發生的點, 可以持續呼叫到所有的都找出來為止, generator是可以被iterate的, 因此這邊可以用

```javascript
for(var i of gen) {
    console.log(i);
}
```

不需要一直call next(), next()回傳回來的會是{value:..., done:...}, 因此可以看done知道是否已經結束

下面這範例則是一個質數產生器, 一直call next就可以產生下一個質數:

{{< gist julianshen 0c283f6f76abf258f9c0d4292a1e14f9 >}}