package model

type Server struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Address  string `json:"address"`
    Port     int    `json:"port"`
    Protocol string `json:"protocol"`
    URI      string `json:"uri"`
    Latency  int    `json:"latency"`
}

type Subscription struct {
    ID     string `json:"id"`
    URL    string `json:"url"`
    Name   string `json:"name"`
    Update int64  `json:"update"`
}
