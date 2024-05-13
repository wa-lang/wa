// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package dap

import (
	"encoding/json"
	"reflect"
	"testing"
)

// -------- Responses & Requests --------

var errorResponseString = `{"seq":11,"type":"response","request_seq":9,"command":"stackTrace","success":false,"message":"Unable to produce stack trace: \"{e}\"","body":{"error":{"id":2004,"format":"Unable to produce stack trace: \"{e}\"","variables":{"e":"Unknown goroutine 1"},"showUser":true}}}`
var errorResponseStruct = ErrorResponse{
	Response: Response{
		ProtocolMessage: ProtocolMessage{
			Seq:  11,
			Type: "response",
		},
		Command:    "stackTrace",
		RequestSeq: 9,
		Success:    false,
		Message:    "Unable to produce stack trace: \"{e}\"",
	},
	Body: ErrorResponseBody{
		Error: &ErrorMessage{
			Id:        2004,
			Format:    "Unable to produce stack trace: \"{e}\"",
			Variables: map[string]string{"e": "Unknown goroutine 1"},
			ShowUser:  true,
		},
	},
}

// -------- CancelRequest

var cancelRequestString = `{"seq":25,"type":"request","command":"cancel","arguments":{"requestId":24}}`
var cancelRequestStruct = CancelRequest{
	Request:   *newRequest(25, "cancel"),
	Arguments: &CancelArguments{RequestId: 24},
}

var cancelResponseString = `{"seq":25,"type":"response","request_seq":26,"command":"cancel","success":true}`
var cancelResponseStruct = CancelResponse{Response: *newResponse(25, 26, "cancel", true)}

// -------- RunInTerminalRequest

var runInTerminalRequestString = `{"seq":45,"type":"request","command":"runInTerminal","arguments":{"kind":"integrated","title":"Some Title","cwd":"/working/dir","args":["mycommand","arg1","arg2"]}}`
var runInTerminalRequestStruct = RunInTerminalRequest{
	Request: *newRequest(45, "runInTerminal"),
	Arguments: RunInTerminalRequestArguments{
		Kind:  "integrated",
		Title: "Some Title",
		Cwd:   "/working/dir",
		Args:  []string{"mycommand", "arg1", "arg2"},
	},
}

var runInTerminalResponseString = `{"seq":45,"type":"response","request_seq":46,"command":"runInTerminal","success":true,"body":{"processId":123456}}`
var runInTerminalResponseStruct = RunInTerminalResponse{
	Response: *newResponse(45, 46, "runInTerminal", true),
	Body:     RunInTerminalResponseBody{ProcessId: 123456},
}

// -------- Initialize

var initializeRequestString = `{"seq":1,"type":"request","command":"initialize","arguments":{"clientID":"vscode","clientName":"Visual Studio Code","adapterID":"go","pathFormat":"path","linesStartAt1":true,"columnsStartAt1":true,"supportsVariableType":true,"supportsVariablePaging":true,"supportsRunInTerminalRequest":true,"locale":"en-us"}}`
var initializeRequestOmitDefaultsString = `{"seq":1,"type":"request","command":"initialize","arguments":{"clientID":"vscode","clientName":"Visual Studio Code","adapterID":"go","supportsVariableType":true,"supportsVariablePaging":true,"supportsRunInTerminalRequest":true,"locale":"en-us"}}`
var initializeRequestStruct = InitializeRequest{
	Request: *newRequest(1, "initialize"),
	Arguments: InitializeRequestArguments{
		ClientID:                     "vscode",
		ClientName:                   "Visual Studio Code",
		AdapterID:                    "go",
		PathFormat:                   "path",
		LinesStartAt1:                true,
		ColumnsStartAt1:              true,
		SupportsVariableType:         true,
		SupportsVariablePaging:       true,
		SupportsRunInTerminalRequest: true,
		Locale:                       "en-us",
	},
}

var initializeRequestNotDefaultsString = `{"seq":1,"type":"request","command":"initialize","arguments":{"clientID":"vscode","clientName":"Visual Studio Code","adapterID":"go","pathFormat":"url","linesStartAt1":false,"columnsStartAt1":false,"supportsVariableType":true,"supportsVariablePaging":true,"supportsRunInTerminalRequest":true,"locale":"en-us"}}`
var initializeRequestNotDefaultsStruct = InitializeRequest{
	Request: *newRequest(1, "initialize"),
	Arguments: InitializeRequestArguments{
		ClientID:                     "vscode",
		ClientName:                   "Visual Studio Code",
		AdapterID:                    "go",
		PathFormat:                   "url",
		LinesStartAt1:                false,
		ColumnsStartAt1:              false,
		SupportsVariableType:         true,
		SupportsVariablePaging:       true,
		SupportsRunInTerminalRequest: true,
		Locale:                       "en-us",
	},
}

var initializeResponseString = `{"seq":1,"type":"response","request_seq":2,"command":"initialize","success":true,"body":{"supportsConfigurationDoneRequest":true,"supportsSetVariable":true}}`
var initializeResponseStruct = InitializeResponse{
	Response: *newResponse(1, 2, "initialize", true),
	Body: Capabilities{
		SupportsConfigurationDoneRequest: true,
		SupportsSetVariable:              true,
	},
}

// -------- ConfigurationDone

var configurationDoneRequestString = `{"seq":2,"type":"request","command":"configurationDone"}`
var configurationDoneRequestStruct = ConfigurationDoneRequest{
	Request:   *newRequest(2, "configurationDone"),
	Arguments: nil,
}

var configurationDoneResponseString = `{"seq":2,"type":"response","request_seq":3,"command":"configurationDone","success":true}`
var configurationDoneResponseStruct = ConfigurationDoneResponse{Response: *newResponse(2, 3, "configurationDone", true)}

// -------- Launch

var launchRequestString = `{"seq":3,"type":"request","command":"launch","arguments":{"noDebug":true,"name":"Launch"}}`
var launchRequestStruct = LaunchRequest{
	Request:   *newRequest(3, "launch"),
	Arguments: json.RawMessage(`{"noDebug":true,"name":"Launch"}`),
}

var launchResponseString = `{"seq":3,"type":"response","request_seq":4,"command":"launch","success":true}`
var launchResponseStruct = LaunchResponse{Response: *newResponse(3, 4, "launch", true)}

// -------- Attach

var attachRequestString = `{"seq":4,"type":"request","command":"attach","arguments":{}}`
var attachRequestStruct = AttachRequest{
	Request:   *newRequest(4, "attach"),
	Arguments: json.RawMessage(`{}`),
}

var attachResponseString = `{"seq":4,"type":"response","request_seq":5,"command":"attach","success":true}`
var attachResponseStruct = AttachResponse{Response: *newResponse(4, 5, "attach", true)}

// -------- Restart

var restartRequestString = `{"seq":5,"type":"request","command":"restart"}`
var restartRequestStruct = RestartRequest{
	Request:   *newRequest(5, "restart"),
	Arguments: nil,
}

var restartResponseString = `{"seq":5,"type":"response","request_seq":6,"command":"restart","success":true}`
var restartResponseStruct = RestartResponse{Response: *newResponse(5, 6, "restart", true)}

// -------- Disconnect

var disconnectRequestString = `{"seq":6,"type":"request","command":"disconnect","arguments":{"restart":true}}`
var disconnectRequestStruct = DisconnectRequest{
	Request:   *newRequest(6, "disconnect"),
	Arguments: &DisconnectArguments{Restart: true},
}

var disconnectResponseString = `{"seq":6,"type":"response","request_seq":7,"command":"disconnect","success":true}`
var disconnectResponseStruct = DisconnectResponse{Response: *newResponse(6, 7, "disconnect", true)}

// -------- Terminate

var terminateRequestString = `{"seq":7,"type":"request","command":"terminate","arguments":{"restart":true}}`
var terminateRequestStruct = TerminateRequest{
	Request:   *newRequest(7, "terminate"),
	Arguments: &TerminateArguments{Restart: true},
}

var terminateResponseString = `{"seq":7,"type":"response","request_seq":8,"command":"terminate","success":true}`
var terminateResponseStruct = TerminateResponse{Response: *newResponse(7, 8, "terminate", true)}

// -------- BreakpointLocations

var breakpointLocationsRequestString = `{"seq":8,"type":"request","command":"breakpointLocations","arguments":{"source":{"name":"hello.go","path":"/Users/foo/go/src/hello/hello.go"},"line":10}}`
var breakpointLocationsRequestStruct = BreakpointLocationsRequest{
	Request: *newRequest(8, "breakpointLocations"),
	Arguments: &BreakpointLocationsArguments{
		Source: Source{Name: "hello.go", Path: "/Users/foo/go/src/hello/hello.go"},
		Line:   10,
	},
}

var breakpointLocationsResponseString = `{"seq":8,"type":"response","request_seq":9,"command":"breakpointLocations","success":true,"body":{"breakpoints":[{"line":14}]}}`
var breakpointLocationsResponseStruct = BreakpointLocationsResponse{
	Response: *newResponse(8, 9, "breakpointLocations", true),
	Body: BreakpointLocationsResponseBody{
		Breakpoints: []BreakpointLocation{{Line: 14}},
	},
}

// -------- SetBreakpoints

var setBreakpointsRequestString = `{"seq":9,"type":"request","command":"setBreakpoints","arguments":{"source":{"name":"hello.go","path":"/Users/foo/go/src/hello/hello.go"},"lines":[14],"breakpoints":[{"line":14}],"sourceModified":false}}`
var setBreakpointsRequestStruct = SetBreakpointsRequest{
	Request: *newRequest(9, "setBreakpoints"),
	Arguments: SetBreakpointsArguments{
		Source:         Source{Name: "hello.go", Path: "/Users/foo/go/src/hello/hello.go"},
		Breakpoints:    []SourceBreakpoint{{Line: 14}},
		Lines:          []int{14},
		SourceModified: false,
	},
}

var setBreakpointsResponseString = `{"seq":9,"type":"response","request_seq":10,"command":"setBreakpoints","success":true,"body":{"breakpoints":[{"verified":true,"line":14}]}}`
var setBreakpointsResponseStruct = SetBreakpointsResponse{
	Response: *newResponse(9, 10, "setBreakpoints", true),
	Body: SetBreakpointsResponseBody{
		Breakpoints: []Breakpoint{{Verified: true, Line: 14}},
	},
}

// -------- SetFunctionBreakpoints

var setFunctionBreakpointsRequestString = `{"seq":10,"type":"request","command":"setFunctionBreakpoints","arguments":{"breakpoints":[{"name":"blah"}]}}`
var setFunctionBreakpointsRequestStruct = SetFunctionBreakpointsRequest{
	Request: *newRequest(10, "setFunctionBreakpoints"),
	Arguments: SetFunctionBreakpointsArguments{
		Breakpoints: []FunctionBreakpoint{{Name: "blah"}},
	},
}

var setFunctionBreakpointsResponseString = `{"seq":10,"type":"response","request_seq":11,"command":"setFunctionBreakpoints","success":true, "body":{"breakpoints":[{"verified":true,"line":20}]}}`
var setFunctionBreakpointsResponseStruct = SetFunctionBreakpointsResponse{
	Response: *newResponse(10, 11, "setFunctionBreakpoints", true),
	Body: SetFunctionBreakpointsResponseBody{
		Breakpoints: []Breakpoint{{Verified: true, Line: 20}},
	},
}

// -------- SetExceptionBreakpoints

var setExceptionBreakpointsRequestString = `{"seq":11,"type":"request","command":"setExceptionBreakpoints","arguments":{"filters":[]}}`
var setExceptionBreakpointsRequestStruct = SetExceptionBreakpointsRequest{
	Request:   *newRequest(11, "setExceptionBreakpoints"),
	Arguments: SetExceptionBreakpointsArguments{Filters: []string{}},
}

var setExceptionBreakpointsResponseString = `{"seq":11,"type":"response","request_seq":12,"command":"setExceptionBreakpoints","success":true}`
var setExceptionBreakpointsResponseStruct = SetExceptionBreakpointsResponse{Response: *newResponse(11, 12, "setExceptionBreakpoints", true)}

// -------- DataBreakpointInfo

var dataBreakpointInfoRequestString = `{"seq":12,"type":"request","command":"dataBreakpointInfo","arguments":{"name":"fuzzybunny"}}`
var dataBreakpointInfoRequestStruct = DataBreakpointInfoRequest{
	Request:   *newRequest(12, "dataBreakpointInfo"),
	Arguments: DataBreakpointInfoArguments{Name: "fuzzybunny"},
}

var dataBreakpointInfoResponseString = `{"seq":12,"type":"response","request_seq":13,"command":"dataBreakpointInfo","success":true, "body":{"dataId":null,"description":"some description"}}`
var dataBreakpointInfoResponseStruct = DataBreakpointInfoResponse{
	Response: *newResponse(12, 13, "dataBreakpointInfo", true),
	Body: DataBreakpointInfoResponseBody{
		DataId:      nil,
		Description: "some description",
	},
}

// -------- SetDataBreakpoints

var setDataBreakpointsRequestString = `{"seq":13,"type":"request","command":"setDataBreakpoints","arguments":{"breakpoints":[{"dataId":"dataid"}]}}`
var setDataBreakpointsRequestStruct = SetDataBreakpointsRequest{
	Request: *newRequest(13, "setDataBreakpoints"),
	Arguments: SetDataBreakpointsArguments{
		Breakpoints: []DataBreakpoint{{DataId: "dataid"}},
	},
}

var setDataBreakpointsResponseString = `{"seq":13,"type":"response","request_seq":14,"command":"setDataBreakpoints","success":true, "body":{"breakpoints":[{"verified":true,"line":100}]}}`
var setDataBreakpointsResponseStruct = SetDataBreakpointsResponse{
	Response: *newResponse(13, 14, "setDataBreakpoints", true),
	Body: SetDataBreakpointsResponseBody{
		Breakpoints: []Breakpoint{{Verified: true, Line: 100}},
	},
}

// -------- Continue

var continueRequestString = `{"seq":14,"type":"request","command":"continue","arguments":{"threadId":1}}`
var continueRequestStruct = ContinueRequest{
	Request:   *newRequest(14, "continue"),
	Arguments: ContinueArguments{ThreadId: 1},
}

var continueResponseString = `{"seq":14,"type":"response","request_seq":15,"command":"continue","success":true, "body":{"allThreadsContinued": true}}`
var continueResponseStruct = ContinueResponse{
	Response: *newResponse(14, 15, "continue", true),
	Body:     ContinueResponseBody{AllThreadsContinued: true},
}

// -------- Next

var nextRequestString = `{"seq":15,"type":"request","command":"next","arguments":{"threadId":1}}`
var nextRequestStruct = NextRequest{
	Request:   *newRequest(15, "next"),
	Arguments: NextArguments{ThreadId: 1},
}

var nextResponseString = `{"seq":15,"type":"response","request_seq":16,"command":"next","success":true}`
var nextResponseStruct = NextResponse{Response: *newResponse(15, 16, "next", true)}

// -------- StepIn

var stepInRequestString = `{"seq":16,"type":"request","command":"stepIn","arguments":{"threadId":1}}`
var stepInRequestStruct = StepInRequest{
	Request:   *newRequest(16, "stepIn"),
	Arguments: StepInArguments{ThreadId: 1},
}

var stepInResponseString = `{"seq":16,"type":"response","request_seq":17,"command":"stepIn","success":true}`
var stepInResponseStruct = StepInResponse{Response: *newResponse(16, 17, "stepIn", true)}

// -------- StepOut

var stepOutRequestString = `{"seq":17,"type":"request","command":"stepOut","arguments":{"threadId":1}}`
var stepOutRequestStruct = StepOutRequest{
	Request:   *newRequest(17, "stepOut"),
	Arguments: StepOutArguments{ThreadId: 1},
}

var stepOutResponseString = `{"seq":17,"type":"response","request_seq":18,"command":"stepOut","success":true}`
var stepOutResponseStruct = StepOutResponse{Response: *newResponse(17, 18, "stepOut", true)}

// -------- StepBack

var stepBackRequestString = `{"seq":18,"type":"request","command":"stepBack","arguments":{"threadId":1}}`
var stepBackRequestStruct = StepBackRequest{
	Request:   *newRequest(18, "stepBack"),
	Arguments: StepBackArguments{ThreadId: 1},
}

var stepBackResponseString = `{"seq":18,"type":"response","request_seq":19,"command":"stepBack","success":true}`
var stepBackResponseStruct = StepBackResponse{Response: *newResponse(18, 19, "stepBack", true)}

// -------- ReverseContinue

var reverseContinueRequestString = `{"seq":19,"type":"request","command":"reverseContinue","arguments":{"threadId":1}}`
var reverseContinueRequestStruct = ReverseContinueRequest{
	Request:   *newRequest(19, "reverseContinue"),
	Arguments: ReverseContinueArguments{ThreadId: 1},
}

var reverseContinueResponseString = `{"seq":19,"type":"response","request_seq":20,"command":"reverseContinue","success":true}`
var reverseContinueResponseStruct = ReverseContinueResponse{Response: *newResponse(19, 20, "reverseContinue", true)}

// -------- RestartFrame

var restartFrameRequestString = `{"seq":20,"type":"request","command":"restartFrame","arguments":{"frameId":5}}`
var restartFrameRequestStruct = RestartFrameRequest{
	Request:   *newRequest(20, "restartFrame"),
	Arguments: RestartFrameArguments{FrameId: 5},
}

var restartFrameResponseString = `{"seq":20,"type":"response","request_seq":21,"command":"restartFrame","success":true}`
var restartFrameResponseStruct = RestartFrameResponse{Response: *newResponse(20, 21, "restartFrame", true)}

// -------- Goto

var gotoRequestString = `{"seq":21,"type":"request","command":"goto","arguments":{"threadId":1,"targetId":2}}`
var gotoRequestStruct = GotoRequest{
	Request:   *newRequest(21, "goto"),
	Arguments: GotoArguments{ThreadId: 1, TargetId: 2},
}

var gotoResponseString = `{"seq":21,"type":"response","request_seq":22,"command":"goto","success":true}`
var gotoResponseStruct = GotoResponse{Response: *newResponse(21, 22, "goto", true)}

// -------- Pause

var pauseRequestString = `{"seq":22,"type":"request","command":"pause","arguments":{"threadId":1}}`
var pauseRequestStruct = PauseRequest{
	Request:   *newRequest(22, "pause"),
	Arguments: PauseArguments{ThreadId: 1},
}

var pauseResponseString = `{"seq":22,"type":"response","request_seq":23,"command":"pause","success":true}`
var pauseResponseStruct = PauseResponse{Response: *newResponse(22, 23, "pause", true)}

// -------- StackTrace

var stackTraceRequestString = `{"seq":23,"type":"request","command":"stackTrace","arguments":{"threadId":1,"startFrame":0,"levels":20}}`
var stackTraceRequestStruct = StackTraceRequest{
	Request: *newRequest(23, "stackTrace"),
	Arguments: StackTraceArguments{
		ThreadId:   1,
		StartFrame: 0,
		Levels:     20,
	},
}

var stackTraceResponseString = `{"seq":23,"type":"response","request_seq":24,"command":"stackTrace","success":true,"body":{"stackFrames":[{"id":1000,"source":{"name":"hello.go","path":"/Users/foo/go/src/hello/hello.go","sourceReference":0},"line":6,"column":0,"name":"main.main"},{"id":1001,"source":{"name":"proc.go","path":"/usr/local/go/src/runtime/proc.go","sourceReference":0},"line":203,"column":0,"name":"runtime.main"},{"id":1002,"source":{"name":"asm_amd64.s","path":"/usr/local/go/src/runtime/asm_amd64.s","sourceReference":0},"line":1357,"column":0,"name":"runtime.goexit"}],"totalFrames":3}}`
var stackTraceResponseStruct = StackTraceResponse{
	Response: *newResponse(23, 24, "stackTrace", true),
	Body: StackTraceResponseBody{
		StackFrames: []StackFrame{
			{
				Id: 1000,
				Source: &Source{
					Name:            "hello.go",
					Path:            "/Users/foo/go/src/hello/hello.go",
					SourceReference: 0,
				},
				Line:   6,
				Column: 0,
				Name:   "main.main",
			},
			{
				Id: 1001,
				Source: &Source{
					Name:            "proc.go",
					Path:            "/usr/local/go/src/runtime/proc.go",
					SourceReference: 0,
				},
				Line:   203,
				Column: 0,
				Name:   "runtime.main",
			},
			{
				Id: 1002,
				Source: &Source{
					Name:            "asm_amd64.s",
					Path:            "/usr/local/go/src/runtime/asm_amd64.s",
					SourceReference: 0,
				},
				Line:   1357,
				Column: 0,
				Name:   "runtime.goexit",
			},
		},
		TotalFrames: 3,
	},
}

// -------- Scopes

var scopesRequestString = `{"seq":24,"type":"request","command":"scopes","arguments":{"frameId":1000}}`
var scopesRequestStruct = ScopesRequest{
	Request:   *newRequest(24, "scopes"),
	Arguments: ScopesArguments{FrameId: 1000},
}

var scopesResponseString = `{"seq":24,"type":"response","request_seq":25,"command":"scopes","success":true,"body":{"scopes":[{"name":"Local","variablesReference":1000,"expensive":false},{"name":"Global","variablesReference":1001,"expensive":false}]}}`
var scopesResponseStruct = ScopesResponse{
	Response: *newResponse(24, 25, "scopes", true),
	Body: ScopesResponseBody{
		Scopes: []Scope{
			{
				Name:               "Local",
				VariablesReference: 1000,
				Expensive:          false,
			},
			{
				Name:               "Global",
				VariablesReference: 1001,
				Expensive:          false,
			},
		},
	},
}

// -------- Variables

var variablesRequestString = `{"seq":25,"type":"request","command":"variables","arguments":{"variablesReference":1001}}`
var variablesRequestStruct = VariablesRequest{
	Request:   *newRequest(25, "variables"),
	Arguments: VariablesArguments{VariablesReference: 1001},
}

var variablesResponseString = `{"seq":25,"type":"response","request_seq":26,"command":"variables","success":true,"body":{"variables":[{"name":"x","value":"824634220368","evaluateName":"x","variablesReference":0}]}}`
var variablesResponseStruct = VariablesResponse{
	Response: *newResponse(25, 26, "variables", true),
	Body: VariablesResponseBody{
		Variables: []Variable{
			{
				Name:               "x",
				Value:              "824634220368",
				EvaluateName:       "x",
				VariablesReference: 0,
			},
		},
	},
}

// -------- SetVariable

var setVariableRequestString = `{"seq":26,"type":"request","command":"setVariable","arguments":{"variablesReference":1008,"name":"x","value":"55"}}`
var setVariableRequestStruct = SetVariableRequest{
	Request: *newRequest(26, "setVariable"),
	Arguments: SetVariableArguments{
		VariablesReference: 1008,
		Name:               "x",
		Value:              "55",
	},
}

var setVariableResponseString = `{"seq":26,"type":"response","request_seq":27,"command":"setVariable","success":true,"body":{"value":"55"}}`
var setVariableResponseStruct = SetVariableResponse{
	Response: *newResponse(26, 27, "setVariable", true),
	Body:     SetVariableResponseBody{Value: "55"},
}

// -------- Source

var sourceRequestString = `{"seq":27,"type":"request","command":"source","arguments":{"sourceReference":123}}`
var sourceRequestStruct = SourceRequest{
	Request:   *newRequest(27, "source"),
	Arguments: SourceArguments{SourceReference: 123},
}

var sourceResponseString = `{"seq":27,"type":"response","request_seq":28,"command":"source","success":true,"body":{"content":"somecontent"}}`
var sourceResponseStruct = SourceResponse{
	Response: *newResponse(27, 28, "source", true),
	Body:     SourceResponseBody{Content: "somecontent"},
}

// -------- Threads

var threadsRequestString = `{"seq":28,"type":"request","command":"threads","arguments":{}}`
var threadsRequestStruct = ThreadsRequest{Request: *newRequest(28, "threads")}

var threadsResponseString = `{"seq":28,"type":"response","request_seq":29,"command":"threads","success":true,"body":{"threads":[{"id":1,"name":"Dummy"}]}}`
var threadsResponseStruct = ThreadsResponse{
	Response: *newResponse(28, 29, "threads", true),
	Body:     ThreadsResponseBody{Threads: []Thread{{Id: 1, Name: "Dummy"}}},
}

// -------- TerminateThreads

var terminateThreadsRequestString = `{"seq":29,"type":"request","command":"terminateThreads","arguments":{"threadIds":[1]}}`
var terminateThreadsRequestStruct = TerminateThreadsRequest{
	Request:   *newRequest(29, "terminateThreads"),
	Arguments: TerminateThreadsArguments{ThreadIds: []int{1}},
}

var terminateThreadsResponseString = `{"seq":29,"type":"response","request_seq":30,"command":"terminateThreads","success":true}`
var terminateThreadsResponseStruct = TerminateThreadsResponse{Response: *newResponse(29, 30, "terminateThreads", true)}

// -------- Modules

var modulesRequestString = `{"seq":30,"type":"request","command":"modules","arguments":{"startModule":1,"moduleCount":3}}`
var modulesRequestStruct = ModulesRequest{
	Request: *newRequest(30, "modules"),
	Arguments: ModulesArguments{
		StartModule: 1,
		ModuleCount: 3,
	},
}

var modulesResponseString = `{"seq":30,"type":"response","request_seq":31,"command":"modules","success":true,"body":{"totalModules":2,"modules":[{"id":1,"name":"one"}]}}`
var modulesResponseStruct = ModulesResponse{
	Response: *newResponse(30, 31, "modules", true),
	Body: ModulesResponseBody{
		TotalModules: 2,
		Modules:      []Module{{Id: 1.0, Name: "one"}},
	},
}

// -------- LoadedSources

var loadedSourcesRequestString = `{"seq":31,"type":"request","command":"loadedSources"}`
var loadedSourcesRequestStruct = LoadedSourcesRequest{
	Request:   *newRequest(31, "loadedSources"),
	Arguments: nil,
}

var loadedSourcesResponseString = `{"seq":31,"type":"response","request_seq":32,"command":"loadedSources","success":true,"body":{"sources":[{"name":"hello.go","path":"/Users/foo/go/src/hello/hello.go"}]}}`
var loadedSourcesResponseStruct = LoadedSourcesResponse{
	Response: *newResponse(31, 32, "loadedSources", true),
	Body: LoadedSourcesResponseBody{
		Sources: []Source{
			{
				Name: "hello.go",
				Path: "/Users/foo/go/src/hello/hello.go",
			},
		},
	},
}

// -------- Evaluate

var evaluateRequestString = `{"seq":32,"type":"request","command":"evaluate","arguments":{"expression":"x==1","frameId":1000,"context":"repl"}}`
var evaluateRequestStruct = EvaluateRequest{
	Request: *newRequest(32, "evaluate"),
	Arguments: EvaluateArguments{
		Expression: "x==1",
		FrameId:    1000,
		Context:    "repl",
	},
}

var evaluateResponseString = `{"seq":32,"type":"response","request_seq":33,"command":"evaluate","success":true,"body":{"result":"false","variablesReference":1}}`
var evaluateResponseStruct = EvaluateResponse{
	Response: *newResponse(32, 33, "evaluate", true),
	Body: EvaluateResponseBody{
		Result:             "false",
		VariablesReference: 1,
	},
}

// -------- SetExpression

var setExpressionRequestString = `{"seq":33,"type":"request","command":"setExpression","arguments":{"expression":"x==1","value":"true"}}`
var setExpressionRequestStruct = SetExpressionRequest{
	Request: *newRequest(33, "setExpression"),
	Arguments: SetExpressionArguments{
		Expression: "x==1",
		Value:      "true",
	},
}

var setExpressionResponseString = `{"seq":33,"type":"response","request_seq":34,"command":"setExpression","success":true,"body":{"value":"true"}}`
var setExpressionResponseStruct = SetExpressionResponse{
	Response: *newResponse(33, 34, "setExpression", true),
	Body:     SetExpressionResponseBody{Value: "true"},
}

// -------- StepInTargets

var stepInTargetsRequestString = `{"seq":34,"type":"request","command":"stepInTargets","arguments":{"frameId":1000}}`
var stepInTargetsRequestStruct = StepInTargetsRequest{
	Request: *newRequest(34, "stepInTargets"),
	Arguments: StepInTargetsArguments{
		FrameId: 1000,
	},
}

var stepInTargetsResponseString = `{"seq":34,"type":"response","request_seq":35,"command":"stepInTargets","success":true,"body":{"targets":[{"id":123,"label":"somelabel"}]}}`
var stepInTargetsResponseStruct = StepInTargetsResponse{
	Response: *newResponse(34, 35, "stepInTargets", true),
	Body: StepInTargetsResponseBody{
		Targets: []StepInTarget{
			{Id: 123, Label: "somelabel"},
		},
	},
}

// -------- GotoTargets

var gotoTargetsRequestString = `{"seq":35,"type":"request","command":"gotoTargets","arguments":{"source":{"name":"hello.go","path":"/Users/foo/go/src/hello/hello.go"},"line":10}}`
var gotoTargetsRequestStruct = GotoTargetsRequest{
	Request: *newRequest(35, "gotoTargets"),
	Arguments: GotoTargetsArguments{
		Source: Source{Name: "hello.go", Path: "/Users/foo/go/src/hello/hello.go"},
		Line:   10,
	},
}

var gotoTargetsResponseString = `{"seq":35,"type":"response","request_seq":36,"command":"gotoTargets","success":true,"body":{"targets":[{"id":123,"label":"somelabel","line":10}]}}`
var gotoTargetsResponseStruct = GotoTargetsResponse{
	Response: *newResponse(35, 36, "gotoTargets", true),
	Body: GotoTargetsResponseBody{
		Targets: []GotoTarget{
			{Id: 123, Label: "somelabel", Line: 10},
		},
	},
}

// -------- Completions

var completionsRequestString = `{"seq": 36,"type":"request","command":"completions","arguments":{"text":"sometext","column":123}}`
var completionsRequestStruct = CompletionsRequest{
	Request:   *newRequest(36, "completions"),
	Arguments: CompletionsArguments{Text: "sometext", Column: 123},
}

var completionsResponseString = `{"seq":36,"type":"response","request_seq":37,"command":"completions","success":true,"body":{"targets":[{"label":"somelabel"}]}}`
var completionsResponseStruct = CompletionsResponse{
	Response: *newResponse(36, 37, "completions", true),
	Body: CompletionsResponseBody{
		Targets: []CompletionItem{
			{Label: "somelabel"},
		},
	},
}

// -------- ExceptionInfo

var exceptionInfoRequestString = `{"seq":36,"type":"request","command":"exceptionInfo","arguments":{"threadId":1}}`
var exceptionInfoRequestStruct = ExceptionInfoRequest{
	Request:   *newRequest(36, "exceptionInfo"),
	Arguments: ExceptionInfoArguments{ThreadId: 1},
}

var exceptionInfoResponseString = `{"seq":36,"type":"response","request_seq":37,"command":"exceptionInfo","success":true,"body":{"exceptionId":"someid","breakMode":"somebreakmode"}}`
var exceptionInfoResponseStruct = ExceptionInfoResponse{
	Response: *newResponse(36, 37, "exceptionInfo", true),
	Body: ExceptionInfoResponseBody{
		ExceptionId: "someid",
		BreakMode:   "somebreakmode",
	},
}

// -------- ReadMemory

var readMemoryRequestString = `{"seq":37,"type":"request","command":"readMemory","arguments":{"memoryReference":"someref","count":123}}`
var readMemoryRequestStruct = ReadMemoryRequest{
	Request: *newRequest(37, "readMemory"),
	Arguments: ReadMemoryArguments{
		MemoryReference: "someref",
		Count:           123,
	},
}

var readMemoryResponseString = `{"seq":37,"type":"response","request_seq":38,"command":"readMemory","success":true,"body":{"address":"someaddr"}}`
var readMemoryResponseStruct = ReadMemoryResponse{
	Response: *newResponse(37, 38, "readMemory", true),
	Body:     ReadMemoryResponseBody{Address: "someaddr"},
}

// -------- Disassemble

var disassembleRequestString = `{"seq":38,"type":"request","command":"disassemble","arguments":{"memoryReference":"someref","instructionCount":123}}`
var disassembleRequestStruct = DisassembleRequest{
	Request: *newRequest(38, "disassemble"),
	Arguments: DisassembleArguments{
		MemoryReference:  "someref",
		InstructionCount: 123,
	},
}

var disassembleResponseString = `{"seq":38,"type":"response","request_seq":39,"command":"disassemble","success":true,"body":{"instructions":[{"address":"someaddr","instruction":"someinstr"}]}}`
var disassembleResponseStruct = DisassembleResponse{
	Response: *newResponse(38, 39, "disassemble", true),
	Body: DisassembleResponseBody{
		Instructions: []DisassembledInstruction{
			{
				Address:     "someaddr",
				Instruction: "someinstr",
			},
		},
	},
}

// -------- StartDebugging

var startDebuggingRequestString = `{"seq":39,"type":"request","command":"startDebugging","arguments":{"request":"launch","configuration":{"any":true}}}`
var startDebuggingRequestStruct = StartDebuggingRequest{
	Request: *newRequest(39, "startDebugging"),
	Arguments: StartDebuggingRequestArguments{
		Request:       "launch",
		Configuration: map[string]interface{}{"any": true},
	},
}

var startDebuggingResponseString = `{"seq":39,"type":"response","request_seq":40,"command":"startDebugging","success":true}`
var startDebuggingResponseStruct = StartDebuggingResponse{
	Response: *newResponse(39, 40, "startDebugging", true),
}

// -------- Events --------

var initializedEventString = `{"seq":1,"type":"event","event":"initialized"}`
var initializedEventStruct = InitializedEvent{
	Event: *newEvent(1, "initialized"),
}

var stoppedEventString = `{"seq":2,"type":"event","event":"stopped","body":{"reason":"breakpoint","threadId":1,"allThreadsStopped":true}}`
var stoppedEventStruct = StoppedEvent{
	Event: *newEvent(2, "stopped"),
	Body:  StoppedEventBody{Reason: "breakpoint", ThreadId: 1, AllThreadsStopped: true},
}

var continuedEventString = `{"seq":3,"type":"event","event":"continued","body":{"threadId":123}}`
var continuedEventStruct = ContinuedEvent{
	Event: *newEvent(3, "continued"),
	Body:  ContinuedEventBody{ThreadId: 123},
}

var exitedEventString = `{"seq":4,"type":"event","event":"exited","body":{"exitCode":123}}`
var exitedEventStruct = ExitedEvent{
	Event: *newEvent(4, "exited"),
	Body:  ExitedEventBody{ExitCode: 123},
}

var terminatedEventString = `{"seq":5,"type":"event","event":"terminated","body":{"restart":true}}`
var terminatedEventStruct = TerminatedEvent{
	Event: *newEvent(5, "terminated"),
	Body:  TerminatedEventBody{Restart: true},
}

var threadEventString = `{"seq":6,"type":"event","event":"thread","body":{"reason":"started","threadId":18}}`
var threadEventStruct = ThreadEvent{
	Event: *newEvent(6, "thread"),
	Body:  ThreadEventBody{Reason: "started", ThreadId: 18},
}

var outputEventString = `{"seq":7,"type":"event","event":"output","body":{"category":"stdout","output":"something that got logged"}}`
var outputEventStruct = OutputEvent{
	Event: *newEvent(7, "output"),
	Body:  OutputEventBody{Category: "stdout", Output: "something that got logged"},
}

var breakpointEventString = `{"seq":8,"type":"event","event":"breakpoint","body":{"reason":"new","breakpoint":{"verified":true}}}`
var breakpointEventStruct = BreakpointEvent{
	Event: *newEvent(8, "breakpoint"),
	Body:  BreakpointEventBody{Reason: "new", Breakpoint: Breakpoint{Verified: true}},
}

var moduleEventString = `{"seq":9,"type":"event","event":"module","body":{"reason":"removed","module":{"id":"id"}}}`
var moduleEventStruct = ModuleEvent{
	Event: *newEvent(9, "module"),
	Body:  ModuleEventBody{Reason: "removed", Module: Module{Id: "id"}},
}

var loadedSourceEventString = `{"seq":10,"type":"event","event":"loadedSource","body":{"reason":"changed","source":{"name":"hello.go","path":"/Users/foo/go/src/hello/hello.go"}}}`
var loadedSourceEventStruct = LoadedSourceEvent{
	Event: *newEvent(10, "loadedSource"),
	Body:  LoadedSourceEventBody{Reason: "changed", Source: Source{Name: "hello.go", Path: "/Users/foo/go/src/hello/hello.go"}},
}

var processEventString = `{"seq":11,"type":"event","event":"process","body":{"name":"/home/example/myproj/program.js"}}`
var processEventStruct = ProcessEvent{
	Event: *newEvent(11, "process"),
	Body:  ProcessEventBody{Name: "/home/example/myproj/program.js"},
}

var capabilitiesEventString = `{"seq":12,"type":"event","event":"capabilities","body":{"capabilities":{"supportsFunctionBreakpoints":true}}}`
var capabilitiesEventStruct = CapabilitiesEvent{
	Event: *newEvent(12, "capabilities"),
	Body:  CapabilitiesEventBody{Capabilities: Capabilities{SupportsFunctionBreakpoints: true}},
}

func TestDecodeProtocolMessage(t *testing.T) {
	// Sometimes partial messages can be returned on error, but
	// the user should not rely on those and just check err itself.
	// Hence the test will not check those.
	var msgIgnoredOnError Message
	const noError = ""
	tests := []struct {
		data    string
		wantMsg Message
		wantErr string
	}{
		// ProtocolMessage
		{``, msgIgnoredOnError, "unexpected end of JSON input"},
		{`,`, msgIgnoredOnError, "invalid character ',' looking for beginning of value"},
		{`{}`, msgIgnoredOnError, "ProtocolMessage type '' is not supported (seq: 0)"},
		{`{"a": 1}`, msgIgnoredOnError, "ProtocolMessage type '' is not supported (seq: 0)"},
		{`{"type":"foo", "seq": 2}`, msgIgnoredOnError, "ProtocolMessage type 'foo' is not supported (seq: 2)"},
		// Request
		{`{"type":"request"}`, msgIgnoredOnError, "Request command '' is not supported (seq: 0)"},
		{cancelRequestString, &cancelRequestStruct, noError},
		{runInTerminalRequestString, &runInTerminalRequestStruct, noError},
		{initializeRequestString, &initializeRequestStruct, noError},
		{initializeRequestOmitDefaultsString, &initializeRequestStruct, noError},
		{initializeRequestNotDefaultsString, &initializeRequestNotDefaultsStruct, noError},
		{configurationDoneRequestString, &configurationDoneRequestStruct, noError},
		{launchRequestString, &launchRequestStruct, noError},
		{attachRequestString, &attachRequestStruct, noError},
		{restartRequestString, &restartRequestStruct, noError},
		{disconnectRequestString, &disconnectRequestStruct, noError},
		{terminateRequestString, &terminateRequestStruct, noError},
		{breakpointLocationsRequestString, &breakpointLocationsRequestStruct, noError},
		{setBreakpointsRequestString, &setBreakpointsRequestStruct, noError},
		{setFunctionBreakpointsRequestString, &setFunctionBreakpointsRequestStruct, noError},
		{setExceptionBreakpointsRequestString, &setExceptionBreakpointsRequestStruct, noError},
		{dataBreakpointInfoRequestString, &dataBreakpointInfoRequestStruct, noError},
		{setDataBreakpointsRequestString, &setDataBreakpointsRequestStruct, noError},
		{continueRequestString, &continueRequestStruct, noError},
		{nextRequestString, &nextRequestStruct, noError},
		{stepInRequestString, &stepInRequestStruct, noError},
		{stepOutRequestString, &stepOutRequestStruct, noError},
		{stepBackRequestString, &stepBackRequestStruct, noError},
		{reverseContinueRequestString, &reverseContinueRequestStruct, noError},
		{restartFrameRequestString, &restartFrameRequestStruct, noError},
		{gotoRequestString, &gotoRequestStruct, noError},
		{pauseRequestString, &pauseRequestStruct, noError},
		{stackTraceRequestString, &stackTraceRequestStruct, noError},
		{scopesRequestString, &scopesRequestStruct, noError},
		{variablesRequestString, &variablesRequestStruct, noError},
		{setVariableRequestString, &setVariableRequestStruct, noError},
		{sourceRequestString, &sourceRequestStruct, noError},
		{threadsRequestString, &threadsRequestStruct, noError},
		{terminateThreadsRequestString, &terminateThreadsRequestStruct, noError},
		{modulesRequestString, &modulesRequestStruct, noError},
		{loadedSourcesRequestString, &loadedSourcesRequestStruct, noError},
		{evaluateRequestString, &evaluateRequestStruct, noError},
		{setExpressionRequestString, &setExpressionRequestStruct, noError},
		{stepInTargetsRequestString, &stepInTargetsRequestStruct, noError},
		{gotoTargetsRequestString, &gotoTargetsRequestStruct, noError},
		{completionsRequestString, &completionsRequestStruct, noError},
		{exceptionInfoRequestString, &exceptionInfoRequestStruct, noError},
		{readMemoryRequestString, &readMemoryRequestStruct, noError},
		{disassembleRequestString, &disassembleRequestStruct, noError},
		{startDebuggingRequestString, &startDebuggingRequestStruct, noError},
		// Response
		{`{"type":"response","success":true, "seq": 77}`, msgIgnoredOnError, "Response command '' is not supported (seq: 77)"},
		{errorResponseString, &errorResponseStruct, noError},
		{cancelResponseString, &cancelResponseStruct, noError},
		{runInTerminalResponseString, &runInTerminalResponseStruct, noError},
		{initializeResponseString, &initializeResponseStruct, noError},
		{configurationDoneResponseString, &configurationDoneResponseStruct, noError},
		{launchResponseString, &launchResponseStruct, noError},
		{attachResponseString, &attachResponseStruct, noError},
		{restartResponseString, &restartResponseStruct, noError},
		{disconnectResponseString, &disconnectResponseStruct, noError},
		{terminateResponseString, &terminateResponseStruct, noError},
		{breakpointLocationsResponseString, &breakpointLocationsResponseStruct, noError},
		{setBreakpointsResponseString, &setBreakpointsResponseStruct, noError},
		{setFunctionBreakpointsResponseString, &setFunctionBreakpointsResponseStruct, noError},
		{setExceptionBreakpointsResponseString, &setExceptionBreakpointsResponseStruct, noError},
		{dataBreakpointInfoResponseString, &dataBreakpointInfoResponseStruct, noError},
		{setDataBreakpointsResponseString, &setDataBreakpointsResponseStruct, noError},
		{continueResponseString, &continueResponseStruct, noError},
		{nextResponseString, &nextResponseStruct, noError},
		{stepInResponseString, &stepInResponseStruct, noError},
		{stepOutResponseString, &stepOutResponseStruct, noError},
		{stepBackResponseString, &stepBackResponseStruct, noError},
		{reverseContinueResponseString, &reverseContinueResponseStruct, noError},
		{restartFrameResponseString, &restartFrameResponseStruct, noError},
		{gotoResponseString, &gotoResponseStruct, noError},
		{pauseResponseString, &pauseResponseStruct, noError},
		{stackTraceResponseString, &stackTraceResponseStruct, noError},
		{scopesResponseString, &scopesResponseStruct, noError},
		{variablesResponseString, &variablesResponseStruct, noError},
		{setVariableResponseString, &setVariableResponseStruct, noError},
		{sourceResponseString, &sourceResponseStruct, noError},
		{threadsResponseString, &threadsResponseStruct, noError},
		{terminateThreadsResponseString, &terminateThreadsResponseStruct, noError},
		{modulesResponseString, &modulesResponseStruct, noError},
		{loadedSourcesResponseString, &loadedSourcesResponseStruct, noError},
		{evaluateResponseString, &evaluateResponseStruct, noError},
		{setExpressionResponseString, &setExpressionResponseStruct, noError},
		{stepInTargetsResponseString, &stepInTargetsResponseStruct, noError},
		{gotoTargetsResponseString, &gotoTargetsResponseStruct, noError},
		{completionsResponseString, &completionsResponseStruct, noError},
		{exceptionInfoResponseString, &exceptionInfoResponseStruct, noError},
		{readMemoryResponseString, &readMemoryResponseStruct, noError},
		{disassembleResponseString, &disassembleResponseStruct, noError},
		{startDebuggingResponseString, &startDebuggingResponseStruct, noError},
		// Event
		{`{"type":"event", "seq": 8}`, msgIgnoredOnError, "Event event '' is not supported (seq: 8)"},
		{initializedEventString, &initializedEventStruct, noError},
		{stoppedEventString, &stoppedEventStruct, noError},
		{continuedEventString, &continuedEventStruct, noError},
		{exitedEventString, &exitedEventStruct, noError},
		{terminatedEventString, &terminatedEventStruct, noError},
		{threadEventString, &threadEventStruct, noError},
		{outputEventString, &outputEventStruct, noError},
		{breakpointEventString, &breakpointEventStruct, noError},
		{moduleEventString, &moduleEventStruct, noError},
		{loadedSourceEventString, &loadedSourceEventStruct, noError},
		{processEventString, &processEventStruct, noError},
		{capabilitiesEventString, &capabilitiesEventStruct, noError},
	}

	for _, test := range tests {
		t.Run(test.data, func(t *testing.T) {
			msg, err := DecodeProtocolMessage([]byte(test.data))
			if err != nil { // Decoding error
				if err.Error() != test.wantErr { // Was it the right error?
					t.Errorf("got error=%#v, want %q", err, test.wantErr)
				}
			} else { // No decoding error
				if test.wantErr != "" { // Did we expect one?
					t.Errorf("got error=nil, want %#q", test.wantErr)
				}
				got, _ := json.Marshal(msg)
				want, _ := json.Marshal(test.wantMsg)
				if !reflect.DeepEqual(msg, test.wantMsg) { // Check result
					t.Errorf("\ngot message\n%s\nwant\n%s", got, want)
				}
			}
		})
	}
}

// -------- Custom Request/Response and Event --------

type customRequest struct {
	Request
	Body string `json:"body"`
}

func (r *customRequest) GetRequest() *Request { return &r.Request }

type customResponse struct {
	Response
	Body string `json:"body"`
}

func (r *customResponse) GetResponse() *Response { return &r.Response }

var customRequestString = `{"seq":40,"type":"request","command":"customReq","body":"242424"}`
var customRequestStruct = customRequest{
	Request: *newRequest(40, "customReq"),
	Body:    "242424",
}

var customResponseString = `{"seq":40,"type":"response","request_seq":41,"command":"customReq","success":true,"body":"424242"}`
var customResponseStruct = customResponse{
	Response: *newResponse(40, 41, "customReq", true),
	Body:     "424242",
}

type customEvent struct {
	Event
	Body int `json:"body"`
}

func (e *customEvent) GetEvent() *Event { return &e.Event }

var customEventString = `{"seq":13,"type":"event","event":"customEvt","body":42}`
var customEventStruct = customEvent{
	Event: *newEvent(13, "customEvt"),
	Body:  42,
}

func TestDecodeProtocolMessage_Custom(t *testing.T) {
	tests := []struct {
		data    string
		wantMsg Message
	}{
		{customRequestString, &customRequestStruct},
		{customResponseString, &customResponseStruct},
		{customEventString, &customEventStruct},
	}

	codec := NewCodec()
	codec.RegisterRequest("customReq", func() Message { return new(customRequest) }, func() Message { return new(customResponse) })
	codec.RegisterEvent("customEvt", func() Message { return new(customEvent) })
	for _, test := range tests {
		t.Run(test.data, func(t *testing.T) {
			msg, err := codec.DecodeMessage([]byte(test.data))
			if err != nil { // Decoding error
				t.Fatalf("codec.DecodeMessage() failed with %v", err)
			}
			got, _ := json.Marshal(msg)
			want, _ := json.Marshal(test.wantMsg)
			if !reflect.DeepEqual(msg, test.wantMsg) { // Check result
				t.Errorf("\ngot message\n%s\nwant\n%s", got, want)
			}
		})
	}
}

// newRequest builds a Request struct with the specified fields.
func newRequest(seq int, command string) *Request {
	return &Request{
		ProtocolMessage: ProtocolMessage{
			Type: "request",
			Seq:  seq,
		},
		Command: command,
	}
}

// newEvent builds an Event struct with the specified fields.
func newEvent(seq int, event string) *Event {
	return &Event{
		ProtocolMessage: ProtocolMessage{
			Seq:  seq,
			Type: "event",
		},
		Event: event,
	}
}

// newResponse builds a Response struct with the specified fields.
func newResponse(seq int, requestSeq int, command string, success bool) *Response {
	return &Response{
		ProtocolMessage: ProtocolMessage{
			Seq:  seq,
			Type: "response",
		},
		Command:    command,
		RequestSeq: requestSeq,
		Success:    success,
	}
}
