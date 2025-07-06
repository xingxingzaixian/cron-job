import type { AppRouteRecordRaw } from '@/router/types';

const routes: AppRouteRecordRaw[] = [
  {
    path: '/install',
    name: 'Install',
    component: () => import('@/views/install/index.vue'),
    meta: {
      title: 'route.install',
      ignoreAuth: true,
      hideInMenu: true
    }
  }
];

export default routes;