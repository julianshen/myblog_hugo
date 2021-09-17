---
date: 2021-09-17T12:20:38+08:00
title: "用D3js建立台股板塊圖"
images: 
- "https://og.jln.co/jlns1/55SoRDNqc-W7uueri-WPsOiCoeadv-WhiuWclg"
---

這題目瑣碎的東西太多了, 所以這篇打算只是做個紀錄, 做這東西原因是看到[Finviz](https://finviz.com/map.ashx?t=sec_all)這個板塊圖覺得還蠻有趣的, 想說該怎去做到這樣的圖表

![stockmap](/images/posts/2021-09-17-12-25-23.png)

一眼就可以大略看市場的狀況, 感覺還蠻酷的, 查了一下, 這東西叫[Treemapping](https://en.wikipedia.org/wiki/Treemapping), 想到資料視覺化, 我是先想到[D3.js](https://d3js.org/), 雖然說[Highcharts](https://www.highcharts.com/)也可以達到一樣的目的, 不過[D3.js](https://d3js.org/)使用上跟jQuery類似, 比較簡單, 所以選擇它來實現

先給結果[https://fviz.jln.co/marketmap](https://fviz.jln.co/marketmap)

![f](/images/posts/2021-09-17-12-45-59.png)

這邊使用到的東西有:

* D3.js
* Next.js (deploy to GitHub page)

純粹靜態網頁, 沒資料庫, 不過目前資料只抓到2021/09/14, 定期抓資料的部分還懶得弄

## 資料來源

這邊所需要的資料有幾個:

* 上市股票的收盤資訊
* 上櫃股票的收盤資訊
* 各股所屬的分類資訊

雖然都是容易爬的到的資料, 但兩個市場資料格式不是那麼的統一

抓取集中市場的歷史資料用這個 URL : "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=json&date=%s&type=ALL&_=%s" , date的格式用"20060102"這樣, "_"可以用timestamp即可, 我要的當然是json, 會比csv來的好處理點

個股的交易資訊在`data9`這個欄位, 欄位定義是`fields9`, 所以用`jq`來看一下

```
curl "https://www.twse.com.tw/exchangeReport/MI_INDEX?response=json&date=20210910&type=ALL&_=1631369120214" | jq ".fields9"
  % Total    % Received % Xferd  Average Speed   Time    Time     Time  Current
                                 Dload  Upload   Total   Spent    Left  Speed
100 3122k    0 3122k    0     0  1526k      0 --:--:--  0:00:02 --:--:-- 1526k
[
  "證券代號",
  "證券名稱",
  "成交股數",
  "成交筆數",
  "成交金額",
  "開盤價",
  "最高價",
  "最低價",
  "收盤價",
  "漲跌(+/-)",
  "漲跌價差",
  "最後揭示買價",
  "最後揭示買量",
  "最後揭示賣價",
  "最後揭示賣量",
  "本益比"
]
```

來看一下`.data9`內的單筆資料, 基本上就是都放到array去, 算好處理

```json
[
    "9944",
    "新麗",
    "92,813",
    "53",
    "1,950,673",
    "21.05",
    "21.10",
    "20.85",
    "21.10",
    "<p style= color:red>+</p>",
    "0.10",
    "20.75",
    "2",
    "21.15",
    "3",
    "23.19"
  ]
```

這邊"漲跌(+/-)"的部分, 其實不是只有+和-, 而居然是html tags, 包含四種狀況, +/-/ /X, +/-很好懂, 就是漲跟跌, 空白就是平盤了, X的狀況通常發生在除權息, 增減資這類狀況

那櫃檯市場呢? URL是這個`https://www.tpex.org.tw/web/stock/aftertrading/daily_close_quotes/stk_quote_result.php?l=zh-tw&d=110/09/01&_=1631603049`, 日期格式跟集中市場不同, 是用"/"隔開, 並且是民國紀年, 這邊資料也是array放在`aaData`這欄位

```json
["9960","\u9081\u9054\u5eb7","27.10","+0.35","27.00","27.25","26.85","27.06","56,004","1,515,312","37","27.10","1","27.15","5","33,592,500","27.10","29.80","24.40"]
```

那個股基本資料呢?這邊就神奇了, 居然有Open API document: [https://openapi.twse.com.tw/v1/swagger.json](https://openapi.twse.com.tw/v1/swagger.json), 可以用"/v1/opendata/t187ap03_L"取得基本資料, 這邊雖然也有API可以取得當日交易資訊, 但只有當日並無歷史資料

上櫃股票的資料也有一樣的東西, 在 [https://www.tpex.org.tw/openapi/swagger.json](https://www.tpex.org.tw/openapi/swagger.json)

抓到的分類類別是代號, 所以要對應到正確的類別名稱可以用這表:

```golang
var Categories = map[string]string{
	"01": "水泥",
	"02": "食品",
	"03": "塑膠",
	"04": "紡織纖維",
	"05": "電機機械",
	"06": "電器電纜",
	"21": "化工",
	"22": "生技",
	"08": "玻璃陶瓷",
	"09": "造紙",
	"10": "鋼鐵",
	"11": "橡膠",
	"12": "汽車",
	"24": "半導體",
	"25": "電腦及週邊",
	"26": "光電",
	"27": "通信網路",
	"28": "電子零組件",
	"29": "電子通路",
	"30": "資訊服務",
	"31": "其他電子",
	"14": "建材營造",
	"15": "航運",
	"16": "觀光",
	"17": "金融保險",
	"18": "貿易百貨",
	"23": "油電燃氣",
	"19": "綜合",
	"20": "其他",
	"32": "文創",
	"33": "農業科技",
	"34": "電商",
	"80": "管理股票",
	"91": "存託憑證",
}
```

把以上資料, 整合起來, 我需要的是這樣的資料:

```json
{
    "name": "台股版塊",
    "children": [{
        "name": "集中市場",
        "children": [{
            "name": "水泥",
            "children": [{
                "name": "1101",
                "data": {
                    "Code": "1101",
                    "Name": "台泥",
                    "TradeVolume": "14853294",
                    "Transaction": "6367",
                    "TradeValue": "715327703",
                    "OpeningPrice": "48.35",
                    "HighestPrice": "48.40",
                    "LowestPrice": "47.85",
                    "ClosingPrice": "48.40",
                    "Change": "-0.05",
                    "Time": "2021-09-01T00:00:00+08:00"
                }
            }]
        }]
    }]
}
```

顧名思義, Treemap就是一個樹狀的結構而來的, 因此需要的資料結構就需要有個階層, 這邊設計成 "台股板塊->市場別->分類->個股"

## D3.js + React JS

因為我用Next.js, 就想要把這個treemap包裝在一個react component

使用D3.js要先引入這幾個packages (我用typescript開發):

* @types/d3
* @types/d3-hierarchy
* d3
* d3-hierarchy

`d3-hierarchy`是用來畫treemap的, 只用d3基本功能是不需要含進來的

先給這個Treemap的component一個殼:

```typescript
const Treemap = (props: { width:number, height:number, date:Date }) => {
    const svgRef = useRef(null);
    
    const dataFile = "data/" + props.date.getFullYear() + "-" 
      + (props.date.getMonth() + 1).toString().padStart(2, "0") 
      + "-" + props.date.getDate().toString().padStart(2, "0") + ".json";

    const renderTreemap = async () => {
        const svg = d3.select(svgRef.current).style("font", "10px sans-serif");
        svg.attr('width', props.width).attr('height', props.height);
        svg.selectAll("*").remove();
        
        var stockData:StockData

        try {
            stockData = await d3.json(dataFile) as StockData;
        } catch(e) {
            svg.append("text")
                .text("本日無資料, 請按左上角按鈕選取時間")
                .attr("x", 6)
                .attr("y", 22)
                .attr("stroke", "white");
            return;
        }
    };

    useEffect(() => {
        renderTreemap();
    });
    
    return (
        <Box>
            <svg ref={svgRef} />
        </Box>
      );
}

export default Treemap
```

d3的用法跟jQuery類似, 因此這邊跟React包裝一起的方法也很簡單, 就是用`useRef`給它有個reference可以select, 這邊要畫圖, 所以就包到一個svg去, 實際render出來的也是svg

抓取資料可以用`d3.json(dataFile)`, 其實也有`d3.csv`, 有點類似用`fetch`

Treemap的實做就稍微有點複雜了, 可以參考這邊 "[Nested Treemap](https://observablehq.com/@d3/nested-treemap)", 這篇也不錯: "[D3.js 實戰 － 利用 Treemap Layout 將政府預算視覺化](http://blog.infographics.tw/2015/10/d3js-tutorial-treemap-and-budget/)", 這邊我使用了交易量跟交易總額去做面積跟排序

## 加上Tooltip

```typescript
svg.selectAll("rect")
    .data(root.leaves())
    .enter()
    .append("rect")
    .attr('x', d => { return (d as HierarchyRectangularNode<StockData>).x0; })
    .attr('y', d => { return (d as HierarchyRectangularNode<StockData>).y0; })
    .attr('width', d => { return (d as HierarchyRectangularNode<StockData>).x1 - (d as HierarchyRectangularNode<StockData>).x0; })
    .attr('height', d => { 
            const h = (d as HierarchyRectangularNode<StockData>).y1 - (d as HierarchyRectangularNode<StockData>).y0; 
            return h;
        })
    .style("stroke", "black")
    .style("fill", d => ccolor(d.data))
    .on('mouseover', (event, dataNode)=>{
        mouseOver(event, dataNode);
    }).on('mouseleave', () => {
        tooltip().style("opacity", 0);
    });
```

透過`on('mouseover')`和`on('mouseleave')`就可以來加上tooltip的效果

## 發佈到Github pages

next.js由於ssg, ssr的關係, 需要跑個server, 但其實也有機會發布成全靜態網頁(只要沒需要有在server跑的部分), 步驟如下:

* next build
* next export

結果就會在`out/`目錄, 拿這個目錄的內容放github pages就可以了

最後附上寫這東西時來搗亂的傢伙

![](/images/posts/2021-09-17-13-43-04.png)