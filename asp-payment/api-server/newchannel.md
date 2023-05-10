## 添加新渠道事项

### 管理后台
- 添加渠道
- 添加渠道商户
- 添加渠道配置列表 配置相关key 状态为正常
- 添加渠道支付配置 开启对用的支付方式 H5 Payout Wappay......


### 接口添加
- config.yaml 添加配置文件
- common/pkg/config/config.go 添加对应的结构体 配置
- common/service/supplier/config.go 添加对应渠道映射关系
- common/service/supplier/config.go => SupplierErrorMap 添加成功时候的对应的字符串映射
- common/supplier/impl 复制一份目录实现对应的方法

### 渠道内部实现
- constant.go 地址变更
- model.go 返回参数结构体解析