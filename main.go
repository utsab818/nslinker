package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strconv"

	"github.com/containerd/containerd"
	"github.com/containerd/containerd/cio"
	"github.com/containerd/containerd/namespaces"
)

type ContainerInfo struct {
	ID      string
	PID     uint32
	PodName string
}

func getRunningContainers() ([]ContainerInfo, error) {
	ctx := namespaces.WithNamespace(context.Background(), "test")

	client, err := containerd.New("/run/containerd/containerd.sock")
	if err != nil {
		fmt.Printf("Error connecting to Containerd: %v\n", err)
		return nil, err
	}

	defer client.Close()

	containers, err := client.Containers(ctx)
	if err != nil {
		fmt.Printf("Error listing containers: %v\n", err)
		return nil, err
	}

	var containerInfos []ContainerInfo

	for _, container := range containers {
		info, err := container.Info(ctx)
		if err != nil {
			fmt.Printf("Error getting container info: %v\n", err)
			continue
		}

		task, err := container.Task(ctx, cio.NewAttach()) //attaches to the running container's task and retrieves its process details.
		podName := getPodNameWithCID(info.ID)

		containerInfos = append(containerInfos, ContainerInfo{
			ID:      info.ID,
			PID:     task.Pid(),
			PodName: podName,
		})

		fmt.Printf("Container ID: %s. PID:%d\n", info.ID, task.Pid())
	}

	return containerInfos, nil
}

func getPodNameWithCID(cid string) string {
	jsonFile, err := ioutil.ReadFile("/var/run/containerd/io.containerd.runtime.v2.task/test/" + cid + "/config.json")
	if err != nil {
		panic(err)
	}

	var jsonData map[string]interface{}
	if err := json.Unmarshal(jsonFile, &jsonData); err != nil {
		panic(err)
	}

	hostname, ok := jsonData["hostname"].(string)
	if !ok {
		fmt.Println("Failed to parse hostname as a string")
		return ""
	}

	return hostname
}

func createSymLinks(cn ContainerInfo) error {
	pidStr := strconv.Itoa(int(cn.PID))
	targetPath := "/proc/" + pidStr + "/ns/net"
	symlinkPath := "/var/run/netns/" + cn.PodName

	_, err := os.Stat(symlinkPath)

	if err == nil {
		return fmt.Errorf("Already exists")
	} else if os.IsNotExist(err) {
		cmdName := "ln"
		cmdArgs := []string{"-sfT", targetPath, symlinkPath}

		cmd := exec.Command(cmdName, cmdArgs...)

		_, err := cmd.CombinedOutput()
		if err != nil {
			return err
		}
	} else {
		fmt.Printf("Error checking path %s: %v\n", symlinkPath, err)
	}

	fmt.Println("Symlink created:", symlinkPath)
	return nil
}

func main() {
	cns, err := getRunningContainers()
	if err != nil {
		log.Fatal(err)
	}

	for _, c := range cns {
		err := createSymLinks(c)
		if err != nil {
			fmt.Printf("Error processing container: %s, %s \n", c.ID, err.Error())
		}
	}
}
