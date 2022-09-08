## Go script to cull old request buckets
## Set up instructions without Go installation

# Install Go (Ubuntu)
- Download file wherever you want:
  - Ubuntu: `curl -OL https://go.dev/dl/go1.18.3.linux-amd64.tar.gz`
  - MacOS: :shrug: AMD or ARM? https://go.dev/dl/ to decide.
    - `curl -OL https://go.dev/dl/go1.18.3.darwin-amd64.pkg`
    - `curl -OL https://go.dev/dl/go1.18.3.darwin-arm64.pkg`
- Extract contents (My suggested destination) (Ubuntu):
  - `tar -C $HOME -xvf go1.18.3.linux-amd64.tar.gz`
- Add Go to $PATH:
  - `echo "export PATH=$PATH:$HOME/go/bin" >> $HOME/.profile`
  - `source $HOME/.profile`
  - verify install with `go version` giving you appropriate response

# Set up script
- put `requestbucket` Go script folder in `$HOME/go/src/.`
- cd to script directory and run `go install`

# Implement ON CASCADE DELETE constraint on postgres
this just constrains the request records to be deleted when parent bucket is
`psql -d requestbucket < requestuckets/schemamigrate.sql`

# Set up .env file
- Script runs on a handful of environment variables you can see in `environment.go`
- make .env file in environment subfolder: `touch $HOME/go/src/requestbucket/environment/.env`
- Something like this, fill in relevant missing pieces:
    > HOST=localhost
      PORT=5432
      USER=
      PASSWORD=""
      PGDBNAME=requestbucket
      PGTABLE=buckets
      LOGFILE="$HOME/go/src/requestbucket/deletelog.md"
      MONGODB_URI=mongodb://localhost:27017/
      MONGODB=
      MONGODB_COLL=requests
- Script is not meant to utilize postgres password, so if you figured that out, you'll just have to add the parameter to the connection on line 32

# Set up cron
- set up the crontab file with `crontab -e`
- append the job to the bottom - will link to project bash script
> `* */1 * * * /home/rjg/go/src/requestbucket/cron.sh` to run every hour