import type { RouteLocationNormalizedLoaded } from 'vue-router';
import type { AppRouteRecordRaw } from '@/router/types';
import { $t } from '@/locales';
import { useSvgIcon } from '@/hooks';

/**
 * Get global menu by route
 *
 * @param route
 */
function getGlobalMenuByBaseRoute(route: RouteLocationNormalizedLoaded) {
  const { name, path } = route;
  const { title } = route.meta ?? {};

  const label = $t(title as string);
  const { SvgIconVNode } = useSvgIcon();
  const menu: Global.Menu = {
    key: name as string,
    label,
    i18nKey: title as string,
    routePath: path,
    order: Number(route.meta?.order) ?? 0,
    icon: SvgIconVNode({
      icon: (route.meta?.icon as string) ?? import.meta.env.VITE_APP_MENU_ICON,
      fontSize: 20
    })
  };

  return menu;
}

/**
 * Transform menu to breadcrumb
 *
 * @param menu
 */
function transformMenuToBreadcrumb(menu: Global.Menu) {
  const { children, ...rest } = menu;

  const breadcrumb: Global.Breadcrumb = {
    ...rest
  };

  if (children?.length) {
    breadcrumb.options = children.map(transformMenuToBreadcrumb);
  }

  return breadcrumb;
}

/**
 * Get breadcrumbs by route
 *
 * @param route
 * @param menus
 */
export function getBreadcrumbsByRoute(route: RouteLocationNormalizedLoaded, menus: Global.Menu[]): Global.Breadcrumb[] {
  const key = route.name as string;
  const activeKey = route.meta?.activeMenu;

  const menuKey = activeKey || key;

  for (const menu of menus) {
    if (menu.key === menuKey) {
      const breadcrumbMenu = menuKey !== activeKey ? menu : getGlobalMenuByBaseRoute(route);

      return [transformMenuToBreadcrumb(breadcrumbMenu)];
    }

    if (menu.children?.length) {
      const result = getBreadcrumbsByRoute(route, menu.children);
      if (result.length > 0) {
        return [transformMenuToBreadcrumb(menu), ...result];
      }
    }
  }

  return [];
}

/**
 * Get selected menu key path
 *
 * @param selectedKey
 * @param menus
 */
export function getSelectedMenuKeyPathByKey(selectedKey: string, menus: Global.Menu[]) {
  const keyPath: string[] = [];

  menus.some((menu) => {
    const path = findMenuPath(selectedKey, menu);

    const find = Boolean(path?.length);

    if (find) {
      keyPath.push(...path!);
    }

    return find;
  });

  return keyPath;
}

/**
 * Find menu path
 *
 * @param targetKey Target menu key
 * @param menu Menu
 */
function findMenuPath(targetKey: string, menu: Global.Menu): string[] | null {
  const path: string[] = [];

  function dfs(item: Global.Menu): boolean {
    path.push(item.key);

    if (item.key === targetKey) {
      return true;
    }

    if (item.children) {
      for (const child of item.children) {
        if (dfs(child)) {
          return true;
        }
      }
    }

    path.pop();

    return false;
  }

  if (dfs(menu)) {
    return path;
  }

  return null;
}

/**
 * sort route by order
 *
 * @param route route
 */
function sortMenuByOrder(route: Global.Menu) {
  if (route.children?.length) {
    route.children.sort((next, prev) => (next.order || 0) - (prev.order || 0));
    route.children.forEach(sortMenuByOrder);
  }

  return route;
}

/**
 * sort routes by order
 *
 * @param routes routes
 */
export function sortMenusByOrder(routes: Global.Menu[]) {
  if (!routes) {
    return [];
  }

  routes.sort((next, prev) => (next.order || 0) - (prev.order || 0));
  routes.forEach(sortMenuByOrder);

  return routes;
}

export function parseRoutesToMenus(routes: AppRouteRecordRaw[]) {
  const menus: Global.Menu[] = [];
  for (const route of routes) {
    if (route.meta?.hideInMenu) {
      continue;
    }

    const menu = getGlobalMenuByBaseRoute(route as unknown as RouteLocationNormalizedLoaded);
    if (route.children?.length) {
      const children = parseRoutesToMenus(route.children);

      if (children?.length > 0) {
        menu.children = children;
      }
    }

    menus.push(menu);
  }
  return menus;
}

export function getCacheRouteNames(routes: AppRouteRecordRaw[]) {
  const cacheNames: string[] = [];

  routes.forEach((route) => {
    // only get last two level route, which has component
    route.children?.forEach((child: AppRouteRecordRaw) => {
      if (child.component && child.meta?.keepAlive) {
        cacheNames.push(child.name);
      }
    });
  });

  return cacheNames;
}

export function updateLocaleOfGlobalMenus(menus: Global.Menu[]) {
  const result: Global.Menu[] = [];

  menus.forEach((menu) => {
    const { i18nKey, label, children } = menu;

    const newLabel = i18nKey ? $t(i18nKey) : label;

    const newMenu: Global.Menu = {
      ...menu,
      label: newLabel
    };

    if (children?.length) {
      newMenu.children = updateLocaleOfGlobalMenus(children);
    }

    result.push(newMenu);
  });

  return result;
}
