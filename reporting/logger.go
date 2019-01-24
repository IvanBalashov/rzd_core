package reporting

import "log"

type Logger struct {
	LogChan chan string
	AppName string
}

func NewLogger(log chan string, name string) Logger {
	return Logger{
		LogChan: log,
		AppName: name,
	}
}

func (l *Logger) Start() {
	//log.SetFlags(log.LstdFlags)
	log.Printf("%s_%s\n", l.AppName, "Logger: Start logging")
	go func() {
		for {
			select {
			case val, ok := <-l.LogChan:
				if !ok {
					break
				}
				//add logic here if wana smth more, like sentry
				log.Printf("%s_%s\n", l.AppName, val)
			}
		}
	}()
}

func (l *Logger) Write(logline string) {
	log.Printf("%s_%s\n", l.AppName, logline)
}
