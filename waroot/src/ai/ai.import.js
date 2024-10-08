
aiproxy: new function () {

  //---------------------------------------------------------------

  this.request_session = (tid, ob, od, ol) => {
    if (!ai) {
      throw Error('ai not supported.');
    }

    let option = app._mem_util.get_string(od, ol);

    ai.assistant.create().then((session) => {
      if (!session) {
        alert('Couldn\'t request session.');
        throw Error('Couldn\'t request session.');
      }
      let sh = app._extobj.insert_obj(session);
      app._wasm_inst.exports["ai.onSessionRequested"](tid, sh);
    })
    .catch((err) => {
      console.log(`ai.assistant.create(): err = ${err}`);
      app._wasm_inst.exports["ai.onSessionRequested"](tid, 0);
    })
  }

  //---------------------------------------------------------------

  this.prompt = (tid, sh, keyb, keyd, keyl) => {
    let session = app._extobj.get_obj(sh);
    let key = app._mem_util.get_string(keyd, keyl);
    session.prompt(key)
      .then((res) =>{
        let params = [];
        params.push(tid);
        let s = app._mem_util.set_string(res);
        params = params.concat(s);
        
        app._wasm_inst.exports["ai.onPrompted"](...params);        
        
        app._mem_util.block_release(s[0]);
      })
      .catch((err) => {
        console.log(`ai.prompt: err = ${err}`)
      });
  }

  //---------------------------------------------------------------

},
