
package timer

import (
  "fmt"
  "time"
)

var endTime time.Time
var timerActive bool
var startUpTime time.Time = time.Now()



func timerStart(duration int) {
  startTime := time.Now()
  endTime = startTime.Add(time.Second*time.Duration(duration))
  timerActive = true
  fmt.Println(endTime)
}


func stopTimer() {
  timerActive = false
}


func pollTimeOut(recieve chan<- bool) {
  for {
    if timerActive && endTime.Sub(time.Now()) < 0 {
      receive <- true
      stopTimer()
    }
  }
}


func getTimeStamp() time.Time {
  return time.Now()
}