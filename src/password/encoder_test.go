package password

import (
	"fmt"
)

// Test output generated from: https://approsto.com/sha-generator/

func ExampleEncode() {
	encodedPassword := Encode("P@ssW0rd!")
	fmt.Println(encodedPassword)
	// Output: 62+j0x1/W8bCgSgF3YggMtf+AfOqb28xuOXvKvTXBs8iDZDwQci9cGBiNdHvHHyywclJeKIhPWoftStSNJdf5g==
}

func ExampleEncodeEmptyString() {
	fmt.Println(Encode(""))
	// Output: z4PhNX7vuL3xVChQ1m2AB9Yg5AULVxXcg/SpIdNs6c5H0NE8XYXysP+DGNKHfuwvY7kxvUdBeoGlODJ6+SfaPg==
}
