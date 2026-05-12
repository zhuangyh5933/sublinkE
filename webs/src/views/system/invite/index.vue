<template>
  <el-card class="page-card">
    <template #header>
      <div class="header-row">
        <span>邀请码管理</span>
        <el-button type="primary" @click="createInvite">创建邀请码</el-button>
      </div>
    </template>
    <el-table :data="invites" border>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="code" label="邀请码" />
      <el-table-column prop="description" label="说明" />
      <el-table-column prop="usedCount" label="使用次数" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="scope">
          <el-switch :model-value="scope.row.enabled" @change="(val) => toggleInvite(scope.row, Boolean(val))" />
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { addInvite, getInvites, updateInvite, type InviteItem } from '@/api/admin'

const invites = ref<InviteItem[]>([])

const fetchInvites = async () => {
  const res = await getInvites()
  invites.value = res.data || []
}

const createInvite = async () => {
  await addInvite({})
  ElMessage.success('创建成功')
  await fetchInvites()
}

const toggleInvite = async (row: InviteItem, enabled: boolean) => {
  await updateInvite({ id: row.id, enabled })
  ElMessage.success('更新成功')
  await fetchInvites()
}

onMounted(fetchInvites)
</script>

<style scoped>
.page-card { margin: 10px; }
.header-row { display: flex; justify-content: space-between; align-items: center; }
</style>
