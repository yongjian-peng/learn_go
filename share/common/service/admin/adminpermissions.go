package admin

import (
	"github.com/gofiber/fiber/v2"
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

var adminPermissionsService *AdminPermissionsService
var AdminPermissionsServiceOnce sync.Once

type AdminPermissionsService struct {
	*Service
}

func AdminPermissions(c *fiber.Ctx) *AdminPermissionsService {
	return &AdminPermissionsService{Service: NewService(c, "adminPermissionsService")}
}

// List 列表
func (s *AdminPermissionsService) List() error {
	//绑定参数并校验
	reqBody := new(req.AdminPermissionsList)
	if err := s.ValidateGet(reqBody); err != nil {
		return s.Error(err)
	}
	//获取顶级权限列表
	adminPermissionsList, pagination, _ := repository.AdminPermissionsRepository.FindPageByCdn(database.SqlCdn().Eq("pid", 0).Page(reqBody.PageNum, reqBody.PageSize))
	rspPageResult := &rsp.PageResult{}
	rspPageResult.Page = pagination
	rspAdminPermissionsList := make([]*rsp.AdminPermissions, 0)
	for _, adminPermissions := range adminPermissionsList {
		//获取下级权限
		children := s.getChildrenRspPermissions(adminPermissions.Id)
		rspAdminPermissionsList = append(rspAdminPermissionsList, s.getRspAdminPermissions(adminPermissions, children))
	}
	rspPageResult.List = rspAdminPermissionsList
	return s.Success(rspPageResult)
}

// GetRolePermissionsIds 获取角色的权限ids
func (s *AdminPermissionsService) GetRolePermissionsIds(roleId int) []int {
	permissionsIds := make([]int, 0)
	rolePermissionsList, _ := repository.AdminRolePermissionsRepository.Find(database.SqlCdn().Eq("role_id", roleId))
	for _, rolePermissions := range rolePermissionsList {
		permissionsIds = append(permissionsIds, rolePermissions.PermissionsId)
	}
	return permissionsIds
}

// getRspAdminPermissions 组装RspAdminPermissions
func (s *AdminPermissionsService) getRspAdminPermissions(adminPermissions *model.AdminPermissions, children []*rsp.AdminPermissions) *rsp.AdminPermissions {
	return &rsp.AdminPermissions{
		Id:         adminPermissions.Id,
		Pid:        adminPermissions.Pid,
		Name:       adminPermissions.Name,
		RouteUrl:   adminPermissions.RouteUrl,
		Sort:       adminPermissions.Sort,
		CreateTime: goutils.TimeFormat("Y-m-d H:i:s", adminPermissions.CreateTime),
		UpdateTime: goutils.TimeFormat("Y-m-d H:i:s", adminPermissions.UpdateTime),
		Children:   children,
	}
}

// getChildrenPermissions 获取下级权限列表
func (s *AdminPermissionsService) getChildrenRspPermissions(pid int) []*rsp.AdminPermissions {
	rspAdminPermissionsList := make([]*rsp.AdminPermissions, 0)
	adminPermissionsList := s.getChildrenPermissions(pid)
	for _, adminPermissions := range adminPermissionsList {
		children := s.getChildrenRspPermissions(adminPermissions.Id)
		rspAdminPermissionsList = append(rspAdminPermissionsList, s.getRspAdminPermissions(adminPermissions, children))
	}
	return rspAdminPermissionsList
}

// getChildrenPermissions 获取下级权限列表
func (s *AdminPermissionsService) getChildrenPermissions(pid int) []*model.AdminPermissions {
	AdminPermissionsList, _ := repository.AdminPermissionsRepository.Find(database.SqlCdn().Eq("pid", pid))
	return AdminPermissionsList
}

// Create 创建
func (s *AdminPermissionsService) Create() error {
	//绑定参数并校验
	reqBody := new(req.CreateAdminPermissions)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	pid := reqBody.Pid

	//0为顶级id
	if pid != 0 {
		parentAdminPermission, _ := repository.AdminPermissionsRepository.Take("id = ?", pid)
		if parentAdminPermission == nil {
			return s.Error(appError.NewError("上级权限不存在"))
		}
	}

	//查询权限名称是否已经存在
	name := reqBody.Name
	nameAdminPermission, _ := repository.AdminPermissionsRepository.Take("name = ?", name)
	if nameAdminPermission != nil {
		return s.Error(appError.NewError("权限名称已经存在"))
	}

	//查询权限路由是否已经存在
	routeUrl := reqBody.RouteUrl
	routeUrlAdminPermission, _ := repository.AdminPermissionsRepository.Take("route_url", routeUrl)
	if routeUrlAdminPermission != nil {
		return s.Error(appError.NewError("权限路由已经存在"))
	}

	//新增
	adminPermissions := &model.AdminPermissions{
		Pid:        reqBody.Pid,
		Name:       reqBody.Name,
		RouteUrl:   reqBody.RouteUrl,
		Sort:       reqBody.Sort,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	err := repository.AdminPermissionsRepository.Create(adminPermissions)
	if err != nil {
		return s.Error(appError.NewError(err.Error()))
	}

	return s.Success(nil)
}

// Update 更新
func (s *AdminPermissionsService) Update() error {
	//绑定参数并校验
	reqBody := new(req.UpdateAdminPermissions)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}

	//0为顶级id
	if reqBody.Pid != 0 {
		parentAdminPermission, _ := repository.AdminPermissionsRepository.Take("id = ?", reqBody.Pid)
		if parentAdminPermission == nil {
			return s.Error(appError.NewError("上级权限不存在"))
		}
	}

	adminPermission, _ := repository.AdminPermissionsRepository.Take("id = ?", reqBody.Id)
	if adminPermission == nil {
		return s.Error(appError.NewError("权限不存在"))
	}

	//查询权限名称是否已经存在,非自己
	nameAdminPermission, _ := repository.AdminPermissionsRepository.Take("name = ? and id <> ?", reqBody.Name, reqBody.Id)
	if nameAdminPermission != nil {
		return s.Error(appError.NewError("权限名称已经存在"))
	}

	//查询权限路由是否已经存在,非自己
	routeUrlAdminPermission, _ := repository.AdminPermissionsRepository.Take("route_url = ? and id <> ?", reqBody.RouteUrl, reqBody.Id)
	if routeUrlAdminPermission != nil {
		return s.Error(appError.NewError("权限路由已经存在"))
	}

	//更新信息
	err := repository.AdminPermissionsRepository.Updates(map[string]interface{}{
		"pid":         reqBody.Pid,
		"name":        reqBody.Name,
		"route_url":   reqBody.RouteUrl,
		"sort":        reqBody.Sort,
		"update_time": goutils.GetCurTimeStr(),
	}, "id = ?", reqBody.Id)
	if err != nil {
		return err
	}

	return s.Success(nil)
}

// Delete 删除
func (s *AdminPermissionsService) Delete() error {
	//绑定参数并校验
	reqBody := new(req.DeleteAdminPermissions)
	if err := s.ValidatePost(reqBody); err != nil {
		return s.Error(err)
	}
	adminPermission, _ := repository.AdminPermissionsRepository.Take("id = ?", reqBody.Id)
	if adminPermission == nil {
		return s.Error(appError.NewError("权限不存在"))
	}

	//权限是否已被使用，如果被使用不能删除
	if s.GetAdminRolePermissionsService().CheckPermissionsExist(adminPermission.Id) {
		return s.Error(appError.NewError("权限还存在于角色中，请先删除角色中的权限"))
	}

	//删除权限
	err := repository.AdminPermissionsRepository.Delete(&model.AdminPermissions{}, "id = ?", adminPermission.Id)
	if err != nil {
		return s.Error(appError.NewError("删除失败"))
	}

	return s.Success(nil)
}
