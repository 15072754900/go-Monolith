FROM golang:alpine as builder
# 设置工作目录
WORKDIR /gvb
# 将当前目录内容拷到工作目录 (相对路径)
COPY . .
# 配置 golang 环境
RUN go mod tidy \
    && go build -o server .

FROM alpine:latest
ENV WORK_PATH /gvb
WORKDIR ${WORK_PATH}
COPY --from=0 ${WORK_PATH}/server .
COPY --from=0 ${WORK_PATH}/config/config.toml .
COPY --from=0 ${WORK_PATH}/assets/ip2region.xdb ./assets/ip2region.xdb
COPY --from=0 ${WORK_PATH}/assets/wait-for .

RUN chmod a+x ./wait-for

# 后台接口
EXPOSE 8081

# 在 docker-compose 中使用 wait-for 依赖 mysql 再启动
# ENTRYPOINT ./server -c config.docker.toml
# CMD sleep 5 && ./server -c config.docker.toml