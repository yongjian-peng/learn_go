package admin

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"share/common/model"
	"share/common/pkg/appError"
	"share/common/pkg/database"
	"share/common/pkg/goRedis"
	"share/common/pkg/goutils"
	"share/common/repository"
	"share/common/req"
	"share/common/rsp"
	"sync"
	"time"
)

var adminRoleService *AdminRoleService
var AdminRoleServiceOnce sync.Once

type AdminRoleService struct {
	*Service
}

func AdminRole(c *fiber.Ctx) *AdminRoleService {
	return &AdminRoleService{Service: NewService(c, "adminRoleService")}
}

func (s *AdminRoleService) Options() error {

	//获取列表
	adminRoleList, _ := repository.AdminRoleRepository.Find(database.SqlCdn())
	options := make([]*rsp.Option, 0)

	for _, adminRole := range adminRoleList {
		options = append(options, &rsp.Option{
			Label:    adminRole.Name,
			Value:    adminRole.Id,
			Disabled: goutils.IfBool(adminRole.Id == 1, true, false),
		})
	}

	return s.Success(options)
}

// List 列表
func (s *AdminRoleService) List() error {
	//绑定参数并校验
	reqBody := new(req.AdminRoleList)
	if err := s.ValidateGet(reqBody); err != nil {
		return s.Error(err)
	}
	//构造搜索条件
	cdn := database.SqlCdn()
	if reqBody.Name != "" {
		cdn.Like("name", reqBody.Name)
	}

	if len(reqBody.CreateTime) > 0 {
		cdn.Gte("create_time", reqBody.CreateTime[0])
	}

	if len(reqBody.CreateTime) > 1 {
		cdn.Lte("create_time", reqBody.CreateTime[1])
	}

	cdn.Page(reqBody.PageNum, reqBody.PageSize)
	//获取列表
	adminRoleList, pagination, _ := repository.AdminRoleRepository.FindPageByCdn(cdn)
	rspPageResult := &rsp.PageResult{}
	rspPageResult.Page = pagination

	rspAdminRoleList := make([]*rsp.AdminRole, 0)
	for _, adminRole := range adminRoleList {
		rspAdminRoleList = append(rspAdminRoleList, s.getRspAdminRole(adminRole))
	}
	rspPageResult.List = rspAdminRoleList

	return s.Success(rspPageResult)
}

func (s *AdminRoleService) getRspAdminRole(adminRole *model.AdminRole) *rsp.AdminRole {
	//获取角色的菜单ids
	return &rsp.AdminRole{
		Id:          adminRole.Id,
		Name:        adminRole.Name,
		Menus:       s.GetAdminMenuService().GetRoleMenuIds(adminRole.Id),
		Permissions: s.GetAdminPermissionsService().GetRolePermissionsIds(adminRole.Id),
		MenuButtons: s.GetAdminMenuService().GetRoleMenuButtons(adminRole.Id),
		CreateTime:  goutils.TimeFormat("Y-m-d H:i:s", adminRole.CreateTime),
		UpdateTime:  goutils.TimeFormat("Y-m-d H:i:s", adminRole.UpdateTime),
	}
}

// Create 创建
func (s *AdminRoleService) Create() error {
	//绑定参数并校验
	reqBody := new(req.CreateAdminRole)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//查询名称是否已经存在
	name := reqBody.Name
	nameAdminRole, _ := repository.AdminRoleRepository.Take("name = ?", name)
	if nameAdminRole != nil {
		return s.Error(appError.NewError("角色名称已经存在"))
	}

	//新增
	adminRole := &model.AdminRole{
		Name:       reqBody.Name,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	errC := repository.AdminRoleRepository.Create(adminRole)
	if errC != nil {
		return s.Error(appError.NewError(errC.Error()))
	}

	return s.Success(nil)
}

// Edit 编辑
func (s *AdminRoleService) Edit() error {
	//绑定参数并校验
	reqBody := new(req.UpdateAdminRole)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	adminRole, _ := repository.AdminRoleRepository.Take("id = ?", reqBody.Id)
	if adminRole == nil {
		return s.Error(appError.NewError("角色不存在"))
	}

	//查询名称是否已经存在,非自己
	nameAdminRole, _ := repository.AdminRoleRepository.Take("name = ? and id <> ?", reqBody.Name, reqBody.Id)
	if nameAdminRole != nil {
		return s.Error(appError.NewError("角色名称已经存在"))
	}

	//更新信息
	err := repository.AdminRoleRepository.Updates(map[string]interface{}{
		"name":        reqBody.Name,
		"update_time": goutils.GetCurTimeStr(),
	}, "id = ?", reqBody.Id)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// Delete 删除
func (s *AdminRoleService) Delete() error {
	//绑定参数并校验
	reqBody := new(req.DeleteAdminRole)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}
	adminRole, _ := repository.AdminRoleRepository.Take("id = ?", reqBody.Id)
	if adminRole == nil {
		return s.Error(appError.NewError("角色不存在"))
	}

	if reqBody.Id == 1 {
		return s.Error(appError.NewError("超级管理员无法删除"))
	}

	userAdminRole, _ := repository.AdminUserRoleRepository.Take("role_id = ?", adminRole.Id)
	if userAdminRole != nil {
		return s.Error(appError.NewError("角色还存在用户中，请先删除用户中的角色"))
	}

	menuAdminRole, _ := repository.AdminMenuRepository.Take("role_id = ?", adminRole.Id)
	if menuAdminRole != nil {
		return s.Error(appError.NewError("角色还存在菜单中，请先删除菜单中的角色"))
	}

	//删除角色
	err := repository.AdminRoleRepository.Delete(&model.AdminRole{}, "id = ?", adminRole.Id)
	if err != nil {
		return s.Error(appError.NewError("删除失败"))
	}

	return s.Success(nil)
}

// SetPermissions 设置菜单权限和接口权限
func (s *AdminRoleService) SetPermissions() error {
	//绑定参数并校验
	reqBody := new(req.SetAdminPermissions)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}
	adminRole, _ := repository.AdminRoleRepository.Take("id = ?", reqBody.Id)
	if adminRole == nil {
		return s.Error(appError.NewError("角色不存在"))
	}

	menus := reqBody.Menus
	permissions := reqBody.Permissions
	menuButtons := reqBody.MenuButtons
	if len(menus) == 0 && len(permissions) == 0 {
		return s.Error(appError.NewError("角色菜单和角色权限不能为空"))
	}

	//在事务中进行
	txErr := database.Db.Transaction(func(tx *gorm.DB) error {

		adminRoleMenuRepository := repository.NewRepository[*model.AdminRoleMenu](tx, goRedis.Redis)
		adminRolePermissionsRepository := repository.NewRepository[*model.AdminRolePermissions](tx, goRedis.Redis)
		//等于空，删除角色的所有菜单权限
		if len(menus) == 0 {
			err := adminRoleMenuRepository.Delete(&model.AdminRoleMenu{}, "role_id = ?", reqBody.Id)
			if err != nil {
				return err
			}
		} else {
			//新的菜单与原来的菜单比较
			preRoleMenus, _ := adminRoleMenuRepository.Find(database.SqlCdn().Eq("role_id", reqBody.Id))
			for _, preRoleMenu := range preRoleMenus {
				//原来的不在新增的里面
				if !goutils.InSlice(preRoleMenu.MenuId, menus) {
					err := adminRoleMenuRepository.Delete(&model.AdminRoleMenu{}, "id = ?", preRoleMenu.Id)
					if err != nil {
						return err
					}
				}
			}

			//获取原来的菜单id
			preRoleMenuIds := make([]int, 0)
			for _, preRoleMenu := range preRoleMenus {
				preRoleMenuIds = append(preRoleMenuIds, preRoleMenu.MenuId)
			}

			for _, menuId := range menus {
				//不在原来的角色菜单中，进行添加
				if !goutils.InSlice(menuId, preRoleMenuIds) {
					err := adminRoleMenuRepository.Create(&model.AdminRoleMenu{
						RoleId:     reqBody.Id,
						MenuId:     menuId,
						CreateTime: time.Now(),
						UpdateTime: time.Now(),
					})
					if err != nil {
						return err
					}
				}
			}

			for _, menuButton := range menuButtons {
				//更新角色菜单的菜单按钮
				jsonStr, _ := goutils.JsonEncode(menuButton.Buttons)
				err := adminRoleMenuRepository.Updates(map[string]interface{}{
					"btn": jsonStr,
				}, "role_id = ? and menu_id = ?", reqBody.Id, menuButton.MenuId)
				if err != nil {
					return err
				}
			}
		}

		if len(permissions) == 0 {
			err := adminRolePermissionsRepository.Delete(&model.AdminRolePermissions{}, "role_id = ?", reqBody.Id)
			if err != nil {
				return err
			}
		} else {
			//新的权限与原来的权限比较
			preRolePermissions, _ := adminRolePermissionsRepository.Find(database.SqlCdn().Eq("role_id", reqBody.Id))
			for _, preRolePermission := range preRolePermissions {
				//原来的不在新增的里面
				if !goutils.InSlice(preRolePermission.PermissionsId, permissions) {
					err := adminRolePermissionsRepository.Delete(&model.AdminRolePermissions{}, "id = ?", preRolePermission.Id)
					if err != nil {
						return err
					}
				}
			}

			//获取原来的权限id
			preRolePermissionsIds := make([]int, 0)
			for _, preRolePermission := range preRolePermissions {
				preRolePermissionsIds = append(preRolePermissionsIds, preRolePermission.PermissionsId)
			}

			for _, permissionsId := range permissions {
				//不在原来的角色权限中，进行添加
				if !goutils.InSlice(permissionsId, preRolePermissionsIds) {
					err := adminRolePermissionsRepository.Create(&model.AdminRolePermissions{
						RoleId:        reqBody.Id,
						PermissionsId: permissionsId,
						CreateTime:    time.Now(),
						UpdateTime:    time.Now(),
					})
					if err != nil {
						return err
					}
				}
			}
		}

		return nil
	})

	if txErr != nil {
		return s.Error(appError.NewError("设置失败，请重试"))
	}

	return s.Success(nil)
}
