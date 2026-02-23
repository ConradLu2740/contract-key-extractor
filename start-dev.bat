@echo off
echo Starting Contract Key Information Extractor...

echo.
echo [1/3] Starting AI Service...
start "AI Service" cmd /k "cd ai-service && python main.py"

timeout /t 3 /nobreak > nul

echo.
echo [2/3] Starting Go Backend...
start "Go Backend" cmd /k "go run cmd/server/main.go"

timeout /t 3 /nobreak > nul

echo.
echo [3/3] Starting Frontend Development Server...
start "Frontend" cmd /k "cd web && npm run dev"

echo.
echo All services started!
echo.
echo Access the application at: http://localhost:3000
echo API Health Check: http://127.0.0.1:8080/health
echo AI Service Health: http://127.0.0.1:8000/health
echo.
