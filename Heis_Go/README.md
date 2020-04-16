# ReadMe for elevator project spring 2020

## Using the simulator

To run a new elevator type in the following:

Choose a 5 digit connection port, i.e. 12345.

In the terminal running the program:
go run -ldflags="-X main.TCP_ConnectionPort=12345" main.go

In terminal running the simulator:

./SimElevatorServer --port 12345


## Project description
This project aims to control a system of n elevators operating at m floors. The object is that the elevators should be able to distribute orders amongst themselves, such that the system ideally performs better than with just one single elevator. The system should be fault-tolerant, meaning that all orders should be taken, even when errors appear.

In this specific software, we run 2 elevators, both operating at 4 floors. Although, the system is scalable such that it is possible to alter both the number of elevators and number of floors. We assume that at least one elevator works at each time.

The software uses a “momentary master” procedure, where all elevators has the responsibility of distributing orders registered by themselves. To do this, all elevators must have the knowledge of the current position, direction and order list of all other elevators at all time. In addition, all elevators must know which orders are active, that is which orders that already have been assigned and are to be taken. To achieve this, we utilize UDP broadcasting. All orders that are registered are broadcasted, and the status of the elevator is broadcasted each tie there is a change of status. These “messages” are broadcasted to specific ports assigned to each broadcast type.

To make sure the system is fault-tolerant, all elevators has a list of timestamps, that are activated/deactivated whenever an order is assigned/executed. If an elevator detects that the time since a timestamp exceeds a certain amount, the corresponding order is taken by the local elevator. In that way, we make sure that all orders are taken, however, we may observe that some orders are taken by several elevators. 


## Communication sequences
### Status update
For the system to be able to correctly distribute orders, all elevators must know the status of all other elevators. The relevant information is each elevators current position, current direction and their order list.

A status update is to be sent every time there has happend a status update. That is, every time the elevator reaches a new floor or changes direction. When this happens, the FSM module sends True over the channel ch_status_updated, by that telling the Elevator object module that there has been a change of state. The elevator object module then sends a copy of the elevator object to the order distributer, using the channel ch_internal_status_update.

The order distributer registers the status update correctly by placing the elevator object in the list ELEVATOR_STATUS_LIST and the corresponding ID on the same place (same index) in the list ADDED_ELEVATORS. When the status update has been correctly registerd, it is sent over the channel ch_bcast_elev_Status, to be broadcasted.

All status updates are broadcasted to the same port (20000). The function StartSendingAndReceivingStatusUpdates starts two goroutines which at all time broadcasts all status updates sent on the channel ch_bcast_elev_Status, and sends all messages received on port 20000 on the channel ch_receive_elev_status. When an status update from another elevator has been received, the status lists are updated as described earlier.

### Order registered


### Order executed


## Data types
