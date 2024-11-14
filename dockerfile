FROM golang:1.23.2 as build

WORKDIR /app

COPY . .

RUN go mod tidy && go test ./... -v
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o weather-api ./cmd/api

FROM scratch

WORKDIR /app

COPY --from=alpine /etc/ssl/certs /etc/ssl/certs
COPY --from=build /app/weather-api .

EXPOSE 8080
ENTRYPOINT [ "./weather-api" ]