#!/bin/bash

echo "Starting Contract Key Information Extractor..."

echo ""
echo "[1/3] Starting AI Service..."
cd ai-service
python main.py &
AI_PID=$!
cd ..

sleep 3

echo ""
echo "[2/3] Starting Go Backend..."
go run cmd/server/main.go &
BACKEND_PID=$!

sleep 3

echo ""
echo "[3/3] Starting Frontend Development Server..."
cd web
npm run dev &
FRONTEND_PID=$!
cd ..

echo ""
echo "All services started!"
echo ""
echo "Access the application at: http://localhost:3000"
echo "API Health Check: http://127.0.0.1:8080/health"
echo "AI Service Health: http://127.0.0.1:8000/health"
echo ""
echo "Press Ctrl+C to stop all services"
echo ""

trap "kill $AI_PID $BACKEND_PID $FRONTEND_PID" EXIT

wait
