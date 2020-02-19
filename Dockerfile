FROM my_go:latest

ENV PORT 8080
ENV HOST 0.0.0.0

ADD . /go/src/github.com/YakovBudnikov/TimetableBot
WORKDIR /go/src/github.com/YakovBudnikov/TimetableBot/cmd
RUN go build .

ENTRYPOINT /go/src/github.com/YakovBudnikov/TimetableBot/cmd/cmd