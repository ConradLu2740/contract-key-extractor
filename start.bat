@echo off
chcp 65001 >nul
title 合同关键信息提取系统

echo ========================================
echo    合同关键信息提取系统 - 一键启动
echo ========================================
echo.

:: 检查Python是否安装
python --version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到Python，请先安装Python 3.9+
    echo 下载地址: https://www.python.org/downloads/
    pause
    exit /b 1
)

:: 检查Go是否安装
go version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到Go，请先安装Go 1.21+
    echo 下载地址: https://go.dev/dl/
    pause
    exit /b 1
)

:: 检查Node.js是否安装
node --version >nul 2>&1
if errorlevel 1 (
    echo [错误] 未检测到Node.js，请先安装Node.js 18+
    echo 下载地址: https://nodejs.org/
    pause
    exit /b 1
)

:: 检查.env文件是否存在
if not exist "ai-service\.env" (
    echo [警告] 未找到ai-service\.env文件
    echo 请创建ai-service\.env文件并配置ZHIPU_API_KEY
    echo.
    echo 示例内容:
    echo ZHIPU_API_KEY=your_api_key_here
    echo.
    pause
    exit /b 1
)

:: 创建必要目录
if not exist "uploads" mkdir uploads
if not exist "outputs" mkdir outputs

echo [1/4] 安装Python依赖...
cd ai-service
pip install -r requirements.txt -q
cd ..

echo [2/4] 安装Go依赖...
go mod tidy

echo [3/4] 安装前端依赖...
cd web
call npm install --silent
cd ..

echo [4/4] 启动所有服务...
echo.
echo ========================================
echo   服务启动中，请稍候...
echo ========================================
echo.

:: 启动AI服务
start "AI Service" cmd /k "cd ai-service && python main.py"

:: 等待AI服务启动
timeout /t 3 /nobreak >nul

:: 启动Go后端
start "Go Backend" cmd /k "go run cmd/server/main.go"

:: 等待Go后端启动
timeout /t 3 /nobreak >nul

:: 启动前端
start "Web Frontend" cmd /k "cd web && npm run dev"

:: 等待前端启动
timeout /t 5 /nobreak >nul

echo.
echo ========================================
echo   所有服务已启动！
echo ========================================
echo.
echo   前端地址: http://localhost:3000
echo   后端地址: http://localhost:8080
echo   AI服务地址: http://localhost:8000
echo.
echo   请在浏览器中打开 http://localhost:3000
echo.
echo   关闭此窗口不会停止服务
echo   如需停止，请关闭对应的命令行窗口
echo ========================================
echo.

:: 自动打开浏览器
start http://localhost:3000

pause
