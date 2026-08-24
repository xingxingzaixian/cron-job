import type { AppRouteRecordRaw } from '@/router/types';
import { BaseLayout } from '../types';

const routes: AppRouteRecordRaw[] = [
  {
    path: '/',
    name: 'Task',
    component: BaseLayout,
    meta: {
      title: 'route.task',
      icon: 'eos-icons:service-plan'
    },
    redirect: '/task/list',
    children: [
      {
        path: '/task/list',
        name: 'TaskList',
        component: () => import('@/views/task/list/index.vue'),
        meta: {
          title: 'route.task_list',
          icon: 'ic:baseline-format-list-bulleted'
        }
      },
      {
        path: '/task/log',
        name: 'TaskLog',
        component: () => import('@/views/task/log/index.vue'),
        meta: {
          title: 'route.task_log',
          icon: 'material-symbols-light:nest-clock-farsight-analog-outline-rounded'
        }
      },
      {
        path: '/task/edit',
        name: 'TaskEdit',
        component: () => import('@/views/task/edit/index.vue'),
        meta: {
          title: 'route.task_edit',
          icon: 'ic:baseline-edit-note',
          hideInMenu: true,
          // 编辑页为隐藏路由，高亮"任务列表"菜单项
          activeMenu: 'TaskList'
        }
      }
    ]
  }
];

export default routes;
