FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY *.go ./
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/api .

FROM scratch
COPY --from=build /out/api /api
EXPOSE 8080
ENTRYPOINT ["/api"]
