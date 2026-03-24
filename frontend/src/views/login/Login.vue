<template>
  <div class="login-container">
    <div class="login-card">
      <div class="login-header">
        <h2 class="title">DMP Admin V2</h2>
        <p class="subtitle">智能行为分析平台</p>
      </div>
      
      <el-form
        ref="loginFormRef"
        :model="loginForm"
        :rules="loginRules"
        class="login-form"
        label-position="top"
        @keyup.enter="handleLogin"
      >
        <el-form-item prop="username" label="账号">
          <el-input
            v-model="loginForm.username"
            placeholder="请输入账号"
            prefix-icon="User"
          />
        </el-form-item>
        
        <el-form-item prop="password" label="密码">
          <el-input
            v-model="loginForm.password"
            type="password"
            placeholder="请输入密码"
            prefix-icon="Lock"
            show-password
          />
        </el-form-item>
        
        <div class="form-options">
          <el-checkbox v-model="rememberMe">记住我</el-checkbox>
          <el-button type="primary" link>忘记密码？</el-button>
        </div>
        
        <el-button
          :loading="loading"
          type="primary"
          class="login-button"
          @click="handleLogin"
        >
          登 录
        </el-button>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { User, Lock } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules } from 'element-plus'

const router = useRouter()
const userStore = useUserStore()

const loginFormRef = ref<FormInstance>()
const loading = ref(false)
const rememberMe = ref(false)

const loginForm = reactive({
  username: '',
  password: ''
})

const loginRules = reactive<FormRules>({
  username: [{ required: true, message: '请输入账号', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

const handleLogin = async () => {
  if (!loginFormRef.value) return
  
  await loginFormRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        // 模拟登录
        // const res = await login(loginForm)
        // userStore.setToken(res.access_token)
        
        // 暂时使用模拟数据
        if (loginForm.username === 'admin' && loginForm.password === '123456') {
          userStore.setToken('mock-token-' + Date.now())
          userStore.setUserInfo({
            id: '1',
            name: '管理员',
            avatar: ''
          })
          ElMessage.success('登录成功')
          router.push('/')
        } else {
          ElMessage.error('账号或密码错误（测试账号: admin/123456）')
        }
      } catch (error: any) {
        console.error(error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.login-container {
  height: 100vh;
  display: flex;
  justify-content: center;
  align-items: center;
  background-color: #f0f2f5;
  background-image: radial-gradient(#409eff 0.5px, transparent 0.5px);
  background-size: 20px 20px;
}

.login-card {
  width: 400px;
  padding: 40px;
  background: #fff;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.login-header {
  text-align: center;
  margin-bottom: 40px;
}

.title {
  margin: 0;
  font-size: 28px;
  color: #303133;
}

.subtitle {
  margin: 10px 0 0;
  font-size: 14px;
  color: #909399;
}

.login-form {
  margin-top: 20px;
}

.form-options {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.login-button {
  width: 100%;
  height: 40px;
  font-size: 16px;
}
</style>
