FROM golang:latest

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN sed -i -e 's/\r$//' *.sh
RUN chmod +x *.sh

# go dependencies
RUN go mod download -x

# swagger
RUN go get -u github.com/swaggo/swag/cmd/swag
RUN make swag

# build go app
RUN make build

CMD ["./bin/app.out"]
