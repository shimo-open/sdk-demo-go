import { FileHeader } from '../component/FileHeader'
import { useEffect, useState } from "react";
import { useParams, useModel } from "umi";
import { fileService } from '@/services/file.service';
import { fileConstants } from '@/constants';
import { Spin } from 'antd';
import styles from '../ShimoFile/ShimoFile.less';
import { START_PARAMS_FIELD } from 'shimo-js-sdk'

export default function PreviewFile() {
  const [getFileState, setGetFileState] = useState('')
  const [currentFile, setCurrentFile] = useState({} as any)
  const { id } = useParams()
  const { lang } = useModel('user')
  const language = lang.current

  const loading = !getFileState || getFileState === fileConstants.GET_FILE_REQUEST
  const smParams = new URLSearchParams(location.search)?.get(START_PARAMS_FIELD) || ""

  useEffect(() => {
    if (!getFileState) {
      setGetFileState(fileConstants.GET_FILE_REQUEST)
      fileService.getShimoFile({ id: id as string, mode: 'preview', smParams: smParams, lang: language }).then(res => {
        setGetFileState(fileConstants.GET_FILE_SUCCESS)
        setCurrentFile(res.data)
      }).catch(() => {
        setGetFileState(fileConstants.GET_FILE_FAILURE)
      })
    }
  })

  return (
    <div className={styles.previewFile}>
      <FileHeader file={currentFile} mode="preview" />
      {loading && <Spin className={styles.spin} />}
      {currentFile && (
        <>
          <iframe
            title={currentFile.name}
            className={styles.preview}
            src={currentFile.previewUrl}
          ></iframe>
        </>
      )}
    </div>
  )
}