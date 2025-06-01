#!/bin/bash

# 设置错误时退出
set -e

# 设置版本号 (可通过环境变量覆盖)
VERSION=${VERSION:-"1.0.0"}

# 设置输出目录
OUTPUT_DIR="release"

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
check_command() {
    if ! command -v $1 &> /dev/null; then
        print_error "$1 未找到，请确保已安装并在 PATH 中"
        exit 1
    fi
}

# 清理旧的构建文件
print_info "清理旧的构建文件..."
rm -rf $OUTPUT_DIR
mkdir -p $OUTPUT_DIR

# 检查必要的工具
print_info "检查构建环境..."
check_command "go"
check_command "pnpm"

# 首先构建前端
print_info "正在构建前端..."
cd ../web

if [ ! -f "package.json" ]; then
    print_error "未找到 package.json 文件"
    exit 1
fi

print_info "安装前端依赖..."
pnpm install

print_info "构建前端项目..."
pnpm run build

# 检查前端构建结果
if [ ! -d "dist" ]; then
    print_error "前端构建完成但未找到 dist 目录"
    exit 1
fi

cd ../build

# 设置需要构建的平台
PLATFORMS=(
    "windows-amd64"
    "linux-amd64"
    "linux-arm64"
    "darwin-amd64"
    "darwin-arm64"
)

print_info "开始交叉编译，目标平台: ${PLATFORMS[*]}"
print_info "版本号: $VERSION"
echo

# 获取构建时间
BUILD_TIME=$(date '+%Y-%m-%d %H:%M:%S')

# 遍历平台进行构建
for platform in "${PLATFORMS[@]}"; do
    # 分割平台字符串
    IFS='-' read -r OS ARCH <<< "$platform"

    echo "========================================"
    print_info "正在为 $OS-$ARCH 构建..."
    echo "========================================"

    # 设置环境变量
    export GOOS=$OS
    export GOARCH=$ARCH
    export CGO_ENABLED=0

    # 设置输出文件名
    if [ "$OS" = "windows" ]; then
        BINARY_NAME="cronJob.exe"
    else
        BINARY_NAME="cronJob"
    fi

    # 创建平台特定的输出目录
    PLATFORM_DIR="$OUTPUT_DIR/$OS-$ARCH"
    mkdir -p $PLATFORM_DIR

    # 编译
    print_info "编译 Go 二进制文件..."
    cd ..

    if ! go build -ldflags "-s -w -X cronJob/cmd.AppVersion=$VERSION -X 'cronJob/cmd.BuildTime=$BUILD_TIME'" -o "$PLATFORM_DIR/$BINARY_NAME" .; then
        print_error "$OS-$ARCH 编译失败"
        cd build
        continue
    fi

    # 验证二进制文件是否生成
    if [ ! -f "$PLATFORM_DIR/$BINARY_NAME" ]; then
        print_error "二进制文件未生成"
        cd build
        continue
    fi

    # 复制配置文件
    print_info "复制配置文件..."
    if ! cp config.example.yaml "$PLATFORM_DIR/config.yaml"; then
        print_warning "配置文件复制失败"
    fi

    # 复制其他必要文件
    [ -f "README.md" ] && cp README.md "$PLATFORM_DIR/"
    [ -f "LICENSE" ] && cp LICENSE "$PLATFORM_DIR/"

    # 创建启动脚本
    if [ "$OS" = "windows" ]; then
        cat > "$PLATFORM_DIR/start.bat" << EOF
@echo off
echo 启动 CronJob 服务...
$BINARY_NAME server
pause
EOF
    else
        cat > "$PLATFORM_DIR/start.sh" << EOF
#!/bin/bash
echo "启动 CronJob 服务..."
./$BINARY_NAME server
EOF
        chmod +x "$PLATFORM_DIR/start.sh"
    fi

    # 创建压缩包
    print_info "创建压缩包..."
    cd $PLATFORM_DIR

    if [ "$OS" = "windows" ]; then
        ARCHIVE_NAME="cronJob-$OS-$ARCH-v$VERSION.zip"
        if command -v zip &> /dev/null; then
            zip -r "../../$OUTPUT_DIR/$ARCHIVE_NAME" * > /dev/null
        else
            print_warning "zip 命令未找到，跳过压缩包创建"
        fi
    else
        ARCHIVE_NAME="cronJob-$OS-$ARCH-v$VERSION.tar.gz"
        tar -czf "../../$OUTPUT_DIR/$ARCHIVE_NAME" *
    fi

    cd ../../build

    if [ -f "$OUTPUT_DIR/$ARCHIVE_NAME" ]; then
        print_success "压缩包创建成功: $ARCHIVE_NAME"
    else
        print_warning "压缩包创建失败"
    fi

    print_success "$OS-$ARCH 构建完成！"
    echo
done

echo "========================================"
print_success "构建摘要"
echo "========================================"
print_info "版本号: $VERSION"
print_info "构建时间: $BUILD_TIME"
print_info "输出目录: $OUTPUT_DIR"
echo

print_info "生成的文件:"
ls -la $OUTPUT_DIR/*.{zip,tar.gz} 2>/dev/null || print_warning "未找到压缩包文件"
echo

print_success "所有平台构建完成！"
print_info "输出文件位于 $OUTPUT_DIR 目录下"
echo

print_info "使用说明:"
echo "1. 解压对应平台的压缩包"
echo "2. 修改 config.yaml 配置文件"
echo "3. 运行启动脚本或直接执行二进制文件"
echo
