import req from '@/utils/axios'

export const systemMessageService = {
  getSystemMessages,
  errorCallback
}

async function getSystemMessages(appID: string, from: string, to: string) {
  return await req.get(`api/events/system-messages`, {
    params: {
      appID: appID,
      from: from,
      to: to
    }
  })
}

async function errorCallback() {
  return await req.get('api/events/error_callback')
}
