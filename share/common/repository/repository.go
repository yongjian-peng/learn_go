package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
	"share/common/model"
	"share/common/pkg/database"
	"share/common/pkg/goRedis"
	"share/common/pkg/goutils"
	"time"
)

var AdminPermissionsRepository *Repository[*model.AdminPermissions]
var AdminRolePermissionsRepository *Repository[*model.AdminRolePermissions]
var AdminMenuRepository *Repository[*model.AdminMenu]
var AdminRoleMenuRepository *Repository[*model.AdminRoleMenu]
var AdminRoleRepository *Repository[*model.AdminRole]
var AdminUserRoleRepository *Repository[*model.AdminUserRole]
var AdminUserRepository *Repository[*model.AdminUser]
var AppRepository *Repository[*model.App]
var AppConfigRepository *Repository[*model.AppConfig]
var UserRepository *Repository[*model.User]

func Init() {
	AdminPermissionsRepository = NewRepository[*model.AdminPermissions](database.Db, goRedis.Redis)
	AdminRolePermissionsRepository = NewRepository[*model.AdminRolePermissions](database.Db, goRedis.Redis)
	AdminMenuRepository = NewRepository[*model.AdminMenu](database.Db, goRedis.Redis)
	AdminRoleRepository = NewRepository[*model.AdminRole](database.Db, goRedis.Redis)
	AdminRoleMenuRepository = NewRepository[*model.AdminRoleMenu](database.Db, goRedis.Redis)
	AdminUserRoleRepository = NewRepository[*model.AdminUserRole](database.Db, goRedis.Redis)
	AdminUserRepository = NewRepository[*model.AdminUser](database.Db, goRedis.Redis)
	AppRepository = NewRepository[*model.App](database.Db, goRedis.Redis)
	AppConfigRepository = NewRepository[*model.AppConfig](database.Db, goRedis.Redis)
	UserRepository = NewRepository[*model.User](database.Db, goRedis.Redis)
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
		PageNum:  cnd.Paging.PageNum,
		PageSize: cnd.Paging.PageSize,
		Total:    count,
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
