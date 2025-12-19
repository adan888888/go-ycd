-- 添加庄占比字段
-- 执行时间: 2025-12-18

-- 为 table_yanchendao1 表添加 column_zhuang_zhan_bi 字段
ALTER TABLE `table_yanchendao1` 
ADD COLUMN `column_zhuang_zhan_bi` INT NOT NULL DEFAULT 50 COMMENT '庄占比(0-100)' AFTER `column_liushui_index`;

