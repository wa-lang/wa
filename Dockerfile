# 版权 @2024 凹语言 作者。保留所有权利。

FROM --platform=linux/amd64 ubuntu:22.04

LABEL org.opencontainers.image.source="https://github.com/wa-lang/wa"

WORKDIR /root

# install git
RUN apt update
RUN apt install -y git

# go run ./builder
COPY ./wa /usr/local/wa
ENV PATH=${PATH}:/usr/local/bin/wa

# docker run --platform linux/amd64 --rm -it wa-lang/wa
CMD ["/bin/bash"]
