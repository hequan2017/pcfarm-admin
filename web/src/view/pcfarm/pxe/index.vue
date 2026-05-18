<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="refresh" @click="refreshConfig">刷新配置</el-button>
        <el-button icon="circle-check" @click="loadStatus">校验状态</el-button>
      </div>
      <el-descriptions :column="1" border>
        <el-descriptions-item label="Ubuntu Live镜像路径">
          /srv/pcfarm/ubuntu-live
        </el-descriptions-item>
        <el-descriptions-item label="TFTP根目录">
          /srv/tftp
        </el-descriptions-item>
        <el-descriptions-item label="HTTP/NFS地址">
          http://pcfarm.local/ubuntu-live
        </el-descriptions-item>
        <el-descriptions-item label="dnsmasq服务状态">
          <el-tag>{{ status }}</el-tag>
        </el-descriptions-item>
      </el-descriptions>
    </div>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { getPcfarmPXEStatus, refreshPcfarmPXE } from '@/api/pcfarm'

  defineOptions({ name: 'PcfarmPXE' })

  const status = ref('unknown')

  const loadStatus = async () => {
    const res = await getPcfarmPXEStatus()
    if (res.code === 0) {
      status.value = res.data.status
    }
  }

  const refreshConfig = async () => {
    const res = await refreshPcfarmPXE()
    if (res.code === 0) {
      ElMessage.success('刷新请求已提交')
      loadStatus()
    }
  }

  loadStatus()
</script>
