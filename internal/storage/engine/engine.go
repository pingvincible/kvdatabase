package engine

type Engine struct {
	Storage map[string]string
}

func New() *Engine {
	return &Engine{
		Storage: make(map[string]string),
	}
}

func (e *Engine) Set(key, value string) {
	e.Storage[key] = value
}

func (e *Engine) Get(key string) string {
	return e.Storage[key]
}

func (e *Engine) Delete(key string) {
	delete(e.Storage, key)
}
