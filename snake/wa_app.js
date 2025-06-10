(() => {
  class WaApp {
    constructor() {
      this._inst = null;
    }

    init(url) {
      let app = this;
      let importsObject = {
        syscall_js: new function () {
          this.print_bool = (b) => { console.log(b) },
          this.print_u32 = (i) => { console.log(i) },
          this.print_i32 = (i) => { console.log(i) },
          this.print_u64 = (i) => { console.log(i) },
          this.print_u64 = (i) => { console.log(i) },
          this.print_i64 = (i) => { console.log(i) },
          this.print_f32 = (f) => { console.log(f) },
          this.print_f64 = (f) => { console.log(f) },
          this.print_ptr = (p) => {  console.log(p) },
          this.print_str = (d, l) => {
            const mem = app._inst.exports.memory;
            const mem_view = new DataView(mem.buffer, d, l);
            const s = new TextDecoder("utf-8").decode(mem_view);
            console.log(s);
          }
          this.proc_exit = (code) => { alert(code) },
          this.print_rune = (c) => {
            let ch = String.fromCodePoint(c);
            console.log(ch)
          }
        },
        snake_game: new function () {
          this.rand = (m) => {
            return parseInt(Math.random() * m)
          }

          this.newCanvas = (w, h) => {
            let canvas = document.createElement('canvas');
            canvas.width = w;
            canvas.height = h;
            canvas.id = 0;  //!!!!!

            const waContent = document.getElementById('game__screen-content');
            waContent.appendChild(canvas);

            function getPointOnCanvas(x, y) {
              var bbox = canvas.getBoundingClientRect();
              return {
                x: parseInt((x - bbox.left) * (canvas.width / bbox.width)),
                y: parseInt((y - bbox.top) * (canvas.height / bbox.height))
              };
            }

            function onMouseDown(ev) {
              let pt = getPointOnCanvas(ev.clientX, ev.clientY);
              app._inst.exports['snake.Canvas_OnMouseDown'](canvas.id, pt.x, pt.y);
            }

            function onMouseUp(ev) {
              let pt = getPointOnCanvas(ev.clientX, ev.clientY);
              app._inst.exports['snake.Canvas_OnMouseUp'](canvas.id, pt.x, pt.y);
            }

            function onKeyDown(ev) {
              app._inst.exports['snake.Canvas_OnKeyDown'](canvas.id, ev.keyCode);
            }

            function onKeyUp(ev) {
              app._inst.exports['snake.Canvas_OnKeyUp'](canvas.id, ev.keyCode);
            }

            if (IS_MOBILE) {
              MOBILE_DIR_MAP.forEach((dir) => {
                const el = document.getElementById(dir.id);
                el.addEventListener('touchstart', (ev) => onKeyDown({ keyCode: dir.keyCode }));
                el.addEventListener('touchend', (ev) => onKeyUp({ keyCode: dir.keyCode }));
              });
            }

            canvas.addEventListener('mousedown', onMouseDown, true);
            canvas.addEventListener('mouseup', onMouseUp, true);
            canvas.addEventListener('keydown', onKeyDown, true);
            canvas.addEventListener('keyup', onKeyUp, true);
            canvas.tabIndex = -1;  //tabindex
            canvas.focus();

            this._ctx = canvas.getContext('2d');
            this._canvas = canvas;
            return canvas.id;
          }
          this.updateCanvas = (id, block, data) => {
            let img = this._ctx.createImageData(this._canvas.width, this._canvas.height);
            let buf_len = this._canvas.width * this._canvas.height * 4
            let buf = app.memUint8Array(data, buf_len);
            for (var i = 0; i < buf_len; i++) {
              img.data[i] = buf[i];
            }
            this._ctx.putImageData(img, 0, 0);
          }
        }
      }
      WebAssembly.instantiateStreaming(fetch(url), importsObject).then(res => {
        this._inst = res.instance;
        this._inst.exports._start();
        this._inst.exports['snake.main']();
        const timer = setInterval(gameLoop, 150);
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

  function gameLoop() {
    window['waApp']._inst.exports['snake.Step']();
  }

  window['waApp'] = new WaApp();
  window['waApp'].init("./snake.wasm")
})()