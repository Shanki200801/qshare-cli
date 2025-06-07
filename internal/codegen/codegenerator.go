package codegen

import (
	"fmt"
	"math/rand"
)

var words = []string{
	"apple", "banana", "cherry", "date", "elderberry", "fig", "grape", "honeydew", "kiwi", "lemon",
	"mango", "nectarine", "orange", "pear", "plum", "raspberry", "strawberry", "tangerine", "watermelon",
	"zebra", "lion", "tiger", "elephant", "giraffe", "hippo", "kangaroo", "leopard", "monkey", "panda",
	"rabbit", "snake", "turtle", "wolf", "zebra",
}

func GenerateCode() string {
	num1 := rand.Intn(10) // 0-9
	num2 := rand.Intn(10) // 0-9
	w1 := words[rand.Intn(len(words))]
	w2 := words[rand.Intn(len(words))]
	return fmt.Sprintf("%d-%s-%d-%s", num1, w1, num2, w2)
}
