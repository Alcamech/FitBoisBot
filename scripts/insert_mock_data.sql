-- Mock Data for FitBoisBot Testing
-- Run this script to populate the database with test data

-- Insert mock groups (test group)
INSERT IGNORE INTO `groups` (id, timezone, created_at, updated_at) VALUES 
(-834417029, 'America/New_York', NOW(), NOW());

-- Insert mock users
INSERT IGNORE INTO `users` (id, name, group_id, created_at, updated_at) VALUES 
(707903149, 'Lawton', -834417029, NOW(), NOW()),
(123456789, 'Alice', -834417029, NOW(), NOW()),
(987654321, 'Bob', -834417029, NOW(), NOW()),
(555666777, 'Charlie', -834417029, NOW(), NOW()),
(888999000, 'Diana', -834417029, NOW(), NOW());

-- Insert mock activities (spread across different days and months)
INSERT IGNORE INTO `activities` (user_id, group_id, message_id, activity, month, day, year, created_at, updated_at) VALUES 
-- Lawton's activities
(707903149, -834417029, NULL, 'deadlift', '09', '01', '2025', '2025-09-01 12:00:00', '2025-09-01 12:00:00'),
(707903149, -834417029, NULL, 'squat', '09', '03', '2025', '2025-09-03 12:00:00', '2025-09-03 12:00:00'),
(707903149, -834417029, NULL, 'bench', '09', '05', '2025', '2025-09-05 12:00:00', '2025-09-05 12:00:00'),
(707903149, -834417029, NULL, 'run', '09', '07', '2025', '2025-09-07 12:00:00', '2025-09-07 12:00:00'),

-- Alice's activities  
(123456789, -834417029, NULL, 'yoga', '09', '01', '2025', '2025-09-01 12:00:00', '2025-09-01 12:00:00'),
(123456789, -834417029, NULL, 'pilates', '09', '02', '2025', '2025-09-02 12:00:00', '2025-09-02 12:00:00'),
(123456789, -834417029, NULL, 'cardio', '09', '04', '2025', '2025-09-04 12:00:00', '2025-09-04 12:00:00'),
(123456789, -834417029, NULL, 'strength', '09', '06', '2025', '2025-09-06 12:00:00', '2025-09-06 12:00:00'),
(123456789, -834417029, NULL, 'swimming', '09', '07', '2025', '2025-09-07 12:00:00', '2025-09-07 12:00:00'),

-- Bob's activities (most active)
(987654321, -834417029, NULL, 'crossfit', '09', '01', '2025', '2025-09-01 12:00:00', '2025-09-01 12:00:00'),
(987654321, -834417029, NULL, 'running', '09', '02', '2025', '2025-09-02 12:00:00', '2025-09-02 12:00:00'),
(987654321, -834417029, NULL, 'cycling', '09', '03', '2025', '2025-09-03 12:00:00', '2025-09-03 12:00:00'),
(987654321, -834417029, NULL, 'lifting', '09', '04', '2025', '2025-09-04 12:00:00', '2025-09-04 12:00:00'),
(987654321, -834417029, NULL, 'boxing', '09', '05', '2025', '2025-09-05 12:00:00', '2025-09-05 12:00:00'),
(987654321, -834417029, NULL, 'basketball', '09', '06', '2025', '2025-09-06 12:00:00', '2025-09-06 12:00:00'),
(987654321, -834417029, NULL, 'tennis', '09', '07', '2025', '2025-09-07 12:00:00', '2025-09-07 12:00:00'),

-- Charlie's activities  
(555666777, -834417029, NULL, 'hiking', '09', '02', '2025', '2025-09-02 12:00:00', '2025-09-02 12:00:00'),
(555666777, -834417029, NULL, 'climbing', '09', '05', '2025', '2025-09-05 12:00:00', '2025-09-05 12:00:00'),

-- Diana's activities
(888999000, -834417029, NULL, 'dancing', '09', '01', '2025', '2025-09-01 12:00:00', '2025-09-01 12:00:00'),
(888999000, -834417029, NULL, 'zumba', '09', '04', '2025', '2025-09-04 12:00:00', '2025-09-04 12:00:00'),
(888999000, -834417029, NULL, 'martial-arts', '09', '07', '2025', '2025-09-07 12:00:00', '2025-09-07 12:00:00');

-- Insert mock GG counts (fastest GG competition)
INSERT IGNORE INTO `ggs` (user_id, group_id, year, fast_gg_count, created_at, updated_at) VALUES 
(987654321, -834417029, '2025', 15, NOW(), NOW()),  -- Bob is fastest
(123456789, -834417029, '2025', 8, NOW(), NOW()),   -- Alice second  
(707903149, -834417029, '2025', 5, NOW(), NOW()),   -- Lawton third
(888999000, -834417029, '2025', 3, NOW(), NOW()),   -- Diana
(555666777, -834417029, '2025', 2, NOW(), NOW());   -- Charlie

-- Insert mock token balances (rewards system)
INSERT IGNORE INTO `tokens` (user_id, group_id, year, balance, created_at, updated_at) VALUES 
(987654321, -834417029, '2025', 250, NOW(), NOW()),  -- Bob highest (most active)
(123456789, -834417029, '2025', 150, NOW(), NOW()),  -- Alice second
(707903149, -834417029, '2025', 100, NOW(), NOW()),  -- Lawton 
(888999000, -834417029, '2025', 75, NOW(), NOW()),   -- Diana
(555666777, -834417029, '2025', 50, NOW(), NOW());   -- Charlie

-- Display summary of inserted data
SELECT 'Mock data inserted successfully!' as status;
SELECT 'Groups' as table_name, COUNT(*) as count FROM `groups` WHERE id = -834417029
UNION ALL
SELECT 'Users', COUNT(*) FROM `users` WHERE group_id = -834417029  
UNION ALL
SELECT 'Activities', COUNT(*) FROM `activities` WHERE group_id = -834417029
UNION ALL  
SELECT 'GGs', COUNT(*) FROM `ggs` WHERE group_id = -834417029
UNION ALL
SELECT 'Tokens', COUNT(*) FROM `tokens` WHERE group_id = -834417029;