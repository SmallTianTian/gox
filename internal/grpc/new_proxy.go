package grpc

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/SmallTianTian/fresh-go/config"
	"github.com/SmallTianTian/fresh-go/internal/templates"
	"github.com/SmallTianTian/fresh-go/model"
	"github.com/SmallTianTian/fresh-go/utils"
	ast_util "github.com/SmallTianTian/fresh-go/utils/ast"
	config_util "github.com/SmallTianTian/fresh-go/utils/config"
)

// NewDemoProxy 特殊的存在，只应该初始化项目的时候使用
func NewDemoProxy() {
	NewProxy("helloworld", "greeter")
}

func NewProxy(srv, alias string) {
	initProxy()

	if alias == "" {
		alias = srv
	}

	dir := config.DefaultConfig.Project.Path
	proxyPath := filepath.Join(dir, "internal", "server", "proxy.go")
	lines := utils.ReadTxtFileEachLine(proxyPath)
	sb := strings.Builder{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		// 在 import 处导入新包
		if strings.HasPrefix(line, "import (") {
			sb.WriteString(line + "\n")
			org := config.DefaultConfig.Project.Org
			name := config.DefaultConfig.Project.Name
			sb.WriteString(fmt.Sprintf(`%s "%s/%s/api/%s/v1"`, srv, org, name, srv) + "\n")
			continue
		}

		if strings.Contains(line, "return []error{") {
			sb.WriteString(line + "\n")
			sb.WriteString(fmt.Sprintf("%s.Register%sServiceHandlerFromEndpoint(ctx, mux, port, opts),\n", srv, utils.FirstUp(alias))) // nolint
			continue
		}

		sb.WriteString(line + "\n")
	}
	utils.OverwritingFile(proxyPath, "", sb.String())
}

// nolint
func initProxy() {
	dir := config.DefaultConfig.Project.Path
	proxyPath := filepath.Join(dir, "internal", "server", "proxy.go")
	if utils.IsExist(proxyPath) {
		return
	}

	// 初始化 server/proxy.go
	mod := filepath.Join(config.DefaultConfig.Project.Org, config.DefaultConfig.Project.Name)
	base := []*model.FileTemp{
		{Name: "pkg/application/proxy_server.go", Content: templates.ReadTemplateFile("pkg/application/proxy_server.go.tmpl")},
		{Name: "internal/server/proxy.go", Content: templates.ReadTemplateFile("internal/server/proxy.go.tmpl")},
	}
	utils.WriteByTemplate(dir, map[string]interface{}{"module": mod}, base...)

	// 初始化 application 中的数据

	// 写配置文件
	pg := filepath.Join(dir, "internal/conf/config.go")
	fga := utils.File2GoAST(pg)

	config_util.WriteConfig(dir, "proxy", 8089, []string{"port"})
	ast_util.AddField2AstFile(fga, "Proxy", "int", []string{"Config", "Port"})
	utils.WriteAstFile(pg, "", fga)

	// 加入 server/wire.go
	serverWirePath := filepath.Join(dir, "internal", "server", "wire.go")
	lines := utils.ReadTxtFileEachLine(serverWirePath)
	sb := strings.Builder{}

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		tmp := "var ServerSet = wire.NewSet("
		if strings.HasPrefix(line, tmp) {
			sb.WriteString(tmp + "GetProxy, " + line[len(tmp):] + "\n")
			continue
		}

		sb.WriteString(line + "\n")
	}
	utils.OverwritingFile(serverWirePath, "", sb.String())

	// 加入 cmd/server/main.go
	cmdMainPath := filepath.Join(dir, "cmd", "server", "main.go")
	lines = utils.ReadTxtFileEachLine(cmdMainPath)
	sb = strings.Builder{}
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		// 在 import 处导入新包
		if strings.HasPrefix(line, "import (") {
			sb.WriteString(line + "\n")
			sb.WriteString(`"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"` + "\n")
			continue
		}

		// 在 NewApplication 加入对应的参数
		tmp := "func NewApplication("
		if strings.HasPrefix(line, tmp) {
			sb.WriteString(tmp + "proxyMux *runtime.ServeMux, " + line[len(tmp):] + "\n")
			continue
		}

		if strings.Contains(line, "return application.NewApplication") {
			sb.WriteString(line + "\n")
			if line[len(line)-1] != '.' {
				sb.WriteString(line + ".\n")
				sb.WriteString("WithProxy(proxyMux, config.Port.Proxy)\n")
			} else {
				sb.WriteString("WithProxy(proxyMux, config.Port.Proxy).\n")
			}
			continue
		}
		sb.WriteString(line + "\n")
	}
	utils.OverwritingFile(cmdMainPath, "", sb.String())
	utils.GoWireGen(dir)
}
