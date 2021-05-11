version = v0.0.0
go_flags = -ldflags "-w -s -X 'main.Version=$(version)' -X 'main.BuildTime=`date "+%Y-%m-%d %H:%M:%S"`' -X 'main.GoVersion=`go version`' -X 'main.GitRevision=`git rev-parse HEAD`'"
ARMS = 5 6 7
MIPSS = mips mipsle

.PHONY: i b c install build compile clean

# short comment
i: install
b: build
c: clean

install:
	@go install ${go_flags}

build:
	@go build ${go_flags} -o ./bin/fresh-go

compile:
# mac
	@GOOS=darwin GOARCH=amd64 go build $(go_flags) -o ./bin/fresh-go_darwin_amd64_${version}
# windows
	@GOOS=windows GOARCH=amd64 go build $(go_flags) -o ./bin/fresh-go_windows_amd64_${version}.exe
	@GOOS=windows GOARCH=386 go build $(go_flags) -o ./bin/fresh-go_windows_386_${version}.exe
# freebsd
	@GOOS=freebsd GOARCH=amd64 go build $(go_flags) -o ./bin/fresh-go_freebsd_amd64_${version}
	@GOOS=freebsd GOARCH=386 go build $(go_flags) -o ./bin/fresh-go_freebsd_386_${version}
# linux
	@GOOS=linux GOARCH=amd64 go build $(go_flags) -o ./bin/fresh-go_linux_amd64_${version}
	@GOOS=linux GOARCH=386 go build $(go_flags) -o ./bin/fresh-go_linux_386_${version}
	@for v in ${ARMS}; do \
		CGO_ENABLED=0 GOOS=linux GOARCH=arm GOARM=$$v go build $(go_flags) -o ./bin/fresh-go_linux_arm$${v}_${version}; \
	done;
	@CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $(go_flags) -o ./bin/fresh-go_linux_arm64_${version}
	@for v in ${MIPSS}; do \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$v go build $(go_flags) -o ./bin/fresh-go_linux_$${v}_${version}; \
		CGO_ENABLED=0 GOOS=linux GOARCH=$$v GOMIPS=softfloat go build $(go_flags) -o ./bin/fresh-go_linux_$${v}_sf_${version}; \
	done

clean:
	@rm -rf ./bin