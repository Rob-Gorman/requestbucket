# Build Go binary for cron executable
FROM golang AS cronexe

WORKDIR /go/src/reqcron
COPY ./requestcron .
RUN go mod download
RUN CGO_ENABLED=0 go build -o /requestcron

# Build cronjob image
FROM ubuntu

RUN apt-get -qq update && apt-get -y install -qq --force-yes cron

# --- did you know about the access to env vars in cron jobs? wew lad
COPY ./.env .
RUN cat ./.env >> /etc/environment

COPY --from=cronexe requestcron ./requestcron
RUN chmod 0744 requestcron

# Create cron file and give appropriate permissions
RUN touch /etc/cron.d/crontask
RUN chmod 0744 /etc/cron.d/crontask
# Set up cron job
RUN echo '* * * * * /requestcron >> /var/log/cron.log 2>&1' > /etc/cron.d/crontask
# Apply cronjob
RUN crontab /etc/cron.d/crontask

RUN touch /var/log/cron.log

CMD cron && tail -f /var/log/cron.log