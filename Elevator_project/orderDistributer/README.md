
# orderDistributer

Note: this module has a help-function module containing lower-level help functions. This is done to try and improve qode quality.

## Module description:
The order distributer is in charge of distributing orders that are registered locally, as well as handling orders that has been distributed by other elevators. It also introduces some fault-tolerance.

### Distributing orders
Whenever the order handler receives a hall order from the fsm, it is sent directly to the order distributer. Here, the order is assigned to one of the elevators in our system, using the cost function "*BestChoice*" seen in orderDistributer-helpfunc. This function iterates through the elevators we are connected with, and checks who is best fitted to execute the order. This is done utilizing a cost function.

When an owner of the order is found, we send a message to the network module, for it to broadcast. This message is of type ExtOrder, which is a struct declared in orderDistributer-helpfunc. It is on the following form:

```
type ExtOrder struct {
  Id string
  Floor int
  Button elevio.ButtonType
}
```

Whenever a message is broadcasted by any elevator, it is detected by the network module, which sends the message to the order distributer. If the owner ID of the order is the local ID, which means that I am to execute the order, the order is converted from an external order (ExtOrder) to an internal order (ButtonEvent), and is sent to the order handler.



### Status update

To be able to distribute orders, this module must have information about the status of the other elevators. To keep track of this information we have made two lists, that can be found in the help-function module. These are:

**ELEVATOR_STATUS_LIST:**

|copy elevator 1 | copy elevator 2 | ... |
|----------------|-----------------|-----|


**ADDED_ELEVATORS:**
|ID elevator 1 | ID elevator 2 | ... |
|--------------|---------------|-----|


#### Updating the lists

Each time the status of an elevator changes, a copy of the local elevator object is sent to the order distributer. Firstly, this status update is broadcasted, so that all the other elevators are aware of the new status. Secondly, the order distributer scans through the list ADDED_ELEVATORS, to check wether or not the new status update has already been added to the list. If it has not, the ID is added in the first empty spot and the copy of the elevator object is placed in the list ELEVATOR_STATUS_LIST with the same index. If the ID already lies in the ADDED_ELEVATOR list, it merely replaces the object that is on the same index in the ELEVATOR_STATUS_LIST with the new, updated copy. That way, we make sure that the ID of the object in place i in the ELEVATOR_STATUS_LIST, lies in place i in the ADDED_ELEVATORS list.

When a status update is registered by the network module, it is sent to the order distributer. Then the status update is added to the two lists as described above.

## Fault tolerance - Keeping track on active orders
In order to introduce fault tolerance to the system, we will at all time keep track of all orders that are active. That is, we want to know each time an order is being taken, and if the order has been executed.

To keep track of this, we have constructed the following list:

**ACTIVE_ORDERS:**

|floor\Button type | BT_HallUp | BT_HallDown | 
|       :----:     |   :----:  |   :----:    |
|       0          |    0/1    |      0/1    |
|       1          |    0/1    |      0/1    |
|       2          |    0/1    |      0/1    |
|       3          |    0/1    |      0/1    |


Each time an order is distributed by any elevator, we set the corrseponding element in the list ACTIVE_ORDERS to 1. When the order is executed, it is set to 0.

### Registering executed orders:
In this project we have made the assumption that if an elevator stops at a floor, all orders are taken at that floor. That is, when any elevator executes an order on a floor, the elements in the ACTIVE_ORDER list corresponding to that floor, shall be set to 0 in all elevators.

Whenever an order has been executed by the local elevator, a message is sent from the order handler to the order distributer via the channel ch_int_order_exec, telling which floor the order has been executed at. This message is directly sent to the network module where it is distributed.

Every time the network module registers that an order-executed message has been broadcasted, it sends it to the order distributer via the channel ch_rec_order_exec. When this happens, the order distributer sets all elements in the ACTIVE_ORDERS list, corresponding to the floor where the order has been executed, to 0.

### Timestamps
To be able to act when something wrong has happend, we must know how long it has been since an order was activated. To do this, we implement another list called TIMER_ACTIVE_ORDERS where each element is 0 when there are no active orders, and a Unix timestamp from the time the order was activated if when an order is active.

**TIMER_ACTIVE_ORDERS:**

|floor\Button type | BT_HallUp | BT_HallDown | 
|       :----:     |   :----:  |   :----:    |
|       0          |0/timestamp|0/timestamp  |
|       1          |0/timestamp|0/timestamp  |
|       2          |0/timestamp|0/timestamp  |
|       3          |0/timestamp|0/timestamp  |


The function *PollOrderTimeout* constantly checks wether an active order has reached a certain time limit. In this implementation this limit is set to 40 seconds. If an order exceeds the time limit, the corresponding order is reasigned to the local elevator.


