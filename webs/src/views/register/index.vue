<template>
  <div class="login-container">
    <div class="absolute-lt flex-x-end p-3 w-full">
      <el-switch
        v-model="isDark"
        inline-prompt
        :active-icon="Moon"
        :inactive-icon="Sunny"
        @change="toggleTheme"
      />
      <lang-select class="ml-2 cursor-pointer" />
    </div>
    <el-card class="!border-none !bg-transparent !rounded-4% w-100 <sm:w-85">
      <div class="text-center relative">
        <h2>{{ defaultSettings.title }} - 注册</h2>
        <el-tag class="ml-2 absolute-rt">{{ version }}</el-tag>
      </div>

      <el-form :model="form" class="login-form">
        <el-form-item>
          <div class="flex-y-center w-full">
            <svg-icon icon-class="user" class="mx-2" />
            <el-input v-model="form.username" placeholder="用户名" size="large" class="h-[48px]" />
          </div>
        </el-form-item>
        <el-form-item>
          <div class="flex-y-center w-full">
            <svg-icon icon-class="lock" class="mx-2" />
            <el-input v-model="form.password" placeholder="密码" type="password" size="large" class="h-[48px] pr-2" show-password />
          </div>
        </el-form-item>
        <el-form-item>
          <div class="flex-y-center w-full">
            <svg-icon icon-class="user" class="mx-2" />
            <el-input v-model="form.nickname" placeholder="昵称(可选)" size="large" class="h-[48px]" />
          </div>
        </el-form-item>
        <el-form-item>
          <div class="flex-y-center w-full">
            <svg-icon icon-class="key" class="mx-2" />
            <el-input v-model="form.inviteCode" placeholder="邀请码(启用时必填)" size="large" class="h-[48px]" />
          </div>
        </el-form-item>

        <el-form-item>
          <div class="flex-y-center w-full">
            <svg-icon icon-class="captcha" class="mx-2" />
            <el-input
              v-model="form.captchaCode"
              auto-complete="off"
              size="large"
              class="flex-1"
              placeholder="验证码"
            />
            <el-image
              @click="getCaptcha"
              :src="captchaBase64"
              class="rounded-tr-md rounded-br-md cursor-pointer h-[48px]"
            />
          </div>
        </el-form-item>

        <div class="action-group">
          <el-button :loading="loading" type="primary" size="large" class="w-full action-btn" @click.prevent="handleRegister">注册账号</el-button>
          <el-button plain size="large" class="w-full action-btn" @click="router.push('/login')">返回登录</el-button>
        </div>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { registerApi, GetVersion, getCaptchaApi } from "@/api/auth";
import { LoginData } from "@/api/auth/types";
import { Sunny, Moon } from "@element-plus/icons-vue";
import { ElMessage } from "element-plus";
import router from "@/router";
import defaultSettings from "@/settings";
import { useSettingsStore } from "@/store";
import { ThemeEnum } from "@/enums/ThemeEnum";

const settingsStore = useSettingsStore();
const isDark = ref(settingsStore.theme === ThemeEnum.DARK);
const loading = ref(false);
const version = ref('');
const captchaBase64 = ref('');
const form = ref<LoginData>({
  username: '',
  password: '',
  nickname: '',
  inviteCode: '',
  captchaKey: '',
  captchaCode: ''
});

GetVersion().then((res) => {
  version.value = res.data;
});

const getCaptcha = () => {
  getCaptchaApi().then(({ data }) => {
    form.value.captchaKey = data.captchaKey;
    captchaBase64.value = data.captchaBase64;
  });
};

const handleRegister = async () => {
  if (!form.value.username || !form.value.password) {
    ElMessage.error('请输入用户名和密码');
    return;
  }
  if (form.value.password.length < 6) {
    ElMessage.error('密码至少6位');
    return;
  }
  if (!form.value.captchaCode) {
    ElMessage.error('请输入验证码');
    return;
  }
  try {
    loading.value = true;
    await registerApi(form.value);
    ElMessage.success('注册成功，请前往登录');
    await router.push('/login');
  } catch (error: any) {
    const message = typeof error === 'string' ? error : (error?.message || '注册失败');
    ElMessage.error(message);
    form.value.captchaCode = '';
    getCaptcha();
  } finally {
    loading.value = false;
  }
};

const toggleTheme = () => {
  const newTheme = settingsStore.theme === ThemeEnum.DARK ? ThemeEnum.LIGHT : ThemeEnum.DARK;
  settingsStore.changeTheme(newTheme);
};

onMounted(() => {
  getCaptcha();
});
</script>

<style lang="scss" scoped>
html.dark .login-container {
  background: url("@/assets/images/login-bg-dark.jpg") no-repeat center right;
}
.login-container {
  overflow-y: auto;
  background: url("@/assets/images/login-bg.jpg") no-repeat center right;
  @apply wh-full flex-center;
  .login-form {
    padding: 30px 10px;
  }
}
.el-form-item {
  background: var(--el-input-bg-color);
  border: 1px solid var(--el-border-color);
  border-radius: 5px;
}
.action-group {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.action-btn {
  margin-left: 0 !important;
}
:deep(.el-input) {
  .el-input__wrapper {
    padding: 0;
    background-color: transparent;
    box-shadow: none;
    &.is-focus,
    &:hover {
      box-shadow: none !important;
    }
  }
}
</style>
