package admin

import (
	"context"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/samber/lo"
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

var adminUserService *AdminUserService
var AdminUserServiceOnce sync.Once

type AdminUserService struct {
	*Service
}

func AdminUser(c *fiber.Ctx) *AdminUserService {
	return &AdminUserService{Service: NewService(c, "adminService")}
}

// Login 绑定用户
func (s *AdminUserService) Login() error {
	//绑定参数并校验
	reqBody := new(req.AdminUserLogin)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//查询用户名是否存在
	adminUser, _ := repository.AdminUserRepository.Take("username = ?", reqBody.Username)
	if adminUser == nil {
		return s.Error(appError.NewError("账号不存在"))
	}

	hashPassword := s.getHashPassword(reqBody.Password, adminUser.Salt)

	if adminUser.Password != hashPassword {
		return s.Error(appError.NewError("密码错误"))
	}

	if adminUser.Status == 0 {
		return s.Error(appError.NewError("账号已被禁用"))
	}

	//密码正确 获取token
	accessToken := s.createToken(adminUser.Id, hashPassword)
	//设置token到redis
	s.Redis().Set(context.Background(), goRedis.GetKey(fmt.Sprintf("%s:%s", constant.AdminAccessToken, accessToken)), adminUser.Id, time.Hour*72)

	return s.Success(map[string]interface{}{
		"access_token": accessToken,
		"userInfo": rsp.LoginAdminUser{
			Id:       adminUser.Id,
			Username: adminUser.Username,
			Avatar:   adminUser.Avatar,
		},
	})
}

// Logout 登出
func (s *AdminUserService) Logout() error {
	tokenKey := goRedis.GetKey(fmt.Sprintf("access_token:%s", s.AccessToken))
	s.Redis().Del(context.Background(), tokenKey)
	return s.Success(nil)
}

// Options 登出
func (s *AdminUserService) Options() error {

	return s.Success(map[string]any{
		"status": s.GetStatusOptions(),
	})
}

func (s *AdminUserService) GetStatusOptions() []*rsp.Option {
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

func (s *AdminUserService) StatusOptions() error {
	return s.Success(s.GetStatusOptions())
}

func (s *AdminUserService) getHashPassword(password, salt string) string {
	return goutils.Md5(fmt.Sprintf("%s%s", password, salt))
}

func (s *AdminUserService) createToken(id int, password string) string {
	return goutils.Md5(fmt.Sprintf("%d%s%d", id, password, goutils.GetCurTimeMillisecond()))
}

// List 用户列表
func (s *AdminUserService) List() error {
	//绑定参数并校验
	reqBody := new(req.AdminUserList)
	if err := s.ValidateGet(reqBody); err != nil {
		return s.Error(err)
	}

	//构造搜索条件
	cdn := database.SqlCdn()
	if reqBody.UserName != "" {
		cdn.Like("username", reqBody.UserName)
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
	adminUserList, pagination, _ := repository.AdminUserRepository.FindPageByCdn(cdn.Page(reqBody.PageNum, reqBody.PageSize))
	rspPageResult := &rsp.PageResult{}
	rspPageResult.Page = pagination

	rspAdminUserList := make([]*rsp.AdminUser, 0)
	for _, adminUser := range adminUserList {
		rspAdminUserList = append(rspAdminUserList, s.getRspAdminUser(adminUser))
	}
	rspPageResult.List = rspAdminUserList

	return s.Success(rspPageResult)
}

func (s *AdminUserService) getRspAdminUser(adminUser *model.AdminUser) *rsp.AdminUser {
	return &rsp.AdminUser{
		Id:         adminUser.Id,
		Username:   adminUser.Username,
		Avatar:     adminUser.Avatar,
		Status:     adminUser.Status,
		Roles:      s.GetUserRoleIds(adminUser.Id),
		CreateTime: goutils.TimeFormat("Y-m-d H:i:s", adminUser.CreateTime),
	}
}

// Create 创建用户
func (s *AdminUserService) Create() error {
	//绑定参数并校验
	reqBody := new(req.CreateAdminUser)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//查询名称是否已经存在
	name := reqBody.Username
	nameAdminUser, _ := repository.AdminUserRepository.Take("username = ?", name)
	if nameAdminUser != nil {
		return s.Error(appError.NewError("用户名已经存在"))
	}

	salt := goutils.RandomString(6)
	password := s.getHashPassword(reqBody.Password, salt)

	//新增
	adminUser := &model.AdminUser{
		Username:   reqBody.Username,
		Password:   password,
		Salt:       salt,
		Avatar:     reqBody.Avatar,
		Status:     reqBody.Status,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err := repository.AdminUserRepository.Create(adminUser)
	if err != nil {
		return s.Error(appError.NewError("创建失败，请重试"))
	}

	return s.Success(nil)
}

// Edit 更新
func (s *AdminUserService) Edit() error {
	//绑定参数并校验
	reqBody := new(req.UpdateAdminUser)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminUser, _ := repository.AdminUserRepository.Take("id = ?", reqBody.Id)
	if adminUser == nil {
		return s.Error(appError.NewError("账号不存在"))
	}

	if adminUser.Id == 1 {
		adminUser.Status = 1 //超级管理员账号不能被禁用
	}

	//查询名称是否已经存在,非自己
	nameAdminUser, _ := repository.AdminUserRepository.Take("username = ? and id <> ?", reqBody.Username, reqBody.Id)
	if nameAdminUser != nil {
		return s.Error(appError.NewError("账号名称已经存在"))
	}

	data := map[string]interface{}{
		"username": reqBody.Username,
		"avatar":   reqBody.Avatar,
		"status":   reqBody.Status,
	}

	//密码不为空时，修改密码
	if reqBody.Password != "" {
		//重新生成密码
		salt := adminUser.Salt
		password := s.getHashPassword(reqBody.Password, salt)
		data["password"] = password
	}

	//更新信息
	err := repository.AdminUserRepository.Updates(data, "id = ?", reqBody.Id)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// ResetPwd 重置密码
func (s *AdminUserService) ResetPwd() error {
	//绑定参数并校验
	reqBody := new(req.ResetAdminUserPwd)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminUser, _ := repository.AdminUserRepository.Take("id = ?", reqBody.Id)
	if adminUser == nil {
		return s.Error(appError.NewError("账号不存在"))
	}
	//重新生成密码
	salt := adminUser.Salt
	password := goutils.Md5(fmt.Sprintf("%s%s", "123456", salt))
	data := map[string]interface{}{
		"password": password,
	}

	//更新信息
	err := repository.AdminUserRepository.Updates(data, "id = ?", reqBody.Id)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// ChangeStatus 切换状态
func (s *AdminUserService) ChangeStatus() error {
	//绑定参数并校验
	reqBody := new(req.SetAdminUserStatus)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminUser, _ := repository.AdminUserRepository.Take("id = ?", reqBody.Id)
	if adminUser == nil {
		return s.Error(appError.NewError("账号不存在"))
	}

	if adminUser.Id == 1 {
		return s.Error(appError.NewError("不能修改超级管理员账号状态"))
	}

	data := map[string]interface{}{
		"status": reqBody.Status,
	}

	//更新信息
	err := repository.AdminUserRepository.Updates(data, "id = ?", reqBody.Id)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// SetRole 设置角色
func (s *AdminUserService) SetRole() error {
	//绑定参数并校验
	reqBody := new(req.SetAdminUserRole)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminUser, _ := repository.AdminUserRepository.Take("id = ?", reqBody.Id)
	if adminUser == nil {
		return s.Error(appError.NewError("账号不存在"))
	}

	roles := reqBody.Roles
	//在事务中进行
	txErr := database.Db.Transaction(func(tx *gorm.DB) error {

		adminUserRoleRepository := repository.NewRepository[*model.AdminUserRole](tx, goRedis.Redis)

		preUserRoles, _ := adminUserRoleRepository.Find(database.SqlCdn().Eq("admin_id", adminUser.Id))

		preUserRoleIds := make([]int, 0)
		for _, preUserRole := range preUserRoles {
			preUserRoleIds = append(preUserRoleIds, preUserRole.RoleId)
		}

		//原角色，不在新角色，刪除
		for _, preUserRole := range preUserRoles {
			if !goutils.InSlice[int](preUserRole.RoleId, roles) {
				err := adminUserRoleRepository.Delete(&model.AdminUserRole{}, "admin_id = ? and role_id = ?", adminUser.Id, preUserRole.RoleId)
				if err != nil {
					return err
				}
			}
		}

		//添加用户角色
		for _, roleId := range roles {

			//新角色id 在原id中，不添加
			if goutils.InSlice[int](roleId, preUserRoleIds) {
				continue
			}

			//非超级管理员，但是设置超级管理员角色，忽略添加
			if adminUser.Id != 1 && roleId == 1 {
				continue
			}

			adminUserRole := &model.AdminUserRole{
				AdminId:    adminUser.Id,
				RoleId:     roleId,
				CreateTime: time.Now(),
				UpdateTime: time.Now(),
			}
			err := adminUserRoleRepository.Create(adminUserRole)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if txErr != nil {
		return s.Error(appError.NewError("设置失败，请重试"))
	}

	return s.Success(nil)
}

// Delete 删除
func (s *AdminUserService) Delete() error {
	//绑定参数并校验
	reqBody := new(req.DeleteAdminUser)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}
	adminRole, _ := repository.AdminUserRepository.Take("id = ?", reqBody.Id)
	if adminRole == nil {
		return s.Error(appError.NewError("账号不存在"))
	}

	if reqBody.Id == 1 {
		return s.Error(appError.NewError("超级管理员无法删除"))
	}

	//在事务中进行
	txErr := database.Db.Transaction(func(tx *gorm.DB) error {

		adminUserRepository := repository.NewRepository[*model.AdminUser](tx, goRedis.Redis)
		adminUserRoleRepository := repository.NewRepository[*model.AdminUserRole](tx, goRedis.Redis)

		//删除用户
		err := adminUserRepository.Delete(&model.AdminUser{}, "id = ?", reqBody.Id)
		if err != nil {
			return err
		}

		//删除用户角色
		err = adminUserRoleRepository.Delete(&model.AdminUserRole{}, "admin_id = ?", reqBody.Id)
		if err != nil {
			return err
		}

		return nil
	})

	if txErr != nil {
		return s.Error(appError.NewError("删除失败，请重试"))
	}

	return s.Success(nil)
}

// GetUserRoleIds 获取用户角色Ids
func (s *AdminUserService) GetUserRoleIds(adminId int) []int {

	adminUserRoles, _ := repository.AdminUserRoleRepository.Find(database.SqlCdn().Eq("admin_id", adminId))
	roleIds := make([]int, 0)
	for _, role := range adminUserRoles {
		roleIds = append(roleIds, role.RoleId)
	}
	return roleIds
}

// GetUserPermissionsRoutes 获取用户权限路由
func (s *AdminUserService) GetUserPermissionsRoutes(adminId int) []string {

	//用户所有的角色
	roleIds := s.GetUserRoleIds(adminId)
	//通过角色获取所有角色的权限列表
	adminRolePermissionsList, _ := repository.AdminRolePermissionsRepository.Find(database.SqlCdn().In("role_id", roleIds))

	rolePermissionsIds := make([]int, 0)
	for _, adminRolePermissions := range adminRolePermissionsList {
		rolePermissionsIds = append(rolePermissionsIds, adminRolePermissions.PermissionsId)
	}
	//取唯一权限id
	uniqRolePermissionsIds := lo.Uniq[int](rolePermissionsIds)
	adminPermissions, _ := repository.AdminPermissionsRepository.Find(database.SqlCdn().In("id", uniqRolePermissionsIds))

	permissionsRoutes := make([]string, 0)
	for _, adminPermission := range adminPermissions {
		permissionsRoutes = append(permissionsRoutes, adminPermission.RouteUrl)
	}

	return permissionsRoutes
}

// GetUserMenuButtons 获取用户权限菜单按钮
func (s *AdminUserService) GetUserMenuButtons(adminId int) map[string][]string {

	//超级管理员
	if adminId == 1 {
		//超级管理员直接获取所有菜单
		adminMenus, _ := repository.AdminMenuRepository.Find(database.SqlCdn())
		//获取菜单按钮
		buttons := make(map[string][]string, 0)
		for _, adminMenu := range adminMenus {
			name := adminMenu.Name
			adminMenuAuthBtnList := make([]req.AdminMenuAuthBtn, 0)
			_ = goutils.JsonDecode(adminMenu.AuthBtn, &adminMenuAuthBtnList)
			if len(adminMenuAuthBtnList) > 0 {
				for _, menuAuthBtn := range adminMenuAuthBtnList {
					if nameButtons, ok := buttons[name]; ok {
						nameButtons = append(nameButtons, menuAuthBtn.Name)
						buttons[name] = nameButtons
					} else {
						buttons[name] = []string{
							menuAuthBtn.Name,
						}
					}
				}
			}
		}
		return buttons
	}

	//普通用户所有的角色
	roleIds := s.GetUserRoleIds(adminId)
	//通过角色获取所有角色的菜单列表
	adminRoleMenus, _ := repository.AdminRoleMenuRepository.Find(database.SqlCdn().In("role_id", roleIds))
	//获取菜单按钮
	buttons := make(map[string][]string, 0)
	roleMenuButtons := make(map[int][]string, 0)
	for _, adminRoleMenu := range adminRoleMenus {
		menuId := adminRoleMenu.MenuId //获取当前菜单
		menuBtn := make([]string, 0)   //获取当前角色菜单的权限按钮组
		//判断所有菜单组是否已存在此菜单，如果存在叠加按钮组
		_ = goutils.JsonDecode(adminRoleMenu.Btn, &menuBtn)
		if menuButtons, ok := roleMenuButtons[menuId]; ok {
			menuButtons = append(menuButtons, menuBtn...) //叠加
			menuButtons = lo.Uniq(menuButtons)            //去重
			roleMenuButtons[menuId] = menuButtons
		} else {
			roleMenuButtons[menuId] = menuBtn
		}
	}

	//获取
	menuIds := lo.Keys(roleMenuButtons)
	adminMenus, _ := repository.AdminMenuRepository.Find(database.SqlCdn().In("id", menuIds))

	for _, adminMenu := range adminMenus {
		buttons[adminMenu.Name] = roleMenuButtons[adminMenu.Id]
	}

	return buttons
}

// GetUserMenuIds 获取用户权限菜单ids
func (s *AdminUserService) GetUserMenuIds(adminId int) []int {
	//用户所有的角色
	roleIds := s.GetUserRoleIds(adminId)
	//通过角色获取所有角色的菜单列表
	adminRoleMenuList, _ := repository.AdminRoleMenuRepository.Find(database.SqlCdn().In("role_id", roleIds))

	roleMenuIds := make([]int, 0)
	for _, adminRoleMenu := range adminRoleMenuList {
		roleMenuIds = append(roleMenuIds, adminRoleMenu.MenuId)
	}
	//取唯一权限id
	return lo.Uniq[int](roleMenuIds)
}

// GetUserMenus 获取用户权限菜单
func (s *AdminUserService) GetUserMenus(adminId int) []*model.AdminMenu {

	adminMenus := make([]*model.AdminMenu, 0)
	if adminId == 1 {
		//超级管理员
		adminMenus, _ = repository.AdminMenuRepository.Find(database.SqlCdn())
	} else {
		uniqRoleMenuIds := s.GetUserMenuIds(adminId)
		adminMenus, _ = repository.AdminMenuRepository.Find(database.SqlCdn().In("id", uniqRoleMenuIds))
	}

	return adminMenus
}
