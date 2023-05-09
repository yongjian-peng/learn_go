package router

import (
	"github.com/gofiber/fiber/v2"
	"share/api-server/middleware"
	"share/common/service/admin"
)

func Admin(app *fiber.App) {

	api := app.Group("/admin")
	{

		// 公开路由
		publicApi := api.Group("")
		{
			publicApi.Post("/account/login", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).Login()
			})
		}

		//需要登录的路由
		loginApi := api.Group("", middleware.AdminAuth())
		{
			loginApi.Post("/account/logout", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).Logout()
			})

			loginApi.Get("/account/status_options", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).StatusOptions()
			})

			loginApi.Get("/account/options", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).Options()
			})

			loginApi.Get("/permissions/tree", func(ctx *fiber.Ctx) error {
				return admin.AdminPermissions(ctx).List()
			})

			loginApi.Get("/menu/tree", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).List()
			})

			loginApi.Get("/menu/button_tree", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).ButtonList()
			})

			//登录后的权限菜单
			loginApi.Get("/permissions/menu", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).PermissionsMenu()
			})

			//登录后的权限菜单按钮
			loginApi.Get("/permissions/buttons", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).Buttons()
			})

			//上传图片
			loginApi.Post("/file/upload_img", func(ctx *fiber.Ctx) error {
				return admin.Upload(ctx).UploadImg()
			})

			loginApi.Get("/role/options", func(ctx *fiber.Ctx) error {
				return admin.AdminRole(ctx).Options()
			})
		}

		//需要登录，并授权的路由
		//权限管理
		permissionsApi := api.Group("permissions", middleware.PermissionsAuth())
		{
			permissionsApi.Get("/list", func(ctx *fiber.Ctx) error {
				return admin.AdminPermissions(ctx).List()
			})
			permissionsApi.Post("/create", func(ctx *fiber.Ctx) error {
				return admin.AdminPermissions(ctx).Create()
			})
			permissionsApi.Post("/update", func(ctx *fiber.Ctx) error {
				return admin.AdminPermissions(ctx).Update()
			})
			permissionsApi.Post("/delete", func(ctx *fiber.Ctx) error {
				return admin.AdminPermissions(ctx).Delete()
			})
		}

		//权限管理
		menuApi := api.Group("menu", middleware.PermissionsAuth())
		{
			menuApi.Get("/list", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).List()
			})
			menuApi.Post("/create", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).Create()
			})
			menuApi.Post("/edit", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).Edit()
			})
			menuApi.Post("/delete", func(ctx *fiber.Ctx) error {
				return admin.AdminMenu(ctx).Delete()
			})
		}

		//角色管理
		roleApi := api.Group("role", middleware.PermissionsAuth())
		{
			roleApi.Get("/list", func(ctx *fiber.Ctx) error {
				return admin.AdminRole(ctx).List()
			})
			roleApi.Post("/create", func(ctx *fiber.Ctx) error {
				return admin.AdminRole(ctx).Create()
			})
			roleApi.Post("/edit", func(ctx *fiber.Ctx) error {
				return admin.AdminRole(ctx).Edit()
			})
			roleApi.Post("/delete", func(ctx *fiber.Ctx) error {
				return admin.AdminRole(ctx).Delete()
			})
			roleApi.Post("/set_permissions", func(ctx *fiber.Ctx) error {
				return admin.AdminRole(ctx).SetPermissions()
			})
		}

		accountApi := api.Group("account", middleware.PermissionsAuth())
		{

			accountApi.Get("/list", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).List()
			})

			accountApi.Post("/create", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).Create()
			})

			accountApi.Post("/edit", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).Edit()
			})

			accountApi.Post("/reset_pwd", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).ResetPwd()
			})

			accountApi.Post("/change_status", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).ChangeStatus()
			})

			accountApi.Post("/set_role", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).SetRole()
			})

			accountApi.Post("/delete", func(ctx *fiber.Ctx) error {
				return admin.AdminUser(ctx).Delete()
			})
		}

		appApi := api.Group("app", middleware.PermissionsAuth())
		{
			appApi.Get("/list", func(ctx *fiber.Ctx) error {
				return admin.App(ctx).List()
			})

			appApi.Post("/create", func(ctx *fiber.Ctx) error {
				return admin.App(ctx).Create()
			})

			appApi.Post("/edit", func(ctx *fiber.Ctx) error {
				return admin.App(ctx).Edit()
			})

			appApi.Post("/change_status", func(ctx *fiber.Ctx) error {
				return admin.App(ctx).ChangeStatus()
			})

			appApi.Post("/reset_secret", func(ctx *fiber.Ctx) error {
				return admin.App(ctx).ResetSecret()
			})

			appApi.Post("/change_config", func(ctx *fiber.Ctx) error {
				return admin.App(ctx).ChangeConfig()
			})
		}
	}

}
