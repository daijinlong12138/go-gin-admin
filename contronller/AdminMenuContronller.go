package contronller

import (
	"github.com/gin-gonic/gin"
	"go-gin-admin/common"
	"go-gin-admin/model"
	"go-gin-admin/response"
	"gorm.io/gorm"
	"sort"
	"strconv"
	"strings"
)

func MenuNew(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	parent_id := c.PostForm("parent_id") //父级ID -- 顶级默认0
	title := c.PostForm("title")         //菜单名
	header := c.PostForm("header")       //header
	icon := c.PostForm("icon")           //图标
	url := c.PostForm("url")             //路径
	role_ids := c.PostForm("role_ids")   //角色,逗号分割

	order := 0
	if ordert, isExist := c.GetPostForm("order"); isExist == true {
		orderti, _ := strconv.Atoi(ordert)
		if orderti < 1 {
			order = 0
		} else {
			order = orderti
		}
	}
	//校验
	parent_idInt, _ := strconv.Atoi(parent_id)
	if parent_idInt < 0 {
		response.Fail(c, "父级ID不能为负数", nil)
		return
	}
	if len(title) == 0 {
		response.Fail(c, "菜单名不能为空", nil)
		return
	}

	var roleIdArr []string
	if len(role_ids) != 0 {
		roleIdArr = strings.Split(role_ids, `,`)
	}

	if len(roleIdArr) > 0 {
		for _, roleId := range roleIdArr {
			roleIdint, _ := strconv.Atoi(roleId)
			roles := findRolesById(DB, roleIdint)
			if roles.ID == 0 {
				response.Fail(c, roleId+"不存在角色", nil)
				return
			}
		}
	}

	tx := DB.Begin()
	//创建用户

	newMenu := model.AdminMenu{
		ParentId: parent_idInt,
		Title:    title,
		Header:   header,
		Icon:     icon,
		Url:      url,
		Order:    order,
	}
	err := tx.Create(&newMenu).Error
	if err != nil || newMenu.ID == 0 {
		response.Fail(c, "创建失败", nil)
		tx.Rollback()
		return
	}

	if len(roleIdArr) > 0 {
		//批量添加
		var arr []model.AdminRoleMenu
		arr = make([]model.AdminRoleMenu, len(roleIdArr))
		for k, roleId := range roleIdArr {
			roleIdint, _ := strconv.Atoi(roleId)
			arr[k].RoleId = roleIdint
			arr[k].MenuId = int(newMenu.ID)
		}
		err = model.AddRoleMenus(tx, arr)
		if err != nil {
			response.Fail(c, "用户添加角色失败", nil)
			common.LogError(c, "用户添加角色失败: "+err.Error())
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	response.Success(c, nil, "用户创建成功")
}

func MenuDetail(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id") //ID
	var menu model.AdminMenu

	DB.Where("id = ?", id).Find(&menu)
	if menu.ID == 0 {
		response.Fail(c, "不存在", nil)
		return
	}
	//角色
	DB.Raw("select admin_roles.id as id, admin_roles.name as name, admin_roles.slug as slug, admin_roles.created_at as created_at,admin_roles.updated_at as updated_at "+
		"from admin_role_menus "+
		"left join admin_roles on admin_roles.id = admin_role_menus.role_id where admin_role_menus.menu_id = ?", menu.ID).Scan(&menu.Roles)

	response.Success(c, model.ToAdminMenuTodo(menu), "成功")
}

func MenuDelete(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id") //角色ID
	var menu model.AdminMenu
	Id, _ := strconv.Atoi(id)
	DB.Where("id = ?", id).Find(&menu)
	if menu.ID == 0 {
		response.Fail(c, "不存在", nil)
		return
	}

	tx := DB.Begin()
	//菜单表
	err := tx.Where("id = ?", Id).Unscoped().Delete(&model.AdminMenu{}).Error
	if err != nil {
		tx.Rollback()
		response.Fail(c, "删除失败", nil)
		common.LogError(c, "删除失败: "+err.Error())
		return
	}
	//删除自身及所有下级菜单
	check, err := dealSubMenu(tx, Id)
	if !check {
		tx.Rollback()
		response.Fail(c, "删除失败:"+err.Error(), nil)
		common.LogError(c, "删除失败: "+err.Error())
		return
	}

	tx.Commit()
	response.Success(c, nil, "菜单删除成功")
}

func MenuInfo(c *gin.Context) {

	menuList := make([]model.AdminMenu, 0)
	db := common.GetDB()
	//展示菜单的树形结构
	//1.获取全部菜单数据，order 默认0，越大排序越前
	err := db.Find(&menuList).Error
	if err != nil {
		response.Fail(c, "查询失败", nil)
		common.LogError(c, "查询失败: "+err.Error())
		return
	}
	//2.数据处理
	// 查找相同上级目录--组成数组，并排序
	tree := MenuList(menuList, 0)

	response.Success(c, tree, "成功")
}

func MenuEdit(c *gin.Context) {
	DB := common.GetDB()
	//获取参数
	id := c.PostForm("id")                        //菜单id
	parent_id := c.DefaultQuery("parent_id", "0") //父级ID -- 顶级默认0
	title := c.PostForm("title")                  //菜单名
	header := c.PostForm("header")                //header
	icon := c.PostForm("icon")                    //图标
	url := c.PostForm("url")                      //路径
	order := c.PostForm("order")                  //排序
	role_ids := c.PostForm("role_ids")            //角色,逗号分割

	parent_idInt, _ := strconv.Atoi(parent_id)
	idInt, _ := strconv.Atoi(id)

	//校验
	if parent_idInt < 0 {
		response.Fail(c, "父级ID不能为负数", nil)
		return
	}
	if len(title) == 0 {
		response.Fail(c, "菜单名不能为空", nil)
		return
	}

	menu := findMenuById(DB, idInt)
	if menu.ID != uint(idInt) {
		response.Fail(c, "不存在菜单", nil)
		return
	}

	var roleIdArr []string
	if len(role_ids) != 0 {
		roleIdArr = strings.Split(role_ids, `,`)
	}

	if len(roleIdArr) > 0 {
		for _, roleId := range roleIdArr {
			roleIdint, _ := strconv.Atoi(roleId)
			roles := findRolesById(DB, roleIdint)
			if roles.ID == 0 {
				response.Fail(c, roleId+"不存在角色", nil)
				return
			}
		}
	}

	//不能对移动到当前自身及其子菜单下
	menuList := make([]model.AdminMenu, 0)
	db := common.GetDB()
	//展示菜单的树形结构
	//1.获取全部菜单数据，order 默认0，越大排序越前
	err := db.Find(&menuList).Error
	if err != nil {
		response.Fail(c, "查询失败", nil)
		common.LogError(c, "查询失败: "+err.Error())
	}
	//2.数据处理
	// 查找相同上级目录--组成数组，并排序
	tree := MenuList(menuList, idInt)
	/*fmt.Println(tree)
	fmt.Println(findIdForTree(tree,parent_idInt))*/
	if findIdForTree(tree, parent_idInt) {
		response.Fail(c, "查询失败", nil)
	}

	tx := DB.Begin()

	data := make(map[string]interface{})
	data["parent_id"] = parent_id
	data["title"] = title
	data["header"] = header
	data["icon"] = icon
	data["url"] = url
	data["order"] = order

	if err = tx.Model(&model.AdminMenu{}).Where("id = ?", id).Updates(&data).Error; err != nil {
		response.Fail(c, "更新失败", nil)
		common.LogError(c, "更新失败: "+err.Error())
		tx.Rollback()
		return
	}

	//批量删除菜单对应角色
	err = tx.Where("menu_id = ?", id).Unscoped().Delete(&model.AdminRoleMenu{}).Error
	if err != nil {
		response.Fail(c, "删除菜单角色失败", nil)
		common.LogError(c, "删除菜单角色失败: "+err.Error())
		tx.Rollback()
		return
	}
	if len(roleIdArr) > 0 {
		//批量添加
		var arr []model.AdminRoleMenu
		arr = make([]model.AdminRoleMenu, len(roleIdArr))
		for k, rolesId := range roleIdArr {
			rolesIdint, _ := strconv.Atoi(rolesId)
			arr[k].RoleId = rolesIdint
			MenuIdint, _ := strconv.Atoi(id)
			arr[k].MenuId = MenuIdint
		}
		err = model.AddRoleMenus(tx, arr)
		if err != nil {
			response.Fail(c, "更新角色失败", nil)
			common.LogError(c, "更新角色失败: "+err.Error())
			tx.Rollback()
			return
		}
	}

	tx.Commit()
	response.Success(c, nil, "成功")
}

func findMenuById(db *gorm.DB, id int) model.AdminMenu {
	var menu model.AdminMenu
	db.Where("id = ?", id).Find(&menu)
	return menu
}

// 删除成功返回true,失败false
func dealSubMenu(tx *gorm.DB, pid int) (bool, error) {
	var menus []model.AdminMenu
	//删除自身
	err := tx.Where("id = ?", pid).Unscoped().Delete(&model.AdminMenu{}).Error
	if err != nil {
		return false, err
	}
	//删除关联表
	err = tx.Where("menu_id = ?", pid).Unscoped().Delete(&model.AdminRoleMenu{}).Error
	if err != nil {
		return false, err
	}
	//查找下级
	err = tx.Where("parent_id = ?", pid).Find(&menus).Error
	if err != nil {
		return false, err
	}

	if len(menus) != 0 {
		for _, menu := range menus {
			//fmt.Println(menu,menu.ID)
			check, err := dealSubMenu(tx, int(menu.ID))
			if !check {
				return false, err
			}
		}
	}
	return true, nil
}

func MenuList(menus []model.AdminMenu, pid int) []model.AdminMenuTodo {

	var menuPidArr []model.AdminMenuTodo
	//1.把相同父类id菜单找出来
	for _, menu := range menus {
		if menu.ParentId == pid {
			child := MenuList(menus, int(menu.ID))
			node := model.AdminMenuTodo{
				Id:        int(menu.ID),
				ParentId:  menu.ParentId,
				Type:      menu.Type,
				Order:     menu.Order,
				Title:     menu.Title,
				Icon:      menu.Icon,
				Url:       menu.Url,
				Header:    menu.Header,
				CreatedAt: menu.CreatedAt.Format("2006-01-02 15:04:05"),
				UpdatedAt: menu.UpdatedAt.Format("2006-01-02 15:04:05"),
			}
			node.Sub = child
			menuPidArr = append(menuPidArr, node)
		}
	}
	//2.然后排序
	sort.Slice(menuPidArr, func(i, j int) bool {
		return menuPidArr[i].Order > menuPidArr[j].Order
	})

	return menuPidArr
}

func findIdForTree(menus []model.AdminMenuTodo, id int) bool {
	for _, menu := range menus {
		if menu.Id == id {
			return true
		}
	}
	return false
}
