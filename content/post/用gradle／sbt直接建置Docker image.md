---
date: 2021-08-20T19:11:48+08:00
title: "用gradle／sbt直接建置Docker Image"
images: 
- "https://og.jln.co/jlns1/55SoZ3JhZGxl77yPc2J055u05o6l5bu6572uRG9ja2VyIEltYWdl"
---

現在流行甚麼東西都要包成Docker image來部署, 每一個都要寫一個Dockerfile, 每個又很類似, 建置完程式碼又得立刻跑Docker來建置Docker image, 如果有更簡單的方法可以一次做完, 應該會省事很多

以下提到的這些方法可以一行Dockerfile都不用寫

## 使用gradle把Java application建置並包裝成Docker image

要達到這個目的, 可以使用 `com.bmuschko.docker-java-application` - https://bmuschko.github.io/gradle-docker-plugin/#java-application-plugin 這個Gradle plugin, 這個plugin會去呼叫Docker client API, 所以使用前確定你有啟動dockerd

先來看一下底下這範例, 這是一個用[Ktor](https://ktor.io/)寫的server application:

```kotlin
val ktor_version: String by project
val kotlin_version: String by project
val logback_version: String by project

plugins {
    application
    kotlin("jvm") version "1.5.21"
    id("com.bmuschko.docker-java-application") version "6.7.0"
}

group = "co.jln"
version = "0.0.1"
application {
    mainClass.set("co.jln.ApplicationKt")
}

docker {
    javaApplication {
        baseImage.set("openjdk:11-jre")
        ports.set(listOf(8080))
        images.set(setOf("julianshen/myexamplekt:" + version, "julianshen/myexamplekt:latest"))
        jvmArgs.set(listOf("-Xms256m", "-Xmx2048m"))
    }
}

repositories {
    mavenCentral()
}

dependencies {
    implementation("io.ktor:ktor-server-core:$ktor_version")
    implementation("io.ktor:ktor-server-netty:$ktor_version")
    implementation("ch.qos.logback:logback-classic:$logback_version")
    testImplementation("io.ktor:ktor-server-tests:$ktor_version")
    testImplementation("org.jetbrains.kotlin:kotlin-test:$kotlin_version")
}
```

雖然是kotlin寫的, 不過這是一個標準的Java Application, 有它的main class, 我們要加上的只有 `id("com.bmuschko.docker-java-application") version "6.7.0"` 和 `docker {}` 內的內容而已, 基本主要就是baseImage跟image的名字就好了

執行 `gradle dockerBuildImage` 就可以直接幫你把程式建置好並包裝成docker image了

如果是要把image給push到repository上的話, 執行 `gradle dockerPushImage`

如果你的不是一般的Java application 而是Spring boot application的話, 則是可以用 `com.bmuschko.docker-spring-boot-application` 這個plugin而不是上面那個, 不過如果是Spring boot 2.3之後, 還有另一個方法

## 直接把Spring Boot應用程式建置成Docker image

如果是Spring Boot 2.3之後, 因為內建就有支援 [Cloud Native  Buildpacks](https://buildpacks.io/) , 所以直接就可以建置成docker image , 蠻簡單的, 只要執行

```
gradle bootBuildImage
```

不過, 它image的名字會是 `library/project_name` , 所以如果你需要用其他的名字取代的話, 有兩種方法, 一種是加上 `--imageName` 給定名字, 像是: 

```
gradle bootBuildImage --imageName=julianshen/springsample
```

另一種是把這段加到 `build.gradle.kts` 去(這邊以kotlin當範例):

```kotlin
tasks.getByName<org.springframework.boot.gradle.tasks.bundling.BootBuildImage>("bootBuildImage") {
	docker {
		imageName = "julianshen/sprintsmaple"
	}
}
```

這個的缺點是, 不像前面提到的plugin有支援push, 如果你需要把建置好的結果放到repository上的話, 就得自己執行 `docker push`

## 把Scala application包裝成Docker image

如果是Scala就需要用到 [SBT Native Packager](https://www.scala-sbt.org/sbt-native-packager/index.html)

用法也不難, 首先先把下面這段加到 `plugins.sbt` 去:

```
addSbtPlugin("com.typesafe.sbt" % "sbt-native-packager" % "1.7.6")
```

在 `build.sbt` 內加入:

```
enablePlugins(DockerPlugin)
```

然後執行 `sbt docker:publishLocal`即可, 相關設定可以參考 [Docker plugin的文件](https://www.scala-sbt.org/sbt-native-packager/formats/docker.html)