FROM centos:7
WORKDIR /app
ADD . /app
ENV RUN_MODE pro
ENV GIN_MODE release
EXPOSE 8090
CMD ["/app/dx-server"]