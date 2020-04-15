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

func StartSendingAndReceivingStatusUpdates(sendStatusUpdates chan elevator.Elev, receiveStatusUpdates chan elevator.Elev) {
  go bcast.Transmitter(PORT_STATUS_UPDATES, sendStatusUpdates)
  go bcast.Receiver(PORT_STATUS_UPDATES, receiveStatusUpdates)
}

func StartSendingAndReceivingOrders(broadcastMessage chan orderDistributerHF.ExtOrder, receiveMessage chan orderDistributerHF.ExtOrder) {
  go bcast.Transmitter(PORT_ORDERS, broadcastMessage)
  go bcast.Receiver(PORT_ORDERS, receiveMessage)
}

func StartSendingAndReceivingOrderUpdates(broadcastOrderExecuted chan []int, listenOrderExecuted chan []int) {
  go bcast.Transmitter(PORT_ORDER_UPDATES, broadcastOrderExecuted)
  go bcast.Receiver(PORT_ORDER_UPDATES, listenOrderExecuted)
}
