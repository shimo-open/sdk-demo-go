import { TeamOutlined, BellFilled, SwapOutlined, SlidersOutlined, BellOutlined, LogoutOutlined, FundOutlined, DatabaseOutlined } from '@ant-design/icons'
import { Dropdown, Popover } from "antd"
import { Link } from 'umi'
import { useState } from 'react'
import { LanguageModal } from './LanguageModal'
import styles from './UserMenu.module.less'
import type { MenuProps } from 'antd'
import { userService } from '@/services/user.service'

export function UserMenu(props: any) {
  const { user, lang } = props
  const [openLang, setOpenLang] = useState(false)

  function logout() {
    userService.logout()
  }

  const userClick: MenuProps['onClick'] = ({ key }) => {
    switch (key) {
      case 'log-out':
        logout()
        break;
      case 'switch-lang':
        setOpenLang(true)
    }
  };

  const items: MenuProps['items'] = [
    {
      key: 'team-info',
      label: (
        <Link to={`/team`} className={styles.linkWrapper}>
          <span>团队信息</span>
        </Link>
      ),
      icon: <TeamOutlined />
    },
    {
      key: 'switch-lang',
      label: (
        <span>切换语言</span>
      ),
      icon: <SwapOutlined />
    },
    {
      key: 'system-info',
      label: (
        <Link to={`/system-info`} className={styles.linkWrapper}>
          <span>系统信息</span>
        </Link>
      ),
      icon: <FundOutlined />
    },
    {
      key: 'system-event',
      label: (
        <Link to={`/system-messages`} className={styles.linkWrapper}>
          <span>系统通知记录</span>
        </Link>
      ),
      icon: <BellOutlined />
    },
    {
      key: 'test-api',
      label: (
        <Link to={`/test-api`} className={styles.linkWrapper}>
          <span>系统测试</span>
        </Link>
      ),
      icon: <SlidersOutlined />
    },
    {
      key: 'ai-knowledge-base',
      label: (
        <Link to={`/knowledge-base`} className={styles.linkWrapper}>
          <span>AI 知识库</span>
        </Link>
      ),
      icon: <DatabaseOutlined />,
    },
    {
      key: 'log-out',
      icon: <LogoutOutlined />,
      label: '登出'
    }
  ]

  return (
    <>
      <div className={styles.eventsWrap}>
        <Link
          to={`/events`}
          className={styles.linkWrapper}
        >
          <Popover content='事件列表'>
            <BellFilled />
          </Popover>
        </Link>
        <Dropdown menu={{ items, onClick: userClick }} trigger={['click']} placement="bottomRight">
          <span className={styles.eventsIcon}>
            <img src={user?.avatar} />
            {user?.name || 'shimo'}
          </span>
        </Dropdown>
      </div>
      <LanguageModal open={openLang} currentLang={lang.current} onClose={() => setOpenLang(false)} />
    </>
  )
}
