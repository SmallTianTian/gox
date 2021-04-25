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
		New     bool   // 是否是新项目
		EnvPath string // 环境变量中的项目路径
		Path    string // 项目本地路径
		EnvOrg  string // 环境变量中的组织
		Org     string // 项目组织
		Name    string // 项目名称
		Vendor  bool   // 是否启用 vendor
	}
	// HTTP 相关配置
	HTTP struct {
		Port int // HTTP 监听端口
	}
	GRPC struct {
		Port  int // GRPC 监听端口
		Proxy int // GRPC 代理监听端口
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

	DefaultConfig.Project.Path = dir

	// 如果没有从配置文件中读到路径，则用当前地址填充
	if DefaultConfig.Project.EnvPath == "" {
		DefaultConfig.Project.EnvPath = DefaultConfig.Project.Path
	}
	// 当不是 Go 项目，将直接返回
	if !utils.CheckGoProject(dir) {
		if DefaultConfig.Project.EnvOrg == "" {
			DefaultConfig.Project.EnvOrg = "github.com"
		}
		return
	}

	// Go 目录将获取 go.mod 中组织和项目名。
	org, pro := utils.GetOrganizationAndProjectName(dir)
	DefaultConfig.Project.Name = pro
	DefaultConfig.Project.Org = org
	DefaultConfig.Project.New = false
	DefaultConfig.Project.Vendor = utils.CheckUseVendor(dir)

	if DefaultConfig.Project.EnvOrg == "" && DefaultConfig.Project.Org != "" {
		DefaultConfig.Project.EnvOrg = DefaultConfig.Project.Org
	}
}
