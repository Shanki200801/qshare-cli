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
	num := rand.Intn(10)
	w1 := words[rand.Intn(len(words))]
	w2 := words[rand.Intn(len(words))]
	return fmt.Sprintf("%d-%s-%s", num, w1, w2)
}