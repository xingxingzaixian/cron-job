import apiHttp from '@/request';
import type { HttpResult } from '@/types/api';
import type { InstallCheckResult, DatabaseTestParams, DatabaseTestResult, InstallParams, InstallResult } from './types';

/** 检查安装状态 */
export function checkInstall() {
  return apiHttp.get<HttpResult<InstallCheckResult>>({
    url: '/api/install/check'
  });
}

/** 测试数据库连接 */
export function testDatabase(data: DatabaseTestParams) {
  return apiHttp.post<HttpResult<DatabaseTestResult>>({
    url: '/api/install/test-db',
    data
  });
}

/** 执行安装 */
export function install(data: InstallParams) {
  return apiHttp.post<HttpResult<InstallResult>>({
    url: '/api/install/install',
    data
  });
}