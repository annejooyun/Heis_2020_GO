package stateMachine

import (
	"../stateMachine-helpfunc"
	"../elevator"
	"../elevio"
	"fmt"
	)



func RunStateMachine(elev *elevator.Elev,
										 order_registered chan elevio.ButtonEvent,
										 new_order chan elevio.ButtonEvent,
										 status_updated chan bool,
										 order_executed chan bool,
										 drv_buttons chan elevio.ButtonEvent,
										 drv_floors chan int,
										 ch_timer chan bool) {



	for {
        select {

				//There has been added a new order to the order list.
				case order := <- new_order:
					if elev.CurrState == elevator.Idle {
						stateMachineHF.ButtonPressedWhileIdle(elev,order)
					}


				//A button has been pressed
        case buttonPressed := <- drv_buttons:
						//The order is sent on the channel order_registered to the order handler
						fmt.Printf("Button pushed\n")
						order_registered <- buttonPressed
						fmt.Printf("Button registered\n")


				//A new floor is detected
				case floor := <- drv_floors:
					fmt.Printf("Floor detected\n")

						//Store the previous floor we were on
            prevFloor := elev.Floor

						//Update the current floor to the floor detected
            elevator.UpdateFloor(elev, floor)
            elevio.SetFloorIndicator(floor)

            if floor != -1 && elev.Floor != prevFloor {
              stateMachineHF.ArrivedOnFloor(elev)
            }
						status_updated <- true


					//Timeout
	        case <- ch_timer:
						fmt.Printf("Timeout\n")
	          stateMachineHF.DoorTimeout(elev)
						order_executed <- true
						status_updated <- true

        }
    }
	}
