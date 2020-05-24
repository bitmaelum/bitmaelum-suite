package server

type OutputStatus struct {
    Error bool `json:"error"`
    Status string `json:"status"`
}

func StatusOk(status string) *OutputStatus {
    return &OutputStatus{
        Error: false,
        Status: status,
    }
}
func StatusError(status string) *OutputStatus {
    return &OutputStatus{
        Error: true,
        Status: status,
    }
}
