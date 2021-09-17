package config

import (
	"fmt"
	"os"
)

type LogType string

const (
	ZapLog    LogType = "zap"
	LogrusLog LogType = "logrus"
)

// Config 项目相关配置文件.
type Config struct {
	// 项目相关配置
	Logger struct {
		Level string // [全局配置] 本项目运行日志级别
	}
	GoEnv struct {
		Module string  // [代码生成] module
		Dir    string  // [全局配置] 生成代码所在的目录
		Vendor bool    // [代码生成;全局配置] 启用 vendor
		NoGit  bool    // [全局配置] 不启用 git
		Logger LogType // [全局配置] 使用什么类型的日志，目前支持 zap、logrus
		// 默认是 github.com，用户提供项目名后，可自动生成 module，
		// 例如在默认前缀时，用户项目名称为 test，module 则为 github.com/test
		ModulePre string // [全局配置] module 前缀
	}
}

func GetDefaultConfig() *Config {
	c := &Config{}

	// 默认日志为 error 级别
	c.Logger.Level = "error"

	// 默认代码生成目录为当前目录
	dir, err := os.Getwd()
	mustNotError(err)
	c.GoEnv.Dir = dir

	// 默认使用 zap 作为日志框架
	c.GoEnv.Logger = ZapLog

	// 默认使用 github.com 作为 module 前缀
	c.GoEnv.ModulePre = "github.com"
	return c
}

// 不使用 util 里面的方法，
// 避免循环引用
func mustNotError(err error) {
	if err != nil {
		panic(fmt.Sprintf("Shouldn't get error. %v", err))
	}
}
