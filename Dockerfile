FROM php:7.1-fpm-alpine
LABEL maintainer = "DockerFile generator by fp <alexwolk01@gmail.com>" 

ENV php_conf /usr/local/etc/php-fpm.conf
ENV fpm_conf /usr/local/etc/php-fpm.d/www.conf
ENV php_vars /usr/local/etc/php/conf.d/docker-vars.ini
ENV LD PRELOAD /usr/lib/preloadable_libconv.so php

ARG MSGPACK_TAG=msgpack-2.0.2
ARG IMAGICK_TAG = "3.4.2"
ARG MEMCACHED_TAG=v3.0.4

RUN apk add --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing gnu-libiconv && \
echo @testing http://nl.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories && \
echo @main http://mirror.yandex.ru/mirrors/alpine/edge/main >>  /etc/apk/repositories && \
echo @community http://mirror.yandex.ru/mirrors/alpine/edge/community >>  /etc/apk/repositories && \
echo /etc/apk/repositories && \
apk update && \
apk add --no-cache bash \
wget \
supervisor \
curl \
libcurl \
git \
python \
python-dev \
py-pip \
augeas-dev \
openssl-dev \
ca-certificates \
dialog \
autoconf \
make \
gcc \
musl-dev \
linux-headers \
libmcrypt-dev \
libpng-dev \
icu-dev \
libpq \
libxslt-dev \
libffi-dev \
freetype-dev \
sqlite-dev \
bzip2-dev \
libmemcached-dev \
libjpeg-turbo-dev \
&& \
docker-php-ext-configure gd \
--with-gd \
--with-freetype-dir=/usr/include/ \
--with-png-dir=/usr/include/ \
--with-jpeg-dir=/usr/include/ && \
docker-php-ext-install iconv pdo_mysql pdo_sqlite mysqli gd exif intl xsl json soap dom zip opcache xml mbstring bz2 calendar ctype && \
docker-php-source delete && \
git clone -o ${MSGPACK_TAG} --depth 1 https://github.com/msgpack/msgpack-php.git /tmp/msgpack-php && \
cd /tmp/msgpack-php && \
phpize &&\
./configure && \
make && \
make install &&\
apk add --no-cache --virtual .imagick-build-dependencies \
  autoconf \
  g++ \
  gcc \
  git \
  imagemagick-dev \
  libtool \
  make \
  tar && \
apk add --virtual .imagick-runtime-dependencies \
  imagemagick &&\
git clone -o ${IMAGICK_TAG} --depth 1 https://github.com/mkoppanen/imagick.git /tmp/imagick &&\
cd /tmp/imagick && \
phpize &&\
./configure && \
make && \
make install &&\
echo "extension=imagick.so" > /usr/local/etc/php/conf.d/ext-imagick.ini && \
apk del .imagick-build-dependencies && \
apk add --virtual .memcached-build-dependencies \
	libmemcached-dev \
	cyrus-sasl-dev && \
apk add --virtual .memcached-runtime-dependencies \
libmemcached &&\
git clone -o ${MEMCACHED_TAG} --depth 1 https://github.com/php-memcached-dev/php-memcached.git /tmp/php-memcached && \
cd /tmp/php-memcached &&\
phpize &&\
./configure \
    --disable-memcached-sasl \
    --enable-memcached-msgpack \
    --enable-memcached-json && \
make && \
make install && \
apk del .memcached-build-dependencies && \
pip install -U pip && \
pip install -U certbot && \
mkdir -p /etc/letsencrypt/webrootauth && \
apk del gcc musl-dev linux-headers libffi-dev augeas-dev python-dev make autoconf 
