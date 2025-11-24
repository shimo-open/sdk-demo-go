import { getFileAnchor } from './fileLink'

export function Revision({ data: eventData, files, users }: { data: any, files: any, users: any }) {

  function renderByData(event: any) {
    const currentUserId = Number(event.userId)
    const u = users ? users[currentUserId] || { name: '不存在用户' } : { name: '不存在用户' }
    const file = files ? files[event.fileId] || { name: '不存在文件', deleted: true } : { name: '不存在文件', deleted: true }

    switch (event.action) {
      case 'create':
        return (
          <div>
            <p>
              {`${u.name} 为 `}
              {getFileAnchor`${event.fileId}`(file)}
              {` 创建了新版本：${event?.revision?.title ?? '未知版本标题'}`}
            </p>
          </div>
        )
      case 'delete':
        return (
          <div>
            <p>
              {`${u.name} 删除了 `}
              {getFileAnchor`${event.fileId}`(file)}
              {` 的版本：${event?.revision?.title ?? '未知版本标题'}`}
            </p>
          </div>
        )
      default:
        return <div>无法识别的协作者状态事件</div>
    }
  }

  return <div>{renderByData(eventData)}</div>
}
