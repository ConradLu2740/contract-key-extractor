# Quick Start Guide

## Prerequisites

Before starting, ensure you have:

1. **Go 1.21+** - [Download](https://golang.org/dl/)
2. **Python 3.8+** - [Download](https://www.python.org/downloads/)
3. **Node.js 16+** - [Download](https://nodejs.org/)
4. **Zhipu AI API Key** - Get free tokens at [https://open.bigmodel.cn/](https://open.bigmodel.cn/)

## Setup Steps

### 1. Get Zhipu AI API Key

1. Visit [https://open.bigmodel.cn/](https://open.bigmodel.cn/)
2. Register for a free account
3. Navigate to API Keys section
4. Create a new API key
5. Copy the key for later use

### 2. Configure API Key

Create a `.env` file in the `ai-service` directory:

```bash
cd ai-service
echo "ZHIPU_API_KEY=YOUR_API_KEY_HERE" > .env
```

**Important**: Replace `YOUR_API_KEY_HERE` with your actual API key from step 1.

### 3. Install Dependencies

```bash
# Install Go dependencies
go mod download

# Install Python dependencies
cd ai-service
pip install -r requirements.txt
cd ..

# Install Node dependencies
cd web
npm install
cd ..
```

### 4. Start Services

#### Option A: Quick Start (Windows)

```bash
start.bat
```

#### Option B: Manual Start

Open 3 separate terminal windows:

**Terminal 1 - AI Service:**
```bash
cd ai-service
python main.py
```

**Terminal 2 - Go Backend:**
```bash
go run cmd/server/main.go
```

**Terminal 3 - Frontend:**
```bash
cd web
npm run dev
```

### 5. Access the Application

Open your browser and navigate to:

- **Application**: http://localhost:3000
- **Backend API**: http://127.0.0.1:8080/health
- **AI Service**: http://127.0.0.1:8000/health

## Usage

1. **Upload Files**: Drag and drop your contract files (PDF, Word, Excel) to the upload area
2. **Start Extraction**: Click "Start Extraction" button
3. **Monitor Progress**: Watch the real-time progress bar
4. **View Results**: Click "View Results" when processing completes
5. **Export**: Download results as Excel file

## Troubleshooting

### Port Already in Use

If you see "port already in use" error:

```bash
# Find process using port 8080 (Windows)
netstat -ano | findstr :8080
taskkill /PID <PID> /F

# Find process using port 8080 (Linux/Mac)
lsof -ti:8080 | xargs kill -9
```

### API Key Error

Ensure your `.env` file is correctly formatted:

```env
ZHIPU_API_KEY=YOUR_API_KEY_HERE
```

**Note**: No quotes around the key!

### Module Not Found Errors

```bash
# Reinstall Go dependencies
go mod tidy
go mod download

# Reinstall Python dependencies
cd ai-service
pip install --upgrade pip
pip install -r requirements.txt

# Reinstall Node dependencies
cd web
rm -rf node_modules package-lock.json
npm install
```

## Next Steps

- Read the full [README.md](README.md) for detailed documentation
- Check [API Documentation](README.md#api-documentation) for integration
- Explore [Configuration](README.md#configuration) options

## Support

For issues or questions:
1. Check the [Troubleshooting](README.md#troubleshooting) section in README
2. Open an issue on GitHub
3. Contact the development team
