import { ref, computed } from 'vue';
import { defineStore } from 'pinia';
import store from '../../index';
import router, { router as globalRouter } from '@/router';
import {
  getBreadcrumbsByRoute,
  getSelectedMenuKeyPathByKey,
  sortMenusByOrder,
  parseRoutesToMenus,
  getCacheRouteNames,
  updateLocaleOfGlobalMenus
} from './shared';

const routeStore = defineStore('route-store', () => {
  /** Global menus */
  const menus = ref<Global.Menu[]>([]);
  const cacheRoutes = ref<string[]>([]);

  /** Global breadcrumbs */
  const breadcrumbs = computed(() => getBreadcrumbsByRoute(router.currentRoute.value, menus.value));

  function getSelectedMenuKeyPath(selectedKey: string) {
    return getSelectedMenuKeyPathByKey(selectedKey, menus.value);
  }

  function updateGlobalMenusByLocale() {
    menus.value = updateLocaleOfGlobalMenus(menus.value);
  }

  async function initAuthRoute() {
    // 如果有动态路由，可以在这里请求服务接口返回动态路由
    const routes = [...globalRouter.getRoutes()];
    // 1. 解析路由为菜单
    const menuList = parseRoutesToMenus(routes);
    cacheRoutes.value = getCacheRouteNames(routes);

    // 2. 按照order进行排序
    menus.value = sortMenusByOrder(menuList);
  }

  return {
    menus,
    cacheRoutes,
    breadcrumbs,
    initAuthRoute,
    getSelectedMenuKeyPath,
    updateGlobalMenusByLocale
  };
});

const useRouteStore = () => {
  return routeStore(store);
};

export default useRouteStore;
