package network

import (
  "../network-helpfunc/bcast"
  "../elevator"
  "../orderDistributer-helpfunc"
)

//Sending and receiving messages

const PORT_STATUS_UPDATES = 20000
const PORT_ORDERS = 20203
const PORT_ORDER_UPDATES = 20194

func StartSendingAndReceivingStatusUpdates(bcast_status_update chan elevator.Elev, receive_status_update chan elevator.Elev) {
  go bcast.Transmitter(PORT_STATUS_UPDATES, bcast_status_update)
  go bcast.Receiver(PORT_STATUS_UPDATES, receive_status_update)
}

func StartSendingAndReceivingOrders(bcast_order chan orderDistributerHF.ExtOrder, receive_ext_order chan orderDistributerHF.ExtOrder) {
  go bcast.Transmitter(PORT_ORDERS, bcast_order)
  go bcast.Receiver(PORT_ORDERS, receive_ext_order)
}

func StartSendingAndReceivingOrderUpdates(bcast_order_executed chan int, receive_order_executed chan int) {
  go bcast.Transmitter(PORT_ORDER_UPDATES, bcast_order_executed)
  go bcast.Receiver(PORT_ORDER_UPDATES, receive_order_executed)
}
