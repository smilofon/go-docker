package docker

import (
	"bytes"
	"fmt"
	"math/rand"
	"os/exec"
	"strings"
	"time"
)

type DockerContext struct {
	DockerName string
	DockerIp   string
}

// randString generates a random string of 32 characters
func randString() string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune("qwertyuiopasdfghjklzxcvbnm")
	b := make([]rune, 32)
	for i := range b {
		b[i] = runes[rand.Intn(len(runes))]
	}
	return string(b)
}

// Get Ip address of a running container
func GetDockerIp(name string) (string, error) {
	for i := 0; i < 10; i++ {
		var out bytes.Buffer
		cmd := exec.Command("docker", "inspect", "-f", "{{range.NetworkSettings.Networks}}{{.IPAddress}}{{end}}", name)
		cmd.Stdout = &out
		err := cmd.Run()
		if err == nil {
			result := strings.Trim(out.String(), " \r\n")
			if len(result) > 0 {
				return result, nil
			}
		}
		time.Sleep(1 * time.Second)
	}
	return "", fmt.Errorf("could not get IP address of Docker container")
}

// Setup starts a container in detach mode. Additionnal arguments (like volume or environment variables) and command can be given.
func SetUp(imageName string, arg []string, cmd []string) (*DockerContext, error) {
	ctx := &DockerContext{
		DockerName: randString(),
	}
	args := []string{"run", "--rm", "--name", ctx.DockerName, "-d"}
	if len(arg) != 0 {
		args = append(args, arg...)
	}
	args = append(args, imageName)
	if len(cmd) != 0 {
		args = append(args, cmd...)
	}
	c := exec.Command("docker", args...)
	err := c.Run()
	if err != nil {
		return nil, err
	}
	ctx.DockerIp, err = GetDockerIp(ctx.DockerName)
	if err != nil {
		return nil, err
	}
	return ctx, nil
}

// TearDown destroys the running container
func TearDown(ctx *DockerContext) error {
	cmd := exec.Command("docker", "kill", ctx.DockerName)
	err := cmd.Run()
	if err != nil {
		return err
	}
	cmd = exec.Command("docker", "rm", ctx.DockerName)
	_ = cmd.Run() //ignore error. It should be useless since we started container with "--rm" option, but do it again :)
	return nil
}
