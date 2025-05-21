package commands

type PutCommandRequestPayload struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type PutCommandResponsePayload struct{}

type GetCommandRequestPayload struct {
	Key string `json:"key"`
}

type GetCommandResponsePayload struct {
	Value string `json:"value"`
	Ok    bool   `json:"ok"`
}

type DeleteCommandRequestPayload struct {
	Key string `json:"key"`
}

type DeleteCommandResponsePayload struct {
	Ok bool `json:"ok"`
}

type ListCommandRequestPayload struct {}

type ListCommandResponsePayload struct {
	Value []string `json:"value"`
	Ok    bool   `json:"ok"`
}

const (
	PutCommandName    string = "put"
	GetCommandName    string = "get"
	DeleteCommandName string = "delete"
	ListCommandName string = "list"
)
