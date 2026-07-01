ALTER TABLE users ADD COLUMN email varchar(255) NOT NULL;
ALTER TABLE users ADD CONSTRAINT up_users_email UNIQUE(email);
ALTER TABLE users ADD COLUMN password_hash varchar(255) NOT NULL;