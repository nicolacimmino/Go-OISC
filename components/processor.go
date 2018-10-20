package components

/*
 *
 */
type processor interface {
	initProcessor(bus *Bus)
	processResetLine(resetSubscriber BusSubscriber)
}

/**
 *
 */
type processorAbstract struct {
	bus *Bus
	processor
}

/**
 *
 */
func (processor *processorAbstract) initProcessor(bus *Bus, onReset func()) {
	processor.bus = bus

	resetSubscriber := make(BusSubscriber)
	bus.SubscribeToReset(resetSubscriber)

	go processor.processResetLine(resetSubscriber, onReset)
}

/**
 * Watches the bus reset line and resets the processor on a positive transition.
 */
func (processor *processorAbstract) processResetLine(resetSubscriber BusSubscriber, onReset func()) {
	for {
		reset, ok := <-resetSubscriber
		if !ok {
			break
		}

		if reset {
			onReset()
		}
	}
}
