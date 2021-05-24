package config

import (
	"errors"
	"os"
	"strings"

	"github.com/SmallTianTian/fresh-go/utils"
	"github.com/spf13/viper"
)

var (
	// DefaultConfig 默认配置文件.
	DefaultConfig *Config
)

// 项目相关配置文件.
type Config struct {
	// 日志相关配置
	Logger struct {
		Level string // 日志级别
	}
	// 项目相关配置
	Project struct {
		New    bool   // 是否是新项目
		Path   string // 项目本地路径
		Remote string // 远端地址
		Owner  string // 项目所有者，可以是 group/user
		Name   string // 项目名称
		Vendor bool   // 是否启用 vendor
	}
	// 框架支持的能力
	FrameEnable struct {
		HTTP  bool
		GRPC  bool
		Proxy bool
	}
	// 适配器，可以替换的组件变量应该在这里标识
	Adapter struct {
		Logger string // 使用哪个日志系统，目前支持 zap、logrus
	}
}

// 初始化配置文件.
func init() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(os.Getenv("FRESH_GO_PATH"))
	viper.AddConfigPath("~/.fresh_go")

	err := viper.ReadInConfig()
	if !errors.As(err, &viper.ConfigFileNotFoundError{}) {
		utils.MustNotError(err)
	}

	// bind config value from env
	viper.SetEnvPrefix("fresh_go")
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	viper.AutomaticEnv()

	err = viper.Unmarshal(&DefaultConfig)
	utils.MustNotError(err)

	// 设置默认日志级别
	if DefaultConfig.Logger.Level == "" {
		DefaultConfig.Logger.Level = "error"
	}

	complementProjectInfo()
}

// 默认补全当前环境信息.
func complementProjectInfo() {
	dir, err := os.Getwd()
	utils.MustNotError(err)

	// 如果没有从配置文件中读到路径，则用当前地址填充
	if DefaultConfig.Project.Path == "" {
		DefaultConfig.Project.Path = dir
	}
	// 当不是 Go 项目，将直接返回
	if !utils.CheckGoProject(dir) {
		return
	}

	// Go 目录将获取 go.mod 中远端、组织和项目名。
	remote, owner, name := utils.GetRemoteOwnerAndProjectName(dir)
	DefaultConfig.Project.Name = name
	DefaultConfig.Project.Remote = remote
	DefaultConfig.Project.Owner = owner
	DefaultConfig.Project.New = false
	DefaultConfig.Project.Vendor = utils.CheckUseVendor(dir)

	// 当前是 go 项目，则重新设置 path 为当前路径
	DefaultConfig.Project.Path = dir
}
