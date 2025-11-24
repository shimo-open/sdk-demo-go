import { useEffect } from 'react'
import { Outlet, useModel } from 'umi';
import { ConfigProvider } from 'antd';
import zhCN from 'antd/locale/zh_CN';
import { Header } from '@/pages/HeaderPage/Header'
import { userService } from '@/services/user.service';
import { userConstants } from "@/constants/user.constants";
import { setToken } from '@/utils/axios'

export default function Layout() {
  const { userInfo, lang, setUserInfo, setHostPath } = useModel("user");
  const { setAvailableFileTypes, setPremiumFileTypes, setAllowBoardAndMindmap } = useModel("file")

  useEffect(() => {
    // Only fetch on page refresh (skip during explicit login)
    let loginFlag = localStorage.getItem(userConstants.LOGIN_KEY)
    if ((!userInfo || Object.keys(userInfo).length == 0) && !loginFlag) {
      userService.auth().then((res) => {
        if (res.data.user.id < 0) {
          userService.logout()
        } else {
          const { data: { user, availableFileTypes, premiumFileTypes, hostPath, token, allowBoardAndMindmap } } = res;
          setUserInfo(user)
          setAvailableFileTypes(availableFileTypes)
          setPremiumFileTypes(premiumFileTypes)
          setAllowBoardAndMindmap(allowBoardAndMindmap)
          setHostPath(hostPath)
          setToken(token)
        }
      }, (error) => {
        userService.logout()
      })
    }
    // Remove the login flag so a refresh will trigger the fetch again
    localStorage.removeItem(userConstants.LOGIN_KEY)
  }, [userInfo])

  return (
    <ConfigProvider locale={zhCN}>
      <Header user={userInfo} lang={lang} />
      <main style={{ flex: 1 }}>
        <Outlet />
      </main>
    </ConfigProvider>
  );
}