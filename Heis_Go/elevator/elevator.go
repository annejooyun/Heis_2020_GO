package elevator

import (
  "../elevio"
  "fmt"
)

const N_FLOORS = 4
const N_BUTTONS = 3 // Number of button types (cab,hall up, hall down)
//const N_ELEVS = 3

var LOCAL_ELEV_ID string

type State int
const (
  Idle   State = 0
  Moving       = 1
  DoorOpen     = 2
)

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

  //Turn of all lights
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


func UpdateElevatorDirectionsAndStates(elevator *Elev, direction elevio.MotorDirection, state State) {
  UpdatePrevDirection(elevator)
  UpdateCurrDirection(elevator, direction)
  UpdatePrevState(elevator)
  UpdateCurrState(elevator, state)
}


func UpdateCurrState(elevator *Elev, state State){
  elevator.CurrState = state
  //fmt.Println("New CurrState")
  //fmt.Printf("%+v\n", elevator.CurrState)
}


func UpdatePrevState(elevator *Elev){
  elevator.PrevState = elevator.CurrState
  //fmt.Println("New PrevState")
  //fmt.Printf("%+v\n", elevator.PrevState)
}


func UpdateCurrDirection(elevator *Elev, direction elevio.MotorDirection) {
  elevator.CurrDirection = direction
  //fmt.Println("New CurrDirection")
  //fmt.Printf("%+v\n", elevator.CurrDirection)
}


func UpdatePrevDirection(elevator *Elev) {
  elevator.PrevDirection = elevator.CurrDirection
  //fmt.Println("New PrevDirection")
  //fmt.Printf("%+v\n", elevator.PrevDirection)
}


func UpdateFloor(elevator *Elev, floor int) {
  elevator.Floor = floor
}


func PollInternalElevatorStatus(elevator *Elev, status_updated chan bool, send_status_update chan Elev) {
  for {
    select {
    case shouldUpdate := <- status_updated:
      if shouldUpdate {
        send_status_update <- *elevator
      }
    }
  }
}
