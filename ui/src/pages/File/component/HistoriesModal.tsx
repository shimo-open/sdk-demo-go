import { useModel } from "umi";
import { Modal, Spin, List, Avatar, Button } from "antd";
import { fileConstants } from "@/constants";
import { HistoryOutlined } from '@ant-design/icons';
import styles from '../FileList.module.less'
import { useState, useEffect } from 'react';
import { fileService } from '@/services/file.service'

export function HistoriesModal({ open, onClose, fileId }: { open: boolean, onClose: () => void, fileId: string }) {
  const unkownUser = '不存在用户'
  const [histories, setHistories] = useState([] as any)
  const [loading, setLoading] = useState('')
  const [userMap, setUserMap] = useState({} as any)

  function renderActionHistory(history: { content: string; userId: string | number; }) {
    if (!history.content) {
      return null
    }
    let actionInfo: any = {}
    try {
      actionInfo = JSON.parse(history.content)
    } catch (error) {
      return <span>操作历史数据异常</span>
    }
    switch (actionInfo.action) {
      case 'createRevision':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 创建了版本
          </span>
        )
      case 'renameRevision':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 将版本 "
            {actionInfo.before}" 重命名为 "{actionInfo.after}"
          </span>
        )
      case 'deleteRevision':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 删除了版本 "
            {actionInfo.label}"
          </span>
        )
      case 'lock_cell':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 锁定了单元格 "
            {actionInfo.name}-{actionInfo.range.join('')}"
          </span>
        )

      case 'unlock_cell':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 解锁了单元格 "
            {actionInfo.name}-{actionInfo.range.join('')}"
          </span>
        )
      case 'update_lock_cell':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 更新了单元格锁定 "
            {actionInfo.name}-{actionInfo.range.join('')}"
          </span>
        )
      case 'lock_sheet':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 锁定了工作表 "
            {actionInfo.name}"
          </span>
        )

      case 'unlock_sheet':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 解锁了工作表 "
            {actionInfo.name}"
          </span>
        )

      case 'create_sheet':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 创建了工作表 "
            {actionInfo.name}"
          </span>
        )

      case 'copy_sheet':
        return (
          <span>
            <a>{userMap[history.userId] || unkownUser}</a> 创建了工作表 "
            {actionInfo.name}" 的副本
          </span>
        )

      default:
        return <span>无法识别的操作历史</span>
    }
  }

  useEffect(() => {
    if (fileId && open) {
      setLoading('loading')
      fileService.getHistories(fileId).then((res) => {
        setHistories(res.data?.histories || [])
        setUserMap(res.data?.users || {})
        setLoading('loaded')
      }).catch(() => {
        setHistories([])
        setLoading('error')
      })
    }
  }, [fileId, open])

  return (
    <Modal
      open={open}
      title="查看历史列表"
      okText=""
      footer={<Button onClick={onClose}>关闭</Button>}
      onCancel={onClose}
    >
      {loading === 'loading' && (
        <Spin style={{ display: 'block' }} />
      )}
      {loading !== 'loading' && (
        <p> demo 仅用于演示，仅显示第一页侧边栏历史 </p>
      )}
      {histories.length > 0 && (
        <List
          dataSource={histories}
          renderItem={(item: any, index) => (
            <List.Item>
              <List.Item.Meta
                avatar={<Avatar icon={<HistoryOutlined />} />}
                title={<>
                  {item.historyType == 2 && (
                    <span>{item.userId.split(',').map((i: string) => userMap[i] || unkownUser).map((name: string, idx: number) => {
                      if (idx === 0) {
                        return <a> {name} </a>
                      }
                      return (
                        <span>
                          、 <a>{name} </a>
                        </span>
                      )
                    })}
                      编辑了文件</span>
                  )}
                  {item.historyType == 1 && (
                    <span>{renderActionHistory(item)}</span>
                  )}
                  <span className={styles.date}>{new Date(item.createdAt).toLocaleString()}</span>
                </>}
                description={item.label}
              />
            </List.Item>
          )}
        />
      )}
      {histories.length === 0 && loading === 'error' ? (
        <span>获取历史失败</span>
      ) : null}

      {histories.length === 0 && loading === 'loaded' ? (
        <span>暂无历史</span>
      ) : null}
    </Modal>
  )
}