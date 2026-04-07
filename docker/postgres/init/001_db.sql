SELECT 'CREATE DATABASE beldb' 
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'beldb') \gexec;