this package is used to record logger.
copy the config.template.toml to your workdir, change the default value to what you want,
and pass it the log.Init function to init the logger parameter.

**this package only supports 4 log level,from debug to error.**

## how to use 
firstly, you should call Init() to initialize the log object:

```go
err := log.Init(your config file path)
```

this function may fail if your config file is invalid, so receive the return value and check.

secondly, you can call different level function to record your log:
```go
log.Error("something wrong","pleasse check")
log.Infof("info is %d",100)
```

if you want to change level,you can call SetLevel():
```go
log.SetLevel(log.ErrorLevel)
fmt.Println(log.GetLevel())
fmt.Println(log.GetLevelStr())
```

if you want to flush all buffered logs, call Flush():
```go
log.Flush()
```