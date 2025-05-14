FROM ubuntu:latest
LABEL authors="polzovatel"

ENTRYPOINT ["top", "-b"]