package control

import "../elevio"
import "fmt"

const N_FLOORS = 4
const N_BUTTONS = 3 // Number of button types (cab,hall up, hall down)
const Elev_ID = 1

type State int
const (
  Idle   State = 0
  Moving       = 1
  DoorOpen     = 2
)

type Elev struct {
  Id int
  CurrState State
  PrevState State
  CurrDirection elevio.MotorDirection
  PrevDirection elevio.MotorDirection
  Floor int  //current or previous floor
  OrderList [N_FLOORS][N_BUTTONS]int
}

//Moves the elevator to a known floor. Makes elevator object. Turns off lights
func InitializeElevator() Elev {
  elevio.Init("localhost:15657", N_FLOORS)
  elevio.SetMotorDirection(elevio.MD_Down)
  floor := -1
  for floor == -1 {
    floor = elevio.GetFloor()
  }
  elevio.SetMotorDirection(elevio.MD_Stop)
  initializeLights()
  elevator := Elev{Elev_ID,Idle,Idle,elevio.MD_Stop,elevio.MD_Stop,floor,[N_FLOORS][N_BUTTONS]int{}}

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


func UpdateElevatorValues (elevator *Elev, direction elevio.MotorDirection, state State) {

  UpdatePrevDirection(elevator)
  UpdateCurrDirection(elevator, direction)
  UpdatePrevState(elevator)
  UpdateCurrState(elevator, state)

}





func UpdateCurrState(elevator *Elev, state State){
  elevator.CurrState = state
  fmt.Println("New CurrState")
  fmt.Printf("%+v\n", elevator)
}

func UpdatePrevState(elevator *Elev){
  elevator.PrevState = elevator.CurrState
  fmt.Println("New PrevState")
  fmt.Printf("%+v\n", elevator)
}

func UpdateCurrDirection(elevator *Elev, direction elevio.MotorDirection) {
  elevator.CurrDirection = direction
  fmt.Println("New CurrDirection")
  fmt.Printf("%+v\n", elevator)
}

func UpdatePrevDirection(elevator *Elev) {
  elevator.PrevDirection = elevator.CurrDirection
  fmt.Println("New PrevDirection")
  fmt.Printf("%+v\n", elevator)
}

func UpdateFloor(elevator *Elev, floor int) {
  elevator.Floor = floor
}
