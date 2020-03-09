package control

import "../elevio"

const N_FLOORS = 4
const N_BUTTONS = 3 // Number of button types (cab,hall up, hall down)
const Elev_ID = 1

type State int
const (
  idle   State = 0
  moving       = 1
  doorOpen     = 2
)

type Elev struct {
  Id int
  State State
  Direction elevio.MotorDirection
  PrevFloor int  //prevFloor might be current floor
  OrderList [N_FLOORS][N_BUTTONS]int
}


func InitializeElevator() Elev {
  elevio.Init("localhost:15657", N_FLOORS)
  elevio.SetMotorDirection(elevio.MD_Down)
  floor := -1
  for floor == -1 {
    floor = elevio.GetFloor()
  }
  elevio.SetMotorDirection(elevio.MD_Stop)
  elevator := Elev{Elev_ID, idle,elevio.MD_Stop,floor,[N_FLOORS][N_BUTTONS]int{}}
  return elevator
}
