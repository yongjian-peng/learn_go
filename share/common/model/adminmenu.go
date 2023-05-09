package model

import (
	"time"
)

// AdminMenu 菜单列表
type AdminMenu struct {
	Id          int       `gorm:"column:id;type:int(10);primary_key;AUTO_INCREMENT;comment:菜单id" json:"id"`
	Pid         int       `gorm:"column:pid;type:int(11);default:0;comment:父级菜单;NOT NULL" json:"pid"`
	Path        string    `gorm:"column:path;type:varchar(255);comment:路由地址;NOT NULL" json:"path"`
	Name        string    `gorm:"column:name;type:varchar(255);comment:路由名称;NOT NULL" json:"name"`
	Redirect    string    `gorm:"column:redirect;type:varchar(255);comment:重定向地址;NOT NULL" json:"redirect"`
	Component   string    `gorm:"column:component;type:varchar(255);comment:视图文件路径;NOT NULL" json:"component"`
	Icon        string    `gorm:"column:icon;type:varchar(30);comment:菜单图标;NOT NULL" json:"icon"`
	Title       string    `gorm:"column:title;type:varchar(255);comment:菜单标题;NOT NULL" json:"title"`
	ActiveMenu  string    `gorm:"column:active_menu;type:varchar(255);comment:当前路由为详情页时，需要高亮的菜单;NOT NULL" json:"active_menu"`
	Link        string    `gorm:"column:link;type:varchar(255);comment:外链地址;NOT NULL" json:"link"`
	IsHide      int       `gorm:"column:is_hide;type:tinyint(4);default:0;comment:是否隐藏;NOT NULL" json:"is_hide"`
	IsFull      int       `gorm:"column:is_full;type:tinyint(4);default:0;comment:是否全屏;NOT NULL" json:"is_full"`
	IsAffix     int       `gorm:"column:is_affix;type:tinyint(4);default:0;comment:是否固定在 tabs nav;NOT NULL" json:"is_affix"`
	IsKeepAlive int       `gorm:"column:is_keep_alive;type:tinyint(4);default:0;comment:是否缓存;NOT NULL" json:"is_keep_alive"`
	AuthBtn     string    `gorm:"column:auth_btn;type:text;comment:菜单对应的页面的按钮组标识 json数组;NOT NULL" json:"auth_btn"`
	Sort        int       `gorm:"column:sort;type:int(11);default:0;comment:排序;NOT NULL" json:"sort"`
	CreateTime  time.Time `gorm:"column:create_time;type:datetime;comment:创建时间" json:"create_time"`
	UpdateTime  time.Time `gorm:"column:update_time;type:datetime;comment:更新时间" json:"update_time"`
}

func (m *AdminMenu) TableName() string {
	return "admin_menu"
}
