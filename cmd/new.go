package cmd

import (
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"tianxu.xin/gox/internal/config"
	"tianxu.xin/gox/internal/util"
	"tianxu.xin/gox/task"
	"tianxu.xin/gox/task/env"
	"tianxu.xin/gox/task/env/log"
	"tianxu.xin/gox/task/env/mod"
	"tianxu.xin/gox/task/git"
)

var (
	ProjectPath string
	Module      string
	ModulePre   string
	Vendor      bool
	NoGit       bool
	Log         string
	OpenDebug   bool
)

var newCmd = &cobra.Command{
	Use:   "new [project name]",
	Short: "创建一个新的项目。",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		conf := config.GetDefaultConfig()
		setDebug(conf)

		// 设置项目地址
		var pp string
		if pp = ProjectPath; pp == "" {
			pp = conf.GoEnv.Dir
		}
		pp = filepath.Join(pp, args[0])
		t := task.New(conf, pp)

		addNewTask(conf, t, args[0])

		t.AppendChildTask(mod.NewGoModInitTask())
		// 按需添加 vendor 任务
		if Vendor || conf.GoEnv.Vendor {
			t.AppendChildTask(mod.NewGoModVendorTask())
		}

		err := t.Exec()
		util.MustNotError(err)

		util.GoFmtCode(pp)
	},
}

func init() {
	newRegist(newCmd.Flags())
}

func newRegist(fs *pflag.FlagSet) {
	fs.StringVarP(&ProjectPath, "path", "p", "", "设置该项目创建的路径，归功于 go.mod，我们可以在任何地方创建项目。")
	fs.StringVar(&Module, "module", "", "项目的 module，将在初始化 go.mod 的时候将其写入。")
	fs.StringVar(&ModulePre, "module_pre", "", "module 的前缀，拼接项目名称作为 module。如果已设置 module，本参数将失效。")
	fs.BoolVar(&Vendor, "vendor", false, "使用 vendor，不推荐。")
	fs.BoolVar(&NoGit, "no_git", false, "不初始化 git，不推荐。")
	fs.StringVar(&Log, "log", "", "选择日志框架，默认 zap。目前支持：zap、logrus")
}

func addNewTask(conf *config.Config, t *task.GoXRootTask, projectName string) {
	var module string
	var logT config.LogType

	// 设置 config 中的 module
	if module = Module; module == "" {
		var pre string
		if pre = ModulePre; pre == "" {
			pre = conf.GoEnv.ModulePre
		}
		module = filepath.Join(pre, projectName)
	}
	conf.GoEnv.Module = module

	// 添加新建新项目的任务
	t.AppendChildTask(env.NewFreshProjectTask())

	// 添加日志任务
	switch strings.ToLower(Log) {
	case string(config.LogrusLog):
		logT = config.LogrusLog
	case string(config.ZapLog):
		logT = config.ZapLog
	}
	if logT == "" {
		logT = conf.GoEnv.Logger
	}
	switch logT {
	case config.LogrusLog:
		t.AppendChildTask(log.NewLogrusTask())
	case config.ZapLog:
		t.AppendChildTask(log.NewZapTask())
	}

	// 按需添加 git 任务
	if !(NoGit || conf.GoEnv.NoGit) {
		t.AppendChildTask(git.NewInitTask())
	}
}
