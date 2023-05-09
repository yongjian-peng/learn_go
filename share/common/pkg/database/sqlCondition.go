package database

import (
	"share/common/pkg/logger"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type SqlCondition struct {
	SelectCols []string      // 要查询的字段，如果为空，表示查询所有字段
	Params     []QueryParams // 参数
	Orders     []OrderBy     // 排序
	Paging     *Pagination   // 分页
}

func SqlCdn() *SqlCondition {
	return &SqlCondition{}
}

func (s *SqlCondition) Cols(selectCols ...string) *SqlCondition {
	if len(selectCols) > 0 {
		s.SelectCols = append(s.SelectCols, selectCols...)
	}
	return s
}

func (s *SqlCondition) Eq(column string, args ...interface{}) *SqlCondition {
	s.Where(column+" = ?", args)
	return s
}

func (s *SqlCondition) NotEq(column string, args ...interface{}) *SqlCondition {
	s.Where(column+" <> ?", args)
	return s
}

func (s *SqlCondition) Gt(column string, args ...interface{}) *SqlCondition {
	s.Where(column+" > ?", args)
	return s
}

func (s *SqlCondition) Gte(column string, args ...interface{}) *SqlCondition {
	s.Where(column+" >= ?", args)
	return s
}

func (s *SqlCondition) Lt(column string, args ...interface{}) *SqlCondition {
	s.Where(column+" < ?", args)
	return s
}

func (s *SqlCondition) Lte(column string, args ...interface{}) *SqlCondition {
	s.Where(column+" <= ?", args)
	return s
}

func (s *SqlCondition) Like(column string, str string) *SqlCondition {
	s.Where(column+" LIKE ?", "%"+str+"%")
	return s
}

func (s *SqlCondition) Starting(column string, str string) *SqlCondition {
	s.Where(column+" LIKE ?", str+"%")
	return s
}

func (s *SqlCondition) Ending(column string, str string) *SqlCondition {
	s.Where(column+" LIKE ?", "%"+str)
	return s
}

func (s *SqlCondition) In(column string, params interface{}) *SqlCondition {
	s.Where(column+" in (?) ", params)
	return s
}

func (s *SqlCondition) Where(query string, args ...interface{}) *SqlCondition {
	s.Params = append(s.Params, QueryParams{query, args})
	return s
}

func (s *SqlCondition) Asc(column string) *SqlCondition {
	s.Orders = append(s.Orders, OrderBy{Column: column, Order: "ASC"})
	return s
}

func (s *SqlCondition) Desc(column string) *SqlCondition {
	s.Orders = append(s.Orders, OrderBy{Column: column, Order: "DESC"})
	return s
}

func (s *SqlCondition) Limit(limit int) *SqlCondition {
	s.Page(1, limit)
	return s
}

func (s *SqlCondition) Page(page, limit int) *SqlCondition {
	if s.Paging == nil {
		s.Paging = &Pagination{PageNum: page, PageSize: limit}
	} else {
		s.Paging.PageNum = page
		s.Paging.PageSize = limit
	}
	return s
}

func (s *SqlCondition) Build(db *gorm.DB) *gorm.DB {
	ret := db

	if len(s.SelectCols) > 0 {
		ret = ret.Select(s.SelectCols)
	}

	// where
	if len(s.Params) > 0 {
		for _, param := range s.Params {
			ret = ret.Where(param.Query, param.Args...)
		}
	}

	// order
	if len(s.Orders) > 0 {
		for _, order := range s.Orders {
			ret = ret.Order(order.Column + " " + order.Order)
		}
	}

	// limit
	if s.Paging != nil && s.Paging.PageSize > 0 {
		ret = ret.Limit(s.Paging.PageSize)
	}

	// offset
	if s.Paging != nil && s.Paging.Offset() > 0 {
		ret = ret.Offset(s.Paging.Offset())
	}
	return ret
}

func (s *SqlCondition) Find(db *gorm.DB, out interface{}) error {
	if err := s.Build(db).Find(out).Error; err != nil {
		logger.GetLogger("db").Error("db find error", zap.Any("errMsg", err.Error()))
		return err
	}
	return nil
}

func (s *SqlCondition) Take(db *gorm.DB, out interface{}) error {
	if err := s.Build(db).Take(out).Error; err != nil {
		return err
	}
	return nil
}

func (s *SqlCondition) First(db *gorm.DB, out interface{}) error {
	if err := s.Build(db).First(out).Error; err != nil {
		return err
	}
	return nil
}

func (s *SqlCondition) Count(db *gorm.DB, model interface{}) int {
	ret := db.Model(model)

	// where
	if len(s.Params) > 0 {
		for _, query := range s.Params {
			ret = ret.Where(query.Query, query.Args...)
		}
	}

	var count int64
	if err := ret.Count(&count).Error; err != nil {
		logger.GetLogger("db").Error("db find error", zap.Any("errMsg", err.Error()))
	}
	return int(count)
}
