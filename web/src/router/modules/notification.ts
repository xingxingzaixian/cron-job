import type { AppRouteRecordRaw } from '@/router/types';
import { BaseLayout } from '../types';

const routes: AppRouteRecordRaw[] = [
  {
    path: '/notification',
    name: 'Notification',
    component: BaseLayout,
    meta: {
      title: 'route.notification',
      icon: 'mdi:bell-outline'
    },
    redirect: '/notification/list',
    children: [
      {
        path: '/notification/list',
        name: 'NotificationList',
        component: () => import('@/views/notification/index.vue'),
        meta: {
          title: 'route.notification_list',
          icon: 'ic:baseline-format-list-bulleted'
        }
      }
    ]
  }
];

export default routes;
