export interface ServiceListOutput {
  list: ServiceItemOutput[];
  total: number;
}

export interface ServiceItemOutput {
  id: number;
  rule: string;
  serviceDesc: string;
  serviceName: string;
  urlRewrite: string;
}

export interface ServiceSearchParam {
  info: string;
  pageNo: number;
  pageSize: number;
}

export interface ServiceAddHTTPInput {
  needHttps: 0 | 1;
  rule: string;
  serviceDesc: string;
  serviceName: string;
  urlRewrite: string;
}

export interface ServiceUpdateHTTPInput extends ServiceAddHTTPInput {
  id: number | null;
}

export interface HttpRule {
  id: number;
  needHttps: 0 | 1;
  rule: string;
  serviceId: number;
  urlRewrite: string;
}

export interface ServiceInfo {
  id: number;
  serviceName: string;
  serviceDesc: string;
  isDelete: number;
}

export interface ServiceDetail {
  httpRule: HttpRule;
  info: ServiceInfo;
}

export interface ServiceRequestLogItem {
  id: number;
  originUrl: string;
  proxyUrl: string;
  serviceName: string;
}

export interface ServiceRequestLogList {
  totle: number;
  list: ServiceRequestLogItem[];
}

export interface ServiceStatisticsRes {
  key: string;
  value: number;
}
