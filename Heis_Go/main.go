package main

import (
  "./elevio"
  "./elevator"
  "./orderHandler"
  "./stateMachine"
  "./network"
  //"./timer"
  //"fmt"
  //"./Network-go/network/bcast"
  //"./messageHandler"
  "./orderDistributer"
  "./orderDistributer-helpfunc"
)

var TCP_ConnectionPort = "15657"

func main(){

    //CREATING ELEVATOR OBJECT
    elev := elevator.InitializeElevator(TCP_ConnectionPort)
    //Declaring what ports to UpdateCurrState


    //CREATING CHANNELS
    //Internal elevator channels

    //Whenever a new order is registered by the order handler, the order is sent to the state machine
    ch_new_order := make(chan elevio.ButtonEvent)

    //All orders detected by the local elevator is sent to the order handler
    ch_order_registered := make(chan elevio.ButtonEvent)

    //Message from fsm that a order has been executed at current floor
    ch_order_executed := make(chan bool)

    //Whenever a message is sent on this channel, a status update is sent to order distributer
    ch_stat_updated := make(chan bool)


    //Communication between elevator (control) and order distributer

    //Channel for sending status updates
    ch_int_stat_update := make(chan elevator.Elev)


    //Communication between order handler and order distributer

    //Orders delegated to local elevator is sent to order handler
    ch_order_to_exec := make(chan elevio.ButtonEvent)

    //Orders to be delegated is sent from order handler to distributer
    ch_order_to_distribute := make(chan elevio.ButtonEvent)


    //Communication between order distributer and message handler
    ch_bcast_stat_update := make(chan elevator.Elev)
    ch_ext_stat_update := make(chan elevator.Elev)

    ch_bcast_order := make(chan orderDistributerHF.ExtOrder)
    ch_rec_ext_order := make(chan orderDistributerHF.ExtOrder)

/*
    ch_bcast_order_exec := make(chan []int)
    ch_rec_order_exec := make(chan []int)

    ch_int_order_exec := make(chan []int)*/

    ch_bcast_order_exec := make(chan int)
    ch_rec_order_exec := make(chan int,2)

    ch_int_order_exec := make(chan int)

    ch_order_timeout := make(chan elevio.ButtonEvent)



    //STARTING GOROUTINES
    go elevator.PollInternalElevatorStatus(&elev, ch_stat_updated, ch_int_stat_update) //Whenever the status of the local elevator is updated, send an elevator copy to the order distributer

    go orderHandler.DistributeInternalOrders(&elev,ch_order_registered, ch_order_to_exec, ch_order_to_distribute, ch_new_order)

    go orderHandler.RegisterExecutedOrders(&elev,ch_order_executed, ch_int_order_exec)

    go orderDistributer.PollStatusUpdates(ch_int_stat_update, ch_ext_stat_update, ch_bcast_stat_update)

    go orderDistributer.DistributeOrders(ch_order_to_distribute, ch_order_to_exec, ch_bcast_order, ch_order_timeout)

    go orderDistributer.ReceiveOrders(ch_rec_ext_order, ch_order_to_exec)

    //go orderDistributer.NewDistributeOrders(&elev,ch_order_registered, ch_new_order, ch_rec_ext_order, ch_bcast_order, ch_order_timeout)

    go orderDistributer.RegisterExecutedOrders(ch_rec_order_exec, ch_int_order_exec, ch_bcast_order_exec)

    go orderDistributer.PollOrderTimeout(ch_order_timeout)


    //BEGIN SENDING AND RECEIVING BROADCASTS
    network.InitializeStatusUpdates(ch_bcast_stat_update, ch_ext_stat_update)

    network.InitializeOrders(ch_bcast_order, ch_rec_ext_order)

    network.InitializeOrderUpdates(ch_bcast_order_exec, ch_rec_order_exec)


    //START THE ELEVATOR STATE MACHINE
    stateMachine.RunStateMachine(&elev, ch_order_registered, ch_new_order, ch_stat_updated, ch_order_executed)


}
