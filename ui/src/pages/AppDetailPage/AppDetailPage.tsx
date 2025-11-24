import { Card, Button, message, Divider, Space } from "antd";
import { appService } from "@/services/app.service";
import { useEffect, useState } from "react";
import styles from './AppDetail.module.less'


export default function AppDetailPage() {
  const [detail, setDetail] = useState({} as any)

  const appName = detail && detail.appName ? detail.appName : '-'
  const fileTypes =
    detail && detail.availableFileTypes ? detail.availableFileTypes : []
  const activatedCount =
    detail && !isNaN(detail.activatedUserCount) ? detail.activatedUserCount : 0
  const userCount = detail && !isNaN(detail.userCount) ? detail.userCount : 0
  const memberLimit =
    detail && !isNaN(detail.memberLimit) ? detail.memberLimit : 0
  const validFrom =
    detail && detail.validFrom
      ? new Date(detail.validFrom).toLocaleString()
      : '-'
  const validUntil =
    detail && detail.validUntil
      ? new Date(detail.validUntil).toLocaleString()
      : '-'
  const endpointUrl = detail && detail.endpointUrl ? detail.endpointUrl : '-'

  useEffect(() => {
    appService.getAppDetail().then(res => {
      setDetail(res.data)
    })
  }, [])

  function doUpdateEndpointUrl() {
    appService.updateEndpointUrl().then(res => {
      message.success('更新 endpoint url 成功')
    })
  }

  return (
    <div className={styles.infoWrapper}>
      <Card>
        <h3>应用名称</h3>
        <span>{appName}</span>
        <Divider></Divider>
        <h3>席位信息</h3>
        <Space>
          <div className={styles.title}>
            <span>{activatedCount}</span>
            <span>已激活数</span>
          </div>
          <div className={styles.title}>
            <span>{userCount}</span>
            <span>用户总数</span>
          </div>
          <div className={styles.title}>
            <span>{memberLimit}</span>
            <span>席位总数</span>
          </div>
        </Space>
        <Divider></Divider>
        <h3>可用套件</h3>
        <span>{fileTypes.join('、')}</span>
        <Divider></Divider>
        <h3>License 有效期</h3>
        <span>{`${validFrom} 至 ${validUntil}`}</span>
        <Divider></Divider>
        <h3>回调地址</h3>
        <span>{endpointUrl}</span>
      </Card>
      <Card>
        <Button onClick={doUpdateEndpointUrl}>测试更新回调地址</Button>
      </Card>
    </div>
  )
}