FROM node:19-alpine AS webui-build

WORKDIR /usr/src/radius-webui

COPY radius-webui/package*.json ./

RUN npm install

COPY radius-webui/vite.config.js radius-webui/index.html ./
COPY radius-webui/src/ ./src

RUN npm run build

FROM golang:latest AS server-build

WORKDIR /usr/src/radius-server
COPY radius-server/go.* .

RUN go mod download

COPY radius-server/cmd cmd/
COPY radius-server/internal internal/

RUN CGO_ENABLED=0 go build -o radius-server cmd/radius-server/main.go

FROM scratch

WORKDIR /

COPY --from=server-build /usr/src/radius-server/radius-server .
COPY --from=webui-build /usr/src/radius-webui/dist ./dist

ENV HOST=0.0.0.0
ENV PORT=80

EXPOSE 80

CMD [ "/radius-server" ]
