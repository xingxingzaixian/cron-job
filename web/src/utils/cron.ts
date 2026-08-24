/**
 * cron表达式（6字段：秒 分 时 日 月 周）前端校验与辅助工具
 * 规则与后端 internal/service/cron/task_manager/cron_spec.go 保持一致
 */

const FIELD_LABELS = ['秒', '分', '时', '日', '月', '周'];

const FIELD_RANGES: [number, number][] = [
  [0, 59],
  [0, 59],
  [0, 23],
  [1, 31],
  [1, 12],
  [0, 6]
];

const MONTH_NAMES: Record<string, number> = {
  january: 1, jan: 1,
  february: 2, feb: 2,
  march: 3, mar: 3,
  april: 4, apr: 4,
  may: 5,
  june: 6, jun: 6,
  july: 7, jul: 7,
  august: 8, aug: 8,
  september: 9, sep: 9,
  october: 10, oct: 10,
  november: 11, nov: 11,
  december: 12, dec: 12
};

const WEEK_NAMES: Record<string, number> = {
  sunday: 0, sun: 0,
  monday: 1, mon: 1,
  tuesday: 2, tue: 2,
  wednesday: 3, wed: 3,
  thursday: 4, thu: 4,
  friday: 5, fri: 5,
  saturday: 6, sat: 6
};

/** 常用频率快捷选项 */
export const CRON_PRESETS: { label: string; value: string }[] = [
  { label: '每秒', value: '* * * * * *' },
  { label: '每10秒', value: '*/10 * * * * *' },
  { label: '每分钟', value: '0 * * * * *' },
  { label: '每5分钟', value: '0 */5 * * * *' },
  { label: '每小时', value: '0 0 * * * *' },
  { label: '每天零点', value: '0 0 0 * * *' },
  { label: '每周一零点', value: '0 0 0 * * 1' },
  { label: '每月1号零点', value: '0 0 0 1 * *' }
];

function cronFieldValue(value: string, fieldIndex: number): number | null {
  if (/^\d+$/.test(value)) {
    return Number(value);
  }
  if (fieldIndex === 4) {
    const n = MONTH_NAMES[value.toLowerCase()];
    if (n !== undefined) return n;
  }
  if (fieldIndex === 5) {
    const n = WEEK_NAMES[value.toLowerCase()];
    if (n !== undefined) return n;
  }
  return null;
}

/** 校验cron表达式，返回错误信息；合法时返回空串 */
export function validateCronSpec(spec: string): string {
  const s = spec.trim();
  if (!s) {
    return 'cron表达式不能为空';
  }

  if (s.startsWith('@')) {
    if (/^@(yearly|annually|monthly|weekly|daily|midnight|hourly)$/i.test(s)) {
      return '';
    }
    // 时长格式交由后端精确解析（如 @every 5m / @every 1d），前端只要求有参数
    if (/^@every\s+\S+$/i.test(s)) {
      return '';
    }
    return '不支持的预定义格式（如 @daily、@every 5m）';
  }

  const fields = s.split(/\s+/);
  if (fields.length !== 6) {
    return `cron表达式需为6个字段（秒 分 时 日 月 周），当前为${fields.length}个字段`;
  }

  for (let i = 0; i < 6; i++) {
    if (i === 0 && fields[0] === '#') {
      continue;
    }
    const [min, max] = FIELD_RANGES[i];
    for (const item of fields[i].split(',')) {
      let base = item;
      const slashIdx = item.indexOf('/');
      if (slashIdx >= 0) {
        const step = Number(item.slice(slashIdx + 1));
        if (!Number.isInteger(step) || step <= 0) {
          return `第${i + 1}个字段步长无效: ${item}`;
        }
        base = item.slice(0, slashIdx);
      }
      if (base === '*' || base === '?') {
        continue;
      }

      const dashIdx = base.indexOf('-');
      if (dashIdx >= 0) {
        const from = cronFieldValue(base.slice(0, dashIdx), i);
        const to = cronFieldValue(base.slice(dashIdx + 1), i);
        if (from === null || to === null) {
          return `第${i + 1}个字段无效: ${item}`;
        }
        if (from < min || to > max || from > to) {
          return `第${i + 1}个字段取值范围错误: ${item}（合法范围 ${min}-${max}）`;
        }
        continue;
      }

      const val = cronFieldValue(base, i);
      if (val === null || val < min || val > max) {
        return `第${i + 1}个字段取值范围错误: ${item}（合法范围 ${min}-${max}）`;
      }
    }
  }
  return '';
}

/** 生成字段拆分提示，用于可视化展示；表达式不合法时返回 null */
export function cronSpecBreakdown(spec: string): string | null {
  if (validateCronSpec(spec) !== '') {
    return null;
  }
  const s = spec.trim();
  if (s.startsWith('@')) {
    return s;
  }
  const fields = s.split(/\s+/);
  if (fields.length !== 6) {
    return null;
  }
  return fields.map((field, i) => `${FIELD_LABELS[i]}:${field}`).join('  ');
}
