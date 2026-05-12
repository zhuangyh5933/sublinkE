<template>
  <el-card class="page-card">
    <template #header>
      <div class="header-row">
        <span>用户管理</span>
        <div class="header-actions">
          <el-input v-model="keyword" placeholder="搜索用户名/昵称" clearable style="width: 220px" />
          <el-button type="primary" @click="fetchUsers">刷新</el-button>
        </div>
      </div>
    </template>
    <el-table :data="pagedUsers" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="username" label="用户名" />
      <el-table-column prop="nickname" label="昵称" />
      <el-table-column prop="role" label="角色" width="110" />
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-tag :type="scope.row.disabled ? 'danger' : 'success'">{{ scope.row.disabled ? '禁用' : '启用' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="subscriptionId" label="订阅ID" width="100" />
      <el-table-column prop="allowedRegions" label="允许地区(逗号分隔)" />
      <el-table-column label="操作" width="300">
        <template #default="scope">
          <el-button type="primary" link @click="openEdit(scope.row)">编辑</el-button>
          <el-button type="success" link @click="showPullLogs(scope.row)">拉取记录</el-button>
          <el-button type="warning" link @click="toggleDisabled(scope.row)">{{ scope.row.disabled ? '启用' : '禁用' }}</el-button>
          <el-button type="danger" link @click="resetToken(scope.row)">重置订阅Token</el-button>
          <el-button type="danger" link @click="removeUser(scope.row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>
    <div style="margin-top: 16px; display: flex; justify-content: flex-end;">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        layout="total, prev, pager, next"
        :total="filteredUsers.length"
      />
    </div>
  </el-card>

  <el-dialog v-model="pullLogDialogVisible" title="用户拉取记录" width="900px">
    <el-table :data="currentPullLogs" border>
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
  </el-dialog>

  <el-dialog v-model="dialogVisible" title="编辑用户" width="520px">
    <el-form :model="form" label-width="110px">
      <el-form-item label="角色">
        <el-select v-model="form.role">
          <el-option label="管理员" value="admin" />
          <el-option label="普通用户" value="user" />
        </el-select>
      </el-form-item>
      <el-form-item label="订阅ID">
        <el-select v-model="form.subscriptionId" placeholder="请选择订阅">
          <el-option label="不分配" :value="0" />
          <el-option v-for="item in subscriptions" :key="item.ID" :label="`${item.ID} - ${item.Name}`" :value="item.ID" />
        </el-select>
      </el-form-item>
      <el-form-item label="订阅Token">
        <el-input v-model="form.subscriptionToken" />
      </el-form-item>
      <el-form-item label="允许地区">
        <el-input v-model="form.allowedRegions" type="textarea" placeholder="例如: 中国,广东,香港" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="dialogVisible = false">取消</el-button>
      <el-button type="primary" @click="saveUser">保存</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { deleteAdminUser, getAdminUsers, resetSubscriptionToken, updateAdminUser, type AdminUserItem } from '@/api/admin'
import { getSubs } from '@/api/subcription/subs'

const users = ref<AdminUserItem[]>([])
const subscriptions = ref<{ ID: number; Name: string }[]>([])
const dialogVisible = ref(false)
const pullLogDialogVisible = ref(false)
const currentPullLogs = ref<any[]>([])
const form = ref<Partial<AdminUserItem>>({})
const keyword = ref('')
const currentPage = ref(1)
const pageSize = ref(10)

const filteredUsers = computed(() => {
  const q = keyword.value.trim().toLowerCase()
  if (!q) return users.value
  return users.value.filter((item) =>
    item.username.toLowerCase().includes(q) ||
    item.nickname.toLowerCase().includes(q)
  )
})

const pagedUsers = computed(() => {
  const start = (currentPage.value - 1) * pageSize.value
  return filteredUsers.value.slice(start, start + pageSize.value)
})

const fetchUsers = async () => {
  const res = await getAdminUsers()
  users.value = res.data || []
}

const fetchSubscriptions = async () => {
  const res = await getSubs()
  subscriptions.value = (res.data || []).map((item: any) => ({ ID: item.ID, Name: item.Name }))
}

const openEdit = (row: AdminUserItem) => {
  form.value = { ...row, subscriptionId: row.subscriptionId ?? 0 }
  dialogVisible.value = true
}

const saveUser = async () => {
  await updateAdminUser({
    id: Number(form.value.id),
    role: form.value.role,
    subscriptionId: Number(form.value.subscriptionId || 0),
    allowedRegions: form.value.allowedRegions || '',
    subscriptionToken: form.value.subscriptionToken || ''
  })
  ElMessage.success('保存成功')
  dialogVisible.value = false
  await fetchUsers()
}

const toggleDisabled = async (row: AdminUserItem) => {
  await updateAdminUser({
    id: row.id,
    disabled: !row.disabled,
    allowedRegions: row.allowedRegions || ''
  })
  ElMessage.success(row.disabled ? '已启用' : '已禁用')
  await fetchUsers()
}

const resetToken = async (row: AdminUserItem) => {
  const res = await resetSubscriptionToken(row.id)
  ElMessage.success(`重置成功: ${res.data.subscriptionToken}`)
  await fetchUsers()
}

const showPullLogs = (row: AdminUserItem) => {
  currentPullLogs.value = row.pullLogs || []
  pullLogDialogVisible.value = true
}

const removeUser = (row: AdminUserItem) => {
  ElMessageBox.confirm(`确认删除用户 ${row.username} 吗？`, '提示', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  }).then(async () => {
    await deleteAdminUser(row.id)
    ElMessage.success('删除成功')
    await fetchUsers()
  })
}

onMounted(async () => {
  await fetchUsers()
  await fetchSubscriptions()
})
</script>

<style scoped>
.page-card { margin: 10px; }
.header-row { display: flex; justify-content: space-between; align-items: center; }
.header-actions { display: flex; gap: 12px; align-items: center; }
</style>
