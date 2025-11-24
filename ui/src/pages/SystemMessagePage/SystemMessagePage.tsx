import { Space, Input, DatePicker, Button, Table, message } from "antd";
import { useEffect, useState } from "react"
import { systemMessageService } from '@/services/systemMessage.service'
import dayjs from "dayjs";
import weekday from "dayjs/plugin/weekday";
import localeData from "dayjs/plugin/localeData";

dayjs.extend(weekday);
dayjs.extend(localeData);
dayjs.locale("en");
dayjs.locale("vi");

const { RangePicker } = DatePicker;
export default function SystemMessagePage() {
  const [tableLoading, setTableLoading] = useState(false);
  const [errorCallLoading, setErrorCallLoading] = useState(false);
  const [appIdValue, setAppIdValue] = useState<string>('');
  const [tableData, setTableData] = useState<any[]>([])
  const columns = [
    {
      title: 'AppID',
      dataIndex: 'app_id',
      key: 'app_id'
    },
    {
      title: '事件类型',
      dataIndex: 'event_type',
      key: 'event_type'
    },
    {
      title: '通知时间',
      dataIndex: 'notify_time',
      key: 'notify_time'
    }
  ]
  useEffect(() => {
    handleChange([dayjs().subtract(1, 'day').format('YYYY-MM-DD hh:mm:ss'), dayjs().format('YYYY-MM-DD hh:mm:ss')])
  }, [])

  function handleChange(timeRange: any) {
    const appID = appIdValue
    setTableLoading(true)
    systemMessageService.getSystemMessages(appID, dayjs(timeRange[0]).format('YYYY-MM-DD hh:mm:ss'), dayjs(timeRange[1]).format('YYYY-MM-DD hh:mm:ss')).then((res) => {
      if (res && res.data?.length && Array.isArray(res.data)) {
        let data = res.data
        for (let i = 0; i < data.length; i++) {
          data[i]['key'] = i
          data[i]['notify_time'] = new Date(
            data[i]['notify_time']
          ).toLocaleString()
        }
        setTableData(data?.reverse())
      }
    }).finally(() => {
      setTableLoading(false)
    })
  }

  function errorCallback() {
    setErrorCallLoading(true)
    systemMessageService.errorCallback().then((res) => {
      message.success('success')
    }).finally(() => {
      setErrorCallLoading(false)
    })
  }
  return (
    <>
      <Space style={{ margin: '10px' }}>
        <Input placeholder="App ID" value={appIdValue} onChange={(e: any) => setAppIdValue(e)}></Input>
        <RangePicker showTime onChange={handleChange} defaultValue={[dayjs().subtract(1, 'day'), dayjs()]} />
        <Button type="primary"
          loading={errorCallLoading}
          onClick={errorCallback}>生成一条错误回调</Button>
      </Space>
      <Table dataSource={tableData}
        columns={columns}
        pagination={{ position: ['bottomCenter'] }}
        loading={tableLoading}
        expandable={{
          expandedRowRender: (record) => (
            <p style={{ margin: 0 }} >
              {record.message}
            </p>
          )
        }}
      />
    </>
  )
}