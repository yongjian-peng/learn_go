package constant

const (
	AUDIT_RECORDS_AUDIT_TYPE_ORDER          = 1 // 代收
	AUDIT_RECORDS_AUDIT_TYPE_PAYOUT         = 2 // 代付
	AUDIT_RECORDS_AUDIT_TYPE_STATUS_PASS    = 1 // 通过
	AUDIT_RECORDS_AUDIT_TYPE_STATUS_RETURN  = 2 // 退回
	AUDIT_RECORDS_AUDIT_TYPE_STATUS_WAITING = 3 // 等待
	AUDIT_RECORDS_AUDIT_TYPE_STATUS_DOING   = 4 // 进行中

	CHANNEL_CONFIG_STATUS_CLOSED = 0 // 关闭
	CHANNEL_CONFIG_STATUS_NORMAL = 1 // 正常

	DEPART_STATUS_NOT_APPLY = 0 // 草稿 待提交审核(填写资料未提交审核)
	DEPART_STATUS_NORMAL    = 1 // 正常(商户已审核 & 渠道已申请)
	DEPART_STATUS_CLOSED    = 3 // 已关闭
	DEPART_STATUS_RETURNED  = 4 // 已退回

	ID_STATUS_CLOSED = 0 // 关闭
	ID_STATUS_NORMAL = 1 // 正常

	IP_TYPE_MERCHANT         = "merchant"
	IP_TYPE_MERCHANT_PROJECT = "merchant_project"

	MERCHANT_PROJECT_USER_STATUS_NORMAL = 1 // cp 产品用户状态 正常
	MERCHANT_PROJECT_USER_STATUS_CLOSED = 0 // cp 产品用户状态 关闭

	ORDER_FEE_TYPE_INR           = "INR"
	ORDER_TRADE_STATE_PENDING    = "PENDING"    // 未支付
	ORDER_TRADE_STATE_SUCCESS    = "SUCCESS"    // 支付成功
	ORDER_TRADE_STATE_USERPAYING = "USERPAYING" // 支付中 需要输入密码 等待结果
	ORDER_TRADE_STATE_PAYERROR   = "PAYERROR"   // 支付异常
	ORDER_TRADE_STATE_FAILED     = "FAILED"     // 支付失败
	ORDER_REGION                 = "india"      // 地区
	ORDER_INFO                   = "asp_order"

	PAYOUT_FEE_TYPE_INR           = "INR"
	PAYOUT_TRADE_STATE_PENDING    = "PENDING"    // 未提现
	PAYOUT_TRADE_STATE_USERPAYING = "USERPAYOUT" // 提现中 需要输入密码 等待结果

	PAYOUT_TRADE_STATE_APPLY           = "APPLY"           // 0 待审核
	PAYOUT_TRADE_STATE_RETURN          = "RETURN"          // 1:拒绝审核 （管理员拒绝）
	PAYOUT_TRADE_STATE_FREEZE_SUCCESS  = "FREEZE_SUCCESS"  // 2:冻结成功 方法 事务处理 更改余额 修改订单状态
	PAYOUT_TRADE_STATE_CHANNEL_PENDING = "CHANNEL_PENDING" // 3:请求上游成功 返回支付中 更新 刷新订单
	PAYOUT_TRADE_STATE_CHANNEL_FAILED  = "CHANNEL_FAILED"  // 4:请求上游成功 返回失败 更新
	PAYOUT_TRADE_STATE_CHANNEL_SUCCESS = "CHANNEL_SUCCESS" // 5:请求上游成功 返回成功 更新
	PAYOUT_TRADE_STATE_SUCCESS         = "SUCCESS"         // 6:解冻成功 + 代付成功 function {包含事务，}
	PAYOUT_TRADE_STATE_FAILED          = "FAILED"          // 7:解冻成功 + 代付失败 function {包含事务，}
	PAYOUT_TRADE_STATE_REVOKE          = "REVOKE"          // 8:解冻成功 + 代付取消 （管理员取消）function {包含事务，}
	PAYOUT_IS_CHECKOUT_APPLY           = 0                 // 申请 初始化
	PAYOUT_IS_CHECKOUT_SUCCESS         = 1                 // 已完成
	PAYOUT_IS_CHECKOUT_PENDING         = 2                 // 进行中

	MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYORDER               = 1 // 业务类型 1:代收 asp_order
	MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_FREEZE          = 2 // 业务类型 2:代付冻结 asp_payout
	MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_RECHARGE               = 3 // 业务类型 3:充值
	MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_SETTLEMENT             = 4 // 业务类型 4:月结算 asp_merchant_project_month_settlement
	MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_PAYOUT_UNFREEZE        = 5 // 业务类型 5:代付解冻 asp_payout
	MERCHANT_PROJECT_CAPITAL_FLOW_BUSINESS_TYPE_TO_AVAILABLE_TOTAL_FEE = 6 // 业务类型 6:cp项目待结算余额转到可用余额 asp_merchant_project_transfers_day_flow
	MERCHANT_PROJECT_CAPITAL_FLOW_CASH_TYPE_IN                         = 1 // 类资金方向类型 (1.转入
	MERCHANT_PROJECT_CAPITAL_FLOW_CASH_TYPE_OUT                        = 2 // 资金方向类型 (2.转出

	MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_FREEZE_SUCCESS         = "FREEZE_SUCCESS"           // 具体业务相关备注 FREEZE_SUCCESS
	MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_SUCCESS       = "UNFREEZE_SUCCESS_SUCCESS" // 具体业务相关备注 UNFREEZE_SUCCESS
	MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_REVOKE        = "UNFREEZE_REVOKE"          // 具体业务相关备注 取消 UNFREEZE_REVOKE
	MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_TO_AVAILABLE_TOTAL_FEE = "TO_AVAILABLE_TOTAL_FEE"   // 具体业务相关备注 cp项目待结算余额转到可用余额中 不写入到数据库中，业务类型判断需要用到的
	MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_UNFREEZE_FAILED        = "UNFREEZE_SUCCESS_FAILED"  // 具体业务相关备注 失败 UNFREEZE_FAILED
	MERCHANT_PROJECT_CAPITAL_FLOW_REMARK_SUCCESS                = "SUCCESS"                  // 具体业务相关备注 SUCCESS

	MERCHANT_PROJECT_PRE_FLOW_BUSINESS_STATUS_SUCCESS = 1 // 冻结记录锁定的状态(1.锁定中
	MERCHANT_PROJECT_PRE_FLOW_BUSINESS_STATUS_FAILD   = 2 // 冻结记录锁定的状态(2.已解锁

	PARAMS_PAY_TYPE_BANK     = "bank"                 // 代付收款类型 bank 银行卡
	PARAMS_PAY_TYPE_UPI      = "upi"                  // 代付收款类型 upi 印度支付方式
	SYSTEM_PAYOUT_STATUS     = "system_payout_status" // 系统是否开启代付 key
	SYSTEM_PAYOUT_STATUS_OFF = "off"                  // 系统是否开启代付 key 关闭

	SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_CODE    = "sunny-failed"   // code 上游渠道金额和系统金额不一致
	SYSTEMCTL_MONEY_ORDER_TOTAL_IS_DIFF_MESSAGE = "上游渠道金额和系统金额不一致" // message 上游渠道金额和系统金额不一致

	BENEFICIARY_TRADE_STATE_PENDING = "PENDING" // 未支付
	BENEFICIARY_TRADE_STATE_SUCCESS = "SUCCESS" // 支付成功

	CASHIERDESK_UPIPREFIX        = "upi://"  // 支付upi 公共前缀标识
	CASHIERDESK_PAY_TYPE_UPI     = "upi"     // 支付类型 upi
	CASHIERDESK_PAY_TYPE_GPAY    = "gpay"    // 支付类型 gpay
	CASHIERDESK_PAY_TYPE_PAYTM   = "paytm"   // 支付类型 paytm
	CASHIERDESK_PAY_TYPE_PHONEPE = "phonepe" // 支付类型 phonepe
)
