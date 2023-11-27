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
	"sort"
	"strconv"
	"strings"
	"time"
)

// 这里就是直接和数据库对接的地方、

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

func (*User) Update(req req.UpdateUser) int {
	userInfo := model.UserInfo{
		Universal: model.Universal{ID: req.UserInfoId},
		Nickname:  req.Nickname,
	}
	dao.Update(&userInfo)
	// 清空user_role关系并更新user_role列表
	dao.Delete(model.UserRole{}, "user_id = ?", req.UserInfoId)
	// 先删除旧的，再添加新的
	var userRoles []model.UserRole
	for _, id := range req.RoleIds {
		userRoles = append(userRoles, model.UserRole{
			RoleId: id,
			UserId: req.UserInfoId,
		})
	}
	dao.Create(&userRoles)
	return r.OK
}

func (*User) UpdateDisable(id, isDisable int) {
	dao.UpdatesMap(&model.UserInfo{}, map[string]any{"is_disable": isDisable}, "id", id)
}

func (*User) UpdatePassword(req req.UpdatePassword) int {
	// 判断用户是否存在
	if exist := checkUserExistByName(req.Username); !exist {
		return r.ERROR_USER_NOT_EXIST
	}
	// 执行存在时的更新服务
	m := map[string]any{"password": utils.Encryptor.BcryptHash(req.Password)}
	dao.UpdatesMap(&model.UserAuth{}, m, "username = ?", req.Password)
	return r.OK
}

func checkUserExistByName(username string) bool {
	existUser := dao.GetOne(model.UserAuth{}, "username = ?", username)
	return existUser.ID != 0
}

func (*User) UpdateCurrentPassword(req req.UpdateAdminPassword, id int) int {
	user := dao.GetOne(model.UserAuth{}, "id", id) // 获取用户信息
	if !user.IsEmpty() && utils.Encryptor.BcryptCheck(req.NewPassword, req.OldPassword) {
		user.Password = utils.Encryptor.BcryptHash(req.NewPassword)
		dao.Update(&user, "password")
		return r.OK
	} else {
		return r.ERROR_OLD_PASSWORD
	}
}

func (*User) UpdateCurrent(req req.UpdateCurrentUser) (code int) {
	user := utils.CopyProperties[model.UserInfo](req)
	dao.Update(&user, "nickname", "intro", "website", "avatar", "email")
	return r.OK
}

// GetOnlineList 查询当前在线用户，分页+条件搜索
func (*User) GetOnlineList(req req.PageQuery) resp.PageResult[[]resp.UserOnline] {
	onlineList := make([]resp.UserOnline, 0)

	keys := utils.Redis.Keys(KEY_USER + "*")
	for _, key := range keys {
		var sessionInfo dto.SessionInfo
		utils.Json.Unmarshal(utils.Redis.GetVal(key), &sessionInfo) // 这里进行解码并将数据保存

		// 查询关键字不为空，且不符查询条件
		if req.KeyWord != "" && !strings.Contains(sessionInfo.Nickname, req.KeyWord) {
			continue
		}

		onlineUser := utils.CopyProperties[resp.UserOnline](sessionInfo)
		onlineUser.UserIndoId = sessionInfo.UserInfoId // 一个个获取并保存
		onlineList = append(onlineList, onlineUser)
	}

	// 根据上次登陆时间进行排序
	sort.Slice(onlineList, func(i, j int) bool {
		return onlineList[i].LastLoginTime.Unix() > onlineList[j].LastLoginTime.Unix()
	})
	return resp.PageResult[[]resp.UserOnline]{
		Total: int64(len(keys)),
		List:  onlineList,
	}
}

// ForceOffline 在线与离线服务计算存储是在Redis中设置的
func (*User) ForceOffline(req req.ForceOfflineUser) (code int) {
	uuid := utils.Encryptor.MD5(req.IpAddress + req.Browser + req.OS)
	var sessionInfo dto.SessionInfo
	utils.Json.Unmarshal(utils.Redis.GetVal(KEY_USER+uuid), &sessionInfo)
	sessionInfo.IsOffline = 1
	utils.Redis.Del(KEY_USER + uuid)
	// 设置强制离线之后 Redis 中存储的 delete:xxx 时间和 Token 过期时间一致
	utils.Redis.Set(KEY_DELETE+uuid, utils.Json.Marshal(sessionInfo), time.Duration(config.Cfg.JWT.Expire)*time.Hour)
	return r.OK
}
