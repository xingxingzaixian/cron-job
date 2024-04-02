import type { Component } from 'vue';
import type { RouteMeta, RouteRecordRaw } from 'vue-router';

export interface AppRouteMeta extends RouteMeta {
  title: string;
  ignoreAuth?: boolean;
  permission?: string;
  // 是否在菜单中隐藏
  hideInMenu?: boolean;
  keepAlive?: boolean;
  icon?: string;
}

// 类型和接口在扩展联合类型方面存在不足，因此可以用类型来扩展，而不是接口
export type AppRouteRecordRaw = RouteRecordRaw & {
  name: string;
  meta?: AppRouteMeta;
  component?: Component | string;
  components?: Component;
  children?: AppRouteRecordRaw[];
  props?: Recordable;
  fullPath?: string;
};

// 菜单类型
export interface MenuType {
  name: string;
  title: string;
  icon: string;
  children?: MenuType[];
}

export const BaseLayout = () => import('@/layouts/base-layout/index.vue');
export const BlankLayout = () => import('@/layouts/blank-layout/index.vue');
