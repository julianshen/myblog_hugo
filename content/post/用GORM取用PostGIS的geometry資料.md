---
date: 2021-10-29T23:58:47+08:00
title: "用GORM取用PostGIS的geometry資料"
images: 
- "https://og.jln.co/jlns1/55SoR09STeWPlueUqFBvc3RHSVPnmoRnZW9tZXRyeeizh-aWmQ"
---

[PostGIS](https://postgis.net/) 是讓PostgresSQL可以支援地理資訊資料的一個擴充, 而[geometry](https://www.postgis.net/workshops/postgis-intro/geometries.html)是[PostGIS](https://postgis.net/)定義來儲存地理資料的資料型態, 包含了座標點(POINT), 線, 多邊形等等

而[GORM](https://gorm.io/)則是Golang界算蠻有名的ORM套件, 支援了蠻多不同的關聯式資料庫

不過[GORM](https://gorm.io/)是沒有直接支援PostGIS的geometry這個資料型態的, 畢竟geometry並非一般SQL標準的資料型別, 因此如果要讓[GORM](https://gorm.io/)可以來存取geometry, 就得使用它所提供的[Customize Data Type](https://gorm.io/docs/data_types.html)

要存取這種自訂的資料型別, 需要實做`Scanner`和`Valuer`這兩個介面:

```go
type Scanner interface {
	// Scan assigns a value from a database driver.
	//
	// The src value will be of one of the following types:
	//
	//    int64
	//    float64
	//    bool
	//    []byte
	//    string
	//    time.Time
	//    nil - for NULL values
	//
	// An error should be returned if the value cannot be stored
	// without loss of information.
	//
	// Reference types such as []byte are only valid until the next call to Scan
	// and should not be retained. Their underlying memory is owned by the driver.
	// If retention is necessary, copy their values before the next call to Scan.
	Scan(src interface{}) error
}

type Valuer interface {
	// Value returns a driver Value.
	// Value must not panic.
	Value() (Value, error)
}
```

那在實做前, 可以先來看看geometry實際存到資料庫是長怎麼樣? 它長的就會像是一串這樣的文字`0101000020E61000001E64C73CEA905EC09CD6C962A7C84240`, 看起來就是十六進位編碼過的文字, 根據[文件](http://postgis.net/docs/using_postgis_dbmanagement.html#PostGIS_Geography), 它是經過EWKB編碼過的(EWKB是Postgis從[WKB](https://en.wikipedia.org/wiki/Well-known_text_representation_of_geometry)延伸來的), 然後編碼過的binary資料再轉成Hex字串

這部分的解碼, 可以不用自己寫, 使用open source的好處就是可以踩在前人的肩膀上, 可以利用[go-goem](https://github.com/twpayne/go-geom)來幫我們處理這件事, 這邊以處理座標點為例 (採用常見的[EPSG:4326](https://spatialreference.org/ref/epsg/4326/)座標系統), 定義一個叫做`GeoPoint`的型別給GORM使用, 而這個型別的實作可以是這樣:

```go
type GeoPoint geom.Point

func (g *GeoPoint) Scan(val interface{}) error {
	pt, err := ewkbhex.Decode(val.(string))

	if err == nil {
		if p, ok := pt.(*geom.Point); ok {
			*g = GeoPoint(*p)
		} else {
			return errors.New(fmt.Sprint("Failed to unmarshal geometry:", val))
		}
	}

	return err
}

func (g GeoPoint) Value() (driver.Value, error) {
	pt := &g
	toPt, err := ewkbhex.Encode((*geom.Point)(pt).SetSRID(4326), binary.BigEndian)

	return toPt, err
}

```

其實非常簡單的利用了[go-goem](https://github.com/twpayne/go-geom)提供的`ewkbhex`來編解碼而已