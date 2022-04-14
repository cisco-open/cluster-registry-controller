// Copyright (c) 2021, and 2022 Cisco and/or its affiliates. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package signals

import (
	"context"
	"os"
	"os/signal"
	"syscall"
)

func NotifyContext(parentContext context.Context) context.Context {
	ctx, cancel := context.WithCancel(parentContext)
	c := make(chan os.Signal, 1)
	signal.Notify(c, []os.Signal{os.Interrupt, syscall.SIGTERM}...)

	go func() {
		select {
		case <-c:
			cancel()
			<-c
			os.Exit(1) // second signal. Exit directly.
		case <-ctx.Done():
		}
	}()

	return ctx
}
