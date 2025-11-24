import { getFileAnchor } from './fileLink'

export function Discussion({ data: eventData, files, users }: { data: any, files: any, users: any }) {

  function renderByData(event: any) {
    let action = ''
    let id = ''
    let content = ''
    let time: any = ''

    const currentUserId = Number(event.userId)
    const file = files ? files[event.fileId] || { name: '不存在文件', deleted: true } : { name: '不存在文件', deleted: true }

    switch (event.action) {
      case 'create':
        action = '发送'
        id = event.discussion.id
        content = event.discussion.content
        time = new Date(event.discussion.unixus / 1000)
        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "不存在用户"} 在 `}
              {getFileAnchor`${event.fileId}`(file)}
              {` 中${action}了讨论消息: “${content}”`}
            </p>
            <p>讨论 ID: {id}</p>
            <p>发送时间: {time.toLocaleString()}</p>
          </div>
        )
      case 'like':
        action = '点赞'
        id = event.likeDiscussion.id
        content = event.likeDiscussion.content
        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "不存在用户"} 在`}
              {getFileAnchor`${event.fileId}`(file)}
              {` 中${action}了讨论消息: “${content}”`}
            </p>
            <p>讨论 ID: {id}</p>
          </div>
        )
      default:
        return <div>无法识别的讨论类型</div>
    }
  }

  return <div>{renderByData(eventData)}</div>
}
