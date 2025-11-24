import styles from './FileList.module.less'
import { Table, Dropdown, Space, Button, Popover, Popconfirm, message, Tag, Modal } from 'antd';
import type { TableProps, MenuProps } from 'antd';
import { DeleteOutlined, EllipsisOutlined, LoadingOutlined, CheckCircleOutlined, MinusCircleOutlined, UserOutlined } from '@ant-design/icons';
import { fileConstants } from '@/constants';
import { useEffect, useState } from 'react';
import { history, useModel } from "umi";
import { RevisionModal } from "./component/RevisionModal";
import { FileCollaborators } from "./component/FileCollaborators";
import { HistoriesModal } from "./component/HistoriesModal";
import { TitleSolt } from "./component/TitleSolt";
import { DataType, fileService } from "@/services/file.service"
import { FileType } from 'shimo-js-sdk'
type TableRowSelection<T extends object> = TableProps<T>['rowSelection'];
const tip = `仅支持以下类型：${fileConstants.PREVIEWABLE_EXTNAMES.join(', ')}`

export function FileList(props: { className: string; }) {
  // File list row selection
  const [rowSelection, setRowSelection] = useState<TableRowSelection<DataType> | undefined>(undefined);
  const [selectedRowKeys, setSelectedRowKeys] = useState<string[]>([]);
  // Confirm dialog helper
  const { confirm } = Modal;
  // Modal states
  const [openRevisionModal, setOpenRevisionModal] = useState(false)
  const [openHistoryModal, setOpenHistoryModal] = useState(false)
  // Currently selected file ID
  const [currentFileId, setCurrentFileId] = useState('')

  const [loading, setLoading] = useState(false)
  // Delete state
  const [removeState, setRemoveState] = useState({} as any)

  const [actionLoding, setActionLoading] = useState({} as any)
  // Collaborator modal state
  const [selected, setSelected] = useState({
    editCollaborators: false,
    file: {} as DataType
  })

  const { createState, importState, uploadState, importUrlState, tableItems, setTableItems } = useModel('file')

  // Fetch the file list
  useEffect(() => {
    getFiles()
  }, [])

  useEffect(() => {
    if (
      uploadState === fileConstants.UPLOAD_FILE_SUCCESS ||
      createState === fileConstants.CREATE_SHIMO_FILE_SUCCESS ||
      importState === fileConstants.IMPORT_FILE_SUCCESS ||
      importUrlState == fileConstants.IMPORT_URL_SUCCESS
    ) {
      getFiles()
    }
  }, [uploadState, createState, importState, importUrlState])

  // Fetch the file list
  function getFiles() {
    setLoading(true)
    fileService.getFiles().then(files => {
      if (files && files.status === 200) {
        setTableItems(files.data)
      } else {
        setTableItems([])
      }
      setLoading(false)
    })
  }
  // Derive the display file type
  function formatFileType(type: string) {
    if (!type) return '其他文件'
    switch (type) {
      case fileConstants.TYPE_DOCUMENT:
        return '轻文档'
      case fileConstants.TYPE_DOCUMENT_PRO:
        return '传统文档'
      case fileConstants.TYPE_SPREADSHEET:
        return '表格'
      case fileConstants.TYPE_PRESENTATION:
        return '幻灯片'
      case fileConstants.TYPE_TABLE:
        return '应用表格'
      case fileConstants.TYPE_FORM:
        return '表单'
      case fileConstants.TYPE_BOARD:
        return '白板'
      case fileConstants.TYPE_MINDMAP:
        return '思维导图'
      case fileConstants.TYPE_FLOWCHART:
        return '流程图'
    }
    const [main, sub] = type.split('/')
    switch (main) {
      case 'audio':
        return '音频'
      case 'video':
        return '视频'
      case 'image':
        return '图片'
    }
    const [_main] = type.split(';')
    return (fileConstants.TYPES as any)[type] || (fileConstants.TYPES as any)[sub] || (fileConstants.TYPES as any)[_main] || '其他文件'
  }
  // Build the export option list
  function getExportButtons(file: DataType): { key: string, label: string }[] {
    let buttons: { text: string; type: string; }[]
    switch (file.shimoType) {
      case FileType.Document:
        buttons = [
          { text: 'DOCX', type: 'docx' },
          { text: 'Markdown', type: 'md' },
          { text: 'JPEG', type: 'jpg' },
          { text: 'PDF', type: 'pdf' }
        ]
        break
      case FileType.DocumentPro:
        buttons = [
          { text: 'DOCX', type: 'docx' },
          { text: 'PDF', type: 'pdf' },
          { text: 'WPS', type: 'wps' }
        ]
        break
      case FileType.Presentation:
        buttons = [
          { text: 'PPTX', type: 'pptx' },
          { text: 'PDF', type: 'pdf' }
        ]
        break
      case FileType.Spreadsheet:
        buttons = [{ text: 'XLSX', type: 'xlsx' }]
        break
      // SDK 0730: application tables do not yet support import/export
      case FileType.Table:
        buttons = [{ text: 'XLSX', type: 'xlsx' }]
        break
      case fileConstants.TYPE_MINDMAP:
        buttons = [{ text: '图片', type: 'jpg' }, { text: 'Xmind', type: 'xmind'}]
        break
      default:
        buttons = []
    }
    return buttons.map((item, i) => {
      return {
        key: item.type,
        label: item.text
      }
    })
  }
  // Action items
  const getItems = (file: DataType): MenuProps['items'] => {
    let _items: MenuProps['items'] = [
      { key: 'preview', label: '预览' },
      { key: 'createCopy', label: '创建副本' },
      { key: 'collaborator', label: '协作者' },
    ]
    let exportButtons = getExportButtons(file)
    if (exportButtons && exportButtons.length > 0) {
      _items.push({ key: 'export', label: '导 出', children: exportButtons })
    }
    if ([FileType.Document, FileType.DocumentPro, FileType.Spreadsheet].includes(file.shimoType as FileType)) {
      _items.push({ key: 'exportText', label: '下载纯文本' })
    }
    return _items
  }

  // Open the revision modal
  function handleOpenRevisionModal(file: DataType) {
    setOpenRevisionModal(true)
    setCurrentFileId(file.id)
  }
  // Open the history modal
  function handleOpenHistoryModal(file: DataType) {
    setOpenHistoryModal(true)
    setCurrentFileId(file.id)
  }
  // Delete a file
  function handleRemove(fileId: string) {
    setRemoveState({ ...removeState, [fileId]: fileConstants.REMOVE_FILE_REQUEST })
    fileService.removeFile(fileId).then(res => {
      setRemoveState({ [fileId]: fileConstants.REMOVE_FILE_SUCCESS })
      setTableItems(Array.isArray(tableItems)
        ? tableItems.filter((f: { id: number | string }) => f.id !== fileId)
        : [])
      message.success('删除成功');
    }).catch(() => {
      delete removeState[fileId]
      setRemoveState({ ...removeState })
    })
  }
  /**
   * Handle preview button rendering
  */
  function PreviewButton({ file }: { file: DataType }) {
    if (file.isShimoFile === 0) {
      const index = file.name.lastIndexOf(".");
      const ext = file.name.substring(index);
      const previewable = fileConstants.PREVIEWABLE_EXTNAMES.includes(ext)
      const button = (
        <Button
          onClick={() => previewable && history.push(`/preview/${file.id}`)}
          className={previewable ? '' : styles.unpreviewable}
        >
          预览
        </Button>
      )
      if (!previewable) {
        return <Popover content={tip}>
          {button}
        </Popover>
      } else {
        return button
      }
    }
  }
  // Table column definitions
  const columns: TableProps<DataType>['columns'] = [
    {
      title: '标题',
      dataIndex: 'name',
      key: 'name',
      render: (_, record) => (
        <TitleSolt key={record.id} file={record} handleOpenHistoryModal={handleOpenHistoryModal} handleOpenRevisionModal={handleOpenRevisionModal}></TitleSolt>
      ),
      width: '350px'
    },
    {
      title: '类型',
      dataIndex: 'shimoType',
      key: 'shimoType',
      width: '100px',
      render: (_, record) => (<span>{formatFileType(record.shimoType ? record.shimoType : record.type)}</span>)
    },
    {
      title: '角色',
      dataIndex: 'role',
      key: 'role',
      width: '100px',
      render: (text, record) => (
        <Popover content={
          <Space size={[6, 12]} wrap>
            {((process.env.USE_NEW_FILE_PERMISSION == "true") ? fileConstants.NEW_PERMISSIONS : fileConstants.PERMISSIONS).map((m: { label: string, value: string }, i: number) => {
              return (
                <Tag key={i}
                  color={record.permissions[m.value as keyof typeof record.permissions] ? 'success' : 'default'}
                  icon={record.permissions[m.value as keyof typeof record.permissions] ? <CheckCircleOutlined /> : <MinusCircleOutlined />}
                >
                  {m.label.replace(/可/g, '')}
                </Tag>
              )
            })}
          </Space>
        }>
          <Tag
            color='default'
            style={{ borderStyle: text == 'owner' ? '' : 'dashed', padding: '5px' }}
            icon={text == 'owner' ? <UserOutlined /> : null}
          >
            {text == 'collaborator' ? '协作者' : text == 'owner' ? '管理者' : '其他'}
          </Tag>
        </Popover>
      )
    },
    {
      title: '',
      key: 'action',
      align: 'right',
      render: (_, record) => (
        <Space.Compact size="middle">
          {record.isShimoFile > 0 ? (
            <Button onClick={() => history.push(`/shimo-files/${record.id}`)}>打开</Button>
          ) : (
            <PreviewButton file={record} />
          )}
          {record.isShimoFile > 0 && (
            <Dropdown menu={{ items: getItems(record), onClick: (e) => actionClick(e, record) }} trigger={['click']} disabled={actionLoding[record.id]}>
              <Button onClick={(e) => e.preventDefault()} >
                {actionLoding[record.id] ? <LoadingOutlined /> : <EllipsisOutlined />}
              </Button>
            </Dropdown>
          )}
          <Popconfirm
            title="删除"
            description="确认删除？"
            onConfirm={() => handleRemove(record.id)}
            okText="确认"
            cancelText="取消"
          >
            <Button disabled={removeState[record.id] != null || record.permissions?.manageable !== true}>
              {removeState[record.id] != null ? <LoadingOutlined /> : <DeleteOutlined />}
            </Button>
          </Popconfirm>
        </Space.Compact>
      ),
      width: '150px'
    },
  ];
  /**
   * Handle clicks on the action column dropdown
  */
  const actionClick = (e: any, file: DataType) => {
    console.log(e, file);
    // Export by type
    if (e.keyPath.length === 2) {
      setActionLoading({ [file.id]: true })
      if (file.shimoType === FileType.Table) {
        fileService.exportTableFile(file.id).then((res) => {
          const a = document.createElement('a')
          a.download = ''
          a.href = res.data.downloadUrl
          a.click()
        }).finally(() => {
          setActionLoading({ [file.id]: false })
        })
      } else {
        fileService.exportFile(file.id, e.key).then(async (res) => {
          do {
            try {
              if (!(res.data && res.data.data && res.data.data.taskId)) {
                console.log("taskid获取失败");
                break
              }
              let resp: any = null
              await fileService.getExportProgress(res.data.data.taskId).then((r) => {
                resp = r
              }).catch((err) => { console.log(err) })
              if (!resp) {
                console.log("导出进度获取失败", resp.data);
                break
              }
              if (resp.data?.status !== 0) {
                throw new Error(resp.data?.message)
              }
              console.log("import processing...");

              if (resp.data.data.progress === 100) {
                const a = document.createElement('a')
                a.download = ''
                a.href = resp.data.data.downloadUrl
                a.click()
                break
              }
            } catch (e) {
              console.log(e + "导出进度获取失败");
              break
            }
            await new Promise((resolve) => setTimeout(resolve, 1000))
          } while (true)
        }).finally(() => {
          setActionLoading({ [file.id]: false })
        })
      }
    } else {
      switch (e.key) {
        // Create a copy
        case 'createCopy':
          setActionLoading({ [file.id]: true })
          fileService.duplicateFile(file.id).then((res) => {
            setTableItems([res.data, ...tableItems])
          }).finally(() => {
            setActionLoading({ [file.id]: false })
          })
          break;
        // Manage collaborators
        case 'collaborator':
          setSelected({ editCollaborators: true, file: file })
          break;
        case 'exportText':
          window.location.href = fileService.getPlainTextAPIUrl(file.id)
          break;
        case 'preview':
          history.push(`/preview/${file.id}`)
          break;
      }
    }
  };
  /**
   * Persist collaborator management changes
  */
  function editCollaborators(fileId: string, data: any) {
    setActionLoading({ [fileId]: true })
    setSelected({ editCollaborators: false, file: {} as DataType })
    fileService.saveCollaborators(fileId, data).then(() => {
      message.success('保存成功')
    }).finally(() => {
      setActionLoading({ [fileId]: false })
    })
  }
  // Toggle batch selection
  function handleRowSelectionChange() {
    setSelectedRowKeys([])
    setRowSelection(rowSelection == undefined ? {
      onChange: (selectedRowKeys: React.Key[], selectedRows: DataType[]) => {
        setSelectedRowKeys(selectedRowKeys.map((v) => String(v)))
      },
      // Only managers can delete
      getCheckboxProps: (record: DataType) => ({
        disabled: record.permissions?.manageable !== true
      }),
    } : undefined)
  }
  // Batch delete
  function handleBatchDelete() {
    console.log(selectedRowKeys)
    if (!selectedRowKeys.length) {
      message.warning("请选择要删除的文件")
      return
    }
    confirm({
      title: '确认删除',
      content: `确认删除选中的 ${selectedRowKeys.length} 个文件`,
      onOk() {
        fileService.batchDeleteFile(selectedRowKeys).then(() => {
          setTableItems(Array.isArray(tableItems)
            ? tableItems.filter((f: { id: string }) => !selectedRowKeys.includes(f.id))
            : [])
          setSelectedRowKeys([])
          message.success('删除成功');
        })
      }
    })
  }

  return (
    <>
      {
        rowSelection == undefined ? (
          <Button onClick={handleRowSelectionChange} type='link' className={styles.batchDelete}>批量删除</Button>
        ) : (
          <div className={styles.batchDelete}>
            <Button onClick={handleBatchDelete} type='primary' style={{ marginRight: "10px" }}>删除</Button>
            <Button onClick={handleRowSelectionChange}>取消</Button>
          </div>
        )
      }
      <Table rowSelection={rowSelection} columns={columns} dataSource={tableItems} rowKey="id" pagination={false} className={styles.TableProps} loading={loading} scroll={{ y: 525 }} />
      <RevisionModal open={openRevisionModal} onClose={() => setOpenRevisionModal(false)} fileId={currentFileId}></RevisionModal>
      <HistoriesModal open={openHistoryModal} onClose={() => setOpenHistoryModal(false)} fileId={currentFileId}></HistoriesModal>
      <FileCollaborators open={selected.editCollaborators} onClose={() => setSelected({ editCollaborators: false, file: {} as DataType })} onSave={editCollaborators} file={selected.file}></FileCollaborators>
    </>
  )
}
