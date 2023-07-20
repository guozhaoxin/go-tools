package log

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
	"time"
)

func TestInit(t *testing.T) {
	path := "./config.template.toml"
	t.Logf("file is %s",path)
	err := Init(path)
	assert.Nil(t,err,"init should succ")
	if err != nil{
		t.Errorf("err is %s",err)
	}

	path = "./wrong-config.template.toml"
	t.Logf("invalid file is %s",path)
	err = Init(path)
	assert.NotNil(t,err,"init should fail")
	if err != nil{
		t.Logf("err is %s",err)
	}

	path = "./abc/test.file"
	t.Logf("file is %s and not existed",path)
	err = Init(path)
	assert.NotNil(t,err,"init should fail")
	if err != nil{
		t.Logf("err is %s",err)
	}
}

func TestSetLevel(t *testing.T) {
	t.Log("init log with empty str")
	err := Init(" ")
	assert.Nil(t,err,"init shoul succ")
	if err != nil{
		t.Fatalf("unreached case, err is %s",err)
	}

	t.Logf("current level is %d",GetLevel())

	levels := []Level{ErrorLevel,ErrorLevel,WarnLevel,InfoLevel,DebugLevel}
	for _, level := range levels{
		t.Logf("set level to %d",level)
		assert.Nil(t,SetLevel(level),"set should succ")
		assert.Equal(t,level,GetLevel(),"level should equal.")
	}

	var level,oldLevel Level

	level = -3
	oldLevel = GetLevel()
	t.Logf("set level to %d",level)
	assert.NotNil(t,SetLevel(level),"set should fail")
	assert.Equal(t,oldLevel,GetLevel(),"level should not changed.")

	level = 7
	oldLevel = GetLevel()
	t.Logf("set level to %d",level)
	assert.NotNil(t,SetLevel(level),"set should fail")
	assert.Equal(t,oldLevel,GetLevel(),"level should not changed.")
}

func TestGetLevelStr(t *testing.T) {
	_ = Init("")
	levelStr := GetLevelStr()
	assert.Equal(t,"debug",levelStr)
	_ = SetLevel(InfoLevel)
	levelStr = GetLevelStr()
	assert.Equal(t,"info",levelStr)
	_ = SetLevel(ErrorLevel)
	levelStr = GetLevelStr()
	assert.Equal(t,"error",levelStr)
}

func TestLogs(t *testing.T) {
	t.Log("test log functions")
	_ = Init("./config.template.toml")
	writeLogs(t)
	_ = Init("")
	writeLogs(t)
}

func writeLogs(t *testing.T){
	t.Log("create goroutine to run")
	for i := 0; i < 3; i++{
		wg := sync.WaitGroup{}
		for j := 0; j < 3; j++{
			wg.Add(1)
			go func(index int) {
				defer wg.Done()
				Info("goroutine info ",index)
				Infof("goroutine infof %d",index)
				Debug("goroutine debug ",index)
				Debugf("goroutine debugf %d",index)
				Warn("goroutine warn ",index)
				Warnf("goroutine warnf %d",index)
				Error("goroutine error ",index)
				Errorf("goroutine errorf %d",index)
			}(j)
		}
		wg.Wait()
	}
}

func TestLogsWithModeChanged(t *testing.T){
	levels := []Level{DebugLevel,InfoLevel,WarnLevel,ErrorLevel,DebugLevel}

	_ = Init("")
	count := 1
	for _, level := range levels{
		_ = SetLevel(level)
		Debug("debug ",count)
		Debugf("debugf %d",count)
		Info("info ",count)
		Infof("infof %d",count)
		Warn("warn ",count)
		Warnf("warnf %d",count)
		Error("error ",count)
		Errorf("errorf %d",count)
		count++
		time.Sleep(time.Second * 3)
	}

	_ = Init("./config.template.toml")
	count = 1
	for _, level := range levels{
		_ = SetLevel(level)
		Debug("debug ",count)
		Debugf("debugf %d",count)
		Info("info ",count)
		Infof("infof %d",count)
		Warn("warn ",count)
		Warnf("warnf %d",count)
		Error("error ",count)
		Errorf("errorf %d",count)
		count++
		time.Sleep(time.Second * 3)
	}
}

func TestWriteToFile(t *testing.T){
	err := Init("./config.template.toml")
	if err != nil{
		t.Fatalf("unreached case, but failed: %s",err)
	}

	wg :=  sync.WaitGroup{}
	for i:=0; i <10; i++{
		wg.Add(1)
		go func(name string) {
			defer wg.Done()
			for{
				Info("info ",name)
				Infof("info %s",name)
				Debug("debug ",name)
				Debugf("debug %s",name)
				Warn("warn ",name)
				Warnf("warn %s",name)
				Error("error ",name)
				Errorf("error %s",name)
			}
		}(fmt.Sprintf("goroutine %d",i))
	}
	wg.Wait()
}