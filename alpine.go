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
	return software
}