<template>
  <div>
    <div class="gva-search-box">
      <el-form :inline="true" :model="searchInfo">
        <el-form-item label="关键词">
          <el-input v-model="searchInfo.keyword" clearable placeholder="资产编号 / 序列号 / MAC / IP" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchInfo.status" clearable placeholder="请选择">
            <el-option label="离线" value="offline" />
            <el-option label="PXE就绪" value="pxe_ready" />
            <el-option label="启动中" value="booting" />
            <el-option label="在线" value="online" />
            <el-option label="心跳丢失" value="heartbeat_lost" />
            <el-option label="远控失败" value="power_failed" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" icon="search" @click="onSearch">查询</el-button>
          <el-button icon="refresh" @click="onReset">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="openCreateDialog">新增服务器</el-button>
        <el-button icon="switch-button" :disabled="!selectedRows.length" @click="runBatchAction('on')">批量开机</el-button>
        <el-button icon="video-pause" :disabled="!selectedRows.length" @click="runBatchAction('off')">批量关机</el-button>
        <el-button icon="refresh-right" :disabled="!selectedRows.length" @click="runBatchAction('reboot')">批量重启</el-button>
        <el-button icon="connection" :disabled="!selectedRows.length" @click="runBatchAction('boot_pxe')">批量进PXE</el-button>
      </div>

      <el-table :data="tableData" row-key="ID" @selection-change="handleSelectionChange">
        <el-table-column type="selection" width="55" />
        <el-table-column label="资产编号" prop="assetCode" min-width="130" />
        <el-table-column label="序列号" prop="serialNumber" min-width="150" />
        <el-table-column label="PXE MAC" prop="pxeMac" min-width="160" />
        <el-table-column label="固定IP" prop="fixedIp" min-width="130" />
        <el-table-column label="远控协议" prop="powerProtocol" width="100" />
        <el-table-column label="启动策略" min-width="150">
          <template #default="scope">
            <el-select
              v-model="scope.row.bootPolicy"
              size="small"
              @change="(value) => changeBootPolicy(scope.row, value)"
            >
              <el-option label="本地盘" value="local_disk" />
              <el-option label="Ubuntu Live" value="ubuntu_live" />
              <el-option label="维护模式" value="maintenance" />
            </el-select>
          </template>
        </el-table-column>
        <el-table-column label="状态" prop="status" min-width="120" />
        <el-table-column label="最后心跳" min-width="180">
          <template #default="scope">
            <span>{{ scope.row.lastHeartbeatAt || '-' }}</span>
          </template>
        </el-table-column>
        <el-table-column label="操作" fixed="right" min-width="220">
          <template #default="scope">
            <el-button link type="primary" icon="switch-button" @click="runPowerAction(scope.row, 'on')">开机</el-button>
            <el-button link type="primary" icon="video-pause" @click="runPowerAction(scope.row, 'off')">关机</el-button>
            <el-button link type="primary" icon="refresh-right" @click="runPowerAction(scope.row, 'reboot')">重启</el-button>
            <el-button link type="primary" icon="connection" @click="runPowerAction(scope.row, 'boot_pxe')">PXE</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="gva-pagination">
        <el-pagination
          :current-page="page"
          :page-size="pageSize"
          :page-sizes="[10, 30, 50, 100]"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="handleCurrentChange"
          @size-change="handleSizeChange"
        />
      </div>
    </div>

    <el-dialog v-model="createDialogVisible" title="新增服务器" width="560px">
      <el-form :model="form" label-width="110px">
        <el-form-item label="资产编号"><el-input v-model="form.assetCode" /></el-form-item>
        <el-form-item label="序列号"><el-input v-model="form.serialNumber" /></el-form-item>
        <el-form-item label="PXE MAC"><el-input v-model="form.pxeMac" /></el-form-item>
        <el-form-item label="BMC地址"><el-input v-model="form.bmcAddress" /></el-form-item>
        <el-form-item label="BMC用户名"><el-input v-model="form.bmcUsername" /></el-form-item>
        <el-form-item label="BMC密码"><el-input v-model="form.bmcPassword" show-password /></el-form-item>
        <el-form-item label="远控协议">
          <el-select v-model="form.powerProtocol">
            <el-option label="IPMI" value="ipmi" />
            <el-option label="Redfish" value="redfish" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="createDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createServer">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage, ElMessageBox } from 'element-plus'
  import {
    createPcfarmServer,
    executePcfarmPowerAction,
    getPcfarmServerList,
    updatePcfarmBootPolicy
  } from '@/api/pcfarm'

  defineOptions({ name: 'PcfarmServer' })

  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const tableData = ref([])
  const selectedRows = ref([])
  const searchInfo = ref({ keyword: '', status: '' })
  const createDialogVisible = ref(false)
  const form = ref({
    assetCode: '',
    serialNumber: '',
    pxeMac: '',
    bmcAddress: '',
    bmcUsername: '',
    bmcPassword: '',
    powerProtocol: 'ipmi'
  })

  const getTableData = async () => {
    const res = await getPcfarmServerList({
      page: page.value,
      pageSize: pageSize.value,
      ...searchInfo.value
    })
    if (res.code === 0) {
      tableData.value = res.data.list
      total.value = res.data.total
      page.value = res.data.page
      pageSize.value = res.data.pageSize
    }
  }

  const handleSizeChange = (val) => {
    pageSize.value = val
    getTableData()
  }

  const handleCurrentChange = (val) => {
    page.value = val
    getTableData()
  }

  const handleSelectionChange = (rows) => {
    selectedRows.value = rows
  }

  const onSearch = () => {
    page.value = 1
    getTableData()
  }

  const onReset = () => {
    searchInfo.value = { keyword: '', status: '' }
    onSearch()
  }

  const openCreateDialog = () => {
    form.value = {
      assetCode: '',
      serialNumber: '',
      pxeMac: '',
      bmcAddress: '',
      bmcUsername: '',
      bmcPassword: '',
      powerProtocol: 'ipmi'
    }
    createDialogVisible.value = true
  }

  const createServer = async () => {
    const res = await createPcfarmServer(form.value)
    if (res.code === 0) {
      ElMessage.success('创建成功')
      createDialogVisible.value = false
      getTableData()
    }
  }

  const changeBootPolicy = async (row, bootPolicy) => {
    const res = await updatePcfarmBootPolicy({ id: row.ID, bootPolicy })
    if (res.code === 0) {
      ElMessage.success('启动策略已更新')
      getTableData()
    }
  }

  const runPowerAction = async (row, action) => {
    if (['off', 'reboot', 'boot_pxe'].includes(action)) {
      await ElMessageBox.confirm('该操作会影响服务器电源或启动状态，是否继续？', '危险操作确认', {
        confirmButtonText: '继续',
        cancelButtonText: '取消',
        type: 'warning'
      })
    }
    const res = await executePcfarmPowerAction({ id: row.ID, action })
    if (res.code === 0) {
      ElMessage.success('操作已提交')
    }
  }

  const runBatchAction = async (action) => {
    await ElMessageBox.confirm('该操作会影响选中的服务器电源状态，是否继续？', '危险操作确认', {
      confirmButtonText: '继续',
      cancelButtonText: '取消',
      type: 'warning'
    })
    for (const row of selectedRows.value) {
      await executePcfarmPowerAction({ id: row.ID, action })
    }
    ElMessage.success('批量操作已提交')
  }

  getTableData()
</script>
