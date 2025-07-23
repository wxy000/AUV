#!/bin/sh
##################################
# Golang一键打包 macOS, Linux, Windows 应用程序
# 使用方法:
#   sh build.sh [-n appname] [-a arch] [-A]
#
# 参数说明:
#   -n appname: 设置应用程序名称，默认为 "AUV"
#   -a arch: 指定架构，默认为所有支持的架构
#   -A: 打包所有架构
#
# 示例:
#   sh build.sh -n AUV -A
#   将在 target 目录下生成以下 9 个可执行文件:
#   AUV-darwin-amd64.bin  AUV-darwin-arm64.bin
#   AUV-linux-amd64.bin  AUV-linux-arm64.bin  AUV-linux-arm.bin  AUV-linux-386.bin
#   AUV-windows-amd64.exe  AUV-windows-arm64.exe  AUV-windows-386.exe
#
# 支持的架构: amd64, arm64, arm (32位), 386
##################################
# 检查是否安装了 Go
if ! command -v go &> /dev/null; then
    echo "Go 编译器未安装，请先安装 Go！"
    exit 1
fi
# 获取当前时间戳
timestamp=$(date +"%Y-%m-%d_%H-%M-%S")
# 获取用户输入参数
while getopts ":n:a:A" opt; do
    case $opt in
    n)
        APPNAME=$OPTARG
        ;;
    a)
        ARCH=$OPTARG
        ;;
    A)
        ALL_ARCHS=true
        ;;
    ?)
        echo "未知参数"
        exit 1
        ;;
    esac
done
# 默认应用程序名称为 "AUV"
APPNAME=${APPNAME:-"AUV"}
# 默认架构为所有架构
ARCH=${ARCH:-"all"}
# 创建 bin 和 log 目录
mkdir -p bin
mkdir -p logs
# 日志文件路径
log_file="logs/build-${timestamp}.log"
# 通用变量
export CGO_ENABLED=0  # 关闭 CGO
# 定义支持的架构列表
ALL_ARCHS_LIST=("amd64" "arm64" "arm" "386")
# 如果没有指定特定架构，则打包所有架构
if [ "$ARCH" = "all" ] || [ -n "$ALL_ARCHS" ]; then
    ARCHS=("${ALL_ARCHS_LIST[@]}")
else
    IFS=',' read -r -a ARCHS <<< "$ARCH"
fi
# 构建函数
build_for_os() {
    local os=$1
    for arch in "${ARCHS[@]}"; do
        export GOOS=$os
        export GOARCH=$arch
        case $arch in
            amd64)
                build_target="${APPNAME}-${os}-amd64.bin"
                ;;
            arm64)
                build_target="${APPNAME}-${os}-arm64.bin"
                ;;
            arm)
                build_target="${APPNAME}-${os}-arm.bin"
                ;;
            386)
                build_target="${APPNAME}-${os}-386.bin"
                ;;
            *)
                echo "不支持的架构: $arch"
                exit 1
                ;;
        esac
        if [ "$os" = "windows" ]; then
            build_target="${build_target%.bin}.exe"
        fi
        # 添加时间戳并记录到日志文件
        echo "$(date +%Y-%m-%d\ %H:%M:%S) - 构建 $os ($arch): $build_target" | tee -a "$log_file"
        go build -ldflags "-s -w" -o bin/"$build_target" 2>&1 | tee -a "$log_file"
        if [ $? -eq 0 ]; then
            echo "$(date +%Y-%m-%d\ %H:%M:%S) - $os ($arch) 可执行程序 $build_target 打包成功!" | tee -a "$log_file"
            echo "$os ($arch) 可执行程序 $build_target 打包成功!"
        else
            echo "$(date +%Y-%m-%d\ %H:%M:%S) - $os ($arch) 可执行程序 $build_target 打包失败!" | tee -a "$log_file"
            echo "$os ($arch) 可执行程序 $build_target 打包失败!"
            exit 1
        fi
    done
}
# 构建不同操作系统版本
build_for_os darwin
build_for_os linux
build_for_os windows
