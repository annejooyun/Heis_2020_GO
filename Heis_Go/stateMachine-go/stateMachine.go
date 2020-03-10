package stateMachine
import "../control-go"
import "../elevio"
import "../orderHandler"
import "fmt"


func ButtonPressedWhileIdle(elevator *control.Elev, firstButton elevio.ButtonEvent) {
	if firstButton.Floor == elevator.PrevFloor{
		//elevio.SetDoorOpenLamp(true)
		//timer 3 sek
		fmt.Printf("I am already here")
	} else {
		direction := ChooseDirection(elevator)
		elevio.SetMotorDirection(direction)
		control.UpdateDirection(elevator, direction)
		control.UpdateState(elevator, control.Moving)
	}
}


func ArrivedOnFloor(elevator *control.Elev) {


	if shouldIStop(elevator) {
		lastDirection := elevator.Direction
		elevio.SetMotorDirection(elevio.MD_Stop)
		//elevio.SetDoorOpenLamp(true)
		control.UpdateState(elevator, control.DoorOpen)
		//timer 3 sek
		orderhandler.ClearOrdersAtCurrentFloor(elevator)

		//Deciding where the elevator should go next
		switch lastDirection {
		case elevio.MD_Up:
			//Are there any orders above?
			if orderhandler.OrdersAbove(elevator) {
				elevio.SetMotorDirection(elevio.MD_Up)
				control.UpdateState(elevator, control.Moving)
			//Are there any order below?
			} else if orderhandler.OrdersBelow(elevator) {
				elevio.SetMotorDirection(elevio.MD_Down)
				control.UpdateState(elevator, control.Moving)
			//If not, then stand still
			} else {
				elevio.SetMotorDirection(elevio.MD_Stop)
				control.UpdateState(elevator, control.Idle)
			}

		case elevio.MD_Down:
			//Are there any order below?
			if orderhandler.OrdersBelow(elevator) {
				elevio.SetMotorDirection(elevio.MD_Down)
				control.UpdateState(elevator, control.Moving)
			//Are there any orders above?
			} else if orderhandler.OrdersAbove(elevator) {
				elevio.SetMotorDirection(elevio.MD_Up)
				control.UpdateState(elevator, control.Moving)
			//If not, then stand still
			} else {
				elevio.SetMotorDirection(elevio.MD_Stop)
				control.UpdateState(elevator, control.Idle)
			}
		default:
			elevio.SetMotorDirection(elevio.MD_Stop)
			control.UpdateState(elevator, control.Idle)
		}
	}
}



func doorTimeout(elevator *control.Elev) {

	orderhandler.ClearOrdersAtCurrentFloor(elevator)
}









func ChooseDirection(elevator *control.Elev) elevio.MotorDirection {
	switch elevator.Direction {
	case elevio.MD_Stop:
		if orderhandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		} else if orderhandler.OrdersBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}
	case elevio.MD_Up:
		if orderhandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		} else if orderhandler.OrdersBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:
		if orderhandler.OrdersBelow(elevator) {
			return elevio.MD_Down
		} else if orderhandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		} else {
			return elevio.MD_Stop
		}
	}
	return elevio.MD_Stop
}


func shouldIStop(elevator *control.Elev) bool {
	switch elevator.Direction {
	case elevio.MD_Up:
		//If motor direction is UP, stop if there is a button call up,
		//or cab call on floor, or not orders above.
		return elevator.OrderList[elevator.PrevFloor][elevio.BT_HallUp] == 1 ||
			elevator.OrderList[elevator.PrevFloor][elevio.BT_Cab] == 1 ||
			!orderhandler.OrdersAbove(elevator)
	case elevio.MD_Down:
		//If motor direction is DOWN, stop if there is a button call down,
		//or cab call on floor, or no orders below.
		return elevator.OrderList[elevator.PrevFloor][elevio.BT_HallDown] == 1 ||
			elevator.OrderList[elevator.PrevFloor][elevio.BT_Cab] == 1 ||
			!orderhandler.OrdersBelow(elevator)
	case elevio.MD_Stop:
	default:
	}
	return false
}
