package stateMachine

import (
	"../stateMachine-helpfunc"
	"../elevator"
	"../elevio"
	"../timer"
	"fmt"
	//"../messageHandler"
	)



func RunStateMachine(elev *elevator.Elev, send_to_order_handler chan elevio.ButtonEvent, new_order chan elevio.ButtonEvent, status_updated chan bool, order_executed chan bool) {
	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	//drv_obstr   := make(chan bool)
	//drv_stop    := make(chan bool)
	ch_timer     := make(chan bool)

	go elevio.PollButtons(drv_buttons)
  go elevio.PollFloorSensor(drv_floors)
  //go elevio.PollObstructionSwitch(drv_obstr)
  //go elevio.PollStopButton(drv_stop)
  go timer.PollTimeOut(ch_timer)

	for {
        select {
				case order := <- new_order: //There has been added a new order to the order list.
					if elev.CurrState == elevator.Idle {
						stateMachineHF.ButtonPressedWhileIdle(elev,order)
					}


        case buttonPressed := <- drv_buttons:
					send_to_order_handler <- buttonPressed


				case floor := <- drv_floors: //A new floor is detected
						fmt.Printf("Floor detected:\n")
						fmt.Printf("%+v\n", floor)
            prevFloor := elev.Floor //Store the previous floor we were on
            elevator.UpdateFloor(elev, floor) //Update the current floor to the floor detected.
            elevio.SetFloorIndicator(floor)

            if floor != -1 && floor != prevFloor {
              stateMachineHF.ArrivedOnFloor(elev) //We have arrived on a new floor
							//status_updated <- true //There has been a change in status
            }
						status_updated <- true


	        case timeOut := <- ch_timer:
	          fmt.Printf("TIMEOUT = ")
	          fmt.Printf("%+v\n", timeOut)
	          stateMachineHF.DoorTimeout(elev)
						order_executed <- true
						status_updated <- true //There may have been a change in status

        }
    }
	}
