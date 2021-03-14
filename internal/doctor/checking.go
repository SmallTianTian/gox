package doctor

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/SmallTianTian/fresh-go/utils"
)

func CheckingEnv() {
	checks := []struct {
		Name   string
		Desc   string
		IsMust bool
		F      func() (bool, string)
	}{{
		"Go", "Go is the foundation of your programming.", true, CheckGo,
	}, {
		"Make", "Make helps you run commands quickly.", true, CheckMake,
	}, {
		"Docker", "Docker helps you build image.", false, CheckDocker,
	}, {
		"Grpc", "Grpc is the fast protocol.", false, CheckGrpcCompile,
	}, {
		"Grpc gateway", "Grpc gatewat helps you convert http to grpc.", false, CheckGrpcCompile,
	}}

	for _, v := range checks {
		have, p := v.F()
		switch {
		case have:
			println("✅", "\t", v.Name)
		case v.IsMust:
			println("❌", "\t", v.Name)
		default:
			println("⚠️", "\t", v.Name)
		}

		println("\t", v.Desc)
		println(p)
		println()
	}
}

func CheckGo() (bool, string) {
	cmd := exec.Command("go", "version")
	bs, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	return true, string(bs)
}

func CheckMake() (bool, string) {
	cmd := exec.Command("make", "-v")
	bs, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	return true, string(bs)
}

func CheckDocker() (bool, string) {
	cmd := exec.Command("docker", "-v")
	bs, err := cmd.Output()
	if err != nil {
		return false, ""
	}
	return true, string(bs)
}

func CheckGrpcCompile() (bool, string) {
	cmd := exec.Command("protoc", "--version")
	bs, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	var gp string
	if gp = os.Getenv("GOPATH"); gp == "" {
		gp = os.Getenv("HOME") + "/go"
	}

	exist1 := utils.IsExist(filepath.Join(gp, "bin", "protoc-gen-go"))
	exist2 := utils.IsExist(filepath.Join(gp, "bin", "protoc-gen-go-grpc"))
	if exist1 && exist2 {
		return true, string(bs)
	}
	return false, `Cloudn't compile grpc. Should exec:
	go get google.golang.org/protobuf/cmd/protoc-gen-go \
	google.golang.org/grpc/cmd/protoc-gen-go-grpc

	More doc: https://grpc.io/docs/languages/go/quickstart/`
}

func CheckGrpcGatewayCompile() (bool, string) {
	cmd := exec.Command("protoc", "--version")
	bs, err := cmd.Output()
	if err != nil {
		return false, ""
	}

	var gp string
	if gp = os.Getenv("GOPATH"); gp == "" {
		gp = os.Getenv("HOME") + "/go"
	}

	exist1 := utils.IsExist(filepath.Join(gp, "bin", "protoc-gen-grpc-gateway"))
	exist2 := utils.IsExist(filepath.Join(gp, "bin", "protoc-gen-openapiv2"))
	exist3 := utils.IsExist(filepath.Join(gp, "bin", "protoc-gen-go"))
	exist4 := utils.IsExist(filepath.Join(gp, "bin", "protoc-gen-go-grpc"))
	if exist1 && exist2 && exist3 && exist4 {
		return true, string(bs)
	}
	return false, `Cloudn't compile grpc. Should exec:
	go install \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    google.golang.org/protobuf/cmd/protoc-gen-go \
    google.golang.org/grpc/cmd/protoc-gen-go-grpc

	More doc: https://github.com/grpc-ecosystem/grpc-gateway`
}
