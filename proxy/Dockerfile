FROM alpine AS binary-selector
COPY dist/gostc-open_linux_386_sse2/gostc-proxy /binary/gostc-proxy-386
COPY dist/gostc-open_linux_amd64_v1/gostc-proxy /binary/gostc-proxy-amd64
COPY dist/gostc-open_linux_arm_7/gostc-proxy /binary/gostc-proxy-armv7
COPY dist/gostc-open_linux_arm64_v8.0/gostc-proxy /binary/gostc-proxy-arm64
COPY dist/gostc-open_linux_mips64le_hardfloat/gostc-proxy /binary/gostc-proxy-mips64le
COPY dist/gostc-open_linux_s390x/gostc-proxy /binary/gostc-proxy-s390x

ARG TARGETPLATFORM
RUN mkdir /app && if [ "$TARGETPLATFORM" = "linux/386" ]; then \
        cp /binary/gostc-proxy-386 /app/gostc-proxy; \
    elif [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
        cp /binary/gostc-proxy-amd64 /app/gostc-proxy; \
    elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then \
        cp /binary/gostc-proxy-armv7 /app/gostc-proxy; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
        cp /binary/gostc-proxy-arm64 /app/gostc-proxy; \
    elif [ "$TARGETPLATFORM" = "linux/mips64le" ]; then \
        cp /binary/gostc-proxy-mips64le /app/gostc-proxy; \
    elif [ "$TARGETPLATFORM" = "linux/s390x" ]; then \
        cp /binary/gostc-proxy-s390x /app/gostc-proxy; \
    else \
        echo "Unsupported architecture: $TARGETPLATFORM" && exit 1; \
    fi

FROM alpine
COPY --from=binary-selector /app/gostc-proxy /app/gostc-proxy
WORKDIR /app/
ENTRYPOINT ["./gostc-proxy"]