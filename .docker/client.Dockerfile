FROM golang:alpine

# Prepare testing environment with build tooling.
RUN apk update && apk --no-cache add gcc g++ git make postgresql-client

WORKDIR /src

# Gather source dependencies.
COPY go.mod /src
RUN go mod download && go mod verify

COPY . /src
# Check if the source and tests will actually compile at all. Should quit early
# if one of them doesn't work.
RUN go build -v . && go test -c .

COPY .docker/client.sh /
ENTRYPOINT /client.sh
