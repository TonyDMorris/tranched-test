package app

import (
	"log"
)

var DefaultLogger = log.New(log.Writer(), "App", log.LstdFlags)
