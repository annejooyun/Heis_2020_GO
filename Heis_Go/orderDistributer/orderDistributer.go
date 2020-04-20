package orderDistributer

import (
  "../elevator"
	"../elevio"
  "../orderDistributer-helpfunc"

  "fmt"
  "time"
)




//Main function controling where to send orders
func DistributeOrders(order_to_distribute chan elevio.ButtonEvent,
                      order_to_execute chan elevio.ButtonEvent,
                      broadcast_order chan orderDistributerHF.ExtOrder,
                      order_from_timeout chan elevio.ButtonEvent) {

  for {
    select{
    case orderReceived := <- order_to_distribute:

      if !orderDistributerHF.AlreadyActiveOrder(orderReceived) {
        owner := orderDistributerHF.BestChoice(orderReceived)
        externalOrder := orderDistributerHF.ConvertToExternalOrder(orderReceived,owner)
        broadcast_order <- externalOrder
      }

    case orderTimeout:= <- order_from_timeout:
      order_to_execute <- orderTimeout
    }
  }
}


func ReceiveOrders(receive_external_order chan orderDistributerHF.ExtOrder, order_to_execute chan elevio.ButtonEvent){
  for {
    select {
    case extOrderReceived := <- receive_external_order:

      internalOrder := orderDistributerHF.ConvertToInternalOrder(extOrderReceived)
      orderDistributerHF.SetOrderActive(internalOrder, true)

      fmt.Printf("Active orders table is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)

      if extOrderReceived.Id == elevator.LOCAL_ELEV_ID {
        order_to_execute <- internalOrder
        fmt.Printf("Executing external order on floor: %d\n",internalOrder.Floor)
      }
    }
  }
}


func RegisterExecutedOrders(receive_orders_executed chan int, internal_order_executed chan int, bcast_order_executed chan int){
  for {
    select{
    case floor := <- receive_orders_executed: //ordersExecuted = floor
      orderDistributerHF.RemoveOrdersInActiveOrders(floor)


      case floor := <- internal_order_executed: //ordersExecuted = [order up, order down,floor]
        bcast_order_executed <- floor
    }
  }
}


func PollStatusUpdates(internal_status_update chan elevator.Elev,
                       external_status_update chan elevator.Elev,
                       broadcast_status_update chan elevator.Elev) {

  for {
    select {

    case elevator_status := <- internal_status_update:
      broadcast_status_update <- elevator_status
      orderDistributerHF.UpdateElevatorStatusList(elevator_status)

    case elevator_status := <- external_status_update:
      orderDistributerHF.UpdateElevatorStatusList(elevator_status)
    }
  }
}


func PollOrderTimeout(order_timeout chan elevio.ButtonEvent) {

  for{

    time.Sleep(orderDistributerHF.POLL_RATE*time.Millisecond)

    for floor, orderlist := range orderDistributerHF.TIMER_ACTIVE_ORDERS {
      for button, timestamp := range orderlist {
        if orderDistributerHF.IsOrderTimeout(timestamp) {
          order := elevio.CreateButtonEvent(floor,elevio.IntToButtonType(button))
          order_timeout <- order
          orderDistributerHF.SetOrderActive(order,false)
        }
      }
    }
  }
}
