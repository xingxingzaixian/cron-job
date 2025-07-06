import type { RouteRecordRaw, RouteLocationNormalized, NavigationGuardNext, Router, RouterHistory } from 'vue-router';
import { useTitle } from '@vueuse/core';
import { $t } from '@/locales';
import { createRouter, createWebHashHistory } from 'vue-router';
import { loadingBar } from '@/utils/message';
import useUserStore from '@/store/modules/user';
import type { AppRouteRecordRaw } from './types';

class RouteView {
  // 路由对象
  private router: Router | unknown = undefined;
  private staticRoutes: AppRouteRecordRaw[] = [];

  public constructor() {
    this.createBasicRoutes();
    this.router = this.createRouter();
    this.beforeRouteChange();
    this.afterRouteChange();
  }

  // 获取路由对象
  public getRouter(): Router {
    return this.router as Router;
  }

  public getRoutes() {
    return this.staticRoutes;
  }

  // 根据环境变量中的配置生成路由模式
  // eslint-disable-next-line @typescript-eslint/member-ordering
  private static createHistory = (): RouterHistory => {
    return createWebHashHistory();
  };

  // 动态获取 modules 目录下的所有 .ts 文件生成基础路由
  private createBasicRoutes = () => {
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
  private createRouter(): Router {
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
    curRouter.beforeEach(async (to, _, next) => {
      loadingBar.start();

      // 检查安装状态（除了安装页面本身）
      if (to.name !== 'Install') {
        try {
          const { checkInstall } = await import('@/api/install');
          const { data } = await checkInstall();
          if (!data.installed) {
            // 未安装，跳转到安装页面
            next({ name: 'Install' });
            return;
          }
        } catch (error) {
          console.error('检查安装状态失败:', error);
          // 如果检查失败，假设需要安装
          if (to.name !== 'Install') {
            next({ name: 'Install' });
            return;
          }
        }
      }

      // 如果不需要登录认证，路由直接切换
      if (to.meta.ignoreAuth) {
        next();
        return;
      }

      // 获取是否启用认证，检查是否启用认证
      try {
        const checkAuth = await userStore.getCheckAuth();
        if (checkAuth && !checkAuth.authEnabled) {
          // 如果未启用认证，直接放行
          next();
          return;
        }
      } catch (error) {
        console.error('Failed to get check auth:', error);
      }

      // 如果没登录
      if (!userStore.isLoggedIn) {
        this.toLogin(to, next);
        return;
      } else {
        userStore.getUserInfo();
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
      name: 'Login',
      query: {
        redirect: to.fullPath
      }
    });
  }
}

export const router = new RouteView();

export default router.getRouter();
