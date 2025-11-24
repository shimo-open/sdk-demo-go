import req from '@/utils/axios'

export const appService = {
  getAppDetail,
  updateEndpointUrl
}

async function getAppDetail() {
  return req.get(`api/apps/detail`)
}

async function updateEndpointUrl() {
  return req.put(`api/apps/endpoint-url`)

}
