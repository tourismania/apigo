ALTER TABLE "users"
    ADD COLUMN login VARCHAR(255);

UPDATE "users" SET login = email WHERE login IS NULL;

ALTER TABLE "users"
    ALTER COLUMN login SET NOT NULL;

CREATE UNIQUE INDEX IF NOT EXISTS users_email_unique_idx ON "users" (email);
