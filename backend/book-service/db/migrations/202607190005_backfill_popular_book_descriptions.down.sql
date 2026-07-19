-- No-op: migration 202607190004 now seeds the same descriptions.
-- Removing them here would corrupt a fresh database rolled back to version 004.
SELECT 1;
