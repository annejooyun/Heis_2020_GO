package orderDistributer

import (
  "../control-go"
	"../elevio"
  //"../stateMachine-go"
  //"../orderHandler"
  "../Network-go/network/bcast"
  "fmt"
)

type ExtOrder struct {
  Id string
  Floor int
  Button elevio.ButtonType
}

const N_ELEVATORS = 3
const TRAVEL_TIME = 2
const DOOR_OPEN_TIME = 3

var ELEVATOR_STATUS_LIST [N_ELEVATORS] control.Elev //Contains the previous registered states of all connected elevators
var ADDED_ELEVATORS [N_ELEVATORS] string //Contains the Id's of added elevators


var ACTIVE_ORDERS [control.N_FLOORS][2] int //If an order is beeing executed the element is set to 1



func UpdateElevatorStatusList(elevator_status control.Elev) {
  for index,element := range(ADDED_ELEVATORS) { //The elevators have the same order in both lists
    if element == "" || element == elevator_status.Id {
      ADDED_ELEVATORS[index] = elevator_status.Id
      ELEVATOR_STATUS_LIST[index] = elevator_status

      //Printout to check what happens
      fmt.Printf("Status updated for %s\n", elevator_status.Id)
      fmt.Printf("The elevator list now contains the following elevators: %v\n",ADDED_ELEVATORS)
      break
    }
  }
}

func PollStatusUpdates(internal_status_update chan control.Elev, external_status_update chan control.Elev, broadcast_status_update chan control.Elev) {
  for {
    select {
    case elevator_status := <- internal_status_update://elevator_status is an elevator object
      broadcast_status_update <- elevator_status
      UpdateElevatorStatusList(elevator_status)

    case elevator_status := <- external_status_update:
      UpdateElevatorStatusList(elevator_status)
    }
  }
}

/*
//This was from the description given
func costFunction(elevator control.Elev, order elevio.ButtonEvent) int {
  test_elevator := elevator//To make sure we do not alter anything in the real elevator
  duration := 0
  switch test_elevator.CurrState {
  case control.Idle:
    test_elevator.CurrDirection = stateMachine.ChooseDirection(&test_elevator)
    if test_elevator.CurrDirection == elevio.MD_Stop {
      return duration
    }
    break

  case control.Moving:
    duration += TRAVEL_TIME/2
    switch test_elevator.CurrDirection {
    case elevio.MD_Up:
      test_elevator.Floor += 1
    case elevio.MD_Down:
      test_elevator.Floor -= 1
    }
    break

  case control.DoorOpen:
    duration += DOOR_OPEN_TIME/2
  }

  for {
    if stateMachine.ShouldIStop(&test_elevator) {
      orderHandler.ClearOrdersAtCurrentFloor(&test_elevator)
      duration += DOOR_OPEN_TIME
      test_elevator.CurrDirection = stateMachine.ChooseDirection(&test_elevator)
      if test_elevator.CurrDirection == elevio.MD_Stop {
        return duration
      }
    }
    switch test_elevator.CurrDirection {
    case elevio.MD_Up:
      test_elevator.Floor += 1
    case elevio.MD_Down:
      test_elevator.Floor -= 1
    }
    duration += TRAVEL_TIME
  }
  return duration
}
*/

//Finds the absolute value of an int
func Abs(value int) int{
  if value < 0 {
    return -value
  } else {
    return value
  }
}


//Counts the number of orders in an elevators order list
func NumOrders(elevator control.Elev) int {
  nOrders := 0
  for _,row := range elevator.OrderList{
    for _,element := range row {
      nOrders += element
    }
  }
  return nOrders
}


//Counts the number of floors the elevator must pass (worst case) to get to a specified order.
func NumFloors(elevator control.Elev, order elevio.ButtonEvent) int{
  nFloors := 0
  maxFloor := control.N_FLOORS - 1

  currentFloor := elevator.Floor
  destinationFloor := order.Floor

  //Convert direction and button types to comparable ints
  //Set default values for direction
  currentDir := 1
  destinationDir := 1

  switch elevator.CurrDirection {
  case elevio.MD_Up:
    currentDir = 1
  case elevio.MD_Down:
    currentDir = -1
  }

  switch order.Button {
  case elevio.BT_HallUp:
    destinationDir = 1
  case elevio.BT_HallDown:
    destinationDir = -1
  }
  for {
    if currentFloor == destinationFloor && currentDir == destinationDir {
      break
    } else {
      if currentFloor == maxFloor || currentFloor == 0 {
        currentDir = -currentDir
      }
      currentFloor += currentDir
      nFloors += 1
    }
  }
  return nFloors
}


//Have to update cost function
func simpleCostFunction(elevator control.Elev, order elevio.ButtonEvent) int {
  cost := NumOrders(elevator) * DOOR_OPEN_TIME
  if elevator.CurrDirection == elevio.MD_Stop{
    cost += Abs(order.Floor - elevator.Floor) * TRAVEL_TIME
  } else {
    cost += NumFloors(elevator, order) * TRAVEL_TIME
  }
  return cost
}


//Finds the best elevator to execute an order
func bestChoice(order elevio.ButtonEvent) string {
  best_choice := ""
  best_cost  := 10000
  for index,element := range(ELEVATOR_STATUS_LIST) {
    if ADDED_ELEVATORS[index] == "" {
      break
    } else if simpleCostFunction(element, order) < best_cost {
      best_cost = simpleCostFunction(element,order)
      best_choice = element.Id
    }
  }
  fmt.Printf("Best cost: %d\n",best_cost)
  return best_choice
}

func alreadyActiveOrder(order elevio.ButtonEvent) bool {
  return ACTIVE_ORDERS[order.Floor][order.Button] == 1
}

func setOrderActive(order elevio.ButtonEvent, active bool){
  if active{
    ACTIVE_ORDERS[order.Floor][order.Button] = 1
  } else{
    ACTIVE_ORDERS[order.Floor][order.Button] = 0
  }
}


//convert from type internalOrder (elevio.ButtonEvent) to type ExtOrder
func convertToExternalOrder(order elevio.ButtonEvent, owner string) ExtOrder{
  var externalOrder ExtOrder
  externalOrder.Id = owner
  externalOrder.Floor = order.Floor
  externalOrder.Button = order.Button
  return externalOrder
}


//Convert from type ExtOrder to internalOrder (elevio.ButtonEvent)
func convertToInternalOrder(order ExtOrder) elevio.ButtonEvent {
  var internalOrder elevio.ButtonEvent
  internalOrder.Floor = order.Floor
  internalOrder.Button = order.Button

  return internalOrder
}


func removeOrdersOnFloor(floor int) {
  for index,_ := range ACTIVE_ORDERS[floor] {
    ACTIVE_ORDERS[floor][index] = 0
  }
}
//Main function controling where to send orders
func DistributeOrders(order_to_distribute chan elevio.ButtonEvent, external_order chan ExtOrder, send_to_self chan elevio.ButtonEvent, broadcast_order chan ExtOrder) {
  for {
    select{
    case orderReceived := <- order_to_distribute:

      //fmt.Printf("ch_order_to_distribute empty\n\n")
      /*
      if orderReceived.Button == elevio.BT_Cab {
        send_to_self <- orderReceived

        //fmt.Printf("ch_order_to_execute full\n\n")

      } else {*/
        if !alreadyActiveOrder(orderReceived) {
          fmt.Printf("Active orders: %v\n", ACTIVE_ORDERS)
          setOrderActive(orderReceived,true)
          owner := bestChoice(orderReceived)
          fmt.Printf("Owner of order is %s\n",owner)
          if owner == control.LOCAL_ELEV_ID {
            send_to_self <- orderReceived

            //fmt.Printf("ch_order_to_execute full\n\n")

          }
          fmt.Printf("Active orders is updated: %v\n", ACTIVE_ORDERS)
          //create new external order
          externalOrder := convertToExternalOrder(orderReceived,owner)
          //broadcast external order
          broadcast_order <- externalOrder

          //fmt.Printf("ch_broadcast_order full\n\n")
          //fmt.Printf("Broadcasting order to elevator %s on floor %d \n",externalOrder.Id, externalOrder.Floor)
        }


      //}
    case extOrderReceived := <- external_order:

      //fmt.Printf("ch_receive_external_order empty\n\n")

      internalOrder := convertToInternalOrder(extOrderReceived)
      setOrderActive(internalOrder, true)
      fmt.Printf("Active orders is updated: %v", ACTIVE_ORDERS)
      if extOrderReceived.Id == control.LOCAL_ELEV_ID {
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
      //fmt.Printf("These orders are executed at floor %d: %v\n\n",floor,ordersExecuted)
      //fmt.Printf("These orders are now registered at floor %d: %v\n\n", floor, ACTIVE_ORDERS[floor])
      fmt.Printf("Active orders is updated: %v\n", ACTIVE_ORDERS)
      for index,element := range ACTIVE_ORDERS[floor] {
        if ordersExecuted[index] == 1 && element == 1 {
          ACTIVE_ORDERS[floor][index] = 0
          fmt.Printf("Active orders is updated: %v\n", ACTIVE_ORDERS)
        }
      }
      //fmt.Printf("Active orders is updated: %v", ACTIVE_ORDERS)

      case ordersExecuted := <- order_executed_by_me: //ordersExecuted = [order up, order down,floor]

        //fmt.Printf("ch_order_executed_by_me empty\n\n")

        fmt.Printf("I just received : %v\n", ordersExecuted)
        floor := ordersExecuted[2]
        fmt.Printf("Active orders is updated: %v\n", ACTIVE_ORDERS)
        //removeOrdersOnFloor(floor)

        for index,element := range ACTIVE_ORDERS[floor] {
          if ordersExecuted[index] == 1 && element == 1 {
            ACTIVE_ORDERS[floor][index] = 0
            fmt.Printf("Active orders is updated: %v", ACTIVE_ORDERS)
          }
        }
    }
  }
}



//Sending and receiving messages

const PORT_STATUS_UPDATES = 20000
const PORT_ORDERS = 20203
const PORT_ORDER_UPDATES = 20194

func StartSendingAndReceivingStatusUpdates(sendStatusUpdates chan control.Elev, receiveStatusUpdates chan control.Elev) {
  go bcast.Transmitter(PORT_STATUS_UPDATES, sendStatusUpdates)
  go bcast.Receiver(PORT_STATUS_UPDATES, receiveStatusUpdates)
}

func StartSendingAndReceivingOrders(broadcastMessage chan ExtOrder, receiveMessage chan ExtOrder) {
  go bcast.Transmitter(PORT_ORDERS, broadcastMessage)
  go bcast.Receiver(PORT_ORDERS, receiveMessage)
}

func StartSendingAndReceivingOrderUpdates(broadcastOrderExecuted chan []int, listenOrderExecuted chan []int) {
  go bcast.Transmitter(PORT_ORDER_UPDATES, broadcastOrderExecuted)
  go bcast.Receiver(PORT_ORDER_UPDATES, listenOrderExecuted)
}
