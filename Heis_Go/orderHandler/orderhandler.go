package orderHandler

import (
  "../elevator"
  "../elevio"
  "../orderHandler-helpfunc"
)


func DistributeInternalOrders(elev *elevator.Elev,
                              order_from_fsm chan elevio.ButtonEvent,
                              order_from_order_distributer chan elevio.ButtonEvent,
                              distribute_order chan elevio.ButtonEvent,
                              new_order chan elevio.ButtonEvent) {

  for {
    select {

    //Order from its own elevator's state machine
    case order := <- order_from_fsm:
      if order.Button == elevio.BT_Cab {
        orderHandlerHF.TakeOrder(elev,order)
        new_order <- order
      } else {
        distribute_order <- order
      }

    //Order from another elevator
    case order := <- order_from_order_distributer:
      orderHandlerHF.TakeOrder(elev, order)
      new_order <- order
    }
  }
}


func RegisterExecutedOrders(elev *elevator.Elev, order_executed chan bool, internal_order_executed chan int){

  for {
    select {
      
    case <- order_executed:
      //create list of orders executed on the form [button up, button down, floor]
      internal_order_executed <- elev.Floor

      orderHandlerHF.ClearOrdersAtCurrentFloor(elev)
    }
  }
}
