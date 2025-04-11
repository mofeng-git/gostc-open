package fs

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"sync"
	"time"
)

// FsStorage 有序存储接口
type FsStorage interface {
	Insert(key string, value interface{}) error
	Update(key string, value interface{}) error
	Delete(key string) error
	Get(key string) (interface{}, bool)
	List() []KeyValuePair
	ListKeys() []string
	SetAutoPersistInterval(interval time.Duration)
	PersistNow() error
}

// KeyValuePair 键值对结构
type KeyValuePair struct {
	Key   string
	Value interface{}
}

// FileStorage 有序内存文件存储实现
type FileStorage struct {
	data              []KeyValuePair // 保持顺序的切片
	index             map[string]int // 快速查找索引
	filename          string
	mu                sync.RWMutex
	autoPersistTicker *time.Ticker
	stopChan          chan struct{}
	persistInterval   time.Duration
	dirty             bool
}

// NewFileStorage 创建有序存储实例
func NewFileStorage(filename string) (*FileStorage, error) {
	storage := &FileStorage{
		data:            make([]KeyValuePair, 0),
		index:           make(map[string]int),
		filename:        filename,
		stopChan:        make(chan struct{}),
		persistInterval: 0,
		dirty:           false,
	}

	if err := storage.loadFromFile(); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to load data from file: %v", err)
		}
	}

	return storage, nil
}

// 从文件加载数据
func (m *FileStorage) loadFromFile() error {
	//file, err := os.Open(m.filename)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	file, err := os.ReadFile(m.filename)
	if err != nil {
		return err
	}

	// 临时结构体用于加载
	var tempData []KeyValuePair
	decoder := json.NewDecoder(bytes.NewReader(file))
	if err := decoder.Decode(&tempData); err != nil {
		return err
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	m.data = tempData
	// 重建索引
	m.index = make(map[string]int)
	for i, pair := range m.data {
		m.index[pair.Key] = i
	}

	return nil
}

// SetAutoPersistInterval 设置自动持久化间隔
func (m *FileStorage) SetAutoPersistInterval(interval time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.autoPersistTicker != nil {
		m.autoPersistTicker.Stop()
		close(m.stopChan)
		m.stopChan = make(chan struct{})
	}

	m.persistInterval = interval

	if interval > 0 {
		m.autoPersistTicker = time.NewTicker(interval)
		go func() {
			for {
				select {
				case <-m.autoPersistTicker.C:
					m.mu.Lock()
					if m.dirty {
						if err := m.persist(); err != nil {
							fmt.Printf("Auto persist failed: %v\n", err)
						} else {
							m.dirty = false
						}
					}
					m.mu.Unlock()
				case <-m.stopChan:
					return
				}
			}
		}()
	}
}

// Insert 插入新对象
func (m *FileStorage) Insert(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if _, exists := m.index[key]; exists {
		return fmt.Errorf("key %s already exists", key)
	}

	m.data = append(m.data, KeyValuePair{Key: key, Value: value})
	m.index[key] = len(m.data) - 1
	m.dirty = true

	return m.maybePersist()
}

// Update 更新现有对象
func (m *FileStorage) Update(key string, value interface{}) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	idx, exists := m.index[key]
	if !exists {
		return fmt.Errorf("key %s not found", key)
	}

	m.data[idx].Value = value
	m.dirty = true

	return m.maybePersist()
}

// Delete 删除对象
func (m *FileStorage) Delete(key string) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	idx, exists := m.index[key]
	if !exists {
		return fmt.Errorf("key %s not found", key)
	}

	// 从切片中删除
	m.data = append(m.data[:idx], m.data[idx+1:]...)
	// 从索引中删除
	delete(m.index, key)
	// 更新索引中受影响的位置
	for i := idx; i < len(m.data); i++ {
		m.index[m.data[i].Key] = i
	}
	m.dirty = true

	return m.maybePersist()
}

// maybePersist 根据配置决定是否持久化
func (m *FileStorage) maybePersist() error {
	if m.persistInterval == 0 {
		return m.persist()
	}
	return nil
}

// Get 获取单个对象
func (m *FileStorage) Get(key string) (interface{}, bool) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	idx, exists := m.index[key]
	if !exists {
		return nil, false
	}
	return m.data[idx].Value, true
}

// List 获取所有对象(按插入顺序)
func (m *FileStorage) List() []KeyValuePair {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 返回副本
	result := make([]KeyValuePair, len(m.data))
	copy(result, m.data)
	return result
}

// ListKeys 获取所有键(按插入顺序)
func (m *FileStorage) ListKeys() []string {
	m.mu.RLock()
	defer m.mu.RUnlock()

	keys := make([]string, len(m.data))
	for i, pair := range m.data {
		keys[i] = pair.Key
	}
	return keys
}

// PersistNow 立即持久化到磁盘
func (m *FileStorage) PersistNow() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if !m.dirty {
		return nil
	}

	if err := m.persist(); err != nil {
		return err
	}
	m.dirty = false
	return nil
}

// persist 内部持久化方法
func (m *FileStorage) persist() error {
	//file, err := os.Create(m.filename)
	//if err != nil {
	//	return err
	//}
	//defer file.Close()
	marshal, _ := json.Marshal(m.data)
	return os.WriteFile(m.filename, marshal, 0644)
	//encoder := json.NewEncoder(file)
	//encoder.SetIndent("", "  ")
	//return encoder.Encode(m.data)
}

// Close 关闭存储
func (m *FileStorage) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.autoPersistTicker != nil {
		m.autoPersistTicker.Stop()
		close(m.stopChan)
	}

	if m.dirty {
		if err := m.persist(); err != nil {
			fmt.Printf("Failed to persist on close: %v\n", err)
		}
	}
}
