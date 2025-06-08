package service

import (
	"github.com/lesismal/arpc/codec"
	arpcLog "github.com/lesismal/arpc/log"
	"gopkg.in/yaml.v2"
)

type Service interface {
	Start() error
	Stop()
	IsRunning() bool
}

func init() {
	codec.SetCodec(&YAMLCodec{})
	arpcLog.SetLevel(arpcLog.LevelError)
}

type YAMLCodec struct {
}

func (y *YAMLCodec) Marshal(v interface{}) ([]byte, error) {
	return yaml.Marshal(v)
}

func (y *YAMLCodec) Unmarshal(data []byte, v interface{}) error {
	return yaml.Unmarshal(data, v)
}
