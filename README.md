# RequestBucket

Application for observing web requests -- method, headers, body, etc

## Components

- Node.js webserver that generates custom endpoints and captures data from requests made to those endpoints; as well as serving the static build of the React.js frontend
- PostgreSQL database to store metadata on each endpoint
- MongoDB database to store data of every individual request made to all endpoints
- Go application serving as a cron task to cull old endpoint data from the databases

## Local Deploy

Can run locally with `bash localdeploy.sh`
This script will build the two Docker images for the Node.js server and the Go cronjob
It will also run `docker compose up` to run the required containers with the necessary env variables
May need to run as `sudo` on some distributions
