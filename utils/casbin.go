package utils

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"log"
	"sync"
)

const CASBIN_UTIL_ERR_PREFIX = "utils/casbin.go -> "

var (
	cachedEnforcer *casbin.CachedEnforcer
	casbin_db      *gorm.DB
	once           sync.Once
)

type _casbin struct{}

var Casbin = new(_casbin)

func InitCasbin(db *gorm.DB) *casbin.CachedEnforcer {
	var err error

	once.Do(func() {
		casbin_db = db
		adapter, _ := gormadapter.NewAdapterByDB(db)

		// 方法一：从字符串中加载
		text := `
		[request_definition]
		r = sub, obj, act

		[policy_definition]
		p = sub, obj, act

		[role_definition]
		g = _, _

		[policy_effect]
		e = some(where (p.eft == allow))

		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
		`
		m, _ := model.NewModelFromString(text)
		cachedEnforcer, err = casbin.NewCachedEnforcer(m, adapter)
		if err != nil {
			log.Panic("Casbin 初始化失败：", err)
		}

		// 方法二：从配置文件中加载
		// cachedEnforcer, _ = casbin.NewCachedEnforcer("config/casbin.conf",adapter)

		cachedEnforcer.SetExpireTime(60 * 60)
		cachedEnforcer.EnableAutoSave(true)
		cachedEnforcer.LoadPolicy()
	})

	return cachedEnforcer
}

func (*_casbin) Enforcer() *casbin.CachedEnforcer {
	return cachedEnforcer
}

// LoadPolicy 重新加载策略，使得更新的策略生效
func (*_casbin) LoadPolicy() {
	// log.Println("重加载策略，使得更新的策略生效...")
	cachedEnforcer.LoadPolicy()
}

// UpdatePolicy 更新 API 权限
func (c *_casbin) UpdatePolicy(old []string, new []string) {
	_, err := cachedEnforcer.UpdatePolicy(old, new)
	if err != nil {
		Logger.Error(CASBIN_UTIL_ERR_PREFIX+"UpdatePolicy: ", zap.Error(err))
		panic(err)
	}
	c.LoadPolicy()

}

// DeletePermission 删除并更新 API 权限
func (c *_casbin) DeletePermission(permission ...string) {
	_, err := cachedEnforcer.DeletePermission(permission...)
	if err != nil {
		Logger.Error(CASBIN_UTIL_ERR_PREFIX+"DeletePermission: ", zap.Error(err))
		panic(err)
	}
	c.LoadPolicy()
}

func (c *_casbin) DeletePermissionForRole(roleLabel string, permissions ...string) {
	_, err := cachedEnforcer.DeletePermissionForUser(roleLabel, permissions...)
	if err != nil {
		Logger.Error(CASBIN_UTIL_ERR_PREFIX+"DeletePermissionForRole: ", zap.Error(err))
		panic(err)
	}
	c.LoadPolicy()
}

func (c *_casbin) AddPermissionForRole(roleLabel string, permissions ...string) {
	_, err := cachedEnforcer.AddPermissionForUser(roleLabel, permissions...)
	if err != nil {
		Logger.Error(CASBIN_UTIL_ERR_PREFIX+"AddPermissionForRole: ", zap.Error(err))
		panic(err)
	}
	c.LoadPolicy()
}

func (c *_casbin) BatchDeletePermissions(roleLabels []string) {
	for _, roleLabel := range roleLabels {
		_, err := cachedEnforcer.DeletePermissionsForUser(roleLabel) // 删除该角色对应的权限API
		if err != nil {
			Logger.Error(CASBIN_UTIL_ERR_PREFIX+"BatchDeletePermissions: ", zap.Error(err))
			panic(err)
		}
	}
	c.LoadPolicy()
}

func (c *_casbin) AddRoleForUser(user, role string) {
	_, err := cachedEnforcer.AddRoleForUser(user, role)
	if err != nil {
		Logger.Error(CASBIN_UTIL_ERR_PREFIX+"AddRoleForUser: ", zap.Error(err))
		panic(err)
	}
}

func (c *_casbin) UpdateRoleLabels(oldLabel, newLabel string) {
	err := casbin_db.Model(&gormadapter.CasbinRule{}).Where("v0 = ?", oldLabel).Update("v0", newLabel).Error
	if err != nil {
		Logger.Error(CASBIN_UTIL_ERR_PREFIX+"UpdateRoleLabels: ", zap.Error(err))
		panic(err)
	}
	c.LoadPolicy()
}

func (c *_casbin) AddPolicies(rules [][]string) {
	if len(rules) == 0 {
		return
	}
	_, err := cachedEnforcer.AddPolicies(rules)
	if err != nil {
		Logger.Error(CASBIN_UTIL_ERR_PREFIX+"AddPolicies: ", zap.Error(err))
		panic(err)
	}
	c.LoadPolicy()
}
