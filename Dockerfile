FROM arm32v7/ruby:alpine
ADD . /
RUN apk add --update \
    build-base \
  && gem install epoll \
  && rm -rf /var/cache/apk/*
CMD ruby /gpio.rb
