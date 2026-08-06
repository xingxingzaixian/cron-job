-- ============================================================
-- 任务日志表按月分区脚本（MySQL 8.0+）
--
-- 背景：任务日志表会随执行次数无限增长，按月分区后可以
--   - 快速删除过期分区（DROP PARTITION，秒级清理，不锁全表）
--   - 单分区数据量可控，查询只扫描相关分区
--
-- 注意：MySQL 要求分区键必须包含在表的主键/唯一键中，
--       存量表需要先修改主键为 (id, start_time)，请先在测试环境验证。
-- ============================================================

-- 方案一：新表（推荐，数据量可控时）
CREATE TABLE sched_task_log (
    id          BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    created_at  DATETIME(3) NULL,
    updated_at  DATETIME(3) NULL,
    deleted_at  DATETIME(3) NULL,
    task_id     INT NOT NULL DEFAULT 0,
    task_name   VARCHAR(32) NOT NULL DEFAULT '',
    protocol    TINYINT NOT NULL DEFAULT 1,
    retry_times TINYINT NOT NULL DEFAULT 0,
    status      TINYINT NOT NULL DEFAULT 0,
    result      MEDIUMTEXT,
    start_time  DATETIME NULL,
    end_time    DATETIME NULL,
    duration    BIGINT NOT NULL DEFAULT 0,
    PRIMARY KEY (id, start_time),
    KEY idx_task_log_task_id (task_id),
    KEY idx_task_log_task_name (task_name),
    KEY idx_task_log_status (status),
    KEY idx_task_log_start_time (start_time)
) ENGINE=InnoDB
  PARTITION BY RANGE COLUMNS(start_time) (
    PARTITION p202401 VALUES LESS THAN ('2024-02-01'),
    PARTITION p202402 VALUES LESS THAN ('2024-03-01'),
    PARTITION p202403 VALUES LESS THAN ('2024-04-01'),
    PARTITION p202404 VALUES LESS THAN ('2024-05-01'),
    PARTITION p202405 VALUES LESS THAN ('2024-06-01'),
    PARTITION p202406 VALUES LESS THAN ('2024-07-01'),
    PARTITION p202407 VALUES LESS THAN ('2024-08-01'),
    PARTITION p202408 VALUES LESS THAN ('2024-09-01'),
    PARTITION p202409 VALUES LESS THAN ('2024-10-01'),
    PARTITION p202410 VALUES LESS THAN ('2024-11-01'),
    PARTITION p202411 VALUES LESS THAN ('2024-12-01'),
    PARTITION p202412 VALUES LESS THAN ('2025-01-01'),
    PARTITION p202501 VALUES LESS THAN ('2025-02-01'),
    PARTITION p202502 VALUES LESS THAN ('2025-03-01'),
    PARTITION p202503 VALUES LESS THAN ('2025-04-01'),
    PARTITION p202504 VALUES LESS THAN ('2025-05-01'),
    PARTITION p202505 VALUES LESS THAN ('2025-06-01'),
    PARTITION p202506 VALUES LESS THAN ('2025-07-01'),
    PARTITION p202507 VALUES LESS THAN ('2025-08-01'),
    PARTITION p202508 VALUES LESS THAN ('2025-09-01'),
    PARTITION p202509 VALUES LESS THAN ('2025-10-01'),
    PARTITION p202510 VALUES LESS THAN ('2025-11-01'),
    PARTITION p202511 VALUES LESS THAN ('2025-12-01'),
    PARTITION p202512 VALUES LESS THAN ('2026-01-01'),
    PARTITION p_future VALUES LESS THAN (MAXVALUE)
);

-- 方案二：存量表改造（先备份！需要修改主键）
-- ALTER TABLE sched_task_log DROP PRIMARY KEY, ADD PRIMARY KEY (id, start_time);
-- ALTER TABLE sched_task_log
--     PARTITION BY RANGE COLUMNS(start_time) (
--         PARTITION p202401 VALUES LESS THAN ('2024-02-01'),
--         ... 同上列出历史月份 ...
--         PARTITION p_future VALUES LESS THAN (MAXVALUE)
--     );

-- 日常维护（每月执行一次，可放入调度任务）：
-- 1. 新增下月分区
-- ALTER TABLE sched_task_log REORGANIZE PARTITION p_future INTO (
--     PARTITION p202602 VALUES LESS THAN ('2026-03-01'),
--     PARTITION p_future VALUES LESS THAN (MAXVALUE)
-- );
--
-- 2. 删除过期分区（如保留12个月，删除2025-01及更早）
-- ALTER TABLE sched_task_log DROP PARTITION p202401, p202402, p202403;
--
-- 说明：
-- - 分区键使用 start_time（任务开始时间），与保留期清理口径一致；
-- - 应用层的 log.retention_days 清理仍然生效，分区主要用于快速归档/删除；
-- - 若使用 AUTO 主键且不改造主键，MySQL 无法直接按 start_time 分区，
--   请先按方案二改造主键。
