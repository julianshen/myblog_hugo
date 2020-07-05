---
date: "2020-07-05T14:54:23+08:00"
images:
- /images/posts/用golang解析json裡的數字欄位.md.jpg
title: 用golang解析json裡的數字欄位
---

在寫網路相關應用的時候, 應該常常會碰到需要去解析JSON格式的資料, 而Go在這邊也已經內建了一個蠻方便的套件 - **"encoding/json"** 可以讓我們輕易地來處理這類型的資料

先來看看下面這段範例:

```golang
package main

import (
	"fmt"
	"encoding/json"
)

type Sample struct {
   Name string
   Age int 
}

const sampleData = `{
   "name":"julian",
   "age": -1
}`

func main() {
	var sample Sample
	err := json.Unmarshal([]byte(sampleData), &sample)
	
	if err != nil {
	   panic(err)
	}
	fmt.Println(sample)
}
```
[[執行](https://play.golang.org/p/YfOPBMXGGl5)]

從這個範例可以看到, 我們可以用很簡單的程式碼, 把下面這段JSON內容給對應到```Sample```這個結構裡面

```json
{
   "name":"julian",
   "age": -1
}
```

但這邊其實有一個問題, 如果你把這個JSON資料, 改成下面這樣子:

```json
{
   "name":"julian",
   "age": "-1"
}
```

這在現實世界應該蠻常看到的, 只是多加個雙引號而已, 大家應該也會預期這邊應該也會沒問題的解析出一樣的結果吧? 但你如果實際改了資料[執行看看](https://play.golang.org/p/nxMnQFY3CQ-), 你得到的結果應該會是:

```sh
panic: json: cannot unmarshal string into Go struct field Sample.Age of type int

goroutine 1 [running]:
main.main()
	/tmp/sandbox247755092/prog.go:23 +0x162
```

這其實是 ***"encoding/json"*** 把雙引號的內容都當作字串來看, 所以當我要把它塞到一個 *int* 欄位時, 就會出事了

解決方法有好幾種, 下面就一一來看看:

## 用Field tag來解決

如果把```Sample```的定義改成下面這樣:

```golang
type Sample struct {
   Name string
   Age int `json:",string"`
}
```
[[執行範例](https://play.golang.org/p/qexEcHzamM3)]

喔耶, 沒問題了耶!!!可以正常的解出資料了耶!! 慢著, 先別高興太早!! 再試試把雙引號拿掉看看(參考這個[範例](https://play.golang.org/p/Cz17VTbCBBX))

呃, 是怎樣啦!! 換成這個錯誤了!! 

```sh
panic: json: invalid use of ,string struct tag, trying to unmarshal unquoted value into int

goroutine 1 [running]:
main.main()
	/tmp/sandbox912379425/prog.go:23 +0x162
```

在現實案例中, 的確是有可能碰到有時送來的資料有雙引號, 有時候沒有, 這方法是沒法一次滿足的

## 利用Unmarshaler界面來解決

***"encoding/json"*** 是可以讓開發者自行指定怎去解析JSON內容的, 只需要定義一個自定義的型別並實做Unmarshaler界面就可以了, 為了解決這個問題, 我們可以定義一個新的```MyInt```的型別, 並幫它實做```UnmarshalJSON```的方法, 參考下面範例:

```golang
type MyInt int

func (m *MyInt) UnmarshalJSON(data []byte) error {
	str := string(data)
	
	if unquoted, err := strconv.Unquote(str); err == nil {
	   str = unquoted
	}
	
	result, err := strconv.Atoi(str)
	if err != nil {
		return err
	}
	*m = MyInt(result)
	return nil
}

type Sample struct {
	Name string
	Age  MyInt
}
```
[[執行範例](https://play.golang.org/p/FAqSuBcDXgl)]

在Sample這個結構中的Age欄位, 從原本的int改成MyInt, 這樣```json.Unmarshal```碰到Age這欄位的話, 就會用```MyInt```的```UnmarshalJSON```方法去解析

這方法是麻煩了點, 而且可能針對不同型別要去個別做, 但卻是可以同時處理掉前述兩種型態的資料, 程式沒那好看就是了

## 使用[json.Number](https://golang.org/pkg/encoding/json/#Number)

在 ***"encoding/json"*** 裡其實還提供一個資料型態[json.Number](https://golang.org/pkg/encoding/json/#Number), 這個應該是這個問題的比較正規的解法了, 把Sample的定義改成下面這樣:

```golang
type Sample struct {
	Name string
	Age  json.Number
}
```
[[執行範例](https://play.golang.org/p/sWbmNX1aZwi)]

這方法也是可以無誤的解析兩種的型態, 然後當你要取用```Age```這欄位的```int```型態時, 你可以用```sample.Age.Int64()```去取得, 要多一層是麻煩點

那它是怎做到的呢? 來看一下它的原始碼:

```golang
// A Number represents a JSON number literal.
type Number string

// String returns the literal text of the number.
func (n Number) String() string { return string(n) }

// Float64 returns the number as a float64.
func (n Number) Float64() (float64, error) {
	return strconv.ParseFloat(string(n), 64)
}

// Int64 returns the number as an int64.
func (n Number) Int64() (int64, error) {
	return strconv.ParseInt(string(n), 10, 64)
}
```

會發現, 其實沒啥了不起的, Number本身就是一個string, 所以它是當字串在處理, 而當你需要數字表示(int或float)時, 再呼叫ParseInt或ParseFloat來處理

## 使用 [Jsoniter](https://github.com/json-iterator/go) 來取代 "encoding/json"

[Jsoniter](https://github.com/json-iterator/go)是一個號稱比原生的***"encoding/json"***效能還要來的更好的JSON處理套件, 除了Go的版本外, 它也有Java的版本

效能不在這邊的討論範圍, 但除效能外, 它也是最簡單解決這問題的方法, 先來看看完整的程式碼吧:

```golang
package main

import (
	"fmt"
	jsoniter "github.com/json-iterator/go"
	extra "github.com/json-iterator/go/extra"
)

type Sample struct {
	Name string
	Age  int
}

const sampleData = `{
   "name":"julian",
   "age": 45
}`

const sampleData2 = `{
   "name":"julian",
   "age": "45"
}`

func main() {
	extra.RegisterFuzzyDecoders()
	var sample Sample
	err := jsoniter.Unmarshal([]byte(sampleData), &sample)

	if err != nil {
		panic(err)
	}
	fmt.Println(sample)

	err = jsoniter.Unmarshal([]byte(sampleData2), &sample)

	if err != nil {
		panic(err)
	}
	fmt.Println(sample)

}
```
[[執行範例](https://play.golang.org/p/vKl1YcH6lt8)]

為啥說是最簡單的方法呢? 首先, 它有跟 ***"encoding/json"*** 完全一樣的使用方法, 把套件換掉後, 程式幾乎一樣, 所以也不需要改啥程式, 但如果只有改這樣, 會發現問題都還是在, 並沒有解決掉, 這時候就要帶入它額外的功能 ```FuzzyDecoders``` 了, ```FuzzyDecoders```是在它額外的套件裡面, 所以只需要加入import就可以用了

```golang
import extra "github.com/json-iterator/go/extra"
```

然後在開始使用前, 註冊 ```FuzzyDecoders``` 即可(```extra.RegisterFuzzyDecoders()```), 這樣, 不管有沒雙引號都不會有問題了!!
