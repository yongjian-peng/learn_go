package rsp

import (
	"share/common/pkg/database"
)

// PageResult 分页返回数据
type PageResult struct {
	Page *database.Pagination `json:"page"` // 分页信息
	List interface{}          `json:"list"` // 数据
}

// CursorResult Cursor分页返回数据
type CursorResult struct {
	List    interface{} `json:"list"`    // 数据
	Cursor  string      `json:"cursor"`  // 下一页
	HasMore bool        `json:"hasMore"` // 是否还有数据
}

type AdminPermissions struct {
	Id         int                 `json:"id"`
	Pid        int                 `json:"pid"`
	Name       string              `json:"name"`
	RouteUrl   string              ` json:"route_url"`
	Sort       int                 `json:"sort"`
	CreateTime string              `json:"create_time"`
	UpdateTime string              `json:"update_time"`
	Children   []*AdminPermissions `json:"children"`
}

type AdminMenuMeta struct {
	Icon        string `json:"icon"`
	Title       string `json:"title"`
	ActiveMenu  string `json:"activeMenu"`
	Link        string `json:"link"`
	IsHide      bool   `json:"isHide"`
	IsFull      bool   `json:"isFull"`
	IsAffix     bool   `json:"isAffix"`
	IsKeepAlive bool   `json:"isKeepAlive"`
}

type AdminMenu struct {
	Id         int                `json:"id"`
	Pid        int                `json:"pid"`
	Path       string             `json:"path"`
	Name       string             `json:"name"`
	Redirect   string             `json:"redirect"`
	Component  string             `json:"component"`
	Meta       AdminMenuMeta      `json:"meta"`
	CreateTime string             `json:"createTime"`
	UpdateTime string             `json:"updateTime"`
	MenuButton []*AdminMenuButton `json:"menuButton"`
	Children   []*AdminMenu       `json:"children"`
}

type AdminMenuButton struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type AdminRoleMenuButton struct {
	MenuId  int      `json:"menuId"`
	Buttons []string `json:"buttons"`
}

type AdminRole struct {
	Id          int                    `json:"id"`
	Name        string                 `json:"name"`
	Menus       []int                  `json:"menus"`
	MenuButtons []*AdminRoleMenuButton `json:"menuButtons"`
	Permissions []int                  `json:"permissions"`
	CreateTime  string                 `json:"createTime"`
	UpdateTime  string                 `json:"updateTime"`
}

// LoginAdminUser 管理员表
type LoginAdminUser struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

type AdminUser struct {
	Id         int    `json:"id"`
	Username   string `json:"username"`
	Avatar     string `json:"avatar"  label:"avatar" `
	Status     int    `json:"status" label:"status"`
	Roles      []int  `json:"roles" label:"roles"`
	CreateTime string `json:"createTime"`
}

type User struct {
	Id         int64  `json:"id"`
	InviteCode string `json:"inviteCode"  label:"inviteCode" `
	Level      int    `json:"level" label:"level"`
}

type AdminApp struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Secret      string `json:"Secret"`
	PackageName string `json:"package_name"`
	CreateTime  string `json:"createTime"`
	UpdateTime  string `json:"updateTime"`
}

type Option struct {
	Label    string `json:"label"`
	Value    int    `json:"value"`
	Disabled bool   `json:"disabled"`
}
