import { notification } from 'antd';

const Notification: React.ReactNode = (
  <a onClick={() => {
    notification.destroy();
  }} style={{ color: "#6a6868" }}>
    全部关闭
  </a>
)

export default Notification