import { Button, Select, Space, message, Tabs, Progress } from "antd";
import { fileConstants } from "@/constants"
import { apiTestService } from "@/services/apiTest.service"
import { useEffect, useState } from "react";
import { TestApiResult } from "./TestApiResult";
import { testApiConstants } from "@/constants/testApi.constants"
import styles from "./TestApi.module.less"
import { ApiOutlined } from '@ant-design/icons';

export default function TestApi() {
  const [type, setType] = useState(['all'])
  const [data, setData] = useState({} as any)
  const [loading, setLoading] = useState(false)
  const [percent, setPercent] = useState<number>(0);
  const [activeKey, setActiveKey] = useState('BaseTestResMap')

  const fileType = [
    { value: "all", label: '全部' },
    { value: fileConstants.TYPE_DOCUMENT, label: '轻文档' },
    { value: fileConstants.TYPE_DOCUMENT_PRO, label: '传统文档' },
    { value: fileConstants.TYPE_SPREADSHEET, label: '表格' },
    { value: fileConstants.TYPE_PRESENTATION, label: '专业幻灯片' },
    { value: fileConstants.TYPE_TABLE, label: '应用表格' }
  ]
  // API categories
  const sheetNameMap = testApiConstants.SHEET_NAME_MAP.map((item) => {
    return {
      key: item.KEY,
      label: item.VALUE,
      children: <TestApiResult key={item.KEY} type={item.KEY} data={data ? data[item.KEY] : {}}></TestApiResult>,
      icon: <span className={styles.noteStyle} style={{ background: getColor(item.KEY) }}></span>
    }
  })

  useEffect(() => {
    apiTestService.getTestApiList().then(res => {
      if (res.data) {
        setData(res.data.result)
      }
    })
  }, [])

  function getColor(type: string) {
    if (data[type]) {
      let success = 0
      let failure = 0
      data[type].some((item: any) => {
        if (item.success) {
          success++
        } else {
          failure++
        }
        if (success > 0 && failure > 0) return true
      })
      if (success > 0 && failure > 0) {
        return 'radial-gradient(8px, rgb(255, 197, 66), rgba(255, 197, 66, .6)'
      } else if (success > 0) {
        return 'radial-gradient(8px, rgb(76, 217,100), rgba(76, 217,100,.6))'
      } else if (failure > 0) {
        return 'radial-gradient(8px, rgb(237,71,58), rgba(237,71,58, .6))'
      } else {
        return 'radial-gradient(8px, rgb(153,153,153), rgba(153,153,153, .6)'
      }
    } else {
      return 'radial-gradient(8px, rgb(153,153,153), rgba(153,153,153, .6)'
    }
  }

  // Start the test run
  function StartTestApi() {
    let typeParam = ""
    if (type.length) {
      typeParam = type.join(",")
    } else {
      message.warning("请选择测试文件类型")
      return
    }
    setPercent(0)
    setLoading(true)

    apiTestService.allApiTest(typeParam).then(async res => {
      if (res.status && res.status == 200) {
        localStorage.setItem(testApiConstants.TASK_ID_LOCATION_KEY, res.data.taskId)
        do {
          try {
            const resp = await apiTestService.apiTestProcess(res.data.taskId)
            if (!resp.data || resp.data?.status === 2) {
              throw new Error()
            }
            if (resp.data?.status !== 0) {
              console.log("processing...");
              if (resp.data.progress >= testApiConstants.SHEET_NAME_MAP.length) {
                setActiveKey(testApiConstants.SHEET_NAME_MAP[testApiConstants.SHEET_NAME_MAP.length - 1].KEY)
                setPercent(100)
              } else {
                // Switch tabs based on progress
                resp.data.progress == 0 ? null : setActiveKey(testApiConstants.SHEET_NAME_MAP[Number(resp.data.progress)].KEY)
                setPercent(resp.data.progress == 0 ? 0 : resp.data.progress / testApiConstants.SHEET_NAME_MAP.length * 100)
              }
            } else {
              setData(resp.data.result)
              setPercent(100)
              setLoading(false)
              console.log("测试完成", resp.data);
              break
            }
          } catch (e) {
            message.error("测试失败")
            setLoading(false)
            break
          }
          await new Promise((resolve) => setTimeout(resolve, 1000))
        } while (true)
      } else {
        setLoading(false)
        message.error("测试失败")
      }
    })
  }

  function onChangeType(e: any) {
    if (e.some((v: string) => v == "all")) {
      setType(["all"])
    } else {
      setType(e)
    }
  }

  return (
    <div className={styles.testContanier}>
      <Space className={styles.topSpace}>
        <span className={styles.text}>测试文件类型:</span>
        <Select mode="multiple" placeholder="测试文件类型(可多选)" options={fileType} value={type} onChange={(e) => onChangeType(e)}
          style={{ minWidth: "200px" }} />
        <Button onClick={StartTestApi} type="primary" loading={loading}><ApiOutlined />开始测试</Button>
      </Space>
      <div className={styles.resultItem}>
        <Tabs activeKey={activeKey} items={sheetNameMap} onChange={(e) => setActiveKey(e)} />
        {loading &&
          <Progress
            type="circle"
            percent={percent}
            steps={{ count: 6, gap: 8 }}
            trailColor="rgba(0, 0, 0, 0.06)"
            strokeColor="rgb(0, 0, 0)"
            strokeWidth={20}
            showInfo={false}
            size="small"
          />
        }
      </div>
    </div>
  )
}