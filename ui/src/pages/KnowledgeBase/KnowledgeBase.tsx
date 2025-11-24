import { useState, useEffect } from 'react';
import { Button, Input, Modal, Table, Space, Tag, Popconfirm, message, Tabs } from 'antd';
import { PlusOutlined, DeleteOutlined, FileOutlined, FolderOutlined } from '@ant-design/icons';
import { history } from 'umi';
import './KnowledgeBase.module.less';
import { useModel } from 'umi'
import styles from './KnowledgeBase.module.less';

const { Search } = Input;

export default function KnowledgeBase() {
  const [isCreateModalVisible, setIsCreateModalVisible] = useState(false);
  const [newKnowledgeBaseGuid, setNewKnowledgeBaseGuid] = useState('');
  const [searchText, setSearchText] = useState('');

  const {
    knowledgeBases,
    loading,
    getKnowledgeBases,
    createKnowledgeBase,
    deleteKnowledgeBase,
  } = useModel("knowledgeBase");

  useEffect(() => {
    getKnowledgeBases();
  }, [getKnowledgeBases]);

  const handleCreateKnowledgeBase = async () => {
    if (!newKnowledgeBaseGuid.trim()) {
      message.error('请输入知识库ID');
      return;
    }

    const success = await createKnowledgeBase(newKnowledgeBaseGuid.trim());
    if (success) {
      setIsCreateModalVisible(false);
      setNewKnowledgeBaseGuid('');
    }
  };

  const handleDeleteKnowledgeBase = async (guid: string) => {
    await deleteKnowledgeBase(guid);
    getKnowledgeBases();
  };

  const handleKnowledgeBaseClick = (guid: string) => {
    history.push(`/knowledge-base/${guid}`);
  };

  const columns = [
    {
      title: '知识库',
      dataIndex: 'name',
      key: 'name',
      render: (name: string, record: any) => (
        <Button
          type="link"
          onClick={() => handleKnowledgeBaseClick(record.guid)}
          style={{ padding: 0, height: 'auto' }}
        >
          <FolderOutlined style={{ marginRight: 8 }} />
          {name}
        </Button>
      ),
    },
    {
      title: '文件数量',
      dataIndex: 'fileCount',
      key: 'fileCount',
      render: (count: number) => (
        <Tag color="blue">
          <FileOutlined style={{ marginRight: 4 }} />
          {count}
        </Tag>
      ),
    },
    {
      title: '创建时间',
      dataIndex: 'createAt',
      key: 'createAt',
      render: (timestamp: number) => new Date(timestamp * 1000).toLocaleString('zh-CN'),
    },
    {
      title: '操作',
      key: 'action',
      render: (_: any, record: any) => (
        <Space>
          <Button
            type="link"
            onClick={() => handleKnowledgeBaseClick(record.guid)}
          >
            查看
          </Button>
          <Popconfirm
            title="确定要删除这个知识库吗？"
            onConfirm={() => handleDeleteKnowledgeBase(record.guid)}
            okText="确定"
            cancelText="取消"
          >
            <Button type="link" danger>
              <DeleteOutlined />
              删除
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div className="knowledge-base-page">
      <div className={styles.knowledgeBaseHeader}>
        <div className={styles.headerTitle}>知识库管理</div>
        <div className={styles.headerExtra}>
          <Button
            type="primary"
            icon={<PlusOutlined />}
            onClick={() => setIsCreateModalVisible(true)}
          >
            创建知识库
          </Button>
        </div>
      </div>

      <div className={styles.knowledgeBaseContent}>
        <div className={styles.searchSection}>
          <Search
            placeholder="搜索知识库ID"
            allowClear
            value={searchText}
            onChange={(e) => setSearchText(e.target.value)}
            style={{ width: 300 }}
          />
        </div>

        <Table
          columns={columns}
          dataSource={knowledgeBases}
          rowKey="guid"
          loading={loading}
          pagination={{
            pageSize: 10,
            showSizeChanger: true,
            showQuickJumper: true,
            showTotal: (total) => `共 ${total} 条记录`,
          }}
        />
      </div>

      <Modal
        title="创建知识库"
        open={isCreateModalVisible}
        onOk={handleCreateKnowledgeBase}
        onCancel={() => {
          setIsCreateModalVisible(false);
          setNewKnowledgeBaseGuid('');
        }}
        okText="创建"
        cancelText="取消"
      >
        <div style={{ marginBottom: 16 }}>
          <label>知识库名称:</label>
          <Input
            placeholder="请输入知识库名称"
            value={newKnowledgeBaseGuid}
            onChange={(e) => setNewKnowledgeBaseGuid(e.target.value)}
            style={{ marginTop: 8 }}
          />
        </div>
      </Modal>
    </div>
  );
}