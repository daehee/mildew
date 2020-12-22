package mildew

import (
	"bufio"
	"fmt"
	"os"
)

func (mw *Mildew) OutputFile(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("error creating file: %v", err)
	}

	w := bufio.NewWriter(f)
	for _, v := range mw.Subs.Slice() {
		w.WriteString(fmt.Sprintf("%s\n", v))
	}
	w.Flush()

	return nil
}

func (mw *Mildew) OutputScreen() {
	f := os.Stdout
	w := bufio.NewWriter(f)
	for _, v := range mw.Subs.Slice() {
		w.WriteString(fmt.Sprintf("%s\n", v))
	}
	w.Flush()
}
