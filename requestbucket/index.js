const express = require("express");
const app = express();
require('dotenv').config();
const path = require('path');
const PORT = process.env.PORT;
const cors = require('cors');
const mongoose = require('mongoose');
const morgan = require('morgan');
const PgSidekick = require("./utils/pg-sidekick")
const pg = new PgSidekick()
const { urlGenerator } = require("./services/url-generator");
const { sortRequests } =  require("./services/sortFunction");
const Request = require('./models/request')
const MAX_REQUESTS_PER_BIN = 20

// IO is a server engine instance that manages Sockets 
const { createServer } = require("http");
const httpServer = createServer(app)
const { Server } = require("socket.io");
const io = new Server(httpServer, {
  cors: {
    origin: '*',
    methods: ["GET", "POST"]
  }
});

// MIDDLEWARE
app.use(cors());
app.use(express.static(path.join(__dirname, 'build')));
app.use(express.json())
app.use(express.urlencoded({ extended: true }))
app.use(morgan('dev'))

// Add middleware so that all routes have access to io server
app.use((req, res, next) => {
  req.io = io;
  return next();
});

mongoose.connect(process.env.MONGODB_URI)
    .then(() => {
        console.log('connected to MongoDB');
    })
    .catch((error) => {
        console.error('error connecting to MongoDB:', error.message);
    });

// Server listens for connection events from incoming sockets 
io.on('connection', (socket) => {
  // initial connections always HTTP polling
  // upgraded to websocket protocol if handshake successful
  console.log(`--user connected: ${socket.id} via: ${socket.conn.transport.name}`); 
});

// - Home page
app.get('/', (req, res) => {
  res.sendFile(path.join(__dirname+'/static/home.html'));
})

app.post('/create', (req, res) => {
 // creates a new bucket and returns page with new bucket url
  let newUrl = urlGenerator();

  try {
    ;(async function () {
      pg.addBucket(newUrl); // add to postgres db
      res.send(newUrl); // direct to 'created' page
    })()
  } catch (error) {
    console.error(error)
  }
})

// - Bucket "page" that collects all incoming requests
app.all(`/:bucketUrl`, (req, res) => {
  const bucketUrl = req.params.bucketUrl
  const {headers, method, body } = req

  ;(async function () {
    try {
      const bucketId = await pg.getBucketId(bucketUrl);

      if (!bucketId) {
        res.status(404).send('Bucket URL not found')
      } else {

        const bucketRequests = await pg.loadRequests(bucketUrl)
        const countRequests = bucketRequests.length

        // Delete oldest request from bucket
        if (countRequests >= MAX_REQUESTS_PER_BIN) {
          bucketRequests.sort(sortRequests)
          const oldestId = bucketRequests[0]['id']
          await pg.deleteRequest(oldestId);
        }

        // Write to MongoDB & Postgres
        const requestObj = new Request({
          headers: headers,
          payload: body,
          requestType: method
        })
        
        const savedRequest = await requestObj.save()
        await pg.createRequest(bucketId, method, savedRequest._id)

        // Emit new request data to connected sockets
        req.io.emit("NEW_REQUEST_IN_BUCKET", {
          bucketUrl: bucketUrl,
          data: savedRequest
        })
      
        res.send('200')
      }

    } catch (error) {
      console.error(error);
    }
  })() 

})

// - Bin History page
app.get('/stash/:bucketUrl', (req, res) => {
  const bucketUrl = req.params.bucketUrl
  
  ;(async function () {
    try {
      const bucketId = await pg.getBucketId(bucketUrl);

      if (!bucketId) {
        res.status(404).send('Bucket Url not found')
      } else {

        const bucketRequests = await pg.loadRequests(bucketUrl)
        const mongoIds = bucketRequests.map(req => req['mongo_document_ref'])

        if (mongoIds.length == 0) {
          return res.send(`No requests in bucket ${bucketUrl}`)
        }
        
        const data = await Request.find({
          '_id': {$in: mongoIds}
        });

        return res.json(data.reverse())
      }

    } catch (error) {
      console.error(error);
    }
  })() 
  
})

app.get('*', (req, res) => {
  res.sendFile(path.join(__dirname+'/build/index.html'))
})

// Must use http server instead of express server
httpServer.listen(PORT, () => {
  console.log(`httpServer is running on port ${PORT}`)
})
