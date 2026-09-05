// 定义 API 响应体结构
export interface ApiResponse<T> {
  code: number
  message: string
  data?: T
}

// 定义插件设置结构
export interface Setting {
  api_url: string
  timeout: number
  threshold: number
}

// 定义测试模型请求结构
export interface TestModel {
  setting: Setting
  content: string
}

// 定义测试模型响应结构
export interface TestModelData {
  predictions: [[number, number, number]]
}
