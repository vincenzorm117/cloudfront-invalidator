FROM amazonlinux

RUN mkdir -p /app

WORKDIR /app

COPY . .

RUN yum -y update
RUN yum -y install go