-- Verify eurus-backend-db-auth:user_sessions_alter_token on pg

BEGIN;

-- XXX Add verifications here.
SELECT character_maximum_length / CASE  
WHEN  character_maximum_length < 1024 THEN 0
ELSE 1
END
  FROM information_schema.columns
 WHERE table_schema = 'public'
   AND table_name   = 'user_sessions'
     and column_name  = 'token';
     
ROLLBACK;
