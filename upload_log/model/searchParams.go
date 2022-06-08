package model

type PageInfo struct {
	Page     int    `json:"page" form:"page"`         // 页码
	PageSize int    `json:"pageSize" form:"pageSize"` // 每页大小
	Keyword  string `json:"keyword" form:"keyword"`   //关键字
}

type SearchParams struct {
	ProdLogUploadResultModel
	PageInfo
	Desc bool `json:"desc"` // 排序方式:升序false(默认)|降序true
}
