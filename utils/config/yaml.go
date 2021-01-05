package config_util

import (
	"github.com/SmallTianTian/fresh-go/utils"
	y "gopkg.in/yaml.v2"
)

var yaml = func(content []byte, k string, v interface{}, father []string) []byte {
	var m map[interface{}]interface{}
	if err := y.Unmarshal(content, &m); err != nil {
		panic(err)
	}
	tmp := m
	for _, fat := range father {
		vm := make(map[interface{}]interface{})
		if tv, in := tmp[fat]; in {
			switch vr := tv.(type) {
			case interface{}:
				if vmt, ok := vr.(map[interface{}]interface{}); ok {
					vm = vmt
				}
			case map[interface{}]interface{}:
				vm = vr
			default:
				panic("Not support type.")
			}
			if vm[k] != nil && vm[k] != v {
				panic("Not support update vlaue")
			}
		}
		tmp = vm
		m[fat] = vm
	}
	tmp[k] = v
	bs, err := y.Marshal(m)
	utils.MustNotError(err)
	return bs
}
