# 1 шаг
FROM golang:1.23-alpine AS build_stage
COPY . /go/src/my_super_app
WORKDIR /go/src/my_super_app
RUN go mod download

# 2 шаг
FROM alpine AS run_stage
WORKDIR /app_binary
COPY --from=build_stage /go/bin/my_super_app /app_binary/
RUN chmod +x ./my_super_app
EXPOSE 8080/tcp
ENTRYPOINT ./my_super_app

EXPOSE 8080/tcp
CMD [ "my_super_app" ]

# docker build --file=1_Dockerfile.Multistage --tag=my_super_app-multistage:latest .
# docker run --rm -p 8080:8080 my_super_app-multistage