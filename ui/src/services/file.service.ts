import req, { getUrlPrefix, getToken } from '@/utils/axios';

export const fileService = {
  getFiles,
  getRevisions,
  getHistories,
  removeFile,
  getCommentCount,
  getMentionAt,
  createShimoFile,
  getShimoFile,
  updateFile,
  uploadFile,
  importFile,
  importFileByUrl,
  duplicateFile,
  getAllCollaborators,
  saveCollaborators,
  exportFile,
  getPlainTextAPIUrl,
  getImportProgress,
  getExportProgress,
  exportTableFile,
  batchDeleteFile,
  getImportByUrlProgress,
  getWebInspect
}

export interface DataType {
  createdAt: number;
  creatorId: number;
  filePath: string;
  id: string;
  isShimoFile: number;
  name: string;
  permissions: {
    commentable: boolean;
    copyable: boolean
    editable: boolean
    exportable: boolean
    formFillable: boolean
    lockable: boolean
    manageable: boolean
    readable: boolean
    unlockable: boolean
  };
  shimoType: string;
  type: string;
  updatedAt: number;
}

async function getFiles() {
  return await req.get(`api/files/`)
}

async function getRevisions(fileId: string) {
  return await req.get(`/api/files/${fileId}/revisions`)
}

async function getHistories(fileId: string) {
  return await req.get(`api/files/${fileId}/doc-sidebar-info`)
}

async function removeFile(fileId: string) {
  return await req.delete(`/api/files/${fileId}`)
}

async function getCommentCount(fileId: string) {
  return await req.get(`api/files/${fileId}/comment-count`)
}

async function getMentionAt(fileId: string) {
  return await req.get(`api/files/${fileId}/mention-at-list`)
}

async function createShimoFile({ type, lang }: { type: string, lang: string }) {
  let queryString = ''
  if (lang) {
    queryString = `?lang=${lang}`
  }
  return await req.post(`api/files${queryString}`, { shimoType: type })
}

async function getShimoFile({ id, mode, smParams, lang }: { id: string, mode: string, smParams: string, lang: string }) {
  return await req.get(`api/files/${id}`, {
    params: {
      mode: mode,
      smParams: smParams,
      lang: lang
    }
  })
}

async function updateFile(id: string, data: any) {
  return await req.patch(`api/files/${id}`, data)
}

async function uploadFile(file: File) {
  const formData = new FormData()
  formData.append('file', file)
  return await req.post('api/files/upload', formData)
}

async function importFile(file: File, type: string) {
  const formData = new FormData()
  formData.append('file', file)
  formData.append('name', file.name)
  formData.append('shimoType', type)

  return await req.post('api/files/import', formData)
}

async function getExportProgress(id: string) {
  return await req.post(`api/files/export/progress?taskId=${id}`, { taskId: id })
}

async function getImportProgress(id: string, fileId: string) {
  return await req.post(`api/files/import/progress?taskId=${id}&fileId=${fileId}`, { taskId: id, fileId: fileId })
}

async function getImportByUrlProgress(id: string, fileId: string) {
  return await req.post(`api/files/import_by_url/progress?taskId=${id}&fileId=${fileId}`, { taskId: id, fileId: fileId })
}

async function importFileByUrl(fileUrl: string, fileName: string, type: string) {
  const formData = new FormData()
  formData.append('fileUrl', fileUrl)
  formData.append('fileName', fileName)
  formData.append('shimoType', type)

  return await req.post('api/files/import_by_url', formData)
}

async function duplicateFile(id: string) {
  return await req.post(`api/files/${id}/duplicate`)
}

async function getAllCollaborators(id: string) {
  return await req.get(`api/files/${id}/collaborators?all=1`)
}

async function saveCollaborators(id: string, collab: any) {
  return await req.patch(`api/files/${id}/collaborators`, collab)
}

async function exportFile(id: string, type: string) {
  return await req.post(`api/files/${id}/export?type=${type}`)
}

async function exportTableFile(id: string) {
  return await req.post(`api/files/${id}/export/table-sheets`)
}

async function batchDeleteFile(ids: string[]) {
  return await req.delete(`api/files/batch/delete`, {
    data: ids
  })
}

function getPlainTextAPIUrl(id: string) {
  let prefix = process.env.NODE_ENV === 'development' ? process.env.PROXY_PATH : getUrlPrefix()
  return `${prefix}api/files/${id}/download-plain-text?accessToken=${getToken()}`
}

async function getWebInspect() {
  return await req.get('/api/internal/page')
}
