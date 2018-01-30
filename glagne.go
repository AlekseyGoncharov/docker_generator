// This is simple Dockerfile generator
package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"strings"
)

type ParsingYaml struct {
	From     string `yaml:"FROM"`
	Composer string `yaml:"composer,omitempty"`
	PhpExt   string `yaml:"php_modules"`
}

type Version struct {
	php          string
	distrib      string
	package_name string
}


func Soft_install_apk() string {
	software := "RUN apk add --no-cache --repository http://dl-3.alpinelinux.org/alpine/edge/testing gnu-libiconv && \\\n"
	software += "echo @testing http://nl.alpinelinux.org/alpine/edge/testing >> /etc/apk/repositories && \\\n"
	software += "echo @main http://mirror.yandex.ru/mirrors/alpine/edge/main >>  /etc/apk/repositories && \\\n"
	software += "echo @community http://mirror.yandex.ru/mirrors/alpine/edge/community >>  /etc/apk/repositories && \\\n"
	software += "echo /etc/apk/respositories && \\\n"
	software += "apk update && \\\n"
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
	software += "libpng-dev \\\n"
	software += "icu-dev \\\n"
	software += "libpq \\\n"
	software += "libxslt-dev \\\n"
	software += "libffi-dev \\\n"
	software += "freetype-dev \\\n"
	software += "sqlite-dev \\\n"
	software += "bzip2-dev \\\n"
	software += "libmemcached-dev \\\n"
	software += "libjpeg-tubo-dev \\\n"
	software += "&& \\\n"
	return software
}
// php composer installation
func Php_composer_setup() string {
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
func Std_conf_and_make() string {
	make_and_conf := "phpize &&\\\n"
	make_and_conf += "./configure && \\\n"
	make_and_conf += "make && \\\n"
	make_and_conf += "make install &&\\\n"
	return make_and_conf
}
func Install_memcached() string {
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

func Install_msgpack() string {
	msgpack := "git clone -o ${MSGPACK_TAG} --depth 1 https://github.com/msgpack/msgpack-php.git /tmp/msgpack-php && \\\n"
	msgpack += "cd /tmp/msgpack-php && \\\n"
	msgpack += Std_conf_and_make()
	return msgpack
}

func Install_imagick() string {
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
	imagick += Std_conf_and_make()
	imagick += "echo \"extension=imagick.so\" > /usr/local/etc/php/conf.d/ext-imagick.ini && \\\n"
	imagick += "apk del .imagick-build-dependencies && \\\n"
	return imagick
}
// Setup modules from code(GIT)
func Unstandart_modules_install(module string) (string, string) {
	if module == "memcached" {
		// version memcached
		arg := "ARG MEMCACHED_TAG=v3.0.4"
		return arg, Install_memcached()
	}
	if module == "msgpack" {
		arg := "ARG MSGPACK_TAG=msgpack-2.0.2"
		return arg, Install_msgpack()
	}
	if module == "imagick" {
		arg := "ARG IMAGICK_TAG = \"3.4.2\""
		return arg, Install_imagick()
	}
	return "", ""
}

func main() {
	conf, err := ioutil.ReadFile("config.yml")
	if err != nil {
		log.Fatalf("error: %v", err)
	}
	confYaml := ParsingYaml{}

	err = yaml.Unmarshal([]byte(conf), &confYaml)
	if err != nil {
		log.Fatalf("error: %v", err)
	}

	php_version := make(map[string]Version)
	php_version["7.1-alpine"] = Version{
		php:          "7.1",
		distrib:      "alpine",
		package_name: "php:7.1-fpm-alpine",
	}
	php_version["7.2-alpine"] = Version{
		php:          "7.1",
		distrib:      "alpine",
		package_name: "php:7.2-fpm-alpine",
	}
	php_version["7.1-jessie"] = Version{
		php:          "7.1",
		distrib:      "debian",
		package_name: "php:7.1-fpm-jessie",
	}
	php_version["7.2-jessie"] = Version{
		php:          "7.2",
		distrib:      "debian",
		package_name: "php:7.2-fpm-jessie",
	}

	maintainer := "\"DockerFile generator by fp <alexwolk01@gmail.com>\" \n"
	//composer := true
	var php_modules []string
	php_modules = strings.Split(confYaml.PhpExt, " ")
	//php_modules = append(php_modules, "mysqli")
	//php_modules = append(php_modules, "memcached")
	//php_modules = append(php_modules, "imagick")
	//php_modules = append(php_modules, "gd")
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
				arg, str_module = Unstandart_modules_install(module)
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

	letsencrypt := "pip install -U pip && \\\npip install -U certbot && \\\nmkdir -p /etc/letsencrypt/webrootauth && \\\n"
	docker_php_ext_install += "&& \\\ndocker-php-source delete && \\\n"
	ARG += "\n"
	clean := "apk del gcc musl-dev linux-headers libffi-dev augeas-dev python-dev make autoconf \n"
	ENV := "ENV php_conf /usr/local/etc/php-fpm.conf\n"
	ENV += "ENV fpm_conf /usr/local/etc/php-fpm.d/www.conf\n"
	ENV += "ENV php_vars /usr/local/etc/php/conf.d/docker-vars.ini\n"
	ENV += "ENV LD PRELOAD /usr/lib/preloadable_libconv.so php\n"
	HEAD := "FROM " + php_version[confYaml.From].package_name + "\n" + "LABEL maintainer = " + maintainer + "\n"

	Dockerfile := HEAD
	Dockerfile += ENV
	Dockerfile += ARG
	if php_version[confYaml.From].distrib == "alpine" {
		Dockerfile += Soft_install_apk()
	}
	Dockerfile += docker_php_ext_install
	if confYaml.Composer == "YES" {
		Dockerfile += Php_composer_setup()
	}
	Dockerfile += modules_lines
	Dockerfile += letsencrypt
	Dockerfile += clean
	fmt.Println(Dockerfile)
}
