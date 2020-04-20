# Order Handler
Note: this module has a help-function module containing lower-level help-functions. This is done to try and improve code quality.

## Module description
The orderHandler-module distributes, executes and removes all orders created by, or assigned to the local elevator. It also manages the list of orders in the local elevator (which exists in the elevator object, see elevator-module). 
There are two main routines in the orderHandler: `DistributeInternalOrders()` and `RegisterExecutedOrders()`. These are more closely explained under for the keen reader.

### DistributeInternalOrders()

#### Order_from_fsm
Whenever a button is pushed on the local elevator, the order is sent from the stateMachine to the orderHandler. The orderHandler then decides wether the order must be taken directly (i.e. it is a cab order) or needs to be further examined in order to be distributed correctly (hall orders).

If the order can be taken directly, the order is added to the local order list and the correct lights are turned on. On the other hand, if the order is to be distributed it is sent to the orderDistributer-module. This can be followed by the `orderDistributer.DistributeOrders()`-routine. 

#### Order_from_order_distributer
If an order that was registered by another elevator has been assigned to the local elevator, it is sent from the orderDistributer to the orderHandler. The orderHandler then adds the order to the order list and turns on the correct lights.


When the order has been properly registered by the local elevator, it is sent on the channel `ch_new_order`, thus triggering the stateMachine to do something. This can be read about in `stateMachine.RunStateMachine()`.

### RegisterExecutedOrders()

If an order has been executed, we send the floor on the `ch_internal_order_executed`. 

Lastly, we make sure to clear the internal orders at the floor on the local elevator.
