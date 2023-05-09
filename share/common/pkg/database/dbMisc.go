package database

import "database/sql"

// Pagination 分页请求数据
type Pagination struct {
	PageNum  int `json:"pageNum"`  // 页码
	PageSize int `json:"pageSize"` // 每页条数
	Total    int `json:"total"`    // 总数据条数
}

func (p *Pagination) Offset() int {
	offset := 0
	if p.PageNum > 0 {
		offset = (p.PageNum - 1) * p.PageSize
	}
	return offset
}

func (p *Pagination) TotalPage() int {
	if p.Total == 0 || p.PageSize == 0 {
		return 0
	}
	totalPage := p.Total / p.PageSize
	if p.Total%p.PageSize > 0 {
		totalPage = totalPage + 1
	}
	return totalPage
}

type QueryParams struct {
	Query string        // 查询
	Args  []interface{} // 参数
}

// OrderBy 排序信息
type OrderBy struct {
	Column string // 排序字段
	Order  string // ASC | DESC
}

func SqlNullString(value string) sql.NullString {
	return sql.NullString{
		String: value,
		Valid:  len(value) > 0,
	}
}
