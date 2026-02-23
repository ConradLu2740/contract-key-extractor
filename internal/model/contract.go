package model

import "time"

type ContractType string

const (
	ContractTypePurchase   ContractType = "purchase"
	ContractTypeLease      ContractType = "lease"
	ContractTypeLoan       ContractType = "loan"
	ContractTypeEmployment ContractType = "employment"
	ContractTypeService    ContractType = "service"
	ContractTypeOther      ContractType = "other"
)

type ContractInfo struct {
	ContractType     ContractType `json:"contract_type"`
	ContractNumber   string       `json:"contract_number"`
	SigningDate      string       `json:"signing_date"`
	EffectiveDate    string       `json:"effective_date"`
	ExpiryDate       string       `json:"expiry_date"`
	SigningLocation  string       `json:"signing_location"`
	ContractStatus   string       `json:"contract_status"`
	Confidence       float64      `json:"confidence"`
	SourceReferences []SourceRef  `json:"source_references"`
}

type PartyInfo struct {
	Name                string      `json:"name"`
	Type                string      `json:"type"`
	LegalRepresentative string      `json:"legal_representative"`
	IDNumber            string      `json:"id_number"`
	Address             string      `json:"address"`
	Contact             string      `json:"contact"`
	BankName            string      `json:"bank_name"`
	BankAccount         string      `json:"bank_account"`
	Confidence          float64     `json:"confidence"`
	SourceReferences    []SourceRef `json:"source_references"`
}

type FinancialInfo struct {
	TransactionAmount string      `json:"transaction_amount"`
	Currency          string      `json:"currency"`
	PaymentMethod     string      `json:"payment_method"`
	PaymentSchedule   string      `json:"payment_schedule"`
	TaxInfo           string      `json:"tax_info"`
	Confidence        float64     `json:"confidence"`
	SourceReferences  []SourceRef `json:"source_references"`
}

type ValidityInfo struct {
	EffectiveCondition   string      `json:"effective_condition"`
	TerminationCondition string      `json:"termination_condition"`
	ContractStatus       string      `json:"contract_status"`
	TerminationDate      string      `json:"termination_date"`
	Confidence           float64     `json:"confidence"`
	SourceReferences     []SourceRef `json:"source_references"`
}

type RightsObligations struct {
	PartyAObligations   []string    `json:"party_a_obligations"`
	PartyBObligations   []string    `json:"party_b_obligations"`
	PartyARights        []string    `json:"party_a_rights"`
	PartyBRights        []string    `json:"party_b_rights"`
	PerformancePeriod   string      `json:"performance_period"`
	PerformanceLocation string      `json:"performance_location"`
	Confidence          float64     `json:"confidence"`
	SourceReferences    []SourceRef `json:"source_references"`
}

type BreachLiability struct {
	BreachScenarios    []string    `json:"breach_scenarios"`
	LiquidatedDamages  string      `json:"liquidated_damages"`
	CompensationLimit  string      `json:"compensation_limit"`
	ExemptionClauses   []string    `json:"exemption_clauses"`
	ForceMajeureClause string      `json:"force_majeure_clause"`
	Confidence         float64     `json:"confidence"`
	SourceReferences   []SourceRef `json:"source_references"`
}

type DisputeResolution struct {
	ResolutionMethod    string      `json:"resolution_method"`
	JurisdictionCourt   string      `json:"jurisdiction_court"`
	ArbitrationOrg      string      `json:"arbitration_org"`
	ArbitrationLocation string      `json:"arbitration_location"`
	GoverningLaw        string      `json:"governing_law"`
	Confidence          float64     `json:"confidence"`
	SourceReferences    []SourceRef `json:"source_references"`
}

type ConfidentialityIP struct {
	ConfidentialityClause string      `json:"confidentiality_clause"`
	ConfidentialityPeriod string      `json:"confidentiality_period"`
	IPOwnership           string      `json:"ip_ownership"`
	Confidence            float64     `json:"confidence"`
	SourceReferences      []SourceRef `json:"source_references"`
}

type OtherTerms struct {
	ModificationClause   string      `json:"modification_clause"`
	AssignmentClause     string      `json:"assignment_clause"`
	TerminationProcedure string      `json:"termination_procedure"`
	NoticeClause         string      `json:"notice_clause"`
	ContractCopies       string      `json:"contract_copies"`
	Attachments          []string    `json:"attachments"`
	Confidence           float64     `json:"confidence"`
	SourceReferences     []SourceRef `json:"source_references"`
}

type SignatureInfo struct {
	PartyASignatory  string      `json:"party_a_signatory"`
	PartyASignDate   string      `json:"party_a_sign_date"`
	PartyASeal       bool        `json:"party_a_seal"`
	PartyBSignatory  string      `json:"party_b_signatory"`
	PartyBSignDate   string      `json:"party_b_sign_date"`
	PartyBSeal       bool        `json:"party_b_seal"`
	WitnessName      string      `json:"witness_name"`
	WitnessContact   string      `json:"witness_contact"`
	Confidence       float64     `json:"confidence"`
	SourceReferences []SourceRef `json:"source_references"`
}

type TypeSpecificFields struct {
	EmploymentFields *EmploymentFields `json:"employment_fields,omitempty"`
	LeaseFields      *LeaseFields      `json:"lease_fields,omitempty"`
	LoanFields       *LoanFields       `json:"loan_fields,omitempty"`
	ServiceFields    *ServiceFields    `json:"service_fields,omitempty"`
	PurchaseFields   *PurchaseFields   `json:"purchase_fields,omitempty"`
}

type EmploymentFields struct {
	Position         string  `json:"position"`
	WorkLocation     string  `json:"work_location"`
	WorkHours        string  `json:"work_hours"`
	ProbationPeriod  string  `json:"probation_period"`
	Salary           string  `json:"salary"`
	SocialInsurance  string  `json:"social_insurance"`
	NonCompeteClause string  `json:"non_compete_clause"`
	Confidence       float64 `json:"confidence"`
}

type LeaseFields struct {
	LeasedProperty   string  `json:"leased_property"`
	LeaseArea        string  `json:"lease_area"`
	LeasePurpose     string  `json:"lease_purpose"`
	RentAmount       string  `json:"rent_amount"`
	RentPaymentCycle string  `json:"rent_payment_cycle"`
	Deposit          string  `json:"deposit"`
	MaintenanceResp  string  `json:"maintenance_responsibility"`
	Confidence       float64 `json:"confidence"`
}

type LoanFields struct {
	LoanAmount      string  `json:"loan_amount"`
	LoanPurpose     string  `json:"loan_purpose"`
	LoanTerm        string  `json:"loan_term"`
	InterestRate    string  `json:"interest_rate"`
	RepaymentMethod string  `json:"repayment_method"`
	Collateral      string  `json:"collateral"`
	Guarantor       string  `json:"guarantor"`
	Confidence      float64 `json:"confidence"`
}

type ServiceFields struct {
	ServiceContent     string  `json:"service_content"`
	ServiceStandard    string  `json:"service_standard"`
	ServicePeriod      string  `json:"service_period"`
	ServiceFee         string  `json:"service_fee"`
	AcceptanceCriteria string  `json:"acceptance_criteria"`
	Confidence         float64 `json:"confidence"`
}

type PurchaseFields struct {
	GoodsName        string  `json:"goods_name"`
	GoodsSpec        string  `json:"goods_spec"`
	GoodsQuantity    string  `json:"goods_quantity"`
	GoodsPrice       string  `json:"goods_price"`
	DeliveryLocation string  `json:"delivery_location"`
	DeliveryDate     string  `json:"delivery_date"`
	QualityStandard  string  `json:"quality_standard"`
	WarrantyPeriod   string  `json:"warranty_period"`
	Confidence       float64 `json:"confidence"`
}

type SourceRef struct {
	Page      int    `json:"page"`
	Paragraph int    `json:"paragraph"`
	Text      string `json:"text"`
}

type ExtractionResult struct {
	ID                string             `json:"id"`
	FileName          string             `json:"file_name"`
	ContractInfo      ContractInfo       `json:"contract_info"`
	PartyA            PartyInfo          `json:"party_a"`
	PartyB            PartyInfo          `json:"party_b"`
	Financial         FinancialInfo      `json:"financial"`
	Validity          ValidityInfo       `json:"validity"`
	RightsObligations RightsObligations  `json:"rights_obligations"`
	BreachLiability   BreachLiability    `json:"breach_liability"`
	DisputeResolution DisputeResolution  `json:"dispute_resolution"`
	ConfidentialityIP ConfidentialityIP  `json:"confidentiality_ip"`
	OtherTerms        OtherTerms         `json:"other_terms"`
	Signature         SignatureInfo      `json:"signature"`
	TypeSpecific      TypeSpecificFields `json:"type_specific"`
	Metadata          Metadata           `json:"metadata"`
}

type Metadata struct {
	SourceFile          string    `json:"source_file"`
	PageCount           int       `json:"page_count"`
	ExtractionTime      time.Time `json:"extraction_time"`
	ProcessingDuration  float64   `json:"processing_duration"`
	OverallConfidence   float64   `json:"overall_confidence"`
	OCRRequired         bool      `json:"ocr_required"`
	ContractTypeChinese string    `json:"contract_type_chinese"`
}

type ExtractionRequest struct {
	Files        []string `json:"files"`
	OutputFormat string   `json:"output_format"`
}

type ExtractionResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Results []ExtractionResult `json:"results"`
}

type TaskStatus struct {
	TaskID      string  `json:"task_id"`
	Status      string  `json:"status"`
	Progress    float64 `json:"progress"`
	TotalFiles  int     `json:"total_files"`
	Processed   int     `json:"processed"`
	Failed      int     `json:"failed"`
	ResultPath  string  `json:"result_path"`
	Error       string  `json:"error,omitempty"`
	CreatedAt   string  `json:"created_at"`
	CompletedAt string  `json:"completed_at,omitempty"`
}

type FileType string

const (
	FileTypePDF   FileType = "pdf"
	FileTypeExcel FileType = "excel"
	FileTypeWord  FileType = "word"
	FileTypeImage FileType = "image"
)

type ParsedDocument struct {
	FileName   string   `json:"file_name"`
	FileType   FileType `json:"file_type"`
	Content    string   `json:"content"`
	PageCount  int      `json:"page_count"`
	IsScanned  bool     `json:"is_scanned"`
	ImagePaths []string `json:"image_paths,omitempty"`
}

type AIExtractionRequest struct {
	DocumentText string `json:"document_text"`
	ContractType string `json:"contract_type,omitempty"`
}

type AIExtractionResponse struct {
	ContractInfo      ContractInfo       `json:"contract_info"`
	PartyA            PartyInfo          `json:"party_a"`
	PartyB            PartyInfo          `json:"party_b"`
	Financial         FinancialInfo      `json:"financial"`
	Validity          ValidityInfo       `json:"validity"`
	RightsObligations RightsObligations  `json:"rights_obligations"`
	BreachLiability   BreachLiability    `json:"breach_liability"`
	DisputeResolution DisputeResolution  `json:"dispute_resolution"`
	ConfidentialityIP ConfidentialityIP  `json:"confidentiality_ip"`
	OtherTerms        OtherTerms         `json:"other_terms"`
	Signature         SignatureInfo      `json:"signature"`
	TypeSpecific      TypeSpecificFields `json:"type_specific"`
	OCRRequired       bool               `json:"ocr_required"`
}
