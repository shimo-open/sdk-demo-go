import styles from './CopyrightFotter.module.less'

export function CopyrightFooter(props:any) {
  const year = new Date().getFullYear()

  return (
    <footer className={[props.className, styles.footer].join(' ')}>
      <span>
        © 2014-{year} 武汉初心科技有限公司
        <br />
        本程序仅用于测试演示目的，石墨文档拥有和保留一切权利，且不对本程序安全性、稳定性等方面作任何保证。
      </span>
    </footer>
  )
}
