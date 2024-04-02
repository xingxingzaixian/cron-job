import type { AppRouteRecordRaw } from '@/router/types';

const routes: AppRouteRecordRaw[] = [
  {
    path: '/task',
    name: 'Task',
    component: () => import('@/views/task'),
    meta: {
      title: 'route.task',
      ignoreAuth: true,
      icon: 'eos-icons:service-plan',
      hideInMenu: true
    }
  }
];

export default routes;
