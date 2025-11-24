import React, { useState } from 'react'
import { userConstants } from '@/constants'
import { useModel } from 'umi'
import { Modal, Button, Select, Space } from 'antd'
import { GlobalOutlined } from '@ant-design/icons'

const languageOptions = [
  { key: 'Local', label: '系统默认', value: userConstants.LANGUAGE_LOCAL },
  { key: 'Chinese', label: '简体中文', value: userConstants.LANGUAGE_ZH_CN },
  { key: 'English', label: 'English', value: userConstants.LANGUAGE_EN },
  { key: 'Japanese', label: '日本語', value: userConstants.LANGUAGE_JA },
  { key: 'Arabic', label: 'عربى', value: userConstants.LANGUAGE_AR_SA },
]

export function LanguageModal({ open, currentLang, onClose }: { open: boolean, currentLang: string, onClose: () => void }) {
  const { lang, setLang } = useModel('user')
  const [tempLang, setTempLang] = useState('')

  function handleChange() {
    closeReset()
    if (tempLang == "") {
      return
    }
    try {
      localStorage.setItem(userConstants.LANGUAGE_STORAGE_KEY, tempLang)
    } catch (err) {
      console.warn(`write language value ${tempLang} to localStorage error`)
    }
    setLang({ current: tempLang })
  }

  function closeReset() {
    onClose()
    setTempLang('')
  }

  return (
    <Modal
      open={open}
      title="切换语言"
      footer={[<Button onClick={closeReset} key='cancel'>取消</Button>, <Button onClick={handleChange} key='create'>确认</Button>]}
      onCancel={onClose}
      width={400}
    >
      <p>选择石墨文件编辑器界面语言</p>
      <Space.Compact>
        <Button><GlobalOutlined /></Button>
        <Select options={languageOptions} value={tempLang || currentLang} style={{ width: '100px' }}
          onChange={(e) => {
            setTempLang(e)
          }}></Select>
      </Space.Compact>
    </Modal>
  )
}
