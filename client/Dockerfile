FROM alpine AS binary-selector
COPY dist/gostc-sub_linux_386_sse2/gostc /binary/gostc-386
COPY dist/gostc-sub_linux_amd64_v1/gostc /binary/gostc-amd64
COPY dist/gostc-sub_linux_arm_7/gostc /binary/gostc-armv7
COPY dist/gostc-sub_linux_arm64_v8.0/gostc /binary/gostc-arm64
COPY dist/gostc-sub_linux_mips64le_hardfloat/gostc /binary/gostc-mips64le
COPY dist/gostc-sub_linux_s390x/gostc /binary/gostc-s390x

ARG TARGETPLATFORM
RUN mkdir /app && if [ "$TARGETPLATFORM" = "linux/386" ]; then \
        cp /binary/gostc-386 /app/gostc; \
    elif [ "$TARGETPLATFORM" = "linux/amd64" ]; then \
        cp /binary/gostc-amd64 /app/gostc; \
    elif [ "$TARGETPLATFORM" = "linux/arm/v7" ]; then \
        cp /binary/gostc-armv7 /app/gostc; \
    elif [ "$TARGETPLATFORM" = "linux/arm64" ]; then \
        cp /binary/gostc-arm64 /app/gostc; \
    elif [ "$TARGETPLATFORM" = "linux/mips64le" ]; then \
        cp /binary/gostc-mips64le /app/gostc; \
    elif [ "$TARGETPLATFORM" = "linux/s390x" ]; then \
        cp /binary/gostc-s390x /app/gostc; \
    else \
        echo "Unsupported architecture: $TARGETPLATFORM" && exit 1; \
    fi

FROM alpine
COPY --from=binary-selector /app/gostc /app/gostc
RUN apk add --no-cache iptables
WORKDIR /app/
ENTRYPOINT ["./gostc"]