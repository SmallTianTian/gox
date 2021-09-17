// 方便 context 中设置内容
// idea from https://github.com/grpc-ecosystem/go-grpc-middleware/tree/master/tags
package context

type MapTag struct {
	kv map[string]interface{}
}

func (t *MapTag) Set(key string, value interface{}) *MapTag {
	t.kv[key] = value
	return t
}

func (t *MapTag) Has(key string) bool {
	_, ok := t.kv[key]
	return ok
}

func (t *MapTag) Get(key string) interface{} {
	return t.kv[key]
}

func (t *MapTag) Values() map[string]interface{} {
	return t.kv
}

func NewMapTag() *MapTag {
	return &MapTag{make(map[string]interface{})}
}
