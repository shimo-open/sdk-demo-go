import emailRegex from 'email-regex'

const reEmail = emailRegex({ exact: true })

export const validate = {
  email,
  password,
  teamName,
  departmentName
}

function email(input: string) {
  if (input && !reEmail.test(input)) {
    return Promise.reject('请输入有效邮箱');
  }
  return Promise.resolve();
}

function password(input: string) {
  if (input && !/^.{8,32}$/.test(input)) {
    return Promise.reject('请输入 8-32 位的密码');
  }
  return Promise.resolve();
}

function teamName(input: string) {
  if (input && !/^.{1,20}$/.test(input)) {
    return Promise.reject('请输入 8-32 位的密码');
  }
  return Promise.resolve();
}

function departmentName(input: string) {
  if (input && !/^.{1,30}$/.test(input)) {
    return Promise.reject('请输入 8-32 位的密码');
  }
  return Promise.resolve();
}
