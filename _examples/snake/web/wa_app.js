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
              app._inst.exports['snake$canvas.OnMouseDown'](canvas.id, pt.x, pt.y);
            }

            function onMouseUp(ev) {
              let pt = getPointOnCanvas(ev.clientX, ev.clientY);
              app._inst.exports['snake$canvas.OnMouseUp'](canvas.id, pt.x, pt.y);
            }

            function onKeyDown(ev) {
              app._inst.exports['snake$canvas.OnKeyDown'](canvas.id, ev.keyCode);
            }

            function onKeyUp(ev) {
              app._inst.exports['snake$canvas.OnKeyUp'](canvas.id, ev.keyCode);
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