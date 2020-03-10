package orderhandler

import "../elevio"
import "../control-go"
import "../stateMachine-go"
import "fmt"

func AddOrder(elevator *control.Elev, button elevio.ButtonEvent) {
  elevator.OrderList[button.Floor][button.Button] = 1
}


func removeOrder(elevator *control.Elev, button elevio.ButtonType, floor int) {
  elevator.OrderList[floor][button] = 0
}



func ClearOrdersAtCurrentFloor(elevator control.Elev) control.Elev {
	//Always delete cab order at current floor
	removeOrder(&elevator, elevio.BT_Cab, elevator.PrevFloor)
	switch elevator.Direction {
	case elevio.MD_Up:
		//If direction is up, delete orders of type Hall Up.
    removeOrder(&elevator, elevio.BT_HallUp, elevator.PrevFloor)
		if !stateMachine.OrdersAbove(elevator) {
        removeOrder(&elevator, elevio.BT_HallDown, elevator.PrevFloor)
			}
			break
	case elevio.MD_Down:
		//If direction is down, delete orders of type Hall down.
		removeOrder(&elevator, elevio.BT_HallDown, elevator.PrevFloor)
		if !stateMachine.OrdersBelow(elevator) {
			removeOrder(&elevator, elevio.BT_HallUp, elevator.PrevFloor)
		}
		break
	case elevio.MD_Stop:
		break
	default:
		fmt.Println("Could not erase orders")
	}
	return elevator
}
