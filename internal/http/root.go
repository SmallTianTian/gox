package http

import (
	"fmt"
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/utils"
	ast_util "github.com/SmallTianTian/fresh-go/utils/ast"
	config_util "github.com/SmallTianTian/fresh-go/utils/config"
)

var fileAndTmpl = map[string]string{
	"ui/http/root.go": utils.ReadStatikFile("/http/gin/ui_http_root.go.tmpl"),
	"server/http.go":  utils.ReadStatikFile("/http/server_http.go.tmpl"),
}

func NewHttp(path, organization string, httpPort int) {
	module := filepath.Join(organization, filepath.Base(path))
	var kRv = map[string]interface{}{"module": module}
	utils.WriteByTemplate(path, fileAndTmpl, kRv)
	addConfig(path, httpPort)
	addCmdRun(path, module)
}

func addConfig(path string, httpPort int) {
	pg := filepath.Join(path, "config/config.go")
	fga := utils.File2GoAST(pg)
	ast_util.AddField2AstFile(fga, "Http", "int", []string{"Config", "Port"})
	ast_util.AddField2AstFile(fga, "HttpPrefix", "string", []string{"Config", "Application"})
	utils.WriteAstFile(pg, "", fga)

	config_util.WriteConfig(path, "http", httpPort, []string{"port"})
	config_util.WriteConfig(path, "HttpPrefix", "", []string{"application"})
}

func addCmdRun(path, module string) {
	path = filepath.Join(path, "cmd/base.go")
	fga := utils.File2GoAST(path)
	ast_util.AppendFuncCall2AstFile(fga, "server.RunHttp", []string{}, []string{"start"})
	ast_util.AppendFuncCall2AstFile(fga, "server.StopHttp", []string{}, []string{"stop"})
	ast_util.SetImport2AstFile(fga, fmt.Sprintf("%s/server", module))
	utils.WriteAstFile(path, "", fga)
}
