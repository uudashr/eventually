package eventually

// Recorder record the events.
// The events will be recorded on the Events.
type Recorder struct {
	Events []Event
}

// NewRecorder allocates and returns a new Recorder.
func NewRecorder() *Recorder {
	return &Recorder{}
}

// Publish the event to the recorder.
func (r *Recorder) Publish(event Event) {
	r.Events = append(r.Events, event)
}
