<template>
  <div class="user-management">
    <!-- 页面头部 -->
    <div class="page-header">
      <div class="header-content">
        <div class="header-left">
          <h1 class="page-title">
            <n-icon size="24" class="title-icon">
              <svg viewBox="0 0 24 24">
                <path fill="currentColor" d="M12 12c2.21 0 4-1.79 4-4s-1.79-4-4-4-4 1.79-4 4 1.79 4 4 4zm0 2c-2.67 0-8 1.34-8 4v2h16v-2c0-2.66-5.33-4-8-4z"/>
              </svg>
            </n-icon>
            用户管理
          </h1>
          <p class="page-description">管理系统用户信息和账户状态</p>
        </div>
        <div class="header-actions">
          <n-button type="primary" size="large" @click="handleAdd" class="add-btn">
            <template #icon>
              <n-icon>
                <svg viewBox="0 0 24 24">
                  <path fill="currentColor" d="M19 13h-6v6h-2v-6H5v-2h6V5h2v6h6v2z"/>
                </svg>
              </n-icon>
            </template>
            新增用户
          </n-button>
        </div>
      </div>
    </div>

    <!-- 搜索和筛选区域 -->
    <n-card :bordered="false" class="search-card">
      <n-form inline :label-width="80" class="search-form">
        <n-form-item label="用户名">
          <n-input
            v-model:value="searchForm.username"
            placeholder="请输入用户名"
            clearable
            @keydown.enter="handleSearch"
            class="search-input"
          />
        </n-form-item>
        <n-form-item label="昵称">
          <n-input
            v-model:value="searchForm.name"
            placeholder="请输入昵称"
            clearable
            @keydown.enter="handleSearch"
            class="search-input"
          />
        </n-form-item>
        <n-form-item label="邮箱">
          <n-input
            v-model:value="searchForm.email"
            placeholder="请输入邮箱"
            clearable
            @keydown.enter="handleSearch"
            class="search-input"
          />
        </n-form-item>
        <n-form-item>
          <n-space>
            <n-button type="primary" @click="handleSearch" :loading="loading">
              <template #icon>
                <n-icon>
                  <svg viewBox="0 0 24 24">
                    <path fill="currentColor" d="M15.5 14h-.79l-.28-.27C15.41 12.59 16 11.11 16 9.5 16 5.91 13.09 3 9.5 3S3 5.91 3 9.5 5.91 16 9.5 16c1.61 0 3.09-.59 4.23-1.57l.27.28v.79l5 4.99L20.49 19l-4.99-5zm-6 0C7.01 14 5 11.99 5 9.5S7.01 5 9.5 5 14 7.01 14 9.5 11.99 14 9.5 14z"/>
                  </svg>
                </n-icon>
              </template>
              搜索
            </n-button>
            <n-button @click="handleReset">
              <template #icon>
                <n-icon>
                  <svg viewBox="0 0 24 24">
                    <path fill="currentColor" d="M17.65 6.35C16.2 4.9 14.21 4 12 4c-4.42 0-7.99 3.58-7.99 8s3.57 8 7.99 8c3.73 0 6.84-2.55 7.73-6h-2.08c-.82 2.33-3.04 4-5.65 4-3.31 0-6-2.69-6-6s2.69-6 6-6c1.66 0 3.14.69 4.22 1.78L13 11h7V4l-2.35 2.35z"/>
                  </svg>
                </n-icon>
              </template>
              重置
            </n-button>
          </n-space>
        </n-form-item>
      </n-form>
    </n-card>

    <!-- 数据表格 -->
    <n-card :bordered="false" class="table-card">
      <n-data-table
        :columns="columns"
        :data="userList"
        :loading="loading"
        :pagination="pagination"
        :row-key="(row: any) => row.id"
        @update:page="handlePageChange"
        @update:page-size="handlePageSizeChange"
        class="user-table"
      />
    </n-card>

    <!-- 用户编辑弹窗 -->
    <n-modal v-model:show="showModal" preset="dialog" title="用户信息" style="width: 600px">
      <n-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="用户名" path="username">
          <n-input v-model:value="formData.username" placeholder="请输入用户名" :disabled="!!formData.id" />
        </n-form-item>
        <n-form-item label="昵称" path="nickname">
          <n-input v-model:value="formData.nickname" placeholder="请输入昵称" />
        </n-form-item>
        <n-form-item label="邮箱" path="email">
          <n-input v-model:value="formData.email" placeholder="请输入邮箱" />
        </n-form-item>
        <n-form-item v-if="!formData.id" label="密码" path="password">
          <n-input v-model:value="formData.password" type="password" placeholder="请输入密码" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-space>
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="submitLoading" @click="handleSubmit">
            {{ formData.id ? '更新' : '创建' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 修改密码弹窗 -->
    <n-modal v-model:show="showPasswordModal" preset="dialog" title="修改密码" style="width: 500px">
      <n-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-placement="left"
        label-width="auto"
        require-mark-placement="right-hanging"
      >
        <n-form-item label="原密码" path="oldPass">
          <n-input v-model:value="passwordForm.oldPass" type="password" placeholder="请输入原密码" />
        </n-form-item>
        <n-form-item label="新密码" path="newPass">
          <n-input v-model:value="passwordForm.newPass" type="password" placeholder="请输入新密码" />
        </n-form-item>
        <n-form-item label="确认密码" path="confirmPassword">
          <n-input v-model:value="passwordForm.confirmPassword" type="password" placeholder="请再次输入新密码" />
        </n-form-item>
      </n-form>
      <template #action>
        <n-space>
          <n-button @click="showPasswordModal = false">取消</n-button>
          <n-button type="primary" :loading="passwordSubmitLoading" @click="handlePasswordSubmit">
            确认修改
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, h } from 'vue';
import type { DataTableColumns, FormInst } from 'naive-ui';
import { NButton, NSpace, useDialog, useMessage } from 'naive-ui';
import type {
  UserEditInput,
  SearchUserParams,
  DelUserParams,
  UpdatePasswordInput
} from '@/api/user/model';
import {
  fetchUserList,
  fetchUserEdit,
  fetchUserDelete,
  fetchUpdatePassword
} from '@/api/user';

const message = useMessage();
const dialog = useDialog();

// 响应式数据
const loading = ref(false);
const userList = ref<UserEditInput[]>([]);
const showModal = ref(false);
const showPasswordModal = ref(false);
const submitLoading = ref(false);
const passwordSubmitLoading = ref(false);

// 表单引用
const formRef = ref<FormInst | null>(null);
const passwordFormRef = ref<FormInst | null>(null);

// 搜索表单
const searchForm = reactive<SearchUserParams>({
  pageNo: 1,
  pageSize: 10,
  name: '',
  username: '',
  email: ''
});

// 分页
const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  showQuickJumper: true,
  prefix: (info: any) => `共 ${info.itemCount} 条`
});

// 用户表单
const formData = reactive<UserEditInput>({
  id: undefined,
  username: '',
  nickname: '',
  password: '',
  email: ''
});

// 密码表单
const passwordForm = reactive<UpdatePasswordInput>({
  oldPass: '',
  newPass: '',
  confirmPassword: ''
});

// 表格列定义
const columns: DataTableColumns<UserEditInput> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    align: 'center'
  },
  {
    title: '用户名',
    key: 'username',
    width: 120
  },
  {
    title: '昵称',
    key: 'nickname',
    width: 120
  },
  {
    title: '邮箱',
    key: 'email',
    width: 200
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    render: (row) => {
      return row.created_at ? new Date(row.created_at).toLocaleString() : '-';
    }
  },
  {
    title: '操作',
    key: 'actions',
    width: 200,
    align: 'center',
    render: (row) => {
      return h(NSpace, { justify: 'center' }, {
        default: () => [
          h(NButton, {
            size: 'small',
            type: 'primary',
            secondary: true,
            onClick: () => handleEdit(row)
          }, { default: () => '编辑' }),
          h(NButton, {
            size: 'small',
            type: 'warning',
            secondary: true,
            onClick: () => handleChangePassword(row)
          }, { default: () => '改密' }),
          h(NButton, {
            size: 'small',
            type: 'error',
            secondary: true,
            onClick: () => handleDelete(row)
          }, { default: () => '删除' })
        ]
      });
    }
  }
];

// 表单验证规则
const formRules = {
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 20, message: '用户名长度在 3 到 20 个字符', trigger: 'blur' }
  ],
  nickname: [
    { required: true, message: '请输入昵称', trigger: 'blur' },
    { min: 2, max: 20, message: '昵称长度在 2 到 20 个字符', trigger: 'blur' }
  ],
  email: [
    { type: 'email', message: '请输入正确的邮箱地址', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ]
};

// 密码表单验证规则
const passwordRules = {
  oldPass: [
    { required: true, message: '请输入原密码', trigger: 'blur' }
  ],
  newPass: [
    { required: true, message: '请输入新密码', trigger: 'blur' },
    { min: 6, max: 20, message: '密码长度在 6 到 20 个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认新密码', trigger: 'blur' },
    {
      validator: (rule: any, value: string) => {
        return value === passwordForm.newPass;
      },
      message: '两次输入的密码不一致',
      trigger: 'blur'
    }
  ]
};

// 方法定义
const loadUserList = async () => {
  loading.value = true;
  try {
    const res = await fetchUserList(searchForm);
    if (res.code === 200) {
      userList.value = res.data.list;
      pagination.total = res.data.total;
    } else {
      message.error(res.message || '获取用户列表失败');
    }
  } catch (error) {
    message.error('获取用户列表失败');
  } finally {
    loading.value = false;
  }
};

const handleSearch = () => {
  searchForm.pageNo = 1;
  pagination.page = 1;
  loadUserList();
};

const handleReset = () => {
  Object.assign(searchForm, {
    pageNo: 1,
    pageSize: 10,
    name: '',
    username: '',
    email: ''
  });
  pagination.page = 1;
  loadUserList();
};

const handlePageChange = (page: number) => {
  searchForm.pageNo = page;
  pagination.page = page;
  loadUserList();
};

const handlePageSizeChange = (pageSize: number) => {
  searchForm.pageSize = pageSize;
  pagination.pageSize = pageSize;
  searchForm.pageNo = 1;
  pagination.page = 1;
  loadUserList();
};

const resetForm = () => {
  Object.assign(formData, {
    id: undefined,
    username: '',
    nickname: '',
    password: '',
    email: ''
  });
};

const handleAdd = () => {
  resetForm();
  showModal.value = true;
};

const handleEdit = (row: UserEditInput) => {
  Object.assign(formData, {
    ...row,
    password: '' // 编辑时不显示密码
  });
  showModal.value = true;
};

const handleSubmit = async () => {
  if (!formRef.value) return;

  try {
    await formRef.value.validate();
    submitLoading.value = true;

    const res = await fetchUserEdit(formData);
    if (res.code === 200) {
      message.success(formData.id ? '更新成功' : '创建成功');
      showModal.value = false;
      loadUserList();
    } else {
      message.error(res.message || '操作失败');
    }
  } catch (error) {
    // 验证失败或网络错误
  } finally {
    submitLoading.value = false;
  }
};

const handleDelete = (row: UserEditInput) => {
  dialog.warning({
    title: '确认删除',
    content: `确定要删除用户 "${row.nickname}" 吗？此操作不可恢复。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        const params: DelUserParams = { id: row.id! };
        const res = await fetchUserDelete(params);
        if (res.code === 200) {
          message.success('删除成功');
          loadUserList();
        } else {
          message.error(res.message || '删除失败');
        }
      } catch (error) {
        message.error('删除失败');
      }
    }
  });
};

const handleChangePassword = (row: UserEditInput) => {
  Object.assign(passwordForm, {
    oldPass: '',
    newPass: '',
    confirmPassword: ''
  });
  showPasswordModal.value = true;
};

const handlePasswordSubmit = async () => {
  if (!passwordFormRef.value) return;

  try {
    await passwordFormRef.value.validate();
    passwordSubmitLoading.value = true;

    const res = await fetchUpdatePassword(passwordForm);
    if (res.code === 200) {
      message.success('密码修改成功');
      showPasswordModal.value = false;
    } else {
      message.error(res.message || '密码修改失败');
    }
  } catch (error) {
    // 验证失败或网络错误
  } finally {
    passwordSubmitLoading.value = false;
  }
};

// 生命周期
onMounted(() => {
  loadUserList();
});
</script>

<style scoped>
/* 登录页面风格的渐变背景色彩 */
.user-management {
  min-height: 100vh;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 24px;
}

/* 页面头部 */
.page-header {
  margin-bottom: 24px;
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  border-radius: 16px;
  padding: 24px 32px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}

.header-left {
  flex: 1;
}

.page-title {
  display: flex;
  align-items: center;
  margin: 0 0 8px 0;
  font-size: 28px;
  font-weight: 600;
  color: #2c3e50;
  background: linear-gradient(135deg, #667eea, #764ba2);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.title-icon {
  margin-right: 12px;
  color: #667eea;
}

.page-description {
  margin: 0;
  color: #64748b;
  font-size: 14px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.add-btn {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
  border-radius: 12px;
  padding: 0 24px;
  height: 44px;
  font-weight: 500;
  box-shadow: 0 4px 16px rgba(102, 126, 234, 0.3);
  transition: all 0.3s ease;
}

.add-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 6px 20px rgba(102, 126, 234, 0.4);
}

/* 搜索卡片 */
.search-card {
  margin-bottom: 24px;
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.search-form {
  padding: 8px;
}

.search-input {
  width: 200px;
}

/* 表格卡片 */
.table-card {
  border-radius: 16px;
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(10px);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.1);
}

.user-table {
  border-radius: 12px;
}

/* 响应式设计 */
@media (max-width: 768px) {
  .user-management {
    padding: 16px;
  }

  .header-content {
    flex-direction: column;
    gap: 16px;
    padding: 20px;
  }

  .page-title {
    font-size: 24px;
  }

  .search-input {
    width: 100%;
  }
}

/* 动画效果 */
.search-card,
.table-card {
  animation: fadeInUp 0.6s ease-out;
}

.page-header {
  animation: fadeInDown 0.6s ease-out;
}

@keyframes fadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* 表格行悬停效果 */
:deep(.n-data-table-tbody .n-data-table-tr:hover) {
  background-color: rgba(102, 126, 234, 0.05);
}

/* 按钮样式优化 */
:deep(.n-button--primary-type) {
  background: linear-gradient(135deg, #667eea, #764ba2);
  border: none;
}

:deep(.n-button--primary-type:hover) {
  background: linear-gradient(135deg, #5a6fd8, #6a4190);
}
</style>

