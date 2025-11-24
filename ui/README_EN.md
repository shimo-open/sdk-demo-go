# UI

1. node: 16.x
2. antd: <= 4.20.x
3. umi: 4.x

# Running the Project

- Start the corresponding backend project sdk-demo-go locally, and modify the proxy property to point to the backend project address

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

- Start the project

```
npm run dev
```

