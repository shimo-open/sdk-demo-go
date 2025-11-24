import { getFileAnchor } from './fileLink'

export function MentionAt({ data: eventData, files, users }: { data: any, files: any, users: any }) {

  function renderByData(event: any) {
    const currentUserId = Number(event.userId)
    const file = files ? files[event.fileId] || { name: '不存在文件', deleted: true } : { name: '不存在文件', deleted: true }

    switch (event.type) {
      case 'comment': {
        const userIdsComment = event.comment.userIds || []
        const usersStrComment = userIdsComment
          .map((i: any) => users ? (users[Number(i)] || { name: '不存在用户' }).name : "不存在用户")
          .join(', ')

        return (
          <div>
            <p>
              {`${users[currentUserId].name} 在 `}
              {getFileAnchor`${event.fileId}?mentionId=${event.comment.selectionGuid}`(
                file
              )}
              {`的评论中提到了 ${usersStrComment}`}
            </p>
            <p>评论消息: "{event.comment.content}"</p>
          </div>
        )
      }

      case 'discussion': {
        const userIdsDiscussion = event.discussion.userIds || []
        const usersStrDiscussion = userIdsDiscussion
          .map((i: any) => users ? (users[Number(i)] || { name: '不存在用户' }).name : "不存在用户")
          .join(', ')

        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "不存在用户"} 在 `}
              {getFileAnchor`${event.fileId}`(file)}
              {` 的讨论中提到了 ${usersStrDiscussion}`}
            </p>
            <p>讨论消息: "{event.discussion.content}"</p>
          </div>
        )
      }

      case 'mention_at': {
        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "不存在用户"} 在 `}
              {getFileAnchor`${event.fileId}?mentionId=${event.mentionAt.guid}`(
                file
              )}
              {` 中提到了 ${users ? (
                users[Number(event.mentionAt.userId)] || {
                  name: '不存在用户'
                }
              ).name : "不存在用户"
                }`}
            </p>
            <p>提及内容: "{event.mentionAt.content}"</p>
          </div>
        )
      }
      default:
        return <div>无法识别的提及 at 事件</div>
    }
  }

  return <div>{renderByData(eventData)}</div>
}
