@echo off
chcp 65001 >nul
title 停止所有服务

echo ========================================
echo    停止合同关键信息提取系统所有服务
echo ========================================
echo.

echo 正在停止服务...

:: 停止Python进程 (AI服务)
taskkill /f /im python.exe 2>nul
if errorlevel 1 (
    echo AI服务未运行或已停止
) else (
    echo AI服务已停止
)

:: 停止Go进程 (后端服务)
taskkill /f /im main.exe 2>nul
if errorlevel 1 (
    echo Go后端未运行或已停止
) else (
    echo Go后端已停止
)

:: 停止Node进程 (前端服务)
taskkill /f /im node.exe 2>nul
if errorlevel 1 (
    echo 前端服务未运行或已停止
) else (
    echo 前端服务已停止
)

echo.
echo ========================================
echo   所有服务已停止
echo ========================================
echo.

pause
