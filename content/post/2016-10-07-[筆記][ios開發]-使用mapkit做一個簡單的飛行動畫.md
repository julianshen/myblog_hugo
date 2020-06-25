---
date: "2016-10-07T23:10:04Z"
images:
- /images/posts/2016-10-07-[筆記][ios開發]-使用mapkit做一個簡單的飛行動畫.md.jpg
tags:
- iOS
- Swift
- MapKit
- mobiledev
title: '[筆記][ios開發] 使用MapKit做一個簡單的飛行動畫'
---

小時候很喜歡Indiana Jones系列的電影, 對於它裡面的地圖片段也一直覺得很有趣

{{< youtube FhoPTyMUgX0 >}}

如果這樣的動畫, 用在遊記類的blog上, 應該也蠻酷的, 但好像也沒一個比較好的工具, 因此想說用MapKit來實作一個試試

### 功能需求

先來定義一個簡單的功能需求

1. 在起點跟終點畫一條連結線
1. 一架飛機延這條線飛到終點
1. 地圖視角跟著飛機走

### 實作

#### 在起點跟終點畫一條連結線

這部份要用到MKGeodesicPolyline, 給它兩個點, 它就會自動連結成一條線, 但這條線並不是完美的直線, 因為地球表面是曲面的, 所以它是一條弧線

```swift
let coords = [start, end]

geodesicPolyline = MKGeodesicPolyline(coordinates: coords, count: 2)

print(geodesicPolyline!.pointCount)
mapView.add(geodesicPolyline!)
``` 

這邊coors只要給訂起始點跟終點的位置就好, 印出pointCount就會發現它把經由的點都補足了(實際上印出來的會多出2很多很多)

MKGeodesicPolyline裡的參數coordinates是一個UnsafeMutablePointer, 在Swift 3之前要寫成&coords, 但在Swift 3大改之後, "&"就不需要了

由於MKGeodesicPolyline是一個Overlay, 因此最後只需要用mapView.add (Swift 3之前是addOverlay)加入mapView就可以了, 但加入之後, 會發現, 這條線根本沒被畫出來, 那是因為少寫了一部分

在MapView裡面要畫出Overlay, 就必須要跟MapView說怎麼畫出這個Overlay, 這就要實作MKMapViewDelegate裡的mapView(_ mapView: MKMapView, rendererFor overlay: MKOverlay)

```swift
func mapView(_ mapView: MKMapView, rendererFor overlay: MKOverlay) -> MKOverlayRenderer {
    guard let polyline = overlay as? MKPolyline else {
        return MKOverlayRenderer()
    }
    
    let renderer = MKPolylineRenderer(polyline: polyline)
    renderer.lineWidth = 3.0
    renderer.alpha = 0.5
    renderer.strokeColor = UIColor.blue
    
    return renderer
}
```

由於我們是要畫Poly line, 因此這邊回傳給它一個MKPolylineRenderer, 線寬是3.0, 線的顏色是藍色(alpha = 0.5)

這樣就可以很完美的畫出那條線了

#### 在地圖上畫出飛機

這部份就要借重到MKPointAnnotation

```swift
let thePlane = MKPointAnnotation()
thePlane.coordinate = start //給定起始座標
        
mapView.addAnnotation(thePlane)
```

一樣, 這邊如果沒告訴MapView怎畫這Annotation, 它是會用預設的取代, 因此我們一樣要去實作MKMapViewDelegate

```swift
func mapView(_ mapView: MKMapView, viewFor annotation: MKAnnotation) -> MKAnnotationView? {
    let planeIdentifier = "Plane"
    
    let annotationView = mapView.dequeueReusableAnnotationView(withIdentifier: planeIdentifier)
        ?? MKAnnotationView(annotation: annotation, reuseIdentifier: planeIdentifier)
    
    annotationView.image = UIImage(named: "ic_flight_48pt")
    
    return annotationView
}
```

這邊用一個UIImage來指定飛機的圖標

#### 地圖視角

把飛機置於地圖正中央, 我們才看得到他, 因此, 需要設定可視的區域, 包含中心點跟範圍, 如下:

```swift
let span = MKCoordinateSpanMake(8.0, 8.0)
let region1 = MKCoordinateRegion(center: start, span: span)
self.mapView.setRegion(region1, animated: true)
```

#### 沿著線飛

這時候要利用到 perform(#selector(updatePlane), with: self, afterDelay: 0.4) 讓它每隔0.4秒去更新一次飛機位置, 直到到終點為止, 當然也要一直更新region, 以讓飛機維持在正中央

```swift
func updatePlane() {
    planePositionIndex = planePositionIndex + step;
    
    if (planePositionIndex >= geodesicPolyline!.pointCount)
    {
        //plane has reached end, stop moving
        planePositionIndex = 0
        return;
    }
    
    let s = 8.0
    
    let nextMapPoint = geodesicPolyline?.points()[planePositionIndex]

    thePlane.coordinate = MKCoordinateForMapPoint(nextMapPoint!) 
    mapView.region = MKCoordinateRegionMake(thePlane.coordinate, MKCoordinateSpan(latitudeDelta: s, longitudeDelta: s))

    perform(#selector(updatePlane), with: self, afterDelay: 0.4)
}
```

### 結果

{{< youtube GSAd1Vygkhw >}}

這邊有幾個缺點尚待改進

1. 飛機閃爍
1. 如果把可視區域範圍縮小, 或是位置更新過快, 就會造成地圖來不及載入的現象(可能需要預跑幾次將map tile載入cache內)
1. 飛機頭永遠保持往上, 應該要朝著線方向轉

### 完整程式碼

{{< gist julianshen 229f4ac32b3893816bd7636b96fe6f7d >}}