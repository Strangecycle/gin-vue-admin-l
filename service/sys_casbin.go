package service

import (
	"errors"
	"gin-vue-admin-l/global"
	"gin-vue-admin-l/model"
	"gin-vue-admin-l/model/request"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/util"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"strings"
	"sync"
)

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func Casbin() *casbin.SyncedEnforcer {
	once.Do(func() {
		// casbin gorm 适配模式
		a, _ := gormadapter.NewAdapterByDB(global.GVA_DB)
		// 读取 model 配置文件
		syncedEnforcer, _ = casbin.NewSyncedEnforcer(global.GVA_CONFIG.Casbin.ModelPath, a)
		// 自定义匹配规则
		syncedEnforcer.AddFunction("ParamsMatch", ParamsMatchFunc)
	})

	// 加载自定义匹配规则
	_ = syncedEnforcer.LoadPolicy()
	return syncedEnforcer
}

func ParamsMatch(fullNameKey1 string, key2 string) bool {
	// 拿到中间件定义的 obj, 这里是过滤掉参数的请求路径
	key1 := strings.Split(fullNameKey1, "?")[0]
	// 剥离路径后再使用 casbin 的 keyMatch2 来匹配路径
	return util.KeyMatch2(key1, key2)
}

// 自定义匹配 casbin 规则函数, 入参参考 model 配置文件
func ParamsMatchFunc(args ...interface{}) (interface{}, error) {
	n1 := args[0].(string)
	n2 := args[1].(string)

	return ParamsMatch(n1, n2), nil
}

func ClearCasbin(v int, p ...string) bool {
	e := Casbin()
	success, _ := e.RemoveFilteredPolicy(v, p...)
	return success
}

func GetPolicyPathByAuthorityId(authId string) (pathMaps []request.CasbinInfo) {
	e := Casbin()
	list := e.GetFilteredPolicy(0, authId)
	for _, v := range list {
		pathMaps = append(pathMaps, request.CasbinInfo{
			Path:   v[1],
			Method: v[2],
		})
	}
	return pathMaps
}

func UpdateCasbin(authId string, infos []request.CasbinInfo) (err error) {
	// 先从数据库将当前 authorityId 的权限全部删除
	ClearCasbin(0, authId)
	rules := [][]string{}
	for _, v := range infos {
		cm := model.CasbinModel{
			Ptype:       "p",
			AuthorityId: authId,
			Path:        v.Path,
			Method:      v.Method,
		}
		rules = append(rules, []string{cm.AuthorityId, cm.Path, cm.Method})
	}
	e := Casbin()
	success, _ := e.AddPolicies(rules)
	if !success {
		return errors.New("存在相同api,添加失败,请联系管理员")
	}
	return nil
}

// 更新 Api 时，casbin 表中也同步更新
func UpdateCasbinApi(op string, np string, om string, nm string) (err error) {
	err = global.GVA_DB.Table("casbin_rule").Model(&model.CasbinModel{}).Where("v1 = ? AND v2 = ?", op, om).Updates(map[string]interface{}{
		"v1": np,
		"v2": nm,
	}).Error
	return err
}
