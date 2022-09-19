## Build static frontend
FROM node AS static

WORKDIR /react

COPY ./frontend .

RUN npm install

RUN npm build

## App
FROM node AS requestbucket

WORKDIR /requestbucket

COPY --from=static /react/build .

COPY ./requestbucket .

RUN npm install

CMD ["node", "index.js"]

