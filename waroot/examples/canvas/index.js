// 版权 @2023 hello 作者。保留所有权利。

const myapp = document.getElementById("myapp");
const ctx = myapp.getContext("2d");

const importsObject = {
    syscall_js: new function() {
        this.print_bool = (v) => console.log(v? 'true': 'faslse');
        this.print_i32 = (i) => console.log(i);
        this.print_u32 = (i) =>  console.log(i);
        this.print_ptr = (i) => console.log(i);
        this.print_i64 = (i) => console.log(i);
        this.print_u64 = (i) => console.log(i);
        this.print_f32 = (i) => console.log(i);
        this.print_f64 = (i) => console.log(i);

        this.print_rune = (c) => {
            let ch = String.fromCodePoint(c);
            console.log(ch);
        }
        this.print_str = (prt, len) => {
            //let s = window.waApp.getString(prt, len);
            //console.log(s);
        }
        this.proc_exit = (i) => {
            // exit(i);
        }
    },

    canvas: {
        rect: (x, y, width, height) => ctx.rect(x, y, width, height),
        moveTo: (x, y) => ctx.moveTo(x, y),
        lineTo: (x, y) => ctx.lineTo(x, y),
        stroke: () => ctx.stroke(),
    }
}

WebAssembly.instantiateStreaming(
    fetch("./output/hello.wasm"), importsObject).then((obj) => {
        obj.instance.exports._start();
        obj.instance.exports['myapp.DrawLogo']();
    }
);
