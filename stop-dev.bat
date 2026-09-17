@echo off
setlocal enabledelayedexpansion

rem ==============================================
rem  一键关闭本地前后端服务(按端口查找,存在才关闭)
rem  后端: Go server  端口 8080
rem  前端: Vite       端口 5173
rem  双击运行,或在终端里执行: stop-dev.bat
rem ==============================================

set "BACKEND_PORT=8080"
set "FRONTEND_PORT=5173"

echo 正在查找并关闭服务...
echo.

call :stop_port %BACKEND_PORT% 后端
call :stop_port %FRONTEND_PORT% 前端

echo.
echo 全部完成。
pause
exit /b

:stop_port
set "port=%~1"
set "name=%~2"
set "lastpid="
set "found=0"
for /f "tokens=5" %%a in ('netstat -ano ^| findstr /R /C:":%port% .*LISTENING"') do (
    if not "!lastpid!"=="%%a" (
        set "lastpid=%%a"
        set "found=1"
        taskkill /PID %%a /F >nul 2>&1
        if errorlevel 1 (
            echo [失败] !name! 端口 %port%,PID %%a 关闭失败,可能需要管理员权限
        ) else (
            echo [OK] 已关闭 !name! 进程 PID %%a,端口 %port%
        )
    )
)
if "!found!"=="0" echo [--] !name! 未在运行,端口 %port% 无监听
goto :eof
