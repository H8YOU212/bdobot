package utils

import (
	"fmt"
	"time"
)

// Times the function
func TimeIt(start time.Time, functionName string) {
	elapsed := time.Since(start)
	fmt.Printf("%s took %s\n", functionName, elapsed)
}
