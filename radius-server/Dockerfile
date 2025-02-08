FROM golang:latest AS build

WORKDIR /usr/src/radius-server
COPY go.* .

RUN go mod download

COPY cmd cmd/
COPY internal internal/

RUN CGO_ENABLED=0 go build -o radius-server cmd/radius-server/main.go

FROM scratch

WORKDIR /

COPY --from=build /usr/src/radius-server/radius-server .

ENV HOST=0.0.0.0
ENV PORT=80

EXPOSE 80

CMD [ "/radius-server" ]
