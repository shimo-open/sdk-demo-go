import { getFileAnchor } from './fileLink'
import { parseDate } from '@/utils/index'

export function DateMention({ data: eventData, files, users }: { data: any, files: any, users: any }) {
  function renderByData(event: any) {
    switch (event.action) {
      case 'create':
        const file = files ? files[event.createData.fileId] || {
          name: '不存在文件',
          deleted: true
        } : {
          name: '不存在文件',
          deleted: true
        }
        const createRemindUserIds = event.createData.remindUserIds || []
        const userIdsStr = createRemindUserIds
          .map((id: any) => users ? (users[Number(id)] || { name: '不存在用户' }).name : "不存在用户")
          .join(', ')

        return (
          <div>
            <p>
              {`${users ? (
                users[Number(event.createData.authorId)] || {
                  name: '不存在用户'
                }
              ).name : "不存在用户"
                } 在`}
              {getFileAnchor`${event.createData.fileId}`(file)}
              {`中创建了日期提醒#${event.createData.id}`}
            </p>
            <p>{`到期时间: ${parseDate(event.createData.remindAt)}`}</p>
            <p>{`提醒用户：${userIdsStr}`}</p>
          </div>
        )
      case 'update':
        const updateRemindUserIds = event.updateData.remindUserIds || []
        const userIdsStr1 = updateRemindUserIds
          .map((id: any) => users ? (users[Number(id)] || { name: '不存在用户' }).name : "不存在用户")
          .join(', ')

        return (
          <div>
            <p>{`日期提醒#${event.updateData.id} 发生变更`}</p>
            <p>{`到期时间: ${parseDate(event.updateData.remindAt)}`}</p>
            <p>{`提醒用户：${userIdsStr1}`}</p>
          </div>
        )
      case 'remove':
        return (
          <div>
            <p>{`日期提醒#${event.removeData.id} 被删除`}</p>
          </div>
        )
      default:
        return <div>无法识别的日期提醒事件</div>
    }
  }

  return <div>{renderByData(eventData)}</div>
}
