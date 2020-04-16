package orderDistributer

import (
  "../elevator"
	"../elevio"
  //"../network-helpfunc/bcast"
  "../orderDistributer-helpfunc"
  //"../stateMachine-go"
  //"../orderHandler"
  "time"
  "fmt"
)



func PollStatusUpdates(internal_status_update chan elevator.Elev, external_status_update chan elevator.Elev, broadcast_status_update chan elevator.Elev) {
  for {
    select {
    case elevator_status := <- internal_status_update://elevator_status is an elevator object
      broadcast_status_update <- elevator_status
      orderDistributerHF.UpdateElevatorStatusList(elevator_status)

    case elevator_status := <- external_status_update:
      orderDistributerHF.UpdateElevatorStatusList(elevator_status)
    }
  }
}

func PollOrderTimeout(order_timeout chan elevio.ButtonEvent) {
  var currTime int64
  for{
    currTime = time.Now().Unix()
    for floor, orderlist := range orderDistributerHF.TIMER_ACTIVE_ORDERS {
      for button, timestamp := range orderlist {
        if timestamp != 0 && timestamp + orderDistributerHF.TIMOUT_LIMIT < currTime {
          var order elevio.ButtonEvent
          order.Floor = floor
          order.Button = elevio.IntToButtonType(button)
          order_timeout <- order
          orderDistributerHF.SetOrderActive(order,false)
        }
      }
    }
  }
}

//Main function controling where to send orders
func DistributeOrders(order_to_distribute chan elevio.ButtonEvent, receive_external_order chan orderDistributerHF.ExtOrder, order_to_execute chan elevio.ButtonEvent, broadcast_order chan orderDistributerHF.ExtOrder, order_from_timeout chan elevio.ButtonEvent) {
  for {
    select{
    case orderReceived := <- order_to_distribute:

      if !orderDistributerHF.AlreadyActiveOrder(orderReceived) {
        owner := orderDistributerHF.BestChoice(orderReceived)
        //create new external order
        externalOrder := orderDistributerHF.ConvertToExternalOrder(orderReceived,owner)
        //broadcast external order
        broadcast_order <- externalOrder
        fmt.Printf("I sent an order on floor: %d\n", externalOrder.Floor)
      }

    case extOrderReceived := <- receive_external_order:

      fmt.Printf("Received order at floor: %d\n",extOrderReceived.Floor)

      internalOrder := orderDistributerHF.ConvertToInternalOrder(extOrderReceived)
      orderDistributerHF.SetOrderActive(internalOrder, true)
      fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
      if extOrderReceived.Id == elevator.LOCAL_ELEV_ID {
        order_to_execute <- internalOrder
        fmt.Printf("Executing external order on floor: %d\n",internalOrder.Floor)
      }

    case orderTimeout:= <- order_from_timeout:
      order_to_execute <- orderTimeout
      fmt.Printf("order timout\n")
    }
  }
}

/*
func NewDistributeOrders(elev *elevator.Elev, order_from_fsm chan elevio.ButtonEvent, new_order chan elevio.ButtonEvent, receive_external_order chan orderDistributerHF.ExtOrder, broadcast_order chan orderDistributerHF.ExtOrder, order_from_timeout chan elevio.ButtonEvent) {
  for {
    select{
    case order := <- order_from_fsm:
      if order.Button == elevio.BT_Cab {
        orderHandler.TakeOrder(elev,order)
        new_order <- order
      } else {
        if !orderDistributerHF.AlreadyActiveOrder(order) {
          //find who should execute the orders
          owner := orderDistributerHF.BestChoice(order)
          //create new external order
          externalOrder := orderDistributerHF.ConvertToExternalOrder(order,owner)
          //broadcast external order
          broadcast_order <- externalOrder
          fmt.Printf("I sent an order on floor: %d\n", externalOrder.Floor)
        }
      }

    case order := <- receive_external_order:

      fmt.Printf("Received order at floor: %d\n",order.Floor)

      internalOrder := orderDistributerHF.ConvertToInternalOrder(order)
      orderDistributerHF.SetOrderActive(internalOrder, true)
      //fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
      if order.Id == elevator.LOCAL_ELEV_ID {
        orderHandler.TakeOrder(elev, internalOrder)
        new_order <- internalOrder
        fmt.Printf("Executing external order on floor: %d\n",internalOrder.Floor)
      }
    case order:= <- order_from_timeout:
      orderHandler.TakeOrder(elev, order)
      new_order <- order
      fmt.Printf("order timout\n")
    }
  }
}
*/

func RegisterExecutedOrders(receive_orders_executed chan int, internal_order_executed chan int, bcast_order_executed chan int){
  for {
    select{
    case floor := <- receive_orders_executed: //ordersExecuted = floor
      orderDistributerHF.RemoveOrdersOnFloor(floor)

      //fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)

      case floor := <- internal_order_executed: //ordersExecuted = [order up, order down,floor]
        //fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
        orderDistributerHF.RemoveOrdersOnFloor(floor)
        bcast_order_executed <- floor
    }
  }
}





/*
func RegisterExecutedOrders(receive_orders_executed chan []int, internal_order_executed chan []int){
  for {
    select{
    case ordersExecuted := <- receive_orders_executed: //ordersExecuted = [order up, order down, floor]
      //fmt.Printf("ch_listen_order_executed_at_floor empty\n\n")

      floor := ordersExecuted[2]
      orderDistributerHF.RemoveOrdersOnFloor(floor)
      //fmt.Printf("These orders are executed at floor %d: %v\n\n",floor,ordersExecuted)
      //fmt.Printf("These orders are now registered at floor %d: %v\n\n", floor, orderDistributerHF.ACTIVE_ORDERS[floor])
      fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
      for index,element := range orderDistributerHF.ACTIVE_ORDERS[floor] {
        if ordersExecuted[index] == 1 && element == 1 {
          orderDistributerHF.ACTIVE_ORDERS[floor][index] = 0
          fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
        }
      }
      //fmt.Printf("Active orders is updated: %v", orderDistributerHF.ACTIVE_ORDERS)

      case ordersExecuted := <- internal_order_executed: //ordersExecuted = [order up, order down,floor]

        //fmt.Printf("ch_order_executed_by_me empty\n\n")

        fmt.Printf("I just received : %v\n", ordersExecuted)
        floor := ordersExecuted[2]
        fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
        orderDistributerHF.RemoveOrdersOnFloor(floor)

        for index,element := range orderDistributerHF.ACTIVE_ORDERS[floor] {
          if ordersExecuted[index] == 1 && element == 1 {
            orderDistributerHF.ACTIVE_ORDERS[floor][index] = 0
            fmt.Printf("Active orders is updated: %v", orderDistributerHF.ACTIVE_ORDERS)
          }
        }
    }
  }
}*/
