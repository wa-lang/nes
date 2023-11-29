canvas: new function() {
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

  this.fill_rect = (ctx_handle, x, y, w, h) => {
    if (ctx_handle == 0) return;

    const ctx = app._extobj.get_obj(ctx_handle);
    ctx.fillRect(x, y, w, h);
  }

  app._canvas = this;
},