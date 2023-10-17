package error

import "fmt"

func Error(err error) {
	fmt.Println("An Error Has Occurred!\n\n" + err.Error())
}