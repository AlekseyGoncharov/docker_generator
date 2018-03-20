FROM php:7.1-fpm
LABEL maintainer = "DockerFile generator by fp <alexwolk01@gmail.com>" 

RUN apt-get update && apt-get install -y \
libsqlite3-dev \
libicu-dev \
libxslt-dev \
libbz2-dev \
libmemcached-dev \
zlib1g-dev \
&&docker-php-ext-install iconv pdo_mysql pdo_sqlite mysqli gd exif intl xsl json soap dom zip opcache xml mbstring bz2 calendar ctype 