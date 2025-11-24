import { useModel } from "umi";
import { Modal, Spin, List, Avatar, Button } from "antd";
import { fileConstants } from "@/constants";
import { FileZipOutlined } from '@ant-design/icons';
import styles from '../FileList.module.less'
import { useEffect, useState } from "react";
import { fileService } from "@/services/file.service";

export function RevisionModal({ open, onClose, fileId }: { open: boolean, onClose: () => void, fileId: string }) {
  const [revisions, setRevisions] = useState([] as any)
  const [loading, setLoading] = useState('')

  useEffect(() => {
    if (fileId && open) {
      setLoading('loading')
      fileService.getRevisions(fileId).then((res) => {
        (Array.isArray(res.data) || !res.data) ? setRevisions(res.data) : setRevisions([res.data])
        setLoading('loaded')
      }).catch(() => {
        setRevisions([])
        setLoading('error')
      })
    } else {
      setRevisions([])
    }
  }, [fileId, open])

  return (
    <Modal
      open={open}
      title="查看版本列表"
      okText=""
      footer={<Button onClick={onClose}>关闭</Button>}
      onCancel={onClose}
    >
      {loading === 'loading' && (
        <Spin style={{ display: 'block' }} />
      )}
      {revisions?.length > 0 && (
        <List
          dataSource={revisions}
          renderItem={(item: any, index) => (
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar icon={<FileZipOutlined />} />}
                title={<>
                  <a>{item.user.name}</a> 创建了版本
                  <span className={styles.date}>{new Date(item.createdAt).toLocaleString()}</span>
                </>}
                description={item.label}
              />
            </List.Item>
          )}
        />
      )}
      {revisions.length === 0 && loading === 'error' ? (
        <span>获取版本失败</span>
      ) : null}

      {revisions.length === 0 && loading === 'loaded' ? (
        <span>暂无版本</span>
      ) : null}
    </Modal>
  )
}