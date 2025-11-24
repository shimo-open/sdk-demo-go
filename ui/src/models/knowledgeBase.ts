import { useState, useCallback } from 'react';
import { message } from 'antd';
import { knowledgeBaseService, KnowledgeBaseInfo, FileInfo } from '@/services/knowledgeBase.service';
import { knowledgeBaseConstants } from '@/constants/knowledgeBase.constants';
import { fileService, DataType } from '@/services/file.service';

export default function useKnowledgeBase() {
  const [knowledgeBases, setKnowledgeBases] = useState<KnowledgeBaseInfo[]>([]);
  const [knowledgeBaseFiles, setKnowledgeBaseFiles] = useState<FileInfo[]>([]);
  const [loading, setLoading] = useState(false);
  const [currentKnowledgeBase, setCurrentKnowledgeBase] = useState<string>('');

  // Fetch the list of knowledge bases
  const getKnowledgeBases = useCallback(async () => {
    setLoading(true);
    try {
      const response = await knowledgeBaseService.getKnowledgeBases();
      if (response.data?.data) {
        setKnowledgeBases(response.data.data);
      } else {
        setKnowledgeBases([]);
      }
    } catch (error) {
      message.error('获取知识库列表失败');
      console.error('Failed to get knowledge bases:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  // Create a knowledge base
  const createKnowledgeBase = useCallback(async (name: string) => {
    try {
      const response = await knowledgeBaseService.createKnowledgeBase(name);
      if (response.data?.message === 'Knowledge base created successfully') {
        message.success('知识库创建成功');
        await getKnowledgeBases(); // Refresh the list
        return true;
      }
    } catch (error) {
      message.error('创建知识库失败');
      console.error('Failed to create knowledge base:', error);
    }
    return false;
  }, [getKnowledgeBases]);

  // Delete a knowledge base
  const deleteKnowledgeBase = useCallback(async (guid: string) => {
    try {
      const response = await knowledgeBaseService.deleteKnowledgeBase(guid);
      if (response.data?.message === 'Knowledge base deleted successfully') {
        message.success('知识库删除成功');
        await getKnowledgeBases(); // Refresh the list
        return true;
      }
    } catch (error) {
      message.error('删除知识库失败');
      console.error('Failed to delete knowledge base:', error);
    }
    return false;
  }, [getKnowledgeBases]);

  // Fetch the files inside a knowledge base
  const getKnowledgeBaseFiles = useCallback(async (guid: string) => {
    setLoading(true);
    setCurrentKnowledgeBase(guid);
    try {
      const response = await knowledgeBaseService.getKnowledgeBaseFiles(guid);
      if (response.data?.data) {
        setKnowledgeBaseFiles(response.data.data);
      } else {
        setKnowledgeBaseFiles([]);
      }
    } catch (error) {
      message.error('获取知识库文件列表失败');
      console.error('Failed to get knowledge base files:', error);
    } finally {
      setLoading(false);
    }
  }, []);

  // Import a file into the knowledge base
  const importFileToKnowledgeBase = useCallback(async (params: {
    knowledgeBaseGuid: string;
    importType: string;
    fileGuid?: string;
    fileType: string;
    downloadUrl?: string;
  }) => {
    try {
      const response = await knowledgeBaseService.importFileToKnowledgeBase(params);
      if (response.data?.message === 'File imported successfully') {
        message.success('文件导入成功');
        await getKnowledgeBaseFiles(params.knowledgeBaseGuid); // Refresh the file list
        return true;
      }
    } catch (error: any) {
      message.error('文件导入失败' + error?.data?.error);
      console.error('Failed to import file to knowledge base:', error);
    }
    return false;
  }, [getKnowledgeBaseFiles]);

  // Remove a file from the knowledge base
  const deleteFileFromKnowledgeBase = useCallback(async (guid: string, fileGuid: string) => {
    try {
      const response = await knowledgeBaseService.deleteFileFromKnowledgeBase(guid, fileGuid);
      if (response.data?.message === 'File removed from knowledge base successfully') {
        message.success('文件移除成功');
        await getKnowledgeBaseFiles(guid); // Refresh the file list
        return true;
      }
    } catch (error) {
      message.error('文件移除失败');
      console.error('Failed to delete file from knowledge base:', error);
    }
    return false;
  }, [getKnowledgeBaseFiles]);

  // Fetch all files (used when adding to the knowledge base)
  const getAllFiles = useCallback(async () => {
    try {
      const response = await fileService.getFiles();
      return response.data || [];
    } catch (error) {
      message.error('获取文件列表失败');
      console.error('Failed to get files:', error);
      return [];
    }
  }, []);

  return {
    knowledgeBases,
    knowledgeBaseFiles,
    loading,
    currentKnowledgeBase,
    getKnowledgeBases,
    createKnowledgeBase,
    deleteKnowledgeBase,
    getKnowledgeBaseFiles,
    importFileToKnowledgeBase,
    deleteFileFromKnowledgeBase,
    getAllFiles,
  };
}
