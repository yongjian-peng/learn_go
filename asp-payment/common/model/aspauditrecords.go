package model

// AspAuditRecords 审核记录表
type AspAuditRecords struct {
	Id            int    `gorm:"column:id;type:int(11);primary_key;AUTO_INCREMENT;comment:主键" json:"id"`
	AuditType     int    `gorm:"column:audit_type;type:tinyint(1);default:1;comment:审核类型（1代付）;NOT NULL" json:"audit_type"`
	ProjectId     int    `gorm:"column:project_id;type:int(11);default:0;comment:关联的id主键;NOT NULL" json:"project_id"`
	Sn            string `gorm:"column:sn;type:varchar(45);comment:系统订单号;NOT NULL" json:"sn"`
	DepartId      int    `gorm:"column:depart_id;type:int(11);default:0;comment:asp（平台）商户id 关联 departs表id;NOT NULL" json:"depart_id"`
	ChannelId     uint   `gorm:"column:channel_id;type:int(10) unsigned;default:0;comment:渠道id;NOT NULL" json:"channel_id"`
	MchId         uint   `gorm:"column:mch_id;type:int(10) unsigned;default:0;comment:外部商家账户id;NOT NULL" json:"mch_id"`
	Status        int    `gorm:"column:status;type:tinyint(1);default:0;comment:审核状态;NOT NULL" json:"status"`
	OperateId     int    `gorm:"column:operate_id;type:int(11);default:0;comment:操作的用户id;NOT NULL" json:"operate_id"`
	OperateRemark string `gorm:"column:operate_remark;type:varchar(255);comment:备注;NOT NULL" json:"operate_remark"`
	OperateTime   int64  `gorm:"column:operate_time;type:bigint(20);default:0;comment:审核时间;NOT NULL" json:"operate_time"`
}

func (AspAuditRecords) TableName() string {
	return "asp_audit_records"
}
