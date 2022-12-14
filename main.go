package main

import (
	"fmt"
	"taskupdate/initializers"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.SyncDatabase()
	//DB:=initializers.ConnectToDb()
}
func main() {

	fmt.Println("Hi hello")
	initialMigratraion()
	initialRouter()

}
