// Helper to write to PostgreSQL

const config = require("./config");
const { Client } = require("pg");

const logQuery = (statement, parameters) => {
  let timeStamp = new Date();
  let formattedTimeStamp = timeStamp.toString().substring(4, 24);
  console.log(formattedTimeStamp, statement, parameters);
};

const PG_ENV = {
  user: config.PGUSER,
  password: config.PGPASSWORD,
  host: config.PGHOST,
  port: config.PGPORT,
  database: config.PGDATABASE,
  ssl: false
}

const CONNECTION = {
  connectionString: config.DATABASE_URL,
  ssl: false,
};

module.exports = {
  async dbQuery(statement, ...parameters) {
    let client = new Client(PG_ENV);

    await client.connect();
    logQuery(statement, parameters);
    let result = await client.query(statement, parameters);
    await client.end();

    return result;
  }
};