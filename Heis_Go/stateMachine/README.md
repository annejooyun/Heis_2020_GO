# State Machine

### Module description:
The stateMachine module is in charge of running one elevator according to the description. It has four triggers: Button pushed, floor detected, order received and timeOut.

**Button pushed:**

When a button is pushed on the elevator, a message with the respective order is sent on the channel "orders_to_order_handler". The order handler then distributes the order to the right elevator.

**Floor detected:**

When the floor sensor detects a new floor, a message with the floor number is sent on the channel "drv_floors". The state machine then saves the previous floor in a variable and then updates the elevator object to contain the floor detected. The floor indicator light is also turned on on the floor detected.

If we have arrived at a new floor we shall do the following:
If we are meant to stop on the floor, we stop the elevator and executes the orders.
Since we now are at a new floor, the elevator object has now changed, and we send a message on the status update channel.

**Order received:**

An order is received from the order handler on the channel "incoming_orders". We should then add the order to our execution list, and set the light in the correct order button. If we are in Idle, we should execute the order immediately.
