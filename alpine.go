// for alpine version
package alpine

import

func soft()string {
	software := "apk update && \\\n"
	software += "apk add --no-cache bash \\\n"
	software += "wget \\\n"
	software += "supervisor \\\n"
	software += "curl \\\n"
	software += "libcurl \\\n"
	software += "git \\\n"
	software += "python \\\n"
	software += "python-dev \\\n"
	software += "py-pip \\\n"
	software += "augeas-dev \\\n"
	software += "openssl-dev \\\n"
	software += "ca-certificates \\\n"
	software += "dialog \\\n"
	software += "autoconf \\\n"
	software += "make \\\n"
	software += "gcc \\\n"
	software += "musl-dev \\\n"
	software += "linux-headers \\\n"
	software += "libmcrypt-dev \\\n"
	software += "libpng-dev \\\"
	software += "icu-dev \\\n"
	software += "libpq \\\n"
	software += "libxslt-dev \\\n"
	software += "libffi-dev \\\n"
	software += "freetype-dev \\\n"
	software += "sqlite-dev \\\n"
	software += "bzip2-dev \\\n"
	software += "libmemcached-dev \\\n"
	software += "libjpeg-tubo-dev "
	software += "&& \\\n"
	return software
}