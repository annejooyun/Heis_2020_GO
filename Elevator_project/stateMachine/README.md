# State Machine

Note: this module has a help-function module containing lower-level help-functions. This is done to try and improve code quality.

## Module description:
The stateMachine-module is in charge of running one elevator according to the description. It has four triggers: New order, Button pushed, Floor detected, and Time-out.

**New order/Order received:**

An order is received from the orderHandler on the channel `new_order`. If the elevator is in the Idle state, it executes the order immediately. If not, it will be taken care of in the Floor detected case.

**Button pushed:**

When a button is pushed on the elevator, a message with the respective order is sent on the channel `order_registered`. The order handler then decides if it should be taken by its own elevator or not. This can be viewed in `orderHandler.DistributeInternalOrders()`.

**Floor detected:**

When the floor sensor detects a new floor, a message with the floor number is sent on the channel `drv_floors`. The state machine saves the previous floor in the class variable `Floor`, and then updates the elevator object to contain the floor detected. The floor indicator light on the current floor is also turned on.

If we have arrived at a new floor and we are meant to stop on the floor (i.e. there are order(s) that haven't been executed on the floor), we stop the elevator and execute the order(s). Since we now are at a new floor, the elevator object has now changed, and we send a message on the status update channel.


**Time-out:**

If we enter this case, it means that the elevator door has been open for >3s and it is time to close it. We turn off the light, and decides which direction the elevator should proceed with.




