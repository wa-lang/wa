canvas: new function() {
  this.set_width_height = (canvas_handle, width, height) => {
    if (canvas_handle == 0) return;
    const canvas = app._extobj.get_obj(canvas_handle);
    canvas.width = width;
    canvas.height = height;
  }

  this.get_context2d = (h) => {
    if (h == 0) return 0;

    const canvas = app._extobj.get_obj(h);
    const ctx = canvas.getContext("2d");
    if (ctx) {
      return app._extobj.insert_obj(ctx);
    }

    return 0;
  }
  this.set_fill_style = (ctx_handle, style_b, style_d, style_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const style = app._mem_util.get_string(style_d, style_l);
    ctx.fillStyle = style;
  }
  this.set_stroke_style = (ctx_handle, style_b, style_d, style_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const style = app._mem_util.get_string(style_d, style_l);
    ctx.strokeStyle = style;
  }
  this.set_shadow_color = (ctx_handle, color_b, color_d, color_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const color = app._mem_util.get_string(color_d, color_l);
    ctx.shadowColor = color;
  }
  this.set_shadow_blur = (ctx_handle, blur) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.shadowBlur = blur;
  }
  this.set_shadow_offset_x = (ctx_handle, offset_x) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.shadowOffsetX = offset_x;
  }
  this.set_shadow_offset_y = (ctx_handle, offset_y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.shadowOffsetY = offset_y;
  }
  this.create_linear_gradient = (ctx_handle, x0, y0, x1, y1) => {
    // TODO: 待完善
  }
  this.create_pattern = (ctx_handle, image_handle, repeat_b, repeat_d, repeat_l) => {
    // TODO: 待完善
  }
  this.create_radial_gradient = (ctx_handle, x0, y0, r0, x1, y1, r1) => {
    // TODO: 待完善
  }
  this.add_color_stop = (ctx_handle, stop, color_b, color_d, color_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const color = app._mem_util.get_string(color_d, color_l);
    ctx.addColorStop(stop, color);
  }
  this.set_line_cap = (ctx_handle, cap_b, cap_d, cap_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const cap = app._mem_util.get_string(cap_d, cap_l);
    ctx.lineCap = cap;
  }
  this.set_line_join = (ctx_handle, join_b, join_d, join_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const join = app._mem_util.get_string(join_d, join_l);
    ctx.lineJoin = join;
  }
  this.set_line_width = (ctx_handle, width) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.lineWidth = width;
  }
  this.set_miter_limit = (ctx_handle, limit) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.miterLimit = limit;
  }
  this.fill = (ctx_handle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.fill();
  }
  this.stroke = (ctx_handle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.stroke();
  }
  this.begin_path = (ctx_handle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.beginPath();
  }
  this.move_to = (ctx_handle, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.moveTo(x, y);
  }
  this.close_path = (ctx_handle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.closePath();
  }
  this.line_to = (ctx_handle, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.lineTo(x, y);
  }
  this.clip = (ctx_handle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.clip();
  }
  this.quadratic_curve_to = (ctx_handle, cpx, cpy, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.quadraticCurveTo(cpx, cpy, x, y);
  }
  this.bezier_curve_to = (ctx_handle, cp1x, cp1y, cp2x, cp2y, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.bezierCurveTo(cp1x, cp1y, cp2x, cp2y, x, y);
  }
  this.arc = (ctx_handle, x, y, r, sAngle, eAngle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.arc(x, y, r, sAngle, eAngle, false);
  }
  this.arc_with_direction = (ctx_handle, x, y, r, sAngle, eAngle, counterclockwise) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.arc(x, y, r, sAngle, eAngle, counterclockwise);
  }

  this.arc_to = (ctx_handle, x1, y1, x2, y2, r) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.arcTo(x1, y1, x2, y2, r);
  }
  this.ellipse = (ctx_handle, x, y, radiusX, radiusY, rotation, startAngle, endAngle, anticlockwise) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.ellipse(x, y, radiusX, radiusY, rotation, startAngle, endAngle, anticlockwise);
  }
  this.is_point_in_path = (ctx_handle, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    return ctx.isPointInPath(x, y);
  }
  this.rect = (ctx_handle, x, y, w, h) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.rect(x, y, w, h);
  }
  this.fill_rect = (ctx_handle, x, y, w, h) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.fillRect(x, y, w, h);
  }
  this.stroke_rect = (ctx_handle, x, y, w, h) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.strokeRect(x, y, w, h);
  }
  this.clear_rect = (ctx_handle, x, y, w, h) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.clearRect(x, y, w, h);
  }
  this.scale = (ctx_handle, scale_width, scale_height) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.scale(scale_width, scale_height);
  }
  this.rotate = (ctx_handle, angle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.rotate(angle);
  }
  this.translate = (ctx_handle, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.translate(x, y);
  }
  this.transform = (ctx_handle, a, b, c, d, e, f) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.transform(a, b, c, d, e, f);
  }
  this.set_transform = (ctx_handle, a, b, c, d, e, f) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.setTransform(a, b, c, d, e, f);
  }
  this.set_font = (ctx_handle, font_b, font_d, font_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const font = app._mem_util.get_string(font_d, font_l);
    ctx.font = font;
  }
  this.set_text_align = (ctx_handle, align_b, align_d, align_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const align = app._mem_util.get_string(align_d, align_l);
    ctx.textAlign = align;
  }
  this.set_text_baseline = (ctx_handle, baseline_b, baseline_d, baseline_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const baseline = app._mem_util.get_string(baseline_d, baseline_l);
    ctx.textBaseline = baseline;
  }
  this.fill_text = (ctx_handle, text_b, text_d, text_l, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const text = app._mem_util.get_string(text_d, text_l);
    ctx.fillText(text, x, y);
  }
  this.fill_text_with_max_width = (ctx_handle, text_b, text_d, text_l, x, y, max_width) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const text = app._mem_util.get_string(text_d, text_l);
    ctx.fillText(text, x, y, max_width);
  }
  this.stroke_text = (ctx_handle, text_b, text_d, text_l, x, y) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const text = app._mem_util.get_string(text_d, text_l);
    ctx.strokeText(text, x, y);
  }
  this.stroke_text_with_max_width = (ctx_handle, text_b, text_d, text_l, x, y, max_width) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const text = app._mem_util.get_string(text_d, text_l);
    ctx.strokeText(text, x, y, max_width);
  }
  this.measure_text = (ctx_handle, text_b, text_d, text_l) => {
    // TODO: 待完善
  }
  this.measure_text_with_max_width = (ctx_handle, text_b, text_d, text_l, max_width) => {
    // TODO: 待完善
  }
  this.draw_image = (ctx_handle, img_data_b, img_data_d, img_data_l, img_data_c, dx, dy, dwidth, dheight) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const buf = new Uint8ClampedArray(app._mem_util.mem().buffer, img_data_d, img_data_l);
    const image = new ImageData(buf, dwidth, dheight);
    ctx.drawImage(image, dx, dy, dwidth, dheight);
  }
  this.draw_image_cropped = (ctx_handle, img_data_b, img_data_d, img_data_l, img_data_c, sx, sy, swidth, sheight, dx, dy, dwidth, dheight) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const buf = new Uint8ClampedArray(app._mem_util.mem().buffer, img_data_d, img_data_l);
    const image = new ImageData(buf, swidth, sheight);
    ctx.drawImage(image, sx, sy, swidth, sheight, dx, dy, dwidth, dheight);
  }
  this.create_image_data = (ctx_handle, w, h) => {
    // TODO: 待完善
  }
  this.get_image_data = (ctx_handle, x, y, w, h, buf_b, buf_d, buf_l, buf_c) => {
    // TODO: 待完善
  }
  this.put_image_data = (ctx_handle, img_data_b, img_data_d, img_data_l, img_data_c, dx, dy, dirty_x, dirty_y, dirty_w, dirty_h) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const buf = new Uint8ClampedArray(app._mem_util.mem().buffer, img_data_d, img_data_l);
    const image = new ImageData(buf, dirty_w, dirty_h);
    ctx.putImageData(image, dx, dy, dirty_x, dirty_y, dirty_w, dirty_h);
  }
  this.set_global_alpha = (ctx_handle, alpha) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.globalAlpha = alpha;
  }
  this.set_global_composite_operation = (ctx_handle, op_b, op_d, op_l) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    const op = app._mem_util.get_string(op_d, op_l);
    ctx.globalCompositeOperation = op;
  }
  this.save = (ctx_handle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.save();
  }
  this.restore = (ctx_handle) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.restore();
  }

  app._canvas = this;
},