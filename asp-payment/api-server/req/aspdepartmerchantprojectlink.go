package req

type DepartMerchantProjectLinkList struct {
	Id           uint   `json:"id"`
	DepartId     int    `json:"depart_id"`
	MchId        int    `json:"mch_id"`
	MchProjectId int    `json:"mch_project_id"`
	Status       int    `json:"status"`
	Sort         int    `json:"sort"`
	CreateTime   uint64 `json:"create_time"`
	UpdateTime   uint64 `json:"update_time"`
}
