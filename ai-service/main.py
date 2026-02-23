import os
from dotenv import load_dotenv
from pathlib import Path

env_path = Path(__file__).parent / ".env"
load_dotenv(env_path)

from fastapi import FastAPI, HTTPException, UploadFile, File
from fastapi.middleware.cors import CORSMiddleware
from contextlib import asynccontextmanager

from models.schemas import ExtractionRequest, ExtractionResponse, OCRResponse
from core.extractor import GLMExtractor
from core.ocr import GLMOCR


@asynccontextmanager
async def lifespan(app: FastAPI):
    api_key = os.getenv("ZHIPU_API_KEY", "")
    print(f"Loaded API Key: {api_key[:20]}..." if api_key else "API Key not found!")
    app.state.extractor = GLMExtractor(api_key)
    app.state.ocr = GLMOCR(api_key)
    yield


app = FastAPI(
    title="Contract Key Information Extraction AI Service",
    description="AI service for extracting key information from legal contracts",
    version="1.0.0",
    lifespan=lifespan
)

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],
    allow_credentials=True,
    allow_methods=["*"],
    allow_headers=["*"],
)


@app.get("/health")
async def health_check():
    return {"status": "healthy", "service": "ai-service"}


@app.post("/api/v1/extract", response_model=ExtractionResponse)
async def extract_contract_info(request: ExtractionRequest):
    try:
        extractor: GLMExtractor = app.state.extractor
        result = await extractor.extract(request)
        return result
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/ocr", response_model=OCRResponse)
async def perform_ocr(file: UploadFile = File(...)):
    try:
        image_data = await file.read()
        ocr: GLMOCR = app.state.ocr
        text = await ocr.extract_text(image_data)
        return OCRResponse(text=text)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/ocr/raw", response_model=OCRResponse)
async def perform_ocr_raw(file: UploadFile = File(...)):
    try:
        image_data = await file.read()
        ocr: GLMOCR = app.state.ocr
        text = await ocr.extract_text(image_data)
        return OCRResponse(text=text)
    except Exception as e:
        raise HTTPException(status_code=500, detail=str(e))


@app.post("/api/v1/ocr/pdf", response_model=OCRResponse)
async def perform_pdf_ocr(file: UploadFile = File(...)):
    try:
        pdf_data = await file.read()
        print(f"[DEBUG] Received PDF data, size: {len(pdf_data)} bytes")
        ocr: GLMOCR = app.state.ocr
        text = await ocr.extract_text_from_pdf(pdf_data)
        print(f"[DEBUG] OCR result length: {len(text)} chars")
        return OCRResponse(text=text)
    except Exception as e:
        print(f"[ERROR] PDF OCR failed: {e}")
        raise HTTPException(status_code=500, detail=str(e))


if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="127.0.0.1", port=8000)
