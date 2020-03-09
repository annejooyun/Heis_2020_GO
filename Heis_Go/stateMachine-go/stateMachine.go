package stateMachine
import "../control-go"
import "../elevio"


func ChooseDirection(elevator control.Elev) elevio.MotorDirection {
	switch elevator.Direction {
	case elevio.MD_Stop:
		if OrdersAbove(elevator) {
			return elevio.MD_Up
		} else if OrdersBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}
	case elevio.MD_Up:
		if OrdersAbove(elevator) {
			return elevio.MD_Up
		} else if OrdersBelow(elevator) {
			return elevio.MD_Down
		} else {
			return elevio.MD_Stop
		}

	case elevio.MD_Down:
		if OrdersBelow(elevator) {
			return elevio.MD_Down
		} else if OrdersAbove(elevator) {
			return elevio.MD_Up
		} else {
			return elevio.MD_Stop
		}
	}
	return elevio.MD_Stop
}




func OrdersAbove(elevator control.Elev) bool {
	for floor := elevator.PrevFloor + 1; floor < control.N_FLOORS; floor++ {
		for btn := 0; btn < control.N_BUTTONS; btn++ {
			if elevator.OrderList[floor][btn] == 1 {
				return true
			}
		}
	}
	return false
}

func OrdersBelow(elevator control.Elev) bool {
	for floor := 0; floor < elevator.PrevFloor; floor++ {
		for btn := 0; btn < control.N_BUTTONS; btn++ {
			if elevator.OrderList[floor][btn] == 1 {
				return true
			}
		}
	}
	return false
}



func shouldIStop(elevator control.Elev) bool {
	switch elevator.Direction {
	case elevio.MD_Up:
		//If motor direction is up, stop if there is a button call up,
		//or cab call on floor, or not orders above.
		return elevator.OrderList[elevator.PrevFloor][elevio.BT_HallUp] == 1 ||
			elevator.OrderList[elevator.PrevFloor][elevio.BT_Cab] == 1 ||
			!OrdersAbove(elevator)
	case elevio.MD_Down:
		//If motor direction is down, stop if there is a button call down,
		//or cab call on floor, or no orders below.
		return elevator.OrderList[elevator.PrevFloor][elevio.BT_HallDown] == 1 ||
			elevator.OrderList[elevator.PrevFloor][elevio.BT_Cab] == 1 ||
			!OrdersBelow(elevator)
	case elevio.MD_Stop:
	default:
	}
	return false
}
