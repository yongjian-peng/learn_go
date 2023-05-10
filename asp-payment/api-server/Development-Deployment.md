
# 项目介绍

# 开发命令
```
go run main.go
```

# 项目编译

需要交叉编译成Linux文件

## cmd下编译
windows下编译linux可执行文件
```
set GOARCH=amd64
set GOOS=linux
go build main.go
go build -o go-admin main.go
```
windows下编译windows可执行文件
```
go build main.go
go build -o go-admin main.go
```
## powershell下编译
windows下编译linux可执行文件
```
$env:GOOS="linux"
$env:GOARCH="amd64"
go build main.go
go build -o go-admin main.go
```
windows下编译windows可执行文件
```
go build main.go
go build -o go-admin main.go
```














