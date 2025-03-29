package replication

type Request struct {
	LastSegment string
}

type Response struct {
	Ok       bool
	Filename string
	Segment  []byte
}

func NewRequest(lastSegment string) Request {
	return Request{LastSegment: lastSegment}
}

func NewSuccessResponse(filename string, segment []byte) Response {
	return Response{
		Ok:       true,
		Filename: filename,
		Segment:  segment,
	}
}

func NewErrorResponse() Response {
	return Response{}
}
