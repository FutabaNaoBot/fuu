FROM ubuntu:latest
WORKDIR /app
RUN apt-get update -y && apt-get upgrade -y
RUN apt-get install -y tzdata && ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime
RUN apt-get install -y ca-certificates
COPY kohme ./
CMD ["./kohme"]