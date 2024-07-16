package main

import "fmt"

func coba(abjad string) bool {
	for i := 0; i < len(abjad); i++ {
		if abjad[i] != abjad[len(abjad)-1-i] {
			return false
		}
	}
	return true
}

func main() {
	for i := 0; i < 5; i++ {
		if i == 0 || i == 5-1 {
			for j := 0; j < 5; j++ {
				fmt.Printf("*")
			}
		} else {
			for k := 0; k < 5; k++ {
				if k == 0 || k == 5-1 {
					fmt.Printf("*")
				} else {
					fmt.Printf(" ")
				}
			}
		}
		fmt.Println()
	}
}
