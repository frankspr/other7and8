package taillog

import (
	"fmt"
	"time"
)

type logEntry struct {
	Path string
}
type logTaskMgr struct {
	logEntry    []*logEntry
	taskMap     map[string]*TailTask
	newConfChan chan []*logEntry
}

func Init(logEntryConf []*logEntry) {
	taskMgr := &logTaskMgr{
		logEntry:    logEntryConf,
		taskMap:     make(map[string]*TailTask, 16),
		newConfChan: make(chan []*logEntry),
	}
	for _, logEntry := range logEntryConf {
		tailObj := NewTailTask(logEntry.Path)
		mk := fmt.Sprintf("%s", logEntry.Path)
		taskMgr.taskMap[mk] = tailObj
	}

}

func (t *logTaskMgr) run() {
	for {
		select {
		case newlog := <-t.newConfChan:
			for _, c := range newlog {
				mk := fmt.Sprintf("%s", c.Path)
				_, ok := t.taskMap[mk]
				if ok {
					continue
				} else {

					tailObj := NewTailTask(c.Path)
					t.taskMap[mk] = tailObj
				}
			}
			//
			for _, c1 := range t.logEntry {
				isDel := true
				for _, c2 := range newlog {
					if c1 == c2 {
						isDel = false
						continue
					}
				}
				if isDel {
					mk := fmt.Sprintf("%s", c1.Path)
					t.taskMap[mk].cancelFunc()
				}
			}
		default:
			time.Sleep(time.Second)
		}

	}
}
