package admin

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"share/common/model"
	"share/common/pkg/appError"
	"share/common/pkg/constant"
	"share/common/pkg/database"
	"share/common/pkg/goRedis"
	"share/common/pkg/goutils"
	"share/common/repository"
	"share/common/req"
	"share/common/rsp"
	"sync"
	"time"
)

var appService *AppService
var AppServiceOnce sync.Once

type AppService struct {
	*Service
}

func App(c *fiber.Ctx) *AppService {
	return &AppService{Service: NewService(c, "appService")}
}

// Options 获取所有options
func (s *AppService) Options() error {
	return s.Success(map[string]any{
		"status": s.GetStatusOptions(),
	})
}

func (s *AppService) GetStatusOptions() []*rsp.Option {
	options := make([]*rsp.Option, 0)
	options = append(options, &rsp.Option{
		Label: "启用",
		Value: constant.AdminUserStatusEnable,
	})
	options = append(options, &rsp.Option{
		Label: "禁用",
		Value: constant.AdminUserStatusDisable,
	})
	return options
}

func (s *AppService) StatusOptions() error {
	return s.Success(s.GetStatusOptions())
}

// List 用户列表
func (s *AppService) List() error {
	//绑定参数并校验
	reqBody := new(req.AdminAppList)
	if err := s.ValidateGet(reqBody); err != nil {
		return s.Error(err)
	}

	//构造搜索条件
	cdn := database.SqlCdn()
	if reqBody.Name != "" {
		cdn.Like("name", reqBody.Name)
	}

	if reqBody.Status != -1 {
		cdn.Eq("status", reqBody.Status)
	}

	if len(reqBody.CreateTime) > 0 {
		cdn.Gte("create_time", reqBody.CreateTime[0])
	}

	if len(reqBody.CreateTime) > 1 {
		cdn.Lte("create_time", reqBody.CreateTime[1])
	}

	//获取列表
	appList, pagination, _ := repository.AppRepository.FindPageByCdn(cdn.Page(reqBody.PageNum, reqBody.PageSize))
	rspPageResult := &rsp.PageResult{}
	rspPageResult.Page = pagination

	rspAppList := make([]*rsp.AdminApp, 0)
	for _, adminUser := range appList {
		rspAppList = append(rspAppList, s.getRspAdminApp(adminUser))
	}
	rspPageResult.List = rspAppList

	return s.Success(rspPageResult)
}

func (s *AppService) getRspAdminApp(adminApp *model.App) *rsp.AdminApp {
	return &rsp.AdminApp{
		Id:          adminApp.Id,
		Name:        adminApp.Name,
		Secret:      adminApp.Secret,
		PackageName: adminApp.PackageName,
		UpdateTime:  goutils.TimeFormat("Y-m-d H:i:s", adminApp.UpdateTime),
		CreateTime:  goutils.TimeFormat("Y-m-d H:i:s", adminApp.CreateTime),
	}
}

// Create 创建用户
func (s *AppService) Create() error {
	//绑定参数并校验
	reqBody := new(req.CreateAdminApp)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//查询名称是否已经存在
	name := reqBody.Name
	nameApp, _ := repository.AppRepository.Take("name = ?", name)
	if nameApp != nil {
		return s.Error(appError.NewError("名称已经存在"))
	}

	salt := goutils.RandomString(6)
	secret := s.getHashSecret(goutils.GetCurTimeStr(), salt)

	//在事务中进行
	txErr := database.Db.Transaction(func(tx *gorm.DB) error {
		adminAppRepository := repository.NewRepository[*model.App](tx, goRedis.Redis)
		adminAppConfigRepository := repository.NewRepository[*model.AppConfig](tx, goRedis.Redis)
		//新增
		adminApp := &model.App{
			Name:        reqBody.Name,
			PackageName: reqBody.PackageName,
			Status:      reqBody.Status,
			Salt:        salt,
			Secret:      secret,
			CreateTime:  time.Now(),
			UpdateTime:  time.Now(),
		}

		err := adminAppRepository.Create(adminApp)
		if err != nil {
			return err
		}

		//添加app config
		appConfig := &model.AppConfig{
			AppId:              adminApp.Id,
			UserDayPayoutLimit: 0,
			PayLimit:           0,
			CreateTime:         time.Now(),
			UpdateTime:         time.Now(),
		}

		appErr := adminAppConfigRepository.Create(appConfig)
		if appErr != nil {
			return appErr
		}

		return nil
	})

	if txErr != nil {
		return s.Error(appError.NewError("创建失败，请重试"))
	}

	return s.Success(nil)
}

func (s *AppService) getHashSecret(curTime, salt string) string {
	return goutils.Md5(fmt.Sprintf("%s%s", curTime, salt))
}

// Edit 更新
func (s *AppService) Edit() error {
	//绑定参数并校验
	reqBody := new(req.UpdateAdminApp)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminApp, _ := repository.AppRepository.Take("id = ?", reqBody.Id)
	if adminApp == nil {
		return s.Error(appError.NewError("应用不存在"))
	}

	//查询名称是否已经存在,非自己
	nameAdminApp, _ := repository.AppRepository.Take("name = ? and id <> ?", reqBody.Name, reqBody.Id)
	if nameAdminApp != nil {
		return s.Error(appError.NewError("应用名称已经存在"))
	}

	data := map[string]interface{}{
		"name":         reqBody.Name,
		"package_name": reqBody.PackageName,
		"status":       reqBody.Status,
	}

	//更新信息
	err := repository.AppRepository.Updates(data, "id = ?", reqBody.Id)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// ResetSecret 重置秘钥
func (s *AppService) ResetSecret() error {
	//绑定参数并校验
	reqBody := new(req.ResetAdminAppSecret)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminApp, _ := repository.AppRepository.Take("id = ?", reqBody.Id)
	if adminApp == nil {
		return s.Error(appError.NewError("应用不存在"))
	}
	//重新生成秘钥
	salt := goutils.RandomString(6)
	secret := goutils.Md5(fmt.Sprintf("%s%s", goutils.GetCurTimeStr(), salt))
	data := map[string]interface{}{
		"secret": secret,
		"salt":   salt,
	}

	//更新信息
	err := repository.AppRepository.Updates(data, "id = ?", reqBody.Id)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// ChangeStatus 切换状态
func (s *AppService) ChangeStatus() error {
	//绑定参数并校验
	reqBody := new(req.SetAdminAppStatus)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminApp, _ := repository.AppRepository.Take("id = ?", reqBody.Id)
	if adminApp == nil {
		return s.Error(appError.NewError("应用不存在"))
	}

	data := map[string]interface{}{
		"status": reqBody.Status,
	}

	//更新信息
	err := repository.AppRepository.Updates(data, "id = ?", reqBody.Id)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// ChangeConfig 修改配置信息
func (s *AppService) ChangeConfig() error {
	//绑定参数并校验
	reqBody := new(req.SetAdminAppConfig)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminAppConfig, _ := repository.AppConfigRepository.Take("app_id = ?", reqBody.AppId)
	if adminAppConfig == nil {
		return s.Error(appError.NewError("应用配置不存在"))
	}

	data := map[string]interface{}{
		"user_day_payout_limit": reqBody.UserDayPayoutLimit,
		"pay_limit":             reqBody.PayLimit,
	}

	//更新信息
	err := repository.AppConfigRepository.Updates(data, "app_id = ?", reqBody.AppId)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}
