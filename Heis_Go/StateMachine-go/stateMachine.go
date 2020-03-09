package 
import "control"


func chooseDirection(elevator Elev) Direction {
	switch elevator.Dir {
	case MD_Stop:
		if ordersAbove(elevator) {
			return MD_Up
		} else if ordersBelow(elevator) {
			return MD_Down
		} else {
			return MD_Stop
		}
	case MD_Up:
		if ordersAbove(elevator) {
			return MD_Up
		} else if ordersBelow(elevator) {
			return MD_Down
		} else {
			return MD_Stop
		}

	case MD_Down:
		if ordersBelow(elevator) {
			return MD_Down
		} else if ordersAbove(elevator) {
			return MD_Up
		} else {
			return MD_Stop
		}
	}
	return MD_Stop
}




func ordersAbove(elevator Elev) bool {
	for floor := elevator.Floor + 1; floor < NumFloors; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}

func ordersBelow(elevator Elev) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < NumButtons; btn++ {
			if elevator.Queue[floor][btn] {
				return true
			}
		}
	}
	return false
}



func shouldIStop(elevator Elev) bool {
	switch elevator.Dir {
	case MD_Up:
		return elevator.Queue[elevator.Floor][BT_HallUp] ||
			elevator.Queue[elevator.Floor][BT_Cab] ||
			!ordersAbove(elevator)
	case MD_Down:
		return elevator.Queue[elevator.Floor][BT_HallDown] ||
			elevator.Queue[elevator.Floor][BT_Cab] ||
			!ordersBelow(elevator)
	case MD_Stop:
	default:
	}
	return false
}


func clearOrdersAtCurrentFloor(elevator Elev) Elev {
	elevator.Queue[elevator.floor][BT_Cab] = 0
	switch elevator.Dir {
	case MD_Up:
		elevator.Queue[elevator.Floor][BT_HallUp] = 0
		if !ordersAbove(elevator) {
				elevator.Queue[elevator.Floor][BT_HallDown] = 0
			}
			break
	case MD_Down:
		elevator.Queue[elevator.Floor][BT_HallDown] = 0
		if !ordersBelow(elevator) {
			elevator.Queue[elevator.Floor][BT_HallDown] = 0
		}
		break
	case MD_Stop:
		break
	default:
		fmt.Println("Could not erase orders")
	}
	return elevator
}
