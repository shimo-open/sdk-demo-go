import { Table, Modal, Button, Popover } from 'antd';
import type { TableProps } from 'antd';
import { DatabaseFilled } from '@ant-design/icons';
import { useEffect, useState } from 'react';
import styles from './ShimoEvents.less'
import { eventService } from '@/services/event.service';
import { Comment } from './Events/Comment'
import { Discussion } from './Events/Discussion'
import { MentionAt } from './Events/MentionAt'
import { DateMention } from './Events/DateMention'
import { FileContent } from './Events/FileContent'
import { Collaborator } from './Events/Collaborator'
import { Revision } from './Events/Revision'
import { System } from './Events/System'

export default function ShimoEventsPage() {
  const [tableItems, setTableItems] = useState([] as any)
  const [loading, setLoading] = useState(false)
  const [page, setPage] = useState(1)
  const [pageSize, setPageSize] = useState(10)
  const [count, setCount] = useState(0)
  const [open, setOpen] = useState(false)
  const [currentData, setCurrentData] = useState({})

  function renderCredentialType(headers: any) {
    let headersData: any = {}
    let invalidEventHeaders = false

    try {
      headersData = JSON.parse(headers)
    } catch (error) {
      console.warn('parse event headers data error', error)
      invalidEventHeaders = true
    }
    if (invalidEventHeaders) {
      return <span>-</span>
    }

    if (!headersData) {
      return <span>无效事件数据</span>
    }

    const credentialHeader = Object.keys(headersData).find(
      (h) => h && String(h).toLowerCase() === 'x-shimo-credential-type'
    )

    if (!credentialHeader) {
      return <span>未找到凭证类型数据</span>
    }

    const credentialType = headersData[credentialHeader]

    return <span>{credentialType}</span>
  }

  function renderByTypes(record: any) {
    let invalidEventData = false
    let eventData = {}
    try {
      eventData = JSON.parse(record.rawData)
    } catch (error) {
      console.warn('parse event raw data error', error)
      invalidEventData = true
    }
    if (invalidEventData) {
      return <span>无效事件数据</span>
    }
    switch (record.type) {
      case 'Comment':
        return <Comment data={eventData} files={record.files} users={record.users} />
      case 'Discussion':
        return <Discussion data={eventData} files={record.files} users={record.users} />
      case 'MentionAt':
        return <MentionAt data={eventData} files={record.files} users={record.users} />
      case 'DateMention':
        return <DateMention data={eventData} files={record.files} users={record.users} />
      case 'FileContent':
        return <FileContent data={eventData} files={record.files} users={record.users} />
      case 'Collaborator':
        return <Collaborator data={eventData} files={record.files} users={record.users} />
      case 'Revision':
        return <Revision data={eventData} files={record.files} users={record.users} />
      case 'System':
        return <System data={eventData} />
      default:
        return <div>暂不支持的事件类型</div>
    }
  }

  const columns: TableProps['columns'] = [
    {
      title: '事件信息',
      dataIndex: 'eventData',
      key: 'eventData',
      render: (_, record) => (
        renderByTypes(record)
      ),
      width: '50%'
    },
    {
      title: '凭证类型',
      dataIndex: 'headers',
      key: 'headers',
      width: '12%',
      render: (headers) => (
        renderCredentialType(headers)
      )
    },
    {
      title: '事件保存时间',
      dataIndex: 'createdAt',
      width: '15%',
      render: (createdAt) => (new Date(Number(createdAt) * 1000).toLocaleString())
    },
    {
      title: '操作',
      dataIndex: 'action',
      render: (_, record) => (
        <>
          <Popover content='查看事件原始数据'>
            <DatabaseFilled onClick={() => { setOpen(true); setCurrentData(record) }} />
          </Popover>
        </>
      ),
      width: '10%',
    },
  ]

  useEffect(() => {
    setLoading(true)
    eventService.getEvents({ page: page, size: pageSize, fileId: '' }).then(res => {
      setTableItems(res.data.list)
      setCount(res.data.count)
    }).finally(() => {
      setLoading(false)
    })
  }, [page, pageSize])

  function handlePaginationChange(page: number, pageSize: number) {
    setPage(page)
    setPageSize(pageSize)
  }

  return (
    <>
      <Table columns={columns} dataSource={tableItems} rowKey="id" className={styles.TableProps} loading={loading}
        pagination={{
          onChange: handlePaginationChange,
          current: page,
          pageSize: pageSize,
          total: count,
          showTotal: (total) => `共 ${total} 条`
        }}
        scroll={{ y: 565 }}
      />
      <DataModal data={currentData} open={open} onClose={() => setOpen(false)} />
    </>
  )
}

function DataModal({ data, open, onClose }: { data: any, open: boolean, onClose: () => void }) {
  const [headersData, setHeadersData] = useState<any>({})
  const [eventData, setEventData] = useState<any>({})

  useEffect(() => {
    if (data && data.headers) {
      setHeadersData(JSON.parse(data.headers))
    }
    if (data && data.rawData) {
      setEventData(JSON.parse(data.rawData))
    }
  }, [data])

  return (
    <Modal
      open={open && data}
      title="事件原始数据"
      footer={<Button onClick={onClose}>关闭</Button>}
      onCancel={onClose}
      width="80vh"
      className={styles.ModalProps}
    >
      <div>
        <h4 className={styles.modalDataTitle}>Event Request Headers</h4>
        <pre className={styles.modalDataBlock}>
          {JSON.stringify(headersData, null, 2)}
        </pre>
        <h4 className={styles.modalDataTitle}>Event Request Body</h4>
        <pre className={styles.modalDataBlock}>
          {JSON.stringify(eventData, null, 2)}
        </pre>
      </div>
    </Modal>
  )
}