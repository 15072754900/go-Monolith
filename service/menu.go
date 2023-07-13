package service

import (
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"sort"
)

type Menu struct{}

// GetUserMenuTree 获取某个用户的菜单列表（树形）
func (s *Menu) GetUserMenuTree(userInfoId int) []resp.UserMenuVo {
	// 从数据库查出用户菜单列表（非树形）
	menuList := menuDao.GetMenusByUserInfoId(userInfoId)
	// 过滤出一级菜单（parent_id = 0）
	firstLevelMenuList := s.getFirstLevelMenus(menuList)
	// 获取存储子菜单的 map
	menuChildrenMap := s.getMenuChildrenMap(menuList)
	// 将一级 Menu 列表转成 UserMenuVo 列表，只要处理其 Children
	return s.menus2UserMenuVos(firstLevelMenuList, menuChildrenMap)
}

// GetTreeList 获取菜单列表（树形）
func (s *Menu) GetTreeList(req req.PageQuery) []resp.MenuVo {
	// 从数据库中查询出菜单列表（非树形）,获取了菜单内容，但是要输出菜单形式
	menuList := menuDao.GetMenus(req)
	// 过滤出一级菜单(parent_id = 0)
	firstLevelMenuList := s.getFirstLevelMenus(menuList)
	// 获取存储子菜单的 map
	menuChildrenMap := s.getMenuChildrenMap(menuList)
	// 将一级 Menu 列表转成 MenuVo 列表，主要处理其 children
	return s.menus2MenuVos(firstLevelMenuList, menuChildrenMap)
}

// 筛选出一级菜单（parentId == 0 的菜单）
func (s *Menu) getFirstLevelMenus(menuList []model.Menu) []model.Menu {
	firstLevelMenus := make([]model.Menu, 0)
	for _, menu := range menuList {
		if menu.ParentId == 0 {
			firstLevelMenus = append(firstLevelMenus, menu)
		}
	}
	s.sortMenu(firstLevelMenus) // 以order降序排列
	return firstLevelMenus
}

// SortMenu 对菜单排序：以 orderNum 字段进行排序
func (*Menu) sortMenu(menus []model.Menu) {
	sort.Slice(menus, func(i, j int) bool {
		return menus[i].OrderNum < menus[j].OrderNum
	})
}

// 获取 map：key是菜单 ID,value 是该菜单对应的子菜单列表
func (*Menu) getMenuChildrenMap(menus []model.Menu) map[int][]model.Menu {
	mcMap := make(map[int][]model.Menu)
	for _, menu := range menus {
		if menu.ParentId != 0 { // 在数据结构里面就表示清楚了，意思是这是一个具备子菜单的内容，然后获取到响应的mcMap里面
			mcMap[menu.ParentId] = append(mcMap[menu.ParentId], menu)
		}
	}
	return mcMap
}

// 构建用户菜单的树形结构数据，[]model.Menu => []resp.UserMenuVo
func (s *Menu) menu2UserMenuVos(firstLevelMenuList []model.Menu, childrenMap map[int][]model.Menu) []resp.UserMenuVo {
	resList := make([]resp.UserMenuVo, 0)

	// 遍历一级 Menu，将其构造成 UserMenu
	for _, firstLevelMenu := range firstLevelMenuList {
		var userMenu resp.UserMenuVo           // 当前 [用户菜单]
		var userMenuChildren []resp.UserMenuVo // 当前 [用户菜单] 的 [子用户菜单]

		children := childrenMap[firstLevelMenu.ID] // 子菜单
		if len(children) > 0 {                     // 存在子菜单
			userMenu = s.menu2UserMenuVo(firstLevelMenu) // [菜单] -> [用户菜单]
			// userMenu.Path = "" // TODO! 外层一定是 “” 吗？
			s.sortMenu(children) // 对子菜单按照 OrderNum 排序
			// 遍历子菜单，将其构造成 用户菜单
			for _, child := range children {
				userMenuChildren = append(userMenuChildren, s.menu2UserMenuVo(child))
			}
		} else { // 没有子菜单，利用一节菜单构造一个用户菜单(Layout)，将原本的一级菜单作为子菜单变成新菜单的 children
			userMenu = resp.UserMenuVo{
				ID:        firstLevelMenu.ID,
				Path:      firstLevelMenu.Path,
				Name:      firstLevelMenu.Name,
				Component: firstLevelMenu.Component,
				OrderNum:  firstLevelMenu.OrderNum,
				Hidden:    firstLevelMenu.IsHidden == 1,
				KeepAlive: firstLevelMenu.KeepAlive == 1,
				Redirect:  firstLevelMenu.Redirect,
			}
			tmpUserMenu := s.menu2UserMenuVo(firstLevelMenu) // 就是转化格式，再进入下一个循环
			// tmpUserMenu.Path = "" // TODO ! 考虑一下
			userMenuChildren = append(userMenuChildren, tmpUserMenu)
		}
		userMenu.Children = userMenuChildren
		resList = append(resList, userMenu)
	}
	return resList
}

// TODO: 修改成使用 MapStruct?
// model.Menu => resp.MenuVo
func (*Menu) menu2MenuVo(menu model.Menu) resp.MenuVo {
	return resp.MenuVo{
		ID:        menu.ID,
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		CreatedAt: menu.CreatedAt,
		OrderNum:  menu.OrderNum,
		IsHidden:  int(menu.IsHidden),
		ParentId:  menu.ParentId,
		KeepAlive: menu.KeepAlive,
		Redirect:  menu.Redirect,
	}
}

// 构建菜单列表的树形数据结构数据，[]model.Menu => []resp.MenuVo
func (s *Menu) menus2MenuVos(firstLevelMenuList []model.Menu, menuChildrenMap map[int][]model.Menu) []resp.MenuVo {
	resMenuVos := make([]resp.MenuVo, 0)
	// 遍历一级菜单
	for _, firstLevelMenu := range firstLevelMenuList {
		// 尝试将一级菜单转 MenuV, 还需处理children
		menuVo := s.menu2MenuVo(firstLevelMenu)
		// 获取到当前菜单的子菜单
		menuChildren := menuChildrenMap[firstLevelMenu.ID]
		s.sortMenu(menuChildren)
		// 对子菜单进行 []menu => []menuVo
		childMenuVos := make([]resp.MenuVo, 0)
		for _, childMenu := range menuChildren {
			childMenuVos = append(childMenuVos, s.menu2MenuVo(childMenu))
		}
		menuVo.Children = childMenuVos
		resMenuVos = append(resMenuVos, menuVo)
		delete(menuChildrenMap, firstLevelMenu.ID) // 从 map 中删除以构建完的菜单
	}
	// 处理剩下的子菜单 TODO： 思考
	if len(menuChildrenMap) > 0 {
		var menuChildren []model.Menu
		// 这整个函数包括那个user的，都需要好好思考
		for _, v := range menuChildrenMap {
			menuChildren = append(menuChildren, v...)
		}
		s.sortMenu(menuChildren)
		for _, menu := range menuChildren {
			resMenuVos = append(resMenuVos, s.menu2MenuVo(menu))
		}
	}
	return resMenuVos
}

// 构建用户菜单的树形结构数据，[]model.Menu => []resp.UserMenuVo
func (s *Menu) menus2UserMenuVos(firstLevelMenuList []model.Menu, childrenMap map[int][]model.Menu) []resp.UserMenuVo {
	resList := make([]resp.UserMenuVo, 0)

	// 遍历一级 Menu，将其构造成 UserMenu
	for _, firstLevelMenu := range firstLevelMenuList {
		var userMenu resp.UserMenuVo           // 当前 [用户菜单]
		var userMenuChildren []resp.UserMenuVo // 当前 [用户菜单] 的 [子菜单用户]

		children := childrenMap[firstLevelMenu.ID] // 子菜单
		if len(children) > 0 {                     // 存在子菜单
			userMenu = s.menu2UserMenuVo(firstLevelMenu) // [菜单] -> [用户菜单]
			// userMenu.Path = ""
			s.sortMenu(children) // 对子菜单按照 OrderNum 排序
			// 遍历子菜单，将其构造成 用户菜单
			for _, child := range children {
				userMenuChildren = append(userMenuChildren, s.menu2UserMenuVo(child))
			}
		} else { // 没有子菜单，利用一级菜单构造一个用户菜单(layout)，将原本的一级菜单作为子菜单变成新菜单的 children
			userMenu = resp.UserMenuVo{
				ID:        firstLevelMenu.ID,
				Path:      firstLevelMenu.Path,
				Name:      firstLevelMenu.Name,
				Component: firstLevelMenu.Component,
				OrderNum:  firstLevelMenu.OrderNum,
				Hidden:    firstLevelMenu.IsHidden == 1,
				KeepAlive: firstLevelMenu.KeepAlive == 1,
				Redirect:  firstLevelMenu.Redirect,
			}
			tmpUserMenu := s.menu2UserMenuVo(firstLevelMenu)
			// tmpUserMenu.Path = "" // TODO ! 考虑一下
			userMenuChildren = append(userMenuChildren, tmpUserMenu)
		}
		userMenu.Children = userMenuChildren
		resList = append(resList, userMenu)
	}
	return resList
}

// model.Menu => resp.UserMenuVo
func (*Menu) menu2UserMenuVo(menu model.Menu) resp.UserMenuVo {
	return resp.UserMenuVo{
		ID:        menu.ID,
		Name:      menu.Name,
		Path:      menu.Path,
		Component: menu.Component,
		Icon:      menu.Icon,
		OrderNum:  menu.OrderNum,
		Hidden:    menu.IsHidden == 1,
		KeepAlive: menu.KeepAlive == 1,
		ParentId:  menu.ParentId,
		Redirect:  menu.Redirect,
	}
}
