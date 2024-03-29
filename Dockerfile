FROM golang:1.17

WORKDIR /go/src
COPY . .

RUN go mod tidy

ENTRYPOINT  ["go", "test", "-v", "./accounts", "./httputils", "./integration_tests", "-coverprofile", "cov.out"]