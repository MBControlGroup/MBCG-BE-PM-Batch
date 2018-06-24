FROM golang:1.10.1
RUN mkdir -p /go/src/mbcs_pm_batch
ADD . /go/src/mbcs_pm_batch
WORKDIR /go/src/mbcs_pm_batch
CMD ["go", "run", "main.go"]
