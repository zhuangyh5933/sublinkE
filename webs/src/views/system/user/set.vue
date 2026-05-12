<script setup lang='ts'>
import { ref,onMounted } from 'vue'
import {useUserStore} from "@/store"
import {getUserPullLogsApi, updateUserPassword} from "@/api/user"
import { useI18n } from 'vue-i18n'
import type { PullLog, UserInfo } from '@/api/user/types'

const { t } = useI18n()
const userinfo = ref<UserInfo>()
const userStore = useUserStore()
const username:Ref<string> = ref('')
const password:Ref<string> = ref('')
const pullLogs = ref<PullLog[]>([])

onMounted( async() => {
  userinfo.value = await userStore.getUserInfo()
  const res = await getUserPullLogsApi()
  pullLogs.value = (res.data || []).map((item: PullLog) => ({
    ...item,
    IP: item.IP || (item as any).ip,
    Region: item.Region || (item as any).region,
    Addr: item.Addr || (item as any).addr,
    Client: item.Client || (item as any).client,
    Status: item.Status || (item as any).status,
    Count: item.Count || (item as any).count,
    Date: item.Date || (item as any).date,
  }))
})

const copyUrl = async (url?: string) => {
  if (!url) {
    ElMessage.warning('暂无订阅链接')
    return
  }
  try {
    await navigator.clipboard.writeText(url)
    ElMessage.success('复制成功')
  } catch (error) {
    ElMessage.error('复制失败')
  }
}

/** 重置密码 */
function resetPassword() {
  if (!username.value || !password.value) {
    ElMessage.error(t('userset.message.xx1'))
    return
  }
  if ((password.value.length < 6)) {
    ElMessage.error(t('userset.message.xx2'))
    return
  }
  ElMessageBox.confirm(
    t('userset.message.xx3'),
    t('userset.message.title'),
    {
      confirmButtonText: 'OK',
      cancelButtonText: 'Cancel',
      type: 'warning',
    }
  )
    .then(() => {
      updateUserPassword({
          username:username.value.trim(),
          password:password.value.trim()
        }
      ).then(() => {
        ElMessage.success(t('userset.message.xx4') + password.value);
        window.location.reload();
      });
    })
}
</script>

<template>
  <div>
    <el-card style="margin: 10px;text-align: center;">
      <el-row :gutter="20" justify="center">
        <el-col :span="18">
          <h2>{{$t('userset.title')}}</h2>
        </el-col>
        <el-col :span="18" v-if="userinfo">
          <el-badge :value="userinfo.username" class="item">
            <el-image :src="userinfo.avatar" />
          </el-badge>
        </el-col>

        <el-col :span="18">
          <el-input
            v-model="username"
            :placeholder="$t('userset.newUsername')"
          />
        </el-col>
        <el-col :span="18">
          <el-input
            type="password"
            v-model="password"
            :placeholder="$t('userset.newPassword')"
          />
        </el-col>
        <el-col :span="18">
          <el-button type="primary" @click="resetPassword">修改</el-button>
        </el-col>
      </el-row>
    </el-card>

    <el-card style="margin: 10px;" v-if="userinfo">
      <template #header>
        <div class="subscription-header">
          <span>我的订阅</span>
          <el-tag v-if="userinfo.subscriptionName" type="success">{{ userinfo.subscriptionName }}</el-tag>
          <el-tag v-else type="warning">未分配订阅</el-tag>
        </div>
      </template>

      <el-descriptions :column="1" border>
        <el-descriptions-item label="订阅Token">{{ userinfo.subscriptionToken || '-' }}</el-descriptions-item>
        <el-descriptions-item label="允许地区">{{ userinfo.allowedRegions || '不限' }}</el-descriptions-item>
        <el-descriptions-item label="自动识别链接">
          <div class="sub-line">
            <span>{{ userinfo.subscriptionUrl || '未生成' }}</span>
            <el-button size="small" type="primary" @click="copyUrl(userinfo.subscriptionUrl)">复制</el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="Clash">
          <div class="sub-line">
            <span>{{ userinfo.clashUrl || '未生成' }}</span>
            <el-button size="small" type="primary" @click="copyUrl(userinfo.clashUrl)">复制</el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="Surge">
          <div class="sub-line">
            <span>{{ userinfo.surgeUrl || '未生成' }}</span>
            <el-button size="small" type="primary" @click="copyUrl(userinfo.surgeUrl)">复制</el-button>
          </div>
        </el-descriptions-item>
        <el-descriptions-item label="V2Ray">
          <div class="sub-line">
            <span>{{ userinfo.v2rayUrl || '未生成' }}</span>
            <el-button size="small" type="primary" @click="copyUrl(userinfo.v2rayUrl)">复制</el-button>
          </div>
        </el-descriptions-item>
      </el-descriptions>
    </el-card>

    <el-card style="margin: 10px;">
      <template #header>
        <div class="subscription-header">
          <span>订阅拉取记录</span>
          <el-tag type="info">{{ pullLogs.length }} 条</el-tag>
        </div>
      </template>
      <el-table :data="pullLogs" border>
        <el-table-column prop="IP" label="IP" />
        <el-table-column prop="Region" label="地区" />
        <el-table-column prop="Addr" label="位置详情" />
        <el-table-column prop="Client" label="客户端" width="100" />
        <el-table-column label="状态" width="140">
          <template #default="scope">
            <el-tag :type="scope.row.Status === 'blocked_region' ? 'danger' : 'success'">
              {{ scope.row.Status === 'blocked_region' ? '地区拦截失败' : '成功' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="Count" label="拉取次数" width="100" />
        <el-table-column prop="Date" label="最近拉取时间" width="180" />
      </el-table>
    </el-card>
  </div>
</template>

<style scoped>
.el-col {
  margin-bottom: 10px
}
.subscription-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}
.sub-line {
  display: flex;
  gap: 12px;
  align-items: center;
  justify-content: space-between;
  word-break: break-all;
}
</style>
