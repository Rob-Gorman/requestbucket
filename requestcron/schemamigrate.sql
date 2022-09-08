-- to purge all requests from table associated with bucket
-- added ON DELETE CASCADE constraint to requests table
ALTER TABLE requests DROP CONSTRAINT "requests_bucket_id_fkey";
ALTER TABLE requests 
  ADD CONSTRAINT "requests_bucket_id_fkey"
  FOREIGN KEY (bucket_id)
  REFERENCES buckets(id)
  ON DELETE CASCADE;
