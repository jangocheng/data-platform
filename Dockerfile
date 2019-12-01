FROM alpine
ADD platform-srv /platform-srv
ENTRYPOINT [ "/platform-srv" ]
