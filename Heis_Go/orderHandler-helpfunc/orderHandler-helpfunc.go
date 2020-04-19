 package orderHandlerHF

import (
  "../elevator"
  "../elevio"
)




func AddOrder(elev *elevator.Elev, button elevio.ButtonEvent) {
  elev.OrderList[button.Floor][button.Button] = 1
}


func RemoveOrder(elev *elevator.Elev, button elevio.ButtonType, floor int) {
  elev.OrderList[floor][button] = 0
}


func OrdersAbove(elev *elevator.Elev) bool {
	for floor := elev.Floor + 1; floor < elevator.N_FLOORS; floor++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if elev.OrderList[floor][btn] == 1 {
				return true
			}
		}
	}
	return false
}


func OrdersBelow(elev *elevator.Elev) bool {
	for floor := 0; floor < elev.Floor; floor++ {
		for btn := 0; btn < elevator.N_BUTTONS; btn++ {
			if elev.OrderList[floor][btn] == 1 {
				return true
			}
		}
	}
	return false
}


func ClearOrdersAtCurrentFloor(elev *elevator.Elev) {
	//Delete orders
	RemoveOrder(elev, elevio.BT_Cab, elev.Floor)
  RemoveOrder(elev, elevio.BT_HallUp, elev.Floor)
  RemoveOrder(elev, elevio.BT_HallDown, elev.Floor)

  //Turn off all lights at floor
  elevio.SetButtonLamp(elevio.BT_Cab, elev.Floor, false)
  elevio.SetButtonLamp(elevio.BT_HallUp, elev.Floor, false)
  elevio.SetButtonLamp(elevio.BT_HallDown, elev.Floor, false)
}


func TakeOrder(elev *elevator.Elev, order elevio.ButtonEvent) {
  AddOrder(elev, order)

  elevio.SetButtonLamp(order.Button,order.Floor, true)
}
