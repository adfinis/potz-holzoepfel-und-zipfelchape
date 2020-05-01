FROM scratch

COPY potz-holzoepfel-und-zipfelchape /
COPY etc/passwd /etc/passwd

USER 65534

ENTRYPOINT ["/potz-holzoepfel-und-zipfelchape"]
