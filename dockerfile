FROM golang:1.23-alpine


ARG binary_filename=bot
ARG project_path=/src/${binary_filename}

COPY . ${project_path}
WORKDIR ${project_path}

RUN go build -o ./font_download font_downloader/font_downloader.go
RUN ./font_download font_downloader/config.json

RUN go get -d -v ./...
RUN go install -v ./...

ARG binary_dir=${project_path}/src/main
ENV binary_filepath ${binary_dir}/${binary_filename}

RUN go build -o ${binary_filepath} ${binary_dir}

RUN apk update
RUN apk add inkscape


CMD $binary_filepath "${CONFIG_PATH}"