FROM alpine AS binary-selector
COPY dist/gostc-open_linux_amd64_v1/server /binary/server-amd64
COPY dist/gostc-open_linux_arm64_v8.0/server /binary/server-arm64

ARG TARGETPLATFORM
RUN mkdir /app && if [ "$TARGETPLATFORM" = "linux/386" ]; then \
        cp /binary/server-386 /app/server; \
    elif [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
        cp /binary/server-amd64 /app/server; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
        cp /binary/server-s390x /app/server; \
    else \
        echo "Unsupported architecture: $TARGETPLATFORM" && exit 1; \
    fi

FROM alpine
COPY --from=binary-selector /app/server /app/server
WORKDIR /app/
ENTRYPOINT ["./server"]