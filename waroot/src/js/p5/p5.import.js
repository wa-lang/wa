p5js: new function() {
    const p5js = this;
    const theApp = app;

    theApp._p5js = p5js;

    this._inited = false;

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

            // 保持定义一致
            const keyCode_Enter = 13;
            const keyCode_Tab = 9;
            const keyCode_Space = 32;
            const keyCode_ArrowUp = 1004;
            const keyCode_ArrowDown = 1005;
            const keyCode_ArrowLeft = 1006;
            const keyCode_ArrowRight = 1007;
            const keyCode_Escape = 1008;
            const keyCode_Backspace = 1009;
            const keyCode_Delete = 1010;
            const keyCode_Shift = 1011;
            const keyCode_Control = 1012;
            const keyCode_Option = 1013;
            const keyCode_Alt = 1014;

            // 保持和凹语言环境定义一致
            switch (event.key) {
            case "Enter":
                theApp.p5js_onKeyDown(keyCode_Enter);
                return;
            case "Tab":
                theApp.p5js_onKeyDown(keyCode_Tab);
                return;
            case " ":
                theApp.p5js_onKeyDown(keyCode_Space);
                return;

            case "ArrowUp":
                theApp.p5js_onKeyDown(keyCode_ArrowUp);
                return;
            case "ArrowDown":
                theApp.p5js_onKeyDown(keyCode_ArrowDown);
                return;
            case "ArrowLeft":
                theApp.p5js_onKeyDown(keyCode_ArrowLeft);
                return;
            case "ArrowRight":
                theApp.p5js_onKeyDown(keyCode_ArrowRight);
                return;

            case "Escape":
                theApp.p5js_onKeyDown(keyCode_Escape);
                return;
            case "Backspace":
                theApp.p5js_onKeyDown(keyCode_Backspace);
                return;
            case "Delete":
                theApp.p5js_onKeyDown(keyCode_Delete);
                return;

            case "Shift":
                theApp.p5js_onKeyDown(keyCode_Shift);
                return;
            case "Control":
                theApp.p5js_onKeyDown(keyCode_Control);
                return;
            case "Meta":
                theApp.p5js_onKeyDown(keyCode_Option);
                return;
            case "Alt":
                theApp.p5js_onKeyDown(keyCode_Alt);
                return;

            default:
                theApp.p5js_onKeyDown(event.keyCode);
                return;
            }
        });
        document.addEventListener('keyup', (event) => {
            theApp.p5js_onKeyUp();
        });

        // 鼠标按键消息
        const canvas = document.getElementById("canvas");
        if(canvas) {

            // 焦点事件
            canvas.addEventListener("focus", (event) => {    
                theApp.p5js_onFocus();
            }, true);
            canvas.addEventListener("blur", (event) => {    
                theApp.p5js_onBlur();
            }, true);

            canvas.addEventListener("mouseenter", (event) => {    
                theApp.p5js_onMouseEnter();
            });
            canvas.addEventListener("mouseleave", (event) => {    
                theApp.p5js_onMouseLeave();
            });

            canvas.addEventListener("mousedown", (event) => {    
                theApp.p5js_onMouseDown(
                    event.button, event.offsetX, event.offsetY
                );
            });
            canvas.addEventListener("mouseup", (event) => {    
                theApp.p5js_onMouseUp();
            });

            canvas.addEventListener("mousemove", (event) => {
                theApp.p5js_onMouseMoved(event.offsetX, event.offsetY);
            });
        }

        // 帧函数
        if(theApp.Draw) {
            let stepAnima = function (timeStamp) {
                theApp.p5js_onDraw_before(timeStamp/1000.0);
                theApp.Draw();
                theApp.p5js_onDraw_after();

                window.requestAnimationFrame(stepAnima);
            }
            window.requestAnimationFrame(stepAnima);
        }
    }
},
