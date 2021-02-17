-- database name, user and password is for development environment only

CREATE DATABASE gingoadmin_db WITH ENCODING = 'UTF8';

BEGIN;

DO $$
BEGIN
    IF NOT EXISTS (SELECT 1 FROM pg_roles WHERE rolname = 'gingoadmin_user') THEN
        CREATE USER gingoadmin_user WITH ENCRYPTED PASSWORD 'gingoadmin_password';
    END IF;
END
$$;

GRANT ALL privileges ON DATABASE gingoadmin_db to gingoadmin_user;

COMMIT;