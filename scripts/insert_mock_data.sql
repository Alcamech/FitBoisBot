-- Mock Data for FitBoisBot Testing
-- Run this script to populate the database with test data
-- Dynamically uses current month and year

-- Set variables for current month and year
SET @current_month = LPAD(MONTH(NOW()), 2, '0');
SET @current_year = YEAR(NOW());

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

-- Insert mock activities (spread across different days of current month)
INSERT IGNORE INTO `activities` (user_id, group_id, message_id, activity, month, day, year, created_at, updated_at) VALUES
-- Lawton's activities
(707903149, -834417029, NULL, 'deadlift', @current_month, '01', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00')),
(707903149, -834417029, NULL, 'squat', @current_month, '03', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-3) DAY), '%Y-%m-03 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-3) DAY), '%Y-%m-03 12:00:00')),
(707903149, -834417029, NULL, 'bench', @current_month, '05', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-5) DAY), '%Y-%m-05 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-5) DAY), '%Y-%m-05 12:00:00')),
(707903149, -834417029, NULL, 'run', @current_month, '07', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00')),

-- Alice's activities
(123456789, -834417029, NULL, 'yoga', @current_month, '01', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00')),
(123456789, -834417029, NULL, 'pilates', @current_month, '02', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-2) DAY), '%Y-%m-02 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-2) DAY), '%Y-%m-02 12:00:00')),
(123456789, -834417029, NULL, 'cardio', @current_month, '04', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-4) DAY), '%Y-%m-04 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-4) DAY), '%Y-%m-04 12:00:00')),
(123456789, -834417029, NULL, 'strength', @current_month, '06', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-6) DAY), '%Y-%m-06 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-6) DAY), '%Y-%m-06 12:00:00')),
(123456789, -834417029, NULL, 'swimming', @current_month, '07', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00')),

-- Bob's activities (most active)
(987654321, -834417029, NULL, 'crossfit', @current_month, '01', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00')),
(987654321, -834417029, NULL, 'running', @current_month, '02', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-2) DAY), '%Y-%m-02 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-2) DAY), '%Y-%m-02 12:00:00')),
(987654321, -834417029, NULL, 'cycling', @current_month, '03', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-3) DAY), '%Y-%m-03 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-3) DAY), '%Y-%m-03 12:00:00')),
(987654321, -834417029, NULL, 'lifting', @current_month, '04', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-4) DAY), '%Y-%m-04 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-4) DAY), '%Y-%m-04 12:00:00')),
(987654321, -834417029, NULL, 'boxing', @current_month, '05', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-5) DAY), '%Y-%m-05 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-5) DAY), '%Y-%m-05 12:00:00')),
(987654321, -834417029, NULL, 'basketball', @current_month, '06', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-6) DAY), '%Y-%m-06 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-6) DAY), '%Y-%m-06 12:00:00')),
(987654321, -834417029, NULL, 'tennis', @current_month, '07', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00')),

-- Charlie's activities
(555666777, -834417029, NULL, 'hiking', @current_month, '02', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-2) DAY), '%Y-%m-02 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-2) DAY), '%Y-%m-02 12:00:00')),
(555666777, -834417029, NULL, 'climbing', @current_month, '05', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-5) DAY), '%Y-%m-05 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-5) DAY), '%Y-%m-05 12:00:00')),

-- Diana's activities
(888999000, -834417029, NULL, 'dancing', @current_month, '01', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-1) DAY), '%Y-%m-01 12:00:00')),
(888999000, -834417029, NULL, 'zumba', @current_month, '04', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-4) DAY), '%Y-%m-04 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-4) DAY), '%Y-%m-04 12:00:00')),
(888999000, -834417029, NULL, 'martial-arts', @current_month, '07', @current_year, DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00'), DATE_FORMAT(DATE_SUB(NOW(), INTERVAL (DAY(NOW())-7) DAY), '%Y-%m-07 12:00:00'));

-- Insert mock GG counts (fastest GG competition)
INSERT IGNORE INTO `ggs` (user_id, group_id, year, fast_gg_count, created_at, updated_at) VALUES
(987654321, -834417029, @current_year, 15, NOW(), NOW()),  -- Bob is fastest
(123456789, -834417029, @current_year, 8, NOW(), NOW()),   -- Alice second
(707903149, -834417029, @current_year, 5, NOW(), NOW()),   -- Lawton third
(888999000, -834417029, @current_year, 3, NOW(), NOW()),   -- Diana
(555666777, -834417029, @current_year, 2, NOW(), NOW());   -- Charlie

-- Insert mock token balances (rewards system)
INSERT IGNORE INTO `tokens` (user_id, group_id, year, balance, created_at, updated_at) VALUES
(987654321, -834417029, @current_year, 250, NOW(), NOW()),  -- Bob highest (most active)
(123456789, -834417029, @current_year, 150, NOW(), NOW()),  -- Alice second
(707903149, -834417029, @current_year, 100, NOW(), NOW()),  -- Lawton
(888999000, -834417029, @current_year, 75, NOW(), NOW()),   -- Diana
(555666777, -834417029, @current_year, 50, NOW(), NOW());   -- Charlie

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