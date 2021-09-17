// 本文件主要以下事情
// 1. 定义任务接口
// 2. 抽象一个任务共同的 struct，方便其他 struct 继承
package task

import (
	"context"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGoxAbstractTask_WithName(t *testing.T) {
	const newName = "new name"
	const oldName = "old name"
	Convey("测试设置新名称", t, func() {
		Convey("新申明指针状态的 task 无法设置名称，且不会 panic", func() {
			var null *GoxAbstractTask

			So(func() { null.WithName(newName) }, ShouldNotPanic)
			So(null, ShouldBeNil)
		})
		Convey("非指针 task 能重置 task 名称", func() {
			var empty GoxAbstractTask

			So(func() { empty.WithName(newName) }, ShouldNotPanic)
			So(empty.Name, ShouldEqual, newName)
		})
		Convey("能多次重置 task 名称", func() {
			var empty GoxAbstractTask

			So(func() { empty.WithName(oldName) }, ShouldNotPanic)
			So(empty.Name, ShouldEqual, oldName)
			So(func() { empty.WithName(newName) }, ShouldNotPanic)
			So(empty.Name, ShouldEqual, newName)
		})
	})
}

func TestGoxAbstractTask_AppendCheck(t *testing.T) {
	constCheck := func(context.Context) bool { return false }

	Convey("测试新增检查", t, func() {
		Convey("新申明指针状态的 task 无法增加检查，且不会 panic", func() {
			var null *GoxAbstractTask

			So(func() { null.AppendCheck(constCheck) }, ShouldNotPanic)
			So(null, ShouldBeNil)
		})
		Convey("非指针 task 能新增检查", func() {
			var empty GoxAbstractTask

			So(len(empty.PreChecks), ShouldEqual, 0)
			So(func() { empty.AppendCheck(constCheck) }, ShouldNotPanic)
			So(len(empty.PreChecks), ShouldEqual, 1)
			So(empty.PreChecks[0], ShouldEqual, constCheck)
		})
		Convey("不能增加 nil 状态的检查", func() {
			var empty GoxAbstractTask

			So(len(empty.PreChecks), ShouldEqual, 0)
			So(func() { empty.AppendCheck(nil) }, ShouldPanic)
		})
	})
}

func TestGoxAbstractTask_AppendChildTask(t *testing.T) {
	type constTask struct {
		GoXTask
	}

	Convey("测试新增检查", t, func() {
		Convey("新申明指针状态的 task 无法增加检查，且不会 panic", func() {
			var null *GoxAbstractTask

			So(func() { null.AppendChildTask(&constTask{}) }, ShouldNotPanic)
			So(null, ShouldBeNil)
		})
		Convey("非指针 task 能新增检查", func() {
			var empty GoxAbstractTask

			So(len(empty.ChildsTask), ShouldEqual, 0)
			So(func() { empty.AppendChildTask(&constTask{}) }, ShouldNotPanic)
			So(len(empty.ChildsTask), ShouldEqual, 1)
			So(empty.ChildsTask[0], ShouldHaveSameTypeAs, &constTask{})
		})
		Convey("不能增加 nil 状态的任务", func() {
			var empty GoxAbstractTask

			So(len(empty.ChildsTask), ShouldEqual, 0)
			So(func() { empty.AppendChildTask(nil) }, ShouldPanic)
		})
	})
}
