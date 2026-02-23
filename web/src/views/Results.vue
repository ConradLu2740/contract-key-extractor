<template>
  <div class="results-container">
    <el-page-header @back="goBack" title="Back">
      <template #content>
        <span class="page-title">Extraction Results</span>
      </template>
    </el-page-header>
    
    <div class="results-content" v-loading="loading">
      <template v-if="results.length > 0">
        <div class="summary-bar">
          <el-tag type="success">Successfully Extracted: {{ results.length }} files</el-tag>
          <el-button type="primary" @click="exportToExcel">
            <el-icon><Download /></el-icon>
            Export to Excel
          </el-button>
        </div>
        
        <el-collapse v-model="activeNames" accordion>
          <el-collapse-item 
            v-for="(result, index) in results" 
            :key="result.id"
            :name="index"
          >
            <template #title>
              <div class="collapse-title">
                <el-icon><Document /></el-icon>
                <span>{{ result.file_name }}</span>
                <el-tag 
                  :type="getConfidenceType(result.metadata?.overall_confidence || 0.8)"
                  size="small"
                >
                  Confidence: {{ ((result.metadata?.overall_confidence || 0.8) * 100).toFixed(1) }}%
                </el-tag>
              </div>
            </template>
            
            <div class="result-detail">
              <el-descriptions title="Contract Basic Info" :column="2" border>
                <el-descriptions-item label="Contract Type">
                  {{ getContractTypeLabel(result.contract_info?.contract_type) }}
                </el-descriptions-item>
                <el-descriptions-item label="Contract Number">
                  {{ result.contract_info?.contract_number || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Signing Date">
                  {{ result.contract_info?.signing_date || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Effective Date">
                  {{ result.contract_info?.effective_date || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Expiry Date">
                  {{ result.contract_info?.expiry_date || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Signing Location">
                  {{ result.contract_info?.signing_location || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Contract Status">
                  {{ result.contract_info?.contract_status || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Party A Information" :column="2" border>
                <el-descriptions-item label="Name">
                  {{ result.party_a?.name || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Type">
                  {{ result.party_a?.type || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Legal Representative">
                  {{ result.party_a?.legal_representative || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="ID Number">
                  {{ result.party_a?.id_number || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Address">
                  {{ result.party_a?.address || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Contact">
                  {{ result.party_a?.contact || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Bank Name">
                  {{ result.party_a?.bank_name || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Bank Account">
                  {{ result.party_a?.bank_account || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Party B Information" :column="2" border>
                <el-descriptions-item label="Name">
                  {{ result.party_b?.name || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Type">
                  {{ result.party_b?.type || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Legal Representative">
                  {{ result.party_b?.legal_representative || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="ID Number">
                  {{ result.party_b?.id_number || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Address">
                  {{ result.party_b?.address || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Contact">
                  {{ result.party_b?.contact || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Bank Name">
                  {{ result.party_b?.bank_name || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Bank Account">
                  {{ result.party_b?.bank_account || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Financial Information" :column="2" border>
                <el-descriptions-item label="Transaction Amount">
                  <span class="amount">{{ result.financial?.transaction_amount || 'Unknown' }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="Currency">
                  {{ result.financial?.currency || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Payment Method">
                  {{ result.financial?.payment_method || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Payment Schedule">
                  {{ result.financial?.payment_schedule || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Tax Info">
                  {{ result.financial?.tax_info || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Contract Validity" :column="2" border>
                <el-descriptions-item label="Effective Condition">
                  {{ result.validity?.effective_condition || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Termination Condition">
                  {{ result.validity?.termination_condition || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Contract Status">
                  {{ result.validity?.contract_status || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Termination Date">
                  {{ result.validity?.termination_date || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Rights and Obligations" :column="1" border>
                <el-descriptions-item label="Party A Obligations">
                  <el-tag v-for="(item, idx) in (result.rights_obligations?.party_a_obligations || [])" :key="'ao'+idx" class="clause-tag">{{ item }}</el-tag>
                  <span v-if="!result.rights_obligations?.party_a_obligations?.length">Unknown</span>
                </el-descriptions-item>
                <el-descriptions-item label="Party B Obligations">
                  <el-tag v-for="(item, idx) in (result.rights_obligations?.party_b_obligations || [])" :key="'bo'+idx" class="clause-tag">{{ item }}</el-tag>
                  <span v-if="!result.rights_obligations?.party_b_obligations?.length">Unknown</span>
                </el-descriptions-item>
                <el-descriptions-item label="Party A Rights">
                  <el-tag v-for="(item, idx) in (result.rights_obligations?.party_a_rights || [])" :key="'ar'+idx" class="clause-tag">{{ item }}</el-tag>
                  <span v-if="!result.rights_obligations?.party_a_rights?.length">Unknown</span>
                </el-descriptions-item>
                <el-descriptions-item label="Party B Rights">
                  <el-tag v-for="(item, idx) in (result.rights_obligations?.party_b_rights || [])" :key="'br'+idx" class="clause-tag">{{ item }}</el-tag>
                  <span v-if="!result.rights_obligations?.party_b_rights?.length">Unknown</span>
                </el-descriptions-item>
                <el-descriptions-item label="Performance Period">
                  {{ result.rights_obligations?.performance_period || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Performance Location">
                  {{ result.rights_obligations?.performance_location || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Breach Liability" :column="1" border>
                <el-descriptions-item label="Breach Scenarios">
                  <el-tag v-for="(item, idx) in (result.breach_liability?.breach_scenarios || [])" :key="'bs'+idx" type="danger" class="clause-tag">{{ item }}</el-tag>
                  <span v-if="!result.breach_liability?.breach_scenarios?.length">Unknown</span>
                </el-descriptions-item>
                <el-descriptions-item label="Liquidated Damages">
                  <span class="highlight">{{ result.breach_liability?.liquidated_damages || 'Unknown' }}</span>
                </el-descriptions-item>
                <el-descriptions-item label="Compensation Limit">
                  {{ result.breach_liability?.compensation_limit || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Exemption Clauses">
                  <el-tag v-for="(item, idx) in (result.breach_liability?.exemption_clauses || [])" :key="'ec'+idx" type="warning" class="clause-tag">{{ item }}</el-tag>
                  <span v-if="!result.breach_liability?.exemption_clauses?.length">Unknown</span>
                </el-descriptions-item>
                <el-descriptions-item label="Force Majeure Clause">
                  {{ result.breach_liability?.force_majeure_clause || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Dispute Resolution" :column="2" border>
                <el-descriptions-item label="Resolution Method">
                  {{ result.dispute_resolution?.resolution_method || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Jurisdiction Court">
                  {{ result.dispute_resolution?.jurisdiction_court || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Arbitration Org">
                  {{ result.dispute_resolution?.arbitration_org || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Arbitration Location">
                  {{ result.dispute_resolution?.arbitration_location || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Governing Law">
                  {{ result.dispute_resolution?.governing_law || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Confidentiality & IP" :column="2" border>
                <el-descriptions-item label="Confidentiality Clause">
                  {{ result.confidentiality_ip?.confidentiality_clause || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Confidentiality Period">
                  {{ result.confidentiality_ip?.confidentiality_period || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="IP Ownership">
                  {{ result.confidentiality_ip?.ip_ownership || 'Unknown' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Other Terms" :column="2" border>
                <el-descriptions-item label="Modification Clause">
                  {{ result.other_terms?.modification_clause || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Assignment Clause">
                  {{ result.other_terms?.assignment_clause || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Termination Procedure">
                  {{ result.other_terms?.termination_procedure || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Notice Clause">
                  {{ result.other_terms?.notice_clause || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Contract Copies">
                  {{ result.other_terms?.contract_copies || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Attachments">
                  <el-tag v-for="(item, idx) in (result.other_terms?.attachments || [])" :key="'att'+idx" class="clause-tag">{{ item }}</el-tag>
                  <span v-if="!result.other_terms?.attachments?.length">None</span>
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Signature Information" :column="2" border>
                <el-descriptions-item label="Party A Signatory">
                  {{ result.signature?.party_a_signatory || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Party A Sign Date">
                  {{ result.signature?.party_a_sign_date || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Party A Seal">
                  <el-tag :type="result.signature?.party_a_seal ? 'success' : 'info'" size="small">
                    {{ result.signature?.party_a_seal ? 'Yes' : 'No' }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="Party B Signatory">
                  {{ result.signature?.party_b_signatory || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Party B Sign Date">
                  {{ result.signature?.party_b_sign_date || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Party B Seal">
                  <el-tag :type="result.signature?.party_b_seal ? 'success' : 'info'" size="small">
                    {{ result.signature?.party_b_seal ? 'Yes' : 'No' }}
                  </el-tag>
                </el-descriptions-item>
                <el-descriptions-item label="Witness Name">
                  {{ result.signature?.witness_name || 'None' }}
                </el-descriptions-item>
                <el-descriptions-item label="Witness Contact">
                  {{ result.signature?.witness_contact || 'None' }}
                </el-descriptions-item>
              </el-descriptions>
              
              <el-divider />
              
              <el-descriptions title="Metadata" :column="2" border size="small">
                <el-descriptions-item label="Source File">
                  {{ result.metadata?.source_file || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Page Count">
                  {{ result.metadata?.page_count || 'Unknown' }}
                </el-descriptions-item>
                <el-descriptions-item label="Processing Time">
                  {{ (result.metadata?.processing_duration || 0).toFixed(2) }}s
                </el-descriptions-item>
                <el-descriptions-item label="OCR Required">
                  <el-tag :type="result.metadata?.ocr_required ? 'warning' : 'success'" size="small">
                    {{ result.metadata?.ocr_required ? 'Yes' : 'No' }}
                  </el-tag>
                </el-descriptions-item>
              </el-descriptions>
            </div>
          </el-collapse-item>
        </el-collapse>
      </template>
      
      <el-empty v-else-if="!loading" description="No extraction results available" />
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { getTaskResults, downloadResult } from '../api'

const router = useRouter()
const route = useRoute()
const loading = ref(true)
const results = ref([])
const activeNames = ref([0])

const contractTypeMap = {
  'purchase': 'Purchase Contract',
  'lease': 'Lease Contract',
  'loan': 'Loan Contract',
  'employment': 'Employment Contract',
  'service': 'Service Contract',
  'other': 'Other Contract'
}

const getContractTypeLabel = (type) => {
  return contractTypeMap[type] || type || 'Unknown'
}

const getConfidenceType = (confidence) => {
  if (confidence >= 0.8) return 'success'
  if (confidence >= 0.6) return 'warning'
  return 'danger'
}

const goBack = () => {
  router.push('/')
}

const fetchResults = async () => {
  const taskId = route.params.taskId
  if (!taskId) {
    ElMessage.error('Task ID not found')
    router.push('/')
    return
  }
  
  try {
    const response = await getTaskResults(taskId)
    results.value = response.results || []
  } catch (error) {
    ElMessage.error('Failed to fetch results: ' + (error.response?.data?.error || error.message))
  } finally {
    loading.value = false
  }
}

const exportToExcel = async () => {
  const taskId = route.params.taskId
  try {
    const blob = await downloadResult(taskId)
    const url = window.URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = `extraction_result_${taskId}.xlsx`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('Export successful')
  } catch (error) {
    ElMessage.error('Export failed: ' + (error.response?.data?.error || error.message))
  }
}

onMounted(() => {
  fetchResults()
})
</script>

<style scoped>
.results-container {
  max-width: 1000px;
  margin: 0 auto;
  background: white;
  border-radius: 12px;
  padding: 20px;
}

.page-title {
  font-size: 20px;
  font-weight: 600;
}

.results-content {
  margin-top: 20px;
}

.summary-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 8px;
}

.collapse-title {
  display: flex;
  align-items: center;
  gap: 10px;
  width: 100%;
}

.collapse-title .el-tag {
  margin-left: auto;
}

.result-detail {
  padding: 20px;
}

.amount {
  font-size: 18px;
  font-weight: 600;
  color: #409eff;
}

.highlight {
  font-weight: 600;
  color: #e6a23c;
}

.clause-tag {
  margin: 4px;
}

:deep(.el-descriptions__title) {
  font-size: 16px;
  color: #333;
}
</style>
