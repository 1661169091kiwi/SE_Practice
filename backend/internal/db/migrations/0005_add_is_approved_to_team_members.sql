-- 添加 is_approved 字段到 team_members 表
-- 注意：如果字段已存在，执行此脚本会报错，可以安全地忽略

-- 检查并添加 is_approved 字段
-- MySQL 不支持 IF NOT EXISTS，需要手动检查或忽略错误
ALTER TABLE team_members 
ADD COLUMN is_approved BOOLEAN DEFAULT FALSE;

-- 更新现有记录：如果 is_active 为 TRUE，则设置为已批准（兼容旧数据）
UPDATE team_members 
SET is_approved = TRUE 
WHERE is_active = TRUE;

