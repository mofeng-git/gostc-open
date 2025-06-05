package gost_engine

import (
	"server/repository/query"
)

func NodeStop(code string, msg string) {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return
	}
	engine.GetNode().Stop(msg)
}

func NodeIngress(tx *query.Query, code string) {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return
	}
	engine.GetNode().Ingress(tx)
}

func NodeConfig(tx *query.Query, code string) {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return
	}
	engine.GetNode().Config(tx)
}

func NodePortCheck(tx *query.Query, code string, port string) (async bool, allow bool) {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return
	}
	return engine.GetNode().PortCheck(tx, "", port)
}

func NodeAddDomain(tx *query.Query, code, domain, cert, key string, forceHttps int) {
	engine, ok := EngineRegistry.Get(code)
	if !ok {
		return
	}
	engine.GetNode().CustomDomain(tx, domain, cert, key, forceHttps)
}
