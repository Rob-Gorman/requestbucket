require('dotenv').config()


const PORT = process.env.PORT

// Mongo config
const MONGODB_URI = process.env.MONGODB_URI

module.exports = {
  MONGODB_URI,
  PORT
}
module.exports = process.env
