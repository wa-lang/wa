// 从 JSON 恢复 File 对象，并加方法

// <script src="fset.js"></script>
// <script>
//   // fsetJson 就是全局变量
//   let fileSet = restoreFileSet(fsetJson);
//   let pos = fileSet.Position(5);
//   console.log(pos);  
//   // => { filename: "main.wa", offset: 4, line: 1, column: 5 }
// </script>

// fset.js 是编译器生成的位置信息文件, 格式如下
// const fsetJson = {
//   base: 1,
//   files: [
//     { name: "main.wa", base: 1, size: 100, lines: [0,10,20], infos: [] }
//   ]
// };

function restoreFileSet(f) {
    return {
        name: f.name,
        base: f.base,
        size: f.size,
        lines: f.lines,

        // PositionFor: 根据位置 p 和 adjusted=true 返回 Position 对象
        PositionFor: function(p, adjusted) {
            // f.base <= int(p) <= f.base+f.size
            let offset = p - this.base; // 相对文件的 offset

            // 计算行号：在 lines 中找到最后一个 <= offset 的行开始位置
            let lineIndex = searchInts(this.lines, offset); // 返回行数组下标
            let lineStart = this.lines[lineIndex] || 0;
            let line = lineIndex + 1; // 行号从 1 开始
            let column = offset - lineStart + 1; // 列号从 1 开始

            return {
                filename: this.name,
                offset: offset,
                line: line,
                column: column
            };
        },

        // Position: 相当于调用 PositionFor(p, true)
        Position: function(p) {
            return this.PositionFor(p, true);
        }
    };
}

// searchInts: 在升序数组 lines 中找最后一个 <= offset 的下标
function searchInts(lines, offset) {
    let i = 0, j = lines.length;
    while (i < j) {
        let h = i + ((j - i) >> 1);
        if (lines[h] <= offset) {
            i = h + 1;
        } else {
            j = h;
        }
    }
    return i - 1;
}
