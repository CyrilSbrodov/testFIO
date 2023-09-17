package main

import "testFIO/internal/app/server"

func main() {
	srv := server.NewServerApp()
	srv.Run()
}
