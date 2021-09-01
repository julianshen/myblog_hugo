---
date: 2021-09-01T14:13:25+08:00
title: "在 Next.js 實作 Feature Toggle"
images: 
- "https://og.jln.co/jlns1/5ZyoIE5leHTvvI5qcyDlr6bkvZwgRmVhdHVyZSBUb2dnbGU"
---

之前一直有在思考, 如果把 Feature toggles或Feature flags帶入開發流程可以有甚麼幫助, 先撇開A/B Test不談(因為這還是要配合產品策略考量), 在一個開發Sprint週期時間越來越短下, 一個feature的完成通常也需要跨越多個sprints, 如果加上某些可能需要配合event時間推出的功能, 這對release也是會帶來一些挑戰, 理想化的狀況就是最新的程式碼一直都在, 測試沒問題後開個開關就可以打開, 甚或搭配一些技巧可以讓QA也能夠在生產環境驗證還沒release的功能

需要feature toggles的場合, 可能前後端都會有機會, 不過, 我想了一下, 後端還是可以從API版本或是其他方法隱藏未釋出功能, 而從前端一次隱藏整個功能區塊或是整個頁面, 或許是一個比較好開始切入的部分

先來介紹一下Feature toggle好了

## Feature Toggle

Feature Toggles (又名Feature Flags)是一個應該算是常見的軟體開發方式, 藉由開關旗標(flag)來開啟或隱藏程式中的功能, 土法一點, 你是可以自訂一個布林變數, 在release之前去打開或關閉它(true or false), 當然這樣彈性比較小, 理想上當然會希望可以彈性隨時開關

根據Martin Fowler這篇[FeatureToggle](https://martinfowler.com/bliki/FeatureToggle.html)以及Pete Hodgson先生寫的這篇: [Feature Toggles (aka Feature Flags)](https://www.martinfowler.com/articles/feature-toggles.html), Feature toggles根據用途可以分做不同種類(後面那篇介紹比較詳細):

1. Release Toggles: 這就如同前面講的, 開關或隱藏程式中功能, 或許是還未完成之功能, 可以搭配 trunk-based development服用, 通常Release toggles的變動頻率應該要是相當小的, 也可能等功能穩定後就會被拿掉
1. Experiment Toggles: 用在做A/B Test實驗, 會需要經過一些條件判斷來做啟用或禁用, 通常會根據request的資訊有變動
1. Ops Toggles: 用在跟系統運維相關, 比如說系統問題臨時需要下架, 需要上公告頁面, 或是像這種狀況需要下架某些功能(一個想法: 或許可以搭配系統監控使用)
1. Permission Toggles: 用在限定特定使用者可以使用, 或許直覺會想到是針對有權限控制的功能, 不過其實還可以用在Campaign/Event相關功能, 如果需要提早內部做dog fooding, 或許就可以考量beta user whitelist的作法, 搭配Permission Toggles

類別只是一種定義而已, 真正使用狀況或許會更複雜搭配, 也可能需要考慮flags之間的相依性跟連動關係, 不過這篇並不是要討論這部分, 這篇會以簡單的實現Release toggle來先做探討

關於feature toggles資源跟相關探討也可參考 [The Hub for Feature Flag Driven Development](https://featureflags.io/)

## Feature toggles的解決方案

有人提出方法, 就會有人提出解決方案和產品, Feature toggles相關的產品其實非常多, 商業產品有:

1. [LaunchDarkly](https://launchdarkly.com/about-us/)
1. [Unleash](https://www.getunleash.io/plans)
1. [HappyKit](https://happykit.dev/)

另外也可以找到不少Open source的

1. [Togglez](https://www.togglz.org/) (Java)
1. [FF4J](https://ff4j.github.io) (Java)
1. [Finagle](https://twitter.github.io/finagle/guide/Configuration.html#:~:text=Feature%20toggles%20are%20a%20commonly%20used%20mechanism%20for,rollout%20functionality%20in%20a%20measured%20and%20controlled%20manner.) (Finagle有內建, 不過不好用)
1. [Feature Flags API for Go ](https://github.com/AntoineAugusti/feature-flags) (Go)

你可以從[The Hub for Feature Flag Driven Development](https://featureflags.io/)找到一大堆

但這些絕大部分作法是把feature toggles的設定集中管理, 不管是放在資料庫, 或是有一個server來提供(後面也是放資料庫)

個人是覺得, Feature toggles並不是主角, 而是跑龍套的, 應該盡可能的輕量化, 做成獨立系統有點多了, 在Cloud native (k8s native)時代, 考慮用檔案來存設定(搭配ConfigMap), 加上lib, 應該是相對輕量的作法

因此針對這想法, 我在Next.js上做了幾個做法來驗證, 這邊把我的作法跟碰到的狀況做一個紀錄

## 基本設計

首先, 我先想了一下, 我在Next.js上要怎使用Feature toggles比較合適, 搭配tag應該是最為易讀的方式, 像是

```typescript
<div>
    <WithFeature feature="feature1">
        <div>This is new feature</div>
    </WithFeature>
    <Link href="/">HOME</Link>
</div>
```

用一個區塊包藏住需要開關的部分, 再由讀取到的設定做為開關

另一個想到的方式是:

```typescript
<PageWithFeature feature="feature3">
    <Link href="/">HOME</Link>
</PageWithFeature>
```

這種跟前一種的不同是, 如果功能未打開, 區塊就不會顯示, 但整頁的其他部分還會顯示, 但第二種應用在, 整頁都是新功能, 不想因為提早被發現URL而露出, 所以功能未打開的情況需要顯示404

另外config的設計希望是一個yaml檔案給程式去讀取, 像這樣:

```yaml
features:
  feature1: false
  feature2: true
  feature3: true
  offline: false
```

## 使用Custom App跟React context來做一個通用設計

為何說是通用設計? 這邊是希望每個需要用到toggle的頁面不用自已寫載入設定的部分(後面會講個範例是有需要的), 而是只要直接用 `<WithFeature feature="feature1">` 就好

Next.js有支援["Custom `App`"](https://nextjs.org/docs/advanced-features/custom-app), 所以我們載入config的部分可以直接放在"`_app.tsx`" 或 "`_app.jsx`"內就可以讓所有頁面共用(詳細請參考文件)載入部分的程式碼, 另外還得搭配[React context](https://reactjs.org/docs/context.html)才能順利的把設定下傳到component

首先, 我們需要先來設計一個`Context`用來傳遞設定給`WithFeature`的元件(Component)

```typescript
import { createContext, useState } from "react";

type FeatureMap = {
    [key:string]:boolean
}

export const FeatureContext = createContext({features: {} as FeatureMap})
```

這邊設計很單純, 假設每個feature都是true或false的開關, 前面的config yaml也是這樣的

在`_app.tsx`我們則需要有這東西:

```typescript
import { FeatureContext } from '../context/featurecontext'
function MyApp({ Component, pageProps }: AppProps) {  
  return <FeatureContext.Provider value={{features:pageProps.features}}><Component {...pageProps}/></FeatureContext.Provider>
}
```

這邊是把pageProps裡的features給帶到`FeatureContext`以讓後面的可以讀到, 那pageProps是又哪裡來的? 這邊我們就得去實作`app.getInitialProps`了, 這個在每個page初始化都會呼叫到它, 並把它產生的pageProps給餵到前面那個function

那是不是就可以把下面這段直接放到`getInitialProps`就可以取得feature flags的設定了呢?

```typescript
const configFile = path.join(process.cwd(), 'features.yaml')
const yaml = await fs.promises.readFile(configFile)
const config = YAML.parse(yaml.toString())
```

答案是"不行!!", 為何? 因為在Next.js的設計上:

1. `getInitialProps`是會有可能在client(瀏覽器), server跑
1. 除非你頁面實作上有`getServerSideProps`, `getInitialProps`才會總是在server端跑
1. 只有`getServerSideProps`,`getStaticProps`才可使用server端(node.js)API, 例如讀檔 (還有一個則是API, 但API實做不同, 我不規在此類)
1. `getServerSideProps`只設計給頁面的實作使用, 所以每頁有自己的 `getServerSideProps`, 但 App沒`getServerSideProps`

因為client跟server都會有機會呼叫到, 如果要讓它們能夠共用, 那就只有做一個API給它, 我們在`page/api`目錄底下開一個`features.tsx`的檔案, 內容是

```typescript
import { NextApiRequest, NextApiResponse } from "next"
import fs from 'fs'
import path from 'path'
import YAML from 'yaml'

type Features = {
    [key:string]:boolean
}

export default async function handler(
    req: NextApiRequest,
    res: NextApiResponse<Features>
  ) {
    const configFile = path.join(process.cwd(), 'features.yaml')
    const yaml = await fs.promises.readFile(configFile)
    const config = YAML.parse(yaml.toString())
    
    res.status(200).json(config.features as Features)
  }
```

那我們`getInitialProps`就可以這樣寫:

```typescript
MyApp.getInitialProps = async (appContext: AppContext) => {
   const appProps = await App.getInitialProps(appContext)
   const req = appContext.ctx.req
   var host = req
    ? req.headers["x-forwarded-host"] || req.headers["host"]
    : window.location.host;
   const resp = await fetch(`http://${host}/api/features`)
   const features = await resp.json()
   appProps.pageProps['features'] = features
   return { ...appProps }
 }
```

這邊有一點需要注意的, 雖然用`fetch`, 但它client端, server端用的是不一樣的, 雖然是長一樣的API, 所以這邊的URL是不可以用相對路徑的("`/api/features`"), 原因就是在server端只有相對路徑是無法知道實際要呼叫哪裡, 所以這邊還是組出一個完整的URL來使用(不過這實做有點不是很好, 需要改, 當PoC就算了 :p)

所以`WithFeature`就可以長這樣:

```typescript
import { ReactNode, useContext} from "react"
import { FeatureContext } from "../context/featurecontext"

export type FeaturesProps = {
    children?:ReactNode
    feature?: string
}

const WithFeature = function(props:FeaturesProps) {
    const {features} = useContext(FeatureContext)

    if(features && props.feature && features[props.feature]) {
        return (
            <>{props.children}</>)   
    }
    return <></>
}

export default WithFeature
```

因為`features`是放在`FeatureContext`中, 所以我們可以透過`useContext`來取值

這作法不算難, 而且也蠻好運用, 只要在每個需要的頁面自行運用`WithFeature`即可, 但缺點是甚麼? 基本上它算是CSR(Client side rendering), 而且要甚麼flag都是透過API跟server要, 會被看的一清二楚外, 或許可能還有方法偽造以至於你的功能提早被發現

## SSR (Server Side Rendering)的作法

Next.js強大的地方就是它有[支援SSR(Sever side rendering)和SSG(Server side generation)](https://nextjs.org/docs/basic-features/data-fetching), 兩者的好處就是在server端就把頁面內容給產生好

假設我們有一頁叫做`Post3`, 實作是這樣的:

```typescript
import { NextPage } from "next";
import { GetServerSideProps } from "next"
import fs from 'fs'
import path from 'path'
import Link from "next/link";
import YAML from 'yaml'
import WithFeature from "../components/feature2";
import PageWithFeature from "../components/feature3";

type Params = {}
type Props = {
    features: any
}

const Post3:NextPage<Props> = (props:Props) => {
    return (
        <FeatureContext.Provider value={{features:props.features}}>
            <PageWithFeature feature="feature3">
                <WithFeature feature="feature4">
                    this is new feature 3
                </WithFeature>                
                <Link href="/">HOME</Link>
            </PageWithFeature>
        </FeatureContext.Provider>
    )
}

export default Post3

export const getServerSideProps:GetServerSideProps<Props, Params> = async (ctx) =>{
    const configFile = path.join(process.cwd(), 'features.yaml')
    const yaml = await fs.promises.readFile(configFile)
    const config = YAML.parse(yaml.toString())

    return {
        props: {
            features: config.features
        }
    }
}
```

這頁面可以透過`/post3`來存取它, 這個頁面由於實做了`getServerSideProps`, 因此Next.js會使用SSR的模式, 這邊`getServerSideProps`就可以直接讀檔案了

但使用SSR這方法的話, 缺點就是

1. 每個頁面都得自己把讀config加入`getServerSideProps`
1. 也都需要自己把內容包在Context Provider如`<FeatureContext.Provider value={{features:props.features}}>`

就等於很多頁面都會有重複的程式碼, 不是那麼簡潔漂亮, 但優點應是內容在server端就已經先處理好了

## PageWithFeature

上面有偷偷藏一個`PageWithFeature`, 這個實做可以是這樣:

```typescript
import { ReactNode, useContext} from "react"
import Error from "next/error"
import { FeatureContext } from "../context/featurecontext"
import Offline from "./offline"

export type FeaturesProps = {
    features?:any
    children?:ReactNode
    feature?: string
}

const PageWithFeature = function(props:FeaturesProps) {
    const {features} = useContext(FeatureContext)

    if(features['offline']) {
        return <Offline></Offline>
    }
    
    if(features && props.feature && features[props.feature]) {
        return (
            <>{props.children}</>)   
    }
    
    return <Error statusCode={404} />
}

export default PageWithFeature
```

用途有兩個:

1. Disable時導到404頁面
1. 如果是系統下架狀態下, 導到一個暫時停止服務頁面(這邊用另一個component來解決)

## 開發流程上的思考

Atlassian這篇介紹[Feature flags](https://www.atlassian.com/continuous-delivery/principles/feature-flags)裡有張圖, 我覺得蠻有趣的, 可以做為開發流程上的一個參考:
![](https://wac-cdn.atlassian.com/dam/jcr:cc55dc99-ddf4-4e61-a989-db6446cfef3c/Feature%20Flag-Driven%20Development@2x.png?cdnVersion=1785)

Facebook這篇[Rapid release at massive scale](https://engineering.fb.com/2017/08/31/web/rapid-release-at-massive-scale/), 雖然跟Feature flags關係比較不大, 但也蠻有參考價值的