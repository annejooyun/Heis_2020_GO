package orderHandler

import (
  "../elevator"
  "../elevio"
  "fmt"
)


func addOrder(elev *elevator.Elev, button elevio.ButtonEvent) {
  elev.OrderList[button.Floor][button.Button] = 1
}


func removeOrder(elev *elevator.Elev, button elevio.ButtonType, floor int) {
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
	//Always delete cab order at current floor
	removeOrder(elev, elevio.BT_Cab, elev.Floor)
  removeOrder(elev, elevio.BT_HallUp, elev.Floor)
  removeOrder(elev, elevio.BT_HallDown, elev.Floor)

  //Turn off all lights at floor
  elevio.SetButtonLamp(elevio.BT_Cab, elev.Floor, false)
  elevio.SetButtonLamp(elevio.BT_HallUp, elev.Floor, false)
  elevio.SetButtonLamp(elevio.BT_HallDown, elev.Floor, false)
  fmt.Printf("%+v\n", elev.OrderList)
}


func takeOrder(elev *elevator.Elev, order elevio.ButtonEvent) {
  addOrder(elev, order)

  fmt.Printf("%+v\n", elev.OrderList)
  fmt.Printf("%+v\n", elev.CurrState)

  elevio.SetButtonLamp(order.Button,order.Floor, true)
}


///HEI FIKS DENNE DET ER FOR MANGE INPUTSS

func StartOrderHandling(elev *elevator.Elev, order_from_fsm chan elevio.ButtonEvent, order_from_order_distributer chan elevio.ButtonEvent, distribute_order chan elevio.ButtonEvent, new_order chan elevio.ButtonEvent, order_executed chan bool, bcast_order_executed_at_floor chan []int, orders_executed chan []int) {
  for {
    select {
    case order := <- order_from_fsm:
      if order.Button == elevio.BT_Cab {
        takeOrder(elev,order)
        new_order <- order
      } else {
        distribute_order <- order
      }

    case order := <- order_from_order_distributer:
      takeOrder(elev, order)
      new_order <- order

    case <- order_executed:
      //create list of orders executed on the form [button up, button down, floor]
      ordersExecuted := []int{elev.OrderList[elev.Floor][0],elev.OrderList[elev.Floor][1]}
      ordersExecuted = append(ordersExecuted, elev.Floor)

      bcast_order_executed_at_floor <- ordersExecuted
      orders_executed<- ordersExecuted

      ClearOrdersAtCurrentFloor(elev)
    }
  }
}
