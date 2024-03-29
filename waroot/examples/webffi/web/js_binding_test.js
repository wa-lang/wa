class WaApp {
  constructor(url) {
    let app = this;
    this._wasm_inst = null;
    this._wa_print_buf = "";

    this._mem_util = new function() {
      this.mem = () => { return app._wasm_inst.exports.memory; }
      this.mem_view = (addr, len) => { return new DataView(this.mem().buffer, addr, len); }
      this.mem_array_u8 = (addr, len) => { return new Uint8Array(this.mem().buffer, addr, len); }

      this.get_string = (d, l) => { return new TextDecoder("utf-8").decode(this.mem_view(d, l)); }
      this.set_string = (s) => {
        const bytes = new TextEncoder("utf-8").encode(s);
        const l = bytes.length;
        const b = app._wasm_inst.exports["runtime.Block.HeapAlloc"](l, 0, 1);
        const d = b + 16;
        this.mem_array_u8(d, l).set(bytes);
        return [b, d, l]
      }

      this.block_release = (addr) => { app._wasm_inst.exports["runtime.Block.Release"](addr); }

      //基本类型直接读写：
      this.bool_load = (addr) => { /*Todo*/ }
      this.bool_store = (addr, v) => { /*Todo*/ }
      this.u8_load = (addr) => { /*Todo*/ }
      this.u8_store = (addr, v) => { /*Todo*/ }
      this.u16_load = (addr) => { /*Todo*/ }
      this.u16_store = (addr, v) => { /*Todo*/ }
      this.u32_load = (addr) => { /*Todo*/ }
      this.u32_store = (addr, v) => { /*Todo*/ }
      this.i32_load = (addr) => { return app._wasm_inst.exports["runtime.i32_load"](addr); }
      this.i32_store = (addr, v) => { app._wasm_inst.exports["runtime.i32_store"](addr, v); }
      this.rune_load = (addr) => { /*Todo*/ }
      this.rune_store = (addr, v) => { /*Todo*/ }
      this.u64_load = (addr) => { /*Todo*/ }
      this.u64_store = (addr, v) => { /*Todo*/ }
      this.i64_load = (addr) => { /*Todo*/ }
      this.i64_store = (addr, v) => { /*Todo*/ }
      this.f32_load = (addr) => { return app._wasm_inst.exports["runtime.f32_load"](addr); }
      this.f32_store = (addr, v) => { app._wasm_inst.exports["runtime.f32_store"](addr, v); }
      this.f64_load = (addr) => { return app._wasm_inst.exports["runtime.f64_load"](addr); }
      this.f64_store = (addr, v) => { app._wasm_inst.exports["runtime.f64_store"](addr, v); }
      this.string_load = (addr) => {
        const d = this.i32_load(addr + 4)
        const l = this.i32_load(addr + 8)
        return this.get_string(d, l);
      }
      this.string_store = (addr, v) => {
        const b = this.i32_load(addr)
        this.block_release(b)
        let ns = this.set_string(v)
        this.i32_store(addr, ns[0])
        this.i32_store(addr + 4, ns[1])
        this.i32_store(addr + 8, ns[2])
      }

      //返回值提取：
      this.extract_string = (arr) => {
        const s = this.get_string(arr[1], arr[2]);
        this.block_release(arr[0])
        arr.splice(0, 3)
        return s
      }
      this.extract_bool = (arr) => {
        //Todo
        arr.splice(0, 1)
        return 0
      }
      this.extract_u8 = (arr) => {
        //Todo
        arr.splice(0, 1)
        return 0
      }
      this.extract_u16 = (arr) => {
        //Todo
        arr.splice(0, 1)
        return 0
      }
      this.extract_u32 = (arr) => {
        //Todo
        arr.splice(0, 1)
        return 0
      }
      this.extract_i32 = (arr) => {
        const v = arr[0]
        arr.splice(0, 1)
        return v
      }
      this.extract_rune = (arr) => {
        //Todo
        arr.splice(0, 1)
        return 0
      }
      this.extract_f32 = (arr) => {
        const v = arr[0]
        arr.splice(0, 1)
        return 0
      }
      this.extract_f64 = (arr) => {
        const v = arr[0]
        arr.splice(0, 1)
        return 0
      }
    }

    let syscall = new function() {
      this.print_bool = (b) => { app._wa_print_buf += Boolean(b) }
      this.print_u32 = (i) => { app._wa_print_buf += i }
      this.print_i32 = (i) => { app._wa_print_buf += i }
      this.print_u64 = (i) => { app._wa_print_buf += i }
      this.print_i64 = (i) => { app._wa_print_buf += i }
      this.print_f32 = (f) => { app._wa_print_buf += f }
      this.print_f64 = (f) => { app._wa_print_buf += f }
      this.print_ptr = (p) => { app._wa_print_buf += p }
      this.print_str = (addr, len) => { app._wa_print_buf += app._mem_util.get_string(addr, len) }
      this.proc_exit = (code) => { alert(code) }
      this.print_rune = (c) => {
        let ch = String.fromCodePoint(c);
        if (ch == '\n') {
          console.log(app._wa_print_buf);
          app._wa_print_buf = "";
        }
        else {
          app._wa_print_buf += ch
        }
      }
    }

    let imports = {
      syscall_js: syscall
    }

    WebAssembly.instantiateStreaming(fetch(url), imports).then(res => {
      this._wasm_inst = res.instance;

      // 全局变量：
      
      Object.defineProperty(this, "I", {
        get: function() { return this._mem_util.i32_load(this._wasm_inst.exports["webffi.I.1"]); },
        set: function (v) { this._mem_util.i32_store(this._wasm_inst.exports["webffi.I.1"], v); },
      });
      
      Object.defineProperty(this, "S", {
        get: function() { return this._mem_util.string_load(this._wasm_inst.exports["webffi.S.1"]); },
        set: function (v) { this._mem_util.string_store(this._wasm_inst.exports["webffi.S.1"], v); },
      });
      

      // 全局函数：
      
      this.Fn1 = function() {
        // 准备参数
        let params = [];
        
        let res = this._wasm_inst.exports["webffi.Fn1"](...params);
        
        
        
      }
      
      this.Fn2 = function(s, i) {
        // 准备参数
        let params = [];
        let p0 = this._mem_util.set_string(s);
params = params.concat(p0);
params.push(i);

        let res = this._wasm_inst.exports["webffi.Fn2"](...params);
        let r0 = this._mem_util.extract_string(res);
let r1 = this._mem_util.extract_i32(res);
let r2 = this._mem_util.extract_string(res);

        this._mem_util.block_release(p0[0]);

        return [r0,r1,r2];
      }
      
      this.Fn3 = function() {
        // 准备参数
        let params = [];
        
        let res = this._wasm_inst.exports["webffi.Fn3"](...params);
        
        
        
      }
      

      this._wasm_inst.exports._start();
    })

  }  // constructor
}  // class WaApp

window['waApp'] = new WaApp("./webffi.wasm");