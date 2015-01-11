package fzf

import "testing"

func TestReadFromCommand(t *testing.T) {
	strs := []string{}
	eb := NewEventBox()
	reader := Reader{
		pusher:   func(s string) { strs = append(strs, s) },
		eventBox: eb}

	// Check EventBox
	if eb.Peak(EvtReadNew) {
		t.Error("EvtReadNew should not be set yet")
	}

	// Normal command
	reader.readFromCommand(`echo abc && echo def`)
	if len(strs) != 2 || strs[0] != "abc" || strs[1] != "def" {
		t.Errorf("%s", strs)
	}

	// Check EventBox again
	if !eb.Peak(EvtReadNew) {
		t.Error("EvtReadNew should be set yet")
	}

	// Wait should return immediately
	eb.Wait(func(events *Events) {
		if _, found := (*events)[EvtReadNew]; !found {
			t.Errorf("%s", events)
		}
		events.Clear()
	})

	// EventBox is cleared
	if eb.Peak(EvtReadNew) {
		t.Error("EvtReadNew should not be set yet")
	}

	// Failing command
	reader.readFromCommand(`no-such-command`)
	strs = []string{}
	if len(strs) > 0 {
		t.Errorf("%s", strs)
	}

	// Check EventBox again
	if eb.Peak(EvtReadNew) {
		t.Error("Command failed. EvtReadNew should be set")
	}
}
