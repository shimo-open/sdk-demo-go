import { Link } from 'react-router-dom'

import { otherConstants } from '@/constants/other.constants'
import styles from './Header.module.less'
import { UserMenu } from './UserMenu'

export function Header(props: { user: any; lang: any }) {
  const { user, lang } = props
  return (
    <header className={styles.header}>
      <h1 className={styles.h1}>
        <Link to={`/`} className={styles.imgWrapper}>
          <img
            src={otherConstants.LOGO}
            className={styles.logo}
            alt="石墨文档"
          />
        </Link>
      </h1>
      <UserMenu user={user} lang={lang} />
    </header>
  )
}
