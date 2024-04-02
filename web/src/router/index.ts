import type { RouteRecordRaw, RouteLocationNormalized, NavigationGuardNext } from 'vue-router';
import { useTitle } from '@vueuse/core';
import { $t } from '@/locales';
import { Router, createRouter, createWebHashHistory, RouterHistory } from 'vue-router';
import { loadingBar } from '@/utils/message';
import useUserStore from '@/store/modules/user';
import { LOGIN_URL } from '@/enum';
import type { AppRouteRecordRaw } from './types';

class RouteView {
  // 路由对象
  private router: Router | unknown = undefined;
  private staticRoutes: AppRouteRecordRaw[] = [];

  constructor() {
    this.createBasicRoutes();
    this.router = this.createRouter();
    this.beforeRouteChange();
    this.afterRouteChange();
  }

  // 根据环境变量中的配置生成路由模式
  static createHistory = (): RouterHistory => {
    return createWebHashHistory();
  };

  // 动态获取 modules 目录下的所有 .ts 文件生成基础路由
  createBasicRoutes = () => {
    const moduleFiles: { [key: string]: any } = import.meta.glob('./modules/**/*.ts', { eager: true });
    const routeModuleList: RouteRecordRaw[] = [];
    Object.keys(moduleFiles).forEach((key) => {
      const mod: { [key: string]: any } = moduleFiles[key].default || {};
      const modList: RouteRecordRaw[] = Array.isArray(mod) ? [...mod] : [mod];
      routeModuleList.push(...modList);
    });

    this.staticRoutes = routeModuleList as AppRouteRecordRaw[];
  };

  // 创建路由对象
  createRouter(): Router {
    return createRouter({
      history: RouteView.createHistory(),
      routes: this.staticRoutes,
      strict: true,
      scrollBehavior: () => ({ left: 0, top: 0 })
    });
  }

  // 路由守卫
  private beforeRouteChange(): void {
    if (!this.router) {
      return;
    }

    const curRouter = this.router as Router;
    const userStore = useUserStore();
    curRouter.beforeEach((to, _, next) => {
      loadingBar.start();

      // 整个网站是否不需要登陆认证
      if (import.meta.env.VITE_APP_NO_AUTH.toLowerCase() === 'true') {
        next();
        return;
      }

      // 如果不需要登录认证，路由直接切换
      if (to.meta.ignoreAuth) {
        next();
        return;
      }

      // 如果没登录
      if (!userStore.isLoggedIn) {
        this.toLogin(to, next);
        return;
      }
      next();
    });
  }

  private afterRouteChange(): void {
    if (!this.router) {
      return;
    }

    const curRouter = this.router as Router;
    curRouter.afterEach((to) => {
      if (to.meta.title) {
        useTitle($t(to.meta.title as string));
      }

      loadingBar.finish();
    });
  }

  private toLogin(to: RouteLocationNormalized, next: NavigationGuardNext): void {
    next({
      path: LOGIN_URL,
      query: {
        redirect: to.fullPath
      }
    });
  }

  // 获取路由对象
  public getRouter(): Router {
    return this.router as Router;
  }

  public getRoutes() {
    return this.staticRoutes;
  }
}

export const router = new RouteView();

export default router.getRouter();
