package utils

import (
	"fmt"
	"strings"
)

func FormatGrpcValidationError(err string) string {
	// in golang i have this three version of string
	// error: validation error:
	//  - service.url: Url is required [url]
	// error: validation error:
	//  - service.path: Path is required [path]
	//  - service.url: Url is required [url]
	// error: validation error:
	//  - service.name: Name is required [name]
	//  - service.path: Path is required [path]
	//  - service.url: Url is required [url]

	// i want to format it as json like this

	// {
	// error emsssage : validation error
	// fields :{
	// service.name : Name is required ... and so on

	// }
	// }

	// split the error by \n

	splitErr := strings.Split(err, "-")
	fmt.Sprintln(splitErr)

	return err
}
