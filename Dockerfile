# docker build -t auto-delegator -f DOCKER/Dockerfile .
# docker run -d -p
FROM golang:1.12-buster as builder

ENV APP_PATH /home/adeleg

COPY . ${APP_PATH}

WORKDIR ${APP_PATH}

RUN make create_vendor && \
    make build && \
    cp ./build/auto-delegator /usr/local/bin/auto-delegator

EXPOSE 15500
CMD ["auto-delegator"]
STOPSIGNAL SIGTERM