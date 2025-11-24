# UI

1. node： 16.x
2. antd: <= 4.20.x
3. umi: 4.x

# 项目运行

- 本地需启动对应后端项目 sdk-demo-go, 将 proxy 属性修改为对应后端项目地址

```
proxy: {
  '/api': {
    target: 'http://localhost:9301',
    changeOrigin: true,
    pathRewrite: { '^/api': '' },
  },
}
```

`npm install`

- 启动项目

```
npm run dev
```

`
