# 小张私厨后端服务

## 项目介绍
> `小张私厨`是送给女朋友的生日礼物，为后端源代码。
> 项目使用 `golang gin` 框架开发，用`简易的jwt`做登录授权。

功能说明

- 数据验证 [validator v10](github.com/go-playground/validator)
- 路由，中间件 [go-gin](github.com/gin-gonic/gin)
- 数据库 [gorm](github.com/jinzhu/gorm)
- ~~文档 swagger~~
- 身份鉴权 自写简易jwt
- 热编译 [realize](https://github.com/oxequa/realize) 

## 快速开始

> 本操作在macos 、 linux 下生效，需要golang 1.11+  编译环境,设置git clone 权限

```
git clone git@github.com:54853315/weapp-MrsZhangPrivateKitchen-Backend.git
export GOPROXY=https://goproxy.cn
export GO111MODULE=on
#后端编译
go build -o FoodBackend
```


## 软链到$GOPATH/src

如果想把这个项目放到`GOPATH`下面，不使用go mod模式的话，只需要把这个项目移到`GOPATH`环境变量包含的任意一个目录下面的src目录里，就可以启用`GOPATH`模式了。

>前提是 GOMODULE 环境变量的值必须是auto 或 off

``` shell
cd FoodBackend
ln -s $(PWD) ~/go/src/
```

## 数据移值

```bash
# 执行 sql 语句
mysql> source ./scripts/init.sql;
```


## 运行方式

可以直接用`go run` 启动程序，也可以使用热编译`realize`命令启动。

### go run

`go run mian.go`

### realize

或 `.realize.yaml` 中的`schema.path`为项目当前路径后：

`realize start`

## 访问

`http://localhost:8080/api/`


# 演示 Demo

（暂未完成开发...)

## 贡献代码

非常欢迎优秀的开发者来贡献。在提Pull Request之前，请首先阅读源码，了解原理和架构。