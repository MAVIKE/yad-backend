FROM golang:latest AS build

RUN go version

COPY ./ /github.com/MAVIKE/yad-backend
WORKDIR /github.com/MAVIKE/yad-backend

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

# copy config
RUN if [ "$MODE" = "stage" ]; then \
    make stage_config; \
    else \
    make config; \
    fi

# build go app
RUN CGO_ENABLED=0 GOOS=linux go build -o ./bin/app.out ./cmd/app/main.go

FROM alpine:latest

RUN apk --no-cache add ca-certificates
WORKDIR /root/

RUN apk add postgresql

COPY --from=0 /github.com/MAVIKE/yad-backend/bin/app.out .
COPY --from=0 /github.com/MAVIKE/yad-backend/configs/ ./configs/
COPY --from=0 /github.com/MAVIKE/yad-backend/docs/ ./docs/
COPY --from=0 /github.com/MAVIKE/yad-backend/img/ ./img/

CMD ["./app.out"]
