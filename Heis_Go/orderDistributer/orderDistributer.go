package orderDistributer

import (
  "../control-go"
	"../elevio"
  "../stateMachine-go"
  "../orderHandler"
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


func funUpdateElevStat(elevator_status control.Elev) {
  for index,element := range(ADDED_ELEVATORS) { //The elevators have the same order in both lists
    if element == "" || element == elevator_status.Id {
      ADDED_ELEVATORS[index] = elevator_status.Id
      ELEVATOR_STATUS_LIST[index] = elevator_status

      //Printout to check what happens
      //fmt.Printf("Status updated for %s\n", elevator_status.Id)
      //fmt.Printf("The elevator list now contains the following elevators: %v\n",ADDED_ELEVATORS)
      break
    }
  }
}

func UpdateElevatorStatus(internal_status_update chan control.Elev, external_status_update chan control.Elev, broadcast_status_update chan control.Elev) {
  var test_order elevio.ButtonEvent
  test_order.Button = elevio.BT_HallDown
  test_order.Floor = 0
  for {
    select {
    case elevator_status := <- internal_status_update://elevator_status is an elevator object
      broadcast_status_update <- elevator_status
      funUpdateElevStat(elevator_status)
      cost := simpleCostFunction(elevator_status, test_order)
      fmt.Printf("The cost of the test order on elevator %s is: %d\n",elevator_status.Id,cost)

    case elevator_status := <- external_status_update:
      funUpdateElevStat(elevator_status)
      //cost := simpleCostFunction(elevator_status, test_order)
      //fmt.Printf("The cost of %s is: %d\n",elevator_status.Id,cost)
    }
  }
}

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


func Abs(value int) int{
  if value < 0 {
    return -value
  } else {
    return value
  }
}

func NumOrders(elevator control.Elev) int {
  nOrders := 0
  for _,row := range elevator.OrderList{
    for _,element := range row {
      nOrders += element
    }
  }
  return nOrders
}

//Have to update cost function
func simpleCostFunction(elevator control.Elev, order elevio.ButtonEvent) int {
  cost := NumOrders(elevator) * DOOR_OPEN_TIME
  switch elevator.CurrDirection {
  case elevio.MD_Stop:
    cost += Abs(order.Floor - elevator.Floor) * TRAVEL_TIME
  case elevio.MD_Up:
    cost += ((control.N_FLOORS - elevator.Floor) + (control.N_FLOORS - order.Floor)) * TRAVEL_TIME
  case elevio.MD_Down:
    cost += (order.Floor + elevator.Floor) * TRAVEL_TIME
  }
  return cost
}


func bestChoice(order elevio.ButtonEvent) string {
  best_choice := ""
  best_cost  := 100
  for index,element := range(ELEVATOR_STATUS_LIST) {
    if ADDED_ELEVATORS[index] == "" {
      break
    } else if simpleCostFunction(element, order) < best_cost {
      best_cost = simpleCostFunction(element,order)
      best_choice = element.Id
    }
  }
  return best_choice
}







//Sending and receiving messages


var PORT_STATUS_UPDATES = 20000
var PORT_MESSAGES = 20203

func StartSendingAndReceivingStatusUpdates(sendStatusUpdates chan control.Elev, receiveStatusUpdates chan control.Elev) {
  go bcast.Transmitter(PORT_STATUS_UPDATES, sendStatusUpdates)
  go bcast.Receiver(PORT_STATUS_UPDATES, receiveStatusUpdates)
}

func StartSendingAndReceivingOrders(broadcastMessage chan ExtOrder, receiveMessage chan ExtOrder) {
  go bcast.Transmitter(PORT_MESSAGES, broadcastMessage)
  go bcast.Receiver(PORT_MESSAGES, receiveMessage)
}




func DistributeOrders(order_to_distribute chan elevio.ButtonEvent, external_order chan ExtOrder, send_to_self chan elevio.ButtonEvent, broadcast_order chan ExtOrder) {
  for {
    select{
    case orderReceived := <- order_to_distribute:
      if orderReceived.Button == elevio.BT_Cab {
        send_to_self <- orderReceived
      } else {
        owner := bestChoice(orderReceived)
        fmt.Printf("Owner of order is %s\n",owner)
        if owner == control.LOCAL_ELEV_ID {
          send_to_self <- orderReceived
        } else {
          var externalOrder ExtOrder
          externalOrder.Id = owner
          externalOrder.Floor = orderReceived.Floor
          externalOrder.Button = orderReceived.Button
          broadcast_order <- externalOrder
          fmt.Printf("Broadcasting order to elevator %s on floor %d \n",externalOrder.Id, externalOrder.Floor)
        }
      }
    case extOrderReceived := <- external_order:
      if extOrderReceived.Id == control.LOCAL_ELEV_ID {
        var order elevio.ButtonEvent
        order.Floor = extOrderReceived.Floor
        order.Button = extOrderReceived.Button
        send_to_self <- order
        fmt.Printf("Executing external order on floor: %d\n",order.Floor)
      }
    }
  }
}
