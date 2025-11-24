import req from '@/utils/axios'

export const apiTestService = {
  allApiTest,
  apiTestProcess,
  getTestApiList
}

async function allApiTest(type: string) {
  return req.get(`api/apiTest/allApiTest`, {
    params: { type: type }
  })
}

async function apiTestProcess(taskId: string) {
  return req.get(`api/apiTest/testProgress`, {
    params: { taskId: taskId }
  })
}

async function getTestApiList() {
  return req.get(`api/apiTest/getTestApiList`)
}

