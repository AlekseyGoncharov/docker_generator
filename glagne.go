// This is simple Dockerfile generator
package main
import (
	//"os"
	//"gopkg.in/yaml.v2"
	"fmt"

)
type version struct {
	php_version string
	distrib string
	package_name string
	}

func php_composer()string {
	compose := "EXPECTED_COMPOSER_SIGNATURE=$(wget -q -O - https://composer.github.io/installer.sig) && \\\n"
	compose += "	php -r \"copy('https://getcomposer.org/installer', 'composer-setup.php');\" && \\\n"
	compose += "	php -r \"if (hash_file('SHA384', 'composer-setup.php') === '${EXPECTED_COMPOSER_SIGNATURE}') "
	compose += "{ echo 'Composer.phar Installer verified'; } else { echo 'Composer.phar Installer corrupt'; "
	compose += "unlink('composer-setup.php'); } echo PHP_EOL;\" && \\\n"
	compose += "php composer-setup.php --install-dir=/usr/bin --filename=composer && \\\n"
	compose += "php -r \"unlink('composer-setup.php');\" && \\\n"
	return compose
}

func main() {
	php_version := make(map[string]version)
	php_version["7.1-alpine"] = version{
		php_version: "7.1",
		distrib: "alpine",
		package_name: "php:7.1-fpm-alpine",
	}
	php_version["7.2-alpine"] = version{
		php_version: "7.1",
		distrib: "alpine",
		package_name: "php:7.2-fpm-alpine",
	}
	php_version["7.1-jessie"] = version{
		php_version: "7.1",
		distrib: "debian",
		package_name: "php:7.1-fpm-jessie",
	}
	php_version["7.2-jessie"] = version{
		php_version: "7.2",
		distrib: "debian",
		package_name: "php:7.2-fpm-jessie",
	}

	maintainer := "\"DockerFile generator by fp <alexwolk01@gmail.com>\""
	composer := true
	var php_modules []string
	php_modules = append(php_modules, "mysqli")
	php_modules = append(php_modules,"memcached")
	modules_nopecl := []string{"memcached", "imagick"}

	var switcher bool
	for _, module := range(php_modules) {
		switcher = true
		for _,nopecl := range(modules_nopecl) {
			if module == nopecl {
				//generate script
				switcher = false
			}
		}
		if switcher {
			//add pecl string
		}
	}

	var ENV []string
	ENV = append(ENV,"php_conf /usr/local/etc/php-fpm.conf\n")
	ENV = append(ENV, "fpm_conf /usr/local/etc/php-fpm.d/www.conf\n")
	ENV = append(ENV, "php_vars /usr/local/etc/php/conf.d/docker-vars.ini\n")
	ENV = append(ENV, "LD PELOAD /us/lib/peloadable_libconv,so php\n")
	HEAD := "FROM " + php_version["7.1-alpine"].package_name + "\n" + "LABEL maintainer = " + maintainer + "\n"
	if composer {
		fmt.Println(HEAD, "10")
	} else {
		fmt.Println(HEAD, "2")
	}

}
