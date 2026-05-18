<template>
  <div>
    <div class="gva-table-box">
      <div class="gva-btn-list">
        <el-button type="primary" icon="plus" @click="dialogVisible = true">新增地址池</el-button>
      </div>
      <el-table :data="tableData" row-key="ID">
        <el-table-column label="名称" prop="name" min-width="140" />
        <el-table-column label="CIDR" prop="cidr" min-width="140" />
        <el-table-column label="起始IP" prop="startIp" min-width="130" />
        <el-table-column label="结束IP" prop="endIp" min-width="130" />
        <el-table-column label="网关" prop="gateway" min-width="130" />
        <el-table-column label="DNS" prop="dns" min-width="150" />
        <el-table-column label="绑定网卡" prop="bindIface" min-width="120" />
        <el-table-column label="启用" width="90">
          <template #default="scope">
            <el-tag :type="scope.row.enabled ? 'success' : 'info'">
              {{ scope.row.enabled ? '启用' : '停用' }}
            </el-tag>
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

    <el-dialog v-model="dialogVisible" title="新增地址池" width="560px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="名称"><el-input v-model="form.name" /></el-form-item>
        <el-form-item label="CIDR"><el-input v-model="form.cidr" placeholder="192.168.50.0/24" /></el-form-item>
        <el-form-item label="起始IP"><el-input v-model="form.startIp" /></el-form-item>
        <el-form-item label="结束IP"><el-input v-model="form.endIp" /></el-form-item>
        <el-form-item label="网关"><el-input v-model="form.gateway" /></el-form-item>
        <el-form-item label="DNS"><el-input v-model="form.dns" /></el-form-item>
        <el-form-item label="绑定网卡"><el-input v-model="form.bindIface" /></el-form-item>
        <el-form-item label="启用"><el-switch v-model="form.enabled" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" @click="createPool">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup>
  import { ref } from 'vue'
  import { ElMessage } from 'element-plus'
  import { createPcfarmIPPool, getPcfarmIPPoolList } from '@/api/pcfarm'

  defineOptions({ name: 'PcfarmIPPool' })

  const page = ref(1)
  const pageSize = ref(10)
  const total = ref(0)
  const tableData = ref([])
  const dialogVisible = ref(false)
  const form = ref({
    name: '',
    cidr: '',
    startIp: '',
    endIp: '',
    gateway: '',
    dns: '',
    bindIface: '',
    enabled: true
  })

  const getTableData = async () => {
    const res = await getPcfarmIPPoolList({ page: page.value, pageSize: pageSize.value })
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

  const createPool = async () => {
    const res = await createPcfarmIPPool(form.value)
    if (res.code === 0) {
      ElMessage.success('创建成功')
      dialogVisible.value = false
      getTableData()
    }
  }

  getTableData()
</script>
