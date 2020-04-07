package main

import (
  "./elevio"
  "./control-go"
  "./orderHandler"
  "./stateMachine-go"
  //"./timer"
  //"fmt"
  //"./Network-go/network/bcast"
  //"./messageHandler"
  "./orderDistributer"
)

var TCP_ConnectionPort = "15657"

func main(){

    //Creating elevator
    elevator := control.InitializeElevator(TCP_ConnectionPort)
    //Declaring what ports to UpdateCurrState


    //Creating channels:

    //Internal elevator channels
    ch_orderHandler_fsm_new_order := make(chan elevio.ButtonEvent) //Whenever a new order is registered by the order handler, the order is sent to the fsm

    ch_fsm_orderHandler_order := make(chan elevio.ButtonEvent) //All orders detected by the local elevator is sent to the order handler

    ch_fsm_orderHandler_order_executed := make(chan bool) //Message from fsm that a order has been executed at current floor

    ch_fsm_control_status_updated := make(chan bool) //Whenever a message is sent on this channel, a status update is sent to order distributer


    //Communication between elevator (control) and order distributer

    ch_internal_status_update := make(chan control.Elev) //channel for sending status updates.


    //Communication between order handler and order distributer
    ch_order_to_execute := make(chan elevio.ButtonEvent) //Orders delegated to local elevator is sent to order handler

    ch_order_to_distribute := make(chan elevio.ButtonEvent) //Orders to be delegated is sent from orderhandler to distributer


    //Communication between order distributer and message handler
    ch_broadcast_status_update := make(chan control.Elev)
    ch_external_status_update := make(chan control.Elev)

    ch_broadcast_order := make(chan orderDistributer.ExtOrder)
    ch_receive_external_order := make(chan orderDistributer.ExtOrder)


    ch_bcast_order_executed_at_floor := make(chan []int)
    ch_listen_order_executed_at_floor := make(chan []int)

    ch_order_executed_by_me := make(chan []int)



    //Starting goroutines
    go control.PollInternalElevatorStatus(&elevator, ch_fsm_control_status_updated, ch_internal_status_update) //Whenever the status of the local elevator is updated, send an elevator copy to the order distributer

    go orderHandler.StartOrderHandling(&elevator,ch_fsm_orderHandler_order, ch_order_to_execute, ch_order_to_distribute, ch_orderHandler_fsm_new_order, ch_fsm_orderHandler_order_executed, ch_bcast_order_executed_at_floor, ch_order_executed_by_me)
    //go messageHandler.Receive(&elevator,ch_status_receive)

    go orderDistributer.PollStatusUpdates(ch_internal_status_update, ch_external_status_update, ch_broadcast_status_update)

    go orderDistributer.DistributeOrders(ch_order_to_distribute, ch_receive_external_order, ch_order_to_execute, ch_broadcast_order)

    go orderDistributer.RegisterExecutedOrders(ch_listen_order_executed_at_floor, ch_order_executed_by_me)




    orderDistributer.StartSendingAndReceivingStatusUpdates(ch_broadcast_status_update, ch_external_status_update)

    orderDistributer.StartSendingAndReceivingOrders(ch_broadcast_order, ch_receive_external_order)

    orderDistributer.StartSendingAndReceivingOrderUpdates(ch_bcast_order_executed_at_floor, ch_listen_order_executed_at_floor)



    stateMachine.RunStateMachine(&elevator, ch_fsm_orderHandler_order, ch_orderHandler_fsm_new_order, ch_fsm_control_status_updated, ch_fsm_orderHandler_order_executed)


}
