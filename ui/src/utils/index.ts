import baseX from "base-x"

export function pick(input: Object, keys: any) {
  const obj: any = {}
  if (input != null) {
    for (const key of Array.isArray(keys) ? keys : []) {
      if (input.hasOwnProperty(key)) {
        obj[key] = input[key as keyof typeof input]
      }
    }
  }
  return obj
}

export function parseSmParams(smParams: string) {
  if (!smParams) return ""
  // Create a Base62 decoder using '0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz'
  const base62 = baseX('0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz');
  // base62Encoded is the serialized Base62 string
  const base62Encoded = smParams;
  // Decode the string
  const decodedBuffer = base62.decode(base62Encoded);
  const enc = new TextDecoder("utf-8");
  const uint8_msg = new Uint8Array(decodedBuffer);
  const decodedJsonParse = JSON.parse(enc.decode(uint8_msg));
  return decodedJsonParse
}

// Convert UTC timestamps to local time (Shanghai)
export function parseDate(date: string) {
  if (!date) return ""
  let utcDate = new Date(date)
  const localTime = utcDate.toLocaleString()
  return localTime.replace(/\//g, '-')
}
