package log

import (
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"strings"
)

type Level zapcore.Level

const (
	// only supports 4 level
	DebugLevel = (Level)(zap.DebugLevel)
	InfoLevel = (Level)(zap.InfoLevel)
	WarnLevel = (Level)(zap.WarnLevel)
	ErrorLevel = (Level)(zap.ErrorLevel)
)

var level Level
var logger *zap.SugaredLogger

type logFileConfig struct {
	FileName string `toml:"fileName" comment:"the full path for log file with log name ant ext"`
	MaxSize int `toml:"maxSize" comment:"size threshold to rotate log file, in megabytes"`
	MaxBackups int `toml:"maxBackups" comment:"maximum number of old log files"`
	MaxAge int `toml:"maxAge" comment:"max days to retain a log file"`
	Compress bool `toml:"compress" comment:"if the rotated log files should be compressed using gzip, default not"`
}


func Init(configFile string) error {
	configFile = strings.TrimSpace(configFile)
	if len(configFile) == 0 {
		return initLogWithTerminal()
	}
	return initLogToFile(configFile)

}

func initLogWithTerminal() error { // return error to keep consistent with log-file mode.
	fmt.Println("all log write to terminal")
	p, _ := zap.NewProduction()
	logger = p.Sugar()
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
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, writeSyncer, zapcore.DebugLevel)

	l := zap.New(core)
	logger = l.Sugar()
	logger.WithOptions()

	return nil
}

func loadConfig(configFile string) (error, *logFileConfig) {
	var err error
	var config logFileConfig
	_, err = toml.DecodeFile(configFile, &config)
	if err != nil{
		fmt.Printf("fail to decode file: %s\n",err)
		return err, nil
	}

	err = checkConfigValid(&config)
	if err != nil{
		fmt.Println("config parameters cannot pass check.")
		return err,nil
	}

	return err, &config
}

func checkConfigValid(config *logFileConfig)(error){
	s := ""

	if len(config.FileName) == 0{
		s +=  fmt.Sprintf("log file path cannot be empty")
	}
	if config.MaxAge < 1{
		s += fmt.Sprintf("invalid maxAge=%d,cannot less than 1.\n",config.MaxAge)
	}
	if config.MaxSize < 1{
		s += fmt.Sprintf("invalid maxSize=%d,cannot less than 1.\n",config.MaxSize)
	}
	if config.MaxBackups < 1{
		s += fmt.Sprintf("invalid maxBackups=%d,cannot less than 1.\n",config.MaxBackups)
	}

	if len(s) != 0{
		return fmt.Errorf("%s",s)
	}

	return nil
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
	}
	return zapcore.AddSync(lumberJackLogger)
	// todo gzx 需要验证一下，日志是否会追加，还是每次运行时，都新打开一个。
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
	logger.Warn(args...)
}

func Errorf(template string, args ...interface{}){
	if level > ErrorLevel{
		return
	}
	logger.Warnf(template,args...)
}
