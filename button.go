package main

const (
	CLICK_MS        = 40
	CLICK_DOUBLE_MS = 300
	CLICK_LONG_MS   = 500
	CLICK_EXLONG_MS = 1500
)

var (
	input     [8]int64
	inputLast [8]int64
)

// state 0: button down
// state 1: normal click (<750ms)
// state 2: double-click (<750ms)
// state 3: long-hold (>750ms)
// state 4: extra-long-hold (>2s)
func buttonEvent(num, state uint8) {
	print("button ")
	print(num)
	switch state {
	case 0:
		println(" down")
	case 1:
		println(" click")
	case 2:
		println(" double-click")
	case 3:
		println(" long-hold")
	default:
		println(" extra-long-hold")
	}
}

func processInputEvent(gpio uint8, ticks int64) {
	for i := uint8(0); i < 8; i++ {
		isActive := (gpio & (1 << i)) > 0
		if isActive {
			// emit button down event
			input[i] = ticks
			go buttonEvent(i, 0)
			continue
		}

		isButtonUp := input[i] > 0
		if !isButtonUp {
			continue
		}

		// process button-up event
		// gpio active -> inactive transition
		elapsed := ticks - input[i]
		input[i] = 0
		switch {
		case elapsed > CLICK_EXLONG_MS:
			print(elapsed)
			print("ms --> ")
			go buttonEvent(i, 4) // extra-long-hold

		case elapsed > CLICK_LONG_MS:
			print(elapsed)
			print("ms --> ")
			go buttonEvent(i, 3) // long-hold

		case elapsed > CLICK_MS: // normal click / double-click
			print(elapsed)
			print("ms --> ")
			elapsed = ticks - inputLast[i]
			inputLast[i] = ticks
			if elapsed < CLICK_DOUBLE_MS { // double-click time
				print(elapsed)
				print("ms --> ")
				go buttonEvent(i, 2)
			} else {
				print(elapsed)
				print("ms --> ")
				go buttonEvent(i, 1)
			}
		}
	}
}
