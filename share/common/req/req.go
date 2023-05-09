package req

type Header struct {
	Version   string `label:"Version" json:"Version" validate:"oneof=1.0"`
	AppId     string `label:"AppId" json:"AppId" validate:"required" comment:"应用ID"`
	Signature string `label:"Signature" json:"Signature" validate:"required" comment:"签名"`
}

type AdminPermissionsList struct {
	PageNum  int `label:"pageNum" query:"pageNum" validate:"numeric,gte=1"`
	PageSize int `label:"pageSize" query:"pageSize" validate:"numeric,gte=1"`
}

type CreateAdminPermissions struct {
	Pid      int    `label:"pid" json:"pid" validate:"numeric,gte=0"`
	Name     string `label:"name" json:"name" validate:"required"`
	RouteUrl string `label:"route_url" json:"route_url" validate:"required"`
	Sort     int    `label:"sort" json:"sort" validate:"numeric,gte=0"`
}

type UpdateAdminPermissions struct {
	Id       int    `label:"id" json:"id" validate:"numeric,gte=1"`
	Pid      int    `label:"pid" json:"pid" validate:"numeric,gte=0"`
	Name     string `label:"name" json:"name" validate:"required"`
	RouteUrl string `label:"route_url" json:"route_url" validate:"required"`
	Sort     int    `label:"sort" json:"sort" validate:"numeric,gte=0"`
}

type DeleteAdminPermissions struct {
	Id int `label:"id" json:"id" validate:"numeric,gte=1"`
}

type AdminMenuAuthBtn struct {
	Name string `json:"name"`
	Desc string `json:"desc"`
}

type CreateAdminMenu struct {
	Pid         int                `json:"pid" label:"pid" validate:"numeric,gte=0"`
	Path        string             `json:"path" label:"path" validate:"required"`
	Name        string             `json:"name" label:"name" validate:"required"`
	Redirect    string             `json:"redirect" label:"redirect"`
	Component   string             `json:"component" label:"component"`
	Icon        string             `json:"icon" label:"icon" validate:"required"`
	Title       string             `json:"title" label:"title" validate:"required"`
	ActiveMenu  string             `json:"activeMenu" label:"activeMenu"`
	Link        string             `json:"string" label:"link"`
	IsHide      int                `json:"isHide" label:"isHide" validate:"numeric,oneof=0 1"`
	IsFull      int                `json:"isFull" label:"isFull" validate:"numeric,oneof=0 1"`
	IsAffix     int                `json:"isAffix" label:"isAffix" validate:"numeric,oneof=0 1"`
	IsKeepAlive int                `json:"isKeepAlive" label:"isKeepAlive" validate:"numeric,oneof=0 1"`
	AuthBtn     []AdminMenuAuthBtn `json:"authBtn" label:"authBtn"`
	Sort        int                `json:"sort" label:"sort" validate:"numeric,gte=0"`
}

type UpdateAdminMenu struct {
	Id          int                `json:"id" label:"id" validate:"numeric,gte=0"`
	Pid         int                `json:"pid" label:"pid" validate:"numeric,gte=0"`
	Path        string             `json:"path" label:"path" validate:"required"`
	Name        string             `json:"name" label:"name" validate:"required"`
	Redirect    string             `json:"redirect" label:"redirect"`
	Component   string             `json:"component" label:"component"`
	Icon        string             `json:"icon" label:"icon" validate:"required"`
	Title       string             `json:"title" label:"title" validate:"required"`
	ActiveMenu  string             `json:"activeMenu" label:"activeMenu"`
	Link        string             `json:"link" label:"link"`
	IsHide      int                `json:"isHide" label:"isHide" validate:"numeric,oneof=0 1"`
	IsFull      int                `json:"isFull" label:"isFull" validate:"numeric,oneof=0 1"`
	IsAffix     int                `json:"isAffix" label:"isAffix" validate:"numeric,oneof=0 1"`
	IsKeepAlive int                `json:"isKeepAlive" label:"isKeepAlive" validate:"numeric,oneof=0 1"`
	AuthBtn     []AdminMenuAuthBtn `json:"authBtn" label:"authBtn"`
	Sort        int                `json:"sort" label:"sort" validate:"numeric,gte=0"`
}

type DeleteAdminMenu struct {
	Id int `label:"id" json:"id" validate:"numeric,gte=1"`
}

type PageList struct {
	PageNum  int `label:"pageNum" query:"pageNum" validate:"numeric,gte=1"`
	PageSize int `label:"pageSize" query:"pageSize" validate:"numeric,gte=1"`
}

type AdminMenuList struct {
	PageList
}

type AdminMenuButtonList struct {
	PageList
}

type AdminRoleList struct {
	PageList
	Name       string   `label:"name" query:"name"`
	CreateTime []string `label:"createTime" query:"createTime"`
}

type CreateAdminRole struct {
	Name string `json:"name" label:"name" validate:"required"`
}

type UpdateAdminRole struct {
	Id   int    `json:"id" label:"id" validate:"numeric,gte=1"`
	Name string `json:"name" label:"name" validate:"required"`
}

type DeleteAdminRole struct {
	Id int `label:"id" json:"id" validate:"numeric,gte=1"`
}

type AdminRoleMenuButton struct {
	MenuId  int      `json:"menuId"`
	Buttons []string `json:"buttons"`
}

type SetAdminPermissions struct {
	Id          int                   `json:"id" label:"id" validate:"numeric,gte=1"`
	Menus       []int                 `json:"menus" label:"menus" `
	MenuButtons []AdminRoleMenuButton `json:"menuButtons" label:"menuButtons" `
	Permissions []int                 `json:"permissions" label:"permissions"`
}

type AdminUserList struct {
	PageNum    int      `label:"pageNum" query:"pageNum" validate:"numeric,gte=1"`
	PageSize   int      `label:"pageSize" query:"pageSize" validate:"numeric,gte=1"`
	UserName   string   `label:"userName" query:"userName"`
	Status     int      `label:"status" query:"status"`
	CreateTime []string `label:"createTime" query:"createTime"`
}

type AdminAppList struct {
	PageList
	Name       string   `label:"name" query:"name"`
	Status     int      `label:"status" query:"status"`
	CreateTime []string `label:"createTime" query:"createTime"`
}

type CreateAdminUser struct {
	Username string `json:"username" label:"username" validate:"required"`
	Password string `json:"password" label:"password" validate:"required"`
	Avatar   string `json:"avatar"  label:"avatar" `
	Status   int    `json:"status" label:"status" validate:"numeric,oneof=0 1"`
}

type CreateAdminApp struct {
	Name        string `json:"name" label:"username" validate:"required"`
	PackageName string `json:"packageName" label:"packageName" validate:"required"`
	Status      int    `json:"status" label:"status" validate:"numeric,oneof=0 1"`
}

type UpdateAdminUser struct {
	Id       int    `json:"id" label:"id" validate:"numeric,gte=1"`
	Username string `json:"username" label:"username" validate:"required"`
	Password string `json:"password" label:"password"`
	Avatar   string `json:"avatar"  label:"avatar" `
	Status   int    `json:"status" label:"status" validate:"numeric,oneof=0 1"`
	Roles    []int  `json:"roles" label:"roles" validate:"required"`
}

type UpdateAdminApp struct {
	Id          int    `json:"id" label:"id" validate:"numeric,gte=1"`
	Name        string `json:"name" label:"username" validate:"required"`
	PackageName string `json:"packageName" label:"packageName" validate:"required"`
	Status      int    `json:"status" label:"status" validate:"numeric,oneof=0 1"`
}

type ResetAdminUserPwd struct {
	Id int `json:"id" label:"id" validate:"numeric,gte=1"`
}

type ResetAdminAppSecret struct {
	Id int `json:"id" label:"id" validate:"numeric,gte=1"`
}

type SetAdminUserStatus struct {
	Id     int `json:"id" label:"id" validate:"numeric,gte=1"`
	Status int `json:"status" label:"status" validate:"numeric,oneof=0 1"`
}

type SetAdminAppStatus struct {
	Id     int `json:"id" label:"id" validate:"numeric,gte=1"`
	Status int `json:"status" label:"status" validate:"numeric,oneof=0 1"`
}

type SetAdminAppConfig struct {
	AppId              int `json:"appId" label:"appId" validate:"numeric,gte=1"`
	UserDayPayoutLimit int `json:"userDayPayoutLimit" label:"userDayPayoutLimit" validate:"numeric"`
	PayLimit           int `json:"payLimit" label:"payLimit" validate:"numeric"`
}

type SetAdminUserRole struct {
	Id    int   `json:"id" label:"id" validate:"numeric,gte=1"`
	Roles []int `json:"roles" label:"roles" validate:"required"`
}

type DeleteAdminUser struct {
	Id int `label:"id" json:"id" validate:"numeric,gte=1"`
}

type AdminUserLogin struct {
	Username string `json:"username" label:"username" validate:"required"`
	Password string `json:"password" label:"password" validate:"required"`
}

type UserRegister struct {
	AppUid       int64  `json:"appUid" label:"appUid" validate:"required"`
	Phone        string `json:"phone" label:"phone" validate:"required"`
	InviteCode   string `json:"inviteCode" label:"inviteCode" validate:"required"`
	DeviceId     string `json:"deviceId" label:"deviceId" validate:"required"`
	RegisterType string `json:"registerType" label:"registerType" validate:"required"`
	RegisterTime int    `json:"registerTime" label:"registerTime" validate:"required"`
}

type UserLogin struct {
	AppUid int64 `json:"appUid" label:"appUid" validate:"required"`
}
