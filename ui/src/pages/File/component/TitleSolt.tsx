import styles from '../FileList.module.less'
import { DataType } from "@/services/file.service"
import { Space, Tag, Popover } from 'antd';
import { fileConstants } from "@/constants"
import { MessageFilled, CreditCardFilled, HistoryOutlined, SyncOutlined } from '@ant-design/icons';
import { useEffect, useState } from 'react';
import { fileService } from "@/services/file.service";

const {
  TYPE_DOCUMENT,
  TYPE_DOCUMENT_PRO,
  TYPE_SPREADSHEET,
  TYPE_PRESENTATION
} = fileConstants

export function TitleSolt({ file, handleOpenRevisionModal, handleOpenHistoryModal }: { file: DataType, handleOpenRevisionModal: (file: DataType) => void, handleOpenHistoryModal: (file: DataType) => void }) {
  const [commentCountStatus, setcommentCountStatus] = useState('unknown')
  const [commentCount, setcommentCount] = useState(0)
  const [mentionAtListStatus, setmentionAtListStatus] = useState('unknown')
  const [mentionAtList, setmentionAtList] = useState([] as any)

  // Fetch comment and mention info on mount
  useEffect(() => {
    // Mentions
    if (file.shimoType === fileConstants.TYPE_SPREADSHEET) {
      getCommentCount(file.id)
    }
    // Comments (only spreadsheets support comment counts)
    if ([TYPE_DOCUMENT, TYPE_DOCUMENT_PRO, TYPE_SPREADSHEET].some((t) => t === file.shimoType)) {
      getMentionAt(file.id)
    }
  }, [file.id])
  // Fetch comment count
  function getCommentCount(fileId: string) {
    setcommentCountStatus('loading')
    fileService.getCommentCount(fileId).then(res => {
      setcommentCountStatus('loaded')
      setcommentCount(res.data?.count)
    }).catch(err => {
      setcommentCountStatus('error')
    })
  }
  // Fetch mention list
  function getMentionAt(fileId: string) {
    setmentionAtListStatus('loading')
    fileService.getMentionAt(fileId).then(res => {
      setmentionAtList(res.data ? res.data : [])
      setmentionAtListStatus('loaded')
    }).catch(err => {
      setmentionAtListStatus('error')
    })
  }

  return (
    <>
      <span className={styles.fileTitle}>{file.name}</span>
      {[TYPE_DOCUMENT, TYPE_DOCUMENT_PRO, TYPE_SPREADSHEET, TYPE_PRESENTATION].includes(file.shimoType) && (
        <Space>
          {file.shimoType === TYPE_SPREADSHEET ? (
            <Popover content={`${commentCount > 0 ? `${commentCount} 条评论` : '无评论'
              }`}>
              <Tag>
                <Space>
                  <MessageFilled />
                  {commentCountStatus === 'loading' ? (
                    <SyncOutlined spin />
                  ) : null}
                  {commentCountStatus === 'loaded' ? (
                    <span>{commentCount}</span>
                  ) : null}
                  {commentCountStatus === 'error' ? <span>加载失败</span> : null}
                </Space>
              </Tag>
            </Popover>
          ) : null}
          {[TYPE_DOCUMENT, TYPE_DOCUMENT_PRO, TYPE_SPREADSHEET].some(
            (t) => t === file.shimoType) ? (
            <Popover content={`${mentionAtList.length > 0
              ? mentionAtList.map((m: { name: string; }) => `@${m.name}`).join(' ')
              : '无 at 记录'
              }`}>
              <Tag>
                <Space>@
                  {mentionAtListStatus === 'loading' ? (
                    <SyncOutlined spin />
                  ) : null}
                  {mentionAtListStatus === 'loaded' ? (
                    <span>{mentionAtList.length || 0}</span>
                  ) : null}
                  {mentionAtListStatus === 'error' ? <span>加载失败</span> : null}
                </Space>
              </Tag>
            </Popover>
          ) : null}
          {[TYPE_DOCUMENT, TYPE_DOCUMENT_PRO, TYPE_SPREADSHEET, TYPE_PRESENTATION].some(
            (t) => t === file.shimoType) ? (
            <Popover content='点击查看版本'>
              <Tag onClick={() => handleOpenRevisionModal(file)}>
                <Space><CreditCardFilled />版本</Space>
              </Tag>
            </Popover>
          ) : null}
          {[TYPE_DOCUMENT, TYPE_DOCUMENT_PRO, TYPE_SPREADSHEET, TYPE_PRESENTATION].some(
            (t) => t === file.shimoType) ? (
            <Popover content='点击查看历史'>
              <Tag onClick={() => handleOpenHistoryModal(file)}>
                <Space><HistoryOutlined />历史</Space>
              </Tag>
            </Popover>
          ) : null}
        </Space>
      )}
    </>
  )
}