package networks

import "log"

func assertError(err error) {
	if err != nil {
		log.Panicln(err)
	}
}

func assertBool(b bool, s string) {
	if b {
		log.Panicln(s)
	}
}
