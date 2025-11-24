import { DataType } from "@/services/file.service";
import { Modal, Spin, Checkbox, Table, Button } from "antd";
import type { TableProps } from 'antd';
import { useState, useEffect } from "react";
import styles from '../FileList.module.less'
import { fileService } from "@/services/file.service";
import { fileConstants } from "@/constants";

export function FileCollaborators({ open, onClose, file, onSave }: { open: boolean, onClose: () => void, file: DataType, onSave: (fileId: string, data: any) => void }) {
  const [tableData, setTableData] = useState<any[]>([])
  const [loading, setLoading] = useState(false)

  const [permissions, setPermissions] = useState({} as any)
  const [checked, setChecked] = useState({} as any)

  useEffect(() => {
    setTableData([])
    setChecked({})
    setPermissions({})
    if (!file || !file.id) {
      return
    }
    setLoading(true)
    fileService.getAllCollaborators(file.id).then(res => {
      setTableData(res.data)
      const permissions: any = {}
      const checkds: any = {}
      for (const user of res.data) {
        permissions[user.id] = user.permissions[0] ? user.permissions[0] : {}
        if (user.permissions[0]) {
          checkds[user.id] = user.permissions[0]
        }
      }
      setPermissions(permissions)
      setChecked(checkds)
    }).finally(() => {
      setLoading(false)
    })
  }, [file])

  const options = (process.env.USE_NEW_FILE_PERMISSION == "true") ? fileConstants.NEW_PERMISSIONS : fileConstants.PERMISSIONS;

  function defaultNewPermissions(value = true) {
    return {
      editable: value,
      readable: value,
      copyable: value,
      commentable: value,
      exportable: value,
      manageable: value,
      cutable: value,
      imageDownloadable: value,
      copyablePasteClipboard: value,
      attachmentCopyable: value,
      attachmentPreviewable: value,
      attachmentDownloadable: value
    }
  }

  function defaultPermissions(value = true) {
    return {
      editable: value,
      readable: value,
      copyable: value,
      commentable: value,
      exportable: value,
      manageable: value
    }
  }

  const columns: TableProps['columns'] = [
    {
      title: '名字',
      dataIndex: 'name',
      key: 'name',
      width: '90px'
    },
    {
      title: '权限',
      dataIndex: 'permissions',
      key: 'permissions',
      render: (_, record) => (
        <>
          {options.map((m: any, i: number) => {
            return (
              <Checkbox key={m.value} value={m.value} onChange={() => onChangeValue(record.id, m.value, !permissions[record.id][m.value])}
                checked={permissions[record.id] ? permissions[record.id][m.value] : false}
                disabled={(m.value !== 'manageable' && permissions[record.id]?.manageable) || record.isCreator}
              >{m.label}</Checkbox>
            )
          })}
        </>
      ),
      width: '350px'
    },
  ]

  function onChangeValue(userId: string, name: string, value: boolean) {
    let perm = permissions[userId] || {}
    perm[name] = value

    if (!checked[userId]) {
      checked[userId] = {}
    }

    checked[userId][name] = value

    if (name !== 'readable') {
      perm.readable = checked[userId].readable = true
      // Editing implies copy permission
      if (name === 'editable' && value) {
        perm.copyable = checked[userId].copyable = true
      } else if (name === 'manageable') {
        // Manager permissions grant every capability
        perm = checked[userId] = (process.env.USE_NEW_FILE_PERMISSION == "true") ? defaultNewPermissions(value) : defaultPermissions(value)
      } else if (name === 'copyable' && !value) {
        // Without copy permission, disallow external paste
        perm.copyablePasteClipboard = checked[userId].copyablePasteClipboard = false
        // Without copy permission, disallow editing
        perm.editable = checked[userId].editable = false
      } else if (name === 'copyablePasteClipboard' && value) {
        // Allow copy when external paste is enabled
        perm.copyable = checked[userId].copyable = true
      } else if (name === 'attachmentDownloadable' && value) {
        // Downloading attachments also enables preview
        perm.attachmentPreviewable = checked[userId].attachmentPreviewable = true
      } else if (name === 'attachmentPreviewable' && !value) {
        // If attachment preview is disabled, also disable download
        perm.attachmentDownloadable = checked[userId].attachmentDownloadable = false
      }
    } else if (value === false) {
      perm = checked[userId] = (process.env.USE_NEW_FILE_PERMISSION == "true") ? defaultNewPermissions(false) : defaultPermissions(false)
    }

    setPermissions({
      ...permissions,
      [userId]: perm
    })
    setChecked(checked)
  }

  return (
    <Modal
      open={open}
      title="管理协作者"
      footer={[<Button onClick={onClose} key='close'>关闭</Button>, <Button onClick={() => onSave(file.id, checked)} key='ok'>确认</Button>]}
      onCancel={onClose}
      width='60vw'
    >
      <Table columns={columns} dataSource={tableData} rowKey="id" pagination={false} className={styles.CollaboratorProps} loading={loading} scroll={{ y: 480 }} />
    </Modal>
  )
}