import req from '@/utils/axios';

export const knowledgeBaseService = {
  getKnowledgeBases,
  createKnowledgeBase,
  deleteKnowledgeBase,
  importFileToKnowledgeBase,
  importFileToKnowledgeBaseV2,
  getImportProgressV2,
  getKnowledgeBaseFiles,
  deleteFileFromKnowledgeBase,
  getAiAssets,
  getKnowledgeBase,
}

export interface KnowledgeBaseInfo {
  guid: string;
  fileCount: number;
  createAt: number;
}

export interface FileInfo {
  id: number;
  guid: string;
  name: string;
  type: string;
  shimoType: string;
  createAt: number;
  createBy: number;
}

async function getKnowledgeBases() {
  return await req.get('/api/knowledge/list');
}

async function createKnowledgeBase(name: string) {
  return await req.post('/api/knowledge/create', { name });
}

async function deleteKnowledgeBase(guid: string) {
  return await req.delete(`/api/knowledge/${guid}/delete`);
}

async function importFileToKnowledgeBase(params: {
  knowledgeBaseGuid: string;
  importType: string;
  fileGuid?: string;
  fileType: string;
  downloadUrl?: string;
}) {
  return await req.post('/api/knowledge/import', params);
}

async function getKnowledgeBaseFiles(guid: string) {
  return await req.get(`/api/knowledge/${guid}/files`);
}

async function deleteFileFromKnowledgeBase(guid: string, fileGuid: string) {
  return await req.delete(`/api/knowledge/${guid}/delete/${fileGuid}`);
}

async function getAiAssets() {
  return await req.get('/api/knowledge/ai-assets');
}

async function getKnowledgeBase(guid: string) {
  return await req.get(`/api/knowledge/${guid}`);
}

async function importFileToKnowledgeBaseV2(params: {
  knowledgeBaseGuid: string;
  importType: string;
  fileGuid?: string;
  fileType: string;
  downloadUrl?: string;
}) {
  return await req.post('/api/knowledge/v2/import', params);
}

async function getImportProgressV2(taskId: string, fileGuid: string, knowledgeBaseGuid: string) {
  return await req.post('/api/knowledge/v2/import/progress', { taskId, fileGuid, knowledgeBaseGuid });
}
