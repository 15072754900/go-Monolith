package dao

import "gin-blog-hufeng/model"

type Resource struct{}

func (*Resource) GetListById(ids []int) (list []model.Resource) {
	return List([]model.Resource{}, "url, request_method", "", "id in ?", ids)
}

func (*Resource) GetListByKeyword(keyword string) (list []model.Resource) {
	// 判断关键字是否存在，两种情况：存在则按需索取，不存在则全部获取
	if keyword != "" {
		list = List([]model.Resource{}, "*", "", "name like ?", "%"+keyword+"%")
	} else {
		list = List([]model.Resource{}, "*", "", "")
	}
	return
}
