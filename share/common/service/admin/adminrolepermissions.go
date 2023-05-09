package admin

import (
	"github.com/gofiber/fiber/v2"
	"share/common/repository"
	"sync"
)

var adminRolePermissionsService *AdminRolePermissionsService
var AdminRolePermissionsServiceOnce sync.Once

type AdminRolePermissionsService struct {
	*Service
}

func AdminRolePermissions(c *fiber.Ctx) *AdminRolePermissionsService {
	return &AdminRolePermissionsService{Service: NewService(c, "adminRolePermissionsService")}
}

// CheckPermissionsExist 检查权限是否存在
func (s *AdminRolePermissionsService) CheckPermissionsExist(id int) bool {
	adminRolePermissions, _ := repository.AdminRolePermissionsRepository.Take("permissions_id = ?", id)
	if adminRolePermissions != nil {
		return true
	}
	return false
}
