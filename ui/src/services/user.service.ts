import { history } from "umi";
import { userConstants } from '../constants'
import req, { getUrlPrefix } from '@/utils/axios';

interface LoginApiType {
  email: string;
  password: string;
}

export const userService = {
  login,
  getById,
  register,
  logout,
  getAll,
  auth
}

async function login(data: LoginApiType) {
  const user = await req.post(`/api/users/signin`, data)
  return user
}

async function getById(id: any) {
  return await req.get(`/api/users/${id}`)
}

async function register(data: LoginApiType) {
  return await req.post(`/api/users/signup`, data)
}

async function getAll() {
  return await req.get(`/api/users`)
}

function logout() {
  localStorage.removeItem(userConstants.TOKEN_STORAGE_KEY)
  localStorage.removeItem(userConstants.LOGIN_KEY)
  const currentPath = window.location.pathname
  // Remove the BASE prefix
  const basePrefix = getUrlPrefix()
  const cleanPath = currentPath.replace(basePrefix, '/')
  history.push(`/login?redirect=${cleanPath}`)
  window.location.reload()
}

async function auth() {
  return await req.post(`api/users/auth`)
}