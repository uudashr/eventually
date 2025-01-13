package eventually

// Recorder record the events.
// The events will be recorded on the Events.
type Recorder struct {
	Events []Event
}

// Publish the event to the recorder.
func (r *Recorder) Publish(event Event) error {
	r.Events = append(r.Events, event)

	return nil
}
