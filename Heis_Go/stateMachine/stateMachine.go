package stateMachine

import (
	"../stateMachine-helpfunc"
	"../elevator"
	"../elevio"
	"../timer"

	//"fmt"
	)



func RunStateMachine(elev *elevator.Elev, order_registered chan elevio.ButtonEvent, new_order chan elevio.ButtonEvent, status_updated chan bool, order_executed chan bool) {

	drv_buttons := make(chan elevio.ButtonEvent)
	drv_floors  := make(chan int)
	ch_timer    := make(chan bool)

	go elevio.PollButtons(drv_buttons)
  go elevio.PollFloorSensor(drv_floors)
  go timer.PollTimeOut(ch_timer)

	for {
        select {

				//There has been added a new order to the order list.
				case order := <- new_order:
					if elev.CurrState == elevator.Idle {
						stateMachineHF.ButtonPressedWhileIdle(elev,order)
					}


				//A button has been pressed
        case buttonPressed := <- drv_buttons:
						order_registered <- buttonPressed


				//A new floor is detected
				case floor := <- drv_floors:
						//fmt.Printf("Floor detected:\n")
						//fmt.Printf("%+v\n", floor)

						//Store the previous floor we were on
            prevFloor := elev.Floor
						//Update the current floor to the floor detected
            elevator.UpdateFloor(elev, floor)
            elevio.SetFloorIndicator(floor)

            if floor != -1 && floor != prevFloor {
              stateMachineHF.ArrivedOnFloor(elev) //We have arrived on a new floor
							//status_updated <- true //There has been a change in status
            }
						status_updated <- true

					//Timeout
	        case <- ch_timer:
	          //fmt.Printf("TIMEOUT = ")
	          //fmt.Printf("%+v\n", timeOut)
	          stateMachineHF.DoorTimeout(elev)
						order_executed <- true
						status_updated <- true //There may have been a change in status

        }
    }
	}
