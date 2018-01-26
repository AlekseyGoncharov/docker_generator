// This is simple Dockerfile generator
package main

type version struct {
	php          string
	distrib      string
	package_name string
}

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

func unstandart_modules_install(module string) (string, string) {
	if module == "memcached" {
		arg := "ARG MEMCACHED_TAG=v3.0.4"
		return arg, install_memcached()
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

	maintainer := "\"DockerFile generator by fp <alexwolk01@gmail.com>\""
	composer := true
	var php_modules []string
	php_modules = append(php_modules, "mysqli")
	php_modules = append(php_modules, "memcached")
	modules_nopecl := []string{"memcached", "imagick"}

	var switcher bool
	ARG := "\n"
	modules_lines := ""
	var arg, str_module string
	for _, module := range (php_modules) {
		switcher = true
		for _, nopecl := range (modules_nopecl) {
			if module == nopecl {
				arg, str_module = unstandart_modules_install(module)
				//generate script
				switcher = false
				modules_lines += str_module
			}
		}
		if switcher {
			//add pecl string
		}
		ARG += arg

	}

	var ENV []string
	ENV = append(ENV, "php_conf /usr/local/etc/php-fpm.conf\n")
	ENV = append(ENV, "fpm_conf /usr/local/etc/php-fpm.d/www.conf\n")
	ENV = append(ENV, "php_vars /usr/local/etc/php/conf.d/docker-vars.ini\n")
	ENV = append(ENV, "LD PRELOAD /usr/lib/preloadable_libconv.so php\n")
	HEAD := "FROM " + php_version["7.1-alpine"].package_name + "\n" + "LABEL maintainer = " + maintainer + "\n"
	Dockerfile := HEAD
	if composer {
		Dockerfile += php_composer_setup()
	}

}
