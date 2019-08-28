// Code by Bartosz Apanasewicz
// Created 2019-08-28
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

// Struktury
// ============================================================================

// SignalBool - zmienna TRUE/FALSE
// Może być I/O lub flaga
// ----------------------------------------------------------------------------
type SignalBool struct {
	ID  string
	Val bool
}

// SignalInt - zmienna integer
// Może być odczyt z czujnika analogowego
// ----------------------------------------------------------------------------
type SignalInt struct {
	ID  string
	Val int
}

// Signals - tablica zmiennych
// ----------------------------------------------------------------------------
type Signals struct {
	SignalsBool []SignalBool
	SignalsInt  []SignalInt
}

// Equipment - część maszyny, np. robot
// Zawiera nazwę i stany sygnałów
// ----------------------------------------------------------------------------
type Equipment struct {
	ID      string
	Signals Signals
}

// Machine - zawiera equipmenty
// ----------------------------------------------------------------------------
type Machine struct {
	ID        string
	Equipment []Equipment
}

// ErrCheck - obsługa błedów
// ============================================================================
func ErrCheck(errNr error) {
	if errNr != nil {
		fmt.Println(errNr)
	}
}

// MAIN
// ============================================================================
func main() {

	// Nasze zmienne
	var machine Machine

	// Otwieranie pliku JSON
	jsonFile, err := os.Open("machine_def.json")
	defer jsonFile.Close()
	ErrCheck(err)

	// Odczyt bajtów z pliku
	byteStream, err := ioutil.ReadAll(jsonFile)
	ErrCheck(err)

	// Wydruk struktury
	fmt.Println("===================================================")
	fmt.Println("Odczytane bajty...")
	fmt.Println(string(byteStream))

	// Parsowanie do obiektu
	json.Unmarshal(byteStream, &machine)

	// Wydruk struktury
	fmt.Println("===================================================")
	fmt.Println("Dostęp do struktury (ID)...")
	fmt.Println(machine.ID)

	// Wykonanie odwrotnej analizy
	jsonStream, err := json.MarshalIndent(machine, "", "\t")
	ErrCheck(err)

	// Wydruk struktury
	fmt.Println("===================================================")
	fmt.Println("Wygenerowana struktura...")
	fmt.Println(string(jsonStream))
}
