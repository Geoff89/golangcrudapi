FROM golang:1.20-alpine3.16

# set destination for copy
WORKDIR /app

# add packages necessary for project
RUN apk update && \
    apk add libc-dev && \
    apk add gcc && \
    apk add make

# Download Go modules
COPY go.mod go.sum ./ 
RUN go mod download && go mod verify  


# Install Compile Daemon for go. We'll use it to watch changes in go files
RUN go get github.com/githubnemo/CompileDaemon

# COpy and build the app
COPY . ./
COPY config.json ./
COPY ./entrypoint.sh ./entrypoint.sh

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-mysql-build

# wait-for-it requires bash, which alpine doesnt ship with by default. wait-for is used instead
ADD https://raw.githubusercontent.com/eficode/wait-for/v2.1.0/wait-for /usr/local/bin/wait-for
RUN chmod +rx /usr/local/bin/wait-for ./entrypoint.sh

ENTRYPOINT [ "./entrypoint.sh" ]

CMD ["/docker-mysql-build"]