#!/bin/bash
set -e

echo "---U/P/B-----------"
echo $POSTGRES_USER
echo $POSTGRES_PASSWORD
echo $PGDATABASE
echo "-------------------"

# Create the record's table if it does not exist
psql postgres://csvreader:csvreader@db/csvreader?sslmode=disable <<-EOSQL
begin;
CREATE TABLE if not exists public.records (
	id varchar(50) NULL CONSTRAINT recordspk PRIMARY KEY,
	first_name varchar(150) NULL,
	last_name varchar(150) NULL,
	email varchar(150) NULL,
	phone varchar(50) NULL
);
commit;
EOSQL

# start the csv reader
go run csvreader/main.go

