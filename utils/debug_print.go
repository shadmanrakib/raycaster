package utils

import "fmt"

func DebugPrint(print bool, a ...any) {
	if print {
		fmt.Println(a...)
	}
}
