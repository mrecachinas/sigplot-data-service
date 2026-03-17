FROM node:18-alpine AS jsbuilder

WORKDIR /app

COPY ui/webapp/package.json ui/webapp/package-lock.json ./

RUN npm ci

COPY ui/webapp .

RUN npm run build

FROM golang:1.26-alpine AS gobuilder

RUN apk add --no-cache make

WORKDIR /opt/sds/app

COPY go.mod go.sum Makefile ./
COPY cmd/ cmd/
COPY internal/ internal/
COPY ui/sds_ui.go ui/sds_ui_stub.go ui/
COPY vendor/ vendor/

# Copy the built JS app from the previous stage
COPY --from=jsbuilder /app/dist ./ui/webapp/dist

RUN CGO_ENABLED=0 go build -a -ldflags '-w -extldflags "-static"' -tags ui -mod vendor cmd/sds/sigplot_data_service.go

FROM busybox:1.36

WORKDIR /opt/sds

COPY --from=gobuilder /opt/sds/app/sigplot_data_service .

EXPOSE 5055

ENTRYPOINT ["/opt/sds/sigplot_data_service"]
