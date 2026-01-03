<template>
  <div class="pmail-spam-block-settings">
    <el-card class="box-card">
      <template #header>
        <div class="card-header">
          <span>插件设置</span>
        </div>
      </template>

      <el-form :model="formData" label-width="180px" label-position="left" v-loading="loading">
        <el-form-item label="模型API接口地址：" prop="api_url">
          <el-input v-model="formData.api_url" placeholder="请输入模型API接口地址" />
        </el-form-item>

        <el-form-item label="接口超时时间（毫秒）：" prop="timeout">
          <el-input-number v-model="formData.timeout" :min="1" :max="50000" :step="1">
            <template #suffix>
              <span>ms</span>
            </template>
          </el-input-number>
        </el-form-item>

        <el-form-item label="模型识别阈值（0-1）：" prop="threshold">
          <el-slider v-model="formData.threshold" :min="0" :max="1" :step="0.01" />
        </el-form-item>
      </el-form>

      <el-form-item>
        <div class="form-item-buttons">
          <el-button
            type="info"
            :disabled="loading"
            style="margin: 0 auto"
            @click="dialogVisible = true"
          >
            <LinkIcon class="icon" />
            测试接口
          </el-button>

          <el-dialog v-model="dialogVisible" title="请输入测试邮件内容" width="50%">
            <el-input
              :autosize="{ minRows: 3, maxRows: 6 }"
              type="textarea"
              v-model="inputValue"
              placeholder="请输入测试邮件内容"
            ></el-input>
            <template #footer>
              <div class="dialog-footer">
                <el-button @click="dialogVisible = false">关闭</el-button>
                <el-button type="primary" @click="testModel"> 测试 </el-button>
              </div>
            </template>
          </el-dialog>

          <el-dialog v-model="testResultVisible" title="模型接口测试结果" width="50%">
            <el-table :data="testResultData" border stripe style="width: 100%">
              <el-table-column prop="type" label="类型" />
              <el-table-column prop="score" label="分数" sortable />
            </el-table>
            <template #footer>
              <div class="dialog-footer">
                <el-button @click="testResultVisible = false">关闭</el-button>
              </div>
            </template>
          </el-dialog>

          <el-button
            type="primary"
            :disabled="loading"
            style="margin: 0 auto"
            @click="confirmSubmit"
          >
            <SaveIcon class="icon" />
            保存设置
          </el-button>
        </div>
      </el-form-item>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import SaveIcon from '@/assets/icons/save-line.svg'
import LinkIcon from '@/assets/icons/link.svg'
import { ref, onMounted } from 'vue'
import resizeIframeHeight from '@/utils/resize'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { Setting, TestModel } from '@/types'
import { postTestModel } from '@/api/model'
import { getSetting, updateSetting } from '@/api/setting'

const loading = ref(false)
const formData = ref<Setting>({
  api_url: '',
  timeout: 5000,
  threshold: 0.5,
})

// 获取插件配置
const fetchSetting = () => {
  loading.value = true
  getSetting()
    .then((response) => {
      if (response.data) {
        formData.value = response.data
      }
    })
    .catch((error) => {
      ElMessage.error('获取设置信息失败')
      console.error(error)
    })
    .finally(() => {
      loading.value = false
    })
}

// 保存插件配置
const saveSetting = () => {
  loading.value = true
  updateSetting(formData.value || {})
    .then(() => {
      ElMessage.success('设置已保存')
    })
    .catch((error) => {
      ElMessage.error('保存设置失败')
      console.error(error)
    })
    .finally(() => {
      loading.value = false
    })
}

// 确认保存设置
const confirmSubmit = () => {
  ElMessageBox.confirm('确认保存设置吗？', '保存设置', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning',
  }).then(() => {
    saveSetting()
  })
}

const dialogVisible = ref(false)
const inputValue = ref('测试邮件')
const testResultVisible = ref(false)
const testResultData = ref([
  {
    type: '正常邮件',
    score: 0,
  },
  {
    type: '广告邮件',
    score: 0,
  },
  {
    type: '垃圾邮件',
    score: 0,
  },
])

// 测试模型接口
const testModel = () => {
  const api_url = formData.value.api_url
  if (!api_url) {
    ElMessage.error('请输入模型API接口地址')
    return
  }

  const content = inputValue.value.trim()
  if (!content) {
    ElMessage.error('请输入测试邮件内容')
    return
  }

  const testModel: TestModel = {
    setting: formData.value,
    content,
  }
  loading.value = true
  dialogVisible.value = false
  postTestModel(testModel)
    .then((response) => {
      if (response.code === 0) {
        ElMessage.success(response.message || '模型接口测试成功')
        const predictions = response.data?.predictions || []
        const prediction = predictions[0] || [0, 0, 0]
        testResultData.value.forEach((item, index) => {
          item.score = prediction[index] || 0
        })
        testResultVisible.value = true
      } else {
        ElMessage.error(response.message || '模型接口测试失败')
      }
    })
    .catch((error) => {
      ElMessage.error('模型接口测试失败')
      console.error(error)
    })
    .finally(() => {
      loading.value = false
    })
}

onMounted(() => {
  resizeIframeHeight()
  fetchSetting()
})
</script>

<style scoped>
.pmail-spam-block-settings {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.box-card {
  margin-top: 20px;
}

.card-header {
  font-size: 18px;
  font-weight: bold;
}

.icon {
  margin-right: 5px;
  width: 20px;
  height: 20px;
  vertical-align: middle;
  color: white;
}

.form-item-buttons {
  margin-top: 20px;
  margin: 0 auto;
  display: flex;
  justify-content: center;
  align-items: center;
  gap: 20px;
}
</style>
