
webgpu: new function () {
  this._texture_format_map = new Map([
    ["r8unorm", 0],
    ["r8snorm", 1],
    ["r8uint", 2],
    ["r8sint", 3],
    ["r16uint", 4],
    ["r16sint", 5],
    ["r16float", 6],
    ["rg8unorm", 7],
    ["rg8snorm", 8],
    ["rg8uint", 9],
    ["rg8sint", 10],
    ["r32uint", 11],
    ["r32sint", 12],
    ["r32float", 13],
    ["rg16uint", 14],
    ["rg16sint", 15],
    ["rg16float", 16],
    ["rgba8unorm", 17],
    ["rgba8unorm-srgb", 18],
    ["rgba8snorm", 19],
    ["rgba8uint", 20],
    ["rgba8sint", 21],
    ["bgra8unorm", 22],
    ["bgra8unorm-srgb", 23],
    ["rgb9e5ufloat", 24],
    ["rgb10a2uint", 25],
    ["rgb10a2unorm", 26],
    ["rg11b10ufloat", 27],
    ["rg32uint", 28],
    ["rg32sint", 29],
    ["rg32float", 30],
    ["rgba16uint", 31],
    ["rgba16sint", 32],
    ["rgba16float", 33],
    ["rgba32uint", 34],
    ["rgba32sint", 35],
    ["rgba32float", 36],
    ["stencil8", 37],
    ["depth16unorm", 38],
    ["depth24plus", 39],
    ["depth24plus-stencil8", 40],
    ["depth32float", 41],
    ["depth32float-stencil8", 42],
    ["bc1-rgba-unorm", 43],
    ["bc1-rgba-unorm-srgb", 44],
    ["bc2-rgba-unorm", 45],
    ["bc2-rgba-unorm-srgb", 46],
    ["bc3-rgba-unorm", 47],
    ["bc3-rgba-unorm-srgb", 48],
    ["bc4-r-unorm", 49],
    ["bc4-r-snorm", 50],
    ["bc5-rg-unorm", 51],
    ["bc5-rg-snorm", 52],
    ["bc6h-rgb-ufloat", 53],
    ["bc6h-rgb-float", 54],
    ["bc7-rgba-unorm", 55],
    ["bc7-rgba-unorm-srgb", 56],
    ["etc2-rgb8unorm", 57],
    ["etc2-rgb8unorm-srgb", 58],
    ["etc2-rgb8a1unorm", 59],
    ["etc2-rgb8a1unorm-srgb", 60],
    ["etc2-rgba8unorm", 61],
    ["etc2-rgba8unorm-srgb", 62],
    ["eac-r11unorm", 63],
    ["eac-r11snorm", 64],
    ["eac-rg11unorm", 65],
    ["eac-rg11snorm", 66],
    ["astc-4x4-unorm", 67],
    ["astc-4x4-unorm-srgb", 68],
    ["astc-5x4-unorm", 69],
    ["astc-5x4-unorm-srgb", 70],
    ["astc-5x5-unorm", 71],
    ["astc-5x5-unorm-srgb", 72],
    ["astc-6x5-unorm", 73],
    ["astc-6x5-unorm-srgb", 74],
    ["astc-6x6-unorm", 75],
    ["astc-6x6-unorm-srgb", 76],
    ["astc-8x5-unorm", 77],
    ["astc-8x5-unorm-srgb", 78],
    ["astc-8x6-unorm", 79],
    ["astc-8x6-unorm-srgb", 80],
    ["astc-8x8-unorm", 81],
    ["astc-8x8-unorm-srgb", 82],
    ["astc-10x5-unorm", 83],
    ["astc-10x5-unorm-srgb", 84],
    ["astc-10x6-unorm", 85],
    ["astc-10x6-unorm-srgb", 86],
    ["astc-10x8-unorm", 87],
    ["astc-10x8-unorm-srgb", 88],
    ["astc-10x10-unorm", 89],
    ["astc-10x10-unorm-srgb", 90],
    ["astc-12x10-unorm", 91],
    ["astc-12x10-unorm-srgb", 92],
    ["astc-12x12-unorm", 93],
    ["astc-12x12-unorm-srgb", 94],
  ])

  //---------------------------------------------------------------

  this.gpu_get_preferred_canvas_format = () => {
    return this._texture_format_map.get(navigator.gpu.getPreferredCanvasFormat());
  }

  this.gpu_request_adapter = (tid, option_h) => {
    if (!navigator.gpu) {
      throw Error('WebGPU not supported.');
    }

    let option = {};
    if (option_h) {
      option = app._extobj.get_obj(option_h);
    }

    navigator.gpu.requestAdapter(option).then((adapter) => {
      if (!adapter) {
        alert('Couldn\'t request WebGPU adapter.');
        throw Error('Couldn\'t request WebGPU adapter.');
      }
      let ah = app._extobj.insert_obj(adapter);
      app._wasm_inst.exports["gpu.onAdapterRequested"](tid, ah);
    })
    .catch((err) => {
      app._wasm_inst.exports["gpu.onAdapterRequested"](tid, 0);
      console.log(err)
    })
  }

  //---------------------------------------------------------------

  this.obj_get_label = (h) => {
    let obj = app._extobj.get_obj(h);
    let str = app._mem_util.set_string(obj.label);
    return str;
  }

  //---------------------------------------------------------------

  this.adapter_get_features = (ah) => {
    return app._extobj.insert_obj(app._extobj.get_obj(ah).features)
  }

  this.adapter_get_is_fallback_adapter = (ah) => {
    return app._extobj.get_obj(ah).isFallbackAdapter ? 1 : 0;
  }

  this.adapter_get_limits = (ah) => {
    return app._extobj.insert_obj(app._extobj.get_obj(ah).limits)
  }

  this.adapter_request_device = (tid, ah, desc_h) => {
    const adapter = app._extobj.get_obj(ah);
    let desc = {};
    if (desc_h) {
      desc = app._extobj.get_obj(desc_h);
    }

    adapter.requestDevice(desc).then((device) => {
      const device_h = app._extobj.insert_obj(device);
      app._wasm_inst.exports["gpu.onDeviceRequested"](tid, device_h);
    })
  }

  //---------------------------------------------------------------

  this.buffer_map_async = (tid, h, mode) => {
    let buffer = app._extobj.get_obj(h);
    buffer.mapAsync(mode)
      .then(() => {
        buffer._waMappedRange = new Uint8Array(buffer.getMappedRange());
        let slice = app._mem_util.set_bytes(buffer._waMappedRange);
        app._wasm_inst.exports["gpu.onBufferMapped"](tid, 1, ...slice);
        app._mem_util.block_release(slice[0]);
      })
      .catch((err) => {
        app._wasm_inst.exports["gpu.onBufferMapped"](tid, 0, 0, 0, 0, 0);
        console.log(err)
      })
  }

  this.buffer_map_state = (b_h) => {
    let buffer = app._extobj.get_obj(b_h);
    let state = buffer.mapState;
    switch (state) {
      case 'unmapped':
        return 0;

      case 'pending':
        return 1;

      case 'mapped':
        return 2;

      default:
        throw new Error(`Unknown mapState: ` + state);
    }
  }

  this.buffer_size = (b_h) => {
    return app._extobj.get_obj(b_h).size;
  }

  this.buffer_usage = (b_h) => {
    return app._extobj.get_obj(b_h).usage;
  }

  this.buffer_get_mapped_range = (h) => {
    let buffer = app._extobj.get_obj(h);
    buffer._waMappedRange = new Uint8Array(buffer.getMappedRange());
    let slice = app._mem_util.set_bytes(buffer._waMappedRange);
    return slice;
  }
  
  this.buffer_unmap = (h, b, d, l, c) => {
    let buffer = app._extobj.get_obj(h);
    let u8a = app._mem_util.get_bytes(d, l);
    buffer._waMappedRange.set(u8a);
    buffer._waMappedRange = null;
    buffer.unmap();
  }

  //---------------------------------------------------------------

  this.canvas_get_contex = (h) => {
    const canvas = app._extobj.get_obj(h);
    const ctx = canvas.getContext('webgpu');
    return app._extobj.insert_obj(ctx);
  }

  this.contex_get_canvas = (h) => {
    const contex = app._extobj.get_obj(h);
    return app._extobj.insert_obj(contex.canvas);
  }

  this.contex_configure = (contex_h, config_h) => {
    const contex = app._extobj.get_obj(contex_h);
    const config = app._extobj.get_obj(config_h);
    contex.configure(config);
  }

  this.contex_get_current_texture = (contex) => {
    let texture = app._extobj.get_obj(contex).getCurrentTexture();
    return app._extobj.insert_obj(texture)
  }

  this.contex_unconfigure = (h) => {
    app._extobj.get_obj(h).contex_unconfigure();
  }

  //---------------------------------------------------------------
  
  this.commandencoder_begin_render_pass = (command_encoder, render_pass_desc) => {
    let render_pass = app._extobj.get_obj(command_encoder).beginRenderPass(app._extobj.get_obj(render_pass_desc));
    return app._extobj.insert_obj(render_pass);
  }

  this.commandencoder_finish = (command_encoder_h) => {
    const command_encoder = app._extobj.get_obj(command_encoder_h);
    const command_buffer = command_encoder.finish();
    return app._extobj.insert_obj(command_buffer);
  }


  //---------------------------------------------------------------

  this.device_get_queue = (device) => {
    let queue = app._extobj.get_obj(device).queue;
    return app._extobj.insert_obj(queue);
  }

  this.device_create_bind_group = (device, bg_desc) => {
    let bind_group = app._extobj.get_obj(device).createBindGroup(app._extobj.get_obj(bg_desc));
    return app._extobj.insert_obj(bind_group);
  }

  this.device_create_buffer = (device, label_b, label_d, label_l, mapped, byteLen, usage) => {
    let buffer = app._extobj.get_obj(device).createBuffer({
      label: app._mem_util.get_string(label_d, label_l),
      mappedAtCreation: mapped != 0,
      size: byteLen,
      usage: usage,
    });
    return app._extobj.insert_obj(buffer);
  }

  this.device_create_command_encoder = (device) => {
    let encoder = app._extobj.get_obj(device).createCommandEncoder();
    return app._extobj.insert_obj(encoder);
  }

  this.device_create_render_bundle_encoder = (device, desc) => {
    let rb_encoder = app._extobj.get_obj(device).createRenderBundleEncoder(app._extobj.get_obj(desc));
    return app._extobj.insert_obj(rb_encoder);
  }

  this.device_create_render_pipeline = (device, pl_desc) => {
    let pipeline = app._extobj.get_obj(device).createRenderPipeline(app._extobj.get_obj(pl_desc));
    return app._extobj.insert_obj(pipeline);
  }

  this.device_create_sampler = (device, dh) => {
    const desc = app._extobj.get_obj(dh)
    let sampler = app._extobj.get_obj(device).createSampler(desc);
    return app._extobj.insert_obj(sampler);
  }

  this.device_create_shader_module = (device, desc_h) => {
    let shader = app._extobj.get_obj(device).createShaderModule(app._extobj.get_obj(desc_h));
    return app._extobj.insert_obj(shader);
  }

  this.device_create_texture = (device, desc) => {
    let texture = app._extobj.get_obj(device).createTexture(app._extobj.get_obj(desc));
    return app._extobj.insert_obj(texture);
  }

  //---------------------------------------------------------------

  this.queue_copy_external_image_to_texture = (qh, src, dest) => {
    const imageBitmap = app._extobj.get_obj(src);
    const tex = app._extobj.get_obj(dest);
    app._extobj.get_obj(qh).copyExternalImageToTexture(
      { source: imageBitmap },
      { texture: tex },
      [imageBitmap.width, imageBitmap.height]
    );
  }

  this.queue_submit = (qh, command_buffer_h) => {
    app._extobj.get_obj(qh).submit([app._extobj.get_obj(command_buffer_h)]);
  }

  this.queue_write_buffer = (qh, bh, offset, data_b, data_d, data_l, data_c) => {
    const queue = app._extobj.get_obj(qh);
    const buffer = app._extobj.get_obj(bh);
    const data = app._mem_util.mem_array_u8(data_d, data_l);
    queue.writeBuffer(buffer, offset, data)
  }

  //---------------------------------------------------------------

  this.render_encoder_draw = (rh, vertexCount, instanceCount, firstVertex, firstInstance) => {
    app._extobj.get_obj(rh).draw(vertexCount, instanceCount, firstVertex, firstInstance);
  }

  this.render_encoder_draw_indexed = (rh, indexCount, instanceCount, firstIndex, baseVertex, firstInstance) => {
    app._extobj.get_obj(rh).drawIndexed(indexCount, instanceCount, firstIndex, baseVertex, firstInstance);
  }

  this.render_encoder_draw_indirect = (rh, bh, offset) => {
    let encoder = app._extobj.get_obj(rh);
    let buffer = app._extobj.get_obj(bh);
    encoder.drawIndirect(buffer, offset);
  }
  
  this.render_encoder_draw_indexed_indirect = (rh, bh, offset) => {
    let encoder = app._extobj.get_obj(rh);
    let buffer = app._extobj.get_obj(bh);
    encoder.drawIndexedIndirect(buffer, offset);
  }

  this.render_encoder_end = (rh) => {
    app._extobj.get_obj(rh).end()
  }

  this.render_encoder_execute_bundles = (rh, render_bundles) => {
    app._extobj.get_obj(rh).executeBundles(app._extobj.get_obj(render_bundles));
  }

  this.render_encoder_finish_bundle = (encoder) => {
    let bundle = app._extobj.get_obj(encoder).finish();
    return app._extobj.insert_obj(bundle);
  }

  this.render_encoder_set_bind_group = (rh, id, bg, dynamicOffsets) => {
    if (dynamicOffsets == 0) {
      app._extobj.get_obj(rh).setBindGroup(id, app._extobj.get_obj(bg));
    } else {
      app._extobj.get_obj(rh).setBindGroup(id, app._extobj.get_obj(bg), app._extobj.get_obj(dynamicOffsets));
    }
  }

  this.render_encoder_set_index_buffer = (rh, buffer, format, offset) => {
    const indexFormat = format ? 'uint32' : 'uint16';
    app._extobj.get_obj(rh).setIndexBuffer(app._extobj.get_obj(buffer), indexFormat, offset);
  }
  this.render_encoder_set_index_buffer_size = (rh, buffer, format, offset, size) => {
    const indexFormat = format ? 'uint32' : 'uint16';
    app._extobj.get_obj(rh).setIndexBuffer(app._extobj.get_obj(buffer), indexFormat, offset, size);
  }

  this.render_encoder_set_pipeline = (rh, ph) => {
    app._extobj.get_obj(rh).setPipeline(app._extobj.get_obj(ph));
  }

  this.render_encoder_set_vertex_buffer = (rh, slot, buffer, offset) => {
    app._extobj.get_obj(rh).setVertexBuffer(slot, app._extobj.get_obj(buffer), offset);
  }
  this.render_encoder_set_vertex_buffer_size = (rh, slot, buffer, offset, size) => {
    app._extobj.get_obj(rh).setVertexBuffer(slot, app._extobj.get_obj(buffer), offset, size);
  }

  this.render_encoder_set_viewport = (rh, x, y, width, height, minDepth, maxDepth) => {
    app._extobj.get_obj(rh).setViewport(x, y, width, height, minDepth, maxDepth);
  }
  
  //---------------------------------------------------------------

  this.renderpipeline_get_bind_group_layout = (pipeline, id) => {
    let layout = app._extobj.get_obj(pipeline).getBindGroupLayout(id);
    return app._extobj.insert_obj(layout);
  }

  //---------------------------------------------------------------

  this.texture_get_depth = (th) => {
    return app._extobj.get_obj(th).depthOrArrayLayers;
  }

  this.texture_get_dimension = (th) => {
    switch (app._extobj.get_obj(th).dimension)
    {
      case "2d":
        return 0;

      case "1d":
        return 1;

      case "3d":
        return 2;
    }
  }

  this.texture_get_format = (th) => {
    return this._texture_format_map.get(app._extobj.get_obj(th).format);
  }

  this.texture_get_width_height = (th) => {
    const t = app._extobj.get_obj(th);
    return [t.width, t.height];
  }

  this.texture_get_mip_level_count = (th) => {
    return app._extobj.get_obj(th).mipLevelCount;
  }

  this.texture_get_sample_count = (th) => {
    return app._extobj.get_obj(th).sampleCount;
  }

  this.texture_get_usage = (th) => {
    return app._extobj.get_obj(th).usage;
  }

  this.texture_create_texture_view = (th) => {
    let view = app._extobj.get_obj(th).createView();
    return app._extobj.insert_obj(view);
  }


},
