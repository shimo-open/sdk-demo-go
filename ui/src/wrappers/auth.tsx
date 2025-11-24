import { Navigate, Outlet } from 'umi'
import { getToken, setToken } from '@/utils/axios'
import { START_PARAMS_FIELD } from 'shimo-js-sdk'
import { parseSmParams } from "@/utils/index"
import { userConstants } from '@/constants'

const AuthRouter = (props: any) => {
  const smParams = new URLSearchParams(location.search)?.get(START_PARAMS_FIELD) || ""
  const transSmParams: any = parseSmParams(smParams)
  // Stay on the page if a token exists or a shimo-files response-share view is requested
  let isLogin = (getToken() || (location.pathname.match(/\/shimo-files/i) && transSmParams.type == "form" && transSmParams.path?.includes("response-share"))) ? true : false;
  // Read the token from the query string
  const token = new URLSearchParams(location.search)?.get(userConstants.ACCESS_TOKEN_KEY) || ""
  if (token !== "") {
    setToken(token)
    isLogin = true
  }
  return (
    isLogin ? <Outlet /> : <Navigate to={`/login`} />
  )
}

export default AuthRouter;

