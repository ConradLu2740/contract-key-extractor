from pydantic import BaseModel
from typing import List, Optional
from enum import Enum


class ContractType(str, Enum):
    PURCHASE = "purchase"
    LEASE = "lease"
    LOAN = "loan"
    EMPLOYMENT = "employment"
    SERVICE = "service"
    OTHER = "other"


class SourceRef(BaseModel):
    page: int
    paragraph: int
    text: str


class ContractInfo(BaseModel):
    contract_type: str
    contract_number: str
    signing_date: str
    effective_date: str
    expiry_date: str
    signing_location: str
    contract_status: str
    confidence: float
    source_references: List[SourceRef] = []


class PartyInfo(BaseModel):
    name: str
    type: str
    legal_representative: str
    id_number: str
    address: str
    contact: str
    bank_name: str
    bank_account: str
    confidence: float
    source_references: List[SourceRef] = []


class FinancialInfo(BaseModel):
    transaction_amount: str
    currency: str
    payment_method: str
    payment_schedule: str
    tax_info: str
    confidence: float
    source_references: List[SourceRef] = []


class ValidityInfo(BaseModel):
    effective_condition: str
    termination_condition: str
    contract_status: str
    termination_date: str
    confidence: float
    source_references: List[SourceRef] = []


class RightsObligations(BaseModel):
    party_a_obligations: List[str]
    party_b_obligations: List[str]
    party_a_rights: List[str]
    party_b_rights: List[str]
    performance_period: str
    performance_location: str
    confidence: float
    source_references: List[SourceRef] = []


class BreachLiability(BaseModel):
    breach_scenarios: List[str]
    liquidated_damages: str
    compensation_limit: str
    exemption_clauses: List[str]
    force_majeure_clause: str
    confidence: float
    source_references: List[SourceRef] = []


class DisputeResolution(BaseModel):
    resolution_method: str
    jurisdiction_court: str
    arbitration_org: str
    arbitration_location: str
    governing_law: str
    confidence: float
    source_references: List[SourceRef] = []


class ConfidentialityIP(BaseModel):
    confidentiality_clause: str
    confidentiality_period: str
    ip_ownership: str
    confidence: float
    source_references: List[SourceRef] = []


class OtherTerms(BaseModel):
    modification_clause: str
    assignment_clause: str
    termination_procedure: str
    notice_clause: str
    contract_copies: str
    attachments: List[str]
    confidence: float
    source_references: List[SourceRef] = []


class SignatureInfo(BaseModel):
    party_a_signatory: str
    party_a_sign_date: str
    party_a_seal: bool
    party_b_signatory: str
    party_b_sign_date: str
    party_b_seal: bool
    witness_name: str
    witness_contact: str
    confidence: float
    source_references: List[SourceRef] = []


class EmploymentFields(BaseModel):
    position: str
    work_location: str
    work_hours: str
    probation_period: str
    salary: str
    social_insurance: str
    non_compete_clause: str
    confidence: float


class LeaseFields(BaseModel):
    leased_property: str
    lease_area: str
    lease_purpose: str
    rent_amount: str
    rent_payment_cycle: str
    deposit: str
    maintenance_responsibility: str
    confidence: float


class LoanFields(BaseModel):
    loan_amount: str
    loan_purpose: str
    loan_term: str
    interest_rate: str
    repayment_method: str
    collateral: str
    guarantor: str
    confidence: float


class ServiceFields(BaseModel):
    service_content: str
    service_standard: str
    service_period: str
    service_fee: str
    acceptance_criteria: str
    confidence: float


class PurchaseFields(BaseModel):
    goods_name: str
    goods_spec: str
    goods_quantity: str
    goods_price: str
    delivery_location: str
    delivery_date: str
    quality_standard: str
    warranty_period: str
    confidence: float


class TypeSpecificFields(BaseModel):
    employment_fields: Optional[EmploymentFields] = None
    lease_fields: Optional[LeaseFields] = None
    loan_fields: Optional[LoanFields] = None
    service_fields: Optional[ServiceFields] = None
    purchase_fields: Optional[PurchaseFields] = None


class ExtractionRequest(BaseModel):
    document_text: str
    contract_type: Optional[str] = None


class ExtractionResponse(BaseModel):
    contract_info: ContractInfo
    party_a: PartyInfo
    party_b: PartyInfo
    financial: FinancialInfo
    validity: ValidityInfo
    rights_obligations: RightsObligations
    breach_liability: BreachLiability
    dispute_resolution: DisputeResolution
    confidentiality_ip: ConfidentialityIP
    other_terms: OtherTerms
    signature: SignatureInfo
    type_specific: TypeSpecificFields
    ocr_required: bool


class OCRRequest(BaseModel):
    pass


class OCRResponse(BaseModel):
    text: str
