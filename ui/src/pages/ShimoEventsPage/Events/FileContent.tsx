import { getFileAnchor } from './fileLink'

export function FileContent({ data: eventData, files, users }: { data: any, files: any, users: any }) {
  function renderByData(event: any) {
    const time = new Date(event.timestamp)
    const currentUserId = Number(event.userId)
    const file = files ? files[event?.fileId] || { name: '不存在文件', deleted: true } : { name: '不存在文件', deleted: true }

    switch (event.type) {
      case 'auto_mention':
        return (
          <div>
            <p>
              {`${users ? (users[currentUserId] || { name: '不存在用户' }).name : "不存在用户"
                } ${time.toLocaleString()} 更新了 `}
              {getFileAnchor`${event.fileId}`(file)}
              {` 的关注选区`}
            </p>
          </div>
        )
      default:
        return (
          <div>
            <p>
              {`${users ? (users[currentUserId] || { name: '不存在用户' }).name : "不存在用户"
                } ${time.toLocaleString()} 更新了 `}
              {getFileAnchor`${event.fileId}`(file)}
              {` 的内容，版本 ${event.fileContent.version}`}
            </p>
          </div>
        )
    }
  }
  return <div>{renderByData(eventData)}</div>
}
