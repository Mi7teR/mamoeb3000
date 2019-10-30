FROM golang:latest
ARG token
ARG gaKey
ENV token $token
ENV gaKey $gaKey
ADD . /go/src/github.com/Mi7teR/mamoeb3000
WORKDIR /go/src/github.com/Mi7teR/mamoeb3000
RUN go get .
RUN go build -v
CMD /go/src/github.com/Mi7teR/mamoeb3000/mamoeb3000 -token=${token} -gaKey=${gaKey}
