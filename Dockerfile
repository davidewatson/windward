FROM golang:1.7.6

RUN     git clone https://github.com/davidewatson/careen.git windward && \
        cd windward && \
        go get github.com/samsung-cnct/windward \
        go build && \
        cp windward /
