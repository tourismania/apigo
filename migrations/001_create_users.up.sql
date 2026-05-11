CREATE TABLE "users" (
    id          SERIAL PRIMARY KEY,
    uuid        UUID NOT NULL,
    first_name  VARCHAR(255),
    last_name   VARCHAR(255),
    email       VARCHAR(255) NOT NULL,
    created_at  TIMESTAMP NOT NULL,
    updated_at  TIMESTAMP NOT NULL,
    phone       VARCHAR(255),
    password    VARCHAR(255) NOT NULL,
    is_active   BOOLEAN NOT NULL
);
