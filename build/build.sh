#!/bin/bash

# 高级交叉编译脚本
# 支持选择性构建、并行编译、Docker构建等功能

set -e

# 默认配置
VERSION=${VERSION:-"1.0.0"}
OUTPUT_DIR="release"
PARALLEL_JOBS=${PARALLEL_JOBS:-4}
SKIP_FRONTEND=${SKIP_FRONTEND:-false}
DOCKER_BUILD=${DOCKER_BUILD:-false}
PLATFORMS_FILTER=${PLATFORMS_FILTER:-""}

# 颜色输出
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m'

# 打印函数
print_info() { echo -e "${BLUE}[INFO]${NC} $1"; }
print_success() { echo -e "${GREEN}[SUCCESS]${NC} $1"; }
print_warning() { echo -e "${YELLOW}[WARNING]${NC} $1"; }
print_error() { echo -e "${RED}[ERROR]${NC} $1"; }
print_header() { echo -e "${PURPLE}[HEADER]${NC} $1"; }

# 显示帮助信息
show_help() {
    cat << EOF
CronJob 交叉编译脚本

用法: $0 [选项]

选项:
    -v, --version VERSION       设置版本号 (默认: 1.0.0)
    -o, --output DIR           设置输出目录 (默认: release)
    -j, --jobs N               并行编译任务数 (默认: 4)
    -p, --platforms PLATFORMS  指定构建平台 (逗号分隔)
    --skip-frontend            跳过前端构建
    --docker                   使用 Docker 构建
    --clean                    清理构建缓存
    -h, --help                 显示此帮助信息

平台列表:
    windows-amd64, linux-amd64, linux-arm64, darwin-amd64, darwin-arm64

示例:
    $0                                    # 构建所有平台
    $0 -v 2.0.0 -p "linux-amd64,darwin-amd64"  # 构建指定平台
    $0 --skip-frontend --docker          # 跳过前端，使用Docker构建
EOF
}

# 解析命令行参数
parse_args() {
    while [[ $# -gt 0 ]]; do
        case $1 in
            -v|--version)
                VERSION="$2"
                shift 2
                ;;
            -o|--output)
                OUTPUT_DIR="$2"
                shift 2
                ;;
            -j|--jobs)
                PARALLEL_JOBS="$2"
                shift 2
                ;;
            -p|--platforms)
                PLATFORMS_FILTER="$2"
                shift 2
                ;;
            --skip-frontend)
                SKIP_FRONTEND=true
                shift
                ;;
            --docker)
                DOCKER_BUILD=true
                shift
                ;;
            --clean)
                CLEAN_BUILD=true
                shift
                ;;
            -h|--help)
                show_help
                exit 0
                ;;
            *)
                print_error "未知选项: $1"
                show_help
                exit 1
                ;;
        esac
    done
}

# 检查依赖
check_dependencies() {
    local deps=("go")
    
    if [ "$SKIP_FRONTEND" != "true" ]; then
        deps+=("pnpm")
    fi
    
    if [ "$DOCKER_BUILD" = "true" ]; then
        deps+=("docker")
    fi
    
    for dep in "${deps[@]}"; do
        if ! command -v $dep &> /dev/null; then
            print_error "$dep 未找到，请确保已安装"
            exit 1
        fi
    done
}

# 获取平台列表
get_platforms() {
    local all_platforms=(
        "windows-amd64"
        "linux-amd64" 
        "linux-arm64"
        "darwin-amd64"
        "darwin-arm64"
    )
    
    if [ -n "$PLATFORMS_FILTER" ]; then
        IFS=',' read -ra FILTERED_PLATFORMS <<< "$PLATFORMS_FILTER"
        echo "${FILTERED_PLATFORMS[@]}"
    else
        echo "${all_platforms[@]}"
    fi
}

# 构建前端
build_frontend() {
    if [ "$SKIP_FRONTEND" = "true" ]; then
        print_info "跳过前端构建"
        return 0
    fi
    
    print_header "构建前端"
    cd ../web
    
    if [ ! -f "package.json" ]; then
        print_error "未找到 package.json"
        exit 1
    fi
    
    print_info "安装依赖..."
    pnpm install --frozen-lockfile
    
    print_info "构建前端..."
    pnpm run build
    
    if [ ! -d "dist" ]; then
        print_error "前端构建失败"
        exit 1
    fi
    
    print_success "前端构建完成"
    cd ../build
}

# 构建单个平台
build_platform() {
    local platform=$1
    IFS='-' read -r OS ARCH <<< "$platform"
    
    print_info "构建 $platform..."
    
    # 设置环境变量
    export GOOS=$OS
    export GOARCH=$ARCH
    export CGO_ENABLED=0
    
    # 设置二进制文件名
    local binary_name="cronJob"
    if [ "$OS" = "windows" ]; then
        binary_name="cronJob.exe"
    fi
    
    # 创建输出目录
    local platform_dir="$OUTPUT_DIR/$platform"
    mkdir -p "$platform_dir"
    
    # 编译
    local build_time=$(date '+%Y-%m-%d %H:%M:%S')
    local ldflags="-s -w -X cronJob/cmd.AppVersion=$VERSION -X 'cronJob/cmd.BuildTime=$build_time'"
    
    cd ..
    if ! go build -ldflags "$ldflags" -o "$platform_dir/$binary_name" .; then
        print_error "$platform 编译失败"
        cd build
        return 1
    fi
    cd build
    
    # 复制文件
    cp ../config.example.yaml "$platform_dir/config.yaml"
    [ -f "../README.md" ] && cp ../README.md "$platform_dir/"
    [ -f "../LICENSE" ] && cp ../LICENSE "$platform_dir/"
    
    # 创建启动脚本
    if [ "$OS" = "windows" ]; then
        cat > "$platform_dir/start.bat" << EOF
@echo off
echo 启动 CronJob 服务...
$binary_name server
pause
EOF
    else
        cat > "$platform_dir/start.sh" << 'EOF'
#!/bin/bash
echo "启动 CronJob 服务..."
./cronJob server
EOF
        chmod +x "$platform_dir/start.sh"
    fi
    
    # 创建压缩包
    cd "$platform_dir"
    if [ "$OS" = "windows" ]; then
        local archive_name="cronJob-$platform-v$VERSION.zip"
        if command -v zip &> /dev/null; then
            zip -r "../$archive_name" * > /dev/null
        fi
    else
        local archive_name="cronJob-$platform-v$VERSION.tar.gz"
        tar -czf "../$archive_name" *
    fi
    cd ../../build
    
    print_success "$platform 构建完成"
    return 0
}

# 主函数
main() {
    parse_args "$@"
    
    print_header "CronJob 交叉编译"
    print_info "版本: $VERSION"
    print_info "输出目录: $OUTPUT_DIR"
    print_info "并行任务: $PARALLEL_JOBS"
    
    # 清理
    if [ "$CLEAN_BUILD" = "true" ]; then
        print_info "清理构建缓存..."
        go clean -cache
        rm -rf $OUTPUT_DIR
    fi
    
    # 检查依赖
    check_dependencies
    
    # 创建输出目录
    rm -rf $OUTPUT_DIR
    mkdir -p $OUTPUT_DIR
    
    # 构建前端
    build_frontend
    
    # 获取平台列表
    local platforms=($(get_platforms))
    print_info "构建平台: ${platforms[*]}"
    
    # 并行构建
    local pids=()
    local failed_platforms=()
    
    for platform in "${platforms[@]}"; do
        # 控制并发数
        while [ ${#pids[@]} -ge $PARALLEL_JOBS ]; do
            for i in "${!pids[@]}"; do
                if ! kill -0 "${pids[$i]}" 2>/dev/null; then
                    wait "${pids[$i]}"
                    if [ $? -ne 0 ]; then
                        failed_platforms+=("${platform_names[$i]}")
                    fi
                    unset pids[$i]
                    unset platform_names[$i]
                fi
            done
            pids=("${pids[@]}")  # 重新索引数组
            platform_names=("${platform_names[@]}")
            sleep 0.1
        done
        
        # 启动新的构建任务
        build_platform "$platform" &
        pids+=($!)
        platform_names+=("$platform")
    done
    
    # 等待所有任务完成
    for pid in "${pids[@]}"; do
        wait $pid
    done
    
    # 显示结果
    print_header "构建摘要"
    print_info "版本: $VERSION"
    print_info "构建时间: $(date)"
    print_info "输出目录: $OUTPUT_DIR"
    
    if [ ${#failed_platforms[@]} -eq 0 ]; then
        print_success "所有平台构建成功！"
    else
        print_warning "以下平台构建失败: ${failed_platforms[*]}"
    fi
    
    print_info "生成的文件:"
    ls -la $OUTPUT_DIR/ 2>/dev/null || true
}

# 运行主函数
main "$@"
