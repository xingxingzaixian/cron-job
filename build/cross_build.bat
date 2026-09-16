@echo off
setlocal enabledelayedexpansion

rem ============================================================
rem CronJob 交叉编译脚本 (Windows)
rem 用法: build\cross_build.bat [full^|quick]
rem   full : 构建前端 + 交叉编译后端（默认）
rem   quick: 跳过前端构建，仅交叉编译后端（复用已有 web\dist）
rem 可在任意目录执行，也可直接双击运行
rem 可用环境变量覆盖版本号: set VERSION=2.0.0 后再执行
rem ============================================================

rem 以脚本所在目录的上一级作为项目根目录，避免相对路径错位
pushd "%~dp0.."
set ROOT=%CD%
popd

rem 设置版本号 (可通过环境变量覆盖)
if "%VERSION%"=="" set VERSION=1.0.0

rem 设置输出目录（项目根目录下的 release）
set OUTPUT_DIR=%ROOT%\release

rem 解析构建模式参数（默认 full）
set BUILD_MODE=full
if /i "%~1"=="quick" set BUILD_MODE=quick
if /i "%~1"=="full" set BUILD_MODE=full
if not "%~1"=="" if /i not "%~1"=="quick" if /i not "%~1"=="full" goto :usage_error

if /i "%BUILD_MODE%"=="quick" (
    echo 构建模式: quick（跳过前端构建，复用已有前端产物）
) else (
    echo 构建模式: full（构建前端 + 交叉编译后端）
)

rem quick 模式先校验前端产物，避免无效运行清空上一次的 release 目录
if /i "%BUILD_MODE%"=="quick" call :check_frontend_dist
if errorlevel 1 exit /b 1

rem 清理旧的构建文件
echo 清理旧的构建文件...
if exist "%OUTPUT_DIR%" rmdir /s /q "%OUTPUT_DIR%"
mkdir "%OUTPUT_DIR%"

rem 检查必要的工具
echo 检查构建环境...
where go >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 Go 编译器，请确保 Go 已安装并在 PATH 中
    exit /b 1
)

rem full 模式构建前端（前端产物由 web/static.go 的 //go:embed dist 打进二进制）
if /i "%BUILD_MODE%"=="full" call :build_frontend
if errorlevel 1 exit /b 1

rem 设置需要构建的平台
set PLATFORMS=windows-amd64 linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

echo 开始交叉编译，目标平台: %PLATFORMS%
echo 版本号: %VERSION%
echo.

set FAILED=0

rem 遍历平台进行构建
for %%p in (%PLATFORMS%) do call :build_platform %%p

echo ========================================
echo 构建摘要
echo ========================================
echo 版本号: %VERSION%
echo 构建模式: %BUILD_MODE%
echo 构建时间: %date% %time%
echo 输出目录: %OUTPUT_DIR%
echo.
echo 生成的文件:
dir /b "%OUTPUT_DIR%\*.zip" 2>nul
dir /b "%OUTPUT_DIR%\*.tar.gz" 2>nul
echo.

if "%FAILED%"=="0" (
    echo 所有平台构建完成！
) else (
    echo 构建失败的平台数量: %FAILED%
)
echo.
echo 使用说明:
echo 1. 解压对应平台的压缩包
echo 2. 修改 config.yaml 配置文件
echo 3. 运行启动脚本或直接执行二进制文件
echo.
echo Linux 平台解压后需先执行: chmod +x cronJob start.sh
echo.

set RC=0
if not "%FAILED%"=="0" set RC=1
endlocal & exit /b %RC%


rem ========================================
rem 构建单个平台，参数为 os-arch
rem ========================================
:build_platform
for /f "tokens=1,2 delims=-" %%a in ("%~1") do (
    set OS=%%a
    set ARCH=%%b
)

echo ========================================
echo 正在为 !OS!-!ARCH! 构建...
echo ========================================

rem 设置环境变量
set GOOS=!OS!
set GOARCH=!ARCH!
set CGO_ENABLED=0

rem 设置输出文件名
set BINARY_NAME=cronJob
if "!OS!"=="windows" set BINARY_NAME=cronJob.exe

rem 创建平台特定的输出目录
set PLATFORM_DIR=%OUTPUT_DIR%\!OS!-!ARCH!
if not exist "!PLATFORM_DIR!" mkdir "!PLATFORM_DIR!"

rem 编译（-ldflags 的值不能含空格，否则会被 go 按空格切分）
echo 编译 Go 二进制文件...
pushd "%ROOT%"
go build -ldflags "-s -w -X cronJob/cmd.AppVersion=%VERSION%" -o "!PLATFORM_DIR!\!BINARY_NAME!" .
set BUILD_ERR=!errorlevel!
popd

if not "!BUILD_ERR!"=="0" (
    echo 错误: !OS!-!ARCH! 编译失败
    set /a FAILED+=1
    exit /b 1
)

rem 验证二进制文件是否生成
if not exist "!PLATFORM_DIR!\!BINARY_NAME!" (
    echo 错误: 二进制文件未生成
    set /a FAILED+=1
    exit /b 1
)

rem 复制配置文件
echo 复制配置文件...
copy /y "%ROOT%\config.example.yaml" "!PLATFORM_DIR!\config.yaml" >nul
if errorlevel 1 echo 警告: 配置文件复制失败

rem 复制其他必要文件
if exist "%ROOT%\README.md" copy /y "%ROOT%\README.md" "!PLATFORM_DIR!\README.md" >nul
if exist "%ROOT%\LICENSE" copy /y "%ROOT%\LICENSE" "!PLATFORM_DIR!\LICENSE" >nul

rem 创建启动脚本
call :write_start_scripts "!PLATFORM_DIR!" "!BINARY_NAME!" "!OS!"

rem 创建压缩包
echo 创建压缩包...
if "!OS!"=="windows" (
    powershell -NoProfile -Command "Compress-Archive -Path '!PLATFORM_DIR!\*' -DestinationPath '%OUTPUT_DIR%\cronJob-!OS!-!ARCH!-v%VERSION%.zip' -Force"
) else (
    tar -czf "%OUTPUT_DIR%\cronJob-!OS!-!ARCH!-v%VERSION%.tar.gz" -C "!PLATFORM_DIR!" !BINARY_NAME! config.yaml README.md LICENSE start.sh
)

if errorlevel 1 (
    echo 警告: 压缩包创建失败
) else (
    echo 压缩包创建成功
)

echo !OS!-!ARCH! 构建完成！
echo.
exit /b 0


rem ========================================
rem 生成启动脚本，参数: 输出目录 二进制名 目标系统
rem 需临时关闭延迟展开，否则 #!/bin/bash 中的 ! 会被吞掉变成 #/bin/bash
rem ========================================
:write_start_scripts
setlocal disabledelayedexpansion
if "%~3"=="windows" (
    >"%~1\start.bat" echo @echo off
    >>"%~1\start.bat" echo echo 启动 CronJob 服务...
    >>"%~1\start.bat" echo %~2 server
    >>"%~1\start.bat" echo pause
) else (
    >"%~1\start.sh" echo #!/bin/bash
    >>"%~1\start.sh" echo echo "Starting CronJob server..."
    >>"%~1\start.sh" echo ./%~2 server
)
endlocal
exit /b 0


rem ========================================
rem 构建前端（full 模式），产物由 web/static.go 的 //go:embed dist 打进二进制
rem ========================================
:build_frontend
where pnpm >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 pnpm，请确保 pnpm 已安装并在 PATH 中
    exit /b 1
)

echo 正在构建前端...
pushd "%ROOT%\web"

if not exist package.json (
    echo 错误: 未找到 package.json 文件
    popd
    exit /b 1
)

echo 安装前端依赖...
call pnpm install
if errorlevel 1 (
    echo 错误: 前端依赖安装失败
    popd
    exit /b 1
)

echo 构建前端项目...
call pnpm run build
if errorlevel 1 (
    echo 错误: 前端构建失败
    popd
    exit /b 1
)

rem 检查前端构建结果
if not exist dist (
    echo 错误: 前端构建完成但未找到 dist 目录
    popd
    exit /b 1
)

popd
exit /b 0


rem ========================================
rem quick 模式：校验已有前端产物是否存在
rem 缺少 web\dist 时 go:embed dist 会编译失败，所以提前拦截
rem ========================================
:check_frontend_dist
if exist "%ROOT%\web\dist\index.html" exit /b 0
echo 错误: quick 模式未找到前端产物 %ROOT%\web\dist
echo 请先执行一次 full 模式: cross_build.bat full
echo 或手动构建前端: cd web ^&^& pnpm run build
exit /b 1


rem ========================================
rem 参数错误提示
rem ========================================
:usage_error
echo 错误: 未知参数 "%~1"
echo 用法: cross_build.bat [full^|quick]
echo   full : 构建前端 + 交叉编译后端（默认）
echo   quick: 跳过前端构建，仅交叉编译后端（复用已有 web\dist）
exit /b 1
