// code by BLACK ZERO
package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    if _, err := os.Stat("doom-cannon"); os.IsNotExist(err) {
        fmt.Println("❌ Error: doom-cannon file not found")
        os.Exit(1)
    }

    
    cmd := exec.Command("python3", "doom-cannon")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    if err != nil {
        fmt.Println("❌ Failed to run DooM-CANNON:", err)
        os.Exit(1)
    }
}
