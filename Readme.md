## 1. 构建方式
### 构建server
```shell
go build -o server ./myServer/*.go
```
### 构建client
```shell
go build -o client ./myClient/*.go
```
<br>

## 2. 实现功能
- 用户上线提示
- 展示在线用户
```shell
    show online
```

- 用户重命名
```shell
    rename>XXX
```

- 公聊模式
- 私聊模式
```shell
    to>XXX>Msg
```

## 其它
- 默认服务端ip端口： **127.0.0.1:8989** ,可在 **myServer/main.go** 中进行更改
- 客户端同样默认连接 **127.0.0.1:8989** 服务端, 可在启动客户端时添加命令进行更改：
```shell
    ./client -ip 指定IP -port 指定端口
```
- 私聊信息直接使用 **>** 符号进行分割，导致私聊时信息中若出现多余的 **>** 符号会发送失败