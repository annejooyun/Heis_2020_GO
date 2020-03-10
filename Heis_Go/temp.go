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
