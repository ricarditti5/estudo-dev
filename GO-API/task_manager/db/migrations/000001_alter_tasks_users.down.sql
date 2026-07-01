ALTER TABLE users DROP COLUMN password_hash;
ALTER TABLE users DROP CONSTRAINT up_users_email;
ALTER TABLE users DROP COLUMN email;