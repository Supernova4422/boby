FROM golang:1.15


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

RUN apt-get -y update
RUN apt-get -y install gdebi wget
ARG wkhtmltox=https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.buster_amd64.deb
# ARG wkhtmltox=https://github.com/wkhtmltopdf/packaging/releases/download/0.12.6-1/wkhtmltox_0.12.6-1.raspberrypi.buster_armhf.deb
RUN wget -O wkhtmltox.deb ${wkhtmltox}
RUN gdebi --n wkhtmltox.deb


CMD $binary_filepath "$config_path"