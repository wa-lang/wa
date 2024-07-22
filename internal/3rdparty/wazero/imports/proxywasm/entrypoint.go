// Copyright 2020-2021 Tetrate
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package proxywasm

import (
	"wa-lang.org/wa/internal/3rdparty/wazero/imports/proxywasm/internal"
	"wa-lang.org/wa/internal/3rdparty/wazero/imports/proxywasm/types"
)

// SetVMContext is the entrypoint for setting up the entire Wasm VM.
// Please make sure to call this entrypoint during "main()" function;
// otherwise, the VM fails.
func SetVMContext(ctx types.VMContext) {
	internal.SetVMContext(ctx)
}
