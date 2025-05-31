import type { AppRouteRecordRaw } from '@/router/types';
import { BaseLayout } from '../types';

const routes: AppRouteRecordRaw[] = [
  {
    path: '/user',
    name: 'User',
    component: BaseLayout,
    meta: {
      title: 'route.user',
      icon: 'mdi:account-group'
    },
    redirect: '/user/list',
    children: [
      {
        path: '/user/list',
        name: 'UserList',
        component: () => import('@/views/user/index.vue'),
        meta: {
          title: 'route.user_list',
          icon: 'mdi:account-multiple'
        }
      }
    ]
  }
];

export default routes;
