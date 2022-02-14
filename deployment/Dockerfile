FROM golang:1.17.6 as build

ENV BIN_FILE /opt/rotator/rotator-app
ENV CODE_DIR /go/src/

WORKDIR ${CODE_DIR}

# Кэшируем слои с модулями
COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . ${CODE_DIR}

# Собираем статический бинарник Go (без зависимостей на Си API),
# иначе он не будет работать в alpine образе.
ARG LDFLAGS
RUN CGO_ENABLED=0 go build \
        -ldflags "$LDFLAGS" \
        -o ${BIN_FILE} ${CODE_DIR}/cmd/rotator/

# На выходе тонкий образ
FROM alpine:3.9

LABEL ORGANIZATION="ITP"
LABEL SERVICE="rotator"
LABEL MAINTAINERS="igortheprogrammer@gmail.com"

ENV BIN_FILE "/opt/rotator/rotator-app"
COPY --from=build ${BIN_FILE} ${BIN_FILE}

ARG CONFIG_FILE_NAME

ENV CONFIG_FILE /etc/rotator/${CONFIG_FILE_NAME}.yaml
COPY ./configs/${CONFIG_FILE_NAME}.yaml ${CONFIG_FILE}

ENV WAIT_VERSION 2.9.0
ADD https://github.com/ufoscout/docker-compose-wait/releases/download/$WAIT_VERSION/wait /wait
RUN chmod +x /wait

CMD /wait && ${BIN_FILE} -config ${CONFIG_FILE}
