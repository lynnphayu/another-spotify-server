FROM golang:1.16-alpine

WORKDIR /app

COPY . .
RUN go mod download
RUN go mod tidy

RUN go build ./boot/server

EXPOSE 8090

CMD [ "./server" ]