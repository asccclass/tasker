# 任務執行器


## 範例

```
package main

import(
   "time"
   "github.com/asccclass/tasker"
)

func main() {
   sch := SryTask.NewTaskScheduler()
   task := SryTask.NewTask(func() error {
      println("Hello world!")
      return nil
   })

   sch.AddTask(time.Second*1, task)
   sch.ListenForGracefullyShutdown()
   sch.Start()
}
```

* 執行結果

```
SryScheduler: 2024/01/30 01:41:47 INFO: Task started id: 4445324a...
Hello world!
SryScheduler: 2024/01/30 01:41:47 INFO: Task completed successfully id: 4445324a...
```
