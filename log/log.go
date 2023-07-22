package log

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
	"strings"
)

type Level zapcore.Level

const (
	DebugLevel = (Level)(zap.DebugLevel)
	InfoLevel = (Level)(zap.InfoLevel)
	WarnLevel = (Level)(zap.WarnLevel)
	ErrorLevel = (Level)(zap.ErrorLevel)
)

var level Level
var logger *zap.SugaredLogger

type logFileConfig struct {
	FileName string `toml:"fileName" comment:"the full path for log file with log name ant ext" validate:"required" err:"fileName must given"`
	MaxSize int `toml:"maxSize" comment:"size threshold to rotate log file, in megabytes" validate:"min=1,max=300" err:"maxSize must in [1,300]"`
	MaxBackups int `toml:"maxBackups" comment:"maximum number of old log files" validate:"min=1,max=100" err:"maxBackups must in [1,100]"`
	MaxAge int `toml:"maxAge" comment:"max days to retain a log file" validate:"min=1,max=365" err:"maxAge must in [1,365]"`
	Compress bool `toml:"compress" comment:"if the rotated log files should be compressed using gzip, default not"`
}


func Init(configFile string) error {
	configFile = strings.TrimSpace(configFile)
	var err error
	if len(configFile) == 0 {
		err = initLogToTerminal()
	}else{
		err = initLogToFile(configFile)
	}
	_ = SetLevel(DebugLevel)
	return err
}

func buildLogger(out zapcore.WriteSyncer){
	core :=zapcore.NewCore(getEncoder(),out, zapcore.DebugLevel)
	l := zap.New(core)
	logger = l.Sugar()
}

func initLogToTerminal() error { // return error to keep consistent with log-file mode.
	fmt.Println("all log write to terminal")
	buildLogger(os.Stdout)
	return nil
}

func initLogToFile(configFilePath string) error {
	fmt.Println("load config from ",configFilePath)
	err,config := loadConfig(configFilePath)
	if err != nil{
		fmt.Printf("fail to load config: %s\n",err)
		return err
	}
	writeSyncer := getLogWriter(config)
	buildLogger(writeSyncer)

	return nil
}

func getEncoder() zapcore.Encoder {
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	return zapcore.NewConsoleEncoder(encoderConfig)
}

func getLogWriter(config *logFileConfig) zapcore.WriteSyncer {

	lumberJackLogger := &lumberjack.Logger{
		Filename:   config.FileName,
		MaxSize:    config.MaxSize,
		MaxBackups: config.MaxBackups,
		MaxAge:     config.MaxAge,
		Compress:   config.Compress,
		LocalTime: true,
	}
	return zapcore.AddSync(lumberJackLogger)
}

func loadConfig(configFile string) (error, *logFileConfig) {
	var err error
	var config logFileConfig
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil{
		fmt.Printf("fail to decode file: %s\n",err)
		return err, nil
	}

	err = checkConfigValid(config)
	if err != nil{
		fmt.Println("config parameters cannot pass check.")
		return err,nil
	}

	return err, &config
}

func SetLevel(newLevel Level) error {
	if newLevel < DebugLevel {
		return fmt.Errorf("level cannot less than %d",DebugLevel)
	}

	if newLevel > ErrorLevel{
		return fmt.Errorf("level cannot more than %d",ErrorLevel)
	}

	level = newLevel
	return nil
}

func GetLevel() Level{
	return level
}

func GetLevelStr() string{
	switch level {
	case DebugLevel:
		return "debug"
	case InfoLevel:
		return "info"
	case WarnLevel:
		return "warn"
	case ErrorLevel:
		return "error"
	}
	return "unknown mode"
}

func Debug(args ...interface{}){
	if level > DebugLevel{
		return
	}
	logger.Debug(args...)
}

func Debugf(template string, args ...interface{}){
	if level > DebugLevel{
		return
	}
	logger.Debugf(template,args...)
}

func Info(args ...interface{}){
	if level > InfoLevel{
		return
	}
	logger.Info(args...)
}

func Infof(template string, args ...interface{}){
	if level > InfoLevel{
		return
	}
	logger.Infof(template,args...)
}

func Warn(args... interface{}){
	if level > WarnLevel{
		return
	}
	logger.Warn(args...)
}

func Warnf(template string, args ...interface{}){
	if level > WarnLevel{
		return
	}
	logger.Warnf(template,args...)
}

func Error(args... interface{}){
	if level > ErrorLevel{
		return
	}
	logger.Error(args...)
}

func Errorf(template string, args ...interface{}){
	if level > ErrorLevel{
		return
	}
	logger.Errorf(template,args...)
}

func Flush() error {
	return logger.Sync()
}

