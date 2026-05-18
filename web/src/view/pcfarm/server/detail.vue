<template>
  <div>
    <div class="gva-table-box">
      <el-descriptions :column="2" border>
        <el-descriptions-item label="服务器ID">{{ serverId || '-' }}</el-descriptions-item>
        <el-descriptions-item label="Agent状态">待接入</el-descriptions-item>
        <el-descriptions-item label="BMC配置">密码仅支持重置，不回显明文</el-descriptions-item>
        <el-descriptions-item label="事件时间线">后端事件列表接口接入后展示</el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup>
  import { computed } from 'vue'
  import { useRoute } from 'vue-router'

  defineOptions({ name: 'PcfarmServerDetail' })

  const route = useRoute()
  const serverId = computed(() => route.query.id || route.params.id)
</script>
