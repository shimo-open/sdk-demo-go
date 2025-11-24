import { fileConstants, userConstants } from '../constants/'
import { fileService } from "@/services/file.service";
import { useState } from 'react';

const fileModel = () => {
  const [tableItems, setTableItems] = useState([] as any)

  const [shimoFileSaveStatus, setShimoFileSaveStatus] = useState('init')

  const [createState, setCreateState] = useState('')

  const [uploadState, setUploadState] = useState('')

  const [importState, setImportState] = useState('')
  const [importUrlState, setImportUrlState] = useState('')
  // Export state
  const [exportState, setExportState] = useState('')
  // Duplicate creation state
  const [duplicateState, setDuplicateState] = useState({} as any)

  const [exportedUrls, setExportedUrls] = useState({} as any)

  const [file, setFile] = useState<any>(null)

  const [availableFileTypes, setAvailableFileTypes] = useState([] as any)
  const [premiumFileTypes, setPremiumFileTypes] = useState([] as any)
  const [allowBoardAndMindmap, setAllowBoardAndMindmap] = useState(false)

  const createShimoFile = (payload: { type: string, lang: string }) => {
    setCreateState(fileConstants.CREATE_SHIMO_FILE_REQUEST)
    setFile(null)
    fileService.createShimoFile(payload).then(res => {
      setCreateState(fileConstants.CREATE_SHIMO_FILE_SUCCESS)
      setFile(res.data)
    }).catch(err => {
      setCreateState(fileConstants.CREATE_SHIMO_FILE_FAILURE)
    })
  }

  return {
    tableItems,
    setTableItems,
    shimoFileSaveStatus,
    createState,
    setCreateState,
    uploadState,
    setUploadState,
    importState,
    setImportState,
    file,
    setFile,
    availableFileTypes, setAvailableFileTypes,
    premiumFileTypes, setPremiumFileTypes,
    exportState, setExportState,
    duplicateState,
    setDuplicateState,
    exportedUrls,
    setExportedUrls,
    importUrlState,
    setImportUrlState,
    setShimoFileSaveStatus,
    createShimoFile,
    allowBoardAndMindmap,
    setAllowBoardAndMindmap
  }
}
export default fileModel
