# State Machine

Note: this module has a help-function module containing lower-level help functions. This is done to try and improve qode quality.

### Module description:
The stateMachine module is in charge of running one elevator according to the description. It has four triggers: New order, Button pushed, Floor detected, and Time-out.

**New order/Order received:**

An order is received from the orderHandler on the channel "incoming_orders". We should then add the order to our execution list, and set the light in the correct order button. If we are in Idle, we should execute the order immediately.


An order is received from the orderHandler on the channel "new_order". If the elevator is in the Idle state, it should execute the order immediately. If not, it will be taken care of in the Floor detected case.

**Button pushed:**

When a button is pushed on the elevator, a message with the respective order is sent on the channel "orders_to_order_handler". The order handler then distributes the order to the right elevator.


When a button is pushed on the elevator, a message with the respective order is sent on the channel "order_registered". The order handler then distributes the order to the right elevator.


**Floor detected:**

When the floor sensor detects a new floor, a message with the floor number is sent on the channel "drv_floors". The state machine then saves the previous floor in a variable and then updates the elevator object to contain the floor detected. The floor indicator light is also turned on on the floor detected.

If we have arrived at a new floor we shall do the following:
If we are meant to stop on the floor, we stop the elevator and executes the orders.
Since we now are at a new floor, the elevator object has now changed, and we send a message on the status update channel.


When the floor sensor detects a new floor, a message with the floor number is sent on the channel "drv_floors". The state machine then saves the previous floor in the class variable Floor, and then updates the elevator object to contain the floor detected. The floor indicator light on the current floor is also turned on.

If we have arrived at a new floor and we are meant to stop on the floor (i.e. there are order(s) that haven't been executed on the floor), we stop the elevator and execute the order(s).
Since we now are at a new floor, the elevator object has now changed, and we send a message on the status update channel.


**Time-out**

If we enter this case, it means that the elevator door has been open for >3s and it is time to close it. We turn off the light, 




