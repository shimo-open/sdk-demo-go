import { LeftOutlined, SyncOutlined, CheckOutlined, CloseOutlined } from '@ant-design/icons'
import { Link } from 'react-router-dom'
import { useRef, useEffect } from 'react'
import '@/pages/HeaderPage/Header.module.less'
import './FileHeader.less'
import { useModel } from 'umi'
import { FileHeaderMenu } from './FileHeaderMenu'
import { message, Button } from 'antd'

let disableSaveLimitedToast = false
export function FileHeader({
  file,
  mode,
  onEditorMethodCall,
  onFileNameChanged,
}: {
  file: any
  mode: string
  onEditorMethodCall?: (methodName: string, ...args: any[]) => void
  onFileNameChanged?: (name: string) => void
}) {
  const inputRef: any = useRef()
  const { shimoFileSaveStatus } = useModel('file')

  useEffect(() => {
    if (file) {
      inputRef.current.textContent = file?.name
    }
  }, [file])
  const saveLimitedToastKey = 'saveLimited'
  useEffect(() => {
    if (disableSaveLimitedToast) return
    if (shimoFileSaveStatus === 'saveLimited') {
      message.open({
        type: 'info',
        style: { marginTop: '50px' },
        content: (
          <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
            <span>在线编辑人数已达上限，所有编辑将会暂存本地</span>
            <Button
              type="link"
              size="small"
              onClick={() => {
                disableSaveLimitedToast = true
                message.destroy(saveLimitedToastKey)
              }}
            >
              知道了
            </Button>
          </div>
        ),
        duration: 0, // Do not auto-close
        key: saveLimitedToastKey,
      })
    } else if (shimoFileSaveStatus !== 'saving') {
      message.destroy(saveLimitedToastKey)
    }
  }, [shimoFileSaveStatus])

  function onKeyEvent(e: { keyCode: number; preventDefault: () => void }) {
    if (e.keyCode === 13) {
      e.preventDefault()
      updateFileName()
    }
  }

  function updateFileName() {
    const name = inputRef.current?.textContent.trim()
    if (name !== file.name) {
      onFileNameChanged ? onFileNameChanged(name) : null
    }
  }

  return (
    <header className="header shimo-file-header">
      <h3 className="file-name">
        <Link to={`/`}>
          {' '}
          <LeftOutlined style={{ color: '#000', marginLeft: '10px' }} />
        </Link>
        <div className="input" contentEditable={mode == 'edit'} ref={inputRef} onKeyDown={onKeyEvent} onKeyUp={onKeyEvent} onBlur={updateFileName}></div>
        {mode === 'edit' ? (
          <div className="file-save-status">
            {shimoFileSaveStatus === 'init' ? <span>文件将会自动保存</span> : null}
            {shimoFileSaveStatus === 'saving' ? (
              <span>
                <SyncOutlined className="save-state-icon" spin />
                正在保存...
              </span>
            ) : null}
            {shimoFileSaveStatus === 'saved' ? (
              <span>
                <CheckOutlined className="save-state-icon" />
                保存成功
              </span>
            ) : null}
            {shimoFileSaveStatus === 'error' ? (
              <span className="save-error">
                <CloseOutlined className="save-state-icon" />
                保存失败
              </span>
            ) : null}
            {shimoFileSaveStatus === 'saveLimited' ? <span>编辑已暂存本地，正在等待同步....</span> : null}
          </div>
        ) : null}
      </h3>
      {!!file && !!file.isShimoFile && onEditorMethodCall && <FileHeaderMenu file={file} onMenuItemClick={onEditorMethodCall} />}
    </header>
  )
}
