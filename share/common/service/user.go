package service

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-module/carbon/v2"
	"github.com/spf13/cast"
	"gorm.io/gorm"
	"share/common/model"
	"share/common/pkg/appError"
	"share/common/pkg/constant"
	"share/common/pkg/database"
	"share/common/pkg/goRedis"
	"share/common/pkg/goutils"
	"share/common/pkg/invitecode"
	"share/common/repository"
	"share/common/req"
	"share/common/rsp"
	"time"
)

type UserService struct {
	*Service
}

func User(c *fiber.Ctx) *UserService {
	return &UserService{Service: NewService(c, "userService")}
}

// Register 注册用户
func (s *UserService) Register() error {
	//绑定参数并校验
	reqBody := new(req.UserRegister)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//根据邀请码获取邀请人信息
	inviteUser, _ := repository.UserRepository.Take("invite_code = ?", reqBody.InviteCode)
	var inviteUid int64 = 0
	if inviteUser != nil {
		inviteUid = inviteUser.Id
	}

	//在事务中进行
	txErr := database.Db.Transaction(func(tx *gorm.DB) error {
		userRepository := repository.NewRepository[*model.User](tx, goRedis.Redis)
		userAccountRepository := repository.NewRepository[*model.UserAccount](tx, goRedis.Redis)
		//新增
		user := &model.User{
			Appid:         cast.ToInt(s.Head.AppId),
			AppUid:        reqBody.AppUid,
			Phone:         reqBody.Phone,
			Status:        constant.UserStatusEnable,
			InviteCode:    "",
			InviteUid:     inviteUid,
			DeviceId:      reqBody.DeviceId,
			RegisterType:  reqBody.RegisterType,
			RegisterTime:  carbon.CreateFromTimestamp(cast.ToInt64(reqBody.RegisterTime)).Carbon2Time(),
			LastLoginTime: time.Now(),
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}

		err := userRepository.Create(user)
		if err != nil {
			return err
		}

		inviteCode := invitecode.Encode(cast.ToUint64(user.Id))
		//修改用户邀请码
		data := map[string]interface{}{
			"invite_code": inviteCode,
		}

		//更新邀请码
		upError := userRepository.Updates(data, "id = ?", user.Id)
		if upError != nil {
			return upError
		}

		//添加app config
		appConfig := &model.UserAccount{
			Uid:                     user.Id,
			TotalRecharge:           0,
			TotalCommission:         0,
			TotalWithdrawCommission: 0,
			Commission:              0,
			FreezeCommission:        0,
			ChildTotalRecharge:      0,
		}

		accountErr := userAccountRepository.Create(appConfig)
		if accountErr != nil {
			return accountErr
		}

		return nil
	})

	if txErr != nil {
		return s.Error(appError.NewError("创建失败，请重试"))
	}

	return s.Success(nil)
}

func (s *UserService) createToken(id int64, salt string) string {
	return goutils.Md5(fmt.Sprintf("%d%s%d", id, salt, goutils.GetCurTimeMillisecond()))
}

// Login 登录
func (s *UserService) Login() error {
	//绑定参数并校验
	reqBody := new(req.UserLogin)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//通过appid ，appUID 找到系统uid
	user, _ := repository.UserRepository.Take("appid = ? and app_uid = ?", s.Head.AppId, reqBody.AppUid)
	if user == nil {
		return s.Error(appError.UserNotRegister)
	}
	appSalt := cast.ToString(s.C.Locals("AppSalt"))
	uid := user.Id //拿到本系统id

	if user.Status == 0 {
		return s.Error(appError.AccountIsDisable)
	}

	//获取token
	accessToken := s.createToken(uid, appSalt)
	//设置token到redis
	s.Redis().Set(ctx, goRedis.GetKey(fmt.Sprintf("%s:%s", constant.ApiAccessToken, accessToken)), uid, time.Hour*72)

	return s.Success(map[string]interface{}{
		"access_token": accessToken,
	})

}

// Info 获取用户信息
func (s *UserService) Info() error {
	userInfo, _ := s.GetCacheUserInfo(s.Uid)
	return s.Success(fiber.Map{
		"userInfo": s.GetRspUserInfo(userInfo),
	})
}

func (s *UserService) GetRspUserInfo(user *model.User) *rsp.User {
	return &rsp.User{
		Id:         user.Id,
		InviteCode: user.InviteCode,
		Level:      user.Level,
	}
}

func (s *UserService) GetUserInfoKey(id int64) string {
	return goRedis.GetKey(fmt.Sprintf("userinfo:%d", id))
}

func (s *UserService) GetCacheUserInfo(id int64) (*model.User, error) {
	key := s.GetUserInfoKey(id)
	return repository.UserRepository.GetCacheInfo(key, func() (*model.User, error) {
		return repository.UserRepository.Take("id = ?", id)
	})
}
