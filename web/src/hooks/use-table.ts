import { ref, reactive, watch, computed } from 'vue';
import type { Ref } from 'vue';
import { $t } from '@/locales';
import type { DataTableBaseColumn, PaginationProps } from 'naive-ui';
import type { TableColumn } from 'naive-ui/es/data-table/src/interface';
import { useLoading, useBoolean } from '@/hooks';
import useAppStore from '@/store/modules/app';

type BaseData = Record<string, any>;
type ApiFn = (args: any) => Promise<unknown>;

/** table config */
export type TableConfig<T, P> = {
  /** api function to get table data */
  apiFn: ApiFn;
  /** api params */
  apiParams?: P;
  /** transform api response to table data */
  transformer: (response: unknown) => any;
  /** pagination */
  pagination?: PaginationProps;
  /**
   * callback when pagination changed
   *
   * @param pagination
   */
  onPaginationChanged?: (pagination: PaginationProps) => void | Promise<void>;
  /**
   * whether to get data immediately
   *
   * @default true
   */
  immediate?: boolean;
  /** columns factory */
  columns: () => TableColumn<T>[];
};

export type FilteredColumn = {
  key: string;
  title: string;
  checked: boolean;
};

const useTable = <T, P extends BaseData & {}>(config: TableConfig<T, P>) => {
  const { loading, startLoading, endLoading } = useLoading();
  const { bool: empty, setBool: setEmpty } = useBoolean();

  const appStore = useAppStore();
  const data = ref<T>() as Ref<T>;

  const { apiFn, apiParams, transformer, onPaginationChanged, immediate = true } = config;
  const { columns, filteredColumns, reloadColumns } = useTableColumn(config.columns);

  const searchParams = reactive<NonNullable<P>>({ ...apiParams } as any);

  const pagination = reactive({
    page: apiParams?.pageNo || 1,
    pageSize: apiParams?.pageSize || 15,
    showSizePicker: true,
    pageSizes: [15, 30, 50, 100],
    // Fix Naive Pagination's outdated API
    onUpdatePage: async (page: number) => {
      pagination.page = page;

      await onPaginationChanged?.(pagination);
    },
    onUpdatePageSize: async (pageSize: number) => {
      pagination.pageSize = pageSize;
      pagination.page = 1;

      await onPaginationChanged?.(pagination);
    },
    ...config.pagination
  }) as PaginationProps;

  function updatePagination(update: Partial<PaginationProps>) {
    Object.assign(pagination, update);
  }

  async function getData() {
    startLoading();

    const formattedParams = formatSearchParams(searchParams);

    const response = await apiFn(formattedParams);

    const { data: tableData, pageNum, pageSize, total } = transformer(response as T);

    data.value = tableData;

    setEmpty(tableData.length === 0);
    updatePagination({ page: pageNum, pageSize, itemCount: total });
    endLoading();
  }

  function formatSearchParams(params: BaseData) {
    const formattedParams: BaseData = {};

    Object.entries(params).forEach(([key, value]) => {
      if (value !== null && value !== undefined) {
        formattedParams[key] = value;
      }
    });

    return formattedParams;
  }

  function updateSearchParams(params: Partial<P>) {
    Object.assign(searchParams, params);
  }

  function resetSearchParams() {
    Object.assign(searchParams, apiParams);
  }

  if (immediate) {
    getData();
  }

  watch(
    () => appStore.locale,
    () => {
      reloadColumns();
    }
  );

  return {
    loading,
    empty,
    data,
    columns,
    filteredColumns,
    reloadColumns,
    pagination,
    updatePagination,
    getData,
    searchParams,
    updateSearchParams,
    resetSearchParams
  };
};

export default useTable;

function useTableColumn<T>(factory: () => TableColumn<T>[]) {
  const SELECTION_KEY = '__selection__';

  const allColumns = ref(factory()) as Ref<TableColumn<T>[]>;

  const filteredColumns: Ref<FilteredColumn[]> = ref(getFilteredColumns(factory()));

  const columns = computed(() => getColumns());

  function reloadColumns() {
    allColumns.value = factory();
    filteredColumns.value = getFilteredColumns(factory());
  }

  function getFilteredColumns(aColumns: TableColumn<T>[]) {
    const cols: FilteredColumn[] = [];

    aColumns.forEach((column) => {
      if (column.type === undefined) {
        cols.push({
          key: column.key as string,
          title: column.title as string,
          checked: true
        });
      }

      if (column.type === 'selection') {
        cols.push({
          key: SELECTION_KEY,
          title: $t('common.check'),
          checked: true
        });
      }
    });

    return cols;
  }

  function getColumns() {
    const cols = filteredColumns.value
      .filter((column) => column.checked)
      .map((column) => {
        if (column.key === SELECTION_KEY) {
          return allColumns.value.find((col) => col.type === 'selection');
        }
        return allColumns.value.find((col) => (col as DataTableBaseColumn).key === column.key);
      });

    return cols as TableColumn<T>[];
  }

  return {
    columns,
    reloadColumns,
    filteredColumns
  };
}
