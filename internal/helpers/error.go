package helpers

import (
	"log"
)

func CheckErrBool(ok bool) {
	if !ok {
		log.Panic(ok)
	}
}
