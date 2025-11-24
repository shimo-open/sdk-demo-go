
import { Dropdown, Popover } from "antd"
import type { MenuProps } from 'antd'
import { fileConstants } from '@/constants'
import { apis } from '@/utils/apis'
import mobile from 'is-mobile'
import { BarsOutlined } from '@ant-design/icons'
import { useEffect, useState } from "react"
const isMobile = mobile()

export function FileHeaderMenu({ file, onMenuItemClick }: { file: any, onMenuItemClick: (key: string) => void }) {
  const [options, setOption] = useState<MenuProps['items']>([])

  useEffect(() => {
    init()
  }, [])

  const init = () => {
    let ops: MenuProps['items'] = []
    if (file.shimoType === fileConstants.TYPE_DOCUMENT) {
      ops.push(
        ...[
          {
            key: 'showToc',
            label: '查看目录',
          },
          {
            key: 'showHistory',
            label: '查看历史',
          },
          {
            key: 'createRevision',
            label: '创建版本',
          },
          {
            key: 'showRevision',
            label: '查看版本',
          },
          {
            key: 'showDiscussion',
            label: '查看讨论',
          },
          {
            key: 'startDemonstration',
            label: '演示模式',
          },
          {
            key: 'print',
            label: '打印',
          }
        ]
      )
    } else if (file.shimoType === fileConstants.TYPE_DOCUMENT_PRO) {
      ops.push(
        ...[
          {
            key: 'showToc',
            label: '查看目录',
          },
          {
            key: 'showHistory',
            label: '查看历史',
          },
          {
            key: 'printAll',
            label: '打印所有页面',
          }
        ]
      )
    } else if (file.shimoType === fileConstants.TYPE_SPREADSHEET) {
      ops.push(
        ...[
          {
            key: 'showComments',
            label: '查看评论',
          },
          {
            key: 'showHistory',
            label: '查看历史',
          },
          {
            key: 'createRevision',
            label: '创建版本',
          },
          {
            key: 'startDemonstration',
            label: '演示模式',
          },
          {
            key: 'print',
            label: '打印',
          }
        ]
      )
    } else if (file.shimoType === fileConstants.TYPE_PRESENTATION) {
      ops.push(
        ...[
          {
            key: 'showHistory',
            label: '查看历史',
          },
          {
            key: 'startDemonstration',
            label: '演示模式',
          }
        ]
      )
    } else if (file.shimoType === fileConstants.TYPE_TABLE) {
      ops.push(
        ...[
          {
            key: 'showRevision',
            label: '查看历史',
          },
          {
            key: 'createRevision',
            label: '创建版本',
          }
        ]
      )
    } else if (file.shimoType === fileConstants.TYPE_FLOWCHART) {
      ops.push(
        ...[
          {
            key: 'showRevision',
            label: '查看版本',
          }
        ]
      )
    }

    ops = ops.filter((o: any) => {
      for (const api of editorApis) {
        for (const item of api.methods) {
          if (item.method === o.key && (!isMobile || item['support mobile'] === 'Y')) {
            return true
          }
        }
      }
      return false
    })

    setOption(ops)
  }

  const editorApis = apis.apis.filter((api) =>
    api.sdk === file.shimoType ||
    api.sdk === 'all' ||
    (file.shimoType === fileConstants.TYPE_SPREADSHEET && api.sdk === 'sheet')
  )

  const itemClicked: MenuProps['onClick'] = ({ key }) => {
    onMenuItemClick(key)
  };

  return (
    <Dropdown menu={{ items: options, onClick: itemClicked }} trigger={['click']}>
      <BarsOutlined className="triggerIcon" style={{ display: options?.length ? "" : "none" }} />
    </Dropdown>
  )
}