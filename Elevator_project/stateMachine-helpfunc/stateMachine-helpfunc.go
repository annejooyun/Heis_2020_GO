
package stateMachineHF

import (
	"../elevator"
	"../elevio"
	"../orderHandler-helpfunc"
	"../timer"
)




func ButtonPressedWhileIdle(elev *elevator.Elev, pressedButton elevio.ButtonEvent) {

	//Is the elevator on the floor that was ordered?
	if pressedButton.Floor == elev.Floor{
		stopOnFloor(elev)
		timer.StartTimer(timer.DOOR_OPEN_TIME)

	} else {
		direction := chooseDirection(elev)
		elevio.SetMotorDirection(direction)
		elevator.UpdateElevatorDirectionsAndStates (elev, direction, elevator.Moving)
	}
}


func ArrivedOnFloor(elev *elevator.Elev) {

	if shouldIStop(elev) {
		stopOnFloor(elev)
		timer.StartTimer(timer.DOOR_OPEN_TIME)
	}
}


func DoorTimeout(elev *elevator.Elev) {

	elevio.SetDoorOpenLamp(false)

	direction := chooseDirection(elev)
	elevio.SetMotorDirection(direction)

	switch direction {
	case elevio.MD_Stop:
		elevator.UpdateElevatorDirectionsAndStates(elev, direction, elevator.Idle)
	default :
		elevator.UpdateElevatorDirectionsAndStates(elev, direction, elevator.Moving)
	}
}


func chooseDirection(elev *elevator.Elev) elevio.MotorDirection {

	//Deciding where the elevator should go next based on the previous direction
	if elev.Floor == 0 && orderHandlerHF.OrdersAbove(elev){
		return elevio.MD_Up

	}else if elev.Floor == (elevator.N_FLOORS - 1) && orderHandlerHF.OrdersBelow(elev){
		return elevio.MD_Down
	}

	switch elev.PrevDirection {
	case elevio.MD_Up:

		if orderHandlerHF.OrdersAbove(elev) {
			return elevio.MD_Up

		} else if orderHandlerHF.OrdersBelow(elev) {
			return elevio.MD_Down

		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:

		if orderHandlerHF.OrdersBelow(elev) {
			return elevio.MD_Down

		} else if orderHandlerHF.OrdersAbove(elev) {
			return elevio.MD_Up

		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Stop:

		if orderHandlerHF.OrdersAbove(elev) {
			return elevio.MD_Up

		} else if orderHandlerHF.OrdersBelow(elev){
			return elevio.MD_Down

		} else{
			return elevio.MD_Stop
		}
	}
	return elevio.MD_Stop
}


func shouldIStop(elev *elevator.Elev) bool {
	switch elev.CurrDirection {

	case elevio.MD_Up:
		return  isOrderOnFloor(elev) || !orderHandlerHF.OrdersAbove(elev)

	case elevio.MD_Down:
		return isOrderOnFloor(elev)|| !orderHandlerHF.OrdersBelow(elev)

	case elevio.MD_Stop:
		return true
	}
	return true
}


func isOrderOnFloor(elev *elevator.Elev) bool {
	if elev.CurrDirection == elevio.MD_Up {
		return elev.OrderList[elev.Floor][elevio.BT_HallUp] == 1 ||
			elev.OrderList[elev.Floor][elevio.BT_Cab] == 1

	} else if elev.CurrDirection == elevio.MD_Down {
		return elev.OrderList[elev.Floor][elevio.BT_HallDown] == 1 ||
			elev.OrderList[elev.Floor][elevio.BT_Cab] == 1
	} else {
		return false
	}
}


func stopOnFloor(elev *elevator.Elev) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)

	elevator.UpdateElevatorDirectionsAndStates(elev, elevio.MD_Stop, elevator.DoorOpen)
}
