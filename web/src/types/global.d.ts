declare type Nullable<T> = T | null;
declare type Optional<T> = T | undefined;

declare type Recordable = Record<string, any>;

// 在vue中使用import.meta.env引入环境变量时增加的提示
interface ImportMetaEnv {
  readonly VITE_APP_TITLE: string;
  readonly VITE_APP_PORT: string;
  readonly VITE_APP_BASE_URL: string;
  readonly VITE_APP_NO_AUTH: string;
  readonly VITE_ICON_LOCAL_PREFIX: string;
  readonly VITE_APP_MENU_ICON: string;
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}

namespace I18n {
  type LangType = 'en-US' | 'zh-CN';

  type LangOption = {
    label: string;
    key: LangType;
  };
}

namespace Global {
  /** The global menu */
  interface Menu {
    /**
     * The menu key
     *
     * Equal to the route key
     */
    key: string;
    /** The menu label */
    label: string;
    /** */
    i18nKey: string;
    /** The menu order */
    order: number;
    /** The route path */
    routePath: string;
    /** The menu icon */
    icon?: () => VNode;
    /** The menu children */
    children?: Menu[];
  }

  type Breadcrumb = Omit<Menu, 'children'> & {
    options?: Breadcrumb[];
  };
}
