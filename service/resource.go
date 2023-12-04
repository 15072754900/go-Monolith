package service

import (
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
)

type Resource struct{}

func (s *Resource) GetTreeList(req req.PageQuery) []resp.ResourceVo {
	resources := resourceDao.GetListByKeyword(req.KeyWord)
	return s.resources2ResourceVos(resources)
}

// 很重要的一步骤
func (s *Resource) resources2ResourceVos(resources []model.Resource) []resp.ResourceVo {
	resList := make([]resp.ResourceVo, 0)
	// 找到父节点列表（parentId == 0）
	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)
	// 转换结构，[]model.Resource -> []resp.ResourceVo
	for _, item := range parentList {
		// 构建数据树
		resourceVo := s.resource2ResourceVo(item)
		resourceVo.Children = make([]resp.ResourceVo, 0)
		for _, child := range childrenMap[item.ID] {
			resourceVo.Children = append(resourceVo.Children, s.resource2ResourceVo(child))
		}
		resList = append(resList, resourceVo)
	}
	return resList
}

func getModuleList(resources []model.Resource) (list []model.Resource) {
	var parentList []model.Resource
	for _, r := range resources {
		if r.ParentId == 0 {
			parentList = append(parentList, r)
		}
	}
	return parentList
}

func getChildrenMap(resource []model.Resource) map[int][]model.Resource {
	childrenMap := make(map[int][]model.Resource)
	for _, r := range resource {
		if r.ParentId != 0 {
			childrenMap[r.ParentId] = append(childrenMap[r.ParentId], r)
		}
	}
	return childrenMap
}

func (*Resource) resource2ResourceVo(r model.Resource) resp.ResourceVo {
	return resp.ResourceVo{
		ID:            r.ID,
		Name:          r.Name,
		Url:           r.Url,
		RequestMethod: r.RequestMethod,
		IsAnonymous:   r.IsAnonymous,
		CreatedAt:     r.CreatedAt,
	}
}

func (*Resource) SaveOrUpdate(req req.SaveOrUpdateResource) (code int) {
	// 检查资源已存在
	existByName := dao.GetOne(model.Resource{}, "name", req.Name)
	if existByName.ID != 0 && existByName.ID != req.ID {
		return r.ERROR_RESOURCE_NAME_EXIST
	}
	if req.ID != 0 {
		// 已存在并更新
		oldResource := dao.GetOne(model.Resource{}, "id", req.ID)              // 用于修改权限
		dao.UpdatesMap(&model.Resource{}, utils.Struct2Map(req), "id", req.ID) // map可以更新零值
		// ！关联更新 casbin_rule 中的信息
		utils.Casbin.UpdatePolicy(
			[]string{"", oldResource.Url, oldResource.RequestMethod},
			[]string{"", req.Url, req.RequestMethod},
		)
	} else { // 新增
		data := utils.CopyProperties[model.Resource](req)
		dao.Create(&data)
		// ** 涉及前端的一个bug，在这里需要删除新增子节点的与其父节点之间的关联关系
		dao.Delete(model.Resource{}, "resource_id", data.ParentId)
	}
	return r.OK
}

// Delete 关联资源的处理问题，存在问题则返回报错，在子资源完全删除后才能删除父资源
func (*Resource) Delete(id int) (code int) {
	// 检查要删除的资源是否存在
	existResourceById := dao.GetOne(model.Resource{}, "id", id)
	if existResourceById.ID == 0 {
		return r.ERROR_RESOURCE_NOT_EXIST
	}
	// 检查 role_resource 下是否有数据
	existRoleResource := dao.GetOne(model.RoleResource{}, "resource_id", id)
	if existRoleResource.ResourceId != 0 {
		return r.ERROR_RESOURCE_USED_BY_ROLE
	}
	// 如果该 resource 是模块，检查其是否有子资源
	if existResourceById.ParentId == 0 {
		if dao.Count(model.Resource{}, "parent_id = ?", id) != 0 {
			return r.ERROR_RESOURCE_HAS_CHILDREN
		}
	}
	// 删除资源
	dao.Delete(model.Resource{}, "id", id)
	if existResourceById.Url != "" && existResourceById.RequestMethod != "" {
		utils.Casbin.DeletePermission(existResourceById.Url, existResourceById.RequestMethod)
	}
	return r.OK
}

func (*Resource) UpdateAnonymous(req req.UpdateAnonymous) (code int) {
	// 检查要更新的资源是否存在
	existById := dao.GetOne(model.Resource{}, "id", req.ID)
	if existById.ID == 0 {
		return r.ERROR_RESOURCE_NOT_EXIST
	}
	// 只更新 is_anonymous 字段
	dao.UpdatesMap(&model.Resource{}, map[string]any{"is_anonymous": *req.IsAnonymous}, "id", req.ID)
	// 关联 casbin_rule 中的 isAnonymous
	if *req.IsAnonymous == 0 {
		// 进行删除权限操作
		utils.Casbin.DeletePermissionForRole("anonymous", req.Url, req.RequestMethod)
	} else {
		utils.Casbin.AddPermissionForRole("anonymous", req.Url, req.RequestMethod)
	}
	return r.OK
}

// GetOptionList 获取树形数据选项，采取的还是遍历父节点的元素，并加入子节点在父节点中的id对应的元素
func (*Resource) GetOptionList() []resp.TreeOptionVo {
	resList := make([]resp.TreeOptionVo, 0)
	resources := dao.List([]model.Resource{}, "id , name, parent_id", "", "is_anonymous = ?", 0)
	parentList := getModuleList(resources)
	childrenMap := getChildrenMap(resources)
	for _, item := range parentList {
		// 构造 Children
		var childrenOptionVos []resp.TreeOptionVo
		for _, re := range childrenMap[item.ID] {
			childrenOptionVos = append(childrenOptionVos, resp.TreeOptionVo{
				ID:    re.ID,
				Label: re.Name,
			})
		}
		resList = append(resList, resp.TreeOptionVo{
			ID:       item.ID,
			Label:    item.Name,
			Children: childrenOptionVos,
		})
	}
	return resList
}
