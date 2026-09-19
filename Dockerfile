FROM golang:1.25-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/agnos-api .

FROM alpine:3.22
RUN adduser -D -H app
COPY --from=build /out/agnos-api /usr/local/bin/agnos-api
USER app
EXPOSE 5000
ENTRYPOINT ["/usr/local/bin/agnos-api"]
