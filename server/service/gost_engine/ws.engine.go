package gost_engine

import (
	"encoding/json"
	"github.com/google/uuid"
	"github.com/lxzan/gws"
	"net/http"
	"sync"
)

var EngineRegistry = make(EngineRegister)

type EngineRegister map[string]*Engine

var engineRegister = &sync.RWMutex{}

func (reg EngineRegister) Set(engine *Engine) {
	engineRegister.Lock()
	defer engineRegister.Unlock()
	if engine.code == "" {
		return
	}
	reg[engine.code] = engine
}

func (reg EngineRegister) Get(code string) (*Engine, bool) {
	engineRegister.RLock()
	defer engineRegister.RUnlock()
	engine, ok := reg[code]
	return engine, ok
}

type Engine struct {
	code string
	conn *gws.Conn
}

func NewEngine(code string, w http.ResponseWriter, r *http.Request, event gws.Event) (*Engine, error) {
	upgrader := gws.NewUpgrader(event, &gws.ServerOption{
		ParallelEnabled:   true,                                 // 开启并行消息处理
		Recovery:          gws.Recovery,                         // 开启异常恢复
		PermessageDeflate: gws.PermessageDeflate{Enabled: true}, // 开启压缩
	})
	upgrade, err := upgrader.Upgrade(w, r)
	if err != nil {
		return nil, err
	}
	msgRegistry.New(code)
	return &Engine{
		code: code,
		conn: upgrade,
	}, nil
}

func (e *Engine) ReadLoop() {
	e.conn.ReadLoop()
}

func WriteMessage(code string, req MessageData) {
	msgRegistry.PushMessage(code, req)
}

func (e *Engine) Close(msg string) {
	marshal, _ := json.Marshal(NewMessage(uuid.NewString(), "stop", map[string]string{
		"reason": msg,
	}))
	_ = e.conn.WriteString(string(marshal))
}
