package SryTask

import(
   "time"
   "github.com/google/uuid"
)

type Task struct {
   //ID is a unique identifier for the task, If ID is not set, a random ID is generated
   ID string
   //Handler is a function that executes the task
   //If Handler is not set, the task is ignored
   //If Handler is set, the task is executed
   //If Handler is set to nil, the task is ignored
   Handler func() error

   //ErrorHandler is a function that handles errors
   //If ErrorHandler is not set, the error is logged
   //If ErrorHandler is set, the error is passed to the function
   //If ErrorHandler is set to nil, the error is ignored
   ErrorHandler func(error)

   //Timeout is the time after which the task is considered to have timed out
   Timeout time.Duration
}

type Logger interface {
   Printf(format string, v ...interface{})
}

type TimeManager interface {
   NewTicker(duration time.Duration) *time.Ticker
}

func NewTask(handler func() error)(Task) {
   return Task{
      ID:      uuid.New().String(),
      Handler: handler,
   }
}
