ALTER TABLE change_history
DROP COLUMN IF EXISTS comment;

ALTER TABLE change_history
DROP COLUMN IF EXISTS admin_comment;

ALTER TABLE change_history
DROP COLUMN IF EXISTS user_comment;