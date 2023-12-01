package v1

import "gin-blog-hufeng/service"

var (
	userService     service.User
	blogInfoService service.BlogInfo
	categoryService service.Category
	tagService      service.Tag
	menuService     service.Menu
	articleService  service.Article
	commentService  service.Comment
	messageService  service.Message
)
