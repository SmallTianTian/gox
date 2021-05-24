package doctor

import (
	"os/exec"
	"regexp"
	"strconv"

	"github.com/SmallTianTian/fresh-go/utils"
)

func CheckingEnv() {
	checks := []struct {
		Name   string
		Desc   string
		IsMust bool
		F      func() (bool, string)
	}{{
		"Go", "Go 是你开始写程序的基础，但我们需要版本大于 1.16。", true, CheckGo,
	}, {
		"Make", "Make 将帮助你更便捷的运行命令。", true, CheckMake,
	}, {
		"Docker", "Docker 能帮助你生成不依赖本地环境的镜像。", false, CheckDocker,
	}, {
		"Grpc", "Grpc 是一个高效的通讯协议，需要依据 proto 文件编译具体代码。", false, CheckGrpcCompile,
	}, {
		"Grpc gateway", "Grpc gatewat 提供了 http 协议和 gRpc 协议转换的能力，同样需要根据 proto 编译。", false, CheckGrpcGatewayCompile,
	}, {
		"Buf", "Buf 是一个简化编译 grpc proto 的工具。", false, CheckBuf,
	}, {
		"Wire", "Wire 是一个依赖注入框架，能便捷的生成代码。", false, CheckWire,
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
		println()
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
	versionR := regexp.MustCompile(`\d+.\d+`)
	v := versionR.FindString(string(bs))
	vf, _ := strconv.ParseFloat(v, 64) // nolint

	return vf >= 1.16, string(bs) // nolint
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
		return false, `Cloudn't compile grpc.
More doc: https://grpc.io/docs/languages/go/quickstart/`
	}

	exist1 := utils.CheckCommandExists("protoc-gen-go")
	exist2 := utils.CheckCommandExists("protoc-gen-go-grpc")

	if exist1 && exist2 {
		return true, string(bs)
	}
	return false, `Cloudn't compile grpc.
More doc: https://grpc.io/docs/languages/go/quickstart/`
}

func CheckGrpcGatewayCompile() (bool, string) {
	cmd := exec.Command("protoc", "--version")
	bs, err := cmd.Output()
	if err != nil {
		return false, `Cloudn't compile grpc gateway.
More doc: https://github.com/grpc-ecosystem/grpc-gateway`
	}

	exist1 := utils.CheckCommandExists("protoc-gen-grpc-gateway")
	exist2 := utils.CheckCommandExists("protoc-gen-openapiv2")
	exist3 := utils.CheckCommandExists("protoc-gen-go")
	exist4 := utils.CheckCommandExists("protoc-gen-go-grpc")
	if exist1 && exist2 && exist3 && exist4 {
		return true, string(bs)
	}
	return false, `Cloudn't compile grpc gateway.
More doc: https://github.com/grpc-ecosystem/grpc-gateway`
}

func CheckBuf() (bool, string) {
	cmd := exec.Command("buf", "--version")
	bs, err := cmd.Output()
	if err != nil {
		return false, `Cloudn't call buf.
More doc: https://docs.buf.build/installation/#from-source`
	}
	return true, string(bs)
}

func CheckWire() (bool, string) {
	exist := utils.CheckCommandExists("wire")
	if exist {
		return true, ""
	}
	return false, `Cloudn't exist command wire.
More doc: https://github.com/google/wire`
}
