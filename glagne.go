// This is simple Dockerfile generator
package main

import (
	"fmt"
	"awesomeProject/alpine"
)


type version struct {
	php          string
	distrib      string
	package_name string
}
// php composer installation
func php_composer_setup() string {
	compose := "EXPECTED_COMPOSER_SIGNATURE=$(wget -q -O - https://composer.github.io/installer.sig) && \\\n"
	compose += "	php -r \"copy('https://getcomposer.org/installer', 'composer-setup.php');\" && \\\n"
	compose += "	php -r \"if (hash_file('SHA384', 'composer-setup.php') === '${EXPECTED_COMPOSER_SIGNATURE}') "
	compose += "{ echo 'Composer.phar Installer verified'; } else { echo 'Composer.phar Installer corrupt'; "
	compose += "unlink('composer-setup.php'); } echo PHP_EOL;\" && \\\n"
	compose += "php composer-setup.php --install-dir=/usr/bin --filename=composer && \\\n"
	compose += "php -r \"unlink('composer-setup.php');\" && \\\n"
	return compose
}
//install memcached
func std_conf_and_make() string {
	make_and_conf := "phpize &&\\\n"
	make_and_conf += "./configure && \\\n"
	make_and_conf += "make && \\\n"
	make_and_conf += "make install &&\\\n"
	return make_and_conf
}
func install_memcached() string {
	memcach := "apk add --virtual .memcached-build-dependencies \\\n"
	memcach += "	libmemcached-dev \\\n"
	memcach += "	cyrus-sasl-dev && \\\n"
	memcach += "apk add --virtual .memcached-runtime-dependencies \\\n"
	memcach += "libmemcached &&\\\n"
	memcach += "git clone -o ${MEMCACHED_TAG} --depth 1 https://github.com/php-memcached-dev/php-memcached.git /tmp/php-memcached && \\\n"
	memcach += "cd /tmp/php-memcached &&\\\n"
	memcach += "phpize &&\\\n"
	memcach += "./configure \\\n"
	memcach += "    --disable-memcached-sasl \\\n"
	memcach += "    --enable-memcached-msgpack \\\n"
	memcach += "    --enable-memcached-json && \\\n"
	memcach += "make && \\\n"
	memcach += "make install && \\\n"
	memcach += "apk del .memcached-build-dependencies && \\\n"
	return memcach
}

func install_msgpack() string {
	msgpack := "git clone -o ${MSGPACK_TAG} --depth 1 https://github.com/msgpack/msgpack-php.git /tmp/msgpack-php && \\\n"
	msgpack += "cd /tmp/msgpack-php && \\\n"
	msgpack += std_conf_and_make()
	return msgpack
}

func install_imagick() string {
	imagick := "apk add --no-cache --virtual .imagick-build-dependencies \\\n"
	imagick += "  autoconf \\\n"
	imagick += "  g++ \\\n"
	imagick += "  gcc \\\n"
	imagick += "  git \\\n"
	imagick += "  imagemagick-dev \\\n"
	imagick += "  libtool \\\n"
	imagick += "  make \\\n"
	imagick += "  tar && \\\n"
	imagick += "apk add --virtual .imagick-runtime-dependencies \\\n"
	imagick += "  imagemagick &&\\\n"
	imagick += "git clone -o ${IMAGICK_TAG} --depth 1 https://github.com/mkoppanen/imagick.git /tmp/imagick &&\\\n"
	imagick += "cd /tmp/imagick && \\\n"
	imagick += std_conf_and_make()
	imagick += "echo \"extension=imagick.so\" > /usr/local/etc/php/conf.d/ext-imagick.ini && \\\n"
	imagick += "apk del .imagick-build-dependencies && \\\n"
	return imagick
}
// Setup modules from code(GIT)
func unstandart_modules_install(module string) (string, string) {
	if module == "memcached" {
		// version memcached
		arg := "ARG MEMCACHED_TAG=v3.0.4"
		return arg, install_memcached()
	}
	if module == "msgpack" {
		arg := "ARG MSGPACK_TAG=msgpack-2.0.2"
		return arg, install_msgpack()
	}
	if module == "imagick" {
		arg := "ARG IMAGICK_TAG = \"3.4.2\""
		return arg, install_imagick()
	}
	return "", ""
}

func main() {
	php_version := make(map[string]version)
	php_version["7.1-alpine"] = version{
		php:          "7.1",
		distrib:      "alpine",
		package_name: "php:7.1-fpm-alpine",
	}
	php_version["7.2-alpine"] = version{
		php:          "7.1",
		distrib:      "alpine",
		package_name: "php:7.2-fpm-alpine",
	}
	php_version["7.1-jessie"] = version{
		php:          "7.1",
		distrib:      "debian",
		package_name: "php:7.1-fpm-jessie",
	}
	php_version["7.2-jessie"] = version{
		php:          "7.2",
		distrib:      "debian",
		package_name: "php:7.2-fpm-jessie",
	}

	maintainer := "\"DockerFile generator by fp <alexwolk01@gmail.com>\" \n"
	composer := true
	var php_modules []string
	php_modules = append(php_modules, "mysqli")
	php_modules = append(php_modules, "memcached")
	php_modules = append(php_modules, "imagick")
	php_modules = append(php_modules, "gd")
	modules_nopecl := []string{"memcached", "imagick", "msgpack"}

	var switcher bool
	ARG := "\n"
	modules_lines := ""
	var arg, str_module string
	docker_modules := []string{"iconv", "pdo", "sqlite", "mysqli", "gd", "exif", "intl", "xsl",
		"json", "soap", "dom", "zip", "opcache", "xml", "mbstring",
		"bz2", "calendar", "ctype", "bcmatch",
	}

	docker_php_ext_install := "docker-php-ext-install "
	for _, module := range (php_modules) {
		switcher = true
		for _, nopecl := range (modules_nopecl) {
			if module == nopecl {
				arg, str_module = unstandart_modules_install(module)
				//generate script
				switcher = false
				modules_lines += str_module
				ARG += arg + "\n"
			}
		}
		if switcher {
			for _, docker_module := range (docker_modules) {
				if module == docker_module {
					docker_php_ext_install += module + " "
				}
			}
		}

	}
	docker_php_ext_install += "&& \\\ndocker-php-source delete && \\\n"
	ARG += "\n"

	ENV := "ENV php_conf /usr/local/etc/php-fpm.conf\n"
	ENV += "ENV fpm_conf /usr/local/etc/php-fpm.d/www.conf\n"
	ENV += "ENV php_vars /usr/local/etc/php/conf.d/docker-vars.ini\n"
	ENV += "ENV LD PRELOAD /usr/lib/preloadable_libconv.so php\n"
	HEAD := "FROM " + php_version["7.1-alpine"].package_name + "\n" + "LABEL maintainer = " + maintainer + "\n"

	Dockerfile := HEAD
	Dockerfile += ENV
	Dockerfile += ARG
	if php_version["7.1-alpine"].distrib == "alpine" {
		Dockerfile += alpine.Soft_install_apk()
	}
	Dockerfile += docker_php_ext_install
	if composer {
		Dockerfile += php_composer_setup()
	}
	Dockerfile += modules_lines
	fmt.Println(Dockerfile)
}
