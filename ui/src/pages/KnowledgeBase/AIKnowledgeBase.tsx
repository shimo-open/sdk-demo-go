import { useEffect, useState, useRef } from "react";
import { knowledgeBaseService } from "@/services/knowledgeBase.service";
import { history, useParams, useModel } from "umi";

interface AIModule {
  css: string[];
  js: string[];
  name: string;
  devices: string[];
}

interface AIAssetsResponse {
  runtimeEnv: any;
  aiModule: AIModule;
  aiAssetsModule: AIModule;
  token: string;
  signature: string;
  user: any;
}

export default function AIKnowledgeBase() {
  const { knowledgeBaseGuid } = useParams<{ knowledgeBaseGuid: string }>();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const loadedResources = useRef<Set<string>>(new Set());
  const { hostPath } = useModel('user')

  useEffect(() => {
    loadAIAssets();
  }, []);

  const loadAIAssets = async () => {
    if (loading) return;

    setLoading(true);
    setError(null);

    try {
      const response = await knowledgeBaseService.getAiAssets();
      const data: AIAssetsResponse = response.data;

      // Expose runtimeEnv on the window object
      if (data.runtimeEnv) {
        (window as any).__RUNTIME_ENV__ = {
          ...data.runtimeEnv,
          SDK_V2_TOKEN: data.token,
          SDK_V2_SIGNATURE: data.signature,
        };
        console.log('runtimeEnv 已设置:', data.runtimeEnv);
      }

      if (data.user) {
        (window as any).user = data.user
      }

      const url = new URL(hostPath)
      const origin = url.origin

      await loadJSFile(origin + '/minio/shimo-assets/static/sdk-iframe-assets/bootstrap.c13f32f2.js')

      if (data.aiModule && data.aiModule.js && data.aiModule.js.length > 0) {
        // Load the CSS files
        if (data.aiModule.css && data.aiModule.css.length > 0) {
          await loadCSSFiles(data.aiModule.css);
        }

        // Load the JS files
        await loadJSFiles(data.aiModule.js);

        // Inject environment variables

        console.log('AI 模块加载完成');
      } else {
        throw new Error('未找到 AI 模块的 JS 文件');
      }

      // Finally load the sdk-iframe-assets bundle
      if (data.aiAssetsModule && data.aiAssetsModule.js && data.aiAssetsModule.js.length > 0) {
        await loadJSFiles([origin + "/minio/shimo-assets/" + data.aiAssetsModule.js]);
      }
    } catch (err) {
      console.error('加载 AI 资源失败:', err);
      setError(err instanceof Error ? err.message : '加载 AI 资源失败');
    } finally {
      setLoading(false);
    }
  };

  const loadCSSFiles = async (cssUrls: string[]): Promise<void> => {
    const promises = cssUrls.map(url => loadCSSFile(url));
    await Promise.all(promises);
  };

  const loadCSSFile = (url: string): Promise<void> => {
    return new Promise((resolve, reject) => {
      if (loadedResources.current.has(url)) {
        resolve();
        return;
      }

      const link = document.createElement('link');
      link.rel = 'stylesheet';
      link.href = url;
      link.setAttribute('data-dynamic', 'ai');

      link.onload = () => {
        loadedResources.current.add(url);
        resolve();
      };

      link.onerror = () => {
        reject(new Error(`CSS 加载失败: ${url}`));
      };

      document.head.appendChild(link);
    });
  };

  const loadJSFiles = async (jsUrls: string[]): Promise<void> => {
    // Load JS files sequentially to preserve dependency order
    for (const url of jsUrls) {
      await loadJSFile(url);
    }
  };

  const loadJSFile = (url: string): Promise<void> => {
    return new Promise((resolve, reject) => {
      if (loadedResources.current.has(url)) {
        resolve();
        return;
      }

      const script = document.createElement('script');
      script.src = url;
      script.async = false;
      script.setAttribute('data-dynamic', 'ai');

      script.onload = () => {
        loadedResources.current.add(url);
        resolve();
      };

      script.onerror = () => {
        reject(new Error(`JS 加载失败: ${url}`));
      };

      document.head.appendChild(script);
    });
  };

  return (
    <div style={{ padding: '16px', width: '100%' }}>
      <h3>AI 知识库模块</h3>
    </div>
  );
}