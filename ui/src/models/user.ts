import { userConstants } from '@/constants'
import { useState } from 'react';

export interface initialState {
  lang: any,
  userInfo: any
}
const userModel = () => {
  const [lang, setLang] = useState<any>({
    current: localStorage.getItem(userConstants.LANGUAGE_STORAGE_KEY) || userConstants.LANGUAGE_LOCAL
  })
  const [userInfo, setUserInfo] = useState<any>({})

  const [hostPath, setHostPath] = useState("")

  return {
    userInfo,
    lang,
    setUserInfo,
    setLang,
    hostPath,
    setHostPath
  }
}

export default userModel;