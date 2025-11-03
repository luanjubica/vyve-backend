-- Fix Alice's password hash
-- This will set the password to what Bob and Carol currently have
-- (which appears to be working for them)

UPDATE users 
SET password_hash = '$2a$10$N9qo8uLOickgx2ZMRZo5i.uG8vH0Hbd.Zo9/6uzRJpcDO/0u5G/'
WHERE username = 'alice';

-- Verify the update
SELECT username, 
       email, 
       LENGTH(password_hash) as hash_length,
       SUBSTRING(password_hash, 1, 15) as hash_prefix
FROM users 
WHERE username IN ('alice', 'bob', 'carol');
