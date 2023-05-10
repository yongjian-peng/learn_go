package repository

import (
	"asp-payment/common/model"
	"asp-payment/common/pkg/database"
	"asp-payment/common/pkg/goRedis"
	"asp-payment/common/pkg/goutils"
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"time"
)

var AspMerchantProjectChannelDepartTradeTypeLinkRepository *Repository[*model.AspMerchantProjectChannelDepartTradeTypeLink]
var AspChannelDepartTradeTypeRepository *Repository[*model.AspChannelDepartTradeType]

func Init() {
	AspMerchantProjectChannelDepartTradeTypeLinkRepository = NewRepository[*model.AspMerchantProjectChannelDepartTradeTypeLink](database.DB, goRedis.Redis)
	AspChannelDepartTradeTypeRepository = NewRepository[*model.AspChannelDepartTradeType](database.DB, goRedis.Redis)
}

func NewRepository[T model.BaseModel](db *gorm.DB, redis *redis.Client) *Repository[T] {
	return &Repository[T]{db: db, redis: redis}
}

type Repository[T model.BaseModel] struct {
	db    *gorm.DB
	redis *redis.Client
}

// Take 获取一条记录，没有指定排序字段
// db.Take(&users, "name <> ? AND age > ?", "name", 20)
func (r *Repository[T]) Take(where ...interface{}) (T, error) {
	var res T
	if err := r.db.Take(&res, where...).Error; err != nil {
		return *new(T), err
	}
	return res, nil
}

func (r *Repository[T]) Find(cnd *database.SqlCondition) ([]T, error) {
	var list []T
	if err := cnd.Find(r.db, &list); err != nil {
		return list, err
	}
	return list, nil
}

func (r *Repository[T]) FindOne(cnd *database.SqlCondition) (T, error) {
	var res T
	if err := cnd.Take(r.db, &res); err != nil {
		return *new(T), err
	}
	return res, nil
}

// First 获取第一条记录（主键升序）
func (r *Repository[T]) First(cnd *database.SqlCondition) (T, error) {
	var res T
	if err := cnd.First(r.db, &res); err != nil {
		return *new(T), err
	}
	return res, nil
}

func (r *Repository[T]) FindPageByCdn(cnd *database.SqlCondition) ([]T, *database.Pagination, error) {
	var list []T
	var pagination *database.Pagination
	if err := cnd.Find(r.db, &list); err != nil {
		return list, pagination, err
	}
	var res T
	count := cnd.Count(r.db, res)
	pagination = &database.Pagination{
		Page:  cnd.Paging.Page,
		Limit: cnd.Paging.Limit,
		Total: count,
	}
	return list, pagination, nil
}

func (r *Repository[T]) Create(t T) (err error) {
	err = r.db.Create(t).Error
	return
}

func (r *Repository[T]) Update(t T) (err error) {
	err = r.db.Save(t).Error
	return
}

func (r *Repository[T]) Updates(columns map[string]interface{}, query interface{}, args ...interface{}) (err error) {
	var t T
	err = r.db.Model(t).Where(query, args...).Updates(columns).Error
	return
}

func (r *Repository[T]) UpdateColumn(name string, value interface{}, query interface{}, args ...interface{}) (err error) {
	var t T
	err = r.db.Model(t).Where(query, args...).UpdateColumn(name, value).Error
	return
}

func (r *Repository[T]) Delete(t T, query interface{}, args ...interface{}) (err error) {
	err = r.db.Where(query, args...).Delete(t).Error
	return
}

func (r *Repository[T]) FindInIds(ids []int) (res []T) {
	if len(ids) == 0 {
		return nil
	}
	r.db.Where("id in (?)", ids).Find(&res)
	return
}

type DbFindCallback[T model.BaseModel] func() (T, error)
type SetInfoToCache func()

func (r *Repository[T]) GetCacheInfo(key string, dbFindCallback DbFindCallback[T]) (T, error) {
	var t T
	ctx := context.Background()
	//查询是否是空值
	nullStr := r.redis.Get(ctx, key).Val()
	if nullStr == "null" {
		return t, errors.New("null")
	}
	//如果key 不存在
	if nullStr == "" {
		//防止缓存失效，大量用户同时请求击穿数据库 先获取数据库的锁
		lockKey := fmt.Sprintf("%s:dblock", key)
		lockFlag := r.redis.SetNX(ctx, lockKey, 1, 30*time.Second).Val() //加锁
		if lockFlag {
			modelInfo, err := dbFindCallback()
			if err == nil {
				// 转成JSON
				jsonStr, errJ := goutils.JsonEncode(modelInfo)
				if errJ == nil {
					r.redis.Set(ctx, key, jsonStr, time.Hour*3)
				}
				r.redis.Del(ctx, lockKey) //解锁
				return modelInfo, nil
			} else {
				r.redis.Set(ctx, key, "null", time.Second*60)
				r.redis.Del(ctx, lockKey) //解锁
				return t, err
			}
		} else {
			//休息0.5s后重新获取
			time.Sleep(time.Millisecond * 500)
			return r.GetCacheInfo(key, dbFindCallback)
		}
	}

	value := r.redis.Get(ctx, key).Val()
	err := goutils.JsonDecode(value, &t)
	if err != nil {
		return t, err
	}

	return t, nil
}
