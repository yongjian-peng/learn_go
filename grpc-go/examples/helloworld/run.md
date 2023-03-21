## 启动服务
- 下载插件 go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28
- go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
- 更新系统环境变量 export PATH="$PATH:$(go env GOPATH)/bin"
- 下载代码 git clone -b v1.53.0 --depth 1 https://github.com/grpc/grpc-go
- 进入目录 cd grpc-go/examples/helloworld
- protoc 命令 protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative helloworld/helloworld.proto
- 启动server 服务 在 examples/helloworld go run greeter_server/main.go
- 启动client 服务 在 examples/helloworld go run greeter_client/main.go

## 相关引用地址：
- https://grpc.io/docs/languages/go/quickstart/