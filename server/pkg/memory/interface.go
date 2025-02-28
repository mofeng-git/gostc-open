package memory

import "time"

type AnyType interface {
	Marshal() []byte
	Unmarshal(data []byte) error
}

type Interface interface {
	SetStruct(key string, data any, duration time.Duration)
	GetStruct(key string, data any) error
	SetString(key string, value string, duration time.Duration)
	GetString(key string) (value string)
	Get(key string, data AnyType) error
	Set(key string, data AnyType, duration time.Duration)
	Del(key string)
	Sync()
}
