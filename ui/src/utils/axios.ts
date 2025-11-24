import axios from 'axios'
import { message } from 'antd'
import { userConstants } from '@/constants'
import { userService } from '@/services/user.service';
import { notification } from 'antd';
import Notification from '@/utils/notification'

type NotificationType = 'success' | 'info' | 'warning' | 'error';

export const getToken = () => localStorage.getItem(userConstants.TOKEN_STORAGE_KEY) || '';

export function getUrlPrefix(noPrefix = false) {
  const prefix = noPrefix ? '' : process.env.BASE || ''
  return prefix.endsWith('/') ? prefix : prefix + '/'
}

export const setToken = (token: string) => localStorage.setItem(userConstants.TOKEN_STORAGE_KEY, token)

const instance = axios.create({
  baseURL: process.env.NODE_ENV === 'development' ? '/api' : process.env.BASE ? process.env.BASE.endsWith('/') ? process.env.BASE.substring(0, process.env.BASE.length - 1) : process.env.BASE : '',
  withCredentials: true,  // Include cookies for cross-origin requests
  timeout: 1800000, // Five minutes = 300000; thirty minutes = 1800000
});

// Request interceptor
instance.interceptors.request.use(function (config: any) {
  // Inject auth headers before the request is sent
  config.headers = {
    'Authorization': `bearer ${getToken()}`,
  }
  return config;
}, function (error) {
  // Handle request errors
  return Promise.reject(error);
});

// Response interceptor
instance.interceptors.response.use(function (response) {
  return response;
}, function (error) {
  // Handle response errors
  const { response } = error;
  if (response) {
    let message = 'Request failed'
    let description = ''
    switch (response.status) {
      case 401:
        description = 'Unauthorized, please sign in again'
        userService.logout()
        break;
      case 504:
        message = 'Gateway timeout'
        description = response.data
        break;
    }
    openNotificationWithIcon('error', message, description || response.data.message || (typeof response.data === 'string' ? response.data : ''));
  }
  return Promise.reject(response);
  // return response ? response : { data: null };
});

const openNotificationWithIcon = (type: NotificationType, message: string, description: string) => {
  notification[type]({
    message: message,
    description: description,
    duration: 2.5,
    style: {
      width: 'max-content',
      maxWidth: 'calc(45vw)',
      minWidth: '260px',
      opacity: 0.85
    },
    btn: Notification
  });
};

export default instance


