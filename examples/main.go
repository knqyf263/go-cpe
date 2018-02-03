package main

import (
	"fmt"

	"github.com/knqyf263/go-cpe/common"
	"github.com/knqyf263/go-cpe/matching"
	"github.com/knqyf263/go-cpe/naming"
)

func main() {
	wfn, err := naming.UnbindURI("cpe:/a:microsoft:internet_explorer%01%01%01%01:8%02:beta")
	if err != nil {
		fmt.Println(err)
		return
	}
	wfn2, err := naming.UnbindFS(`cpe:2.3:a:microsoft:internet_explorer:8.0.6001:beta:*:*:*:*:*:*`)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(wfn, wfn2)
	fmt.Println(matching.IsDisjoint(wfn, wfn2)) // false
	fmt.Println(matching.IsEqual(wfn, wfn2))    // false
	fmt.Println(matching.IsSubset(wfn, wfn2))   // false
	fmt.Println(matching.IsSuperset(wfn, wfn2)) // true

	wfn3 := common.NewWellFormedName()
	wfn3.Set("part", "a")
	wfn3.Set("vendor", "microsoft")
	wfn3.Set("product", "internet_explorer????")
	wfn3.Set("version", "8\\.0\\.6001")
	fmt.Println(naming.BindToURI(wfn3)) // cpe:/a:microsoft:internet_explorer%01%01%01%01:8.0.6001
	fmt.Println(naming.BindToFS(wfn3))  // cpe:2.3:a:microsoft:internet_explorer????:8.0.6001:*:*:*:*:*:*:*
}
