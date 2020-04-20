package elevator

import (
  "../elevio"

  "fmt"
)


const N_FLOORS = 4
const N_BUTTONS = 3 //Cab, Hall up, Hall down)
//const N_ELEVS = 3

var LOCAL_ELEV_ID string


type State int
const (
  Idle   State = 0
  Moving       = 1
  DoorOpen     = 2
)


//Elevator struct
type Elev struct {
  Id string
  CurrState State
  PrevState State
  CurrDirection elevio.MotorDirection
  PrevDirection elevio.MotorDirection
  Floor int  //current or previous floor
  OrderList [N_FLOORS][N_BUTTONS]int
}




//Moves the elevator to a known floor. Makes elevator object. Turns off lights.
func InitializeElevator(ConnectionPort string) Elev {
  connectionMessage := fmt.Sprintf("localhost:%s",ConnectionPort)
  elevio.Init(connectionMessage, N_FLOORS)

  //Move elevator to nearest floor below current point.
  elevio.SetMotorDirection(elevio.MD_Down)
  floor := -1
  for floor == -1 {
    floor = elevio.GetFloor()
  }
  elevio.SetMotorDirection(elevio.MD_Stop)

  //Turn off all lights
  initializeLights()

  //Create elevator object. ID is set to the connection port.
  elevator := Elev{ConnectionPort,Idle,Idle,elevio.MD_Stop,elevio.MD_Stop,floor,[N_FLOORS][N_BUTTONS]int{}}
  LOCAL_ELEV_ID = ConnectionPort

  return elevator
}

//Turns off all lights when initializing the elevator
func initializeLights() {
  for floor:= 0; floor < N_FLOORS; floor++ {
    elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
    elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
    elevio.SetButtonLamp(elevio.BT_Cab, floor, false)
  }
  elevio.SetDoorOpenLamp(false)
  elevio.SetStopLamp(false)
}


func UpdateElevatorDirectionsAndStates(elev *Elev, direction elevio.MotorDirection, state State) {
  updatePrevDirection(elev)
  updateCurrDirection(elev, direction)
  updatePrevState(elev)
  updateCurrState(elev, state)
  fmt.Printf("PrevDir = %d\n", elev.PrevDirection)
  fmt.Printf("CurrDir = %d\n", elev.CurrDirection)
  fmt.Printf("PrevState = %d\n", elev.PrevState)
  fmt.Printf("CurrState = %d\n", elev.CurrState)
}


func updateCurrState(elev *Elev, state State){
  elev.CurrState = state
}


func updatePrevState(elev *Elev){
  elev.PrevState = elev.CurrState
}


func updateCurrDirection(elev *Elev, direction elevio.MotorDirection) {
  elev.CurrDirection = direction
}


func updatePrevDirection(elev *Elev) {
  elev.PrevDirection = elev.CurrDirection
}


func UpdateFloor(elev *Elev, floor int) {
  elev.Floor = floor
}


func PollInternalElevatorStatus(elev *Elev, status_updated chan bool, send_status_update chan Elev) {
  for {
    select {
    case shouldUpdate := <- status_updated:
      if shouldUpdate {
        send_status_update <- *elev
      }
    }
  }
}
