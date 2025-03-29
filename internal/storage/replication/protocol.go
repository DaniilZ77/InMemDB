package replication

type Request struct {
	LastSegment string
}

type Response struct {
	Filename string
	Segment  []byte
}

func NewRequest(lastSegment string) Request {
	return Request{LastSegment: lastSegment}
}

func NewResponse(filename string, segment []byte) Response {
	return Response{
		Filename: filename,
		Segment:  segment,
	}
}
