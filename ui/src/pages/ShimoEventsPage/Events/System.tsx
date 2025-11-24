export function System({ data }: { data: any }) {
  const eventData = data
  if (eventData.type === 'endpointCallback') {
    return <p>回调接口出错</p>
  } else {
    return <p>License即将过期</p>
  }
}
