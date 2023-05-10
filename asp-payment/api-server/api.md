# Sunny Pay Api Document

## 1 接口规则

### 1.1 接口协议

#### 调用Sunny Pay接口须遵守以下规则：


1. 请求方式一律使用 POST 并将请求的数据参数的JSON字符串以 http body的方式传递
2. 传输方式：采用HTTP传输

3. 请求的数据格式统一使用JSON格式，若参数有json字符串请转义双引号（\"）

4. 字符串编码请统一使用UTF-8

5. 签名算法base64_encode(hmac_sha256)

#### 注意：
1. <font color="red">交易金额：默认为卢比交易，单位为分，参数值不能带小数.</font>
2. 接口文档部分字段使用`xxxxxxxxxxx`做了信息脱敏处理，注意甄别
3. <font color="red">参与计算签名的字段 如果为空（长度为0） 则会抛弃 （数字为0的情况 长度为1 不会抛弃）</font>
4. 参数名ASCII码从小到大排序（字典序）
5. 参数名区分大小写
----
### 1.2 参数签名

## <span id="header">Header</span>

| 参数名称  | 必须参数 | 说明                          |
| :-------- | -------- | --------------------------  |
| Version   | 是       | 版本信息 固定为 1.0          |
| AppId     | 是       | 应用标识，通过关联后台获取    |
| Signature | 是       | 数据签名                    |

1.2.1. 假设请求参数如下：

```json
{
    "payment_method":"sunny.h5",
    "order_id":"100867710086771120008",
    "user_id":"1008677",
    "order_currency":"INR",
    "order_amount":100,
    "return_url":"https://xxx.com/api/v1/notify"
}
```

1.2.2. 将参数按照键值对（key=value）的形式排列,按照参数名ASCII字典序排序,并用&连接

```
str = "order_amount=100&order_currency=INR&order_id=100867710086771120008&payment_method=sunny.h5&return_url=https://xxx.com/api/v1/notify&user_id=1008677"
```

1.2.3. 在最后拼接上密钥字符串 `&paySecret=sQkj0RH9qMxdaxo0sJ8xlbki4ssOjvZd`

```
str += "&paySecret=sQkj0RH9qMxdaxo0sJ8xlbki4ssOjvZd"
```
即 str 为:

```
order_amount=100&order_currency=INR&order_id=100867710086771120008&payment_method=sunny.h5&return_url=https://xxx.com/api/v1/notify&user_id=1008677&paySecret=sQkj0RH9qMxdaxo0sJ8xlbki4ssOjvZd
```


1.2.4. 最后计算签名 base64_encode(hmac_sha256) 得到的结果 放入 header Signature 中

```
sign = base64_encode(hmac_sha256('order_amount=100&order_currency=INR&order_id=100867710086771120008&payment_method=sunny.h5&return_url=https://xxx.com/api/v1/notify&user_id=1008677&paySecret=sQkj0RH9qMxdaxo0sJ8xlbki4ssOjvZd', sQkj0RH9qMxdaxo0sJ8xlbki4ssOjvZd))
```

sign的值为:

```
YzBhM2UzOWUwMmViNGJhZmNiMTAyM2I1Mjc2Y2Y0MWYxMWNmZDhjMjQ3NWI5NGIwMDlkNjcxNmU1MjU1NTZiNg==
```

### 1.4 签名示例代码：

``` php
// PHP 代码
$Signature = base64_encode(hash_hmac('sha256', $body, $SerectKey));

```

``` golang
// golang
key := []byte(secret)
h := hmac.New(sha256.New, key)
h.Write([]byte(message))
sha := hex.EncodeToString(h.Sum(nil))
Signature = base64.StdEncoding.EncodeToString([]byte(sha))

```

----

### 1.5 接口域名

在不同的业务区域，使用不同的域名

```
测试地址：
https://hopepatti.fun


```

### 1.6 响应码说明

响应码| 响应码业务说明                  | 逻辑处理
-----|--------------------------|-----------------------------------------------
4000   | 通用错误                     |未归类的错误，需提示用户或做相应的容错处理
200| 成功                       |按照具体业务逻辑处理
400| 请求参数错误, 请求被拒绝            |检查参数，做出调整
401| 签名错误,请求被拒绝               |检查签名逻辑以及相应的密钥是否正确
403| 权限错误,请求被拒绝               |请联系客服，开通对应的权限
409| 请求数请求有冲突。有关详细信息，请参阅随附的错误消息 |请重新生成，请求参数的值
429| 请求超出了速率限制。在指定的时间段后重试     |稍后再重新调用该接口
430| 请求渠道超出了速率限制。在指定的时间段后重试   |稍后再重新调用该接口
500| 服务器开小差了!                 |稍后再重新调用该接口
501| 服务器没有实现该方法               |根据文档检查接口地址
502| 用户处理异常，请重试               |稍后再重新调用该接口
504| 请求已过期                    |稍后再重新调用该接口
505| 头信息参数错误                  |稍后再重新调用该接口
506| 您的操作尚未完成，请稍等             |稍后再重新调用该接口
507| 写入数据失败，请重试               |稍后再重新调用该接口
508| 渠道初始化错误 请求被拒绝            |稍后再重新调用该接口
509| 渠道网络初始化错误 请求被拒绝          |稍后再重新调用该接口
510| 渠道网络错误 请求被拒绝             |稍后再重新调用该接口
511| 渠道分配错误 请求被拒绝             |稍后再重新调用该接口
512| 渠道内部错误 请求被拒绝             |稍后再重新调用该接口
513| 渠道内部错误商户余额错误 请求被拒绝       |稍后再重新调用该接口
514| 渠道内部错误商户错误 请求被拒绝         |稍后再重新调用该接口
515| 渠道内部错误上游错误 请求被拒绝         |稍后再重新调用该接口
518| 请求失败                     |稍后再重新调用该接口
4002| appid不存在                 |请联系客服
4003| appid不可用                 |请联系客服
4004| 商户项目不存在                  |请联系客服
4005| 商户不存在                    |请联系客服
4006| 渠道信息不存在                  |请联系客服
4007| 代收订单不存在                  |检查代收订单号
4008| 代付订单不存在                  |检查代付订单号
4009| ip白名单不存在                 |请联系客服
5000| 商户金额已超出限制，有关详细信息，请参阅商户限额 |请参阅商户限额

### 1.7 通用响应数据结构

#### 注意：
- <font color="red">网络状态 HTTP 状态码 成功都是200， 业务的逻辑判断 根据 业务code 码 来处理</font>

响应数据为Json格式,Key-Map 描述如下:

属性    | 说明 | 示例
-------|------|-------
code   | 业务响应码(用来判断业务的逻辑) | 200（非200 代表异常情况）
msg    | 业务提示消息，根据业务结果，可直接使用该属性值提示用户 | Success
data   | 业务数据，需根据相应接口进行逻辑处理,有时为空(不存在该属性) |

正常示例

```
{
  "code": 200,
  "msg": "success",
  "data": {
    "appid": "1000258",
    "order_id":"100867710086771120008",
    "order_currency":"INR",
  }
}

```


失败示例

```
{
  "code": 40500,
  "msg": "余额不足",
  "status": "error",
  "data": {}
}
```


## 2 接口列表

### 2.1 创建代收


### Api:

  * #### 请求
    ``` URL: /pay/order```  
    ``` Method: POST```  
    ``` Content-Type:  application/json```
  *  #### header
    参考[header](#header)部分

#### body Parameters 请求参数

| 字段 | 变量名           | 必填 | 类型 | 描述                                                   |
|----|---------------|----|----|------------------------------------------------------|
| payment_method | payment_method | 是 | String | 支付类型 固定 sunny.h5                                     |
| 商户自身订单号 | order_id      | 是 | String | 如果商户有自己的订单系统，可以自己生成订单                                |
| 用户ID | user_id       | 是 | String | 用户 Id                                                |
| 币种 | order_currency | 是 | String | 限定支付币种（目前支持 INR）                                     |
| 订单金额 | order_amount  | 是 | Int | 支付金额 单位为"分" 如10即0.10元                                |
| 订单名称 | order_name    | 是 | String | 订单名称                                                 |
| 异步通知url | notify_url    | 可选 | String | 异步通知url                                              |
| return_url | return_url    | 可选 | String | 支付成功后跳转url                                           |
| 客户的名称 | customer_name | 是 | String | 客户的名称。仅使用字母数字值                                       |
| 客户电话号码 | customer_phone | 是 | String | 客户电话号码                                               |
| 客户电子邮件 | customer_email | 是 | String | 客户电子邮件地址。                                            |
| 用户设备ID | device_info     | 是 | String | 客户用户设备ID。                                            |
| order_note | order_note    | 可选 | String | 订单描述                                                 |
| 附加数据 | attach        | 可选 | String | 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据(中文需要url编码) |


#### payment 参数说明

| 参数值 | 描述 |
|----------|----------|
|sunny.h5 | 网页支付 |



支付返回后，检查交易状态trade_state,并根据其结果，决定是否调用订单查询接口进行结果查询处理

#### 订单支付状态 trade_state 说明

```
PENDING:代收中
SUCCESS:代收成功
EXPIRED:已过期
RETURNED:已取消
PAYERROR:代收异常
FAILED:代收失败
```

#### 交易类型 trade_type 说明
```
PAYOUT: 代付
H5: H5支付
```

#### 正确响应数据说明

响应结果response.data数据说明

| 字段 | 变量名 | 类型 | 描述 |
|----|----|---|----|
| Sunny订单编号 | sn | String | 示例: 1120180209xxxxxxxxxxxxxxxxxx 唯一标号，可以用于查询，等操作 |
| 商户Id | out\_trade\_no | String | 商户订单号 如: "1120180209xxxxxxxxxxxxxxxx" |
| appid | appid | String | 示例: 1000258 唯一标号，商户项目ID |
| 下单币种 | fee_type | String | 交易货币 如：INR |
| 实付币种 | cash_fee_type | String | 客户实际代付的币种 如：INR |
| 交易类型 | trade_type | String | 交易类型 如：H5,PAYOUT |
| 支付提供方 | provider | String | 如:sunny,paytm,cashfree| 
| 支付方交易号 | transaction_id | String | P563xxxxxxxxxxxxx |
| 客户IP | client_ip | String | 如:127.0.0.1 |
| 支付token | order_token | String | 如:127.0.0.1 |
| 订单时间 | create_time | Int | 时间戳 如: 1518155270 |
| 支付时间 | finish_time | Int | 成功支付的时间戳 如1518155297 |
| 订单金额 | order_amount | Int | 如:"100" |
| 实际支付金额 | total_fee | Int | 用户需要支付的金额 单位为"分" 如:10 |
| 支付链接 | payment_link | String | 如:"" 收银台h5地址 |
| 交易状态 | trade_state | String| PENDING ||
| 异步通知url | notify_url | String | 异步通知url |
| return_url | return_url | String | 支付成功后跳转url |
| 客户的名称 | customer_name | String | 客户的名称。仅使用字母数字值 |
| 客户电话号码 | customer_phone | String | 客户电话号码 |
| 客户电子邮件 | customer_email | String | 客户电子邮件地址。 |






### 请求示例

#### H5实例：
请求参数

```
{
    "payment_method":"sunny.h5",
    "order_id":"100867710086771120009",
    "user_id":"1008677",
    "order_currency":"INR",
    "order_amount":100,
    "order_name":"order_name",
    "return_url":"https://xxx.com/api/v1/notify",
    "notify_url":"https://xxx.com/api/v1/notify",
    "customer_name":"ericlu",
    "customer_phone":"903683xxxx",
    "customer_email":"ericlu@gmail.com",
    "device_info": "device_infoxxxxx",
    "order_note":"order_note"
}
```
响应结果
```
{
    "code": 200,
    "msg": "下单成功",
    "data": {
        "id": 2647632,
        "sn": "10202210121713495497389001185123",
        "out_trade_no": "100867710086771120010",
        "appid": "1",
        "fee_type": "INR",
        "cash_fee_type": "INR",
        "trade_type": "H5",
        "provider": "sunny",
        "transaction_id": "",
        "client_ip": "127.0.0.1",
        "order_token": "",
        "create_time": 1665566029,
        "finish_time": 0,
        "order_amount": 100,
        "total_fee": 100,
        "payment_link": "",
        "trade_state": "PENDING",
        "notify_url": "https://xxx.com/api/v1/notify",
        "return_url": "https://xxx.com/api/v1/notify",
        "customer_name": "ericlu",
        "customer_email": "ericlu@gmail.com",
        "customer_phone": "903683xxxx"
    }
}
```

### 2.2 创建代付


### Api:

  * #### 请求
    ``` URL: /pay/payout```  
    ``` Method: POST```  
    ``` Content-Type:  application/json```
  *  #### header
    参考[header](#header)部分

#### body Parameters 请求参数

| 字段 | 变量名 | 必填 | 类型 | 描述                                        |
|----|----|----|----|-------------------------------------------|
| 商户自身订单号 | order_id | 是 | String | 如果商户有自己的订单系统，可以自己生成订单                     |
| 用户ID | user_id | 是 | String | 用户 Id                                     |
| 币种 | order_currency | 是 | String | 限定支付币种（目前支持 INR）                          |
| 订单金额 | order_amount | 是 | Int | 支付金额 单位为"分" 如10即0.10元                     |
| 订单名称 | order_name | 是 | String | 订单名称                                      |
| 异步通知url | notify_url | 是 | String | 异步通知url                                   |
| 客户的名称 | customer_name | 是 | String | 客户的名称。仅使用字母数字值                            |
| 客户电话号码 | customer_phone | 是 | String | 客户电话号码                                    |
| 客户电子邮件 | customer_email | 是 | String | 客户电子邮件地址。                                 |
| 用户设备ID | device_info     | 是 | String | 客户用户设备ID。 |
| order_note | order_note | 可选 | String | 订单描述                                      |
| 附加数据 |attach | 可选 | String | 附加数据，在查询API和支付通知中原样返回，该字段主要用于商户携带订单的自定义数据 |
| ifsc | ifsc | 是 | String | 客户银行 IFSC                                 |
| bank_card | bank_card | 是 | String | 客户银行卡号。网银支付需要                             |
| bank_code | bank_code | 是 | String | 银行编码 详见 2.8 银行编码 客户银行代码。            |
| pay_type | pay_type | 是 | String | 支付类型 固定 bank                              |
| address | address | 是 | String | 客户地址、字母数字                                 |
| city | city | 是 | String | 客户城市，只有字母                                 |
| vpa号 | vpa号 | 可选 | String | vpa号,印度upi提现方式需要提供                        |


代付返回后，检查交易状态trade_state,并根据其结果，决定是否调用订单查询接口进行结果查询处理

#### 订单代付状态 trade_state 说明

```
PENDING:代付中
SUCCESS:代付成功
EXPIRED:已过期
RETURNED:已取消
PAYERROR:代付异常
FAILED:代付失败
```

#### 正确响应数据说明

响应结果response.data数据说明

| 字段 | 变量名 | 类型 | 描述 |
|----|----|---|----|
| Sunny订单编号 | sn | String | 示例: 1120180209xxxxxxxxxxxxxxxxxx 唯一标号，可以用于查询，等操作 |
| 商户Id | out\_trade\_no | String | 商户订单号 如: "1120180209xxxxxxxxxxxxxxxx" |
| appid | appid | String | 示例: 1000258 唯一标号，商户项目ID |
| 下单币种 | fee_type | String | 交易货币 如：INR |
| 实付币种 | cash_fee_type | String | 客户实际代付的币种 如：INR |
| 交易类型 | trade_type | String | 交易类型 如：H5,PAYOUT |
| 支付方交易号 | transaction_id | String | P563xxxxxxxxxxxxx |
| 客户IP | client_ip | String | 如:127.0.0.1 |
| 订单时间 | create_time | Int | 时间戳 如: 1518155270 |
| 支付时间 | finish_time | Int | 成功支付的时间戳 如1518155297 |
| 订单金额 | order_amount | Int | 如:"100" |
| 实际支付金额 | total_fee | Int | 用户需要支付的金额 单位为"分" 如:10 |
| 交易状态 | trade_state | String| PENDING ||
| 异步通知url | notify_url | String | 异步通知url |
| 客户的名称 | customer_name | String | 客户的名称。仅使用字母数字值 |
| 客户电话号码 | customer_phone | String | 客户电话号码 |
| 客户电子邮件 | customer_email | String | 客户电子邮件地址。 |


### 请求示例

#### 代付实例：
请求参数
```
{
    "order_id":"62346857429997272",
    "user_id":"1008677",
    "order_currency":"INR",
    "order_amount":900,
    "order_name":"order_name",
    "notify_url":"https://xxx.com/api/v1/notify",
    "customer_name":"ericlu",
    "customer_phone":"903683xxxx",
    "customer_email":"ericlu@gmail.com",
    "device_info": "device_infoxxxx",
    "order_note":"order_note",
    "ifsc":"ICIC000xxxx",
    "bank_card":"66310170xxxx",
    "pay_type":"bank",
    "address":"F2/251-A SANGAM VIHAR,SOUTH xxxx",
    "city":"DELHI xxxx"
}
```
响应结果
```
{
    "code": 200,
    "msg": "下单成功",
    "data": {
        "id": 2647632,
        "sn": "10202210121713495497389001185123",
        "out_trade_no": "100867710086771120010",
        "appid": "1",
        "fee_type": "INR",
        "cash_fee_type": "INR",
        "transaction_id": "",
        "client_ip": "127.0.0.1",
        "create_time": 1665566029,
        "finish_time": 0,
        "order_amount": 100,
        "total_fee": 100,
        "bank_type": "",
        "notify_url": "https://xxx.com/api/v1/notify",
        "trade_state": "PENDING",
        "customer_name": "ericlu",
        "customer_email": "ericlu@gmail.com",
        "customer_phone": "903683xxxx"
    }
}
```
### 2.3 代收回调


### CALLBACK:

  * #### 请求
    ``` URL: 支付下单中的url地址```  
    ``` Method: POST```  
    ``` Content-Type:  application/json```
  *  #### header
    参考[header](#header)部分

#### body Parameters 请求参数


回调参数

| 字段 | 变量名 | 类型 | 描述 |
|----|----|---|----|
| 发送版本 | version | String | Sunny回调版本(商户可忽略) |
| Sunny订单号 | sn | String | Sunny生成的订单号，可用于查询订单 |
| 支付方交易号 | transaction_id | String | 如: P5631VZG299QZN94JD |
| 商户订单号 | out\_trade\_no | String | 商户订单号 如: "11201802091347484054542598" |
| 下单标价币种 | fee_type | String | 币种 如：INR |
| 订单金额 | total_fee | Int | 订单金额 如：20 |
| 支付货币类型 | cash_fee_type | String | 交易货币 如：INR |
| 代付金额 | cash_fee | Int | 用户代付的金额，单位为"分" 如：20 |
| 交易类型 | trade_type | String | 如: NATIVE |
| 附加数据 | attach | String | 如"7FB42F08C85670A86431xxxxxxxxxxxx",用户提交时候的参数 原路返回 |
| 支付完成时间 | finish_time | int | 如:1665566220 |
| 交易类型 | trade_state | String | 如: PENDING  |
| 代付银行 | bank_type | String | 代付银行编码,如:CFT |
| 支付备注 | order_note | String | 用户提交时候的参数 |
| 商户ID | appid | String | appid,由商户后台获取，或者登录获取 |
| 随机字符串 | nonce_str | String | 随机字符串 如:O2r8GjZ46e |


响应参数
```
SUCCESS
```

1. 同样的通知可能会多次发送给商户系统。商户系统必须能够正确处理重复的通知

2. 后台通知交互时，如果平台收到商户的应答不符合规范或超时，平台会判定本次通知失败，按照机制重新发送通知，

3. 参数接收需使用`POST`的`application/json`形式接收

#### 注意：
1. 同样的通知可能会多次发送给商户系统。商户系统必须能够正确处理重复的通知。 推荐的做法是，当商户系统收到通知进行处理时，先检查对应业务数据的状态，并判断该通知是否已经处理。如果未处理，则再进行处理；如果已处理，则直接返回结果成功。在对业务数据进行状态检查和处理之前，要采用数据锁进行并发控制，以避免函数重入造成的数据混乱。

2. 如果在所有通知没有收到Sunny回调，商户应调用《查询订单信息》接口确认订单信息状态。

#### 特别提醒：
1. 回调 header 中 包含需要签名的参数 （和请求中的是一致的）
2. 商户系统对于开启结果通知的内容一定要做签名验证，并校验通知的信息是否与商户侧的信息一致，防止数据泄露导致出现“假通知”。

### 2.4 代付回调


### CALLBACK:

  * #### 请求
    ``` URL: 支付下单中的url地址```  
    ``` Method: POST```  
    ``` Content-Type:  application/json```
  *  #### header
    参考[header](#header)部分

#### body Parameters 请求参数


回调参数

| 字段 | 变量名 | 类型 | 描述 |
|----|----|---|----|
| Sunny订单号 | sn | String | Sunny生成的订单号，可用于查询订单 |
| 支付方交易号 | transaction_id | String | 如: P5631VZG299QZN94JD |
| 商户订单号 | out\_trade\_no | String | 商户订单号 如: "11201802091347484054542598" |
| 下单标价币种 | fee_type | String | 币种 如：INR |
| 订单金额 | total_fee | Int | 订单金额 如：20 |
| 支付货币类型 | cash_fee_type | String | 交易货币 如：INR |
| 代付金额 | cash_fee | Int | 用户代付的金额，单位为"分" 如：20 |
| 交易类型 | trade_type | String | 如: NATIVE |
| 附加数据 | attach | String | 如"7FB42F08C85670A86431xxxxxxxxxxxx",用户提交时候的参数 原路返回 |
| 支付完成时间 | finish_time | int | 如:1665566220 |
| 交易类型 | trade_state | String | 如: PENDING  |
| 代付银行 | bank_type | String | 代付银行编码,如:CFT |
| 支付备注 | order_note | String | 用户提交时候的参数 |
| 商户ID | appid | String | appid,由商户后台获取，或者登录获取 |
| 随机字符串 | nonce_str | String | 随机字符串 如:O2r8GjZ46e |


响应参数
```
SUCCESS
```

1. 同样的通知可能会多次发送给商户系统。商户系统必须能够正确处理重复的通知

2. 后台通知交互时，如果平台收到商户的应答不符合规范或超时，平台会判定本次通知失败，按照机制重新发送通知，

3. 参数接收需使用`POST`的`application/json`形式接收

#### 注意：
1. 同样的通知可能会多次发送给商户系统。商户系统必须能够正确处理重复的通知。 推荐的做法是，当商户系统收到通知进行处理时，先检查对应业务数据的状态，并判断该通知是否已经处理。如果未处理，则再进行处理；如果已处理，则直接返回结果成功。在对业务数据进行状态检查和处理之前，要采用数据锁进行并发控制，以避免函数重入造成的数据混乱。

2. 如果在所有通知没有收到Sunny回调，商户应调用《查询订单信息》接口确认订单信息状态。

#### 特别提醒：
1. 回调 header 中 包含需要签名的参数 （和请求中的是一致的）
2. 商户系统对于开启结果通知的内容一定要做签名验证，并校验通知的信息是否与商户侧的信息一致，防止数据泄露导致出现“假通知”。

### 2.5 代收查询

### Api:

  * #### 请求
    ``` URL: /pay/queryOrder```  
    ``` Method: GET```  
    ``` Content-Type:  application/json```
  *  #### header
    参考[header](#header)部分

#### url Parameters 请求参数

字段|变量名|必填|类型|描述
----|----|----|----|----
订单编号|sn|是|String|sn Sunny 系统订单号

```
  url 以 /pay/queryOrder?sn=xxxx 的方式请求
```


#### 返回结果

字段|变量名|必填|类型|描述
----|----|----|----|----
返回码|code|是|String(32)|返回码，请参考返回码表
返回信息|msg|是|String(256)|返回信息，成功信息或错误信息
返回数据|data|否|Array/String|返回数据集或其他提示信息

###### 如果code=200,data参数：

| 字段 | 变量名 | 类型 | 描述 |
|----|----|---|----|
| 订单ID | id | Int | 如:10357 |
| Sunny订单编号 | sn | String | 示例: 1120180209xxxxxxxxxxxxxxxxxx 唯一标号，可以用于查询，等操作 |
| 商户Id | out\_trade\_no | String | 商户订单号 如: "1120180209xxxxxxxxxxxxxxxx" |
| appid | appid | String | 示例: 1000258 唯一标号，商户项目ID |
| 下单币种 | fee_type | String | 交易货币 如：INR |
| 实付币种 | cash_fee_type | String | 客户实际代付的币种 如：INR |
| 交易类型 | trade_type | String | 交易类型 如：H5,PAYOUT |
| 支付提供方 | provider | String | 如:sunny,paytm,cashfree| 
| 支付方交易号 | transaction_id | String | P563xxxxxxxxxxxxx |
| 客户IP | client_ip | String | 如:127.0.0.1 |
| 支付token | order_token | String | 如:127.0.0.1 |
| 订单时间 | create_time | Int | 时间戳 如: 1518155270 |
| 支付时间 | finish_time | Int | 成功支付的时间戳 如1518155297 |
| 订单金额 | order_amount | Int | 如:"100" |
| 实际支付金额 | total_fee | Int | 用户需要支付的金额 单位为"分" 如:10 |
| 支付链接 | payment_link | String | 如:"" 收银台h5地址 |
| 交易状态 | trade_state | String| PENDING ||
| 异步通知url | notify_url | String | 异步通知url |
| return_url | return_url | String | 支付成功后跳转url |
| 客户的名称 | customer_name | String | 客户的名称。仅使用字母数字值 |
| 客户电话号码 | customer_phone | String | 客户电话号码 |
| 客户电子邮件 | customer_email | String | 客户电子邮件地址。 |


#### 响应示例

```
{
    "code": 200,
    "msg": "查询成功",
    "data": {
        "id": 2647540,
        "sn": "10202209221453474922045457347191",
        "out_trade_no": "62347557421997101",
        "appid": "",
        "fee_type": "INR",
        "cash_fee_type": "INR",
        "trade_type": "H5",
        "provider": "sunny",
        "transaction_id": "OPIN16638296270595132c87",
        "client_ip": "150.242.228.156",
        "order_token": "20220922122348303559801111078865",
        "create_time": 1665362723,
        "finish_time": 0,
        "order_amount": 50,
        "total_fee": 50,
        "payment_link": "https://pg.onion-pay.com/pg/checkout/1712839ba11943ae8aa5c609efb91e25?amount=50&channel=onionpay&email=9036830689%40gmail.com&host=https%3A%2F%2Fpg.onion-payment.com&order_id=OPIN16638296270595132c87&order_token=20220922122348303559801111078865&pay_app_id=405d8df8376445fdb4ee1823cdbdf45e&phone=9036830689&signature=&signature_new=&user_name=ericlu",
        "trade_state": "FAILED",
        "return_url": "https://xxx.com/notify",
        "notify_url": "http://xxx.com/notify",
        "customer_name": "ericlu",
        "customer_email": "ericlu@gmail.com",
        "customer_phone": "903683xxxx"
    }
}

```

### 2.6 代付查询

### Api:

  * #### 请求
    ``` URL: /pay/queryPayout```  
    ``` Method: GET```  
    ``` Content-Type:  application/json```
  *  #### header
    参考[header](#header)部分

#### url Parameters 请求参数

字段|变量名|必填|类型|描述
----|----|----|----|----
订单编号|sn|是|String|sn Sunny 系统订单号

```
  url 以 /pay/queryPayout?sn=xxxx 的方式请求
```


#### 返回结果

字段|变量名|必填|类型|描述
----|----|----|----|----
返回码|code|是|String(32)|返回码，请参考返回码表
返回信息|msg|是|String(256)|返回信息，成功信息或错误信息
返回数据|data|否|Array/String|返回数据集或其他提示信息

###### 如果code=200,data参数：

| 字段 | 变量名 | 类型 | 描述 |
|----|----|---|----|
| 订单ID | id | Int | 如:10357 |
| Sunny订单编号 | sn | String | 示例: 1120180209xxxxxxxxxxxxxxxxxx 唯一标号，可以用于查询，等操作 |
| 商户Id | out\_trade\_no | String | 商户订单号 如: "1120180209xxxxxxxxxxxxxxxx" |
| appid | appid | String | 示例: 1000258 唯一标号，商户项目ID |
| 下单币种 | fee_type | String | 交易货币 如：INR |
| 实付币种 | cash_fee_type | String | 客户实际代付的币种 如：INR |
| 交易类型 | trade_type | String | 交易类型 如：H5,PAYOUT |
| 支付提供方 | provider | String | 如:sunny,paytm,cashfree| 
| 支付方交易号 | transaction_id | String | P563xxxxxxxxxxxxx |
| 客户IP | client_ip | String | 如:127.0.0.1 |
| 订单时间 | create_time | Int | 时间戳 如: 1518155270 |
| 支付时间 | finish_time | Int | 成功支付的时间戳 如1518155297 |
| 订单金额 | order_amount | Int | 如:"100" |
| 实际支付金额 | total_fee | Int | 用户需要支付的金额 单位为"分" 如:10 |
| 交易状态 | trade_state | String| PENDING ||
| 异步通知url | notify_url | String | 异步通知url |
| 客户的名称 | customer_name | String | 客户的名称。仅使用字母数字值 |
| 客户电话号码 | customer_phone | String | 客户电话号码 |
| 客户电子邮件 | customer_email | String | 客户电子邮件地址。 |


#### 响应示例

```
{
    "code": 200,
    "msg": "查询成功",
    "data": {
        "id": 96,
        "sn": "12202209292057347449460002516875",
        "out_trade_no": "62346857428997268",
        "appid": "",
        "fee_type": "INR",
        "cash_fee_type": "",
        "trade_type": "PAYOUT",
        "provider": "sunny",
        "transaction_id": "",
        "client_ip": "::1",
        "create_time": 1665367785,
        "finish_time": 0,
        "order_amount": 900,
        "total_fee": 900,
        "bank_type": "",
        "notify_url": "http://notify-test.co/notify",
        "trade_state": "AUDIT_APPLY",
        "customer_name": "ericlu",
        "customer_email": "ericlu@gmail.com",
        "customer_phone": "903683xxxx"
    }
}

```

### 2.7 商户查询

### Api:

  * #### 请求
    ``` URL: /pay/queryMerchant```  
    ``` Method: GET```  
    ``` Content-Type:  application/json```
  *  #### header
    参考[header](#header)部分

#### url Parameters 请求参数 无



#### 返回结果

字段|变量名|必填|类型|描述
----|----|----|----|----
返回码|code|是|String(32)|返回码，请参考返回码表
返回信息|msg|是|String(256)|返回信息，成功信息或错误信息
返回数据|data|否|Array/String|返回数据集或其他提示信息

###### 如果code=200,data参数：

| 字段 | 变量名 | 类型 | 描述 |
|----|----|---|----|
| 可用余额 | balance | Int | 可用余额 单位为分 |
| 未结算余额 | balance_freeze | Int | 冻结金额 单位为分 |
| 提现中金额 | balance_ing | Int | 提现中金额 单位为分 |
| 待结算余额 | balance_unsettled | Int | - 待结算余额 单位为分 |
| appid | appid | String | 示例: 1000258 唯一标号，商户项目ID |


#### 响应示例

```
{
    "code": 200,
    "msg": "查询成功",
    "data": {
        "balance": 0,
        "balance_freeze": 0,
        "balance_ing": 0
        "balance_unsettled": 0
    }
}

```

### 2.8 银行编码

### 2.8.1印度银行:
| 银行编码（BankCode） |  银行名称（bankName） |
|----------------|----|
| ANDB	           | Andhra Bank |
| AXIS	           | Axis Bank |
| BRDA           | Bank of Baroda |
| BOIN	           | Bank of India |
| MHRT	           | Bank of Maharashtra |
| CANBK           | 	Canara Bank |
| CSYB	           | Catholic Syrian Bank |
| CBOI	           | Central Bank of India |
| CITI	           | Citi Bank |
| CORB	           | Corporation Bank |
| DENB           | Dena Bank |
| DNLM	           | Dhanlaxmi Bank |
| FEDB           | Federal Bank |
| HDFC	           | HDFC Bank |
| ICICI           | 	ICICI Bank |
| IDBI	           | IDBI Bank |
| INDNB           | 	Indian Bank |
| INOB           | Indian Overseas Bank |
| INDU	           | IndusInd Bank |
| JAKB	           | Jammu and Kashmir Bank |
| KARBK           | 	Karnataka Bank |
| KOTBK           | 	Kotak Mahindra Bank |
| OBOC	           | Oriental Bank of Commerce |
| PASB	           | Punjab and Sind Bank |
| PNJB	           | Punjab National Bank |
| SCTRD           | 	Standard Chartered Bank |
| SBBJ	           | State Bank of Bikaner and Jaipur |
| SBOH	           | State Bank of Hyderabad |
| SBIN           | State Bank of India |
| SBOT	           | State Bank of Travancore |
| SIBK	           | South Indian Bank |
| SYNB	           | Syndicate Bank |
| TAMB	           | Tamilnadu Mercantile Bank |
| UCOB	           | UCO Bank |
| UNBOI            | 	Union Bank of India |
| VJYB	           | Vijaya Bank |
| YESB	           | Yes Bank |

### 2.9 支付币种

### 2.9.1印度银行:
| 货币标识          |  货币名称 |
|---------------|----|
| INR	          | 卢比 |

## Update
- By：peter
- Time:2022.12.09













