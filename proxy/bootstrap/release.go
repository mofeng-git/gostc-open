package bootstrap

import "proxy/global"

var releaseFunc []func()

func Release() {
	global.Logger.Info("exit, release")
	for i := 0; i < len(releaseFunc); i++ {
		fn := releaseFunc[i]
		if fn == nil {
			continue
		}
		fn()
	}
}
