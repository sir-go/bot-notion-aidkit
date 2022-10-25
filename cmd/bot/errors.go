package main

func eh(err error, msg ...string) {
	if err == nil {
		return
	}
	LOG.Panic("panic err: ", err, " msg ", msg)
}

//noinspection GoUnusedFunction
func ehSkip(err error, msg ...string) {
	if err == nil {
		return
	}
	if len(msg) > 0 {
		LOG.Println(err, msg)
	} else {
		LOG.Println(err)
	}
}
