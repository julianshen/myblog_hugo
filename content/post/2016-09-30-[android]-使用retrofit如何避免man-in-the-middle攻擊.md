---
date: "2016-09-30T21:47:36Z"
images:
- /images/posts/2016-09-30-[android]-使用retrofit如何避免man-in-the-middle攻擊.md.jpg
tags:
- Android
- Retrofit
- MITM
- Security
- HTTPS
- SSL
title: '[Android] 使用Retrofit如何避免Man in the middle攻擊'
---
前幾天看到朋友寫如何用[Charles Proxy](https://www.charlesproxy.com)來debug REST API, 
突然就想到了這個問題, [Charles Proxy](https://www.charlesproxy.com)是一個不錯的工具, 我自己也蠻常用的
, 除了可以用來debug HTTP以外, 其實HTTPS也可以(參照[SSL proxying](https://www.charlesproxy.com/documentation/proxying/ssl-proxying/))
, 只要在手機上安裝好它的SSL certificate即可

不過這不是這篇要講的重點, 這種debugging的方式原理即是用Proxy作為一個中間人來紀錄HTTP/HTTPS的傳輸內容
, 這也帶來一個問題, 也就是你的資料是可以這樣簡單的暴露出去, 這種即是所謂的MITM (Man in the middle) 中間人攻擊

{{<mermaid>}}
graph LR;
    Client-->Proxy[Proxy - Middle man];
    Proxy[Proxy - Middle man]-->Client;
    Proxy[Proxy - Middle man]-->Server;
    Server-->Proxy[Proxy - Middle man];
{{</mermaid>}}

一般來說, 用了SSL/HTTPS後, 也不是那麼容易安插一個中間人的 ,像[Charles Proxy](https://www.charlesproxy.com)這樣的東西,
還是需要使用者自己去手機上的安全性設定裝信任憑證(trusted certificate), 一般常見的問題反而來自於開發者本身
, 很多人常以為走了HTTPS就安全了, 卻常常為了debug方便, 或是省錢用Self-signed certificate 而去覆寫了系統預設的
X509TrustManager, 以至於整個檢查機制被跳過, 真的漏洞都來自於人的惰性,
更多資訊可以參考:

1. [窃听风暴： Android平台https嗅探劫持漏洞](https://security.tencent.com/index.php/blog/msg/41)
1. [Android安全之Https中间人攻击漏洞](http://yaq.qq.com/blog/13)
2. [Android HTTPS中间人劫持漏洞浅析](http://www.droidsec.cn/android-https中间人劫持漏洞浅析/)
2. [Android客户端安全 -> HTTPS敏感数据劫持漏洞](http://blog.csdn.net/u013107656/article/details/52036158)
2. [Android App 安全的HTTPS 通信](http://pingguohe.net/2016/02/26/Android-App-secure-ssl.html)
2. [手機應用程式開發上被忽略的 SSL 處理](http://devco.re/blog/2014/08/15/ssl-mishandling-on-mobile-app-development/)
3. [How To: Use mitmproxy to read and modify HTTPS traffic](https://blog.heckel.xyz/2013/07/01/how-to-use-mitmproxy-to-read-and-modify-https-traffic-of-your-phone/)
4. [INTERCEPTING ANDROID SSL / HTTPS TRAFFIC](http://www.7xter.com/2015/07/intercepting-android-ssl-traffic.html)

該避免去寫的, 就該避掉, 以免產生不必要的漏洞, 但不過也有可能碰到使用者的手機被(有意/無意)安裝了攻擊者的憑證到系統的信任憑證中, 
那麼碰到這種, 應用程式本身該如何保護自己?

這邊有兩個方式可以採用 - 

1. certificate pinning
2. 在程式中寫死certificate

以下的範例以Android最紅, 常被使用的[Retrofit](https://square.github.io/retrofit/)為範例(雖說是[Retrofit](https://square.github.io/retrofit/), 但其實是在[okhttp](https://github.com/square/okhttp)動的手腳)

#### Certificate Pinning ####

或叫[HTTP Public Key Pinning (HPKP)](https://en.wikipedia.org/wiki/HTTP_Public_Key_Pinning), 
簡單的說, 就是在程式內綁定Certificate的public key, 如果今天中間人存在, certificate不同了, 連接就不會成功

先來看範例

```java
OkHttpClient client = new OkHttpClient.Builder()
        .certificatePinner(new CertificatePinner.Builder().add("cdn.rawgit.com", "sha256/p0962nIqD0mv1APQ1mgmRiuwrhZXBuj+t6dey/Adk0U=").build())
        .build();
Retrofit retrofit = new Retrofit.Builder()
        .baseUrl("https://cdn.rawgit.com/julianshen/cpblschedule/")
        .addConverterFactory(GsonConverterFactory.create())
        .client(client)
        .build();
```

這邊所動的手腳(~~跟Retrofit沒啥鳥關係~~)是在OkHttpClient上加上一個CertificatePinner,
利用了certificate pinner加上了一個host name跟public key的對應,
以這例子來說: __"sha256/p0962nIqD0mv1APQ1mgmRiuwrhZXBuj+t6dey/Adk0U="__ 這個public key對應的是 __"cdn.rawgit.com"__
, 因此碰到"cdn.rawgit.com"來的url, 會檢查public key是否符合,
這邊public key的參數字串是要以"sha1/"或"sha256/"開始的, 依使用的演算法而不同

但又要怎取得這串"天書"呢?

很簡單, 首先確認你的電腦裡面有沒安裝openssl , 有的話就用以下的指令:

```shell
openssl s_client -connect cdn.rawgit.com:443 -servername cdn.rawgit.com | openssl x509 -pubkey -noout | openssl rsa -pubin -outform der | openssl
dgst -sha256 -binary | openssl enc -base64
```

最後產生的內容如下

```
depth=2 /C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Certification Authority
verify error:num=20:unable to get local issuer certificate
verify return:0
writing RSA key
p0962nIqD0mv1APQ1mgmRiuwrhZXBuj+t6dey/Adk0U=
```

所以我們要填的參數就是 __"sha256/p0962nIqD0mv1APQ1mgmRiuwrhZXBuj+t6dey/Adk0U="__

關於HPKP的openssl相關的指令可以參考[Public Key Pinning](https://developer.mozilla.org/en-US/docs/Web/Security/Public_Key_Pinning)

如果你用[Charles Proxy](https://www.charlesproxy.com)當中間人來測試這段程式的話, 你會得到下面這樣的exception

```
javax.net.ssl.SSLPeerUnverifiedException: Certificate pinning failure!
                                            Peer certificate chain:
                                            sha256/ENlqFHtfARof3AK50Hbc1sj47M5hWhc5kQ5Z2vyfxq4=: CN=rawgit.com,OU=PositiveSSL Multi-Domain,OU=Domain Control Validated
                                            sha256/mXmLzo8k5ANwi11PlLSW/b4OC/Anjw1OeACDyZxD/WM=: C=NZ,ST=Auckland,L=Auckland,O=XK72 Ltd,OU=http://charlesproxy.com/ssl,CN=Charles Proxy Custom Root Certificate (built on Jlnmbp-retina.local\, 11 十二月 2015)
                                            Pinned certificates for cdn.rawgit.com:
                                            sha256/p0962nIqD0mv1APQ1mgmRiuwrhZXBuj+t6dey/Adk0U=
```

#### 在程式中寫死certificate ####

這方法聽起來暴力了一點, 但原理其實跟前一個沒啥太大差別, 只是一個寫死public key一個寫死certificate

一樣先來看個範例:

```java
String cert = ""
        +"-----BEGIN CERTIFICATE-----\n"
        +"MIIFdjCCBF6gAwIBAgIRAPbTBHdHHseM4hArA7n6uREwDQYJKoZIhvcNAQELBQAw\n"
        +"gZAxCzAJBgNVBAYTAkdCMRswGQYDVQQIExJHcmVhdGVyIE1hbmNoZXN0ZXIxEDAO\n"
        +"BgNVBAcTB1NhbGZvcmQxGjAYBgNVBAoTEUNPTU9ETyBDQSBMaW1pdGVkMTYwNAYD\n"
        +"VQQDEy1DT01PRE8gUlNBIERvbWFpbiBWYWxpZGF0aW9uIFNlY3VyZSBTZXJ2ZXIg\n"
        +"Q0EwHhcNMTYwMTAxMDAwMDAwWhcNMTcwMTEzMjM1OTU5WjBbMSEwHwYDVQQLExhE\n"
        +"b21haW4gQ29udHJvbCBWYWxpZGF0ZWQxITAfBgNVBAsTGFBvc2l0aXZlU1NMIE11\n"
        +"bHRpLURvbWFpbjETMBEGA1UEAxMKcmF3Z2l0LmNvbTCCASIwDQYJKoZIhvcNAQEB\n"
        +"BQADggEPADCCAQoCggEBALD2sUNSYp2R+mSEMF2qKvMS780qTkltvqP4rwEGKOLV\n"
        +"rR5QQWo8vzSlgZvxVsguRHi0pPBtVAH794L43tD+IuyQlDlNU2qc1aMqBkn3S2wN\n"
        +"Way8BS9w80pgeFnObZiFJtPI9pdwcB72Bgq8Nlc25oVrDVWr/Q8nLIKS/9FkNs+C\n"
        +"MPU00vGFZSFbR7s15ORj9+qPCskWZcHpQ+m9EKmZD3IVKj3QQyBD17cBoVkYoIpj\n"
        +"I8/r+NfVfKetlsB7Pcv8P3yLFVFC4+PGrnW9TLZK9aRtbDTvjElXwCQIPK5B2fKw\n"
        +"SJK26IJPDX0r1JkeuR53Afr509iEx3xAHcRz/kR5uo8CAwEAAaOCAf0wggH5MB8G\n"
        +"A1UdIwQYMBaAFJCvajqUWgvYkOoSVnPfQ7Q6KNrnMB0GA1UdDgQWBBSI83Kqafo7\n"
        +"kg9cmyT1vnceH4fAUjAOBgNVHQ8BAf8EBAMCBaAwDAYDVR0TAQH/BAIwADAdBgNV\n"
        +"HSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTwYDVR0gBEgwRjA6BgsrBgEEAbIx\n"
        +"AQICBzArMCkGCCsGAQUFBwIBFh1odHRwczovL3NlY3VyZS5jb21vZG8uY29tL0NQ\n"
        +"UzAIBgZngQwBAgEwVAYDVR0fBE0wSzBJoEegRYZDaHR0cDovL2NybC5jb21vZG9j\n"
        +"YS5jb20vQ09NT0RPUlNBRG9tYWluVmFsaWRhdGlvblNlY3VyZVNlcnZlckNBLmNy\n"
        +"bDCBhQYIKwYBBQUHAQEEeTB3ME8GCCsGAQUFBzAChkNodHRwOi8vY3J0LmNvbW9k\n"
        +"b2NhLmNvbS9DT01PRE9SU0FEb21haW5WYWxpZGF0aW9uU2VjdXJlU2VydmVyQ0Eu\n"
        +"Y3J0MCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5jb21vZG9jYS5jb20wSwYDVR0R\n"
        +"BEQwQoIKcmF3Z2l0LmNvbYIVY2RuLW9yaWdpbi5yYXdnaXQuY29tgg5jZG4ucmF3\n"
        +"Z2l0LmNvbYINcmF3Z2l0aHViLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEASx3e6y9e\n"
        +"9F67N5NIDausDHDL+/fz6uj2DDNJdaQvALAYDV8hKpz7+QGotlQfI042U2i83J7m\n"
        +"DGZiKDpzaXa7IFfGiq6PFPUjntElHoU2E4vRb5LApg1sJ5YueYmRd3X7x0/jPYg6\n"
        +"TiTtHyXOnlMuqL2FUYJM0BP7cMfOppvUF0R2zHUVA0rVHuLrSStg8bMg8aVUIkKg\n"
        +"n3NS+eg4ofo95jaMRVykhLnylkYBk9dWBM6B19Yw7LgQd94MZSa+Xix4HxeFgvAM\n"
        +"BF+oeDMaY7mEN4Xm5hnrIS1FElRDq/ckpDV8JVpR4SQqB+tzSZPPIiep8EvgRDzd\n"
        +"2EkL187FF3MQ+w==\n"
        +"-----END CERTIFICATE-----\n";

X509TrustManager trustManager = null;
SSLSocketFactory sslSocketFactory = null;

try {
    CertificateFactory certificateFactory = CertificateFactory.getInstance("X.509");
    Collection<? extends Certificate> certificates = certificateFactory.generateCertificates(new Buffer().writeUtf8(cert).inputStream());
    if (certificates.isEmpty()) {
        throw new IllegalArgumentException("expected non-empty set of trusted certificates");
    } else {
        char[] password = "password".toCharArray(); // Any password will work.
        KeyStore keyStore = KeyStore.getInstance(KeyStore.getDefaultType());
        keyStore.load(null, password);

        int index = 0;
        for (Certificate certificate : certificates) {
            String certificateAlias = Integer.toString(index++);
            keyStore.setCertificateEntry(certificateAlias, certificate);
        }

        TrustManagerFactory trustManagerFactory = TrustManagerFactory.getInstance(TrustManagerFactory.getDefaultAlgorithm());
        trustManagerFactory.init(keyStore);

        TrustManager[] trustManagers = trustManagerFactory.getTrustManagers();
        if (trustManagers.length != 1 || !(trustManagers[0] instanceof X509TrustManager)) {
            throw new IllegalStateException("Unexpected default trust managers:"
                    + Arrays.toString(trustManagers));
        }
        trustManager = (X509TrustManager) trustManagers[0];

        SSLContext sslContext = SSLContext.getInstance("TLS");
        sslContext.init(null, new TrustManager[] { trustManager }, null);
        sslSocketFactory = sslContext.getSocketFactory();
    }
} catch (CertificateException e) {
    e.printStackTrace();
} catch (NoSuchAlgorithmException e) {
    e.printStackTrace();
} catch (KeyStoreException e) {
    e.printStackTrace();
} catch (IOException e) {
    e.printStackTrace();
} catch (KeyManagementException e) {
    e.printStackTrace();
}

OkHttpClient client = new OkHttpClient.Builder()
        .sslSocketFactory(sslSocketFactory, trustManager)
        .build();

Retrofit retrofit = new Retrofit.Builder()
        .baseUrl("https://cdn.rawgit.com/julianshen/cpblschedule/")
        .addConverterFactory(GsonConverterFactory.create())
        .client(client)
        .build();
```

看到這麼長的東西大概就不會想看了吧?(老實說我也寫的很懶, 這邊要特別提到有用到一個Class叫做Buffer,這是來自於[okio](https://github.com/square/okio), 蠻方便的東西), 跟前一個不同的地方在於, 前一個利用了CertificatePinner,
而這一個則是改了sslSocketFactory的TrustManager(喂!!!前面不是隱約有提到這樣不太好?!),
這個TrustManager全部也只信任這一個certifcate, 如果碰到其他的, 就會發生以下的exception:

```
javax.net.ssl.SSLHandshakeException: java.security.cert.CertPathValidatorException: Trust anchor for certification path not found.
        at com.android.org.conscrypt.OpenSSLSocketImpl.startHandshake(OpenSSLSocketImpl.java:328)
```

那, 上面那一大坨憑證資料是怎樣取得的? 方法其實差不多:

    openssl s_client -connect cdn.rawgit.com:443 -showcerts

會得到這樣的結果:

```
CONNECTED(00000003)
depth=2 /C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Certification Authority
verify error:num=20:unable to get local issuer certificate
verify return:0
---
Certificate chain
 0 s:/OU=Domain Control Validated/OU=PositiveSSL Multi-Domain/CN=rawgit.com
   i:/C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Domain Validation Secure Server CA
 1 s:/OU=Domain Control Validated/OU=PositiveSSL Multi-Domain/CN=rawgit.com
   i:/C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Domain Validation Secure Server CA
 2 s:/C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Domain Validation Secure Server CA
   i:/C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Certification Authority
 3 s:/C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Certification Authority
   i:/C=SE/O=AddTrust AB/OU=AddTrust External TTP Network/CN=AddTrust External CA Root
---
Server certificate
-----BEGIN CERTIFICATE-----
MIIFdjCCBF6gAwIBAgIRAPbTBHdHHseM4hArA7n6uREwDQYJKoZIhvcNAQELBQAw
gZAxCzAJBgNVBAYTAkdCMRswGQYDVQQIExJHcmVhdGVyIE1hbmNoZXN0ZXIxEDAO
BgNVBAcTB1NhbGZvcmQxGjAYBgNVBAoTEUNPTU9ETyBDQSBMaW1pdGVkMTYwNAYD
VQQDEy1DT01PRE8gUlNBIERvbWFpbiBWYWxpZGF0aW9uIFNlY3VyZSBTZXJ2ZXIg
Q0EwHhcNMTYwMTAxMDAwMDAwWhcNMTcwMTEzMjM1OTU5WjBbMSEwHwYDVQQLExhE
b21haW4gQ29udHJvbCBWYWxpZGF0ZWQxITAfBgNVBAsTGFBvc2l0aXZlU1NMIE11
bHRpLURvbWFpbjETMBEGA1UEAxMKcmF3Z2l0LmNvbTCCASIwDQYJKoZIhvcNAQEB
BQADggEPADCCAQoCggEBALD2sUNSYp2R+mSEMF2qKvMS780qTkltvqP4rwEGKOLV
rR5QQWo8vzSlgZvxVsguRHi0pPBtVAH794L43tD+IuyQlDlNU2qc1aMqBkn3S2wN
Way8BS9w80pgeFnObZiFJtPI9pdwcB72Bgq8Nlc25oVrDVWr/Q8nLIKS/9FkNs+C
MPU00vGFZSFbR7s15ORj9+qPCskWZcHpQ+m9EKmZD3IVKj3QQyBD17cBoVkYoIpj
I8/r+NfVfKetlsB7Pcv8P3yLFVFC4+PGrnW9TLZK9aRtbDTvjElXwCQIPK5B2fKw
SJK26IJPDX0r1JkeuR53Afr509iEx3xAHcRz/kR5uo8CAwEAAaOCAf0wggH5MB8G
A1UdIwQYMBaAFJCvajqUWgvYkOoSVnPfQ7Q6KNrnMB0GA1UdDgQWBBSI83Kqafo7
kg9cmyT1vnceH4fAUjAOBgNVHQ8BAf8EBAMCBaAwDAYDVR0TAQH/BAIwADAdBgNV
HSUEFjAUBggrBgEFBQcDAQYIKwYBBQUHAwIwTwYDVR0gBEgwRjA6BgsrBgEEAbIx
AQICBzArMCkGCCsGAQUFBwIBFh1odHRwczovL3NlY3VyZS5jb21vZG8uY29tL0NQ
UzAIBgZngQwBAgEwVAYDVR0fBE0wSzBJoEegRYZDaHR0cDovL2NybC5jb21vZG9j
YS5jb20vQ09NT0RPUlNBRG9tYWluVmFsaWRhdGlvblNlY3VyZVNlcnZlckNBLmNy
bDCBhQYIKwYBBQUHAQEEeTB3ME8GCCsGAQUFBzAChkNodHRwOi8vY3J0LmNvbW9k
b2NhLmNvbS9DT01PRE9SU0FEb21haW5WYWxpZGF0aW9uU2VjdXJlU2VydmVyQ0Eu
Y3J0MCQGCCsGAQUFBzABhhhodHRwOi8vb2NzcC5jb21vZG9jYS5jb20wSwYDVR0R
BEQwQoIKcmF3Z2l0LmNvbYIVY2RuLW9yaWdpbi5yYXdnaXQuY29tgg5jZG4ucmF3
Z2l0LmNvbYINcmF3Z2l0aHViLmNvbTANBgkqhkiG9w0BAQsFAAOCAQEASx3e6y9e
9F67N5NIDausDHDL+/fz6uj2DDNJdaQvALAYDV8hKpz7+QGotlQfI042U2i83J7m
DGZiKDpzaXa7IFfGiq6PFPUjntElHoU2E4vRb5LApg1sJ5YueYmRd3X7x0/jPYg6
TiTtHyXOnlMuqL2FUYJM0BP7cMfOppvUF0R2zHUVA0rVHuLrSStg8bMg8aVUIkKg
n3NS+eg4ofo95jaMRVykhLnylkYBk9dWBM6B19Yw7LgQd94MZSa+Xix4HxeFgvAM
BF+oeDMaY7mEN4Xm5hnrIS1FElRDq/ckpDV8JVpR4SQqB+tzSZPPIiep8EvgRDzd
2EkL187FF3MQ+w==
-----END CERTIFICATE-----
subject=/OU=Domain Control Validated/OU=PositiveSSL Multi-Domain/CN=rawgit.com
issuer=/C=GB/ST=Greater Manchester/L=Salford/O=COMODO CA Limited/CN=COMODO RSA Domain Validation Secure Server CA
---
No client certificate CA names sent
---
SSL handshake has read 6716 bytes and written 456 bytes
---
New, TLSv1/SSLv3, Cipher is DHE-RSA-AES256-SHA
Server public key is 2048 bit
Secure Renegotiation IS supported
Compression: NONE
Expansion: NONE
SSL-Session:
    Protocol  : TLSv1
    Cipher    : DHE-RSA-AES256-SHA
    Session-ID: ACFED5197F4B5F64DBAAC97A47380CD9283EFD1797E27EF43A215B06982698F4
    Session-ID-ctx: 
    Master-Key: 856DFB81C863F461FE75EE11E633EA5CAB3CD6683563F597F406EF19AB37CD3F45B4FE659AB8726F9291360CBE856F48
    Key-Arg   : None
    Start Time: 1475223830
    Timeout   : 300 (sec)
    Verify return code: 0 (ok)
---
```

BEGIN到END那段剪貼下來即是

## 缺點 ##

這兩個方法有一個極大的缺點, 如果沒有解決方案還是最好別用(~~啊講那麼多, 最後不能用?!~~)
, 這缺點就是SSL Certificate是有時效性的, 因此server端的憑證是會換會改變的,
只要server端一換, 就可能造成client端的失效, 所以如果要使用這兩個方法,
前提必須先解決的是如何更新Client端要檢查的public key或certificate, 
這邊就不討論這部份的機制了, 這應該有很多方法可以來作才對

\#好久沒寫關於Android的內容了 