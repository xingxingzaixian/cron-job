<template>
  <div class="management-container">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="page-title">
            <n-icon size="24" class="title-icon">
              <svg viewBox="0 0 24 24">
                <path fill="currentColor" d="M12 22c1.1 0 2-.9 2-2h-4c0 1.1.9 2 2 2zm6-6v-5c0-3.07-1.63-5.64-4.5-6.32V4c0-.83-.67-1.5-1.5-1.5s-1.5.67-1.5 1.5v.68C7.64 5.36 6 7.92 6 11v5l-2 2v1h16v-1l-2-2z"/>
              </svg>
            </n-icon>
            通知配置管理
          </h1>
          <p class="page-description">管理任务通知配置、查看发送日志和测试通知功能</p>
        </div>
        <div class="header-actions">
          <n-button type="primary" size="large" class="add-btn" @click="handleAdd">
            <template #icon>
              <n-icon>
                <svg viewBox="0 0 24 24">
                  <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                </svg>
              </n-icon>
            </template>
            新增通知配置
          </n-button>
        </div>
      </div>
    </div>

    <!-- 数据表格 -->
    <n-card :bordered="false" class="glass-card table-card">
      <template #header-extra>
        <TableHeaderOperation
          v-model:columns="filteredColumns"
          :canDelete="false"
          :loading="loading"
          @refresh="getData"
        />
      </template>
      <NotificationSearch @search="search" @reset="reset" />
      <NDataTable
        v-model:checked-row-keys="checkedRowKeys"
        :columns="columns"
        :data="data"
        size="small"
        :scroll-x="962"
        :loading="loading"
        remote
        :pagination="pagination"
        :row-key="(item: any) => item.id"
        class="data-table"
        :max-height="tableHeight"
      />
    </n-card>

    <!-- 新增/编辑弹窗 -->
    <NotificationForm
      v-model:show="showFormModal"
      :notification-id="currentNotificationId"
      :task-id="currentTaskId"
      @success="handleFormSuccess"
    />

    <!-- 日志弹窗 -->
    <NotificationLog
      v-model:show="showLogModal"
      :task-id="currentTaskId"
    />
  </div>
</template>

<script lang="tsx" setup>
import { ref, computed, onMounted, onUnmounted } from 'vue';
import { $t } from '@/locales';
import { message } from '@/utils/message';
import { useTable } from '@/hooks';
import NotificationSearch from './modules/notification-search.vue';
import NotificationForm from './modules/notification-form.vue';
import NotificationLog from './modules/notification-log.vue';
import TableHeaderOperation from '@/components/custom/TableHeaderOperation.vue';
import { fetchNotificationList, fetchNotificationDelete } from '@/api/notification';
import type { SearchNotificationResponse, SearchNotificationParams, NotificationConfigOutput } from '@/api/notification/types';
import { NButton, NPopconfirm, NSwitch, NTag } from 'naive-ui';
import { NotificationType, NotificationTrigger, NotificationStatus } from '@/enum/notification';

defineOptions({ name: 'NotificationList' });

const checkedRowKeys = ref<string[]>([]);
const showFormModal = ref(false);
const showLogModal = ref(false);
const currentNotificationId = ref<number | null>(null);
const currentTaskId = ref<number>(0);

const { columns, filteredColumns, data, loading, pagination, updateSearchParams, resetSearchParams, getData } =
  useTable<SearchNotificationResponse, SearchNotificationParams>({
    apiFn: fetchNotificationList,
    apiParams: {
      pageNo: 1,
      pageSize: 15,
      task_id: 0,
    },
    transformer: (res: any) => {
      const { list = [], total = 0 } = res.data || {};

      return {
        data: list || [],
        pageNum: pagination.page,
        pageSize: pagination.pageSize,
        total
      };
    },
    onPaginationChanged(pg) {
      const { page, pageSize } = pg;

      updateSearchParams({
        pageNo: page,
        pageSize: pageSize
      });

      getData();
    },
    columns: () => [
      {
        key: 'id',
        title: 'ID',
        align: 'center',
        width: 70,
      },
      {
        key: 'task_id',
        title: '任务ID',
        align: 'center',
        width: 80,
      },
      {
        key: 'name',
        title: '通知名称',
        minWidth: 120,
      },
      {
        key: 'type',
        title: '通知类型',
        minWidth: 80,
        render: (row: any) => (
          <div>
            {row.type === NotificationType.Email ? (
              <NTag type="info" size="small">邮件</NTag>
            ) : (
              <NTag type="warning" size="small">Webhook</NTag>
            )}
          </div>
        )
      },
      {
        key: 'target',
        title: '目标地址',
        minWidth: 180,
        ellipsis: { tooltip: true },
      },
      {
        key: 'trigger',
        title: '触发条件',
        minWidth: 100,
        render: (row: any) => {
          const triggerMap: Record<string, { label: string; type: string }> = {
            [NotificationTrigger.Success]: { label: '成功时', type: 'success' },
            [NotificationTrigger.Failure]: { label: '失败时', type: 'error' },
            [NotificationTrigger.All]: { label: '所有情况', type: 'info' },
            [NotificationTrigger.Timeout]: { label: '超时时', type: 'warning' },
            [NotificationTrigger.Cancel]: { label: '取消时', type: 'default' },
          };
          const config = triggerMap[row.trigger] || { label: row.trigger, type: 'default' };
          return <NTag type={config.type as any} size="small">{config.label}</NTag>;
        }
      },
      {
        key: 'enabled',
        title: '状态',
        minWidth: 80,
        render: (row: any) => (
          <div>
            <NSwitch
              value={row.enabled}
              on-update:value={(val: boolean) => handleEnable(val, row)}
            />
          </div>
        )
      },
      {
        key: 'operate',
        title: $t('common.operate'),
        align: 'center',
        width: 280,
        render: (row: any) => (
          <div class="flex items-center justify-center gap-8px">
            <NButton type="primary" ghost size="small" class="mr-4px" onClick={() => handleEdit(row.id)}>
              {$t('common.edit')}
            </NButton>
            <NPopconfirm onPositiveClick={() => handleDelete(row.id)}>
              {{
                default: () => $t('common.confirmDelete'),
                trigger: () => (
                  <NButton type="error" ghost size="small" class="mr-4px">
                    {$t('common.delete')}
                  </NButton>
                )
              }}
            </NPopconfirm>
            <NButton ghost size="small" class="mr-4px" onClick={() => handleLog(row.task_id)}>
              {$t('common.log')}
            </NButton>
          </div>
        )
      }
    ]
  });

function handleAdd() {
  currentNotificationId.value = null;
  currentTaskId.value = 0;
  showFormModal.value = true;
}

function handleEdit(id: number) {
  currentNotificationId.value = id;
  showFormModal.value = true;
}

async function handleDelete(id: number) {
  const res = await fetchNotificationDelete(id);
  if (res.code === 200) {
    message.success($t('common.deleteSuccess'));
    getData();
  } else {
    message.error(res.message || $t('common.deleteSuccess'));
  }
}

async function handleEnable(val: boolean, row: NotificationConfigOutput) {
  // TODO: 实现启用/禁用功能
  message.warning('启用/禁用功能待实现');
}

function handleLog(taskId: number) {
  currentTaskId.value = taskId;
  showLogModal.value = true;
}

function handleFormSuccess() {
  getData();
}

const search = (model: Omit<SearchNotificationParams, 'pageNo' | 'pageSize'>) => {
  updateSearchParams({
    pageNo: 1,
    pageSize: 15,
    ...model
  });

  getData();
};

const reset = () => {
  resetSearchParams();
  getData();
};

// 动态计算表格高度
const windowHeight = ref(window.innerHeight);

const tableHeight = computed(() => {
  // 减去页面头部(120px) + 卡片头部(60px) + 搜索区域(80px) + 分页器(60px) + 边距(40px)
  return windowHeight.value - 360;
});

// 监听窗口大小变化
const handleResize = () => {
  windowHeight.value = window.innerHeight;
};

onMounted(() => {
  window.addEventListener('resize', handleResize);
});

onUnmounted(() => {
  window.removeEventListener('resize', handleResize);
});
</script>

<style scoped>
/* 通知配置列表页面特定样式 */
/* 大部分样式已提取到 management-page.css 公共样式文件中 */

/* 头部布局优化 */
.header-content {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

.header-left {
  flex-shrink: 0;
  min-width: 0;
}

.header-actions {
  flex-shrink: 0;
}

/* 响应式设计 */
@media (max-width: 1200px) {
  .header-content {
    flex-direction: column;
    align-items: stretch;
    gap: 16px;
  }
}

@media (max-width: 768px) {
  .header-content {
    gap: 12px;
  }

  .page-title {
    font-size: 20px;
  }

  .page-description {
    font-size: 13px;
  }
}
</style>
