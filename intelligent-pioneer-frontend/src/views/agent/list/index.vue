<template>
  <div class="container">
    <a-card class="general-card" :bordered="false">
      <!-- 页面标题 -->
      <template #title>
        {{ $t('menu.agent.list') }}
      </template>

      <!-- 表格主体 -->
      <a-table
          row-key="id"
          :loading="loading"
          :pagination="pagination"
          :data="agentList"
          :bordered="false"
          :columns="columns"
          @page-change="handlePageChange"
      >
        <!-- 自定义渲染“状态”列 -->
        <template #status="{ record }">
          <a-space>
            <a-badge
                :status="record.status === 'online' ? 'success' : 'normal'"
            />
            <a-typography-text>
              {{ record.status === 'online' ? '在线' : '离线' }}
            </a-typography-text>
          </a-space>
        </template>

        <!-- 自定义渲染“最后心跳时间”列 -->
        <template #updated_at="{ record }">
          {{ formatTime(record.updated_at) }}
        </template>

        <!-- 自定义渲染“创建时间”列 -->
        <template #created_at="{ record }">
          {{ formatTime(record.created_at) }}
        </template>

        <!-- 未来可以添加“操作”列 -->
        <template #operations>
          <a-button type="text" size="small">查看详情</a-button>
        </template>
      </a-table>
    </a-card>
  </div>
</template>

<script lang="ts" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { useI18n } from 'vue-i18n';
import { getAgentList, type Agent } from '@/api/agent';
import { TableColumnData } from '@arco-design/web-vue';
import dayjs from 'dayjs';
import { Pagination } from '@/types/global';
import axios from "axios";

const { t } = useI18n();
const loading = ref(false);
const agentList = ref<Agent[]>([]);
let intervalId: number | undefined;

// --- 分页配置 ---
const basePagination: Pagination = {
  current: 1,
  pageSize: 20,
};
const pagination = ref({ ...basePagination });

const handlePageChange = (current: number) => {
  pagination.value.current = current;
  // 如果是前端分页，这里不需要重新请求
  // 如果是后端分页，需要调用 fetchData
};

// --- 表格列定义 ---
// 使用 computed 是为了让 t 函数在语言切换时能响应
const columns = computed<TableColumnData[]>(() => [
  { title: '主机名', dataIndex: 'hostname', sortable: { sortDirections: ['ascend', 'descend'] } },
  { title: 'IP 地址', dataIndex: 'ip_address' },
  { title: '操作系统', dataIndex: 'os', ellipsis: true, tooltip: true },
  { title: '状态', dataIndex: 'status', slotName: 'status' },
  { title: '最后心跳时间', dataIndex: 'updated_at', slotName: 'updated_at', sortable: { sortDirections: ['ascend', 'descend'] } },
  { title: '注册时间', dataIndex: 'created_at', slotName: 'created_at' },
  // { title: t('searchTable.columns.operations'), slotName: 'operations' }, // 未来启用操作列
]);

// --- 数据获取逻辑 ---
const fetchData = async () => {
  loading.value = true;
  try {
    const { data } = await getAgentList();
    agentList.value = data;
    // 如果是后端分页，需要在这里更新 pagination.total
    // pagination.value.total = res.total;
  } catch (err) {
    console.error('Failed to fetch agent list. Detailed error:', err);
    // 我们可以更进一步，看看错误的具体内容
    if (axios.isAxiosError(err)) {
      console.error('Axios error response:', err.response);
    }
  } finally {
    loading.value = false;
  }
};

// --- 时间格式化 ---
const formatTime = (timeStr: string) => {
  if (!timeStr || !dayjs(timeStr).isValid()) return '-';
  return dayjs(timeStr).format('YYYY-MM-DD HH:mm:ss');
};

// --- 组件生命周期钩子 ---
onMounted(() => {
  fetchData();
  // 设置一个定时器，每5秒自动刷新数据
  intervalId = window.setInterval(fetchData, 5000);
});

onUnmounted(() => {
  // 组件销毁时，必须清除定时器，防止内存泄漏
  if (intervalId) {
    window.clearInterval(intervalId);
  }
});
</script>

<style scoped lang="less">
.container {
  padding: 20px;
  height: 100%;
  :deep(.arco-card-body) {
    height: calc(100% - 60px); // 确保卡片内容区域占满
  }
}
</style>