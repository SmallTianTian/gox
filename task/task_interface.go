// 本文件主要以下事情
// 1. 定义任务接口
// 2. 抽象一个任务共同的 struct，方便其他 struct 继承
package task

import (
	"context"
	"errors"

	"tianxu.xin/gox/internal/util"
)

// GoXTask 任务接口
type GoXTask interface {
	// PreCheck 任务前置检查
	PreCheck(context.Context) bool
	// Exec 任务主体
	Exec(context.Context) error
}

// GoxAbstractTask 对任务的抽象，任务都应该具有这些属性
type GoxAbstractTask struct {
	// Name 任务名称
	Name string

	// PreChecks 前置检查，可能有多个
	PreChecks []func(context.Context) bool
	// ChildsTask 后续任务，可能有多个
	ChildsTask []GoXTask
}

// WithName 设置名字
func (gox *GoxAbstractTask) WithName(name string) {
	if gox != nil {
		gox.Name = name
	}
}

// AppendCheck 添加一个检查方法
func (gox *GoxAbstractTask) AppendCheck(f func(context.Context) bool) {
	if f == nil {
		panic("nil check")
	}
	if gox != nil {
		gox.PreChecks = append(gox.PreChecks, f)
	}
}

// AppendChildTask 添加一个子任务
func (gox *GoxAbstractTask) AppendChildTask(t GoXTask) {
	if t == nil {
		panic("nil task")
	}
	if gox != nil {
		gox.ChildsTask = append(gox.ChildsTask, t)
	}
}

func (gox *GoxAbstractTask) PreCheck(ctx context.Context) bool {
	if gox == nil {
		return false
	}
	log := util.MustExtractLog(ctx)

	for i, check := range gox.PreChecks {
		if !check(ctx) {
			log.Errorf("%d's check failed", i)
			return false
		}
	}
	return true
}

func (gox *GoxAbstractTask) AfterExec(ctx context.Context) error {
	if gox == nil {
		return errors.New("nil task")
	}
	log := util.MustExtractLog(ctx)

	for i, task := range gox.ChildsTask {
		if !task.PreCheck(ctx) {
			log.Errorf("`%s` %d's child task check failed", gox.Name, i)
			return errors.New("check failed")
		}

		if err := task.Exec(ctx); err != nil {
			log.Errorf("`%s` %d's child task exec failed", gox.Name, i)
			return err
		}
	}
	return nil
}
