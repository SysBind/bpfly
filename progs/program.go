/*


Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package progs

import (
	"bufio"
	"log"
	"os/exec"

	redis "github.com/go-redis/redis/v8"
)

type program struct {
	path        string
	redis       string
	cmd         *exec.Cmd
	interpreter func(line string, rdb *redis.Client)
	rdb         *redis.Client
}

type Program interface {
	Start()
}

func (prog program) Start() {
	prog.rdb = redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
	})
	log.Println("Calling exec..")
	prog.cmd = exec.Command(prog.path)
	stdout, err := prog.cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderr, err := prog.cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	err_scanner := bufio.NewScanner(stderr)
	go func() {
		for err_scanner.Scan() {
			log.Printf("stderr said: %s\n", err_scanner.Text())
			log.Printf("bpfly should run with root privileges. (or otherwise sufficient privileges to run BPF programs)")
		}
	}()

	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			prog.interpreter(scanner.Text(), prog.rdb)
		}
	}()

	if err := prog.cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := prog.cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
