<template>
  <div class="home-container">
    <el-card class="upload-card">
      <template #header>
        <div class="card-header">
          <el-icon size="24"><Upload /></el-icon>
          <span>Upload Contract Files</span>
        </div>
      </template>
      
      <el-upload
        ref="uploadRef"
        class="upload-area"
        drag
        multiple
        :auto-upload="false"
        :on-change="handleFileChange"
        :file-list="fileList"
        accept=".pdf,.docx,.doc,.xlsx,.xls"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">
          Drag files here or <em>click to upload</em>
        </div>
        <template #tip>
          <div class="el-upload__tip">
            Supported formats: PDF, Word (.docx, .doc), Excel (.xlsx, .xls)
          </div>
        </template>
      </el-upload>
      
      <div class="file-list" v-if="fileList.length > 0">
        <h4>Selected Files ({{ fileList.length }})</h4>
        <el-tag 
          v-for="file in fileList" 
          :key="file.uid" 
          closable 
          @close="removeFile(file)"
          class="file-tag"
        >
          {{ file.name }}
        </el-tag>
      </div>
      
      <div class="action-buttons">
        <el-button 
          type="primary" 
          size="large" 
          @click="startExtraction"
          :loading="uploading"
          :disabled="fileList.length === 0"
        >
          <el-icon><Search /></el-icon>
          Start Extraction
        </el-button>
      </div>
    </el-card>
    
    <el-card class="info-card" v-if="taskId">
      <template #header>
        <div class="card-header">
          <el-icon size="24"><Clock /></el-icon>
          <span>Processing Status</span>
        </div>
      </template>
      
      <div class="status-content">
        <el-progress 
          :percentage="progress" 
          :status="progressStatus"
          :stroke-width="20"
        />
        <div class="status-info">
          <p><strong>Task ID:</strong> {{ taskId }}</p>
          <p><strong>Status:</strong> {{ statusText }}</p>
          <p><strong>Progress:</strong> {{ processed }} / {{ totalFiles }} files</p>
          <p v-if="failed > 0"><strong>Failed:</strong> {{ failed }} files</p>
        </div>
        <el-button 
          v-if="status === 'completed'" 
          type="success" 
          @click="viewResults"
        >
          View Results
        </el-button>
      </div>
    </el-card>
    
    <el-card class="features-card">
      <template #header>
        <div class="card-header">
          <el-icon size="24"><Document /></el-icon>
          <span>Extractable Information</span>
        </div>
      </template>
      
      <el-row :gutter="20">
        <el-col :span="12">
          <div class="feature-item">
            <el-icon><Tickets /></el-icon>
            <span>Contract Basic Info</span>
            <p>Contract type, number, signing date, effective date, expiry date</p>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="feature-item">
            <el-icon><User /></el-icon>
            <span>Party Information</span>
            <p>Party A/B name, legal representative, address, contact</p>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="feature-item">
            <el-icon><Money /></el-icon>
            <span>Financial Information</span>
            <p>Transaction amount, currency, payment method, schedule</p>
          </div>
        </el-col>
        <el-col :span="12">
          <div class="feature-item">
            <el-icon><DocumentCopy /></el-icon>
            <span>Key Terms</span>
            <p>Breach liability, dispute resolution, governing law</p>
          </div>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup>
import { ref, computed, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { uploadFiles, getTaskStatus } from '../api'

const router = useRouter()
const uploadRef = ref()
const fileList = ref([])
const uploading = ref(false)
const taskId = ref('')
const status = ref('')
const progress = ref(0)
const processed = ref(0)
const totalFiles = ref(0)
const failed = ref(0)

let pollInterval = null

const progressStatus = computed(() => {
  if (status.value === 'completed') return 'success'
  if (status.value === 'failed') return 'exception'
  return ''
})

const statusText = computed(() => {
  const statusMap = {
    'pending': 'Pending',
    'processing': 'Processing',
    'completed': 'Completed',
    'failed': 'Failed'
  }
  return statusMap[status.value] || status.value
})

const handleFileChange = (file, files) => {
  fileList.value = files
}

const removeFile = (file) => {
  fileList.value = fileList.value.filter(f => f.uid !== file.uid)
}

const startExtraction = async () => {
  if (fileList.value.length === 0) {
    ElMessage.warning('Please select files first')
    return
  }
  
  uploading.value = true
  
  try {
    const result = await uploadFiles(fileList.value.map(f => f.raw))
    taskId.value = result.task_id
    status.value = result.status
    totalFiles.value = result.total_files
    
    ElMessage.success('Files uploaded successfully, processing started')
    
    startPolling()
  } catch (error) {
    ElMessage.error('Upload failed: ' + (error.response?.data?.error || error.message))
  } finally {
    uploading.value = false
  }
}

const startPolling = () => {
  pollInterval = setInterval(async () => {
    try {
      const taskStatus = await getTaskStatus(taskId.value)
      status.value = taskStatus.status
      progress.value = taskStatus.progress
      processed.value = taskStatus.processed
      failed.value = taskStatus.failed
      
      if (taskStatus.status === 'completed' || taskStatus.status === 'failed') {
        stopPolling()
      }
    } catch (error) {
      console.error('Failed to get task status:', error)
    }
  }, 2000)
}

const stopPolling = () => {
  if (pollInterval) {
    clearInterval(pollInterval)
    pollInterval = null
  }
}

const viewResults = () => {
  router.push(`/results/${taskId.value}`)
}

onUnmounted(() => {
  stopPolling()
})
</script>

<style scoped>
.home-container {
  max-width: 900px;
  margin: 0 auto;
}

.upload-card, .info-card, .features-card {
  margin-bottom: 20px;
  border-radius: 12px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 10px;
  font-size: 18px;
  font-weight: 600;
}

.upload-area {
  width: 100%;
}

.file-list {
  margin-top: 20px;
  padding: 15px;
  background: #f5f7fa;
  border-radius: 8px;
}

.file-list h4 {
  margin-bottom: 10px;
  color: #333;
}

.file-tag {
  margin: 5px;
}

.action-buttons {
  margin-top: 20px;
  text-align: center;
}

.status-content {
  padding: 20px;
}

.status-info {
  margin-top: 20px;
}

.status-info p {
  margin: 8px 0;
  color: #666;
}

.features-card {
  background: rgba(255, 255, 255, 0.95);
}

.feature-item {
  padding: 15px;
  text-align: center;
}

.feature-item .el-icon {
  font-size: 32px;
  color: #409eff;
  margin-bottom: 10px;
}

.feature-item span {
  display: block;
  font-weight: 600;
  margin-bottom: 8px;
  color: #333;
}

.feature-item p {
  font-size: 12px;
  color: #999;
  line-height: 1.5;
}
</style>
