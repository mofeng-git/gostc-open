package model

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

const (
	ALLOW_EDIT = 1 // 允许修改
	DENY_EDIT  = 2 // 拒绝修改

	ALLOW_DEL = 1 // 允许删除
	DENY_DEL  = 2 // 拒绝删除
)

type Base struct {
	Id        int    `gorm:"primaryKey"`
	Code      string `gorm:"column:code;size:100;uniqueIndex;comment:code"`
	AllowEdit int    `gorm:"column:allow_edit;default:1;size:1;comment:是否可编辑"`
	AllowDel  int    `gorm:"column:allow_del;default:1;size:1;comment:是否可删除"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

func (b *Base) BeforeCreate(tx *gorm.DB) error {
	if b.Code == "" {
		b.Code = uuid.NewString()
	}
	return nil
}

func (b *Base) BeforeUpdate(tx *gorm.DB) (err error) {
	if b.AllowEdit == DENY_EDIT {
		return errors.New("该数据禁止删除")
	}
	return nil
}

func (b *Base) BeforeDelete(tx *gorm.DB) (err error) {
	if b.AllowDel == DENY_DEL {
		return errors.New("该数据禁止修改")
	}
	return nil
}

type ArrayStr []string

func (arr *ArrayStr) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("数据处理错误")
	}
	return json.Unmarshal(bytes, arr)
}

func (arr *ArrayStr) Value() (driver.Value, error) {
	marshal, err := json.Marshal(arr)
	if err != nil {
		return nil, err
	}
	return string(marshal), nil
}

type Map map[string]any

func (m *Map) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("数据处理错误")
	}
	return json.Unmarshal(bytes, m)
}

func (m *Map) Value() (driver.Value, error) {
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, err
	}
	return string(marshal), nil
}
