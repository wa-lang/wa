const importsObject = {
  wa_js_env: new function () {
    this.newCanvas_JS = (w, h, id) => {
      const canvas = document.createElement('canvas');
      canvas.id = id;
      canvas.width = w;
      canvas.height = h;
      const waContent = document.getElementById('wa-content');
      waContent.appendChild(canvas);
    }
  }
}
