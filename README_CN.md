# Shimo SDK Go Demo

## 项目简介

这是一个基于 Ego 的 Shimo SDK Go 演示项目，提供了完整的 Shimo 文档协作平台 API 集成示例。项目展示了如何通过 Go 后端服务与 Shimo SDK 进行交互，实现文档管理、用户管理、团队协作等功能。

SDK 对接的具体架构如下图，本项目实现了图中黄色的部分：
<img width="1088" height="516" alt="image" src="https://github.com/user-attachments/assets/546d36c1-52f2-478a-8d06-db50da665294" />

## 项目结构

```
.
├── cmd/                        # 命令行工具
│   ├── server/                 # 服务器启动命令
│   └── sdkctl/                 # SDK 测试工具
├── config/                     # 配置文件（本地/生产环境配置）
├── database/                   # 数据库初始化脚本
├── pkg/                        # 核心业务包
│   ├── consts/                 # 常量定义（文件类型、权限、API 等）
│   ├── invoker/                # 依赖注入容器（全局组件管理）
│   ├── models/                 # 数据模型
│   │   └── db/                 # 数据库模型（用户、文件、团队等）
│   ├── server/                 # 服务器实现
│   │   └── http/               # HTTP 服务
│   │       ├── api/            # API 接口（用户、文件、团队等管理）
│   │       ├── callback/       # Shimo 回调接口实现
│   │       └── middlewares/    # HTTP 中间件（认证、回调验证等）
│   ├── services/               # 业务服务层
│   │   ├── signature/          # JWT 签名服务
│   │   ├── awos/               # 对象存储服务（S3/MinIO）
│   │   └── inspect/            # Web 巡检服务
│   └── utils/                  # 工具函数（JWT、加密、文件处理等）
├── resources/                  # 资源文件
│   └── import/                 # 导入测试文件（各种格式）
├── scripts/                    # 脚本文件
│   └── build/                  # 构建脚本
├── ui/                         # 前端项目（React + TypeScript）
│   ├── public/                 # 静态资源
│   └── src/                    # 源代码
│       ├── constants/          # 前端常量
│       ├── layouts/            # 页面布局
│       ├── models/             # 前端数据模型
│       ├── pages/              # 页面组件
│       ├── services/           # API 服务调用
│       ├── utils/              # 前端工具函数
│       └── wrappers/           # 路由包装器
├── Dockerfile                  # Docker 镜像构建文件
├── Makefile                    # Make 构建配置
├── go.mod                      # Go 模块定义
├── go.sum                      # Go 依赖校验和
├── main.go                     # 程序入口
├── README.md                   # 项目文档（中文）
└── README_EN.md                # 项目文档（英文）
```


## 前端代码地址

[ui/](ui/)

## 前端技术栈

- React
- Typescript
- Ant-Design

## 后端技术栈

- **框架**: Ego、Gin
- **数据库**: MySQL
- **认证**: JWT
- **文件存储**: AWS S3/MinIO

## 数据库设计

### ER 图

![ER Diagram](docs/ER.webp)

### 主要数据表说明

- **users**: 用户表，存储系统用户信息
- **teams**: 团队表，组织的基本单位
- **team_roles**: 团队角色表，用户在团队中的角色（创建者/管理员/成员）
- **departments**: 部门表，团队下的部门组织结构
- **dept_members**: 部门成员表，用户与部门的关联关系
- **files**: 文件表，存储协同文档和上传文件
- **file_permissions**: 文件权限表，用户对文件的访问权限
- **events**: 事件表，记录系统中的各类事件
- **knowledge_bases**: 知识库表，知识库相关信息
- **app_clients**: 应用客户端表，存储应用凭证
- **test_api**: API 测试记录表，记录 API 测试结果

## 数据库初始化

请执行 [database/1_init.up.sql](database/1_init.up.sql) 文件中的 SQL 语句来初始化数据库。

```bash
mysql -u your_username -p < database/1_init.up.sql
```

或者直接在 MySQL 客户端中执行该文件。

## 服务启动方式

### 后端启动

#### 使用 Makefile 启动（推荐）

```bash
make server
```

#### 直接使用 Go 命令启动

```bash
go run main.go server --config=config/local.toml
```

配置文件：[config/local.toml](config/local.toml)

#### 构建二进制文件

```bash
make build
```

构建完成后，可执行文件位于 `bin/sdk-demo-go` 目录下。

### 前端启动

#### 环境要求

- Node.js: 16.x
- Ant Design: <= 4.20.x
- Umi: 4.x

#### 启动步骤

1. 进入前端目录

```bash
cd ui
```

2. 安装依赖

```bash
npm install
```

3. 配置代理（如需要）

确保 [ui/.umirc.ts](ui/.umirc.ts) 或相关配置文件中 proxy 属性指向后端服务地址（默认：`http://localhost:9301`）

```json
{
  "proxy": {
    "/api": {
      "target": "http://localhost:9301",
      "changeOrigin": true,
      "pathRewrite": { "^/api": "" }
    }
  }
}
```

4. 启动开发服务器

```bash
npm run dev
```

#### 构建生产版本

```bash
npm run build
```

## API 接口说明

### 认证接口

- `POST /api/sign` - 获取 JWT 令牌
- `POST /api/signByClientName` - 通过客户端名称获取 JWT

### 用户管理

- `POST /api/users/signin` - 用户登录
- `POST /api/users/signup` - 用户注册
- `GET /api/users/{userId}` - 获取用户信息
- `GET /api/users/{userId}/teams` - 获取用户团队

### 文件管理

- `GET /api/files` - 获取用户文件列表
- `POST /api/files/upload` - 上传文件
- `POST /api/files/import` - 导入文件
- `POST /api/files/{fileGuid}/export` - 导出文件
- `GET /api/files/{fileGuid}/open` - 打开文件

### 团队管理

- `GET /api/teams` - 获取团队列表
- `POST /api/teams` - 创建团队
- `GET /api/teams/{teamId}/members` - 获取团队成员
- `POST /api/teams/{teamId}/departments` - 创建部门

### 应用管理

- `GET /api/apps/{appId}` - 获取应用详情
- `PUT /api/apps/{appId}/endpoint-url` - 更新应用回调地址

## 常见问题

### 1. 数据库连接失败

- 检查数据库服务是否启动
- 验证连接字符串和凭据
- 确认网络连通性

### 2. Shimo SDK 调用失败

- 验证 AppID 和 AppSecret 是否正确
- 验证回调是否正确
- 检查网络连接和防火墙设置
- 查看 SDK 日志获取详细错误信息

### 3. 文件上传失败

- 检查文件存储服务配置
- 验证文件大小限制
- 确认存储权限设置
