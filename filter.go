package can

type filter struct {
	id      uint32
	handler Handler
}

func newFilter(id uint32, handler Handler) Handler {
	return &filter{
		id:      id,
		handler: handler,
	}
}

func (f *filter) Handle(frame Frame) {
	if frame.ID == f.id {
		f.handler.Handle(frame)
	}
}

type funcFilter struct {
	filter  func(Frame) bool
	handler Handler
}

func newFuncFilter(filterFunc func(Frame) bool, handler Handler) Handler {
	return &funcFilter{
		filter:  filterFunc,
		handler: handler,
	}
}

func (f *funcFilter) Handle(frame Frame) {
	if f.filter(frame) {
		f.handler.Handle(frame)
	}
}
