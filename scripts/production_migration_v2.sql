-- Production Migration from Old Schema to New Schema
-- Handles migration from original schema to current schema with timestamps and message_id

START TRANSACTION;

-- Add timestamp columns to all tables that need them
ALTER TABLE `groups` ADD COLUMN IF NOT EXISTS `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `groups` ADD COLUMN IF NOT EXISTS `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);

ALTER TABLE `users` ADD COLUMN IF NOT EXISTS `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `users` ADD COLUMN IF NOT EXISTS `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);

ALTER TABLE `ggs` ADD COLUMN IF NOT EXISTS `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `ggs` ADD COLUMN IF NOT EXISTS `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);

ALTER TABLE `tokens` ADD COLUMN IF NOT EXISTS `created_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3);
ALTER TABLE `tokens` MODIFY COLUMN `updated_at` datetime(3) DEFAULT CURRENT_TIMESTAMP(3) ON UPDATE CURRENT_TIMESTAMP(3);

-- Add message_id and timestamp columns to activities table (from old schema)
ALTER TABLE `activities` ADD COLUMN `message_id` INT NULL;
ALTER TABLE `activities` ADD COLUMN `created_at` datetime(3) NULL;
ALTER TABLE `activities` ADD COLUMN `updated_at` datetime(3) NULL;

-- Populate message_id with dummy values for historical data
UPDATE activities SET message_id = (-1000000 - id) WHERE message_id IS NULL;

-- Make message_id NOT NULL and add unique constraint
ALTER TABLE activities MODIFY COLUMN message_id INT NOT NULL;
ALTER TABLE activities ADD UNIQUE KEY unique_message (user_id, group_id, message_id);

-- Backfill timestamps with actual activity dates at noon
UPDATE activities 
SET created_at = STR_TO_DATE(CONCAT(year, '-', month, '-', day, ' 12:00:00'), '%Y-%m-%d %H:%i:%s'),
    updated_at = STR_TO_DATE(CONCAT(year, '-', month, '-', day, ' 12:00:00'), '%Y-%m-%d %H:%i:%s')
WHERE created_at IS NULL OR updated_at IS NULL;

COMMIT;

-- Verification queries (run separately after migration)
-- SELECT COUNT(*) as total_activities, COUNT(CASE WHEN message_id < 0 THEN 1 END) as dummy_ids FROM activities;
-- SELECT COUNT(*) as activities_with_timestamps FROM activities WHERE created_at IS NOT NULL AND updated_at IS NOT NULL;
-- DESCRIBE activities;