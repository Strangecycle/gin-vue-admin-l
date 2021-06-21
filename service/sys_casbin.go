package service

import (
	"github.com/casbin/casbin/v2"
	"sync"
)

var (
	syncedEnforcer *casbin.SyncedEnforcer
	once           sync.Once
)

func Casbin() *casbin.SyncedEnforcer {
	// once.Do(func() {
	// 	a, _ := gormadapter.NewAdapterByDB(global.GVA_DB)
	// 	casbin.NewSyncedEnforcer(global.GVA_CONFIG.Casbin, a)
	// })

	return syncedEnforcer
}
