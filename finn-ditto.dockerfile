FROM alpine:latest 

RUN mkdir /app 

COPY toga /app

CMD ["/app/ditto"]