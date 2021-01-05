package grpc

import (
	"fmt"
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/utils"
	ast_util "github.com/SmallTianTian/fresh-go/utils/ast"
	config_util "github.com/SmallTianTian/fresh-go/utils/config"
)

var fileAndTmpl = map[string]string{
	"ui/grpc/demo/proto/demo.proto": utils.ReadStatikFile("/grpc/ui/grpc/demo/proto/demo.proto.tmpl"),
	"ui/grpc/demo/demo.go":          utils.ReadStatikFile("/grpc/ui/grpc/demo/demo.go.tmpl"),
	"server/grpc.go":                utils.ReadStatikFile("/grpc/server/grpc.go.tmpl"),
}

var gatewayFileAndTmpl = map[string]string{
	"server/proxy.go": utils.ReadStatikFile("/grpc/server/proxy.go.tmpl"),
}

func NewGrpc(path, organization string, grpcPort, grpcProxyPort int) {
	module := filepath.Join(organization, filepath.Base(path))
	var kRv = map[string]interface{}{
		"module":  module,
		"gateway": grpcProxyPort > 0,
	}
	utils.WriteByTemplate(path, fileAndTmpl, kRv)

	if grpcProxyPort > 0 {
		utils.WriteByTemplate(path, gatewayFileAndTmpl, kRv)
	}

	addConfig(path, grpcPort, grpcProxyPort)
	addCmdRun(path, module, grpcProxyPort)
	addMakefile(path, grpcProxyPort)
	utils.Exec(path, "make", "proto")
}

func addConfig(path string, grpcPort, grpcProxyPort int) {
	pg := filepath.Join(path, "config/config.go")
	fga := utils.File2GoAST(pg)

	config_util.WriteConfig(path, "grpc", grpcPort, []string{"port"})
	ast_util.AddField2AstFile(fga, "Grpc", "int", []string{"Config", "Port"})

	if grpcProxyPort > 0 {
		config_util.WriteConfig(path, "proxy", grpcProxyPort, []string{"port"})
		ast_util.AddField2AstFile(fga, "Proxy", "int", []string{"Config", "Port"})
	}
	utils.WriteAstFile(pg, "", fga)
}

func addCmdRun(path, module string, grpcProxyPort int) {
	path = filepath.Join(path, "cmd/base.go")
	fga := utils.File2GoAST(path)

	if grpcProxyPort > 0 {
		ast_util.AppendFuncCall2AstFile(fga, "server.RunProxy", []string{}, []string{"start"})
		ast_util.AppendFuncCall2AstFile(fga, "server.StopProxy", []string{}, []string{"stop"})
	}

	ast_util.AppendFuncCall2AstFile(fga, "server.RunGrpc", []string{}, []string{"start"})
	ast_util.AppendFuncCall2AstFile(fga, "server.StopGrpc", []string{}, []string{"stop"})
	ast_util.SetImport2AstFile(fga, fmt.Sprintf("%s/server", module))
	utils.WriteAstFile(path, "", fga)
}

func addMakefile(path string, grpcProxyPort int) {
	path = filepath.Join(path, "Makefile")
	bs := utils.ReadFile(path)
	command := `

proto: $(shell find . -name '*.proto')
	@protoc -I. -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --go_out=. --go-grpc_out=. --go-grpc_opt=paths=source_relative $^`

	if grpcProxyPort > 0 {
		command += `
	@protoc -I. -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --openapiv2_out . --openapiv2_opt=logtostderr=true $^
	@protoc -I. -I$(GOPATH)/src/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis --grpc-gateway_out . --grpc-gateway_opt=logtostderr=true,paths=source_relative $^`
	}
	utils.OverwritingFile(path, "", string(bs)+command)
}
