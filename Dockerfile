FROM golang:1.24

WORKDIR /usr/src/turbo_snail

COPY go.mod go.sum ./
RUN go mod download

COPY . .

EXPOSE 7777 7000 

RUN go build -v -o /usr/local/bin/turbo_snail

CMD ["turbo_snail"]
