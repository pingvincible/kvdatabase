package engine

import "sync"

type Engine struct {
	mutex   sync.RWMutex
	Storage map[string]string
}

func New() *Engine {
	return &Engine{
		Storage: make(map[string]string),
	}
}

func (e *Engine) Set(key, value string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.Storage[key] = value
}

func (e *Engine) Get(key string) string {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return e.Storage[key]
}

func (e *Engine) Delete(key string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.Storage, key)
}
