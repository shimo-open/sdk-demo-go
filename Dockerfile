ARG APP=sdk-demo-go

## Build frontend dist files
FROM node:18-alpine3.17 AS ui-builder

# Override with -build-arg USE_NEW_FILE_PERMISSION=true
ARG USE_NEW_FILE_PERMISSION=false
WORKDIR ${APP}/ui

COPY ui/package.json ui/.npmrc ./

RUN yarn config set registry https://registry.npmmirror.com \
    && yarn install

ENV NODE_ENV production
ENV BASE "/sdk/demo/"
ENV PUBLIC_PATH "/sdk/demo/"
ENV USE_NEW_FILE_PERMISSION=${USE_NEW_FILE_PERMISSION}
ENV CI_COMMIT_SHORT_SHA "v1"

COPY ui .

RUN yarn run build


## Build Go API Server binary
FROM golang:1.24-alpine3.22 AS api-builder

ARG APP
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn,direct
ENV TZ=Asia/Shanghai
WORKDIR ${APP}

ADD go.mod .
ADD go.sum .
RUN go mod download -x

COPY --from=ui-builder /ui/dist ./ui/dist/

COPY . .
ARG TARGETOS
ARG TARGETARCH

RUN apk add make bash && ls -rtlh && CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH make build \
    && mkdir /data/ \
    && mv ./bin/*  /data/ \
    && mv ./config  /data/ \
    && mv ./resources /data/


# Package compiled binary and configuration files into final runtime image
FROM alpine:3.22 as final

ARG APP
ENV APP=${APP}
ENV WORKDIR=/data
# Adjust to your desired timezone
ENV TZ=Asia/Shanghai

WORKDIR ${WORKDIR}
RUN apk add --no-cache lz4-dev curl tzdata
COPY --from=api-builder ${WORKDIR}/ ./

EXPOSE 9301
EXPOSE 9303

CMD ["sh", "-c", "./${APP} server"]
