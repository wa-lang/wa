p5js: new function() {
    const p5js = this;
    const theApp = app;

    theApp._p5js = p5js;

    this._inited = false;
    this._preTimeStamp = 1;

    // 初始化
    this.init = () => {
        if(this._inited) { return; }
        this._inited = true;

        // 有 main 函数则不生成框架代码
        if(theApp.main) { return; }

        // 生成空的 main 函数
        theApp.main = () => {}

        // 键盘消息
        document.addEventListener('keydown', (event) => {
            // console.log(`p5: document.addEventListener`);
        });
        document.addEventListener('keyup', (event) => {
            // console.log(`p5: document.addEventListener`);
        });

        // 鼠标按键消息
        const canvas = document.getElementById("canvas");
        if(canvas) {
            canvas.addEventListener("mousedown", (event) => {    
                theApp.p5js_onMousePressed(true);
            });
            canvas.addEventListener("mouseup", (event) => {    
                theApp.p5js_onMousePressed(false);
            });

            canvas.addEventListener("mousemove", (event) => {
                theApp.p5js_onMouseMoved(event.offsetX, event.offsetY);
            });
        }

        // 帧函数
        if(theApp.Draw) {
            let stepAnima = function (timeStamp) {
                theApp.Draw();
                p5js._preTimeStamp = timeStamp;
                window.requestAnimationFrame(stepAnima);
            }
            window.requestAnimationFrame(stepAnima);
        }
    }
},
