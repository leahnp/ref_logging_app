FROM golang

ADD log_app.go /

RUN go build -o /log_app /log_app.go

EXPOSE 8080
ENTRYPOINT ["/log_app"]