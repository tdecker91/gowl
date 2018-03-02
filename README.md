# GOWL Golang Websocket Logger
This is a simple service that will take in log messages from a websocket and log them.

## Running
clone the repository
`go get github.com/tdecker91/gowl`

run the binary
`./{GOPATH}/bin/gowl`

GOWL will listen on port 8030 by default for websocket connections. It listens for json message formatted like the example below. Once it receives a message it will log it. 

```json
 {
   "message": "the log message",
   "level": "the log level"
 }
```

## Log levels
Supported levels are debug, info, warn, error. Error will be printed to os.Stderr and everything else will be printed to os.Stdout

### Optional flags
GOWL defaults can be overridden with command line flags. Here is an example of running GOWL on a different port with a different message format.
`gowl -port=9990 -format="(%L): [%M]"`

| Flag  | Default | Description |
| ----- | ------- | ----------- |
| format | [%L] %D: %M | Message format. %L = log level. %D = Date. %M = Message. Only used for the text formatter |
| formatter | text | gowl uses logrus as the underlying logging system. use formatter=json to output messages as json |
| level | info | used to set the log level that will be printed. level supported are debug, info, warn, error |
| port | 8030 | the port gowl will listen on for websocket connections |
| route | / | the route the websocket server will listen for. For example if route is set to "/messages" then clients should connect to "ws://{address}:{port}/messages" |
| timeFormat | "2006-01-02 15:04:05" | Format string for the date. Only used for the text formatter. Refer to golang's time package for information on how to create the format string https://golang.org/pkg/time/#Time.Format |

## Usage
Use `gowl -h` to see usage options

```
Usage of gowl:
  -format string
    	text formatter only. format for the output message. %L = log level. %D = Date. %M = Message (default "[%L] %D: %M")
  -formatter string
    	format to output log messages as. [text, json] (default "text")
  -level string
    	logging level. [debug, info, warn, error] (default "info")
  -port int
    	port to listen on (default 8030)
  -route string
    	route to listen for socket connections (default "/")
  -timeFormat string
    	text formatter only. format string to output time logs. Format as defined here https://golang.org/pkg/time/#Time.Format (default "2006-01-02 15:04:05")
```