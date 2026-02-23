import os
import json
import httpx
from typing import Optional
from dotenv import load_dotenv
from pathlib import Path

env_path = Path(__file__).parent.parent / ".env"
load_dotenv(env_path)

from models.schemas import (
    ExtractionRequest,
    ExtractionResponse,
    ContractInfo,
    PartyInfo,
    FinancialInfo,
    ValidityInfo,
    RightsObligations,
    BreachLiability,
    DisputeResolution,
    ConfidentialityIP,
    OtherTerms,
    SignatureInfo,
    TypeSpecificFields,
    EmploymentFields,
    LeaseFields,
    LoanFields,
    ServiceFields,
    PurchaseFields
)


EXTRACTION_PROMPT = """你是一位专业的法律合同分析师。请从以下合同文本中提取关键信息，并以JSON格式返回。

合同文本:
{document_text}

请提取以下信息并严格按照以下JSON格式返回（注意字段名必须完全一致）:

{{
    "contract_info": {{
        "contract_type": "purchase或lease或loan或employment或service或other",
        "contract_number": "合同编号",
        "signing_date": "签订日期(YYYY-MM-DD格式)",
        "effective_date": "生效日期",
        "expiry_date": "到期日期",
        "signing_location": "签订地点",
        "contract_status": "合同状态",
        "confidence": 0.9,
        "source_references": []
    }},
    "party_a": {{
        "name": "甲方名称",
        "type": "企业或个人",
        "legal_representative": "法定代表人",
        "id_number": "身份证号或统一社会信用代码",
        "address": "地址",
        "contact": "联系方式",
        "bank_name": "开户银行",
        "bank_account": "银行账号",
        "confidence": 0.9,
        "source_references": []
    }},
    "party_b": {{
        "name": "乙方名称",
        "type": "企业或个人",
        "legal_representative": "法定代表人",
        "id_number": "身份证号或统一社会信用代码",
        "address": "地址",
        "contact": "联系方式",
        "bank_name": "开户银行",
        "bank_account": "银行账号",
        "confidence": 0.9,
        "source_references": []
    }},
    "financial": {{
        "transaction_amount": "交易金额",
        "currency": "CNY",
        "payment_method": "支付方式",
        "payment_schedule": "付款安排",
        "tax_info": "税务信息",
        "confidence": 0.9,
        "source_references": []
    }},
    "validity": {{
        "effective_condition": "生效条件",
        "termination_condition": "解除条件",
        "contract_status": "合同状态",
        "termination_date": "终止日期",
        "confidence": 0.9,
        "source_references": []
    }},
    "rights_obligations": {{
        "party_a_obligations": ["义务1", "义务2"],
        "party_b_obligations": ["义务1", "义务2"],
        "party_a_rights": ["权利1", "权利2"],
        "party_b_rights": ["权利1", "权利2"],
        "performance_period": "履行期限",
        "performance_location": "履行地点",
        "confidence": 0.9,
        "source_references": []
    }},
    "breach_liability": {{
        "breach_scenarios": ["违约情形1", "违约情形2"],
        "liquidated_damages": "违约金条款",
        "compensation_limit": "赔偿限额",
        "exemption_clauses": ["免责条款1"],
        "force_majeure_clause": "不可抗力条款",
        "confidence": 0.9,
        "source_references": []
    }},
    "dispute_resolution": {{
        "resolution_method": "诉讼或仲裁",
        "jurisdiction_court": "管辖法院",
        "arbitration_org": "仲裁机构",
        "arbitration_location": "仲裁地点",
        "governing_law": "适用法律",
        "confidence": 0.9,
        "source_references": []
    }},
    "confidentiality_ip": {{
        "confidentiality_clause": "保密条款",
        "confidentiality_period": "保密期限",
        "ip_ownership": "知识产权归属",
        "confidence": 0.9,
        "source_references": []
    }},
    "other_terms": {{
        "modification_clause": "变更条款",
        "assignment_clause": "转让条款",
        "termination_procedure": "解除程序",
        "notice_clause": "通知条款",
        "contract_copies": "合同份数",
        "attachments": [],
        "confidence": 0.9,
        "source_references": []
    }},
    "signature": {{
        "party_a_signatory": "甲方签字人",
        "party_a_sign_date": "甲方签字日期",
        "party_a_seal": false,
        "party_b_signatory": "乙方签字人",
        "party_b_sign_date": "乙方签字日期",
        "party_b_seal": false,
        "witness_name": "见证人姓名",
        "witness_contact": "见证人联系方式",
        "confidence": 0.9,
        "source_references": []
    }},
    "type_specific": {{
        "employment_fields": null,
        "lease_fields": null,
        "loan_fields": null,
        "service_fields": {{
            "service_content": "服务内容描述",
            "service_standard": "服务标准",
            "service_period": "服务期限",
            "service_fee": "服务费用",
            "acceptance_criteria": "验收标准",
            "confidence": 0.9
        }},
        "purchase_fields": null
    }}
}}

重要提示:
1. 字段名必须完全按照上面的格式，不能修改
2. 如果是服务合同(service)，type_specific.service_fields必须包含以下字段：service_content, service_standard, service_period, service_fee, acceptance_criteria, confidence
3. 如果无法提取某字段，填写"Unknown"，数组填写[]
4. 只返回JSON，不要额外解释
"""


class GLMExtractor:
    def __init__(self, api_key: Optional[str] = None):
        self.api_key = api_key or os.getenv("ZHIPU_API_KEY", "")
        self.base_url = "https://open.bigmodel.cn/api/paas/v4"
        self.model = "glm-4"
        
    async def extract(self, request: ExtractionRequest) -> ExtractionResponse:
        prompt = EXTRACTION_PROMPT.format(document_text=request.document_text)
        
        print(f"[DEBUG] Document text length: {len(request.document_text)}")
        
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
                            "role": "system",
                            "content": "你是一位专业的法律合同信息提取专家。请严格按照指定的JSON格式提取信息，字段名必须完全一致。"
                        },
                        {
                            "role": "user",
                            "content": prompt
                        }
                    ],
                    "temperature": 0.1,
                    "max_tokens": 4000
                }
            )
            
            if response.status_code != 200:
                raise Exception(f"GLM API error: {response.status_code} - {response.text}")
            
            result = response.json()
            content = result["choices"][0]["message"]["content"]
            
            print(f"[DEBUG] AI Response preview: {content[:500]}...")
            print(f"[DEBUG] AI Response full:\n{content}")
            
            return self._parse_response(content)
    
    def _parse_response(self, content: str) -> ExtractionResponse:
        try:
            json_str = content.strip()
            
            if json_str.startswith("```"):
                lines = json_str.split("\n")
                if lines[0].startswith("```"):
                    lines = lines[1:]
                if lines and lines[-1].strip() == "```":
                    lines = lines[:-1]
                json_str = "\n".join(lines)
            
            json_start = json_str.find("{")
            json_end = json_str.rfind("}") + 1
            if json_start >= 0 and json_end > json_start:
                json_str = json_str[json_start:json_end]
            
            print(f"[DEBUG] Cleaned JSON length: {len(json_str)}")
            
            try:
                data = json.loads(json_str)
            except json.JSONDecodeError as e:
                print(f"[DEBUG] First parse failed, trying to fix JSON: {e}")
                print(f"[DEBUG] Problem area: ...{json_str[max(0, e.pos-50):e.pos+50]}...")
                json_str = self._fix_json_string(json_str)
                print(f"[DEBUG] Fixed JSON length: {len(json_str)}")
                data = json.loads(json_str)
            
            return ExtractionResponse(
                contract_info=self._parse_contract_info(data.get("contract_info", {})),
                party_a=self._parse_party_info(data.get("party_a", {})),
                party_b=self._parse_party_info(data.get("party_b", {})),
                financial=self._parse_financial(data.get("financial", {})),
                validity=self._parse_validity(data.get("validity", {})),
                rights_obligations=self._parse_rights_obligations(data.get("rights_obligations", {})),
                breach_liability=self._parse_breach_liability(data.get("breach_liability", {})),
                dispute_resolution=self._parse_dispute_resolution(data.get("dispute_resolution", {})),
                confidentiality_ip=self._parse_confidentiality_ip(data.get("confidentiality_ip", {})),
                other_terms=self._parse_other_terms(data.get("other_terms", {})),
                signature=self._parse_signature(data.get("signature", {})),
                type_specific=self._parse_type_specific(data.get("type_specific", {})),
                ocr_required=False
            )
        except json.JSONDecodeError as e:
            print(f"[ERROR] JSON decode error: {e}")
            return self._create_default_response()
        except Exception as e:
            print(f"[ERROR] Parse error: {e}")
            return self._create_default_response()
    
    def _fix_json_string(self, json_str: str) -> str:
        result = []
        in_string = False
        escape_next = False
        i = 0
        
        while i < len(json_str):
            char = json_str[i]
            
            if escape_next:
                result.append(char)
                escape_next = False
                i += 1
                continue
            
            if char == '\\':
                result.append(char)
                escape_next = True
                i += 1
                continue
            
            if char == '"':
                if in_string:
                    next_non_space = i + 1
                    while next_non_space < len(json_str) and json_str[next_non_space] in ' \t\n\r':
                        next_non_space += 1
                    
                    if next_non_space < len(json_str) and json_str[next_non_space] in ':,}]':
                        in_string = False
                        result.append(char)
                    else:
                        result.append('\\"')
                else:
                    in_string = True
                    result.append(char)
                i += 1
                continue
            
            if in_string:
                if char == '\n':
                    result.append(' ')
                elif char == '\r':
                    pass
                elif char == '\t':
                    result.append(' ')
                else:
                    result.append(char)
            else:
                result.append(char)
            
            i += 1
        
        return ''.join(result)
    
    def _parse_contract_info(self, data: dict) -> ContractInfo:
        return ContractInfo(
            contract_type=data.get("contract_type", "other"),
            contract_number=data.get("contract_number", "Unknown"),
            signing_date=data.get("signing_date", "Unknown"),
            effective_date=data.get("effective_date", "Unknown"),
            expiry_date=data.get("expiry_date", "Unknown"),
            signing_location=data.get("signing_location", "Unknown"),
            contract_status=data.get("contract_status", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_party_info(self, data: dict) -> PartyInfo:
        return PartyInfo(
            name=data.get("name", "Unknown"),
            type=data.get("type", "Unknown"),
            legal_representative=data.get("legal_representative", "Unknown"),
            id_number=data.get("id_number", "Unknown"),
            address=data.get("address", "Unknown"),
            contact=data.get("contact", "Unknown"),
            bank_name=data.get("bank_name", "Unknown"),
            bank_account=data.get("bank_account", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_financial(self, data: dict) -> FinancialInfo:
        return FinancialInfo(
            transaction_amount=data.get("transaction_amount", "Unknown"),
            currency=data.get("currency", "Unknown"),
            payment_method=data.get("payment_method", "Unknown"),
            payment_schedule=data.get("payment_schedule", "Unknown"),
            tax_info=data.get("tax_info", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_validity(self, data: dict) -> ValidityInfo:
        return ValidityInfo(
            effective_condition=data.get("effective_condition", "Unknown"),
            termination_condition=data.get("termination_condition", "Unknown"),
            contract_status=data.get("contract_status", "Unknown"),
            termination_date=data.get("termination_date", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_rights_obligations(self, data: dict) -> RightsObligations:
        return RightsObligations(
            party_a_obligations=data.get("party_a_obligations", []),
            party_b_obligations=data.get("party_b_obligations", []),
            party_a_rights=data.get("party_a_rights", []),
            party_b_rights=data.get("party_b_rights", []),
            performance_period=data.get("performance_period", "Unknown"),
            performance_location=data.get("performance_location", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_breach_liability(self, data: dict) -> BreachLiability:
        return BreachLiability(
            breach_scenarios=data.get("breach_scenarios", []),
            liquidated_damages=data.get("liquidated_damages", "Unknown"),
            compensation_limit=data.get("compensation_limit", "Unknown"),
            exemption_clauses=data.get("exemption_clauses", []),
            force_majeure_clause=data.get("force_majeure_clause", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_dispute_resolution(self, data: dict) -> DisputeResolution:
        return DisputeResolution(
            resolution_method=data.get("resolution_method", "Unknown"),
            jurisdiction_court=data.get("jurisdiction_court", "Unknown"),
            arbitration_org=data.get("arbitration_org", "Unknown"),
            arbitration_location=data.get("arbitration_location", "Unknown"),
            governing_law=data.get("governing_law", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_confidentiality_ip(self, data: dict) -> ConfidentialityIP:
        return ConfidentialityIP(
            confidentiality_clause=data.get("confidentiality_clause", "Unknown"),
            confidentiality_period=data.get("confidentiality_period", "Unknown"),
            ip_ownership=data.get("ip_ownership", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_other_terms(self, data: dict) -> OtherTerms:
        return OtherTerms(
            modification_clause=data.get("modification_clause", "Unknown"),
            assignment_clause=data.get("assignment_clause", "Unknown"),
            termination_procedure=data.get("termination_procedure", "Unknown"),
            notice_clause=data.get("notice_clause", "Unknown"),
            contract_copies=data.get("contract_copies", "Unknown"),
            attachments=data.get("attachments", []),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_signature(self, data: dict) -> SignatureInfo:
        return SignatureInfo(
            party_a_signatory=data.get("party_a_signatory", "Unknown"),
            party_a_sign_date=data.get("party_a_sign_date", "Unknown"),
            party_a_seal=bool(data.get("party_a_seal", False)),
            party_b_signatory=data.get("party_b_signatory", "Unknown"),
            party_b_sign_date=data.get("party_b_sign_date", "Unknown"),
            party_b_seal=bool(data.get("party_b_seal", False)),
            witness_name=data.get("witness_name", "Unknown"),
            witness_contact=data.get("witness_contact", "Unknown"),
            confidence=float(data.get("confidence", 0.8)),
            source_references=data.get("source_references", [])
        )
    
    def _parse_type_specific(self, data: dict) -> TypeSpecificFields:
        result = TypeSpecificFields()
        
        if data:
            emp = data.get("employment_fields")
            if emp and isinstance(emp, dict):
                result.employment_fields = EmploymentFields(
                    position=emp.get("position", "Unknown"),
                    work_location=emp.get("work_location", "Unknown"),
                    work_hours=emp.get("work_hours", "Unknown"),
                    probation_period=emp.get("probation_period", "Unknown"),
                    salary=emp.get("salary", "Unknown"),
                    social_insurance=emp.get("social_insurance", "Unknown"),
                    non_compete_clause=emp.get("non_compete_clause", "Unknown"),
                    confidence=float(emp.get("confidence", 0.8))
                )
            
            lease = data.get("lease_fields")
            if lease and isinstance(lease, dict):
                result.lease_fields = LeaseFields(
                    leased_property=lease.get("leased_property", "Unknown"),
                    lease_area=lease.get("lease_area", "Unknown"),
                    lease_purpose=lease.get("lease_purpose", "Unknown"),
                    rent_amount=lease.get("rent_amount", "Unknown"),
                    rent_payment_cycle=lease.get("rent_payment_cycle", "Unknown"),
                    deposit=lease.get("deposit", "Unknown"),
                    maintenance_responsibility=lease.get("maintenance_responsibility", "Unknown"),
                    confidence=float(lease.get("confidence", 0.8))
                )
            
            loan = data.get("loan_fields")
            if loan and isinstance(loan, dict):
                result.loan_fields = LoanFields(
                    loan_amount=loan.get("loan_amount", "Unknown"),
                    loan_purpose=loan.get("loan_purpose", "Unknown"),
                    loan_term=loan.get("loan_term", "Unknown"),
                    interest_rate=loan.get("interest_rate", "Unknown"),
                    repayment_method=loan.get("repayment_method", "Unknown"),
                    collateral=loan.get("collateral", "Unknown"),
                    guarantor=loan.get("guarantor", "Unknown"),
                    confidence=float(loan.get("confidence", 0.8))
                )
            
            svc = data.get("service_fields")
            if svc and isinstance(svc, dict):
                result.service_fields = ServiceFields(
                    service_content=svc.get("service_content", svc.get("service_type", "Unknown")),
                    service_standard=svc.get("service_standard", "Unknown"),
                    service_period=svc.get("service_period", "Unknown"),
                    service_fee=svc.get("service_fee", svc.get("fee", "Unknown")),
                    acceptance_criteria=svc.get("acceptance_criteria", "Unknown"),
                    confidence=float(svc.get("confidence", 0.8))
                )
            
            purchase = data.get("purchase_fields")
            if purchase and isinstance(purchase, dict):
                result.purchase_fields = PurchaseFields(
                    goods_name=purchase.get("goods_name", "Unknown"),
                    goods_spec=purchase.get("goods_spec", "Unknown"),
                    goods_quantity=purchase.get("goods_quantity", "Unknown"),
                    goods_price=purchase.get("goods_price", "Unknown"),
                    delivery_location=purchase.get("delivery_location", "Unknown"),
                    delivery_date=purchase.get("delivery_date", "Unknown"),
                    quality_standard=purchase.get("quality_standard", "Unknown"),
                    warranty_period=purchase.get("warranty_period", "Unknown"),
                    confidence=float(purchase.get("confidence", 0.8))
                )
        
        return result
    
    def _create_default_response(self) -> ExtractionResponse:
        default_party = PartyInfo(
            name="Unknown", type="Unknown", legal_representative="Unknown",
            id_number="Unknown", address="Unknown", contact="Unknown",
            bank_name="Unknown", bank_account="Unknown", confidence=0.0,
            source_references=[]
        )
        
        return ExtractionResponse(
            contract_info=ContractInfo(
                contract_type="other", contract_number="Unknown",
                signing_date="Unknown", effective_date="Unknown",
                expiry_date="Unknown", signing_location="Unknown",
                contract_status="Unknown", confidence=0.0, source_references=[]
            ),
            party_a=default_party,
            party_b=default_party,
            financial=FinancialInfo(
                transaction_amount="Unknown", currency="Unknown",
                payment_method="Unknown", payment_schedule="Unknown",
                tax_info="Unknown", confidence=0.0, source_references=[]
            ),
            validity=ValidityInfo(
                effective_condition="Unknown", termination_condition="Unknown",
                contract_status="Unknown", termination_date="Unknown",
                confidence=0.0, source_references=[]
            ),
            rights_obligations=RightsObligations(
                party_a_obligations=[], party_b_obligations=[],
                party_a_rights=[], party_b_rights=[],
                performance_period="Unknown", performance_location="Unknown",
                confidence=0.0, source_references=[]
            ),
            breach_liability=BreachLiability(
                breach_scenarios=[], liquidated_damages="Unknown",
                compensation_limit="Unknown", exemption_clauses=[],
                force_majeure_clause="Unknown", confidence=0.0, source_references=[]
            ),
            dispute_resolution=DisputeResolution(
                resolution_method="Unknown", jurisdiction_court="Unknown",
                arbitration_org="Unknown", arbitration_location="Unknown",
                governing_law="Unknown", confidence=0.0, source_references=[]
            ),
            confidentiality_ip=ConfidentialityIP(
                confidentiality_clause="Unknown", confidentiality_period="Unknown",
                ip_ownership="Unknown", confidence=0.0, source_references=[]
            ),
            other_terms=OtherTerms(
                modification_clause="Unknown", assignment_clause="Unknown",
                termination_procedure="Unknown", notice_clause="Unknown",
                contract_copies="Unknown", attachments=[], confidence=0.0,
                source_references=[]
            ),
            signature=SignatureInfo(
                party_a_signatory="Unknown", party_a_sign_date="Unknown",
                party_a_seal=False, party_b_signatory="Unknown",
                party_b_sign_date="Unknown", party_b_seal=False,
                witness_name="Unknown", witness_contact="Unknown",
                confidence=0.0, source_references=[]
            ),
            type_specific=TypeSpecificFields(),
            ocr_required=False
        )
