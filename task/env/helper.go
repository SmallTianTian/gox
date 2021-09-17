// 本文件不应该被正式代码所引用
// 仅仅是初始化项目的帮助文件
// 方便测试
package env

import (
	"context"

	"tianxu.xin/gox/internal/test"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task/env/mod"
)

const (
	defaultModule = "tianxu.xin/gox/demo"
)

type option struct {
	module   string
	vendor   bool
	resetCtx bool
}

type HelperOption func(*option)

// WithProject 设置项目 module
// 特殊情况，可传递指定 module，取第一个值。
// 否则使用 defaultModule
func WithProject(module ...string) HelperOption {
	return func(o *option) {
		if len(module) != 0 {
			o.module = module[0]
		} else {
			o.module = defaultModule
		}
	}
}

func UseVendor() HelperOption {
	return func(o *option) {
		o.vendor = true
	}
}

func ResetCtx() HelperOption {
	return func(o *option) {
		o.resetCtx = true
	}
}

// ProjectCtxHelper 项目 ctx 助手
// opts 中不包含 WithProject 方法，
// 将不会创建项目，其他设置均无效。
// 在存在 WithProject 方法的时候，
// 其他参数将生效
func ProjectCtxHelper(projectPath string, opts ...HelperOption) (ctx context.Context) {
	o := &option{}
	for _, opt := range opts {
		opt(o)
	}

	ctx = test.GetTestCtx(projectPath)
	if o.module == "" {
		return
	}
	if o.resetCtx {
		defer func() {
			ctx = test.GetTestCtx(projectPath)
		}()
	}

	// 新建项目
	util.MustExtractConf(ctx).GoEnv.Module = o.module
	NewFreshProjectTask().Exec(ctx)  // nolint
	mod.NewGoModInitTask().Exec(ctx) // nolint

	if o.vendor {
		mod.NewGoModVendorTask().Exec(ctx) // nolint
	}
	return
}
