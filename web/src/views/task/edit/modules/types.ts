export interface KvRow {
  key: string;
  value: string;
  enabled: boolean;
  desc: string;
}

export const newKvRow = (): KvRow => ({
  key: '',
  value: '',
  enabled: true,
  desc: ''
});
