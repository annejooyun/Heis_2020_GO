package control

N_FLOORS = 4
N_BUTTONS = 3 // Number of button types (cab,hall up, hall down)

type State int
const (
  idle   State = 0
  moving       = 1
  doorOpen     = 2
)

type Elev struct {
  id int
  state State
  motorDirection MotorDirection
  prevFloor int
  orderList [N_FLOORS][N_BUTTONS]int


}
