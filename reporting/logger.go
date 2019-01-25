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

// idk how this will be work on load, but, think this is good practise coz one entry point for all logs
func (l *Logger) Start() {
	log.Printf("%s__%s\n", l.AppName, "Logger: Start logging")
	go func() {
		for {
			select {
			case val, ok := <-l.LogChan:
				if !ok {
					break
				}
				//add logic here if wana smth more, like sentry
				log.Printf("%s__%s\n", l.AppName, val)
			}
		}
	}()
}
