package main

import (
	"fmt"
	"os"

	"github.com/looplab/fsm"
)

func main() {
	events := []string{"Preparing", "Planned Reboot", "PXE Booting", "Preparing", "Registering", "Waiting", "Installing", "Booting New Kernel", "Phoned Home"}

	m := fsm.NewFSM(
		"PXE Booting",
		fsm.Events{
			{Name: "PXE Booting", Src: []string{"Planned Reboot"}, Dst: "PXE Booting"},
			{Name: "Preparing", Src: []string{"Planned Reboot"}, Dst: "Preparing"},
			{Name: "PXE Booting", Src: []string{"PXE Booting"}, Dst: "PXE Booting"},
			{Name: "Preparing", Src: []string{"PXE Booting"}, Dst: "Preparing"},
			{Name: "Registering", Src: []string{"Preparing"}, Dst: "Registering"},
			{Name: "Planned Reboot", Src: []string{"Preparing"}, Dst: "Planned Reboot"},
			{Name: "Planned Reboot", Src: []string{"Registering"}, Dst: "Planned Reboot"},
			{Name: "Waiting", Src: []string{"Registering"}, Dst: "Waiting"},
			{Name: "Installing", Src: []string{"Waiting"}, Dst: "Installing"},
			{Name: "Planned Reboot", Src: []string{"Waiting"}, Dst: "Planned Reboot"},
			{Name: "Booting New Kernel", Src: []string{"Installing"}, Dst: "Booting New Kernel"},
			{Name: "Planned Reboot", Src: []string{"Installing"}, Dst: "Planned Reboot"},
			{Name: "Phoned Home", Src: []string{"Booting New Kernel"}, Dst: "Phoned Home"},
			{Name: "PXE Booting", Src: []string{"Crashed"}, Dst: "PXE Booting"},
			{Name: "Preparing", Src: []string{"Crashed"}, Dst: "Preparing"},
			{Name: "PXE Booting", Src: []string{"Phoned Home"}, Dst: "PXE Booting"},
			{Name: "Preparing", Src: []string{"Phoned Home"}, Dst: "Preparing"},
			{Name: "PXE Booting", Src: []string{"Reset Fail Count"}, Dst: "PXE Booting"},
		},
		fsm.Callbacks{
			"enter_state": func(e *fsm.Event) {
				fmt.Println(e.Dst)
			},
		},
	)
	readSequence(m, events)

	visualize(m)
}

func readSequence(m *fsm.FSM, eventSequence []string) {
	for _, event := range eventSequence {
		err := m.Event(event)
		if err != nil {
			fmt.Println(err)
		}
	}
}

func visualize(m *fsm.FSM) {
	file, err := os.Create("graph.dot")
	if err != nil {
		fmt.Println(err)
		return
	}

	graph := fsm.Visualize(m)

	_, err = file.WriteString(graph)
	if err != nil {
		fmt.Println(err)
		return
	}

	err = file.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
