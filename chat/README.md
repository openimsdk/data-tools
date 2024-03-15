# MYSQL -> MongoDB
#### 使用对象
* `https://github.com/openimsdk/chat`
* 配合`https://github.com/openimsdk/open-im-server` v3版本
* v1.6版本,v1.6以下版本升级到v1.6及以上版本
* 只用`https://github.com/openimsdk/open-im-server` 自行实现的业务服务器不用转换

```shell
# 编译
go build
# 运行
./chat -c config.yaml
```

输出`run success`表示转换成功，并退出程序