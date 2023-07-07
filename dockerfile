FROM alpine:3.16.2
COPY build/timer-schedule /usr/local
COPY etc/app.yaml /usr/local
WORKDIR /usr/local
CMD /usr/local/timer-schedule