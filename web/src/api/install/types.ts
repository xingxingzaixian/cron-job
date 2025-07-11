/** 安装检查结果 */
export interface InstallCheckResult {
  installed: boolean;
  config_exists: boolean;
  db_configured: boolean;
}

/** 数据库测试参数 */
export interface DatabaseTestParams {
  engine: string;
  host?: string;
  port?: number;
  name: string;
  user?: string;
  password?: string;
  path?: string;
}

/** 数据库测试结果 */
export interface DatabaseTestResult {
  success: boolean;
  message: string;
}

/** 管理员用户配置 */
export interface AdminUserConfig {
  username: string;
  nickname: string;
  password: string;
  email?: string;
}

/** 安装参数 */
export interface InstallParams {
  database: DatabaseTestParams;
  auth_enabled: boolean;
  admin_user?: AdminUserConfig;
  http_port: number;
}

/** 安装结果 */
export interface InstallResult {
  success: boolean;
  message: string;
}