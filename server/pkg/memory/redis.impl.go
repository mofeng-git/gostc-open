package memory

import (
	"time"
)

type Redis struct {
}

func (r Redis) SetStruct(key string, data any, duration time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) GetStruct(key string, data any) error {
	//TODO implement me
	panic("implement me")
}

func (r Redis) SetString(key string, value string, duration time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) GetString(key string) (value string) {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Get(key string, data AnyType) error {
	//TODO implement me
	panic("implement me")
}

func (r Redis) Set(key string, data AnyType, duration time.Duration) {
	//TODO implement me
	panic("implement me")
}
func (r Redis) Del(key string) {
	panic("implement me")
}

func (r Redis) Sync() {
	//TODO implement me
	panic("implement me")
}
