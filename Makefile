.PHONY: all build compress build-macos build-windows compress-macos compress-windows clean test check cover run lint help

BUILD_PATH=./build
APPNAME=ohurlshortener

all: check test build compress
build: build-macos build-windows
compress: compress-macos compress-windows

build-macos:
	@echo "开始 macOS 可执行程序编译..."
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_PATH)/$(APPNAME)
	@echo "macOS 可执行程序编译完成..."

build-windows:
	@echo "开始 Windows 可执行程序编译..."
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o $(BUILD_PATH)/$(APPNAME)_windows
	@echo "Windows 可执行程序编译完成..."

compress-macos:
	@echo "macOS 可执行程序压缩开始..."
	mv $(BUILD_PATH)/$(APPNAME) $(BUILD_PATH)/$(APPNAME)_tmp
	upx --best -o $(BUILD_PATH)/$(APPNAME) $(BUILD_PATH)/$(APPNAME)_tmp
	rm -fr $(BUILD_PATH)/$(APPNAME)_tmp
	@echo "macOS 可执行程序压缩完成..."

compress-windows:
	@echo "Windows 可执行程序压缩开始..."
	mv $(BUILD_PATH)/$(APPNAME)_windows $(BUILD_PATH)/$(APPNAME)_windows_tmp
	upx --best -o $(BUILD_PATH)/$(APPNAME)_windows $(BUILD_PATH)/$(APPNAME)_windows_tmp
	rm -fr $(BUILD_PATH)/$(APPNAME)_windows_tmp
	@echo "Windows 可执行程序压缩完成..."

clean:
	@go clean

test:
	@go test

check:
	@go fmt ./
	@go vet ./

cover:
	@go test -coverprofile xx.out
	@go tool cover -html=xx.out

run:
	$(BUILD_PATH)/$(APPNAME)

lint:
	golangci-lint run --enable-all

help:
	@echo "ohurlshortener 构建命令集，用于快速构建服务。"
	@echo ""
	@echo "用法："
	@echo ""
	@echo "	make <command>"
	@echo ""
	@echo "可用命令如下：:"
	@echo "	make all              格式化代码, 编译生成并压缩全部平台的可执行程序"
	@echo "	make build            编译生成全部可执行程序"
	@echo "	make compress         压缩全部可执行程序"
	@echo "	make build-macos      构建 masOS 可执行程序"
	@echo "	make build-windows    构建 Windows 可执行程序"
	@echo "	make compress-macos   压缩 masOS 可执行程序"
	@echo "	make compress-windows 压缩 Windows 可执行程序"
	@echo "	make clean            清理中间目标文件"
	@echo "	make test             执行测试"
	@echo "	make check            格式化 Go 代码"
	@echo "	make cover            测试覆盖率"
	@echo "	make run              运行程序"
	@echo "	make lint             执行代码检查"
	@echo "	make help             查看帮助"


