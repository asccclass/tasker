package SryTask

import(
   "os"
   "log"
   "time"
   "sync"
   "syscall"
   "os/signal"
)

type TaskScheduler struct {
   tasks	map[time.Duration][]Task	// Map of tasks per interval
   quit		chan struct{}
   done		chan struct{}

   // Sync wait group
   wg		sync.WaitGroup

   // Dependencies
   logger	Logger
   timeManager	TimeManager
   timeout	time.Duration
}

func(s *TaskScheduler) AddTask(interval time.Duration, task Task) {
   s.tasks[interval] = append(s.tasks[interval], task)
}

func(s *TaskScheduler) ListenForGracefullyShutdown() {
   sigc := make(chan os.Signal, 1)
   signal.Notify(sigc, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
   go func() {
      <-sigc
      s.logger.Printf("Scheduler shutting down gracefully...")
      s.Shutdown()
   }()
}

func(s *TaskScheduler) Shutdown() {
   s.logger.Printf("Scheduler shutting down...")

   // Stop the ticker
   s.quit <- struct{}{}

   //s.wg.Done()
   // Wait for all running tasks to finish
   // why do we need this?
   // because we want to wait for all tasks to finish
   s.wg.Wait()

   s.logger.Printf("Scheduler shutdown completed")

   // send done signal
   // this will unblock the Start() method and return
   // why send a struct{}{}?
   // because we want to signal that the scheduler is done
   s.done <- struct{}{}
}

func(s *TaskScheduler) executeTask(task Task) {
   defer func() {
      s.wg.Done()                   // Decrement the WaitGroup counter after the task is done
      if r := recover(); r != nil { // Recover from any panics in the task
         s.logger.Printf("ERROR: Task %s - %v", task.ID, r)
      }
   }()

   s.logger.Printf("INFO: Task started id: %s", task.ID) // Log task start
   if err := task.Handler(); err != nil {                // Execute the task; if it returns an error...
      task.ErrorHandler(err)                               // Handle the error
      s.logger.Printf("ERROR: Task %s - %v", task.ID, err) // Log the error
   } else {
      s.logger.Printf("INFO: Task completed successfully id: %s", task.ID) // Log task completion
   }
}

func(s *TaskScheduler) Start() {
   for interval, tasks := range s.tasks { // Iterate over each interval and its associated tasks
      ticker := s.timeManager.NewTicker(interval) // Create a new ticker for the interval
      s.wg.Add(1)
      go func(ticker *time.Ticker, tasks []Task) { // Start a new goroutine for each interval
         defer s.wg.Done()   // Decrement the WaitGroup counter after the goroutine is done
         defer ticker.Stop() // Stop the ticker after the goroutine is done
         for {
            select {
            case <-ticker.C: // When the ticker ticks...
               for _, t := range tasks { // Iterate over each task for the interval
                  s.wg.Add(1)   // Increment the WaitGroup counter before starting a new task
                  taskCopy := t // Avoid data race by copying the task
                  go s.executeTask(taskCopy)
               }
            case <-s.quit:
               return
            }
         }
      }(ticker, tasks)
   }
   s.wg.Wait() // Wait for all tasks to complete before returning
   <-s.done  // Block until done is closed
}

func NewTaskScheduler() (*TaskScheduler){
   logger := log.New(os.Stdout, "SryScheduler: ", log.LstdFlags)
   logger.SetFlags(log.LstdFlags)
   return &TaskScheduler{
      tasks:       make(map[time.Duration][]Task),
      quit:        make(chan struct{}),
      logger:      logger,
      done:        make(chan struct{}),
      timeManager: &RealTimeManager{},
   }
}
