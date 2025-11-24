import { useEffect } from "react";
import { Table, Popover } from 'antd';
import { useState } from 'react';
import type { TableProps } from 'antd';
import { CheckCircleFilled, CloseCircleFilled, UnorderedListOutlined } from '@ant-design/icons';
import { TestApiDetail } from "./TestApiDetail"

export function TestApiResult({ type, data }: { type: string, data: any }) {
  const [tableItems, setTableItems] = useState([])
  const [open, setOpen] = useState(false)
  const [currrentRecord, setCurrrentRecord] = useState({} as any)

  useEffect(() => {
    if (!data) return
    data.forEach((item: any, index: number) => {
      item.index = index + 1 // Add an index column
    })
    setTableItems(data)
  }, [type, data])

  const fileColumns: TableProps['columns'] = [
    {
      title: '接口名字',
      dataIndex: 'apiName',
      key: 'apiName',
      width: 150
    },
    {
      title: '上传/导出文件后缀',
      dataIndex: 'fileExt',
      key: 'fileExt',
      width: 150
    },
    {
      title: '请求结果',
      dataIndex: 'success',
      key: 'success',
      width: 150,
      render: (success) => (
        <>
          <span style={{ marginRight: "6px" }}>{success ? 'true' : 'false'}</span>
          {success ?
            <CheckCircleFilled style={{ color: "rgb(76, 217,100)" }} /> :
            <CloseCircleFilled style={{ color: "rgb(237,71,58)" }} />}
        </>
      ),
    },
    {
      title: '测试开始时间',
      dataIndex: 'startTime',
      key: 'startTime',
      width: 150,
      render: (_, record) => (
        record.startTime ? new Date(record.startTime * 1000).toLocaleString() : 0
      ),
    },
    {
      title: '耗时',
      dataIndex: 'timeConsuming',
      key: 'timeConsuming',
      width: 100,
      render: (_, record) => (
        record.timeConsuming?.replace(/(\d+\.\d+|\d+)/g, (match: string) => {
          // match is the numeric substring
          let num = parseFloat(match);
          if (num.toString().includes('.') && num.toString().split('.')[1].length > 4) {
            return num.toFixed(4);
          } else {
            return match; // Preserve numbers with <=4 decimal places
          }
        })
      )
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 80,
      render: (_, record) => (
        <Popover content='查看详情'>
          <UnorderedListOutlined onClick={() => openDetail(record)} />
        </Popover>
      ),
    },
  ]
  const columns: TableProps['columns'] = [
    {
      title: '接口名字',
      dataIndex: 'apiName',
      key: 'apiName',
      width: 150
    },
    {
      title: '是否成功',
      dataIndex: 'success',
      key: 'success',
      render: (success) => (
        <>
          <span style={{ marginRight: "6px" }}>{success ? 'true' : 'false'}</span>
          {success ?
            <CheckCircleFilled style={{ color: "rgb(76, 217,100)" }} /> :
            <CloseCircleFilled style={{ color: "rgb(237,71,58)" }} />}
        </>
      ),
      width: 80
    },
    {
      title: '状态码',
      dataIndex: 'httpCode',
      key: 'httpCode',
      width: 100
    },
    {
      title: '返回结果',
      dataIndex: 'httpResp',
      key: 'httpResp',
      width: 200,
      ellipsis: true,
      render: (text) => (
        <Popover content={text}>
          {text}
        </Popover>
      ),
    },
    {
      title: '错误信息',
      dataIndex: 'errMsg',
      key: 'errMsg',
      width: 200,
      ellipsis: true,
      render: (text) => (
        <Popover content={text}>
          {text}
        </Popover>
      ),
    },
    {
      title: '测试开始时间',
      dataIndex: 'startTime',
      key: 'startTime',
      width: 150,
      render: (_, record) => (
        record.startTime ? new Date(record.startTime * 1000).toLocaleString() : 0
      ),
    },
    {
      title: '耗时',
      dataIndex: 'timeConsuming',
      key: 'timeConsuming',
      width: 100,
      render: (_, record) => (
        record.timeConsuming?.replace(/(\d+\.\d+|\d+)/g, (match: string) => {
          // match is the numeric substring
          let num = parseFloat(match);
          if (num.toString().includes('.') && num.toString().split('.')[1].length > 4) {
            return num.toFixed(4);
          } else {
            return match; // Preserve numbers with <=4 decimal places
          }
        })
      )
    },
    {
      title: '操作',
      dataIndex: 'action',
      key: 'action',
      width: 80,
      render: (_, record) => (
        <Popover content='查看详情'>
          <UnorderedListOutlined onClick={() => openDetail(record)} />
        </Popover>
      ),
    },
  ]

  function openDetail(record: any) {
    setCurrrentRecord(record)
    setOpen(true)
  }

  return (
    <>
      <Table
        columns={type == "FileIOResMap" ? fileColumns : columns}
        dataSource={tableItems}
        rowKey="index"
        pagination={false}
        scroll={{ y: "calc(100vh - 280px)" }} />
      <TestApiDetail open={open} onClose={() => setOpen(false)} record={currrentRecord} type={type} />
    </>
  )
}