import { getFileAnchor } from './fileLink'

export function Comment({ data: eventData, files, users }: { data: any, files: any, users: any }) {

  function renderByData(event: any) {
    const currentUserId = Number(event.userId)
    const file = files ? files[event.fileId] || { name: '不存在文件', deleted: true } : { name: '不存在文件', deleted: true }

    switch (event.action) {
      case 'create':
        return (
          <div>
            <p>
              {`${users[currentUserId].name} 在 `}
              {getFileAnchor`${event.fileId}?mentionId=${event.comment.selectionGuid}`(
                file
              )}
              {`中添加了评论 `}
              <br />
              {`"${event.comment.content}"`}
            </p>
            {event.comment.userIds &&
              Array.isArray(event.comment.userIds) &&
              event.comment.userIds.length > 0 ? (
              <p>{`提到用户 ${event.comment.userIds
                .map((i: any) => users ? (users[Number(i)] || { name: '不存在用户' }).name : "不存在用户")
                .join(', ')}`}</p>
            ) : null}
          </div>
        )
      case 'update':
        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "不存在用户"} 在 `}
              {getFileAnchor`${event.fileId}?mentionId=${event.comment.selectionGuid}`(
                file
              )}
              {`中更新了评论: `}
              <br />
              {`"${event.comment.content}"`}
            </p>
            {event.comment.userIds &&
              Array.isArray(event.comment.userIds) &&
              event.comment.userIds.length > 0 ? (
              <p>{`提到用户 ${event.comment.userIds
                .map((i: any) => users ? (users[Number(i)] || { name: '不存在用户' }).name : "不存在用户")
                .join(', ')}`}</p>
            ) : null}
          </div>
        )
      case 'delete':
        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "不存在用户"} 在 `}
              {getFileAnchor`${event.fileId}?mentionId=${event.deleteComment.selectionGuid}`(
                file
              )}
              {`中删除了评论`}
            </p>
          </div>
        )
      case 'closeComments':
        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "users"} 在 `}
              {getFileAnchor`${event.fileId}`(file)}
              {`中结束了评论`}
            </p>
          </div>
        )
      case 'like':
        return (
          <div>
            <p>
              {`${users ? users[currentUserId].name : "users"} 点赞了 `}
              {getFileAnchor`${event.fileId}?mentionId=${event.likeComment.selectionGuid}`(file)}
              {`的评论`}
              <br />
              {`"${event.likeComment.content}"`}
            </p>
          </div>
        )
      default:
        return <div>无法识别的评论类型</div>
    }
  }

  return <div>{renderByData(eventData)}</div>
}
