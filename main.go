// Code by Bartosz Apanasewicz
// Created 2019-08-28
// ============================================================================

package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
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

var m Machine

// ErrCheck - obsługa błedów
// ============================================================================
func ErrCheck(errNr error) {
	if errNr != nil {
		fmt.Println(errNr)
	}
}

// APIGet - Wysłanie obiektu
// ========================================================
func APIGet(c *gin.Context) {

	c.Header("Access-Control-Allow-Origin", "*")

	c.JSON(http.StatusOK, m)
}

// APIPost - Zmiana obiektu
// ========================================================
func APIPost(c *gin.Context) {

	// var newData Machine
	// err := c.BindJSON(&newData)
	// ErrCheck(err)

	c.Header("Access-Control-Allow-Origin", "*")
	c.JSON(http.StatusOK, gin.H{"Status": "Post OK"})
}

// ElapsedMilliseconds - sprawdznie odstepu czasu
// ============================================================================
func ElapsedMilliseconds(t0 int, czas int) bool {

	t1 := time.Now().Nanosecond()

	if ((t1 - t0) / 1000000) > czas {
		return true
	}

	return false
}

// Options - Obsługa request'u OPTIONS (CORS)
// ========================================================
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// ProgSim - generowanie opóznień
// ============================================================================
func ProgSim() {

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
func ProgInit() {
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
func ProgRun() {

	for {

	}

}

// MAIN
// ============================================================================
func main() {

	// Nasze zmienne
	m.sName = "DTP testing machine"
	m.Robot.sName = "KUKA"
	m.Gripper.sName = "Festo"
	m.WeldingMachine.sName = "Dalex"
	ProgInit()

	// Wydruk struktury
	jsonStream, err := json.MarshalIndent(m, "", "\t")
	ErrCheck(err)
	fmt.Println("===================================================")
	fmt.Println("Wygenerowana struktura...")
	fmt.Println(string(jsonStream))

	// Run machine program
	fmt.Println("Start machine program... " + m.sName)
	go ProgSim()
	go ProgRun()

	/// REST API
	r := gin.Default()
	r.Use(Options)
	r.GET("/", APIGet)
	r.POST("/", APIPost)
	r.Run(":8090")
}
