
package timer

import (
  "time"
)




var END_TIME time.Time
var TIMER_ACTIVE bool
var STARTUP_TIME time.Time = time.Now()

var DOOR_OPEN_TIME int = 5




func StartTimer(duration int) {
  startTime := time.Now()
  END_TIME = startTime.Add(time.Second*time.Duration(duration))
  TIMER_ACTIVE = true
}


func stopTimer() {
  TIMER_ACTIVE = false
}


func PollTimeOut(recieve chan <- bool) {
  for {
    if TIMER_ACTIVE && END_TIME.Sub(time.Now()) < 0 {
      recieve <- true
      stopTimer()
    }
  }
}


func getTimeStamp() time.Time {
  return time.Now()
}
