// code by BLACK ZERO
package main

import (
    "fmt"
    "os"
    "os/exec"
)

func main() {
    if _, err := os.Stat("doom_cannon"); os.IsNotExist(err) {
        fmt.Println("❌ Error: doom_cannon file not found")
        os.Exit(1)
    }

    
    cmd := exec.Command("python3", "doom-cannon")
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    cmd.Stdin = os.Stdin

    err := cmd.Run()
    if err != nil {
        fmt.Println("❌ Failed to run doom_cannon:", err)
        os.Exit(1)
    }
}
