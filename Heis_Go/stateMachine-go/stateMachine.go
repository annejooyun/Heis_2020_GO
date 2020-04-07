package stateMachine

import (
	"../control-go"
	"../elevio"
	"../orderHandler"
	"../timer"
	"fmt"
	//"../messageHandler"
	)


func buttonPressedWhileIdle(elevator *control.Elev, firstButton elevio.ButtonEvent) {
	//Is the elevator on the floor that was ordered?
	if firstButton.Floor == elevator.Floor{
		stopOnFloor(elevator)
		timer.StartTimer(timer.DoorOpenTime)

	} else {
		//direction := chooseDirectionIdle(elevator)
		direction := ChooseDirection(elevator)
		elevio.SetMotorDirection(direction)
		control.UpdateElevatorValues (elevator, direction, control.Moving)
	}
}


func arrivedOnFloor(elevator *control.Elev) {
	//Are there any orders this elevator should take at this floor?
	if ShouldIStop(elevator) {
		//Stopping, opens door and waits
		stopOnFloor(elevator)
		timer.StartTimer(timer.DoorOpenTime)
	}
}


func doorTimeout(elevator *control.Elev) {
	//orderhandler.ClearOrdersAtCurrentFloor(elevator)
	elevio.SetDoorOpenLamp(false)

	//direction := ChooseDirectionIdle(elevator)
	direction := ChooseDirection(elevator)
	elevio.SetMotorDirection(direction)
	switch direction {
	case elevio.MD_Stop:
		control.UpdateElevatorValues(elevator, direction, control.Idle)
	default :
		control.UpdateElevatorValues(elevator, direction, control.Moving)
	}


}


func ChooseDirection(elevator *control.Elev) elevio.MotorDirection {
	//Deciding where the elevator should go next based on the previous direction
	switch elevator.PrevDirection {
	case elevio.MD_Up:
		//Are there any orders above?
		if orderHandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		//Are there any order below?
		} else if orderHandler.OrdersBelow(elevator) {
			return elevio.MD_Down
		//If not, then stand still
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:
		//Are there any order below?
		if orderHandler.OrdersBelow(elevator) {
			return elevio.MD_Down
		//Are there any orders above?
		} else if orderHandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		//If not, then stand still
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Stop:
		if orderHandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		} else if orderHandler.OrdersBelow(elevator){
			return elevio.MD_Down
		}
	}
	return elevio.MD_Stop
}


func ShouldIStop(elevator *control.Elev) bool {
	switch elevator.CurrDirection {

	case elevio.MD_Up:
		return  isOrderOnFloor(elevator) || !orderHandler.OrdersAbove(elevator)

	case elevio.MD_Down:
		return isOrderOnFloor(elevator)|| !orderHandler.OrdersBelow(elevator)

	case elevio.MD_Stop:

	default:
	}
	return false
}


func isOrderOnFloor(elevator *control.Elev) bool {
	if elevator.CurrDirection == elevio.MD_Up {
		return elevator.OrderList[elevator.Floor][elevio.BT_HallUp] == 1 ||
			elevator.OrderList[elevator.Floor][elevio.BT_Cab] == 1

	} else if elevator.CurrDirection == elevio.MD_Down {
		return elevator.OrderList[elevator.Floor][elevio.BT_HallDown] == 1 ||
			elevator.OrderList[elevator.Floor][elevio.BT_Cab] == 1
	} else {
		return false
	}
}


func stopOnFloor(elevator *control.Elev) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)

	control.UpdateElevatorValues(elevator, elevio.MD_Stop, control.DoorOpen)
}


func RunStateMachine(elevator *control.Elev, send_to_order_handler chan elevio.ButtonEvent, new_order chan elevio.ButtonEvent, status_updated chan bool, order_executed chan bool) {
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
					if elevator.CurrState == control.Idle {
						buttonPressedWhileIdle(elevator,order)
					}


        case buttonPressed := <- drv_buttons:
					send_to_order_handler <- buttonPressed


				case floor := <- drv_floors: //A new floor is detected
						fmt.Printf("Floor detected:\n")
						fmt.Printf("%+v\n", floor)
            prevFloor := elevator.Floor //Store the previous floor we were on
            control.UpdateFloor(elevator, floor) //Update the current floor to the floor detected.
            elevio.SetFloorIndicator(floor)

            if floor != -1 && floor != prevFloor {
              arrivedOnFloor(elevator) //We have arrived on a new floor
							//status_updated <- true //There has been a change in status
            }
						status_updated <- true


	        case timeOut := <- ch_timer:
	          fmt.Printf("TIMEOUT = ")
	          fmt.Printf("%+v\n", timeOut)
	          doorTimeout(elevator)
						order_executed <- true
						status_updated <- true //There may have been a change in status

        }
    }
	}
