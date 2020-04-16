
# orderDistributer

## Module description:
Insert text here

## Comunication between modules
![Alt text](overview_channels.png)

### Channel descriptions
#### ch_order_registered:
Type: ButtonEvent

Usage:

Whenever the FSM registeres that a button has been pushed, the order (ButtonEvent) corresponding to the pushed button is sent to the Order Handler using the channel ch_order_registered.

#### ch_order_executed:
Type: Bool

Usage:

Whenever the elevator stops on a floor, one or more orders are executed. The FSM then sends a True over the channel ch_order_executed, to the order handler, indicating that all orders on the floor the elevator is currently at, has been executed.

#### ch_status_changed
Type: Bool

Usage:

#### ch_internal_status_update
Type: Elev

Usage:

Whenever the Elevator object module is informed that the status has been changed, a copy of the elevator object is sent to the order distributer over the channel ch_internal_status_update


#### ch_order_to_execute
Type: ButtonEvent

Usage:

If the order distributer distributes an order to the local elevator, it is sent to the order handler via the channel ch_order_to_execute

#### ch_order_to_distribute
Type: ButtonEvent

Usage:

All hall-orders that the order handler receives from the fsm must be sent to the order distributer, to be distributed. Thus, all hall-orders from the FSM is sent to the order distributer via the channel ch_order_to_distribute.

#### ch_internal_order_executed
Type: Int[]

Format: [floor (int), Button Hall Up (1/0), Button Hall Down (1/0)]

Usage:

When an order is executed at a floor, all orders registered at that floor in the order handler are executed. The order distributer must be informed of which orders the elevator has executed, to be able to 

