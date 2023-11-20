/*
Implement the following logic, which involves calling the "proc" function at once per second
while ensuring that the program remains running continuously.

Task:
	- call `proc` at once per second
	- program should not exit
*/

package main

import (
	"fmt"
	"time"
)

func main() {

	go func() {
		t := time.NewTicker(time.Second * 1)

		for {
			select {
			case <-t.C:
				go func() {
					defer func() {
						if err := recover(); err != nil {
							fmt.Println(err)
						}
					}()
					proc()
				}()
			}
		}

	}()

	select {}

}

func proc() {
	panic("ok")
}
