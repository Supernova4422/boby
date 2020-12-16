FROM golang:1.15

ARG project_path=/go/src/github.com/BKrajancic/FLD-Bot

COPY . ${project_path}
WORKDIR ${project_path}

RUN go get -d -v ./...
RUN go install -v ./...

ARG binary_dir=${project_path}/src/main
ARG binary_filename=fld-bot
ARG binary_filepath=${binary_dir}/${binary_filename}

RUN go build -o ${binary_filepath} ${binary_dir}

ARG test_dir=${project_path}/src/test
WORKDIR ${test_dir}
RUN go test
RUN binary_filename=${binary_filename}

ENV PATH=${binary_dir}:$PATH
WORKDIR ${binary_dir}
CMD [${binary_filename}]
