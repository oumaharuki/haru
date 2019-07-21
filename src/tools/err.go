package tools

import "log"

func CheckErr(err error) {
	if err != nil {
		//log.Println(err.(*errors.Error).ErrorStack())
		log.Panic(err)
	}
}
