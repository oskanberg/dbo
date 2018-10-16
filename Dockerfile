FROM golang:alpine as builder
RUN mkdir /build 
RUN apk add --no-cache git
ADD . /build/
WORKDIR /build/cmd/dbo
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-extldflags "-static"' -o main

FROM scratch
COPY --from=builder /build/cmd/dbo/main /app/
WORKDIR /app
EXPOSE 80
CMD ["./main", "-serve", "-port", "80"]