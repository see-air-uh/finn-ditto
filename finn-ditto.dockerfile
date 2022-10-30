FROM alpine:latest 

RUN mkdir /app 

COPY ditto /app

CMD ["/app/ditto"]