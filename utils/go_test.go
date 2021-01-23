package utils_test

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/SmallTianTian/fresh-go/test"
	"github.com/SmallTianTian/fresh-go/utils"
	. "github.com/smartystreets/goconvey/convey"
)

func Test_CheckGoProject(t *testing.T) {
	Convey("", t, func() {
		Convey("当前项目", func() {
			exist := utils.CheckGoProject("..")
			So(exist, ShouldBeTrue)
		})

		Convey("正常包含 go.mod 文件", func() {
			pro := "temp_project"
			org := "github.com/fresh-go"
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.InitGoMod(org, pro, dir)

			exist := utils.CheckGoProject(dir)
			So(exist, ShouldBeTrue)
		})

		Convey("临时，不包含 go.mod 的目录", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)

			exist := utils.CheckGoProject(dir)
			So(exist, ShouldBeFalse)
		})
	})
}

func Test_CheckUseVendor(t *testing.T) {
	Convey("", t, func() {
		Convey("当前项目", func() {
			exist := utils.CheckUseVendor("..")
			So(exist, ShouldBeFalse)
		})

		Convey("正常包含 vendor 目录", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.InitVendor(dir, true)

			exist := utils.CheckUseVendor(dir)
			So(exist, ShouldBeTrue)
		})

		Convey("包含 vendor 文件", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.InitVendor(dir, false)

			exist := utils.CheckUseVendor(dir)
			So(exist, ShouldBeFalse)
		})

		Convey("不包含任何 vendor 文件/目录", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)

			exist := utils.CheckUseVendor(dir)
			So(exist, ShouldBeFalse)
		})
	})
}

func Test_GetOrganizationAndProjectName(t *testing.T) {
	Convey("", t, func() {
		Convey("当前项目", func() {
			org, pro := utils.GetOrganizationAndProjectName("..")
			So(org, ShouldEqual, "github.com/SmallTianTian")
			So(pro, ShouldEqual, "fresh-go")
		})
		Convey("新建目录，包含 go.mod", func() {
			e_pro := "temp_project"
			e_org := "github.com/fresh-go"
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.InitGoMod(e_org, e_pro, dir)

			org, pro := utils.GetOrganizationAndProjectName(dir)
			So(org, ShouldEqual, e_org)
			So(pro, ShouldEqual, e_pro)
		})
		Convey("go.mod 文件不包含任何内容，将返回空字符串", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.WriteFile(filepath.Join(dir, "go.mod"), "")

			org, pro := utils.GetOrganizationAndProjectName(dir)
			So(org, ShouldEqual, "")
			So(pro, ShouldEqual, "")
		})
		Convey("不包含 go.mod 文件，将抛出错误", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)

			So(func() { utils.GetOrganizationAndProjectName(dir) }, ShouldPanic)
		})
	})
}

func Test_GoFmtCode(t *testing.T) {
	Convey("", t, func() {
		Convey("正常格式化代码", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			f := filepath.Join(dir, "main.go")
			content := `package main;import "fmt";func main() {fmt.Println("vim-go")}`
			test.WriteFile(f, content)

			old, err := ioutil.ReadFile(f)
			So(err, ShouldBeNil)
			So(string(old), ShouldContainSubstring, ";")

			flag := utils.GoFmtCode(dir)
			So(flag, ShouldBeTrue)

			fresh, e := ioutil.ReadFile(f)
			So(e, ShouldBeNil)
			So(string(fresh), ShouldNotContainSubstring, ";")
			So(len(old) < len(fresh), ShouldBeTrue)
		})

		Convey("格式化代码失败", func() {
			SkipSo("不知道什么情况下会出错。")
		})
	})
}

func Test_GoModRebuild(t *testing.T) {
	Convey("", t, func() {
		Convey("正常重新整理 mod，不包含 vendor", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.InitGoMod("github.com/fresh-go", "test", dir)
			test.WriteFile(filepath.Join(dir, "main.go"), `package main
			import (
				"fmt"
				"github.com/spf13/cast"
			)
			func main() { fmt.Println(cast.ToInt("1")) }
			`)
			mod := filepath.Join(dir, "go.mod")

			old, err := ioutil.ReadFile(mod)
			So(err, ShouldBeNil)
			So(string(old), ShouldNotContainSubstring, "github.com/spf13/cast")

			flag := utils.GoModRebuild(dir)
			So(flag, ShouldBeTrue)

			fresh, e := ioutil.ReadFile(mod)
			So(e, ShouldBeNil)
			So(string(fresh), ShouldContainSubstring, "github.com/spf13/cast")
		})
		Convey("正常重新整理 mod，包含 vendor", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.InitGoMod("github.com/fresh-go", "test", dir)
			test.WriteFile(filepath.Join(dir, "main.go"), `package main
			import (
				"fmt"
				"github.com/spf13/cast"
			)
			func main() { fmt.Println(cast.ToInt("1")) }
			`)
			test.InitVendor(dir, true)
			mod := filepath.Join(dir, "go.mod")
			github := filepath.Join(dir, "vendor/github.com")

			old, err := ioutil.ReadFile(mod)
			So(err, ShouldBeNil)
			So(string(old), ShouldNotContainSubstring, "github.com/spf13/cast")
			So(utils.IsExist(github), ShouldBeFalse)

			flag := utils.GoModRebuild(dir)
			So(flag, ShouldBeTrue)

			fresh, e := ioutil.ReadFile(mod)
			So(e, ShouldBeNil)
			So(string(fresh), ShouldContainSubstring, "github.com/spf13/cast")
			So(utils.IsExist(github), ShouldBeTrue)
		})
		Convey("整理不包含 go.mod 的项目将失败", func() {
			dir := test.TempDir()
			defer os.RemoveAll(dir)

			flag := utils.GoModRebuild(dir)
			So(flag, ShouldBeFalse)
		})
	})
}

func Test_FirstMod(t *testing.T) {
	Convey("", t, func() {
		Convey("首次创建 mod，不包含 vendor", func() {
			pro := "temp_project"
			org := "github.com/fresh-go"

			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.WriteFile(filepath.Join(dir, "main.go"), `package main;import (	"fmt";	"github.com/spf13/cast");func main() { fmt.Println(cast.ToInt("1")) }`)
			mod := filepath.Join(dir, "go.mod")
			vendor := filepath.Join(dir, "vendor")

			So(utils.IsExist(mod), ShouldBeFalse)
			So(utils.IsExist(vendor), ShouldBeFalse)

			flag := utils.FirstMod(dir, filepath.Join(pro, org), false)
			So(flag, ShouldBeTrue)

			So(utils.IsExist(mod), ShouldBeTrue)
			So(utils.IsExist(vendor), ShouldBeFalse)
		})

		Convey("首次创建 mod，包含 vendor", func() {
			pro := "temp_project"
			org := "github.com/fresh-go"

			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.WriteFile(filepath.Join(dir, "main.go"), `package main;import (	"fmt";	"github.com/spf13/cast");func main() { fmt.Println(cast.ToInt("1")) }`)
			mod := filepath.Join(dir, "go.mod")
			vendor := filepath.Join(dir, "vendor")

			So(utils.IsExist(mod), ShouldBeFalse)
			So(utils.IsExist(vendor), ShouldBeFalse)

			flag := utils.FirstMod(dir, filepath.Join(pro, org), true)
			So(flag, ShouldBeTrue)

			So(utils.IsExist(mod), ShouldBeTrue)
			So(utils.IsExist(vendor), ShouldBeTrue)
		})

		Convey("首次创建 mod，包含 vendor，不包含任何 go 文件", func() {
			pro := "temp_project"
			org := "github.com/fresh-go"

			dir := test.TempDir()
			defer os.RemoveAll(dir)
			mod := filepath.Join(dir, "go.mod")
			vendor := filepath.Join(dir, "vendor")

			So(utils.IsExist(mod), ShouldBeFalse)
			So(utils.IsExist(vendor), ShouldBeFalse)

			flag := utils.FirstMod(dir, filepath.Join(pro, org), true)
			So(flag, ShouldBeTrue)

			So(utils.IsExist(mod), ShouldBeTrue)
			// 这里为 false 的原因是，没有依赖，go mod vendor 命令将删除本地 vendor 目录
			So(utils.IsExist(vendor), ShouldBeFalse)
		})

		Convey("非首次创建 mod，将不会做任何操作", func() {
			pro := "temp_project"
			org := "github.com/fresh-go"
			b_pro := "bad_pro"
			b_org := "bad_org"

			dir := test.TempDir()
			defer os.RemoveAll(dir)
			test.InitGoMod(org, pro, dir)

			mod := filepath.Join(dir, "go.mod")

			old, err := ioutil.ReadFile(mod)
			So(err, ShouldBeNil)
			So(string(old), ShouldContainSubstring, pro)
			So(string(old), ShouldContainSubstring, org)

			flag := utils.FirstMod(dir, filepath.Join(pro, org), true)
			So(flag, ShouldBeTrue)

			fresh, e := ioutil.ReadFile(mod)
			So(e, ShouldBeNil)
			So(string(fresh), ShouldContainSubstring, pro)
			So(string(fresh), ShouldContainSubstring, org)
			So(string(fresh), ShouldNotContainSubstring, b_pro)
			So(string(fresh), ShouldNotContainSubstring, b_org)
		})
	})
}
