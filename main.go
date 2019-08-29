// Code by Bartosz Apanasewicz
// Created 2019-08-28
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"time"
)

// Machine - struktura maszyny
// ============================================================================
type Machine struct {
	sName string
	Robot struct {
		sName          string
		bHomePos       bool
		bRobotMove     bool
		bProgramActive bool
	}
	WeldingMachine struct {
		sName               string
		bValve              bool
		bPosUp              bool
		bPosDown            bool
		bPreasureSensorBool bool
		iCylinderPosition   int
	}
	Gripper struct {
		sName     string
		bValve    bool
		bPosOpen  bool
		bPosClose bool
	}
}

// ErrCheck - obsługa błedów
// ============================================================================
func ErrCheck(errNr error) {
	if errNr != nil {
		fmt.Println(errNr)
	}
}

// ElapsedMilliseconds - sprawdznie odstepu czasu
// ============================================================================
func ElapsedMilliseconds(t0 int, czas int) bool {

	t1 := time.Now().Nanosecond()

	if ((t1 - t0) / 1000000) > czas {
		return true
	} else {
		return false
	}
}

// ProgSim - generowanie opóznień
// ============================================================================
func ProgSim(m *Machine) {

	GripperValveTime0 := time.Now()

	GripperValvePrev := m.Gripper.bValve

	for {
		t := time.Now()

		// -------------------
		// Symulacja grippera
		// -------------------

		// Gdy chcemy zamknąć pojawia się jedynka na zaworze
		if !GripperValvePrev && m.Gripper.bValve {
			GripperValveTime0 = t
			m.Gripper.bPosClose = false
			m.Gripper.bPosOpen = false
		}

		// Saprawdzmy czy minał czas
		if ElapsedMilliseconds(GripperValveTime0.Nanosecond(), 5000) {
			m.Gripper.bPosClose = true
		}

		GripperValvePrev = m.Gripper.bValve
	}
}

// ProgInit - inicjacja zmiennych
// ============================================================================
func ProgInit(m *Machine) {
	m.Gripper.bValve = false
	m.Gripper.bPosClose = false
	m.Gripper.bPosOpen = true

	m.Robot.bHomePos = true
	m.Robot.bProgramActive = true
	m.Robot.bRobotMove = false

	m.WeldingMachine.bPosDown = false
	m.WeldingMachine.bPosUp = true
	m.WeldingMachine.bPreasureSensorBool = false
	m.WeldingMachine.bValve = false
	m.WeldingMachine.iCylinderPosition = 100
}

// ProgRun - program realizowany przez maszynę (PLC+Robot)
// ============================================================================
func ProgRun(m *Machine) {

	for {

	}

}

// MAIN
// ============================================================================
func main() {

	// Nasze zmienne
	var m Machine
	m.sName = "DTP testing machine"
	m.Robot.sName = "KUKA"
	m.Gripper.sName = "Festo"
	m.WeldingMachine.sName = "Dalex"
	ProgInit(&m)

	// Wydruk struktury
	jsonStream, err := json.MarshalIndent(m, "", "\t")
	ErrCheck(err)
	fmt.Println("===================================================")
	fmt.Println("Wygenerowana struktura...")
	fmt.Println(string(jsonStream))

	// Run machine program
	fmt.Println("Start machine program... " + m.sName)
	go ProgSim(&m)
	ProgRun(&m)
}
