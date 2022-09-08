const mongoose = require('mongoose')

// Schema for request document
const requestSchema = new mongoose.Schema({
  headers: {
    type: Object,
  },
  payload: {
    type: Object,
  },
  requestType: {
    type: String,
  }
})

requestSchema.set('toJSON', {
  transform: (document, returnedObject) => {
    returnedObject.id = returnedObject._id.toString()
    delete returnedObject._id
    delete returnedObject.__v
  }
})

module.exports = mongoose.model('Request', requestSchema)
