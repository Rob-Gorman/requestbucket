# Requestbin

## Pages:
- Home page
- Created page
- Bin page
- Bin History page

## Requirements:
- User requests Bin (GET)
  - Event button triggers function to generate URL???
  - Function: generate a URL route
  - Function: store the bin URL route to Postgres & metadata
  - Send Response back to user with bin URL

- Webhook sends HTTP request to Bin URL (any request)
  - Write an Express route logic that takes bin URL as parameter
  - Check if request count for this bin exceeds 20; if so, remove least recent requests until you only have 19
  - Store request information in Requests table + Mongo; Return 200 OK
    - Mongo Query write new record: id, headers, json_payload, raw_data
      - this returns the mongo id to use below vvv
    - Postgres Query for new record: provide bin id, request type, date/time, mongo reference ID
  - GET: Return 200 OK
  - Redirect to bin history page
  
- User views Bin History
  - "/${binURL}/history" route
  - query Postgres for last 20 requests related to bin_url
  - query mongo for data based on request_ids
  - deserialize & return filtered collection of mongo documents as JS object
  - Parse JS object into headers & payload as two separate variables and return in response

## Front End requirements (react):
1. View bin history
  - JS AJAX request based on bin url
  - Loop through response object and populate table-like component
2. Home page
  - create bin button

## Data Maintenance:
  - CRON job: checking life of bin
  - Schedule: 9:00 AM UTC
  - Calculate date/timestamp for 48 hours prior
  - Query Postgres Bins table for bin_ids > 48 hours and store in array
  - if array length > 0, run mongoose command to delete all bins

  - consider different methods of having the job check for and delete expired bins to account for scale (db is on hold while cron job is sweeping)

## Data structures:
- Bin is an object
  - holds array of requests