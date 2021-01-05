package utils

import (
	"log"
	"path/filepath"
	"strings"
)

func FirstMod(path, oAp string, vendor bool) bool {
	GoFmtCode(path)
	// only the first new project should be executed separately.
	if err := Exec(path, "go", "mod", "init", oAp); err != nil {
		log.Printf("Go mod init failed.\n Please exec `go mod init %s`.\n Error: %v", oAp, err)
		return false
	}
	if !GoModRebuild(path) {
		return false
	}
	// only the first new project should be executed separately.
	if vendor {
		if err := Exec(path, "go", "mod", "vendor"); err != nil {
			log.Printf("Go mod vendor failed.\n Please exec `go mod vendor` in `%s`.\n Error: %v", oAp, err)
			return false
		}
	}
	return true
}

func GoModRebuild(path string) bool {
	GoFmtCode(path)
	if err := Exec(path, "go", "mod", "tidy"); err != nil {
		log.Printf("Go mod tidy failed.\n Please exec `go mod tidy` in `%s`.\n Error: %v", path, err)
		return false
	}
	if IsExist(filepath.Join(path, "vendor")) {
		if err := Exec(path, "go", "mod", "vendor"); err != nil {
			log.Printf("Go mod tidy failed.\n Please exec `go mod tidy` in `%s`.\n Error: %v", path, err)
			return false
		}
	}
	return true
}

func GoFmtCode(path string) bool {
	if err := Exec(path, "gofmt", "-s", "-w", "."); err != nil {
		log.Printf("Go fmt code failed.\n Please exec `gofmt -s -w .` in `%s`.\n Error: %v", path, err)
		return false
	}
	return true
}

func GetOrganizationAndProjectName(path string) string {
	MustTrue(CheckGoProject(path), "Not go project.")

	lines := ReadTxtFileEachLine(filepath.Join(path, "go.mod"))
	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "module") {
			return strings.TrimSpace(line)[len("module "):]
		}
	}
	return filepath.Base(path)
}

func CheckGoProject(path string) bool {
	goMod := filepath.Join(path, "go.mod")
	if !IsExist(goMod) {
		log.Printf("`%s` maybe not a go project. Please check `go.mod` is in path.", path)
		return false
	}
	return true
}

func CheckUseVendor(path string) bool {
	return IsExist(filepath.Join(path, "vendor"))
}
