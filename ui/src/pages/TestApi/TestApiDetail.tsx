import { Modal, Descriptions, Button, message } from "antd";
import type { DescriptionsProps } from 'antd';
import styles from "./TestApi.module.less"
import { useEffect, useState } from "react";

export function TestApiDetail({ record, open, onClose, type }: { record: any, open: boolean, onClose: () => void, type: string }) {
  const isFileIOResMap = type === "FileIOResMap"

  const [formKeys, setFormKeys] = useState<DescriptionsProps['items']>([])
  const [fileFormKeys, setFileFormKeys] = useState<DescriptionsProps['items']>([])

  useEffect(() => {
    let keys: DescriptionsProps['items'] = [
      { key: 'testId', label: "测试ID", children: record.testId, span: 3 },
      { key: 'apiName', label: "接口名字", children: record.apiName, span: 3 },
      { key: 'success', label: "是否成功", children: String(record.success ? true : false), span: 3 },
      { key: 'httpCode', label: "状态码", children: record.httpCode, span: 3 },
      { key: 'httpResp', label: "返回结果", children: record.httpResp, span: 3 },
      { key: 'errMsg', label: "错误信息", children: record.errMsg, span: 3 },
      { key: 'startTime', label: "测试开始时间", children: record.startTime ? new Date(record.startTime * 1000).toLocaleString() : 0, span: 3 },
      { key: 'timeConsuming', label: "耗时", children: record.timeConsuming, span: 3 },
      { key: 'pathStr', label: "接口请求地址", children: record.pathStr, span: 3 },
    ]
    if (record.query && record.query !== "") {
      keys.push({ key: 'query', label: "地址传参", children: record.query, span: 3 })
    }
    if (record.bodyReq && record.bodyReq !== "") {
      keys.push({ key: 'bodyReq', label: "请求体传参", children: record.bodyReq, span: 3 })
    }
    if (record.formData && record.formData !== "") {
      keys.push({ key: 'formData', label: "FORM_DATA传参", children: record.formData, span: 3 })
    }
    setFormKeys(keys)

    let fileKeys: DescriptionsProps['items'] = [
      { key: 'apiName', label: "接口名字", children: record.apiName, span: 3 },
      { key: 'success', label: "请求结果", children: String(record.success ? true : false), span: 3 },
      { key: 'fileExt', label: record.apiName?.includes("导入") ? "上传文件后缀" : "导出文件类型", children: record.fileExt, span: 3 },
      { key: 'startTime', label: "测试开始时间", children: record.startTime ? new Date(record.startTime * 1000).toLocaleString() : 0, span: 3 },
      { key: 'timeConsuming', label: "耗时", children: record.timeConsuming, span: 3 },
      { key: 'pathStr', label: "接口请求地址", children: record.pathStr, span: 3 }
    ]

    if (record.formData && record.formData !== "") {
      fileKeys.push({ key: 'formData', label: "FORM_DATA传参", children: record.formData, span: 3 })
    }
    setFileFormKeys(fileKeys)
  }, [record])

  // Handle copy action
  function onCopy() {
    let text = ""
    let data = isFileIOResMap ? fileFormKeys : formKeys
    data?.forEach(v => {
      text += v.label + ": " + v.children + "\n"
    });
    navigator.clipboard.writeText(text).then(() => {
      message.success('复制成功');
    }, () => {
      message.error('复制失败');
    });
  }

  return (
    <Modal
      open={open}
      title="详情"
      okText="复制"
      onOk={onCopy}
      onCancel={onClose}
      width="80vh"
    >
      <Descriptions
        className={styles.detail}
        items={isFileIOResMap ? fileFormKeys : formKeys}
        bordered
        labelStyle={{ width: isFileIOResMap ? "30%" : "25%" }} />
    </Modal>
  )
}