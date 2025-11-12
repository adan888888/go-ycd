-- 添加软删除字段 deleted_at
-- 执行时间: 2025-11-12

-- 为 table_yanchendao1 表添加 deleted_at 字段
ALTER TABLE `table_yanchendao1` 
ADD COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '删除时间' AFTER `uid`,
ADD INDEX `idx_table_yanchendao1_deleted_at` (`deleted_at`);

-- 为 table_yanchendao2 表添加 deleted_at 字段
ALTER TABLE `table_yanchendao2` 
ADD COLUMN `deleted_at` DATETIME(3) NULL DEFAULT NULL COMMENT '删除时间' AFTER `user_id`,
ADD INDEX `idx_table_yanchendao2_deleted_at` (`deleted_at`);

