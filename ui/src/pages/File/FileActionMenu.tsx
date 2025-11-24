import React, { useState, useEffect, useRef } from 'react'
import { Button, Upload, Dropdown, Space, Modal, Form, Input, Select, message } from 'antd';
import type { MenuProps, UploadProps, FormInstance } from 'antd'
import { CaretDownFilled } from '@ant-design/icons';
import { fileConstants, userConstants } from '@/constants'
import { history, useModel } from 'umi'
import { fileService } from '@/services/file.service';
import { convertFileType, FileType } from 'shimo-js-sdk-shared'
interface Values {
  url: string;
  fileName: string;
  type: string;
}
interface CollectionCreateFormProps {
  onFormInstanceReady: (instance: FormInstance<Values>) => void;
  fileType: Array<{ value: string, label: string }>
}
// URL import form
const CollectionCreateForm: React.FC<CollectionCreateFormProps> = ({
  onFormInstanceReady,
  fileType
}) => {
  const [form] = Form.useForm();
  useEffect(() => {
    onFormInstanceReady(form);
  }, []);
  return (
    <Form form={form} labelCol={{ span: 6 }} wrapperCol={{ span: 18 }} style={{ margin: '2rem 0' }}>
      <Form.Item name="url" label="文件下载链接">
        <Input placeholder="https://domain.com/download/test.docx" />
      </Form.Item>
      <Form.Item name="fileName" label="文件名称">
        <Input placeholder="文件名称需带有正确的后缀名( 如：test.docx )" />
      </Form.Item>
      <Form.Item name="type" label="导入文件类型" >
        <Select options={fileType} />
      </Form.Item>
    </Form>
  );
};

export function FileActionMenu(props: { className: any; }) {
  const [open, setOpen] = useState(false)
  const [formValues, setFormValues] = useState({} as any);

  const { file, createState, setCreateState, setFile, uploadState, setUploadState, importState, setImportState, createShimoFile, importUrlState, setImportUrlState, availableFileTypes, premiumFileTypes, allowBoardAndMindmap } = useModel('file')
  const { lang } = useModel('user')

  const fileImportInput = useRef(null)

  // Document creation menu entries
  const fileCreateOptions: MenuProps['items'] = []

  const fileImportOptions: MenuProps['items'] = []
  const urlImportOptions: any = []

  const handleFileTypes = (data: any, suffix: string) => {
    if (!data) data = []
    for (const t of data) {
      switch (convertFileType(t)) {
        case FileType.Document: {
          urlImportOptions.push({ value: 'document', label: '轻文档' + suffix })
          fileCreateOptions.push({
            key: t,
            label: '轻文档' + suffix,
          })
          fileImportOptions.push({
            key: t,
            label: '轻文档' + suffix,
          })
          break
        }

        case FileType.DocumentPro: {
          urlImportOptions.push({ value: 'documentPro', label: '传统文档' + suffix })
          fileCreateOptions.push({
            key: t,
            label: '传统文档' + suffix,
          })
          fileImportOptions.push({
            key: t,
            label: '传统文档' + suffix,
          })
          break
        }

        case FileType.Spreadsheet: {
          urlImportOptions.push({ value: 'spreadsheet', label: '表格' + suffix })
          fileCreateOptions.push({
            key: t,
            label: '表格' + suffix,
          })
          fileImportOptions.push({
            key: t,
            label: '表格' + suffix,
          })
          break
        }

        case FileType.Table: {
          fileCreateOptions.push({
            key: t,
            label: '应用表格' + suffix,
          })
          break
        }

        case FileType.Presentation: {
          urlImportOptions.push({ value: 'presentation', label: '专业幻灯片' + suffix })
          fileCreateOptions.push({
            key: t,
            label: '专业幻灯片' + suffix,
          })
          fileImportOptions.push({
            key: t,
            label: '专业幻灯片' + suffix,
          })
          break
        }

        case FileType.Form: {
          fileCreateOptions.push({
            key: t,
            label: '表单' + suffix,
          })
          break
        }

        case FileType.Flowchart: {
          fileCreateOptions.push({
            key: t,
            label: '流程图' + suffix,
          })
          break
        }

        default:
      }
      // Additional entries
      switch (t) {
        case 'board': {
          fileCreateOptions.push({
            key: t,
            label: '白板' + suffix,
          })
          break
        }
        case 'mindmap': {
          fileCreateOptions.push({
            key: t,
            label: '思维导图' + suffix,
          })
          break
        }
      }
    }
  }

  handleFileTypes(availableFileTypes, "")

  if (premiumFileTypes && premiumFileTypes.length > 0) {
    handleFileTypes(premiumFileTypes, "（增值席位）")
  }

  if (allowBoardAndMindmap) {
    handleFileTypes(['board', 'mindmap'], "（非席位新增）")
  }

  const uploadProps: UploadProps = {
    name: 'file',
    action: '',
    customRequest(option: any) {
      try {
        setUploadState(fileConstants.UPLOAD_FILE_REQUEST)
        fileService.uploadFile(option.file).then(res => {
          setUploadState(fileConstants.UPLOAD_FILE_SUCCESS)
          option.onSuccess(res.data.filePath)
        }).catch(error => {
          setUploadState(fileConstants.UPLOAD_FILE_FAILURE)
          option.onError(error?.data?.message)
        })
      } catch (error) {
        option.onError(error)
        setUploadState(fileConstants.UPLOAD_FILE_FAILURE)
      }
    },
    showUploadList: false
  }

  const language = lang.current

  useEffect(() => {
    if (createState === fileConstants.CREATE_SHIMO_FILE_SUCCESS && file) {
      setCreateState('')
      setFile(null)
      history.push(`/shimo-files/${file.id}`)
      return
    }

    if (importState === fileConstants.IMPORT_FILE_SUCCESS && file) {
      setImportState('')
      history.push(`/shimo-files/${file.id}`)
    }
  }, [createState, importState, file])


  function importFileByUrl() {
    setOpen(true)
  }
  // Handle URL-based imports
  function handleSubmit(event: any) {
    setOpen(false)
    event.preventDefault()
    const formVal = formValues.getFieldsValue(true)
    setImportUrlState(fileConstants.IMPORT_URL_REQUEST)
    fileService.importFileByUrl(formVal.url, formVal.fileName, formVal.type).then(async res => {
      do {
        try {
          if (!(res.data && res.data.taskId)) {
            message.error("导入失败: taskid获取失败")
            throw new Error(res.data)
          }
          let resp: any = await fileService.getImportByUrlProgress(res.data.taskId, res.data.id)

          if (resp.data?.status !== 0) {
            throw new Error(resp.data?.message)
          }
          console.log("import processing...");

          if (resp.data.data.progress === 100) {
            setFile(res.data)
            setImportUrlState(fileConstants.IMPORT_URL_SUCCESS)
            console.log("导入完成", resp.data);
            break
          }
        } catch (e) {
          console.log(e + "导入进度获取失败")
          setImportUrlState(fileConstants.IMPORT_URL_FAILURE)
          break
        }
        await new Promise((resolve) => setTimeout(resolve, 1000))
      } while (true)
    }).catch(error => {
      setImportUrlState(fileConstants.IMPORT_URL_FAILURE)
    })
  }

  const creating = createState === fileConstants.CREATE_SHIMO_FILE_REQUEST

  const createClick: MenuProps['onClick'] = ({ key }) => {
    console.log(`Click on item ${key}`);
    const lang = language && language !== userConstants.LANGUAGE_LOCAL ? language : null
    createShimoFile({ type: key, lang: lang })
  };

  const importClick: MenuProps['onClick'] = ({ key }) => {
    console.log(`Click on item ${key}`);
    let accept = ''
    switch (key) {
      // Lite doc
      case fileConstants.TYPE_DOCUMENT:
        accept = '.docx,.doc,.md,.html,.txt'
        break;
      // Classic doc
      case fileConstants.TYPE_DOCUMENT_PRO:
        accept = '.docx,.doc'
        break;
      // Spreadsheet
      case fileConstants.TYPE_SPREADSHEET:
        accept = '.xlsx,.xls,.csv,.xlsm'
        break;
      // Presentation
      case fileConstants.TYPE_PRESENTATION:
        accept = '.pptx,.ppt'
        break;
    }
    /**
     * @type {HTMLInputElement}
     */
    const input: any = fileImportInput.current
    input.value = ''
    input.accept = accept
    input.dataset.shimoType = key
    input.click()
  };
  // Import a document from disk
  const importFile = () => {
    const input: any = fileImportInput.current
    const file = input?.files[0]
    setImportState(fileConstants.IMPORT_FILE_REQUEST)

    fileService.importFile(file, input.dataset.shimoType).then(async (res) => {
      do {
        try {
          if (!(res.data && res.data.taskId)) {
            message.error("导入失败: taskid获取失败")
            throw new Error(res.data)
          }
          let resp: any = await fileService.getImportProgress(res.data.taskId, res.data.id)

          if (resp.data?.status !== 0) {
            throw new Error(resp.data?.message)
          }

          console.log("import processing...");

          if (resp.data.data.progress === 100) {
            setFile(res.data)
            setImportState(fileConstants.IMPORT_FILE_SUCCESS)
            console.log("导入完成", resp.data);
            break
          }
        } catch (e) {
          console.log(e + "导入进度获取失败")
          setImportState(fileConstants.IMPORT_FILE_FAILURE)
          break
        }
        await new Promise((resolve) => setTimeout(resolve, 1000))
      } while (true)
    }).catch((e) => {
      setImportState(fileConstants.IMPORT_FILE_FAILURE)
    })
  }

  return (
    <div className={['action-menu', props.className || ''].join(' ')}>
      <Space.Compact>
        <Dropdown menu={{ items: fileCreateOptions, onClick: createClick }} trigger={['click']} disabled={creating || importState === fileConstants.IMPORT_FILE_REQUEST}>
          <Button onClick={(e) => e.preventDefault()} loading={creating}>
            新建石墨文档
            <CaretDownFilled />
          </Button>
        </Dropdown>
        <Upload {...uploadProps}>
          <Button loading={uploadState == fileConstants.UPLOAD_FILE_REQUEST}>上传文件</Button>
        </Upload>
        <Dropdown menu={{ items: fileImportOptions, onClick: importClick }} trigger={['click']}
          disabled={creating || uploadState === fileConstants.UPLOAD_FILE_REQUEST || importUrlState == fileConstants.IMPORT_URL_REQUEST}>
          <Button onClick={(e) => e.preventDefault()} loading={importState == fileConstants.IMPORT_FILE_REQUEST}>
            <Space>
              导入到石墨文档
              <CaretDownFilled />
            </Space>
          </Button>
        </Dropdown>
        <Button onClick={importFileByUrl} loading={importUrlState == fileConstants.IMPORT_URL_REQUEST}>Url导入</Button>
      </Space.Compact>
      <Modal open={open} title="导入到石墨文档" okText="导入" cancelText="取消"
        okButtonProps={{ autoFocus: true }} onCancel={() => setOpen(false)}
        destroyOnClose
        onOk={handleSubmit}
      >
        <CollectionCreateForm
          onFormInstanceReady={(instance: any) => {
            setFormValues(instance);
          }}
          fileType={urlImportOptions}
        />
      </Modal>

      <input
        ref={fileImportInput}
        type="file"
        style={{ display: 'none' }}
        onChange={importFile}
      />
    </div>
  )
}

