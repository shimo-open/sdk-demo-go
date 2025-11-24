import { defineConfig } from "umi";

export default defineConfig({
  publicPath: process.env.PUBLIC_PATH || '/',
  base: process.env.BASE || '/',
  define: {
    'process.env': process.env
  },
  routes: [
    { path: '/', redirect: '/home' },
    { path: '/login', component: '@/pages/LoginPage/LoginPage.tsx', layout: false },
    { path: '/home', component: '@/pages/HomePage/HomePage.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/team', component: '@/pages/TeamPage/Team.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/shimo-files/:id', component: '@/pages/File/ShimoFile/ShimoFile.tsx', wrappers: ['@/wrappers/auth'], layout: false },
    { path: '/form/:id/fill', component: '@/pages/File/ShimoFile/FormFile.tsx', layout: false },
    { path: '/preview/:id', component: '@/pages/File/PreviewFile/PreviewFile.tsx', wrappers: ['@/wrappers/auth'], layout: false },
    { path: '/events', component: '@/pages/ShimoEventsPage/ShimoEventsPage.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/system-info', component: '@/pages/AppDetailPage/AppDetailPage.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/system-messages', component: '@/pages/SystemMessagePage/SystemMessagePage.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/test-api', component: '@/pages/TestApi/TestApi.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/knowledge-base', component: '@/pages/KnowledgeBase/KnowledgeBase.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/knowledge-base/:knowledgeBaseGuid', component: '@/pages/KnowledgeBase/KnowledgeBaseDetail.tsx', wrappers: ['@/wrappers/auth'] },
    { path: '/knowledge-base/:knowledgeBaseGuid/ai', component: '@/pages/KnowledgeBase/AIKnowledgeBase.tsx', wrappers: ['@/wrappers/auth'] },
  ],
  npmClient: 'pnpm',
  proxy: {
    '/api': {
      target: 'http://localhost:9001',
      changeOrigin: true,
      pathRewrite: { '^/api': '' },
    },
  },
  plugins: ['@umijs/plugins/dist/model'],
  model: {},
  metas: [
    {
      name: 'version',
      content: process.env.CI_COMMIT_SHORT_SHA,
    },
  ],
});
