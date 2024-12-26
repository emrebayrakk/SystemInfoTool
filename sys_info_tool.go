package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"

	"github.com/shirou/gopsutil/process"
)

type DataItem struct {
	SEMBOL   string  `json:"SEMBOL"`
	ACIKLAMA string  `json:"ACIKLAMA"`
	KAPANIS  float64 `json:"KAPANIS"`
	ALIS     float64 `json:"ALIS"`
	SATIS    float64 `json:"SATIS"`
}

type Response struct {
	Data []DataItem `json:"data"`
}

func getWiFiInfo() {
	cmd := exec.Command("netsh", "wlan", "show", "profiles")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error fetching WiFi profiles:", err)
		return
	}

	fmt.Println("WiFi Profiles and Passwords:")
	for _, line := range strings.Split(string(output), "\n") {
		if strings.Contains(line, "All User Profile") {
			fields := strings.Split(line, ":")
			if len(fields) > 1 {
				profile := strings.TrimSpace(fields[1])
				cmdPass := exec.Command("netsh", "wlan", "show", "profile", profile, "key=clear")
				passOutput, _ := cmdPass.Output()
				password := "<Not Available>"
				for _, passLine := range strings.Split(string(passOutput), "\n") {
					if strings.Contains(passLine, "Key Content") {
						passFields := strings.Split(passLine, ":")
						if len(passFields) > 1 {
							password = strings.TrimSpace(passFields[1])
						}
					}
				}
				fmt.Printf("Profile: %s, Password: %s\n", profile, password)
			}
		}
	}
}

func scheduleShutdown(minutes int) {
	fmt.Printf("System will shut down in %d minutes.\n", minutes)
	cmd := exec.Command("shutdown", "/s", "/t", strconv.Itoa(minutes*60))
	if err := cmd.Run(); err != nil {
		fmt.Println("Error scheduling shutdown:", err)
	}
}

func getResourceIntensiveApp() {
	procs, err := process.Processes()
	if err != nil {
		fmt.Println("Error fetching process information:", err)
		return
	}
	var maxCPUPid, maxMemPid int32
	var maxCPU, maxMem float64
	for _, proc := range procs {
		cpuPercent, _ := proc.CPUPercent()
		memInfo, _ := proc.MemoryInfo()
		if cpuPercent > maxCPU {
			maxCPU = cpuPercent
			maxCPUPid = proc.Pid
		}
		if memInfo != nil && float64(memInfo.RSS) > maxMem {
			maxMem = float64(memInfo.RSS)
			maxMemPid = proc.Pid
		}
	}
	if maxCPUPid != 0 {
		proc, _ := process.NewProcess(maxCPUPid)
		name, _ := proc.Name()
		fmt.Printf("Application using most CPU: %s (%.2f%%)\n", name, maxCPU)
	}
	if maxMemPid != 0 {
		proc, _ := process.NewProcess(maxMemPid)
		name, _ := proc.Name()
		fmt.Printf("Application using most memory: %s (%.2f MB)\n", name, maxMem/1024/1024)
	}
}


func displaySelectedData() {
	cmd := exec.Command("curl", "https://api.bigpara.hurriyet.com.tr/doviz/headerlist/anasayfa")
	output, err := cmd.Output()
	if err != nil {
		fmt.Println("Error fetching data:", err)
		return
	}

	var response Response
	err = json.Unmarshal(output, &response)
	if err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, item := range response.Data {
		switch item.SEMBOL {
		case "XU100":
			fmt.Printf("BIST 100 - KAPANIS: %.2f\n", item.KAPANIS)
		case "EURTRY":
			fmt.Printf("EURO/TURK LIRASI - ALIS: %.4f, SATIS: %.4f\n", item.ALIS, item.SATIS)
		case "GLDGR":
			fmt.Printf("ALTIN GRAM - TL - ALIS: %.4f, SATIS: %.4f\n", item.ALIS, item.SATIS)
		case "USDTRY":
			fmt.Printf("DOLAR/TURK LIRASI - ALIS: %.4f, SATIS: %.4f\n", item.ALIS, item.SATIS)
		}
	}
}

func main() {
	if runtime.GOOS != "windows" {
		fmt.Println("This program is designed to run on Windows.")
		return
	}

	for {
		fmt.Println("\nMenu:")
		fmt.Println("1: Display WiFi profiles and passwords")
		fmt.Println("2: Schedule shutdown")
		fmt.Println("3: Show application using most RAM and CPU")
		fmt.Println("4: Display selected financial data")
		fmt.Println("0: Exit")
		fmt.Print("Enter your choice: ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			getWiFiInfo()
		case 2:
			fmt.Print("Enter shutdown delay in minutes: ")
			var minutes int
			fmt.Scan(&minutes)
			scheduleShutdown(minutes)
		case 3:
			getResourceIntensiveApp()
		case 4:
			displaySelectedData()
		case 0:
			fmt.Println("Exiting the program.")
			return
		default:
			fmt.Println("Invalid choice. Please try again.")
		}
	}
}