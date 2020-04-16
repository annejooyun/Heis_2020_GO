# Module description

The elevator module has responsibility over the local elevator.

Its main functionality is to create and Initialize the elevator object representing the local elevator, as well as updating it when needed.

The elevator structure looks like follows:

```
type Elev struct {
  Id string
  CurrState State
  PrevState State
  CurrDirection elevio.MotorDirection
  PrevDirection elevio.MotorDirection
  Floor int  //current or previous floor
  OrderList [N_FLOORS][N_BUTTONS]int
}
``` 



