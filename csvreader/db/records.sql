begin;
CREATE TABLE if not exists public.records (
	id varchar(50) NULL CONSTRAINT recordspk PRIMARY KEY,
	first_name varchar(150) NULL,
	last_name varchar(150) NULL,
	email varchar(150) NULL,
	phone varchar(50) NULL
);
commit;
