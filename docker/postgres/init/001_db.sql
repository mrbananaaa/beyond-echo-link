CREATE EXTENSION IF NOT EXISTS dblink;

DO $$
BEGIN
   IF NOT EXISTS (SELECT FROM pg_database WHERE datname = 'beldb') THEN
      PERFORM dblink_exec('dbname=postgres', 'CREATE DATABASE beldb');
   END IF;
END
$$;