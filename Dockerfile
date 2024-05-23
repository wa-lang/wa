# 版权 @2024 凹语言 作者。保留所有权利。

FROM --platform=linux/amd64 golang:1.17

LABEL org.opencontainers.image.source="https://github.com/wa-lang/wa"

WORKDIR /root

# go run ./builder
COPY ./_output/wa-docker-linux-amd64 /usr/local/wa
ENV PATH=${PATH}:/usr/local/wa/bin

# docker run --platform linux/amd64 --rm -it wa-lang/wa
CMD ["/bin/bash"]
