net: new function() {
  this.fetch_blob = (tid, res_b, res_d, res_l) => {
    let resource = app._mem_util.get_string(res_d, res_l)
    let status = 0
    fetch(resource)
      .then((response) => {
        status = response.status;
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.arrayBuffer();
      })
      .then((data) => {
        let u8a = new Uint8Array(data);
        let slice = app._mem_util.set_bytes(u8a);
        app._wasm_inst.exports["net.onFetchBlobDone"](tid, 1, status, ...slice);        
        app._mem_util.block_release(slice[0]);
      })
      .catch((err) => {
        app._wasm_inst.exports["net.onFetchBlobDone"](tid, 0, status, 0, 0, 0, 0);
        console.log(err)
      })
  }

  this.fetch_image = (tid, res_b, res_d, res_l) => {
    let resource = app._mem_util.get_string(res_d, res_l)
    let status = 0
    fetch(resource)
      .then((response) => {
        status = response.status;
        if (!response.ok) {
          throw new Error(`HTTP error! Status: ${response.status}`);
        }
        return response.blob();
      })
      .then((blob) => {
        return createImageBitmap(blob);
      })
      .then((img) => {
        let imgid = app._extobj.insert_obj(img);
        app._wasm_inst.exports["net.onFetchImageDone"](tid, 1, status, imgid);
      })
      .catch((err) => {
        app._wasm_inst.exports["net.onFetchImageDone"](tid, 0, status, 0);
        console.log(err)
      })
  }
},