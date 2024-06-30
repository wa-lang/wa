extobj: new function() {
  this._obj_buf = [{}];
  this._free_ids = [];
  this.new_obj = () => {
    let obj = {}
    let h = this._free_ids.length > 0
      ? this._free_ids.pop()
      : this._obj_buf.length
    this._obj_buf[h] = obj;
    return h
  }
  this.free_obj = (h) => {
    let obj = this._obj_buf[h];
    // obj.free(); // TODO: 释放资源
    this._obj_buf[h] = null;
    this._free_ids.push(h);
  }
  this.insert_obj = (obj) => {
    let h = this.new_obj()
    this._obj_buf[h] = obj
    return h
  }
  this.get_obj = (h) => {
    return this._obj_buf[h];
  }
  this.set_obj = (h, obj) => {
    this._obj_buf[h] = obj;
  }
  this.query_selector = (sel_b, sel_d, sel_l) => {
    const selector = app._mem_util.get_string(sel_d, sel_l);
    const obj = document.querySelector(selector);
    if (obj) {
      return this.insert_obj(obj);
    } else {
      return 0;
    }
  }
  this.set_member_bool = (h, member_name_b, member_name_d, member_name_l, value) => {
    let member_name = app._mem_util.get_string(member_name_d, member_name_l);
    if (value === 0){
      this._obj_buf[h][member_name] = false;
    } else {
      this._obj_buf[h][member_name] = true;
    }
  }
  this.set_member_i32 = (h, member_name_b, member_name_d, member_name_l, value) => {
    let member_name = app._mem_util.get_string(member_name_d, member_name_l);
    this._obj_buf[h][member_name] = value;
  }
  this.set_member_f32 = (h, member_name_b, member_name_d, member_name_l, value) => {
    let member_name = app._mem_util.get_string(member_name_d, member_name_l);
    this._obj_buf[h][member_name] = value;
  }
  this.set_member_string = (h, name_b, name_d, name_l, value_b, value_d, value_l) => {
    let name = app._mem_util.get_string(name_d, name_l);
    let value = app._mem_util.get_string(value_d, value_l);
    this._obj_buf[h][name] = value;
  }
  this.set_member_obj = (h, member_name_b, member_name_d, member_name_l, value) => {
    if (value > 0) {
      let member_name = app._mem_util.get_string(member_name_d, member_name_l);
      this._obj_buf[h][member_name] = this._obj_buf[value];
    }    
  }
  this.new_array = () => {
    let arr = [];
    let h = this._free_ids.length > 0
      ? this._free_ids.pop()
      : this._obj_buf.length
    this._obj_buf[h] = arr
    return h
  }
  this.append_i32 = (h, value) => {
    this._obj_buf[h].push(value);
  }
  this.append_string = (h, value_b, value_d, value_l) => {
    let value = app._mem_util.get_string(value_d, value_l);
    this._obj_buf[h].push(value);
  }
  this.append_obj = (h, value) => {
    if (value > 0) {
      this._obj_buf[h].push(this._obj_buf[value]);
    }
  }
  app._extobj = this;
},