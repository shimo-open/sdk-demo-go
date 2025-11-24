import { useParams, history, useModel } from 'umi'
import { useEffect, useState, useRef } from 'react'
import { FileHeader } from '../component/FileHeader'
import { connect, FileType, START_PARAMS_FIELD, UrlSharingType, GenerateUrlInfo } from 'shimo-js-sdk'
import isPlainObject from 'is-plain-obj'
import { userService } from '@/services/user.service'
import { fileConstants, userConstants } from '@/constants'
import { message, Spin } from 'antd'
import './ShimoFile.less'
import { fileService } from '@/services/file.service'
import { pick } from '@/utils/index'
import { setToken } from '@/utils/axios'
import { parseSmParams } from '@/utils/index'
import throttle from 'lodash/throttle'

const SDK_RENDERING = 1
const SDK_RENDERED = 2

export default function ShimoFile({ noHeader }: { noHeader: boolean }) {
  const { lang, hostPath, setUserInfo, setHostPath } = useModel('user')
  const language = lang.current
  const { setShimoFileSaveStatus, tableItems, setTableItems, setAvailableFileTypes, setPremiumFileTypes, setAllowBoardAndMindmap } = useModel('file')

  const [getFileState, setGetFileState] = useState('')
  const [prefix, setPrefix] = useState(hostPath)
  const [currentFile, setCurrentFile] = useState({} as any)

  const { id } = useParams()
  const container = useRef(null)
  const [sdkInfo, setSDKInfo] = useState(null as any)
  const [sdkState, setSDKState] = useState(0)
  const smParams = new URLSearchParams(location.search)?.get(START_PARAMS_FIELD) || ''
  const transSmParams: any = parseSmParams(smParams)

  useEffect(() => {
    if (!prefix) {
      getHostPath()
    }

    if (prefix && !getFileState) {
      setGetFileState(fileConstants.GET_FILE_REQUEST)
      fileService
        .getShimoFile({ id: id as string, mode: location.pathname.match(/\/form\/\w+\/fill/i) ? 'form_fill' : 'shimo', smParams: smParams, lang: language })
        .then((res) => {
          if (res.data.isShimoFile !== 1) {
            const pathPart = location.href.split('/shimo-files/')[1]
            const previewPath = `/preview/${pathPart}`
            history.push(previewPath)
            return
          }
          setGetFileState(fileConstants.GET_FILE_SUCCESS)
          setCurrentFile(res.data)
        })
        .catch((err) => {
          setGetFileState(fileConstants.GET_FILE_FAILURE)
        })
    }

    if (getFileState === fileConstants.GET_FILE_SUCCESS && currentFile && sdkState < SDK_RENDERING && prefix) {
      // Reset the save status each time the editor initializes
      setShimoFileSaveStatus('init')
      setSDKState(SDK_RENDERING)

      // If a non-default language is configured, pass it during editor initialization
      let editorLang = language && language !== userConstants.LANGUAGE_LOCAL ? language : null

      let paramsList: any = []
      const queryParams = new URLSearchParams(location.search)
      const originParams = queryParams.get(START_PARAMS_FIELD)

      if (originParams) {
        paramsList.push(originParams)
      }

      // Jump to the mentioned user's position
      if (queryParams.get('mentionId')) {
        paramsList.push({ hash: queryParams.get('mentionId') })
      }

      if (paramsList.length === 0) {
        paramsList = null
      }

      const options: any = {
        appId: currentFile.config.appId,
        fileId: id,
        endpoint: currentFile.config.endpoint,
        smParams: paramsList,
        signature: currentFile.config.signature,
        userUuid: currentFile.config.userUuid,
        generateUrl: async (
          /**
           * @type {string}
           */
          fileId: string,
          /**
           * @type {import('shimo-js-sdk').GenerateUrlInfo}
           */
          info = {} as GenerateUrlInfo
        ) => {
          if (info.sharingType === UrlSharingType.FormFill) {
            console.log('formFill >> prefix', prefix)
            return `${prefix}/form/${fileId}/fill`
          }
          console.log('default >> prefix', prefix)
          return `${prefix}/shimo-files/${fileId}`
        },
        getSignature() {
          return currentFile.config.signature
        },
        openLink(url: any, target: string) {
          if (target === '_blank') {
            window.open(url, target)
          } else {
            try {
              new URL(url)
              window.location = url
            } catch (e) {
              history.push(url)
            }
          }
        },
        token: currentFile.config.token,
        container: container.current,
        lang: editorLang,
      }

      try {
        Object.assign(options, JSON.parse(localStorage.getItem('SDK_INIT_OPTIONS') || ''))
      } catch (e) { }
      const handleSaveStatusChanged = throttle((data: { status: string }) => setShimoFileSaveStatus(data.status), 300)
      options.getFileInfoFromUrl = (url: string) => {
        let fromId
        const urlWithoutParams = url.split('?')[0]
        // Handle form-fill pages separately
        const match = urlWithoutParams.match(/\/form\/(\w+)\/fill/i)
        if (match) {
          fromId = match[1]
        } else {
          let splitPath = urlWithoutParams.split('/')
          fromId = splitPath[splitPath.length - 1]
        }
        console.log('FromId: ', fromId)
        return Promise.resolve({
          fileId: fromId,
          // type: currentFile.shimoType
        })
      }

      connect(options)
        .then((sdk: any) => {
          ; (window as any).__SDK__ = sdk
          if (sdk.fileType === FileType.Document) {
            sdk.document.on('titleChange', (title: string) => {
              if (currentFile.name !== title) {
                updateFile({ file: currentFile, data: { name: title } })
              }
            })

            // Wait for the lite document event before it becomes effective
            sdk.document.on('saveStatusChanged', handleSaveStatusChanged)

            // Lite documents currently emit the wrong event name with inconsistent payloads
            sdk.document.on('saveStatusDidChange', (status: any) => {
              setShimoFileSaveStatus(status.status)
            })
          } else {
            sdk.editor.on('saveStatusChanged', handleSaveStatusChanged)
            sdk.editor.on('error', handleError)
          }
          if (sdk.fileType === FileType.Form) {
            sdk.form.on('titleChange', (title: string) => {
              if (currentFile.name !== title) {
                updateFile({ file: currentFile, data: { name: title } })
              }
            })
          }

          setSDKInfo(sdk)
          setSDKState(SDK_RENDERED)
        })
        .catch((err: any) => {
          setSDKState(SDK_RENDERED)
          message.error(String(err))
        })
    }

    if (sdkState == SDK_RENDERING) {
      // End the loading state once the iframe finishes; covers cases where connect never throws on errors
      let iframe: any = document.querySelector('iframe')
      if (iframe) {
        if (iframe.attachEvent) {
          iframe.attachEvent('onload', function () {
            setSDKState(SDK_RENDERED)
          })
        } else {
          iframe.onload = function () {
            setSDKState(SDK_RENDERED)
          }
        }
      }
    }
  }, [id, currentFile, sdkState, getFileState, location, history, language, prefix])

  // Keep the title in sync
  useEffect(() => {
    if (currentFile && currentFile.name && sdkInfo) {
      sdkInfo.getEditor?.().setTitle?.(currentFile.name)
    }
  }, [sdkInfo, currentFile])

  // Fetch the hostPath
  function getHostPath() {
    userService.auth().then((res) => {
      if (
        res.data.user.id < 0 &&
        !location.pathname.match(/\/form\/\w+\/fill/i) &&
        !(transSmParams.type == 'form' && transSmParams.path?.includes('response-share'))
      ) {
        userService.logout()
      } else {
        const {
          data: { user, hostPath, token, availableFileTypes, premiumFileTypes, allowBoardAndMindmap },
        } = res
        setUserInfo(user)
        setHostPath(hostPath)
        setToken(token || '')
        setPrefix(hostPath)
        setAvailableFileTypes(availableFileTypes)
        setPremiumFileTypes(premiumFileTypes)
        setAllowBoardAndMindmap(allowBoardAndMindmap)
      }
    })
  }
  // Let the header menu invoke editor methods
  function editorMethodCall(methodName: string, ...args: any[]) {
    if (!sdkInfo) {
      message.error('石墨 SDK 未初始化')
      return
    }

    const editor: any = sdkInfo[currentFile.shimoType]
    if (!editor) {
      message.error(`不支持的在石墨类型 ${currentFile.shimoType} 上调用方法，可能需要更新 shimo-js-sdk`)
      return
    }

    editor[methodName](...args).catch((err: any) => {
      message.error(String(err.message))
    })
  }

  function onFileNameChanged(name: string) {
    if (currentFile.name !== name) {
      updateFile({ file: currentFile, data: { name: name } })
    }
  }

  function handleError({ code, data }: { code: any; data: any }) {
    let error = ''
    if (!isNaN(code)) {
      error = String(code)
    } else {
      error = 'UNKNOWN_ERR_CODE: ' + code
    }

    if (data) {
      if (isPlainObject(data)) {
        error += ': ' + JSON.stringify(data)
      } else {
        error += ': ' + String(data)
      }
    }
    message.error(error)
  }

  const updateFile = (payload: { file: any; data: { name: string } }) => {
    const oldData = pick(payload.file, Object.keys(payload.data))
    setUpdateFile({ fileId: payload.file.id, data: payload.data })
    fileService.updateFile(payload.file.id, payload.data).catch((err) => {
      setUpdateFile({ fileId: payload.file.id, data: oldData })
    })
  }

  const setUpdateFile = (payload: { fileId: string; data: any }) => {
    if (currentFile) {
      setCurrentFile(Object.assign({}, currentFile, payload.data))
    }
    let copyItem = [...tableItems]
    if (Array.isArray(copyItem)) {
      for (const file of tableItems) {
        if (file.id === payload.fileId) {
          Object.assign(file, payload.data)
          break
        }
      }
    }
    setTableItems(copyItem)
  }

  const loading = !getFileState || getFileState === fileConstants.GET_FILE_REQUEST || sdkState < SDK_RENDERED
  return (
    <>
      {loading && <Spin className="spin" />}
      {noHeader !== true &&
        !(
          transSmParams.type == 'form' &&
          (transSmParams.path?.includes('response-share') || transSmParams.path?.includes('fill') || transSmParams.path?.includes('preview'))
        ) &&
        !new URLSearchParams(location.search)?.get(userConstants.ACCESS_TOKEN_KEY) &&
        new URLSearchParams(location.search)?.get('mode') !== 'insert' && (
          <FileHeader file={currentFile} mode="edit" onEditorMethodCall={editorMethodCall} onFileNameChanged={onFileNameChanged} />
        )}
      <div className="shimo-file-container" ref={container}></div>
    </>
  )
}
