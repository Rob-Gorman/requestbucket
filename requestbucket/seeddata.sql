-- Load test data for dev purposes only. 
INSERT INTO buckets (url) VALUES 
  ('http://example.com'),
  ('google.com'),
  ('askjeeves.com')
;

INSERT INTO requests (bucket_id, request_type, mongo_document_ref) VALUES
  (4, 'GET', 'abcd'),
  (5, 'POST', 'xyz'),
  (6, 'PUT', 'defgh')
;