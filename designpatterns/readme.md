## 文章来源
- https://mp.weixin.qq.com/s/9iXdsgtuQh6pge5vSDPoQA
- https://mp.weixin.qq.com/s/NrUWfpyfzX4sag22qF-Y-Q

## go test 命令 demo
- go  test  .\chainofresponsibility.go .\chainofresponsibility_test.go -v  查看日志输出
- go test -run TestChainOfResponsibility -v 查看日志输出
- go test -v -failfast 如果遇到错误的话，则停止执行测试
- go test -short -v  模式允许我们将任何长时间运行的测试标记为在此模式下跳过。

## 目录结构
- chainofresponsibility  责任链模式
- strategy  策略模式  go test .\strategy.go .\abstro.go .\strategy_test.go -v
- abstractfactory 抽象工厂模式  go test -run  TestAbstractFactory -v
- iterator 迭代器模式 go test -run TestIterator -v
- mediator 中介模式 （调解员） go test -run TestMediator -v