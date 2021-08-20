package model

import (
	"gorm.io/gorm"
)

type AdminMenu struct {
	gorm.Model
	Id       int          `gorm:"-;primary_key;AUTO_INCREMENT"`
	ParentId int          `gorm:"type:bigint(20) unsigned;not null;default:0;comment:父级ID"`
	Type     int          `gorm:"type:tinyint(4) unsigned;not null;default:0;comment:类型"`
	Order    int          `gorm:"type:int(11) unsigned;not null;default:0;comment:排序,越大越前"`
	Title    string       `gorm:"type:varchar(50);not null;comment:菜单名"`
	Icon     string       `gorm:"type:varchar(50);not null;comment:图标"`
	Url      string       `gorm:"type:varchar(3000);not null;default:'';comment:路径"`
	Header   string       `gorm:"type:varchar(150);comment:header"`
	Roles    []AdminRoles `gorm:"-"`
	Sub      []AdminMenu  `gorm:"-"`
}

type AdminMenuTodo struct {
	Id        int              `json:"id"`
	ParentId  int              `json:"parent_id"`
	Type      int              `json:"type"`
	Order     int              `json:"order"`
	Title     string           `json:"title"`
	Icon      string           `json:"icon"`
	Url       string           `json:"url"`
	Header    string           `json:"header"`
	Roles     []AdminRolesTodo `json:"roles"`
	Sub       []AdminMenuTodo  `json:"sub"`
	CreatedAt string           `json:"created_at"`
	UpdatedAt string           `json:"updated_at"`
}

func ToAdminMenuTodo(menu AdminMenu) AdminMenuTodo {

	var AdminRolesTodos []AdminRolesTodo

	for i, _ := range menu.Roles {
		AdminRolesTodos = append(AdminRolesTodos, ToAdminRolesTodo(menu.Roles[i]))
	}

	return AdminMenuTodo{
		Id:        int(menu.ID),
		ParentId:  menu.ParentId,
		Type:      menu.Type,
		Order:     menu.Order,
		Title:     menu.Title,
		Icon:      menu.Icon,
		Url:       menu.Url,
		Header:    menu.Header,
		Roles:     AdminRolesTodos,
		CreatedAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: menu.UpdatedAt.Format("2006-01-02 15:04:05"),
	}
}
