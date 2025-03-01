package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"runtime"

	"github.com/fatih/color"
)

func updateNoping(version string) {
	latestVersion, _, _ := getLatestVersion(version)

	if version != "" {
		fmt.Printf("Updating to version %s\n", color.CyanString(version))
		latestVersion = version
	}

	osArch := fmt.Sprintf("%s-%s", runtime.GOOS, runtime.GOARCH)
	ext := ""
	if runtime.GOOS == "windows" {
		ext = ".exe"
	}

	binaryName := fmt.Sprintf("noping-%s%s", osArch, ext)
	downloadURL := fmt.Sprintf("https://github.com/Bastih18/NoPing/releases/download/%s/%s", latestVersion, binaryName)
	fmt.Printf("Downloading noping from: %s\n", color.BlueString(downloadURL))

	execPath, err := os.Executable()
	if err != nil {
		fmt.Printf("%s Failed to get executable path: %v\n", color.RedString("Error:"), err)
		return
	}

	tmpFile := execPath + ".new"
	out, err := os.Create(tmpFile)
	if err != nil {
		fmt.Printf("%s Failed to create temporary file: %v\n", color.RedString("Error:"), err)
		return
	}
	defer out.Close()
	defer cleanupOnError(tmpFile)

	resp, err := http.Get(downloadURL)
	if err != nil {
		fmt.Printf("%s Failed to download file: %v\n", color.RedString("Error:"), err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		fmt.Printf("%s Version %s not found on GitHub.\n", color.RedString("Error:"), latestVersion)
		return
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("%s Unexpected response from GitHub: %d\n", color.RedString("Error:"), resp.StatusCode)
		return
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		fmt.Printf("%s Failed to copy new binary: %v\n", color.RedString("Error:"), err)
		return
	}

	if runtime.GOOS != "windows" {
		err = os.Chmod(tmpFile, 0755)
		if err != nil {
			fmt.Printf("%s Failed to make binary executable: %v\n", color.RedString("Error:"), err)
			return
		}
	}

	defer func() { tmpFile = "" }()

	err = os.Rename(tmpFile, execPath)
	if err != nil {
		fmt.Printf("%s Failed to replace old binary: %v\n", color.RedString("Error:"), err)
		return
	}

	fmt.Printf("%s Successfully updated to version %s!\n", color.GreenString("✅"), color.CyanString(latestVersion))
}

func cleanupOnError(tmpFile string) {
	if tmpFile != "" {
		_ = os.Remove(tmpFile)
	}
}
