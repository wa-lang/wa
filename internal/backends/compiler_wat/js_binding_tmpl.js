{{- $root := . -}}
// Auto generated by Wa Compiler, DONOT EDIT!!!

class WaApp {
  constructor() {
    this._wasm_url = "./{{$root.Filename}}";
    this._mem_util = this._createMemUtil();
    this._wasm_inst = null;
    this._wa_print_buf = "";
  }

  _createMemUtil = () => {
    return {
      mem: () => { return this._wasm_inst.exports.memory; },
      mem_view: (addr, len) => { return new DataView(this._mem_util.mem().buffer, addr, len); },
      mem_array_u8: (addr, len) => { return new Uint8Array(this._mem_util.mem().buffer, addr, len); },
      mem_array_u16: (addr, len) => { return new Uint16Array(this._mem_util.mem().buffer, addr, len); },
      mem_array_u32: (addr, len) => { return new Uint32Array(this._mem_util.mem().buffer, addr, len); },
      mem_array_i32: (addr, len) => { return new Int32Array(this._mem_util.mem().buffer, addr, len); },
      mem_array_f32: (addr, len) => { return new Float32Array(this._mem_util.mem().buffer, addr, len); },
      mem_array_f64: (addr, len) => { return new Float64Array(this._mem_util.mem().buffer, addr, len); },
      get_string: (d, l) => { return new TextDecoder("utf-8").decode(this._mem_util.mem_view(d, l)); },
      set_string: (s) => {
        const bytes = new TextEncoder("utf-8").encode(s);
        const l = bytes.length;
        const b = this._wasm_inst.exports["runtime.Block.HeapAlloc"](l, 0, 1);
        const d = b + 16;
        this._mem_util.mem_array_u8(d, l).set(bytes);
        return [b, d, l];
      },
      get_bytes: (d, l) => { return this._mem_util.mem_array_u8(d, l).slice(0); },
      set_bytes: (bytes) => {
        const l = bytes.length;
        const c = l;
        const b = this._wasm_inst.exports["runtime.Block.HeapAlloc"](l, 0, 1);
        const d = b + 16;
        this._mem_util.mem_array_u8(d, l).set(bytes);
        return [b, d, l, c];
      },
      block_release: (addr) => { this._wasm_inst.exports["runtime.Block.Release"](addr); },
      //基本类型直接读写：
      bool_load: (addr) => { return this._mem_util.mem_array_u8(addr, 1)[0] != 0; },
      bool_store: (addr, v) => {
        if (v) {
          this._mem_util.mem_array_u8(addr, 1)[0] = 1;
        } else {
          this._mem_util.mem_array_u8(addr, 1)[0] = 0;
        }
      },
      u8_load: (addr) => { return this._mem_util.mem_array_u8(addr, 1)[0]; },
      u8_store: (addr, v) => { this._mem_util.mem_array_u8(addr, 1)[0] = v; },
      u16_load: (addr) => { return this._mem_util.mem_array_u16(addr, 1)[0]; },
      u16_store: (addr, v) => { this._mem_util.mem_array_u16(addr, 1)[0] = v; },
      u32_load: (addr) => { return this._mem_util.mem_array_u32(addr, 1)[0]; },
      u32_store: (addr, v) => { this._mem_util.mem_array_u32(addr, 1)[0] = v; },
      i32_load: (addr) => { return this._mem_util.mem_array_i32(addr, 1)[0]; },
      i32_store: (addr, v) => { this._mem_util.mem_array_i32(addr, 1)[0] = v; },
      rune_load: (addr) => { return String.fromCodePoint(this._mem_util.mem_array_u32(addr, 1)[0]); },
      rune_store: (addr, v) => { this._mem_util.mem_array_u32(addr, 1)[0] = v.codePointAt(0); },
      f32_load: (addr) => { return this._mem_util.mem_array_f32(addr, 1)[0]; },
      f32_store: (addr, v) => { this._mem_util.mem_array_f32(addr, 1)[0] = v; },
      f64_load: (addr) => { return this._mem_util.mem_array_f64(addr, 1)[0]; },
      f64_store: (addr, v) => { this._mem_util.mem_array_f64(addr, 1)[0] = v; },
      string_load: (addr) => {
        const d = this._mem_util.i32_load(addr + 4);
        const l = this._mem_util.i32_load(addr + 8);
        return this._mem_util.get_string(d, l);
      },
      string_store: (addr, v) => {
        const b = this._mem_util.i32_load(addr);
        this._mem_util.block_release(b);
        let ns = this._mem_util.set_string(v);
        this._mem_util.i32_store(addr, ns[0]);
        this._mem_util.i32_store(addr + 4, ns[1]);
        this._mem_util.i32_store(addr + 8, ns[2]);
      },
      extract_string: (arr) => {
        const s = this._mem_util.get_string(arr[1], arr[2]);
        this._mem_util.block_release(arr[0]);
        arr.splice(0, 3);
        return s;
      },
      extract_bytes: (arr) => {
        const b = this._mem_util.get_bytes(arr[1], arr[2]);
        this._mem_util.block_release(arr[0]);
        arr.splice(0, 4);
        return b
      },
      extract_bool: (arr) => { const v = arr[0]; arr.splice(0, 1); return v?true:false; },
      extract_rune: (arr) => { const v = arr[0]; arr.splice(0, 1); return String.fromCodePoint(v); },
      extract_number: (arr) => { const v = arr[0]; arr.splice(0, 1); return v; },
    }
  };

  _createSyscall = () => {
    return {
      print_bool: (b) => { this._wa_print_buf += Boolean(b).toString(); },
      print_u32: (i) => {
        if (i < 0) {
          i += 4294967296;
        }
        this._wa_print_buf += i;
      },
      print_i32: (i) => { this._wa_print_buf += i },
      print_u64: (i) => { this._wa_print_buf += i },
      print_u64: (i) => { this._wa_print_buf += i },
      print_i64: (i) => { this._wa_print_buf += i },
      print_f32: (f) => { this._wa_print_buf += f },
      print_f64: (f) => { this._wa_print_buf += f },
      print_ptr: (p) => { this._wa_print_buf += p },
      print_str: (addr, len) => { this._wa_print_buf += this._mem_util.get_string(addr, len);},
      proc_exit: (code) => { alert(code) },
      print_rune: (c) => {
        let ch = String.fromCodePoint(c);
        if (ch == "\n") {
          console.log(this._wa_print_buf);
          this._wa_print_buf = "";
        }
        else {
          this._wa_print_buf += ch;
        }
      }
    }
  };

  async init() {
    const app = this;
    const imports = {
      syscall_js: this._createSyscall(),
      {{$.ImportCode}}
      // ...
    };

    try {
      const source = await fetch(this._wasm_url);
      const result = await WebAssembly.instantiateStreaming(source, imports);
      this._wasm_inst = result.instance;

      // 全局变量：
      {{ range $.Globals }}
      Object.defineProperty(this, "{{.Name}}", {
        get: function() { return this._mem_util.{{.Type}}_load(this._wasm_inst.exports["{{$.Pkg}}.{{.Name}}.1"]); },
        set: function (v) { this._mem_util.{{.Type}}_store(this._wasm_inst.exports["{{$.Pkg}}.{{.Name}}.1"], v); },
      });
      {{ end }}

      // 全局函数：
      {{ range $.Funcs }}
      this.{{.Name}} = function({{.Params}}) {
        // 准备参数
        let params = [];
        {{.PreCall}}
        {{if .ExportName}}
        let res = this._wasm_inst.exports["{{.ExportName}}"](...params);
        {{else}}
        let res = this._wasm_inst.exports["{{$.Pkg}}.{{.Name}}"](...params);
        {{end}}
        {{.GetResults}}
        {{.Release}}
        {{.Return}}
      }
      {{ end }}

      this._wasm_inst.exports._start();
      return this
    } catch (error) {
      console.error('WASM 初始化失败:', error);
    }
  }

}  // class WaApp

