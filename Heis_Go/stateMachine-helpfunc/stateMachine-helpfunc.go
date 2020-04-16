
package stateMachineHF

import (
	"../elevator"
	"../elevio"
	"../timer"
	"../orderHandler-helpfunc"
)

func ButtonPressedWhileIdle(elev *elevator.Elev, pressedButton elevio.ButtonEvent) {
	//Is the elevator on the floor that was ordered?
	if pressedButton.Floor == elev.Floor{
		stopOnFloor(elev)
		timer.StartTimer(timer.DoorOpenTime)

	} else {
		//direction := chooseDirectionIdle(elev)
		direction := chooseDirection(elev)
		elevio.SetMotorDirection(direction)
		elevator.UpdateElevatorValues (elev, direction, elevator.Moving)
	}
}


func ArrivedOnFloor(elev *elevator.Elev) {
	//Are there any orders this elevator should take at this floor?
	if shouldIStop(elev) {
		//Stopping, opens door and waits
		stopOnFloor(elev)
		timer.StartTimer(timer.DoorOpenTime)
	}
}


func DoorTimeout(elev *elevator.Elev) {
	elevio.SetDoorOpenLamp(false)

	//direction := ChooseDirectionIdle(elev)
	direction := chooseDirection(elev)
	elevio.SetMotorDirection(direction)
	switch direction {
	case elevio.MD_Stop:
		elevator.UpdateElevatorValues(elev, direction, elevator.Idle)
	default :
		elevator.UpdateElevatorValues(elev, direction, elevator.Moving)
	}


}


func chooseDirection(elev *elevator.Elev) elevio.MotorDirection {
	//Deciding where the elevator should go next based on the previous direction
	switch elev.PrevDirection {
	case elevio.MD_Up:
		//Are there any orders above?
		if orderHandlerHF.OrdersAbove(elev) {
			return elevio.MD_Up
		//Are there any order below?
		} else if orderHandlerHF.OrdersBelow(elev) {
			return elevio.MD_Down
		//If not, then stand still
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:
		//Are there any order below?
		if orderHandlerHF.OrdersBelow(elev) {
			return elevio.MD_Down
		//Are there any orders above?
		} else if orderHandlerHF.OrdersAbove(elev) {
			return elevio.MD_Up
		//If not, then stand still
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Stop:
		if orderHandlerHF.OrdersAbove(elev) {
			return elevio.MD_Up
		} else if orderHandlerHF.OrdersBelow(elev){
			return elevio.MD_Down
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

	default:
	}
	return false
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

	elevator.UpdateElevatorValues(elev, elevio.MD_Stop, elevator.DoorOpen)
}
