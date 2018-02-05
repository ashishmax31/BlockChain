package cluster

import (
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"strconv"
	"syscall"
)

var pids = []*os.Process{}

func start() {
	binary, _ := exec.LookPath("blockchain")
	// pid, err := syscall.ForkExec(binary, nil, attr)
	// if err != nil {
	// 	log.Println(err.Error())
	// }
	// println(pid)
	// {"4000", "5000", "6000", "3000", "2000"}
	user, _ := user.Current()
	ports := []string{"4000", "5000", "6000", "3000", "2000"}

	for _, port := range ports {
		cmd := exec.Command(binary, fmt.Sprintf("-p=%s", port))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		uid, _ := strconv.ParseInt(user.Uid, 10, 32)
		gid, _ := strconv.ParseInt(user.Gid, 10, 32)
		cmd.SysProcAttr.Credential = &syscall.Credential{Uid: uint32(uid), Gid: uint32(gid)}
		err := cmd.Start()
		if err != nil {
			panic(err)
		}
		pids = append(pids, cmd.Process)
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		for _, process := range pids {
			fmt.Printf("Killed %d \n", process.Pid)
			err := process.Kill()
			if err != nil {
				panic(err)
			}
		}
		os.Exit(1)
	}()
	pids[0].Wait()

}
