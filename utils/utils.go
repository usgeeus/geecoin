package utils

import "log"

func HandlerErr(err error) {
	if err != nil {
		log.Panic(err)
	}
}
