package fetcher

import "testing"

func TestParseVless(t *testing.T) {
    uri := "vless://uuid@1.2.3.4:443?security=tls&type=ws&path=%2Fpath&host=example.com#TestServer"
    server, err := parseURI(uri)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if server.Address != "1.2.3.4" {
        t.Errorf("Expected address 1.2.3.4, got %s", server.Address)
    }
    if server.Port != 443 {
        t.Errorf("Expected port 443, got %d", server.Port)
    }
    if server.Name != "TestServer" {
        t.Errorf("Expected name TestServer, got %s", server.Name)
    }
    if server.Protocol != "vless" {
        t.Errorf("Expected protocol vless, got %s", server.Protocol)
    }
}

func TestParseVmess(t *testing.T) {
    uri := "vmess://eyJ2IjoiMiIsInBzIjoiVGVzdFZtZXNzIiwiYWRkIjoiNS42LjcuOCIsInBvcnQiOiI4MDgwIn0="
    server, err := parseURI(uri)
    if err != nil {
        t.Fatalf("Expected no error, got %v", err)
    }
    if server.Address != "5.6.7.8" {
        t.Errorf("Expected address 5.6.7.8, got %s", server.Address)
    }
    if server.Port != 8080 {
        t.Errorf("Expected port 8080, got %d", server.Port)
    }
    if server.Name != "TestVmess" {
        t.Errorf("Expected name TestVmess, got %s", server.Name)
    }
}

func TestInvalidURI(t *testing.T) {
    _, err := parseURI("trojan://invalid")
    if err == nil {
        t.Fatal("Expected error for unsupported protocol, got nil")
    }
}
