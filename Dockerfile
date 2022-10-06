## Build static frontend
FROM node AS static

WORKDIR /react

COPY ./frontend .

RUN npm i

RUN npm run build

## App
FROM alpine AS requestbucket

RUN apk add --update nodejs npm

WORKDIR /requestbucket

COPY --from=static /react/build ./build

COPY ./requestbucket .

RUN npm ci

CMD ["node", "index.js"]

