package service

import (
	"fmt"
	"gin-blog-hufeng/config"
	"gin-blog-hufeng/dao"
	"gin-blog-hufeng/model"
	"gin-blog-hufeng/model/dto"
	"gin-blog-hufeng/model/req"
	"gin-blog-hufeng/model/resp"
	"gin-blog-hufeng/utils"
	"gin-blog-hufeng/utils/r"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type User struct{}

func (*User) Login(c *gin.Context, username, password string) (loginVo resp.LoginVO, code int) {
	// 检查用户是否存在
	userAuth := dao.GetOne(model.UserAuth{}, "username", username)
	fmt.Println("this is userAuth", userAuth)
	if userAuth.ID == 0 {
		return loginVo, r.ERROR_USER_NOT_EXIST
	}
	// 检查密码是否正确
	if !utils.Encryptor.BcryptCheck(password, userAuth.Password) {
		return loginVo, r.ERROR_USER_NOT_EXIST
	}
	// 获取用户详细信息 DTO
	userDetailDTO := convertUserDetailDTO(userAuth, c)

	// 登录信息正确，生成token
	// TODO: 目前之歌用户设定一个角色，获取第一个值就行，后期优化：给用户设置多个角色
	// UUID 生成方法：IP + 浏览信息 + 操作系统信息
	uuid := utils.Encryptor.MD5(userDetailDTO.IpAddress + userDetailDTO.Browser + userDetailDTO.OS)
	token, err := utils.GetJWT().GetToken(userAuth.ID, userDetailDTO.RoleLabels[0], uuid)
	fmt.Println("登陆成功,开始鉴权", uuid)
	if err != nil {
		utils.Logger.Info("登录时生成 Token 错误: ", zap.Error(err))
		return loginVo, r.ERROR_TOKEN_CREATE
	}
	userDetailDTO.Token = token
	// 更新用户验证信息：ip 信息 + 上次登录时间
	dao.Update(&model.UserAuth{
		Universal:     model.Universal{ID: userAuth.ID},
		IpAddress:     userDetailDTO.IpAddress,
		IpSource:      userDetailDTO.IpSource,
		LastLoginTime: userDetailDTO.LastLoginTime,
	}, "ip_address", "ip_source", "last_login_time")

	// 保存用户信息到 Session 和 Redis 中
	session := sessions.Default(c) // 获取session的快照
	// ! sessions 中只能存储字符串
	sessionInfoStr := utils.Json.Marshal(dto.SessionInfo{UserDetailDTO: userDetailDTO})
	session.Set(KEY_USER+uuid, sessionInfoStr)
	utils.Redis.Set(KEY_USER+uuid, sessionInfoStr, time.Duration(config.Cfg.Session.MaxAge)*time.Second)
	fmt.Println("设置完成 Redis")
	session.Save()
	// 后续必须学会这些的用法

	return userDetailDTO.LoginVO, r.OK
}

// 转化 UserDetailDTO
func convertUserDetailDTO(userAuth model.UserAuth, c *gin.Context) dto.UserDetailDTO {
	// 获取 IP 相关信息 FIXME：好像无法读取到 IP 信息
	ipAddress := utils.IP.GetIpAddress(c)
	ipSource := utils.IP.GetIpSource(ipAddress)
	browser, os := "unknown", "unknown"

	if userAgent := utils.IP.GetUserAgent(c); userAgent != nil { // 用户信息是在获取 IP 地址时得到的
		browser = userAgent.Name + " " + userAgent.Version.String()
		os = userAgent.OS + " " + userAgent.OSVersion.String()
	}

	// 获取用户详细信息
	userInfo := dao.GetOne(&model.UserInfo{}, "id", userAuth.ID)
	// FIXME：获取该用户对应的角色，没有角色默认是"test"
	rolelabels := roleDao.GetLabelsByUserInfoId(userInfo.ID)
	if len(rolelabels) == 0 {
		rolelabels = append(rolelabels, "test")
	}
	// 用户点赞 Set
	articleLikeSet := utils.Redis.SMembers(KEY_ARTICLE_USER_LIKE_SET + strconv.Itoa(userInfo.ID))
	commentLikeSet := utils.Redis.SMembers(KEY_COMMENT_USER_LIKE_SET + strconv.Itoa(userInfo.ID))

	return dto.UserDetailDTO{
		LoginVO: resp.LoginVO{
			ID:             userAuth.ID,
			UserInfoId:     userInfo.ID,
			Email:          userInfo.Email,
			LoginType:      userAuth.LoginType,
			Username:       userAuth.Username,
			Nickname:       userInfo.Nickname,
			Avatar:         userInfo.Avatar,
			Intro:          userInfo.Intro,
			Website:        userInfo.Website,
			IpAddress:      ipAddress,
			IpSource:       ipSource,
			LastLoginTime:  time.Now(),
			ArticleLikeSet: articleLikeSet,
			CommentLikeSet: commentLikeSet,
		},
		Password:   userAuth.Password,
		RoleLabels: rolelabels,
		IsDisable:  userInfo.IsDisable,
		Browser:    browser,
		OS:         os,
	}
}

// GetInfo TODO: 优化
func (*User) GetInfo(id int) resp.UserInfoVO {
	var userInfo model.UserInfo
	dao.GetOne(&userInfo, "id", id)

	data := utils.CopyProperties[resp.UserInfoVO](userInfo)
	data.ArticleLikeSet = utils.Redis.SMembers(KEY_ARTICLE_USER_LIKE_SET + strconv.Itoa(id)) // 这里的strconv使用可以借鉴一下
	data.CommentLikeSet = utils.Redis.SMembers(KEY_COMMENT_USER_LIKE_SET + strconv.Itoa(id))
	return data
}

func (*User) GetList(req req.GetUsers) resp.PageResult[[]resp.UserVO] { // 泛型的使用，结构的泛型
	count := userDao.GetCount(req)
	list := userDao.GetList(req)
	return resp.PageResult[[]resp.UserVO]{
		PageSize: req.PageSize,
		PageNum:  req.PageNum,
		Total:    count,
		List:     list,
	}
}
