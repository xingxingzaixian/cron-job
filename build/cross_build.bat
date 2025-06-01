@echo off
setlocal enabledelayedexpansion

:: 设置版本号 (可通过环境变量覆盖)
if "%VERSION%"=="" set VERSION=1.0.0

:: 设置输出目录
set OUTPUT_DIR=release

:: 清理旧的构建文件
echo 清理旧的构建文件...
if exist %OUTPUT_DIR% rmdir /s /q %OUTPUT_DIR%
mkdir %OUTPUT_DIR%

:: 检查必要的工具
echo 检查构建环境...
where go >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 Go 编译器，请确保 Go 已安装并在 PATH 中
    exit /b 1
)

where pnpm >nul 2>&1
if errorlevel 1 (
    echo 错误: 未找到 pnpm，请确保 pnpm 已安装并在 PATH 中
    exit /b 1
)

:: 首先构建前端
echo 正在构建前端...
cd ..\web
if not exist package.json (
    echo 错误: 未找到 package.json 文件
    exit /b 1
)

echo 安装前端依赖...
call pnpm install
if errorlevel 1 (
    echo 错误: 前端依赖安装失败
    exit /b 1
)

echo 构建前端项目...
call pnpm run build
if errorlevel 1 (
    echo 错误: 前端构建失败
    exit /b 1
)

:: 检查前端构建结果
if not exist dist (
    echo 错误: 前端构建完成但未找到 dist 目录
    exit /b 1
)

cd ..\build

:: 设置需要构建的平台
set PLATFORMS=windows-amd64 linux-amd64 linux-arm64 darwin-amd64 darwin-arm64

echo 开始交叉编译，目标平台: %PLATFORMS%
echo 版本号: %VERSION%
echo.

:: 遍历平台进行构建
for %%p in (%PLATFORMS%) do (
    for /f "tokens=1,2 delims=-" %%a in ("%%p") do (
        set OS=%%a
        set ARCH=%%b

        echo ========================================
        echo 正在为 !OS!-!ARCH! 构建...
        echo ========================================

        :: 设置环境变量
        set GOOS=!OS!
        set GOARCH=!ARCH!
        set CGO_ENABLED=0

        :: 设置输出文件名
        if "!OS!"=="windows" (
            set BINARY_NAME=cronJob.exe
        ) else (
            set BINARY_NAME=cronJob
        )

        :: 创建平台特定的输出目录
        set PLATFORM_DIR=%OUTPUT_DIR%\!OS!-!ARCH!
        if not exist !PLATFORM_DIR! mkdir !PLATFORM_DIR!

        :: 编译
        echo 编译 Go 二进制文件...
        cd ..
        go build -ldflags "-s -w -X cronJob/cmd.AppVersion=%VERSION% -X cronJob/cmd.BuildTime=%date% %time%" -o !PLATFORM_DIR!\!BINARY_NAME! .
        if errorlevel 1 (
            echo 错误: !OS!-!ARCH! 编译失败
            cd build
            continue
        )

        :: 验证二进制文件是否生成
        if not exist !PLATFORM_DIR!\!BINARY_NAME! (
            echo 错误: 二进制文件未生成
            cd build
            continue
        )

        :: 复制配置文件
        echo 复制配置文件...
        copy config.example.yaml !PLATFORM_DIR!\config.yaml >nul
        if errorlevel 1 (
            echo 警告: 配置文件复制失败
        )

        :: 复制其他必要文件
        if exist README.md copy README.md !PLATFORM_DIR!\ >nul
        if exist LICENSE copy LICENSE !PLATFORM_DIR!\ >nul

        :: 创建启动脚本
        if "!OS!"=="windows" (
            echo @echo off > !PLATFORM_DIR!\start.bat
            echo echo 启动 CronJob 服务... >> !PLATFORM_DIR!\start.bat
            echo !BINARY_NAME! server >> !PLATFORM_DIR!\start.bat
            echo pause >> !PLATFORM_DIR!\start.bat
        ) else (
            echo #!/bin/bash > !PLATFORM_DIR!\start.sh
            echo echo "启动 CronJob 服务..." >> !PLATFORM_DIR!\start.sh
            echo ./!BINARY_NAME! server >> !PLATFORM_DIR!\start.sh
            :: 在 Windows 上无法设置 Unix 权限，但在目标系统上需要手动设置
        )

        :: 创建压缩包
        echo 创建压缩包...
        cd build
        if "!OS!"=="windows" (
            powershell -Command "Compress-Archive -Path '..\!PLATFORM_DIR!\*' -DestinationPath '!OUTPUT_DIR!\cronJob-!OS!-!ARCH!-v%VERSION%.zip' -Force"
        ) else (
            powershell -Command "Compress-Archive -Path '..\!PLATFORM_DIR!\*' -DestinationPath '!OUTPUT_DIR!\cronJob-!OS!-!ARCH!-v%VERSION%.tar.gz' -Force"
        )

        if errorlevel 1 (
            echo 警告: 压缩包创建失败
        ) else (
            echo 压缩包创建成功
        )

        echo !OS!-!ARCH! 构建完成！
        echo.
    )
)

echo ========================================
echo 构建摘要
echo ========================================
echo 版本号: %VERSION%
echo 构建时间: %date% %time%
echo 输出目录: %OUTPUT_DIR%
echo.
echo 生成的文件:
dir /b %OUTPUT_DIR%\*.zip 2>nul
dir /b %OUTPUT_DIR%\*.tar.gz 2>nul
echo.
echo 所有平台构建完成！
echo 输出文件位于 %OUTPUT_DIR% 目录下
echo.
echo 使用说明:
echo 1. 解压对应平台的压缩包
echo 2. 修改 config.yaml 配置文件
echo 3. 运行启动脚本或直接执行二进制文件
echo.

endlocal
