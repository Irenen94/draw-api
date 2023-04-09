# 打包依赖阶段使用golang作为基础镜像
FROM golang:1.16.12 as builder

# 启用go module
ENV GO111MODULE=on \
    GOPROXY=https://goproxy.cn,direct

WORKDIR /app

COPY . .

# 指定OS等，并go build
RUN GOOS=linux GOARCH=amd64 go build  -o web-frontend-demo  ./cmd/api/api.go

# 由于我不止依赖二进制文件，还依赖conf文件夹下的配置文件还有etc文件夹下的一些文件
# 所以我将这些文件放到了bin文件夹
RUN cp web-frontend-demo  bin && cp -r config/conf bin

# 运行阶段指定alpine作为基础镜像
FROM alpine

WORKDIR /app

# 将上一个阶段publish文件夹下的所有文件复制进来
COPY --from=builder /app/bin .

RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2

# 指定运行时环境变量
ENV ENV=dev
# 指定时区
RUN echo -e  "http://mirrors.aliyun.com/alpine/v3.4/main\nhttp://mirrors.aliyun.com/alpine/v3.4/community" >  /etc/apk/repositories \
&& apk update && apk add tzdata \
&& cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
&& echo "Shanghai/Asia" > /etc/timezone \
&& apk del tzdata


EXPOSE 2000

EXPOSE 9999

EXPOSE 8280

ENTRYPOINT ["./web-frontend-demo"]