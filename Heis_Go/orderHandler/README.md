# Module description - Order Handler

The orderHandler module distributes all orders created by, or assigned to the local elevator. It also manages the list of orders in the local elevator (which exists in the elevator object, see module elevator).

Whenever a button is pushed on the local elevator, the order is sent to the order handler. The order handler then decides wether the order can be taken directly (cab order) or needs to be distributed (hall orders).

If the order can be taken directly, the order is added to the local order list and the correct lights are turned on. On the other hand, if the order is to be distributed it is sent to the orderDistributer module.

If an order has been assigned to the local elevator, it is sent from the order distributer to the order handler. The order handler then adds the order to the order list and turns on the correct lights.

If an order has been executed, the fsm will send a boolean signal (true) to the message handler, and by that saying: "The orders on the current floor are executed". The order handler then removes all orders on the current floor from the order list.

This module has the main functions "*DistributeInternalOrders*" and "*RegisterExecutedOrders*", which, as the name implies, distributes all internal orders, and register if an order has been executed.
