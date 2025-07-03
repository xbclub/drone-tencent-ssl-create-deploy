from golang:1.24-alpine as prebuild
add . /data
workdir /data
run go mod tidy && go build -v -o app -ldflags "-s -w" --trimpath

from alpine:latest
COPY --from=prebuild /data/app /app/
workdir /app
ENTRYPOINT ["./app"]