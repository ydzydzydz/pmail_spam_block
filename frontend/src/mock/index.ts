import Mock from 'mockjs'
import type { ApiResponse, TestModel, TestModelData, Setting } from '@/types'
import {
  pluginSettingsResource,
  testModelResource,
  getSettingResource,
  updateSettingResource,
} from '@/api/resource'

// 模拟测试模型响应
const Random = Mock.Random

// 模拟超时时长
Mock.setup({
  timeout: '500-1000',
})

// 模拟获取测试模型响应
export const postTestModelMock = Mock.mock(
  `${pluginSettingsResource}/${testModelResource}`,
  'post',
  (options: { body: TestModel }) => {
    const code = Random.integer(0, -1) // 0 表示成功，其他值表示失败
    const message = code === 0 ? '测试模型成功' : '测试模型失败'
    const setting: ApiResponse<TestModelData> = {
      code,
      message,
      data: {
        predictions: [[Random.float(0, 1), Random.float(0, 1), Random.float(0, 1)]],
      },
    }
    console.log('postTestModel mock response', options.body)
    return setting
  },
)

// 模拟获取设置信息
export const getSettingMock = Mock.mock(
  `${pluginSettingsResource}/${getSettingResource}`,
  'get',
  () => {
    const setting: ApiResponse<Setting> = {
      code: 0,
      message: '获取设置信息成功',
      data: {
        api_url: Random.url(),
        timeout: Random.integer(1000, 5000),
        threshold: Random.float(0, 1),
      },
    }
    // console.log('getSetting mock response', setting)
    return setting
  },
)

// 模拟更新设置信息
export const updateSettingMock = Mock.mock(
  `${pluginSettingsResource}/${updateSettingResource}`,
  'post',
  (options: { body: Setting }) => {
    // console.log('updateSetting options', options)
    const code = Random.integer(0, -1) // 0 表示成功，其他值表示失败
    const message = code === 0 ? '更新设置信息成功' : '更新设置信息失败'
    const setting: ApiResponse<Setting> = {
      code,
      message,
    }
    console.log('updateSetting mock response', options.body)
    return setting
  },
)

export default { postTestModelMock, getSettingMock, updateSettingMock }
