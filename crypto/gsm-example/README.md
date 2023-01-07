## 数据加密方式国密2(sm2)

公钥pubKey:
```
0417d818ce56b3ffc803521ebd48e35230f490c33108e4a8d2d3cebd99ac4890f521e32b1b7e4182f36edd87ba5dd022d0d33b2a3ca66528a7e14a425f6e289002
```
私钥prvKey:
```
2beaa280c903d800c99dc0a17195e4e8fc984ba723fc48f19585a1b4e74f788c
```

```shell
### 特别注意
上面的公钥和私钥是不是一对，是两套。套A的公钥 和套B的私钥。
加签的时候用套B的私钥去加签，服务端会用自己的套B的公钥去验签。
在使用下面验签工具的时候，如果公钥，和私钥都填了文档中的公钥和私钥，肯定会验签false。（不用担心，就是如此）
```



- SM2在线加密验签工具 :http://sm.skill86.com/sm2.html
- SM4在线加密工具:http://sm.skill86.com/sm4.html

## 验签步骤
#### 1. 先去sm4 把自己想要加密的数据进行加密。
sm4 key:
```
@PHXGV9Tb9V+8-J&
```
![file](http://qiniu.skill86.com/20230103/2ZUI18v6rdiTfmr3KXZYZhrpYp4eHInu2jJ5Nhl7.png)
#### 2. 拿到加密后的数据去sm2里面生成签名
![file](http://qiniu.skill86.com/20230103/M1shjGiBs24V4hBmkQ8Fzeucji1g9p6aje9LK0BL.png)
#### 3. 在请求连接的header 添加sign 参数
![file](http://qiniu.skill86.com/20230107/vPbSIbCgRTVEublIYuXv9LOQi51MmJBSplm8cVM8.png)
