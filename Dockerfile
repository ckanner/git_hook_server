FROM golang
RUN go get github.com/julienschmidt/httprouter
COPY . /go/src/github.com/ckanner/git_hook_server/
RUN go install github.com/ckanner/git_hook_server
ENTRYPOINT /go/bin/git_hook_server
EXPOSE 8900
