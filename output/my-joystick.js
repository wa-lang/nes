const MyJoystick = function (obj, isRotate = false) {
  this.hammertime = {};
  // 回调信息
  this.callbackInfo = {
    direction: null,
    effective: null,
    status: null,
  };
  this.isRotate = isRotate;
  // 配置
  const config = {
    frontColor: obj.frontColor || "red",
    backColor: obj.backColor || "red",
    backR: obj.backR || 100,
    frontR: obj.frontR || 50,
    distanceX: obj.distanceX || "50%",
    distanceY: obj.distanceY || "50%",
  };

  // 创建摇杆 DOM：back、front
  MyJoystick.prototype.buildUI = function () {
    this.ui = {};
    // 创建 dom 元素
    this.ui.container = document.createElement("div");
    this.ui.back = document.createElement("div");
    this.ui.front = document.createElement("div");

    // 设置 dom 元素对应的类名
    this.ui.container.className = "container";
    this.ui.back.className = "back";
    this.ui.front.className = "front";

    // 给元素添加 id
    this.ui.front.setAttribute("id", "front");

    // 将 back、front 加入到 container 中
    this.ui.container.appendChild(this.ui.back);
    this.ui.container.appendChild(this.ui.front);

    // 将 dom 元素添加至 body
    const parentDom = document.querySelector('.analog')
    parentDom.appendChild(this.ui.container);

    this.stylizeUI();
    this.addArrows();
    // this.tipText();
  };

  // 定义样式
  MyJoystick.prototype.stylizeUI = function () {
    var styles = {};
    styles.container = {
      position: "absolute",
      // opacity: 0.6,
      top: config.distanceY,
      left: config.distanceX,
    };
    styles.back = {
      position: "absolute",
      display: "block",
      width: config.backR + "px",
      height: config.backR + "px",
      marginLeft: -config.backR / 2 + "px",
      marginTop: -config.backR / 2 + "px",
      background: config.backColor,
      // opacity: 0.8,
      borderRadius: "50%",
    };
    styles.front = {
      width: config.frontR + "px",
      height: config.frontR + "px",
      position: "absolute",
      display: "block",
      marginLeft: -config.frontR / 2 + "px",
      marginTop: -config.frontR / 2 + "px",
      background: config.frontColor,
      // opacity: 0.5,
      borderRadius: "50%",
      left: 0 + "px",
      top: 0 + "px",
    };
    this.applyStyles(styles);
  };

  // 应用样式
  MyJoystick.prototype.applyStyles = function (styles) {
    for (var i in this.ui) {
      for (var j in styles[i]) {
        this.ui[i].style[j] = styles[i][j];
      }
    }
  };

  // 添加箭头
  MyJoystick.prototype.addArrows = function () {
    // 创建 img dom
    this.arrowsGroup = {};
    this.arrowsGroup.upArrow = document.createElement("img");
    this.arrowsGroup.downArrow = document.createElement("img");
    this.arrowsGroup.leftArrow = document.createElement("img");
    this.arrowsGroup.rightArrow = document.createElement("img");

    // 添加 dom，设置公共属性
    for (let i in this.arrowsGroup) {
      this.ui.back.appendChild(this.arrowsGroup[i]);
      this.arrowsGroup[i].setAttribute("src", "./arrow.png");
      this.arrowsGroup[i].style.width = config.backR / 5 + "px";
      this.arrowsGroup[i].style.position = "absolute";
    }
    // 箭头样式
    const styleUp = {
      top: 2 + "px",
      left: "40%",
      opacity: 0.5,
    };
    const styleDown = {
      left: "40%",
      bottom: 2 + "px",
      transform: "rotate(180deg)",
      opacity: 0.5,
    };
    const styleLeft = {
      left: 2 + "px",
      top: "40%",
      transform: "rotate(-90deg)",
      opacity: 0.5,
    };
    const styleRight = {
      right: 2 + "px",
      top: "40%",
      transform: "rotate(90deg)",
      opacity: 0.5,
    };
    // 使用自定义函数extend合并style对象
    extend(this.arrowsGroup.upArrow.style, styleUp);
    extend(this.arrowsGroup.downArrow.style, styleDown);
    extend(this.arrowsGroup.leftArrow.style, styleLeft);
    extend(this.arrowsGroup.rightArrow.style, styleRight);
  };

  // 启动摇杆
  MyJoystick.prototype.openJoystick = function (callback) {
    this.buildUI();
    const test = document.getElementById("front");
    this.hammertime = new Hammer(test);
    this.hammertime
      .get("pan")
      .set({ direction: Hammer.DIRECTION_ALL, threshold: 0 });

    var that = this;

    // 按压事件
    this.hammertime.on("press", function (ev) {
      // console.log("press:", ev);
      that.callbackInfo.status = "press-down";
      var { dis, angle } = that.backToCenter(ev);
      that.isEffective(ev, dis, angle, callback);
    });
    // 按压事件抬起
    this.hammertime.on("pressup", function (ev) {
      // console.log("press:", ev);
      that.callbackInfo.status = "press-up";
      that.callbackInfo.effective = false;
      that.callbackInfo.direction = null;
      callback(that.callbackInfo);
      test.style.top = 0;
      test.style.left = 0;
    });
    // 轻触事件
    this.hammertime.on("tap", function (ev) {
      // console.log("tap:", ev);
      var { dis, angle } = that.backToCenter(ev);
      that.callbackInfo.status = "tap";
      that.isEffective(ev, dis, angle, callback);
      setTimeout(function () {
        that.callbackInfo.status = "tap-end";
        that.callbackInfo.effective = false;
        that.callbackInfo.direction = null;
        callback(that.callbackInfo);
        test.style.top = 0;
        test.style.left = 0;
      }, 200);
    });
    // 移动事件开始
    this.hammertime.on("panstart", function (ev) {
      // console.log("panstart:", ev);
      that.callbackInfo.status = "move-start";
      callback(that.callbackInfo);
    });
    // 移动事件
    this.hammertime.on("panmove", function (ev) {
      that.callbackInfo.status = "move";
      var { dis, angle } = that.backToCenter(ev);
      if (that.isRotate) {
        angle = (angle + 90) % 360;
      }
      that.isEffective(ev, dis, angle, callback);
    });
    // 移动事件结束
    this.hammertime.on("panend", function (ev) {
      // console.log("panend:", ev);
      that.callbackInfo.status = "move-end";
      that.callbackInfo.effective = false;
      that.callbackInfo.direction = null;
      callback(that.callbackInfo);
      test.style.top = 0;
      test.style.left = 0;
    });
  };

  // 死区的判断
  MyJoystick.prototype.isEffective = function (ev, distance, angle, callback) {
    if (distance <= 10) {
      // console.log("无效范围！", ev);
      // tipText("无效范围");
      this.callbackInfo.direction = null;
      this.callbackInfo.effective = false;
      callback(this.callbackInfo);
    } else {
      // console.log("有效：", ev);
      var dir = getDirectionByAngle(angle, this.isRotate);
      this.callbackInfo.direction = dir;
      this.callbackInfo.effective = true;
      callback(this.callbackInfo);
    }
    // this.tipText(this.callbackInfo.direction);
  };

  // 修正 center 位置，并返回距离圆盘中心的距离和角度
  MyJoystick.prototype.backToCenter = function (ev) {
    var res = this.ui.container.getBoundingClientRect();
    var dx = ev.center.x - res.x;
    var dy = ev.center.y - res.y;
    if (this.isRotate) {
      [dx, dy] = [dy, -dx];
    }
    var dis = Math.sqrt(dx * dx + dy * dy);
    var angle = Math.atan2(dy, dx) * (180 / Math.PI);
    // console.log(angle);
    this.ui.front.style.top = dy + "px";
    this.ui.front.style.left = dx + "px";
    if (dis >= config.backR / 2) {
      var dyM = (dy * (config.backR / 2)) / dis;
      var dxM = (dx * (config.backR / 2)) / dis;
      this.ui.front.style.top = dyM + "px";
      this.ui.front.style.left = dxM + "px";
    }
    return { dis, angle };
  };

  // 提示文本
  MyJoystick.prototype.tipText = function (text = "这里是提示信息！") {
    if (this.ui.tip) {
      this.ui.tip.innerHTML = `<h1>${text}</h1>`;
    } else {
      this.ui.tip = document.createElement("div");
      this.ui.tip.className = "tip";
      const bodyDom = document.body;
      bodyDom.appendChild(this.ui.tip);
      this.ui.tip.innerHTML = `<h1>${text}</h1>`;
    }
  };

  // 根据 angle 判断方向
  const getDirectionByAngle = function (angle, isRotate) {
    var direction = "";
    if (angle < -45 && angle >= -135) {
      // this.tipText("上");
      direction = isRotate ? "left" : "up";
    } else if (angle < -135 || angle >= 135) {
      // this.tipText("左");
      direction = isRotate ? "down" : "left";
    } else if (angle >= 45 && angle < 135) {
      // this.tipText("下");
      direction = isRotate ? "right" : "down";
    } else if (angle < 45 && angle >= -45) {
      // this.tipText("右");
      direction = isRotate ? "up" : "right";
    }
    return direction;
  };

  // 合并对象
  const extend = function (objA, objB) {
    for (let i in objB) {
      if (objB.hasOwnProperty(i)) {
        objA[i] = objB[i];
      }
    }
    return objA;
  };
};
