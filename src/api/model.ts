import { post } from '@/api'
import type { TestModel, TestModelData, ApiResponse } from '@/types'
import { testModelResource } from '@/api/resource'

// 测试模型
export function postTestModel(request: TestModel): Promise<ApiResponse<TestModelData>> {
  return post(testModelResource, request)
}
