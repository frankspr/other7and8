package taillog

import (
	"context"
	"fmt"
	"github.com/nxadm/tail"
	"time"
)

type TailTask struct {
	path       string
	instance   *tail.Tail
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewTailTask(path string) (tailObj *TailTask) {
	ctx, cancel := context.WithCancel(context.Background())
	tailObj = &TailTask{
		path:       path,
		ctx:        ctx,
		cancelFunc: cancel,
	}
	tailObj.init()
	return
}

func (t *TailTask) init() {
	config := tail.Config{
		ReOpen: true, // 重新打开
		Follow: true, // 是否跟随
		// Whence 0表示相对于文件的原点，1表示相对于当前偏移量，2表示相对于结束。
		Location:  &tail.SeekInfo{Offset: 0, Whence: 1}, // 从文件的哪个地方开始读
		MustExist: false,                                // 文件不存在不报错
		Poll:      true,
	}
	var err error
	t.instance, err = tail.TailFile(t.path, config)
	if err != nil {
		fmt.Println(err)
	}
	go t.run()
}

func (t *TailTask) run() {
	for {
		select {
		case <-t.ctx.Done():
			fmt.Println("读取任务结束")
			return
		case line := <-t.instance.Lines:
			fmt.Println("打印日志:", line)
		default:
			time.Sleep(time.Second)
		}

	}
}
