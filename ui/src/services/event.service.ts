import req from '@/utils/axios';

export const eventService = {
  getEvents
}

async function getEvents(params: { page: number, size: number, fileId: string | number }) {
  return await req.get(`api/events`, {
    params: params
  })
}