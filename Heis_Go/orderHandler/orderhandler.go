package orderHandler

import (
  "../elevio"
  "../control-go"
  "fmt"
)


func addOrder(elevator *control.Elev, button elevio.ButtonEvent) {
  elevator.OrderList[button.Floor][button.Button] = 1
}


func removeOrder(elevator *control.Elev, button elevio.ButtonType, floor int) {
  elevator.OrderList[floor][button] = 0
}


func OrdersAbove(elevator *control.Elev) bool {
	for floor := elevator.Floor + 1; floor < control.N_FLOORS; floor++ {
		for btn := 0; btn < control.N_BUTTONS; btn++ {
			if elevator.OrderList[floor][btn] == 1 {
				return true
			}
		}
	}
	return false
}


func OrdersBelow(elevator *control.Elev) bool {
	for floor := 0; floor < elevator.Floor; floor++ {
		for btn := 0; btn < control.N_BUTTONS; btn++ {
			if elevator.OrderList[floor][btn] == 1 {
				return true
			}
		}
	}
	return false
}


func ClearOrdersAtCurrentFloor(elevator *control.Elev) {
	//Always delete cab order at current floor
	removeOrder(elevator, elevio.BT_Cab, elevator.Floor)
  removeOrder(elevator, elevio.BT_HallUp, elevator.Floor)
  removeOrder(elevator, elevio.BT_HallDown, elevator.Floor)

  //Turn off all lights at floor
  elevio.SetButtonLamp(elevio.BT_Cab, elevator.Floor, false)
  elevio.SetButtonLamp(elevio.BT_HallUp, elevator.Floor, false)
  elevio.SetButtonLamp(elevio.BT_HallDown, elevator.Floor, false)
  fmt.Printf("%+v\n", elevator.OrderList)
}


func handleOrder(elevator *control.Elev, order elevio.ButtonEvent) {
  addOrder(elevator, order)

  fmt.Printf("%+v\n", elevator.OrderList)
  fmt.Printf("%+v\n", elevator.CurrState)

  elevio.SetButtonLamp(order.Button,order.Floor, true)
}


func StartOrderHandling(elevator *control.Elev, order_from_fsm chan elevio.ButtonEvent, order_from_order_distributer chan elevio.ButtonEvent, distribute_order chan elevio.ButtonEvent, new_order chan elevio.ButtonEvent, order_executed chan bool) {
  for {
    select {
    case order := <- order_from_fsm:
      if order.Button == elevio.BT_Cab {
        handleOrder(elevator,order)
        new_order <- order
      } else {
        distribute_order <- order
      }
    case order := <- order_from_order_distributer:
      handleOrder(elevator, order)
      new_order <- order

    case <- order_executed:
      ClearOrdersAtCurrentFloor(elevator)
  }
  }
}
