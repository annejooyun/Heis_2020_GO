# Network

Note: this module has a help-function module containing lower-level help-functions. This is done to try and improve code quality.

## Module description
This module is in charge of broadcasting and receiving messages between elevators using UDP. There are three types of messages that are sent:

### Status updates
Type: `Elev` (see elevator module)

Port: `2000`

All status updates are broadcasted to the port stated. This port is unique and only used for status updates. The network-module broadcasts all status-messages sent to the channel `ch_bcast_stat_update`. All status-messages that are broadcasted to the port stated, are sent to the channel `ch_rec_stat_update`.

### Orders
Type: `ExtOrder` (see orderDistributer module)

Port: `20203`

All orders are broadcasted to the port stated. This port is unique and only used for order-messages. The network module broadcasts all order-messages sent to the channel `ch_bcast_order`. All order-messages that are broadcasted to the port stated, are sent to the channel `ch_rec_ext_order`.

### Order Updates
Type: `int`

Port: `20194`

All order updates are broadcasted to the port stated. This port is unique and only used for order updates. In this project, an order update states at which floor an order has been executed. Whenever an order has been executed, the corresponding floor is sent to the network-module as an `int` via the channel `ch_bcast_order_exec`. All order updates sent to this channel are directly broadcasted to the port stated.

All messages sent to the port stated, are directly sent to the orderDistributer via the channel `ch_rec_order_exec`.




For more specific information, see the readme in network-helpfunc.
