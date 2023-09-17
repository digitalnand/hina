package hina

import (
	"fmt"
)

func (printNode PrintNode) Evaluate() {
	fmt.Println(printNode.Value)
}
