/*
Copyright © 2021 Vu Hoang Phan

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
package main

import (
	"humn.ai/phan/cmd"
)

func main() {
	cmd.Execute()
	//mbc := mapbox.NewClient("pk.eyJ1IjoidmVudGlsbyIsImEiOiJja3Q4Ym8xb2swdTVlMnBwYzZtczA1cnVvIn0.B2rQDdmZ5p3v8P7BSO9Thw")
	//in := make(chan models.Input, 10000)
	//out := make(chan models.Output, 10000)
	//worker.NewPool(mbc, in, out, 10).Run(nil)
	//go stdio.Write(os.Stdout, out)
	//go stdio.Read(os.Stdin, in)
	//for  {}
}
