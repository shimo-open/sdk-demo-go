import React, { useState, useEffect } from 'react';
import { Button, Table, Space, Tag, Popconfirm, message, Modal, Select, Input } from 'antd';
import { ArrowLeftOutlined, PlusOutlined, DeleteOutlined, FileOutlined, RobotOutlined } from '@ant-design/icons';
import useKnowledgeBase from '@/models/knowledgeBase';
import { knowledgeBaseConstants } from '@/constants/knowledgeBase.constants';
import { fileConstants } from '@/constants/file.constants';
import { history, useParams } from 'umi';
import styles from './KnowledgeBase.module.less';
import { knowledgeBaseService } from '@/services/knowledgeBase.service';

const { Option } = Select;
const { TextArea } = Input;

export default function KnowledgeBaseDetail() {
  // Read knowledgeBaseGuid from the route parameters
  const { knowledgeBaseGuid } = useParams<{ knowledgeBaseGuid: string }>();

  const [isAddFileModalVisible, setIsAddFileModalVisible] = useState(false);
  const [isAddUrlModalVisible, setIsAddUrlModalVisible] = useState(false);
  const [selectedFiles, setSelectedFiles] = useState("");
  const [fileType, setFileType] = useState<string>('');
  const [downloadUrl, setDownloadUrl] = useState('');
  const [addFileLoading, setAddFileLoading] = useState(false);

  const [currentKnowledgeBase, setCurrentKnowledgeBase] = useState<any>(null);

  const {
    knowledgeBaseFiles,
    loading,
    getKnowledgeBaseFiles,
    importFileToKnowledgeBase,
    deleteFileFromKnowledgeBase,
    getAllFiles,
  } = useKnowledgeBase();

  const [allFiles, setAllFiles] = useState<any[]>([]);

  useEffect(() => {
    if (knowledgeBaseGuid) {
      knowledgeBaseService.getKnowledgeBase(knowledgeBaseGuid).then(res => {
        setCurrentKnowledgeBase(res.data);
      });
      getKnowledgeBaseFiles(knowledgeBaseGuid);
      loadAllFiles();
    }
  }, [knowledgeBaseGuid, getKnowledgeBaseFiles]);

  const loadAllFiles = async () => {
    const files = await getAllFiles();
    setAllFiles(files.filter((file: any) => file.isShimoFile));
  };

  const handleAddFile = async () => {
    if (selectedFiles === "") {
      message.error('请选择要添加的文件');
      return;
    }
    setAddFileLoading(true);
    let fileType = '';
    allFiles.some(file => {
      if (file.id === selectedFiles) {
        fileType = file.shimoType;
        return true;
      }
    })

    try {
      // Invoke the v2 import API
      const response = await knowledgeBaseService.importFileToKnowledgeBaseV2({
        knowledgeBaseGuid: knowledgeBaseGuid as string,
        importType: 'file',
        fileGuid: selectedFiles,
        fileType,
      });

      if (response.data?.taskId) {
        const taskId = response.data.taskId;
        const fileGuid = response.data.fileGuid;
        const knowledgeBaseGuid = response.data.knowledgeBaseGuid;
        message.loading('文件导入中，请稍候...', 0);

        // Start polling progress
        const pollProgress = async () => {
          try {
            const progressResponse = await knowledgeBaseService.getImportProgressV2(taskId, fileGuid, knowledgeBaseGuid);
            const { status, progress, message: progressMessage } = progressResponse.data;

            if (status === 'completed' && progress === 100) {
              message.destroy();
              message.success('文件导入成功！');
              setIsAddFileModalVisible(false);
              setSelectedFiles("");
              setFileType('');
              setAddFileLoading(false);
              // Refresh the file list
              getKnowledgeBaseFiles(knowledgeBaseGuid as string);
              return;
            } else if (status === 'failed') {
              message.destroy();
              message.error(`导入失败: ${progressMessage || '未知错误'}`);
              setAddFileLoading(false);
              return;
            } else {
              // Continue polling
              setTimeout(pollProgress, 1000);
            }
          } catch (error) {
            message.destroy();
            message.error('查询导入进度失败');
            console.error('Progress polling error:', error);
            setAddFileLoading(false);
          }
        };

        // Kick off polling
        setTimeout(pollProgress, 1000);
      } else {
        message.error('启动导入任务失败');
        setAddFileLoading(false);
      }
    } catch (error) {
      message.error('导入文件失败');
      console.error('Import error:', error);
      setAddFileLoading(false);
    }
  };

  const handleAddUrl = async () => {
    if (!downloadUrl.trim() || !fileType) {
      message.error('请填写完整信息');
      return;
    }
    setAddFileLoading(true);

    try {
      // Invoke the v2 import API
      const response = await knowledgeBaseService.importFileToKnowledgeBaseV2({
        knowledgeBaseGuid: knowledgeBaseGuid as string,
        importType: 'url',
        fileType,
        downloadUrl: downloadUrl.trim(),
      });

      if (response.data?.taskId) {
        const taskId = response.data.taskId;
        const fileGuid = response.data.fileGuid;
        const knowledgeBaseGuid = response.data.knowledgeBaseGuid;
        message.loading('文件导入中，请稍候...', 0);

        // Start polling progress
        const pollProgress = async () => {
          try {
            const progressResponse = await knowledgeBaseService.getImportProgressV2(taskId, fileGuid, knowledgeBaseGuid);
            const { status, progress, message: progressMessage } = progressResponse.data;

            if (status === 'completed' && progress === 100) {
              message.destroy();
              message.success('文件导入成功！');
              setIsAddUrlModalVisible(false);
              setAddFileLoading(false);
              setDownloadUrl('');
              setFileType('');
              // Refresh the file list
              getKnowledgeBaseFiles(knowledgeBaseGuid as string);
              return;
            } else if (status === 'failed') {
              message.destroy();
              message.error(`导入失败: ${progressMessage || '未知错误'}`);
              setAddFileLoading(false);
              return;
            } else {
              // Continue polling
              setTimeout(pollProgress, 1000);
            }
          } catch (error) {
            message.destroy();
            message.error('查询导入进度失败');
            setAddFileLoading(false);
            console.error('Progress polling error:', error);
          }
        };

        // Kick off polling
        setTimeout(pollProgress, 1000);
      } else {
        message.error('启动导入任务失败');
      }
    } catch (error) {
      message.error('导入文件失败');
      console.error('Import error:', error);
      setAddFileLoading(false);
    }
  };

  const handleDeleteFile = async (fileGuid: string) => {
    await deleteFileFromKnowledgeBase(knowledgeBaseGuid as string, fileGuid);
  };

  const handleBack = () => {
    history.push('/knowledge-base');
  };

  const handleGoToAI = () => {
    history.push(`/knowledge-base/${knowledgeBaseGuid}/ai`);
  };

  const columns = [
    {
      title: '文件名',
      dataIndex: 'name',
      key: 'name',
      render: (name: string) => (
        <span>
          <FileOutlined style={{ marginRight: 8 }} />
          {name}
        </span>
      ),
    },
    {
      title: '文件类型',
      dataIndex: 'type',
      key: 'type',
      render: (_: any, record: any) => (
        <Tag color="green">
          {(fileConstants.TYPES as any)[record.type] || record.type || record.shimoType}
        </Tag>
      ),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Popconfirm
          title="确定要移除这个文件吗？"
          onConfirm={() => handleDeleteFile(record.guid)}
          okText="确定"
          cancelText="取消"
        >
          <Button type="link" danger icon={<DeleteOutlined />}>
            移除
          </Button>
        </Popconfirm>
      ),
    },
  ];

  return (
    <div className={styles.knowledgeBaseDetail}>
      {/* Title section */}
      <div className={styles.titleSection}>
        <Space>
          <Button onClick={handleBack}>
            返回
          </Button>
          <span style={{ fontSize: 15, fontWeight: 500 }}>{currentKnowledgeBase?.name}</span>
        </Space>
        <Space>
          <Button
            type="primary"
            onClick={() => setIsAddFileModalVisible(true)}
          >
            添加文件
          </Button>
          <Button
            onClick={() => setIsAddUrlModalVisible(true)}
          >
            添加链接
          </Button>
          <Button
            type="primary"
            icon={<RobotOutlined />}
            onClick={handleGoToAI}

          >
            AI 知识库
          </Button>
        </Space>
      </div>

      {/* File list section */}
      <div className={styles.fileListContainer}>
        <Table
          columns={columns}
          dataSource={knowledgeBaseFiles}
          rowKey="guid"
          loading={loading}
          pagination={{
            size: 'small',
            showTotal: (total) => `共 ${total} 条记录`,
            pageSize: 10,
          }}
          scroll={{ y: 400 }}
        />
      </div>

      {/* Add-file modal */}
      <Modal
        title="添加文件到知识库"
        open={isAddFileModalVisible}
        onOk={handleAddFile}
        onCancel={() => {
          setIsAddFileModalVisible(false);
          setSelectedFiles("");
          setFileType('');
        }}
        confirmLoading={addFileLoading}
        okText="添加"
        cancelText="取消"
        width={600}
      >
        <div style={{ marginBottom: 16 }}>
          <label>选择文件:</label>
          <Select
            placeholder="请选择要添加的文件"
            value={selectedFiles}
            onChange={setSelectedFiles}
            style={{ width: '100%', marginTop: 8 }}
            showSearch
            filterOption={(input, option) => {
              const children = React.Children.toArray(option?.children as React.ReactNode).join('');
              return children.toLowerCase().includes(input.toLowerCase());
            }}
          >
            {allFiles.map(file => (
              <Option key={file.id} value={file.id}>
                {file.name} {file.isShimoFile ? `(${file.shimoType})` : `(${file.type})`}
              </Option>
            ))}
          </Select>
        </div>
      </Modal>

      {/* Add-link modal */}
      <Modal
        title="通过链接添加文件"
        open={isAddUrlModalVisible}
        onOk={handleAddUrl}
        onCancel={() => {
          setIsAddUrlModalVisible(false);
          setDownloadUrl('');
          setFileType('');
        }}
        okText="添加"
        cancelText="取消"
        confirmLoading={addFileLoading}
      >
        <div style={{ marginBottom: 16 }}>
          <label>下载链接:</label>
          <TextArea
            placeholder="请输入文件下载链接"
            value={downloadUrl}
            onChange={(e) => setDownloadUrl(e.target.value)}
            rows={3}
            style={{ marginTop: 8 }}
          />
        </div>
        <div style={{ marginBottom: 16 }}>
          <label>文件类型:</label>
          <Select
            placeholder="请选择文件类型"
            value={fileType}
            onChange={setFileType}
            style={{ width: '100%', marginTop: 8 }}
          >
            <Option value={knowledgeBaseConstants.FILE_TYPE_PDF}>PDF</Option>
            <Option value={knowledgeBaseConstants.FILE_TYPE_RTF}>RTF</Option>
          </Select>
        </div>
      </Modal>
    </div>
  );
}