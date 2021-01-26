FROM golang:rc-alpine3.13 AS build
WORKDIR /src
COPY . .
RUN export CGO_ENABLED=0 && go build -o /out/aka .
FROM scratch AS bin
COPY --from=build /out/aka /app/aka
WORKDIR /app
CMD ["./aka"]
