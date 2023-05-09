package model

type UserAccount struct {
	Uid                     int64  `gorm:"column:uid;type:bigint(20);primary_key;comment:用户id" json:"uid"`
	TotalRecharge           uint64 `gorm:"column:total_recharge;type:bigint(20) unsigned;comment:累计充值金额(分);NOT NULL" json:"total_recharge"`
	TotalCommission         uint64 `gorm:"column:total_commission;type:bigint(20) unsigned;comment:累计佣金金额(分);NOT NULL" json:"total_commission"`
	TotalWithdrawCommission uint64 `gorm:"column:total_withdraw_commission;type:bigint(20) unsigned;comment:累计提现佣金(分);NOT NULL" json:"total_withdraw_commission"`
	Commission              uint64 `gorm:"column:commission;type:bigint(20) unsigned;comment:当前可用佣金(分);NOT NULL" json:"commission"`
	FreezeCommission        uint64 `gorm:"column:freeze_commission;type:bigint(20) unsigned;comment:冻结中的佣金(分);NOT NULL" json:"freeze_commission"`
	ChildTotalRecharge      uint64 `gorm:"column:child_total_recharge;type:bigint(20) unsigned;comment:下级累计充值金额(分);NOT NULL" json:"child_total_recharge"`
	TotalInvite             int    `gorm:"column:total_invite;type:int(11);comment:累计邀请人数量;NOT NULL" json:"total_invite"`
}

func (m *UserAccount) TableName() string {
	return "user_account"
}
