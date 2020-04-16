
# orderDistributer

Note: this module has a help-function module containing lower-level help functions. This is done to try and improve qode quality.

## Module description:
The order distributer is in charge of distributing orders that are registered locally, as well as handling orders that has been distributed by other elevators.

To be able to do this, this module must have information about the status of the other elevators and which orders that are already being executed. To keep track of this information we have made three lists, that can be found in the help-function module. These are:

**ELEVATOR_STATUS_LIST:**

|copy elevator 1 | copy elevator 2 | ... |
|----------------|-----------------|-----|


**ADDED_ELEVATORS:**
|ID elevator 1 | ID elevator 2 | ... |
|--------------|---------------|-----|


**ACTIVE_ORDERS:**

|floor\Button type | BT_HallUp | BT_HallDown | 
|       :----:     |   :----:  |   :----:    |
|       0          |    0/1    |      0/1    |
|       1          |    0/1    |      0/1    |
|       2          |    0/1    |      0/1    |
|       3          |    0/1    |      0/1    |


### Registering status update
Each time the status of an elevator changes, a copy of the local elevator object is sent to the order distributer. Firstly, this status update is broadcasted, so that all the other elevators are aware of the new status. Secondly, the order distributer scans through the list ADDED_ELEVATORS, to check wether or not the new status update has already been added to the list. If it has not, the ID is added in the first empty spot and the copy of the elevator object is placed in the list ELEVATOR_STATUS_LIST with the same index. If the ID already lies in the ADDED_ELEVATOR list, it merely replaces the object that is on the same index in the ELEVATOR_STATUS_LIST with the new, updated copy. That way, we make sure that the ID of the object in place i in the ELEVATOR_STATUS_LIST, lies in place i in the ADDED_ELEVATORS list.

When a status update is registered by the network module, it is sent to the order distributer. Then the status update is added to the two lists as described above.

### Keeping track on active orders
