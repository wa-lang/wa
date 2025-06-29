package wazero

import (
	"bytes"
	"context"
	"errors"

	"wa-lang.org/wa/internal/3rdparty/wazero/api"
	experimentalapi "wa-lang.org/wa/internal/3rdparty/wazero/experimental"
	"wa-lang.org/wa/internal/3rdparty/wazero/internal/version"
	"wa-lang.org/wa/internal/3rdparty/wazero/internal/wasm"
	binaryformat "wa-lang.org/wa/internal/3rdparty/wazero/internal/wasm/binary"
)

// Runtime allows embedding of WebAssembly modules.
//
// The below is an example of basic initialization:
//
//	ctx := context.Background()
//	r := wazero.NewRuntime(ctx)
//	defer r.Close(ctx) // This closes everything this Runtime created.
//
//	module, _ := r.InstantiateModuleFromBinary(ctx, wasm)
type Runtime interface {
	// NewHostModuleBuilder lets you create modules out of functions defined in Go.
	//
	// Below defines and instantiates a module named "env" with one function:
	//
	//	ctx := context.Background()
	//	hello := func() {
	//		fmt.Fprintln(stdout, "hello!")
	//	}
	//	_, err := r.NewHostModuleBuilder("env").
	//		NewFunctionBuilder().WithFunc(hello).Export("hello").
	//		Instantiate(ctx, r)
	NewHostModuleBuilder(moduleName string) HostModuleBuilder

	// CompileModule decodes the WebAssembly binary (%.wasm) or errs if invalid.
	// Any pre-compilation done after decoding wasm is dependent on RuntimeConfig.
	//
	// There are two main reasons to use CompileModule instead of InstantiateModuleFromBinary:
	//   - Improve performance when the same module is instantiated multiple times under different names
	//   - Reduce the amount of errors that can occur during InstantiateModule.
	//
	// # Notes
	//
	//   - The resulting module name defaults to what was binary from the custom name section.
	//   - Any pre-compilation done after decoding the source is dependent on RuntimeConfig.
	//
	// See https://www.w3.org/TR/2019/REC-wasm-core-1-20191205/#name-section%E2%91%A0
	CompileModule(ctx context.Context, binary []byte) (CompiledModule, error)

	// InstantiateModuleFromBinary instantiates a module from the WebAssembly binary (%.wasm) or errs if invalid.
	//
	// Here's an example:
	//	ctx := context.Background()
	//	r := wazero.NewRuntime(ctx)
	//	defer r.Close(ctx) // This closes everything this Runtime created.
	//
	//	module, _ := r.InstantiateModuleFromBinary(ctx, wasm)
	//
	// # Notes
	//
	//   - This is a convenience utility that chains CompileModule with InstantiateModule. To instantiate the same
	//	source multiple times, use CompileModule as InstantiateModule avoids redundant decoding and/or compilation.
	//   - To avoid using configuration defaults, use InstantiateModule instead.
	InstantiateModuleFromBinary(ctx context.Context, source []byte) (api.Module, error)

	// Namespace is the default namespace of this runtime, and is embedded for convenience. Most users will only use the
	// default namespace.
	//
	// Advanced use cases can use NewNamespace to redefine modules of the same name. For example, to allow different
	// modules to define their own stateful "env" module.
	Namespace

	// NewNamespace creates an empty namespace which won't conflict with any other namespace including the default.
	// This is more efficient than multiple runtimes, as namespaces share a compiler cache.
	//
	// In simplest case, a namespace won't conflict if another has a module with the same name:
	//	b := assemblyscript.NewBuilder(r)
	//	m1, _ := b.InstantiateModule(ctx, r.NewNamespace(ctx))
	//	m2, _ := b.InstantiateModule(ctx, r.NewNamespace(ctx))
	//
	// This is also useful for different modules that import the same module name (like "env"), but need different
	// configuration or state. For example, one with trace logging enabled and another disabled:
	//	b := assemblyscript.NewBuilder(r)
	//
	//	// m1 has trace logging disabled
	//	ns1 := r.NewNamespace(ctx)
	//	_ = b.InstantiateModule(ctx, ns1)
	//	m1, _ := ns1.InstantiateModule(ctx, compiled, config)
	//
	//	// m2 has trace logging enabled
	//	ns2 := r.NewNamespace(ctx)
	//	_ = b.WithTraceToStdout().InstantiateModule(ctx, ns2)
	//	m2, _ := ns2.InstantiateModule(ctx, compiled, config)
	//
	// # Notes
	//
	//   - The returned namespace does not inherit any modules from the runtime default namespace.
	//   - Closing the returned namespace closes any modules in it.
	//   - Closing this runtime also closes the namespace returned from this function.
	NewNamespace(context.Context) Namespace

	// CloseWithExitCode closes all the modules that have been initialized in this Runtime with the provided exit code.
	// An error is returned if any module returns an error when closed.
	//
	// Here's an example:
	//	ctx := context.Background()
	//	r := wazero.NewRuntime(ctx)
	//	defer r.CloseWithExitCode(ctx, 2) // This closes everything this Runtime created.
	//
	//	// Everything below here can be closed, but will anyway due to above.
	//	_, _ = wasi_snapshot_preview1.InstantiateSnapshotPreview1(ctx, r)
	//	mod, _ := r.InstantiateModuleFromBinary(ctx, wasm)
	CloseWithExitCode(ctx context.Context, exitCode uint32) error

	// Closer closes all namespace and compiled code by delegating to CloseWithExitCode with an exit code of zero.
	api.Closer
}

// NewRuntime returns a runtime with a configuration assigned by NewRuntimeConfig.
func NewRuntime(ctx context.Context) Runtime {
	return NewRuntimeWithConfig(ctx, NewRuntimeConfig())
}

// NewRuntimeWithConfig returns a runtime with the given configuration.
func NewRuntimeWithConfig(ctx context.Context, rConfig RuntimeConfig) Runtime {
	if v := ctx.Value(version.WazeroVersionKey{}); v == nil {
		ctx = context.WithValue(ctx, version.WazeroVersionKey{}, wazeroVersion)
	}
	config := rConfig.(*runtimeConfig)
	store, ns := wasm.NewStore(config.enabledFeatures, config.newEngine(ctx, config.enabledFeatures))
	return &runtime{
		store:                 store,
		ns:                    &namespace{store: store, ns: ns},
		enabledFeatures:       config.enabledFeatures,
		memoryLimitPages:      config.memoryLimitPages,
		memoryCapacityFromMax: config.memoryCapacityFromMax,
		isInterpreter:         config.isInterpreter,
	}
}

// runtime allows decoupling of public interfaces from internal representation.
type runtime struct {
	store                 *wasm.Store
	ns                    *namespace
	enabledFeatures       api.CoreFeatures
	memoryLimitPages      uint32
	memoryCapacityFromMax bool
	isInterpreter         bool
	compiledModules       []*compiledModule
}

// NewNamespace implements Runtime.NewNamespace.
func (r *runtime) NewNamespace(ctx context.Context) Namespace {
	return &namespace{store: r.store, ns: r.store.NewNamespace(ctx)}
}

// Module implements Namespace.Module embedded by Runtime.
func (r *runtime) Module(moduleName string) api.Module {
	return r.ns.Module(moduleName)
}

// CompileModule implements Runtime.CompileModule
func (r *runtime) CompileModule(ctx context.Context, binary []byte) (CompiledModule, error) {
	if binary == nil {
		return nil, errors.New("binary == nil")
	}

	if len(binary) < 4 || !bytes.Equal(binary[0:4], binaryformat.Magic) {
		return nil, errors.New("invalid binary")
	}

	internal, err := binaryformat.DecodeModule(binary, r.enabledFeatures, r.memoryLimitPages, r.memoryCapacityFromMax)
	if err != nil {
		return nil, err
	} else if err = internal.Validate(r.enabledFeatures); err != nil {
		// TODO: decoders should validate before returning, as that allows
		// them to err with the correct position in the wasm binary.
		return nil, err
	}

	internal.AssignModuleID(binary)

	// Now that the module is validated, cache the function and memory definitions.
	internal.BuildFunctionDefinitions()
	internal.BuildMemoryDefinitions()

	c := &compiledModule{module: internal, compiledEngine: r.store.Engine}

	listeners, err := buildListeners(ctx, internal)
	if err != nil {
		return nil, err
	}

	if err = r.store.Engine.CompileModule(ctx, internal, listeners); err != nil {
		return nil, err
	}

	r.compiledModules = append(r.compiledModules, c)
	return c, nil
}

func buildListeners(ctx context.Context, internal *wasm.Module) ([]experimentalapi.FunctionListener, error) {
	// Test to see if internal code are using an experimental feature.
	fnlf := ctx.Value(experimentalapi.FunctionListenerFactoryKey{})
	if fnlf == nil {
		return nil, nil
	}
	factory := fnlf.(experimentalapi.FunctionListenerFactory)
	importCount := internal.ImportFuncCount()
	listeners := make([]experimentalapi.FunctionListener, len(internal.FunctionSection))
	for i := 0; i < len(listeners); i++ {
		listeners[i] = factory.NewListener(internal.FunctionDefinitionSection[uint32(i)+importCount])
	}
	return listeners, nil
}

// InstantiateModuleFromBinary implements Runtime.InstantiateModuleFromBinary
func (r *runtime) InstantiateModuleFromBinary(ctx context.Context, binary []byte) (api.Module, error) {
	if compiled, err := r.CompileModule(ctx, binary); err != nil {
		return nil, err
	} else {
		compiled.(*compiledModule).closeWithModule = true
		return r.InstantiateModule(ctx, compiled, NewModuleConfig())
	}
}

// InstantiateModule implements Namespace.InstantiateModule embedded by Runtime.
func (r *runtime) InstantiateModule(
	ctx context.Context,
	compiled CompiledModule,
	mConfig ModuleConfig,
) (api.Module, error) {
	return r.ns.InstantiateModule(ctx, compiled, mConfig)
}

// Close implements api.Closer embedded in Runtime.
func (r *runtime) Close(ctx context.Context) error {
	return r.CloseWithExitCode(ctx, 0)
}

// CloseWithExitCode implements Runtime.CloseWithExitCode
func (r *runtime) CloseWithExitCode(ctx context.Context, exitCode uint32) error {
	err := r.store.CloseWithExitCode(ctx, exitCode)
	for _, c := range r.compiledModules {
		if e := c.Close(ctx); e != nil && err == nil {
			err = e
		}
	}
	return err
}
