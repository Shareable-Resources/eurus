-- Deploy eurus-backend-db-admin:20211006_create_admin_user to pg



CREATE OR REPLACE PROCEDURE createAdminUsersTable()
LANGUAGE plpgsql
AS $$
DECLARE
	is_exists INTEGER;
BEGIN


   is_exists := (SELECT  1 FROM information_schema.tables 
   WHERE  table_schema = 'public'
   AND    table_name   = 'admin_users');
  
	IF is_exists = 1 then 
		CREATE TEMP TABLE temp_admin_users AS
		SELECT *
		FROM admin_users;

		ALTER table admin_users RENAME TO admin_users_old;
	end if;
	
	CREATE TABLE IF NOT EXISTS admin_users (
	 	id serial PRIMARY KEY,
	 	username varchar(50) NOT NULL,
	 	password varchar(512) NOT NULL,
	 	secret varchar(512),
	 	status smallint NOT NULL
 	);

 	CREATE UNIQUE INDEX admin_users_username_idx ON admin_users (username);
 	
 	IF is_exists = 1 then
 		insert into admin_users (username, password, status) select username, password, 1 from temp_admin_users;
	end if;
END;
$$;

BEGIN;

call createAdminUsersTable();

drop procedure createAdminUsersTable;
COMMIT;
