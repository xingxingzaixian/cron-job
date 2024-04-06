export interface HttpPro {
  url: string;
  method: string;
}

export interface HttpParams {
  headers: Record<string, any>;
  query: Record<string, any>;
  data: Record<string, any>;
}
