import axios from 'axios'

const api = axios.create({
  baseURL: '/api/v1',
  timeout: 30000
})

export const uploadFiles = async (files) => {
  const formData = new FormData()
  files.forEach(file => {
    formData.append('files', file)
  })
  const response = await api.post('/upload', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
  return response.data
}

export const getTaskStatus = async (taskId) => {
  const response = await api.get(`/task/${taskId}`)
  return response.data
}

export const getTaskResults = async (taskId) => {
  const response = await api.get(`/task/${taskId}/results`)
  return response.data
}

export const downloadResult = async (taskId) => {
  const response = await api.get(`/task/${taskId}/download`, {
    responseType: 'blob'
  })
  return response.data
}

export default api
