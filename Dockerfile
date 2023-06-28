FROM golang:alpine

# 为我们的镜像设置必要的环境变量
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64\
    GOPROXY=https://goproxy.cn

# 设置/usr/local/landlord，进入容器就会直接进入到这个目录下,而不是进入到默认根目录下面
WORKDIR /usr/local/landlord

# 复制项目中的 go.mod 和 go.sum文件并下载依赖信息
COPY go.mod .
COPY go.sum .
RUN go mod tidy

# 将代码复制到容器中
COPY . .

# 构建应用
RUN go build -o landlord landlord.go

# 启动应用
CMD ["./landlord"]