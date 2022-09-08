-- Import SQL file to create tables

CREATE TABLE buckets (
  id serial PRIMARY KEY,
  url text,
  created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE requests (
  id serial PRIMARY KEY,
  bucket_id integer references buckets(id),
  request_type varchar(10),
  mongo_document_ref text
);