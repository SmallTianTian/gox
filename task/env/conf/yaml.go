package conf

import (
	y "gopkg.in/yaml.v2"
	"tianxu.xin/gox/internal/util"
)

func yaml(content []byte, k string, v interface{}, fathers []string) []byte {
	// 反序列化得到当前配置的 map 结果
	var mp map[interface{}]interface{}
	err := y.Unmarshal(content, &mp)
	util.MustNotError(err)
	if mp == nil {
		mp = make(map[interface{}]interface{})
	}

	cur := mp
	for _, fat := range fathers {
		cm, ok := cur[fat].(map[interface{}]interface{})
		if !ok {
			cm = make(map[interface{}]interface{})
			cur[fat] = cm
		}
		cur = cm
	}
	// 如果这里有值，则不更改
	if cur[k] != nil {
		return content
	}

	cur[k] = v
	bs, err := y.Marshal(mp)
	util.MustNotError(err)
	return bs
}
