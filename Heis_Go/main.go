package main

import "./elevio"
import "./control-go"
import "./orderHandler"
import "./stateMachine-go"
//import "./timer"
import "fmt"

func main(){
    elevator := control.InitializeElevator()




    fmt.Println(elevator)

    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    //drv_obstr   := make(chan bool)
    //drv_stop    := make(chan bool)
  //  timer       := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    //go elevio.PollObstructionSwitch(drv_obstr)
    //go elevio.PollStopButton(drv_stop)
    //go timer.PollTimeOut(timer)

    for {
        select {
        case buttonPressed := <- drv_buttons:
            orderhandler.AddOrder(&elevator, buttonPressed)
            //fmt.Printf("%+v\n", buttonPressed)
            //fmt.Printf("%+v\n", elevator)
            //elevio.SetButtonLamp(buttonPressed.Button, buttonPressed.Floor, true)
            //if elevator.State == control.Idle {
            //  stateMachine.ButtonPressedWhileIdle(&elevator,buttonPressed)
            //}


          case floor := <- drv_floors:
            fmt.Printf("%+v\n", floor)
            control.UpdatePrevFloor(&elevator, floor)
            if floor != -1 && floor != elevator.PrevFloor {
              stateMachine.ArrivedOnFloor(&elevator)
            }

            //makes sure the elevator does not go out of bounds
            if floor == control.N_FLOORS-1 {
                elevio.SetMotorDirection(elevio.MD_Down)
            } else if floor == 0 {
                elevio.SetMotorDirection(elevio.MD_Up)
            }

  /*
        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(elevio.MD_Up)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < control.N_FLOORS; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }

        case timeOut := <- timer:
          fmt.Printf("TIMEOUT")*/
        }
    }
}
