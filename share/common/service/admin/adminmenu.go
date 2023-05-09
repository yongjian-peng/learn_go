package admin

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gookit/goutil/dump"
	"github.com/samber/lo"
	"share/common/model"
	"share/common/pkg/appError"
	"share/common/pkg/database"
	"share/common/pkg/goutils"
	"share/common/repository"
	"share/common/req"
	"share/common/rsp"
	"sync"
	"time"
)

var adminMenuService *AdminMenuService
var AdminMenuServiceOnce sync.Once

type AdminMenuService struct {
	*Service
}

func AdminMenu(c *fiber.Ctx) *AdminMenuService {
	return &AdminMenuService{Service: NewService(c, "adminMenuService")}
}

func (s *AdminMenuService) Buttons() error {
	adminId := s.AdminId
	buttons := s.GetAdminUserService().GetUserMenuButtons(adminId)
	return s.Success(buttons)
}

// PermissionsMenu 获取权限菜单列表
func (s *AdminMenuService) PermissionsMenu() error {
	//获取顶级菜单列表
	adminMenuList, _ := repository.AdminMenuRepository.Find(database.SqlCdn().Eq("pid", 0))
	rspAdminMenuList := make([]*rsp.AdminMenu, 0)

	menuIds := make([]int, 0)
	if s.AdminId != 1 {
		//非超级管理员获取用户拥有的菜单ids
		menuIds = s.GetAdminUserService().GetUserMenuIds(s.AdminId)
	}

	//遍历所有菜单判断权限叠加菜单
	for _, adminMenu := range adminMenuList {
		//非超级管理员判断菜单权限
		if s.AdminId != 1 && !lo.Contains(menuIds, adminMenu.Id) {
			continue
		}
		//获取下级权限
		children := s.getChildrenRspMenu(adminMenu.Id, menuIds)
		rspAdminMenuList = append(rspAdminMenuList, s.getRspAdminMenu(adminMenu, children))
	}

	return s.Success(rspAdminMenuList)
}

// List 列表
func (s *AdminMenuService) List() error {
	//绑定参数并校验
	reqBody := new(req.AdminMenuList)
	if err := s.ValidateGet(reqBody); err != nil {
		return s.Error(err)
	}
	//获取顶级列表
	adminMenuList, pagination, _ := repository.AdminMenuRepository.FindPageByCdn(database.SqlCdn().Eq("pid", 0).Page(reqBody.PageNum, reqBody.PageSize))
	rspPageResult := &rsp.PageResult{}
	rspPageResult.Page = pagination
	rspAdminMenuList := make([]*rsp.AdminMenu, 0)
	for _, adminMenu := range adminMenuList {
		//获取下级权限
		children := s.getChildrenRspMenu(adminMenu.Id, nil)
		rspAdminMenuList = append(rspAdminMenuList, s.getRspAdminMenu(adminMenu, children))
	}
	rspPageResult.List = rspAdminMenuList
	return s.Success(rspPageResult)
}

// ButtonList 列表
func (s *AdminMenuService) ButtonList() error {
	//绑定参数并校验
	reqBody := new(req.AdminMenuButtonList)
	if err := s.ValidateGet(reqBody); err != nil {
		return s.Error(err)
	}
	//获取顶级列表
	adminMenuList, pagination, _ := repository.AdminMenuRepository.FindPageByCdn(database.SqlCdn().Eq("pid", 0).Page(reqBody.PageNum, reqBody.PageSize).Desc("sort"))
	rspPageResult := &rsp.PageResult{}
	rspPageResult.Page = pagination
	rspAdminMenuList := make([]*rsp.AdminMenu, 0)
	for _, adminMenu := range adminMenuList {
		//获取下级权限
		children := s.getChildrenRspMenu(adminMenu.Id, nil)
		rspAdminMenuList = append(rspAdminMenuList, s.getRspAdminMenu(adminMenu, children))
	}
	rspPageResult.List = rspAdminMenuList
	return s.Success(rspPageResult)
}

// getRspAdminMenu 组装RspAdminMenu
func (s *AdminMenuService) getRspAdminMenu(adminMenu *model.AdminMenu, children []*rsp.AdminMenu) *rsp.AdminMenu {
	AdminMenuButtons := make([]*rsp.AdminMenuButton, 0)
	_ = goutils.JsonDecode(adminMenu.AuthBtn, &AdminMenuButtons)
	return &rsp.AdminMenu{
		Id:        adminMenu.Id,
		Pid:       adminMenu.Pid,
		Path:      adminMenu.Path,
		Name:      adminMenu.Name,
		Redirect:  adminMenu.Redirect,
		Component: adminMenu.Component,
		Meta: rsp.AdminMenuMeta{
			Icon:        adminMenu.Icon,
			Title:       adminMenu.Title,
			ActiveMenu:  adminMenu.ActiveMenu,
			Link:        adminMenu.Link,
			IsHide:      goutils.IfBool(adminMenu.IsHide > 0, true, false),
			IsFull:      goutils.IfBool(adminMenu.IsFull > 0, true, false),
			IsAffix:     goutils.IfBool(adminMenu.IsAffix > 0, true, false),
			IsKeepAlive: goutils.IfBool(adminMenu.IsKeepAlive > 0, true, false),
		},
		CreateTime: goutils.TimeFormat("Y-m-d H:i:s", adminMenu.CreateTime),
		UpdateTime: goutils.TimeFormat("Y-m-d H:i:s", adminMenu.UpdateTime),
		Children:   children,
		MenuButton: AdminMenuButtons,
	}
}

// getChildrenMenu 获取下级权限列表
func (s *AdminMenuService) getChildrenRspMenu(pid int, menuIds []int) []*rsp.AdminMenu {
	rspAdminMenuList := make([]*rsp.AdminMenu, 0)
	adminMenuList := s.getChildrenMenu(pid)
	for _, adminMenu := range adminMenuList {
		if len(menuIds) > 0 && !lo.Contains(menuIds, adminMenu.Id) {
			continue
		}
		children := s.getChildrenRspMenu(adminMenu.Id, menuIds)
		rspAdminMenuList = append(rspAdminMenuList, s.getRspAdminMenu(adminMenu, children))
	}
	return rspAdminMenuList
}

// getChildrenMenuButton 获取下级权限列表
func (s *AdminMenuService) getChildrenMenuButton(pid int, menuIds []int) []*rsp.AdminMenu {
	rspAdminMenuList := make([]*rsp.AdminMenu, 0)
	adminMenuList := s.getChildrenMenu(pid)
	for _, adminMenu := range adminMenuList {
		if len(menuIds) > 0 && !lo.Contains(menuIds, adminMenu.Id) {
			continue
		}
		children := s.getChildrenRspMenu(adminMenu.Id, menuIds)
		rspAdminMenuList = append(rspAdminMenuList, s.getRspAdminMenu(adminMenu, children))
	}
	return rspAdminMenuList
}

// GetRoleMenuIds 获取角色的菜单ids
func (s *AdminMenuService) GetRoleMenuIds(roleId int) []int {
	menuIds := make([]int, 0)
	roleMenuList, _ := repository.AdminRoleMenuRepository.Find(database.SqlCdn().Eq("role_id", roleId))
	for _, roleMenu := range roleMenuList {
		menuIds = append(menuIds, roleMenu.MenuId)
	}
	return menuIds
}

// GetRoleMenuButtons 获取角色的菜单ids
func (s *AdminMenuService) GetRoleMenuButtons(roleId int) []*rsp.AdminRoleMenuButton {
	adminRoleMenuButtons := make([]*rsp.AdminRoleMenuButton, 0)
	roleMenuList, _ := repository.AdminRoleMenuRepository.Find(database.SqlCdn().Eq("role_id", roleId))
	for _, roleMenu := range roleMenuList {
		adminMenuButton := make([]string, 0)
		_ = goutils.JsonDecode(roleMenu.Btn, &adminMenuButton)
		menuButtons := make([]string, 0)
		for _, menuButton := range adminMenuButton {
			menuButtons = append(menuButtons, menuButton)
		}
		adminRoleMenuButtons = append(adminRoleMenuButtons, &rsp.AdminRoleMenuButton{
			MenuId:  roleMenu.MenuId,
			Buttons: menuButtons,
		})
	}
	return adminRoleMenuButtons
}

// getChildrenMenu 获取下级权限列表
func (s *AdminMenuService) getChildrenMenu(pid int) []*model.AdminMenu {
	AdminMenuList, _ := repository.AdminMenuRepository.Find(database.SqlCdn().Eq("pid", pid))
	return AdminMenuList
}

// Create 创建
func (s *AdminMenuService) Create() error {
	//绑定参数并校验
	reqBody := new(req.CreateAdminMenu)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	dump.P(reqBody)
	pid := reqBody.Pid

	//0为顶级id
	if pid != 0 {
		parentAdminMenu, _ := repository.AdminMenuRepository.Take("id = ?", pid)
		if parentAdminMenu == nil {
			return appError.NewError("上级菜单不存在")
		}
	}

	//查询名称是否已经存在
	name := reqBody.Name
	nameAdminMenu, _ := repository.AdminMenuRepository.Take("name = ?", name)
	if nameAdminMenu != nil {
		return appError.NewError("菜单名称已经存在")
	}

	//新增
	authBtn, _ := goutils.JsonEncode(reqBody.AuthBtn)
	adminMenu := &model.AdminMenu{
		Pid:        reqBody.Pid,
		Name:       reqBody.Name,
		Path:       reqBody.Path,
		Redirect:   reqBody.Redirect,
		Component:  reqBody.Component,
		Icon:       reqBody.Icon,
		Title:      reqBody.Title,
		ActiveMenu: reqBody.ActiveMenu,
		Link:       reqBody.Link,
		IsHide:     reqBody.IsHide,
		IsFull:     reqBody.IsFull,
		IsAffix:    reqBody.IsKeepAlive,
		AuthBtn:    authBtn,
		Sort:       reqBody.Sort,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	errC := repository.AdminMenuRepository.Create(adminMenu)
	if errC != nil {
		return s.Error(appError.NewError(errC.Error()))
	}

	return s.Success(nil)
}

// Edit 更新
func (s *AdminMenuService) Edit() error {
	//绑定参数并校验
	reqBody := new(req.UpdateAdminMenu)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//0为顶级id
	if reqBody.Pid != 0 {
		parentAdminMenu, _ := repository.AdminMenuRepository.Take("id = ?", reqBody.Pid)
		if parentAdminMenu == nil {
			return s.Error(appError.NewError("上级菜单不存在"))
		}
	}

	adminMenu, _ := repository.AdminMenuRepository.Take("id = ?", reqBody.Id)
	if adminMenu == nil {
		return s.Error(appError.NewError("菜单不存在"))
	}

	//查询名称是否已经存在,非自己
	nameAdminMenu, _ := repository.AdminMenuRepository.Take("name = ? and id <> ?", reqBody.Name, reqBody.Id)
	if nameAdminMenu != nil {
		return s.Error(appError.NewError("菜单名称已经存在"))
	}

	authBtn, _ := goutils.JsonEncode(reqBody.AuthBtn)
	dump.P(reqBody)
	dump.P(authBtn)
	//更新信息
	errU := repository.AdminMenuRepository.Updates(map[string]interface{}{
		"pid":           reqBody.Pid,
		"path":          reqBody.Path,
		"name":          reqBody.Name,
		"redirect":      reqBody.Redirect,
		"component":     reqBody.Component,
		"icon":          reqBody.Icon,
		"title":         reqBody.Title,
		"active_menu":   reqBody.ActiveMenu,
		"link":          reqBody.Link,
		"is_hide":       reqBody.IsHide,
		"is_full":       reqBody.IsFull,
		"is_affix":      reqBody.IsAffix,
		"is_keep_alive": reqBody.IsKeepAlive,
		"auth_btn":      authBtn,
		"sort":          reqBody.Sort,
		"update_time":   goutils.GetCurTimeStr(),
	}, "id = ?", reqBody.Id)
	if errU != nil {
		return s.Error(appError.NewError(errU.Error()))
	}

	return s.Success(nil)
}

// Delete 删除
func (s *AdminMenuService) Delete() error {
	//绑定参数并校验
	reqBody := new(req.DeleteAdminMenu)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}
	adminMenu, _ := repository.AdminMenuRepository.Take("id = ?", reqBody.Id)
	if adminMenu == nil {
		return s.Error(appError.NewError("菜单不存在"))
	}

	childrenAdminMenu, _ := repository.AdminMenuRepository.Take("pid = ?", adminMenu.Id)
	if childrenAdminMenu != nil {
		return s.Error(appError.NewError("菜单还存在子菜单，请先删除子菜单"))
	}

	//删除菜单
	errD := repository.AdminMenuRepository.Delete(&model.AdminMenu{}, "id = ?", adminMenu.Id)
	if errD != nil {
		return s.Error(appError.NewError("删除失败"))
	}

	return s.Success(nil)
}
