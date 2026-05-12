// 插件信息类型定义
export interface PluginInfo {
  name: string;
  version: string;
  description: string;
  enabled: boolean;
  filePath: string;
  config?: Record<string, any>;
}

// 插件配置请求参数
export interface PluginConfigRequest {
  name: string;
  config: Record<string, any>;
}
