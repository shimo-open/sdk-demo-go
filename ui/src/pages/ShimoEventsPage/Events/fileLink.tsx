export function getFileAnchor(strings: any, ...keys: any) {
  return (file: any) => {
    if (file.deleted) {
      return <a>{`${file.name} `}</a>
    }
    let prefix = "/"

    if (process.env.NODE_ENV === "production" && process.env.BASE) {
      prefix = process.env.BASE
    }

    const result = [`${prefix}shimo-files/`]

    for (let index = 0; index < strings.length; index++) {
      result.push(strings[index])
      if (index < keys.length) {
        result.push(keys[index])
      }
    }

    return (
      <a target="_blank" href={result.join('')} rel="noreferrer">
        {`${file.name} `}
      </a>
    )
  }
}
