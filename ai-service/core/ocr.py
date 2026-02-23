import base64
import httpx
from typing import Optional, List
import tempfile
import os
from dotenv import load_dotenv
from pathlib import Path
import io

env_path = Path(__file__).parent.parent / ".env"
load_dotenv(env_path)

OCR_PROMPT = """请从这张图片中提取所有文字内容。只返回提取的文字，不要添加任何解释。保持原有的段落结构。"""


class GLMOCR:
    def __init__(self, api_key: Optional[str] = None):
        self.api_key = api_key or os.getenv("ZHIPU_API_KEY", "")
        self.base_url = "https://open.bigmodel.cn/api/paas/v4"
        self.model = "glm-4v-plus-0111"
        
    def _compress_image(self, image_data: bytes, max_size: int = 800) -> bytes:
        from PIL import Image
        
        img = Image.open(io.BytesIO(image_data))
        
        if img.mode in ('RGBA', 'P'):
            img = img.convert('RGB')
        
        ratio = min(max_size / img.width, max_size / img.height)
        if ratio < 1:
            new_size = (int(img.width * ratio), int(img.height * ratio))
            img = img.resize(new_size, Image.LANCZOS)
        
        output = io.BytesIO()
        img.save(output, format='JPEG', quality=85, optimize=True)
        
        return output.getvalue()
        
    async def extract_text(self, image_data: bytes) -> str:
        compressed_data = self._compress_image(image_data)
        print(f"[DEBUG] Image compressed: {len(image_data)} -> {len(compressed_data)} bytes")
        
        base64_image = base64.b64encode(compressed_data).decode('utf-8')
        
        async with httpx.AsyncClient(timeout=120.0) as client:
            response = await client.post(
                f"{self.base_url}/chat/completions",
                headers={
                    "Authorization": f"Bearer {self.api_key}",
                    "Content-Type": "application/json"
                },
                json={
                    "model": self.model,
                    "messages": [
                        {
                            "role": "user",
                            "content": [
                                {
                                    "type": "image_url",
                                    "image_url": {
                                        "url": base64_image
                                    }
                                },
                                {
                                    "type": "text",
                                    "text": OCR_PROMPT
                                }
                            ]
                        }
                    ],
                    "temperature": 0.1,
                    "max_tokens": 4000
                }
            )
            
            if response.status_code != 200:
                raise Exception(f"GLM OCR API error: {response.status_code} - {response.text}")
            
            result = response.json()
            content = result["choices"][0]["message"]["content"]
            
            return content.strip()
    
    async def extract_text_from_images(self, images: List[bytes]) -> str:
        results = []
        for i, image_data in enumerate(images):
            try:
                text = await self.extract_text(image_data)
                results.append(f"--- 第{i+1}页 ---\n{text}")
            except Exception as e:
                print(f"[ERROR] Failed to extract text from image {i+1}: {e}")
                results.append(f"--- 第{i+1}页 ---\n[OCR识别失败]")
        
        return "\n\n".join(results)
    
    async def extract_text_from_pdf(self, pdf_data: bytes) -> str:
        print(f"[DEBUG] Starting PDF to image conversion, data size: {len(pdf_data)} bytes")
        images = await self._pdf_to_images(pdf_data)
        
        if not images:
            raise Exception("Failed to convert PDF to images")
        
        print(f"[DEBUG] Converted PDF to {len(images)} images")
        
        return await self.extract_text_from_images(images)
    
    async def _pdf_to_images(self, pdf_data: bytes) -> List[bytes]:
        import fitz
        
        images = []
        
        with tempfile.NamedTemporaryFile(suffix='.pdf', delete=False) as tmp_pdf:
            tmp_pdf.write(pdf_data)
            tmp_pdf_path = tmp_pdf.name
        
        try:
            doc = fitz.open(tmp_pdf_path)
            
            for page_num in range(len(doc)):
                page = doc[page_num]
                
                mat = fitz.Matrix(1.5, 1.5)
                pix = page.get_pixmap(matrix=mat)
                
                img_data = pix.tobytes("jpeg")
                images.append(img_data)
                
                print(f"[DEBUG] Converted page {page_num + 1} to image, size: {len(img_data)} bytes")
            
            doc.close()
        except Exception as e:
            print(f"[ERROR] PDF to image conversion failed: {e}")
            raise Exception(f"PDF conversion failed: {e}")
        finally:
            os.unlink(tmp_pdf_path)
        
        return images
