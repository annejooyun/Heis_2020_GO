package main

import "./elevio"
import "./control-go"
import "./orderHandler"
import "fmt"

func main(){
    elevator := control.InitializeElevator()


    var d elevio.MotorDirection = elevio.MD_Up
    elevio.SetMotorDirection(d)
    fmt.Println(elevator)

    drv_buttons := make(chan elevio.ButtonEvent)
    drv_floors  := make(chan int)
    drv_obstr   := make(chan bool)
    drv_stop    := make(chan bool)

    go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go elevio.PollObstructionSwitch(drv_obstr)
    go elevio.PollStopButton(drv_stop)


    for {
        select {
        case buttonPressed := <- drv_buttons:
            orderhandler.AddOrder(&elevator, buttonPressed)
            fmt.Printf("%+v\n", buttonPressed)
            fmt.Printf("%+v\n", elevator)
            elevio.SetButtonLamp(buttonPressed.Button, buttonPressed.Floor, true)
            

        case a := <- drv_floors:
            fmt.Printf("%+v\n", a)

            if a == control.N_FLOORS-1 {
                d = elevio.MD_Down
            } else if a == 0 {
                d = elevio.MD_Up
            }
            elevio.SetMotorDirection(d)


        case a := <- drv_obstr:
            fmt.Printf("%+v\n", a)
            if a {
                elevio.SetMotorDirection(elevio.MD_Stop)
            } else {
                elevio.SetMotorDirection(d)
            }

        case a := <- drv_stop:
            fmt.Printf("%+v\n", a)
            for f := 0; f < control.N_FLOORS; f++ {
                for b := elevio.ButtonType(0); b < 3; b++ {
                    elevio.SetButtonLamp(b, f, false)
                }
            }
        }
    }
}
