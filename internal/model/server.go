package model

type Server struct {
    ID       string `json:"id"`
    Name     string `json:"name"`
    Protocol string `json:"protocol"`
    Address  string `json:"address"`
    Port     int    `json:"port"`
    RawURI   string `json:"raw_uri"`
}
