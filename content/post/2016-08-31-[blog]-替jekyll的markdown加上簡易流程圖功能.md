---
date: "2016-08-31T16:33:15Z"
tags:
- Jekyll
- Blog
title: '[Blog] 替Jekyll的markdown加上簡易流程圖功能'
---
對一個developer的blog來說, 流程圖似乎是蠻需要的, 比較能夠清楚來解釋一些東西, 但每個東西都轉圖檔還蠻麻煩的, 下面介紹一個有用的Jekyll plugin, 可以做到像下面這樣的效果:

**第一例**

{{<mermaid>}}
graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
{{</mermaid>}}

**第二例**

{{<mermaid>}}
sequenceDiagram
    participant John
    participant Alice
    Alice->>John: Hello John, how are you?
    John-->>Alice: Great!
{{</mermaid>}}

這是利用一個叫做[Jekyll-mermaid](https://github.com/jasonbellamy/jekyll-mermaid) 來達成的

而這plugin其實也沒做很多事, 它是包裝了[mermaid](https://github.com/knsv/mermaid)這個工具, 而mermaid這工具他是利用了[ds.js](https://d3js.org)來讓你用很簡單的方式來繪製流程, 以上面兩個例子為例

**第一例**

```markdown
graph TD;
    A-->B;
    A-->C;
    B-->D;
    C-->D;
```

**第二例**

```markdown
sequenceDiagram
    participant John
    participant Alice
    Alice->>John: Hello John, how are you?
    John-->>Alice: Great!
```

所以你在markdown裡面只要加上

```markdown
{ % mermaid % }
sequenceDiagram
    participant John
    participant Alice
    Alice->>John: Hello John, how are you?
    John-->>Alice: Great!
{ % endmermaid % }
```

他就會幫你render出相關的流程了

#### 安裝方法 ####
這邊以我自己blog的安裝方法來說明

1. 把jekyll-mermaid.rb放到_plugins目錄去
1. 在_config.yml加上 (這邊以6.0.0的mermaid為例):

```markdown
mermaid:
  src: 'https://cdn.rawgit.com/knsv/mermaid/6.0.0/dist/mermaid.js'
```
1. 還要在head.html加上css (要配合版面顏色), 可以用這個 : https://cdn.rawgit.com/knsv/mermaid/6.0.0/dist/mermaid.css
