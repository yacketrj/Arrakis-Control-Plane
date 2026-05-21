package main

func init() {
	if listenAddr == ":8080" || listenAddr == "" {
		listenAddr = "127.0.0.1:8080"
	}
}
