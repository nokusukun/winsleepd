package main

import tui "winsleepd/cmd/tui/cmd"

func main() {
	Main()
}

func Main() error {
	err := tui.Run()
	if err != nil {
		return err
	}
	return nil
}
