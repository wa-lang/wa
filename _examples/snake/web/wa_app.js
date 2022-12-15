(() => {
    class WaApp {
      constructor() {
        this._inst = null;
        this._wa_print_buf = "";
      }
  
      init(url) {
        let app = this;
        let importsObject = {
          wa_js_env: new function () {
            this.waPrintI32 = (i) => {
              app._wa_print_buf += i
            }
            this.waPrintRune = (c) => {
              let ch = String.fromCodePoint(c);
              if (ch == '\n') {
                console.log(app._wa_print_buf);
                app._wa_print_buf = "";
              }
              else {
                app._wa_print_buf += ch
              }
            }
            this.waPuts = (prt, len) => {
              let s = app.getString(prt, len);
              app._wa_print_buf += s
            }
            this.newCanvas = (w, h) => {
              let canvas = document.createElement('canvas');
              canvas.width = w;
              canvas.height = h;
              const waContent = document.getElementById('wa-content');
              waContent.appendChild(canvas);
              this._ctx = canvas.getContext('2d');
              this._canvas = canvas;
              canvas.id = 0;
              return 0  //!!!!!
            }
            this.updateCanvas = (id, block, data) => {
              let img = this._ctx.createImageData(this._canvas.width, this._canvas.height);
              let buf_len = this._canvas.width * this._canvas.height * 4
              let buf = app.memUint8Array(data, buf_len);
              for (var i = 0; i < buf_len; i++){
                img.data[i] = buf[i];
              }
              this._ctx.putImageData(img, 0, 0);
            }
          }
        }
        WebAssembly.instantiateStreaming(fetch(url), importsObject).then(res => {
          this._inst = res.instance;
          this._inst.exports._start();
        })
      }
  
      mem() {
        return this._inst.exports.memory;
      }
  
      memView(addr, len) {
        return new DataView(this._inst.exports.memory.buffer, addr, len);
      }
  
      memUint8Array(addr, len) {
        return new Uint8Array(this.mem().buffer, addr, len)
      }
  
      getString(addr, len) {
        return new TextDecoder("utf-8").decode(this.memView(addr, len));
      }
  
      setString(addr, len, s) {
        const bytes = new TextEncoder("utf-8").encode(s);
        if (len > bytes.length) { len = bytes.length; }
        this.MemUint8Array(addr, len).set(bytes);
      }
    }
  
    window['waApp'] = new WaApp();
    window['waApp'].init("./snake.wasm")
  })()