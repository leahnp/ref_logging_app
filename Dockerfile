FROM golang

RUN mkdir var
RUN mkdir var/log

ADD log_app.go /
ADD models/* /models/
ADD templates/* /templates/

RUN go build -o /log_app /log_app.go

EXPOSE 8080
ENTRYPOINT ["/log_app"]