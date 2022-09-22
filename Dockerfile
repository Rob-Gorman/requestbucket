## Build static frontend
FROM node AS static

WORKDIR /react

COPY ./frontend .

RUN npm i

RUN npm run build

## App
FROM node AS requestbucket

WORKDIR /requestbucket

COPY --from=static /react/build ./build

COPY ./requestbucket .

RUN npm i

CMD ["node", "index.js"]

