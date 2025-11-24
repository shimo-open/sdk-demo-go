import { useEffect } from 'react'
import { useModel } from 'umi'
import { CopyrightFooter } from '../CopyrightFooter/Footer'
import { FileActionMenu } from '../File/FileActionMenu'
import { FileList } from '../File/FileList'
import styles from './HomePage.module.less'
import { userConstants } from '@/constants'
export default function HomePage() {
  const { lang, setLang } = useModel('user')
  useEffect(() => {
    let valueInStorage: any = ''
    try {
      valueInStorage = localStorage.getItem(userConstants.LANGUAGE_STORAGE_KEY)
    } catch (err) {
      console.warn(`read language value from localStorage error`)
    }
    if (valueInStorage) {
      setLang({ current: valueInStorage })
    }
  }, [])

  return (
    <>
      <FileActionMenu className={styles.fileActionMenu} />
      <FileList className={styles.list} />
      <CopyrightFooter />
    </>
  )
}