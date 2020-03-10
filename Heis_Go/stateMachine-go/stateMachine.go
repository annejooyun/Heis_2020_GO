package stateMachine
import "../control-go"
import "../elevio"
import "../orderHandler"
import "../timer"
import "fmt"


func ButtonPressedWhileIdle(elevator *control.Elev, firstButton elevio.ButtonEvent) {
	//Is the elevator on the floor that was ordered?
	if firstButton.Floor == elevator.Floor{

		stopOnFloor(elevator)

		//SOMETHING WRONG HERE
		newDirection := ChooseDirectionIdle(elevator)
		fmt.Printf("NewDirection is ")
		fmt.Printf("%+v\n", newDirection)
		control.UpdatePrevDirection(elevator)
		control.UpdateCurrDirection(elevator, newDirection)
		timer.TimerStart(timer.DoorOpenTime)
	} else {
		direction := ChooseDirectionIdle(elevator)
		elevio.SetMotorDirection(direction)
		control.UpdatePrevDirection(elevator)
		control.UpdateCurrDirection(elevator, direction)
		control.UpdateCurrState(elevator, control.Moving)
	}
}


func ArrivedOnFloor(elevator *control.Elev) {

	//Are there any orders this elevator should take at this floor?
	if shouldIStop(elevator) {
		//Stopping, opens door and waits
		stopOnFloor(elevator)
		newDirection := ChooseDirection(elevator)
		control.UpdatePrevDirection(elevator)
		control.UpdateCurrDirection(elevator, newDirection)
		timer.TimerStart(timer.DoorOpenTime)
	}
}



func DoorTimeout(elevator *control.Elev) {
	orderhandler.ClearOrdersAtCurrentFloor(elevator)
	elevio.SetDoorOpenLamp(false)
	if elevator.CurrDirection != elevio.MD_Stop {
		control.UpdateCurrState(elevator, control.Moving)
	} else {
		control.UpdateCurrState(elevator, control.Idle)
	}
}




func ChooseDirectionIdle(elevator *control.Elev) elevio.MotorDirection {
	if orderhandler.OrdersAbove(elevator) == true {
		fmt.Println("There are orders above")
		return elevio.MD_Up
	} else if orderhandler.OrdersBelow(elevator) == true {
		fmt.Println("There are orders below")
		return elevio.MD_Down
	} else {
		fmt.Println("There are no orders")
		return elevio.MD_Stop
	}
}


func ChooseDirection(elevator *control.Elev) elevio.MotorDirection {
	//Deciding where the elevator should go next based on the previous direction
	switch elevator.PrevDirection {
	case elevio.MD_Up:
		//Are there any orders above?
		if orderhandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		//Are there any order below?
		} else if orderhandler.OrdersBelow(elevator) {
			return elevio.MD_Down
		//If not, then stand still
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:
		//Are there any order below?
		if orderhandler.OrdersBelow(elevator) {
			return elevio.MD_Down
		//Are there any orders above?
		} else if orderhandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		//If not, then stand still
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Stop:
		if orderhandler.OrdersAbove(elevator) {
			return elevio.MD_Up
		} else if orderhandler.OrdersBelow(elevator){
			return elevio.MD_Down
		}
	}
	return elevio.MD_Stop
}


func shouldIStop(elevator *control.Elev) bool {
	switch elevator.CurrDirection {
	case elevio.MD_Up:
		//If motor direction is UP, stop if there is a button call up,
		//or cab call on floor, or not orders above.
		return elevator.OrderList[elevator.Floor][elevio.BT_HallUp] == 1 ||
			elevator.OrderList[elevator.Floor][elevio.BT_Cab] == 1 ||
			!orderhandler.OrdersAbove(elevator)
	case elevio.MD_Down:
		//If motor direction is DOWN, stop if there is a button call down,
		//or cab call on floor, or no orders below.
		return elevator.OrderList[elevator.Floor][elevio.BT_HallDown] == 1 ||
			elevator.OrderList[elevator.Floor][elevio.BT_Cab] == 1 ||
			!orderhandler.OrdersBelow(elevator)
	case elevio.MD_Stop:
	default:
	}
	return false
}


func stopOnFloor(elevator *control.Elev) {
	elevio.SetMotorDirection(elevio.MD_Stop)
	elevio.SetDoorOpenLamp(true)
	control.UpdatePrevDirection(elevator)
	control.UpdateCurrDirection(elevator, elevio.MD_Stop)
	control.UpdatePrevState(elevator)
	control.UpdateCurrState(elevator, control.DoorOpen)
}
