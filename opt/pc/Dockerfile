FROM golang:alpine

RUN mkdir /app
WORKDIR /app

COPY . .
COPY opt/pc/entrypoint.sh .
RUN apk update
RUN apk add --no-cache make

ENV BRK_USERNAME=summer
ENV BRK_PASSWORD=SuperSecure
ENV DB_USERNAME=summer
ENV DB_PASSWORD=SuperSecure

RUN make bpc
RUN chmod +x entrypoint.sh
CMD ["/app/entrypoint.sh"]