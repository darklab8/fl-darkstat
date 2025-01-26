package flhash

import (
	"fmt"
	"testing"
)

func TestHashFactions(t *testing.T) {
	tests := []struct {
		hash HashCode
		name string
	}{
		{hash: HashCode(45353), name: "fc_rn_grp"},
		{hash: HashCode(4169), name: "fc_freelancer"},
	}

	for _, test := range tests {
		code := HashFaction(test.name)
		if code != test.hash {
			t.Errorf("%s: expected %d but got %d", test.name, test.hash, code)
		}
	}
}

func TestFindInitNumber(t *testing.T) {
	var table [256]uint32
	NotSimplePopulateTable(0x1021, &table)
	_ = table
	fmt.Println()
}
