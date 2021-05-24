package grpc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/templates"
	"github.com/SmallTianTian/fresh-go/model"
	"github.com/SmallTianTian/fresh-go/pkg/logger"
	"github.com/SmallTianTian/fresh-go/utils"
	ast_util "github.com/SmallTianTian/fresh-go/utils/ast"
	config_util "github.com/SmallTianTian/fresh-go/utils/config"
)

// New 新建一个 GRPC 文件
// TODO
func New(srv string) {
	f := func() {
		dir := config.DefaultConfig.Project.Path
		kRv := map[string]interface{}{
			"module": config_util.GetModule(config.DefaultConfig),
		}
		utils.WriteByTemplate(dir, kRv, optionMap["demo_proto"], optionMap["demo_impl"])
	}
	baseGRPCnew(srv, srv, f)
}

// NewDemo 是一个特殊的存在，只在项目创建的时候才应该被调用
func NewDemo() {
	f := func() {
		dir := config.DefaultConfig.Project.Path
		kRv := map[string]interface{}{
			"module": config_util.GetModule(config.DefaultConfig),
		}
		utils.WriteByTemplate(dir, kRv, optionMap["demo_proto"], optionMap["demo_impl"])
	}
	baseGRPCnew("helloworld", "greeter", f)
}

// 基础创建 grpc 的方法
func baseGRPCnew(name, alias string, f func()) {
	logger.Debugf("Will begin create grpc %s.", name)
	// 0. 检查是否存在该 GRPC
	if existGRPC(name) {
		logger.Errorf("exist %s grpc. will skip create %s grpc.", name, name)
		return
	}

	// 1. 调用自定义写 proto 及实现的方法
	f()

	// 2. 调用 buf 生成代码
	withGRPCBufGen()
	// 3. 设置 grpc 文件中的 wire
	setGRPCImplWire(name, alias)
	// 4. 设置 grpc server
	setGRPCServer()
	// 5. 设置一些帮助信息
	setHelper()
}

var optionMap = map[string]*model.FileTemp{
	"impl_wire":   {Name: "internal/ui/grpc/wire.go", Content: templates.ReadTemplateFile("internal/ui/grpc/wire.go.tmpl")},                   // nolint
	"grpc_server": {Name: "internal/server/grpc.go", Content: templates.ReadTemplateFile("internal/server/grpc.go.tmpl")},                     // nolint
	"server_wire": {Name: "internal/server/wire.go", Content: templates.ReadTemplateFile("internal/server/wire.go.tmpl")},                     // nolint
	"app":         {Name: "pkg/application/grpc_server.go", Content: templates.ReadTemplateFile("pkg/application/grpc_server.go.tmpl")},       // nolint
	"demo_proto":  {Name: "api/helloworld/v1/greeter.proto", Content: templates.ReadTemplateFile("api/helloworld/v1/greeter.proto")},          // nolint
	"demo_impl":   {Name: "internal/ui/grpc/helloworld_v1.go", Content: templates.ReadTemplateFile("internal/ui/grpc/helloworld_v1.go.tmpl")}, // nolint
	"buf":         {Name: "buf.yaml", Content: templates.ReadTemplateFile("buf.yaml.tmpl")},                                                   // nolint
	"buf_gen":     {Name: "buf.gen.yaml", Content: templates.ReadTemplateFile("buf.gen.yaml.tmpl")},                                           // nolint
}

func existGRPC(pro string) bool {
	return utils.IsExist(protoPath(pro))
}

func protoPath(pro string) string {
	return filepath.Join(config.DefaultConfig.Project.Path, "api", pro, "v1")
}

func withGRPCBufGen() {
	dir := config.DefaultConfig.Project.Path
	// 最终调用 buf 生成代码
	// defer utils.MustNotError(utils.Exec(dir, "buf", "generate"))

	// buf gen yaml 写入
	bufGenPath := filepath.Join(dir, "buf.gen.yaml")
	if !utils.IsExist(bufGenPath) {
		utils.WriteByTemplate(dir, nil, optionMap["buf_gen"])
	}

	bufPath := filepath.Join(dir, "buf.yaml")
	if utils.IsExist(bufPath) {
		return
	}
	// buf yaml 写入
	var remote, owner, name string
	name = config.DefaultConfig.Project.Name
	name = strings.ReplaceAll(name, "_", "-")
	name = strings.ReplaceAll(name, ".", "-")
	name = strings.ReplaceAll(name, "/", "-")

	if remote = config.DefaultConfig.Project.Remote; remote == "" {
		remote = "buf.build"
	}
	if owner = config.DefaultConfig.Project.Owner; owner == "" {
		owner = name
		if len(owner) < 4 { // nolint
			owner += strings.Repeat("-", (4 - len(owner))) // nolint
		}
	}

	// 不基于 buf.build 的远端就无法运作，所以先写个模板
	kRv := map[string]interface{}{
		"name": "buf.build/beta/test",
	}
	utils.WriteByTemplate(dir, kRv, optionMap["buf"])
	// 调用 buf mod update 命令加载其他 proto
	utils.MustNotError(utils.Exec(dir, "buf", "beta", "mod", "update"))

	// 这里才是真正的 name
	kRv = map[string]interface{}{
		"name": strings.Join([]string{remote, owner, name}, "/"),
	}
	utils.WriteByTemplate(dir, kRv, optionMap["buf"])
	utils.GoBufGen(dir)
	// 修复 go tidy 最多只拉去 v1.21 版本
	utils.MustNotError(utils.Exec(dir, "go", "get", "-u", "-v", "google.golang.org/grpc"))
}

func setGRPCImplWire(srv, alias string) { // nolint
	if alias == "" {
		alias = srv
	}
	wirePath := filepath.Join(config.DefaultConfig.Project.Path, "internal", "ui", "grpc", "wire.go")
	// 如果不存在 wire 文件，先创建 wire
	if !utils.IsExist(wirePath) {
		utils.WriteByTemplate(config.DefaultConfig.Project.Path, nil, optionMap["impl_wire"])
	}

	// 以下内容 ast 写起来太麻烦，直接操作文件。
	lines := utils.ReadTxtFileEachLine(wirePath)

	sb := strings.Builder{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// 在 import 处导入新包
		if strings.HasPrefix(line, "import (") {
			sb.WriteString(line + "\n")
			sb.WriteString(fmt.Sprintf(`%s "%s/api/%s/v1"`, srv, config_util.GetModule(config.DefaultConfig), srv) + "\n")
			continue
		}

		// 在 GRPCsProvider 中引入新的 Provider
		const gpstr = "var GRPCsProvider = wire.NewSet("
		if strings.HasPrefix(line, gpstr) {
			sb.WriteString(gpstr)
			sb.WriteString(fmt.Sprintf("%sProvider, ", srv))
			sb.WriteString(line[len(gpstr):] + "\n")

			// 继续写新的 Provider
			tmp := `var %sProvider = wire.NewSet(NewV1%sService, wire.Bind(new(%s.%sServiceServer), new(*V1%sService)))`
			varLine := fmt.Sprintf(tmp, srv, utils.FirstUp(srv), srv, utils.FirstUp(alias), utils.FirstUp(srv))
			if len(strings.TrimSpace(varLine)) > 120 { // nolint
				varLine += " //nolint"
			}
			sb.WriteString(varLine + "\n")
			continue
		}

		// 在 RegistGRPCS 中注册
		const rgstr = "func RegistGRPCS("
		if strings.HasPrefix(line, rgstr) {
			sb.WriteString(rgstr)
			sb.WriteString(fmt.Sprintf("%sS %s.%sServiceServer, ", srv, srv, utils.FirstUp(alias)))
			sb.WriteString(line[len(rgstr):] + "\n")

			// 写入下一行
			i++
			sb.WriteString(lines[i] + "\n")

			// 写入注册信息
			tmp := `%s.Register%sServiceServer(sr, %sS)`
			sb.WriteString(fmt.Sprintf(tmp, srv, utils.FirstUp(alias), srv) + "\n")
			continue
		}
		sb.WriteString(lines[i] + "\n")
	}

	utils.OverwritingFile(wirePath, "", sb.String())
	utils.GoWireGen(config.DefaultConfig.Project.Path)
}

func setGRPCServer() {
	gsPath := filepath.Join(config.DefaultConfig.Project.Path, "internal", "server", "grpc.go")
	if utils.IsExist(gsPath) {
		return
	}

	dir := config.DefaultConfig.Project.Path
	// write grpc.go
	utils.WriteByTemplate(dir, nil, optionMap["grpc_server"])

	// write wire.go
	wirePath := filepath.Join(config.DefaultConfig.Project.Path, "internal", "server", "wire.go")
	// 如果不存在 wire 文件，先创建 wire
	if !utils.IsExist(wirePath) {
		utils.WriteByTemplate(config.DefaultConfig.Project.Path, nil, optionMap["server_wire"])
	}

	// 以下内容 ast 写起来太麻烦，直接操作文件。
	lines := utils.ReadTxtFileEachLine(wirePath)

	sb := strings.Builder{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		const ss = `var ServerSet = wire.NewSet(`
		if strings.HasPrefix(line, ss) {
			sb.WriteString(ss)
			sb.WriteString("GetGRPC, ")
			sb.WriteString(line[len(ss):] + "\n")
			continue
		}
		sb.WriteString(line + "\n")
	}
	utils.OverwritingFile(wirePath, "", sb.String())
}

// nolint
func setHelper() {
	dir := config.DefaultConfig.Project.Path
	appPath := filepath.Join(dir, "pkg", "application", "grpc_server.go")
	if utils.IsExist(appPath) {
		return
	}
	// 写 pkg/application/grpc_server.go
	utils.WriteByTemplate(dir, nil, optionMap["app"])

	// 写配置文件
	pg := filepath.Join(dir, "internal/conf/config.go")
	fga := utils.File2GoAST(pg)

	config_util.WriteConfig(dir, "grpc", 50051, []string{"port"})
	ast_util.AddField2AstFile(fga, "GRPC", "int", []string{"Config", "Port"})

	utils.WriteAstFile(pg, "", fga)

	cmdWirePath := filepath.Join(dir, "cmd", "server", "wire.go")
	content := string(utils.ReadFile(cmdWirePath))
	//  如果 cmd 的 wire 中有 GRPCsProvider 就代表已经写过了
	if strings.Contains(content, "GRPCsProvider") {
		return
	}

	// 以下内容 ast 写起来太麻烦，直接操作文件。
	lines := utils.ReadTxtFileEachLine(cmdWirePath)
	sb := strings.Builder{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]

		// 在 import 处导入新包
		if strings.HasPrefix(line, "import (") {
			sb.WriteString(line + "\n")
			sb.WriteString(fmt.Sprintf(`"%s/internal/server"`, config_util.GetModule(config.DefaultConfig)) + "\n")
			sb.WriteString(fmt.Sprintf(`gc "%s/internal/ui/grpc"`, config_util.GetModule(config.DefaultConfig)) + "\n")
			continue
		}

		// 在 build 处加入对应的 set
		if strings.Contains(line, "panic(wire.Build(") {
			sb.WriteString(line + "\n")
			sb.WriteString("gc.GRPCsProvider,\n")
			sb.WriteString("gc.RegistGRPCS,\n")
			sb.WriteString("server.ServerSet,\n")
			continue
		}
		sb.WriteString(line + "\n")
	}
	utils.OverwritingFile(cmdWirePath, "", sb.String())

	// 写入 cmd 中的 main.go
	cmdMainPath := filepath.Join(dir, "cmd", "server", "main.go")
	lines = utils.ReadTxtFileEachLine(cmdMainPath)
	sb = strings.Builder{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		// 在 import 处导入新包
		if strings.HasPrefix(line, "import (") {
			sb.WriteString(line + "\n")
			sb.WriteString(`"google.golang.org/grpc"` + "\n")
			continue
		}

		// 在 NewApplication 加入对应的参数
		tmp := "func NewApplication("
		if strings.HasPrefix(line, tmp) {
			sb.WriteString(tmp + "gc *grpc.Server, \n" + line[len(tmp):] + "\n")
			continue
		}

		if strings.Contains(line, "return application.NewApplication") {
			if line[len(line)-1] != '.' {
				sb.WriteString(line + ".\n")
				sb.WriteString("WithGRPC(gc, config.Port.GRPC)\n")
			} else {
				sb.WriteString("WithGRPC(gc, config.Port.GRPC).\n")
			}
			continue
		}
		sb.WriteString(line + "\n")
	}
	utils.OverwritingFile(cmdMainPath, "", sb.String())
	utils.GoWireGen(dir)
}
