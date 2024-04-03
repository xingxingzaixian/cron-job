import type { AppRouteRecordRaw } from '@/router/types';
import { BaseLayout } from '../types';

const routes: AppRouteRecordRaw[] = [
  {
    path: '/login',
    name: 'Login',
    component: () => import('@/views/login'),
    meta: {
      title: 'route.login',
      ignoreAuth: true,
      hideInMenu: true
    }
  },
  {
    name: 'not-found',
    path: '/:pathMatch(.*)*',
    component: () => import('@/views/exception/404/index.vue'),
    meta: {
      title: 'route.not-found',
      ignoreAuth: true,
      hideInMenu: true
    }
  },
  {
    name: 'exception',
    path: '/exception',
    component: BaseLayout,
    meta: {
      title: 'route.exception',
      icon: 'ant-design:exception-outlined',
      order: 7,
      hideInMenu: true,
      ignoreAuth: true
    },
    children: [
      {
        name: 'exception_403',
        path: '/exception/403',
        component: () => import('@/views/exception/403/index.vue'),
        meta: {
          title: 'route.exception_403',
          icon: 'ic:baseline-block'
        }
      },
      {
        name: 'exception_404',
        path: '/exception/404',
        component: () => import('@/views/exception/404/index.vue'),
        meta: {
          title: 'route.exception_404',
          icon: 'ic:baseline-web-asset-off'
        }
      },
      {
        name: 'exception_500',
        path: '/exception/500',
        component: () => import('@/views/exception/500/index.vue'),
        meta: {
          title: 'route.exception_500',
          icon: 'ic:baseline-wifi-off'
        }
      }
    ]
  }
];

export default routes;
