package main

import (
  "./elevio"
  "./elevator"
  "./orderHandler"
  "./stateMachine"
  "./network"
  "./orderDistributer"
  "./orderDistributer-helpfunc"
  "./timer"
)

var TCP_ConnectionPort = "15657"

func main(){

    //CREATING ELEVATOR OBJECT
    elev := elevator.InitializeElevator(TCP_ConnectionPort)
    //Declaring what ports to UpdateCurrState


    //CREATING CHANNELS
    //**Internal elevator channels**

    ch_new_order := make(chan elevio.ButtonEvent)

    ch_order_registered := make(chan elevio.ButtonEvent)

    ch_order_executed := make(chan bool)

    ch_stat_updated := make(chan bool)


    //**Communication between elevator and orderDistributer**

    ch_int_stat_update := make(chan elevator.Elev)


    //**Communication between orderHandler and orderDistributer**

    ch_order_to_exec := make(chan elevio.ButtonEvent)

    ch_order_to_distribute := make(chan elevio.ButtonEvent)


    //**Communication between order distributer and network
    ch_bcast_stat_update := make(chan elevator.Elev)
    ch_ext_stat_update := make(chan elevator.Elev)

    ch_bcast_order := make(chan orderDistributerHF.ExtOrder)
    ch_rec_ext_order := make(chan orderDistributerHF.ExtOrder)

    ch_bcast_order_exec := make(chan int)
    ch_rec_order_exec := make(chan int,2)

    ch_int_order_exec := make(chan int)

    ch_order_timeout := make(chan elevio.ButtonEvent)




    //STARTING GOROUTINES

    go elevator.PollInternalElevatorStatus(&elev, ch_stat_updated, ch_int_stat_update) //Whenever the status of the local elevator is updated, send an elevator copy to the order distributer

    go orderDistributer.PollStatusUpdates(ch_int_stat_update, ch_ext_stat_update, ch_bcast_stat_update)

    go orderDistributer.PollOrderTimeout(ch_order_timeout)

    go orderHandler.DistributeInternalOrders(&elev,ch_order_registered, ch_order_to_exec, ch_order_to_distribute, ch_new_order)

    go orderDistributer.ReceiveOrders(ch_rec_ext_order, ch_order_to_exec)

    go orderDistributer.DistributeOrders(ch_order_to_distribute, ch_order_to_exec, ch_bcast_order, ch_order_timeout)

    go orderHandler.RegisterExecutedOrders(&elev,ch_order_executed, ch_int_order_exec)

    go orderDistributer.RegisterExecutedOrders(ch_rec_order_exec, ch_int_order_exec, ch_bcast_order_exec)




    //BEGIN SENDING AND RECEIVING BROADCASTS
    network.InitializeStatusUpdates(ch_bcast_stat_update, ch_ext_stat_update)

    network.InitializeOrders(ch_bcast_order, ch_rec_ext_order)

    network.InitializeOrderUpdates(ch_bcast_order_exec, ch_rec_order_exec)




    drv_buttons := make(chan elevio.ButtonEvent)
  	drv_floors  := make(chan int)
  	ch_timer    := make(chan bool)

  	go elevio.PollButtons(drv_buttons)
    go elevio.PollFloorSensor(drv_floors)
    go timer.PollTimeOut(ch_timer)

    //START THE ELEVATOR STATE MACHINE
    go stateMachine.RunStateMachine(&elev, ch_order_registered, ch_new_order, ch_stat_updated, ch_order_executed, drv_buttons, drv_floors, ch_timer)

    select{}

}
