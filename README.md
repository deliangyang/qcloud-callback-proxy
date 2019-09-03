### 说明

这是一个处理腾讯`QCloud IM`回调通知的服务。解决账号共享，通知回调地址唯一的问题。根据规则分发到不同的环境上。
回调地址uri为`/cb/qcloud/im`，可在配置文件`configs/proxy.toml`中修改。

### 背景与解决方案

针对个人回调通知，做账号ID前缀映射到不同的环境内(dev|test|dev-next|test-next)。
针对群组的回调通知，如果有`Operator_Account`或者`Owner_Account`
不为空且不为`admin`时，也采取ID映射规则，否者会分发到所有的环境。

### IM回调代理服务器配置

> configs/proxy.toml

```toml
uri = "/cb/qcloud/im"             # 回调地址
method = "POST"                         # 回调请求方法
port = ":3000"                          # 代理服务器端口

[[envs]]
id = "20"                               # 用户ID前缀
url = "https://dev-ck.example.com"      # 回调host

[[envs]]
id = "30"
url = "https://test-ck.example.com"
```

### 准备工作

用户ID是8位，现在增加两位前缀，20表示dev/new-dev，30表示test/new-test环境。

搭建该服务

1. 为new-test环境新建一个passport的数据库
2. 清空test环境passport的数据，
3. 清空users表中的数据，使用新的数据



### Run
```
go run cmd/proxy/main.go -config=configs/proxy.toml
```

### Build
```
go build cmd/proxy/main.go
```

### Service Type
api-service