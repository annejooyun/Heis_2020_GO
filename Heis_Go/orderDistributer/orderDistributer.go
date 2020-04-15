package orderDistributer

import (
  "../elevator"
	"../elevio"
  //"../network-helpfunc/bcast"
  "../orderDistributer-helpfunc"
  //"../stateMachine-go"
  //"../orderHandler"

  "fmt"
)



func PollStatusUpdates(internal_status_update chan elevator.Elev, external_status_update chan elevator.Elev, broadcast_status_update chan elevator.Elev) {
  for {
    select {
    case elevator_status := <- internal_status_update://elevator_status is an elevator object
      broadcast_status_update <- elevator_status
      orderDistributerHF.UpdateElevatorStatusList(orderDistributerHF.ADDED_ELEVATORS,orderDistributerHF.ELEVATOR_STATUS_LIST,elevator_status)

    case elevator_status := <- external_status_update:
      orderDistributerHF.UpdateElevatorStatusList(orderDistributerHF.ADDED_ELEVATORS,orderDistributerHF.ELEVATOR_STATUS_LIST,elevator_status)
    }
  }
}


//Main function controling where to send orders
func DistributeOrders(order_to_distribute chan elevio.ButtonEvent, external_order chan orderDistributerHF.ExtOrder, send_to_self chan elevio.ButtonEvent, broadcast_order chan orderDistributerHF.ExtOrder) {
  for {
    select{
    case orderReceived := <- order_to_distribute:

      //fmt.Printf("ch_order_to_distribute empty\n\n")
      /*
      if orderReceived.Button == elevio.BT_Cab {
        send_to_self <- orderReceived

        //fmt.Printf("ch_order_to_execute full\n\n")

      } else {*/
        if !orderDistributerHF.AlreadyActiveOrder(orderReceived) {
          fmt.Printf("Active orders: %v\n", orderDistributerHF.ACTIVE_ORDERS)
          orderDistributerHF.SetOrderActive(orderReceived,true)
          owner := orderDistributerHF.BestChoice(orderReceived)
          fmt.Printf("Owner of order is %s\n",owner)
          if owner == elevator.LOCAL_ELEV_ID {
            send_to_self <- orderReceived

            //fmt.Printf("ch_order_to_execute full\n\n")

          }
          fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
          //create new external order
          externalOrder := orderDistributerHF.ConvertToExternalOrder(orderReceived,owner)
          //broadcast external order
          broadcast_order <- externalOrder

          //fmt.Printf("ch_broadcast_order full\n\n")
          //fmt.Printf("Broadcasting order to elevator %s on floor %d \n",externalOrder.Id, externalOrder.Floor)
        }


      //}
    case extOrderReceived := <- external_order:

      //fmt.Printf("ch_receive_external_order empty\n\n")

      internalOrder := orderDistributerHF.ConvertToInternalOrder(extOrderReceived)
      orderDistributerHF.SetOrderActive(internalOrder, true)
      fmt.Printf("Active orders is updated: %v", orderDistributerHF.ACTIVE_ORDERS)
      if extOrderReceived.Id == elevator.LOCAL_ELEV_ID {
        send_to_self <- internalOrder

        //fmt.Printf("ch_order_to_execute full\n\n")

        fmt.Printf("Executing external order on floor: %d\n",internalOrder.Floor)
      }
    }
  }
}


func RegisterExecutedOrders(listenExecutedOrders chan []int, order_executed_by_me chan []int){
  for {
    select{
    case ordersExecuted := <- listenExecutedOrders: //ordersExecuted = [order up, order down, floor]
      //fmt.Printf("ch_listen_order_executed_at_floor empty\n\n")

      floor := ordersExecuted[2]
      orderDistributerHF.RemoveOrdersOnFloor(floor)
      //fmt.Printf("These orders are executed at floor %d: %v\n\n",floor,ordersExecuted)
      //fmt.Printf("These orders are now registered at floor %d: %v\n\n", floor, orderDistributerHF.ACTIVE_ORDERS[floor])
      fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
      /*for index,element := range orderDistributerHF.ACTIVE_ORDERS[floor] {
        if ordersExecuted[index] == 1 && element == 1 {
          orderDistributerHF.ACTIVE_ORDERS[floor][index] = 0
          fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
        }
      }*/
      //fmt.Printf("Active orders is updated: %v", orderDistributerHF.ACTIVE_ORDERS)

      case ordersExecuted := <- order_executed_by_me: //ordersExecuted = [order up, order down,floor]

        //fmt.Printf("ch_order_executed_by_me empty\n\n")

        fmt.Printf("I just received : %v\n", ordersExecuted)
        floor := ordersExecuted[2]
        fmt.Printf("Active orders is updated: %v\n", orderDistributerHF.ACTIVE_ORDERS)
        orderDistributerHF.RemoveOrdersOnFloor(floor)

        /*for index,element := range orderDistributerHF.ACTIVE_ORDERS[floor] {
          if ordersExecuted[index] == 1 && element == 1 {
            orderDistributerHF.ACTIVE_ORDERS[floor][index] = 0
            fmt.Printf("Active orders is updated: %v", orderDistributerHF.ACTIVE_ORDERS)
          }
        }*/
    }
  }
}
