import { getFileAnchor } from './fileLink'

export function Collaborator({ data: eventData, files, users }: { data: any, files: any, users: any }) {
  function renderByData(event: any) {
    const time = new Date(event.timestamp)
    const currentUserId = Number(event.userId)
    const u = users[currentUserId] || { name: '不存在用户' }
    const file = files ? files[event?.fileId] || { name: '不存在文件', deleted: true } : { name: '不存在文件', deleted: true }

    switch (event.action) {
      case 'enter':
        return (
          <div>
            <p>
              {`${u.name} ${time.toLocaleString()} 进入协作 `}
              {getFileAnchor`${event.fileId}`(file)}
            </p>
          </div>
        )
      case 'leave':
        return (
          <div>
            <p>
              {`${u.name} ${time.toLocaleString()} 退出协作 `}
              {getFileAnchor`${event.fileId}`(file)}
            </p>
          </div>
        )
      default:
        return <div>无法识别的协作者状态事件</div>
    }
  }

  return <div>{renderByData(eventData)}</div>
}
