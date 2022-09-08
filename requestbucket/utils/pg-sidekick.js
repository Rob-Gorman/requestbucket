const { dbQuery } = require("./db-query");
// PgSidekick class defines methods for making asynchronous PostgreSQL queries

module.exports = class PgSidekick {
  async loadBuckets() {
    const LOAD_BINS = 'SELECT * FROM buckets'
    let buckets = await dbQuery(LOAD_BINS)
    return buckets.rows
  }
  
  async addBucket(url) {
    const ADD_BUCKET = `INSERT INTO buckets (url) VALUES ('${url}')`;
    await dbQuery(ADD_BUCKET);
    return true;
  }

  async loadRequests(bucketUrl) {
    // returns all requests for given bucketUrl
    const LOAD_REQUESTS = `SELECT * FROM requests WHERE bucket_id =(SELECT id FROM buckets WHERE url=('${bucketUrl}'))`;
   
    let requests = await dbQuery(LOAD_REQUESTS)
    return requests.rows
  }

  async getBucketId(bucketUrl) {
    // returns the pg id for the specified bucket url
    const LOAD_BUCKET = `SELECT id FROM buckets WHERE url='${bucketUrl}'`

    let bucket = await dbQuery(LOAD_BUCKET)
    if (bucket.rows.length == 0) {
      return false
    }
    return bucket.rows[0].id;
  }
  
  async createRequest(bucketId, requestType, mongoId) {
    const INSERT_REQUEST = `INSERT INTO requests (bucket_id, request_type, mongo_document_ref) VALUES ('${bucketId}', '${requestType}', '${mongoId}')`;

    let newRequest = await dbQuery(INSERT_REQUEST)
    return newRequest
  }

  async deleteRequest(requestId) {
    const DELETE_REQUEST = `DELETE FROM requests WHERE id=('${requestId}')`

    const deletedRequest = await dbQuery(DELETE_REQUEST)
    console.log(deletedRequest)

    return
  }
}