package cmd

import (
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "将自动安装所有 golang 的依赖，例如：wire、grpc protoc、buf 等。",
	Run: func(cmd *cobra.Command, args []string) {
		getInit()
	},
}

func getInit() {
	inits := []struct {
		name string
		pgs  []string
	}{
		{
			"grpc", []string{
				"google.golang.org/protobuf/cmd/protoc-gen-go",
				"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
			},
		},
		{
			"grpc gateway", []string{
				"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway",
				"github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2",
				"google.golang.org/protobuf/cmd/protoc-gen-go",
				"google.golang.org/grpc/cmd/protoc-gen-go-grpc",
			},
		},
		{
			"buf", []string{
				"github.com/bufbuild/buf/cmd/buf",
				"github.com/bufbuild/buf/cmd/protoc-gen-buf-breaking",
				"github.com/bufbuild/buf/cmd/protoc-gen-buf-lint",
			},
		},
		{
			"wire", []string{
				"github.com/google/wire/cmd/wire",
			},
		},
	}

	for _, v := range inits {
		println("# " + v.name)
		const pre = "GO111MODULE=on go get -u -v "
		for _, pg := range v.pgs {
			needRun := pre + pg
			println(needRun)

			cmd := exec.Command("go", "get", "-u", "-v", pg)
			cmd.Dir = os.Getenv("HOME")
			cmd.Env = append(os.Environ(), "GO111MODULE=on")
			if err := cmd.Run(); err != nil {
				println("ERR:", err.Error())
			}
		}
	}
}
