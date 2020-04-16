package orderDistributerHF

import(
  "../elevator"
  "../elevio"
  "time"
  "fmt"
)

type ExtOrder struct {
  Id string
  Floor int
  Button elevio.ButtonType
}


const TIMEOUT_LIMIT = 40

const N_ELEVATORS = 3
const TRAVEL_TIME = 2
const DOOR_OPEN_TIME = 3


var ELEVATOR_STATUS_LIST [N_ELEVATORS] elevator.Elev //Contains the previous registered states of all connected elevators
var ADDED_ELEVATORS [N_ELEVATORS] string //Contains the Id's of added elevators


var ACTIVE_ORDERS [elevator.N_FLOORS][2] int //If an order is beeing executed the element is set to 1

var TIMER_ACTIVE_ORDERS [elevator.N_FLOORS][2] int64 //List of timestamps for orders


//---------------------------------------------------------------------------------------------------------------------------//



//Finds the best elevator to execute an order
func BestChoice(order elevio.ButtonEvent) string {
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


func AlreadyActiveOrder(order elevio.ButtonEvent) bool {
  return ACTIVE_ORDERS[order.Floor][order.Button] == 1
}


func SetOrderActive(order elevio.ButtonEvent, active bool){
  if active{
    ACTIVE_ORDERS[order.Floor][order.Button] = 1
    TIMER_ACTIVE_ORDERS[order.Floor][order.Button] = time.Now().Unix()
    elevio.SetButtonLamp(order.Button, order.Floor, true)
  } else{
    ACTIVE_ORDERS[order.Floor][order.Button] = 0
    TIMER_ACTIVE_ORDERS[order.Floor][order.Button] = 0
    elevio.SetButtonLamp(order.Button, order.Floor, false)
  }
}


func RemoveOrdersOnFloor(floor int) {
  for index,_ := range ACTIVE_ORDERS[floor] {
    ACTIVE_ORDERS[floor][index] = 0
    TIMER_ACTIVE_ORDERS[floor][index] = 0
  }
  elevio.SetButtonLamp(elevio.BT_HallUp, floor, false)
  elevio.SetButtonLamp(elevio.BT_HallDown, floor, false)
}


//convert from type internalOrder (elevio.ButtonEvent) to type ExtOrder
func ConvertToExternalOrder(order elevio.ButtonEvent, owner string) ExtOrder{
  var externalOrder ExtOrder
  externalOrder.Id = owner
  externalOrder.Floor = order.Floor
  externalOrder.Button = order.Button
  return externalOrder
}


//Convert from type ExtOrder to internalOrder (elevio.ButtonEvent)
func ConvertToInternalOrder(order ExtOrder) elevio.ButtonEvent {
  var internalOrder elevio.ButtonEvent
  internalOrder.Floor = order.Floor
  internalOrder.Button = order.Button

  return internalOrder
}


func UpdateElevatorStatusList(elevator_status elevator.Elev) {
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


//---------------------------------------------------------------------------------------------------------------------------//


//Finds the absolute value of an int
func absInt(value int) int{
  if value < 0 {
    return -value
  } else {
    return value
  }
}


//Counts the number of orders in an elevators order list
func numOrders(elev elevator.Elev) int {
  nOrders := 0
  for _,row := range elev.OrderList{
    for _,element := range row {
      nOrders += element
    }
  }
  return nOrders
}


//Counts the number of floors the elevator must pass (worst case) to get to a specified order.
func numFloors(elev elevator.Elev, order elevio.ButtonEvent) int{
  nFloors := 0
  maxFloor := elevator.N_FLOORS - 1

  currentFloor := elev.Floor
  destinationFloor := order.Floor

  //Convert direction and button types to comparable ints
  //Set default values for direction
  currentDir := 1
  destinationDir := 1

  switch elev.CurrDirection {
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
func simpleCostFunction(elev elevator.Elev, order elevio.ButtonEvent) int {
  cost := numOrders(elev) * DOOR_OPEN_TIME
  if elev.CurrDirection == elevio.MD_Stop{
    cost += absInt(order.Floor - elev.Floor) * TRAVEL_TIME
  } else {
    cost += numFloors(elev, order) * TRAVEL_TIME
  }
  return cost
}


//Checks if timeout has happened
func IsOrderTimeout(timestamp int64) bool {
  var currTime int64
  currTime = time.Now().Unix()
  if timestamp != 0 && timestamp + TIMEOUT_LIMIT < currTime{
    return true
  }else{
    return false
  }
}
