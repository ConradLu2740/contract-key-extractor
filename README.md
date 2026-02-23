# 合同关键信息提取系统

基于AI的合同关键信息自动提取系统，支持Word、Excel、PDF格式合同文件的智能解析。

## 功能特点

- 支持多种文件格式：Word (.docx)、Excel (.xlsx)、PDF
- 智能OCR识别：扫描版PDF自动识别文字
- 结构化信息提取：自动提取合同双方、金额、期限、权利义务等关键信息
- Excel导出：一键导出提取结果到Excel文件
- 批量处理：支持多文件批量上传处理

## 技术栈

- **后端**: Go 1.21+ (Gin框架)
- **AI服务**: Python 3.9+ (FastAPI + 智谱AI GLM-4V)
- **前端**: Vue 3 + TypeScript + Element Plus

## 目录结构

```
contract-key-extractor/
├── cmd/                    # 应用入口
├── internal/               # 内部模块
│   ├── config/            # 配置管理
│   ├── handler/           # HTTP处理器
│   ├── model/             # 数据模型
│   ├── parser/            # 文件解析器
│   └── service/           # 业务逻辑
├── ai-service/            # AI服务
│   ├── core/              # 核心模块
│   ├── models/            # 数据模型
│   └── main.py            # 服务入口
├── web/                   # 前端项目
├── uploads/               # 上传文件目录
├── outputs/               # 输出文件目录
└── start.bat              # 一键启动脚本
```

## 快速开始

### 前置要求

1. 安装 [Go 1.21+](https://go.dev/dl/)
2. 安装 [Python 3.9+](https://www.python.org/downloads/)
3. 安装 [Node.js 18+](https://nodejs.org/)
4. 获取 [智谱AI API Key](https://open.bigmodel.cn/)

### 一键启动（推荐）

1. 双击运行 `start.bat`
2. 等待所有服务启动完成
3. 打开浏览器访问 http://localhost:3000

### 手动启动

#### 1. 配置环境变量

在 `ai-service/.env` 文件中配置智谱AI API Key：

```env
ZHIPU_API_KEY=your_api_key_here
```

#### 2. 启动AI服务

```bash
cd ai-service
pip install -r requirements.txt
python main.py
```

#### 3. 启动Go后端

```bash
go mod tidy
go run cmd/server/main.go
```

#### 4. 启动前端

```bash
cd web
npm install
npm run dev
```

#### 5. 访问应用

打开浏览器访问 http://localhost:3000

## 配置说明

### AI服务配置 (ai-service/.env)

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| ZHIPU_API_KEY | 智谱AI API密钥 | - |

### 后端配置 (config.yaml)

| 配置项 | 说明 | 默认值 |
|--------|------|--------|
| server.address | 服务地址 | 127.0.0.1:8080 |
| ai_service.host | AI服务地址 | 127.0.0.1 |
| ai_service.port | AI服务端口 | 8000 |
| upload.path | 上传目录 | ./uploads |
| output.path | 输出目录 | ./outputs |

## 使用说明

1. **上传文件**: 点击上传区域或拖拽文件上传
2. **等待处理**: 系统自动解析并提取信息
3. **查看结果**: 在页面查看提取的结构化信息
4. **导出Excel**: 点击"导出Excel"按钮下载结果

## 支持的合同类型

- 服务合同
- 买卖合同
- 租赁合同
- 借款合同
- 劳动合同
- 其他类型合同

## 常见问题

### Q: PDF文件提取结果为空？

A: 请确保PDF文件不是加密的，且内容清晰可读。扫描版PDF需要较长的OCR处理时间。

### Q: 启动失败？

A: 请检查：
1. 端口3000、8080、8000是否被占用
2. Python依赖是否正确安装
3. API Key是否正确配置

## 许可证

MIT License
