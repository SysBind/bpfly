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
	"fmt"
	"log"
	"os/exec"
)

type program struct {
	path  string
	redis string
	cmd   *exec.Cmd
}

type Program interface {
	Start()
}

func (prog program) Start() {
	prog.cmd = exec.Command(prog.path)
	stdout, err := prog.cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	scanner := bufio.NewScanner(stdout)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	if err := prog.cmd.Start(); err != nil {
		log.Fatal(err)
	}

	if err := prog.cmd.Wait(); err != nil {
		log.Fatal(err)
	}
}
