package commands

//counterfeiter:generate -o ./fakes/logger.go --fake-name Logger . logger

type logger interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
	Println(v ...interface{})
}
