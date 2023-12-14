package routes

import v1 "gin-blog-hufeng/api/v1"

// 后台接口
var (
	userAuthAPI     v1.UserAuth
	blogInfoAPI     v1.BlogInfo
	userAPI         v1.User
	categoryAPI     v1.Category
	tagAPI          v1.Tag
	menuAPI         v1.Menu
	articleAPI      v1.Article
	commentAPI      v1.Comment
	messageAPI      v1.Message
	linkAPI         v1.Link
	resourceAPI     v1.Resource
	roleAPI         v1.Role
	operationLogAPI v1.OperationLog
	pageAPI         v1.Page
)

// 前台接口
