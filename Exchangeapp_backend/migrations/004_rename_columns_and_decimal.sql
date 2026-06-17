-- 重命名字段并将会计相关列改为 DECIMAL(18,4) / DECIMAL(10,4)
-- 执行前请备份数据库
-- 前置依赖：建议已执行 001/002/003（尤其 003 添加了 column_zhuang_zhan_bi）

-- ========== 数据清洗：table_yanchendao1 ==========
UPDATE table_yanchendao1
SET column_benjin = '0'
WHERE column_benjin IS NULL OR TRIM(column_benjin) = '';

UPDATE table_yanchendao1
SET column_yongjin = '0'
WHERE column_yongjin IS NULL OR TRIM(column_yongjin) = '';

UPDATE table_yanchendao1
SET column_mean = '0'
WHERE column_mean IS NULL OR TRIM(column_mean) = '';

UPDATE table_yanchendao1
SET column_restart_index = '0'
WHERE column_restart_index IS NULL
   OR TRIM(column_restart_index) = ''
   OR column_restart_index NOT REGEXP '^-?[0-9]+$';

UPDATE table_yanchendao1
SET column_liushui_index = '0'
WHERE column_liushui_index IS NULL
   OR TRIM(column_liushui_index) = ''
   OR column_liushui_index NOT REGEXP '^-?[0-9]+$';

-- ========== 数据清洗：table_yanchendao2 ==========
UPDATE table_yanchendao2
SET colmun_shuyingzhi = TRIM(LEADING '+' FROM colmun_shuyingzhi)
WHERE colmun_shuyingzhi LIKE '+%';

UPDATE table_yanchendao2
SET colmun_shuyingzhi_d = TRIM(LEADING '+' FROM colmun_shuyingzhi_d)
WHERE colmun_shuyingzhi_d LIKE '+%';

UPDATE table_yanchendao2
SET column_xiazhujine = TRIM(LEADING '+' FROM column_xiazhujine)
WHERE column_xiazhujine LIKE '+%';

UPDATE table_yanchendao2
SET column_current_jin = TRIM(LEADING '+' FROM column_current_jin)
WHERE column_current_jin LIKE '+%';

UPDATE table_yanchendao2
SET column_xiazhujine = '0'
WHERE column_xiazhujine IS NULL OR TRIM(column_xiazhujine) = '';

UPDATE table_yanchendao2
SET colmun_shuyingzhi = '0'
WHERE colmun_shuyingzhi IS NULL OR TRIM(colmun_shuyingzhi) = '';

UPDATE table_yanchendao2
SET column_current_jin = '0'
WHERE column_current_jin IS NULL OR TRIM(column_current_jin) = '';

-- 原列 colmun_shuyingzhi_d 为 NOT NULL，须先改为可空，才能把空串清成 NULL
ALTER TABLE table_yanchendao2
  MODIFY COLUMN colmun_shuyingzhi_d VARCHAR(255) NULL DEFAULT NULL;

UPDATE table_yanchendao2
SET colmun_shuyingzhi_d = NULL
WHERE colmun_shuyingzhi_d IS NULL OR TRIM(colmun_shuyingzhi_d) = '';

-- ========== table_yanchendao1：重命名 + 改类型 ==========
ALTER TABLE table_yanchendao1
  CHANGE COLUMN column_benjin benjin DECIMAL(18, 4) NOT NULL COMMENT '本金',
  CHANGE COLUMN column_yongjin yongjin DECIMAL(10, 4) NOT NULL COMMENT '佣金',
  CHANGE COLUMN column_mean mean DECIMAL(10, 4) NOT NULL COMMENT '数学期望',
  CHANGE COLUMN column_restart_index restart_index INT NOT NULL DEFAULT 0 COMMENT '重起位置',
  CHANGE COLUMN column_liushui_index liushui_index INT NOT NULL DEFAULT 0 COMMENT '流水的位置';

-- 庄占比：若已执行 003 则重命名；若报错 Unknown column，改执行下方注释中的 ADD 语句
ALTER TABLE table_yanchendao1
  CHANGE COLUMN column_zhuang_zhan_bi zhuang_zhan_bi INT NOT NULL DEFAULT 50 COMMENT '庄占比(0-100)';
-- ALTER TABLE table_yanchendao1
--   ADD COLUMN zhuang_zhan_bi INT NOT NULL DEFAULT 50 COMMENT '庄占比(0-100)' AFTER liushui_index;

-- ========== table_yanchendao2：重命名 + 改类型 ==========
ALTER TABLE table_yanchendao2
  CHANGE COLUMN column_xiazhujine xiazhujine DECIMAL(18, 4) NOT NULL COMMENT '下注的金额',
  CHANGE COLUMN colmun_shuyingzhi shuyingzhi DECIMAL(18, 4) NOT NULL COMMENT '输赢值',
  CHANGE COLUMN colmun_shuyingzhi_d shuyingzhi_xiaoshu DECIMAL(18, 4) NULL DEFAULT NULL COMMENT '消数后的输赢值',
  CHANGE COLUMN colmun_shengfulu shengfulu VARCHAR(10) NOT NULL COMMENT '胜负路',
  CHANGE COLUMN colmun_zx zx VARCHAR(10) NOT NULL COMMENT '开出庄闲',
  CHANGE COLUMN colmun_remark remark VARCHAR(255) NULL COMMENT '输赢标记备注',
  CHANGE COLUMN column_current_jin current_jin DECIMAL(18, 4) NOT NULL COMMENT '当前金额';
