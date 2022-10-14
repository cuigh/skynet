# ---- Build UI----
FROM node:alpine AS node
WORKDIR /app
COPY ui .
RUN yarn install
#RUN yarn run lint && yarn run test
RUN yarn run build

# ---- Build Go----
FROM golang:1.19-alpine AS golang
WORKDIR /app
COPY --from=node /app/dist ui/dist
COPY . .
#RUN apk update && apk add git
RUN go build

# ---- Release ----
FROM alpine
LABEL maintainer="cuigh <noname@live.com>"
WORKDIR /app
COPY --from=golang /app/skynet .
COPY --from=golang /app/config config/
EXPOSE 8001
ENTRYPOINT ["/app/skynet"]