FROM scratch

COPY potz-holzoepfel-und-zipfelchape /
COPY etc/passwd /etc/passwd

USER nobody

ENTRYPOINT ["/potz-holzoepfel-und-zipfelchape"]
