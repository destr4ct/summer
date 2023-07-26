FROM golang:alpine

RUN mkdir /app
WORKDIR /app

COPY . .
RUN apk update
RUN apk add --no-cache make

ENV DB_USERNAME=summer
ENV DB_PASSWORD=SuperSecure

RUN make tg
CMD ["/app/builds/tg"]