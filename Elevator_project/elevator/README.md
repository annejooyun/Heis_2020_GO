# Elevator

## Module description

The elevator module has responsibility over the local elevator.

Its main functionality is to create and initialize the elevator object representing the local elevator, as well as updating it when needed.

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

Most of these struct variables should be self-explanatory, as well as the two main functions in the elevator-module: `InitializeElevator()` and `PollInternalElevatorStatus()`. However, it is important to note that the `OrderList` only consists of the local orders for the elvator. 
