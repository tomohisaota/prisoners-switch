package strategy

import "github.com/tarao/prisoners-switch/rule"

// MyNewStrategy returns a new strategy
func MyNewStrategy() rule.Strategy {
	register := make(chan int)
	isCounter := make(chan bool)
	go func(register chan int, isCounter chan bool) {
		counterNumber := <-register
		isCounter <- true
		for true {
			number := <-register
			isCounter <- number == counterNumber
		}

	}(register, isCounter)

	return &myStrategy{register, isCounter}
}

type myStrategy struct {
	register  chan int
	isCounter chan bool
}

// What if number is not sequence?
func (s *myStrategy) NewPrisoner(number int, shout chan rule.Shout) rule.Prisoner {
	return &prisoner{number: number, register: s.register, isCounter: s.isCounter, shout: shout}
}

type prisoner struct {
	number    int
	register  chan int
	isCounter chan bool
	shout     chan rule.Shout

	isInitialized bool
	counter       int

	haveSeenOn bool
	switched   bool
}

func (p *prisoner) Enter(room rule.Room) {
	if rule.TotalPrisoners == 1 {
		// Am I the only one?
		p.shout <- rule.Triumph
		return
	}
	p.register <- p.number
	iAmCounter := <-p.isCounter
	aIsOn := room.TakeSwitchA().State()
	if !iAmCounter {
		if p.switched {
			return;
		}
		if aIsOn {
			p.haveSeenOn = true
		} else {
			if p.haveSeenOn {
				room.TakeSwitchA().Toggle()
				p.switched = true
			}
		}
		return
	}
	// counter
	if (!p.isInitialized) {
		// make sure that intial state is 1!!
		if (aIsOn) {
			// intial state was 1
			room.TakeSwitchA().Toggle()
		} else {
			/// intial state was 0
			p.counter -= 1
			room.TakeSwitchA().Toggle()
		}
		p.isInitialized = true
		return
	}
	if (aIsOn) {
		p.counter += 1;
		if p.counter == rule.TotalPrisoners-1 {
			p.shout <- rule.Triumph
			return
		}
		room.TakeSwitchA().Toggle()
	} else {
		p.counter -= 1;
		room.TakeSwitchA().Toggle()
	}
}
