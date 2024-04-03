import type { AppRouteRecordRaw } from '@/router/types';
import { BaseLayout } from '../types';

const routes: AppRouteRecordRaw[] = [
  {
    path: '/',
    name: 'Task',
    component: BaseLayout,
    meta: {
      title: 'route.task',
      ignoreAuth: true,
      icon: 'eos-icons:service-plan'
    },
    redirect: '/task/list',
    children: [
      {
        path: '/task/list',
        name: 'TaskList',
        component: () => import('@/views/task/list'),
        meta: {
          title: 'route.task',
          ignoreAuth: true
        }
      },
      {
        path: '/task/log',
        name: 'TaskLog',
        component: () => import('@/views/task/log'),
        meta: {
          title: 'route.task_log',
          ignoreAuth: true
        }
      }
    ]
  }
];

export default routes;
