!(function(t, n) {
  "object" == typeof exports && "undefined" != typeof module
    ? (module.exports = n())
    : "function" == typeof define && define.amd
    ? define(n)
    : (t.VueWordCloud = n());
})(this, function() {
  "use strict";
  function c(t) {
    return "function" == typeof t;
  }
  function A() {}
  var p = function(t) {
    (this.previousValue = t),
      (this.interrupted = !1),
      (this.interruptHandlers = new Set());
  };
  (p.prototype.throwIfInterrupted = function() {
    if (this.interrupted) throw new Error();
  }),
    (p.prototype.interrupt = function() {
      this.interrupted ||
        ((this.interrupted = !0),
        this.interruptHandlers.forEach(function(t) {
          try {
            t();
          } catch (t) {}
        }));
    }),
    (p.prototype.onInterrupt = function(t) {
      if (this.interrupted && !this.interruptHandlers.has(t))
        try {
          t();
        } catch (t) {}
      this.interruptHandlers.add(t);
    });
  var t = "asyncComputed_",
    h = t + "promise_",
    l = t + "trigger_";
  function n(t) {
    return function() {
      return t;
    };
  }
  function e() {
    return [];
  }
  var r = {
    animationDuration: { type: Number, default: 1e3 },
    animationEasing: { type: String, default: "ease" },
    animationOverlap: { type: Number, default: 1 },
    color: { type: [String, Function], default: "Black" },
    createCanvas: {
      type: Function,
      default: function() {
        return document.createElement("canvas");
      }
    },
    createWorker: {
      type: Function,
      default: function(t) {
        return new Worker(URL.createObjectURL(new Blob([t])));
      }
    },
    enterAnimation: { type: [Object, String], default: n({ opacity: 0 }) },
    fontFamily: { type: [String, Function], default: "serif" },
    fontSizeRatio: { type: Number, default: 0 },
    fontStyle: { type: [String, Function], default: "normal" },
    fontVariant: { type: [String, Function], default: "normal" },
    fontWeight: { type: [String, Function], default: "normal" },
    leaveAnimation: { type: [Object, String], default: n({ opacity: 0 }) },
    loadFont: {
      type: Function,
      default: function(t, n, e, r) {
        return document.fonts.load([n, e, "1px", t].join(" "), r);
      }
    },
    rotation: { type: [Number, Function], default: 0 },
    rotationUnit: { type: [String, Function], default: "turn" },
    spacing: { type: Number, default: 0 },
    text: { type: [String, Function], default: "" },
    weight: { type: [Number, Function], default: 1 },
    words: { type: Array, default: e }
  };
  var s = n(null);
  function P(t) {
    return t && "object" == typeof t;
  }
  function B(t) {
    return "string" == typeof t;
  }
  var i = {
    animationOptions: function() {
      var t,
        r,
        i,
        n = this.enterAnimation,
        e = this.leaveAnimation,
        o = this.animationDuration;
      if (P(n) && P(e)) {
        var a =
            ((t = Object.assign({}, n, e)),
            (r = s),
            (i = {}),
            Object.entries(t).forEach(function(t) {
              var n = t[0],
                e = t[1];
              i[n] = r(e, n);
            }),
            i),
          u = function(t) {
            Object.assign(t.style, n);
          },
          f = function(t, n) {
            setTimeout(function() {
              Object.assign(t.style, a), setTimeout(n, o);
            }, 1);
          };
        return {
          props: { css: !1 },
          on: {
            beforeAppear: u,
            appear: f,
            beforeEnter: u,
            enter: f,
            leave: function(t, n) {
              Object.assign(t.style, e), setTimeout(n, o);
            }
          }
        };
      }
      return B(n) && B(e)
        ? {
            props: {
              duration: o,
              appear: !0,
              appearActiveClass: n,
              enterActiveClass: n,
              leaveActiveClass: e
            }
          }
        : {};
    },
    normalizedAnimationOverlap: function() {
      var t = this.animationOverlap;
      return (t = Math.abs(t)) < 1 && (t = 1 / t), t;
    },
    separateAnimationDelay: function() {
      var t = this.animationDuration,
        n = this.cloudWords,
        e = this.separateAnimationDuration;
      return 1 < n.length ? (t - e) / (n.length - 1) : 0;
    },
    separateAnimationDuration: function() {
      var t = this.animationDuration,
        n = this.normalizedAnimationOverlap,
        e = this.cloudWords;
      return 0 < e.length ? t / Math.min(n, e.length) : 0;
    }
  };
  function D(t) {
    return c(t) ? t : n(t);
  }
  function C(t) {
    return void 0 === t;
  }
  function N(t, n) {
    return (
      t.postMessage(n),
      (o = t),
      new Promise(function(e, r) {
        var i,
          t = function(t) {
            var n = t.data;
            i(), e(n);
          },
          n = function(t) {
            var n = t.error;
            i(), r(n);
          };
        (i = function() {
          o.removeEventListener("message", t),
            o.removeEventListener("error", n);
        }),
          o.addEventListener("message", t),
          o.addEventListener("error", n);
      })
    );
    var o;
  }
  function E(t, n, e) {
    return Math.ceil(t * Math.abs(Math.sin(e)) + n * Math.abs(Math.cos(e)));
  }
  function F(t, n, e) {
    return Math.ceil(t * Math.abs(Math.cos(e)) + n * Math.abs(Math.sin(e)));
  }
  function $(t, n, e, r, i) {
    return [t, n, e, r + "px", i].join(" ");
  }
  function O(t, n) {
    return Math.ceil(t / n) * n;
  }
  function j(t, n, e) {
    var r = e().getContext("2d");
    return (r.font = n), r.measureText(t).width;
  }
  var k = function(t, n, e, r, i, o, a) {
      (this.t = t),
        (this.n = n),
        (this.e = e),
        (this.r = r),
        (this.i = i),
        (this.o = o),
        (this.a = a),
        (this.u = 1),
        (this.f = 0),
        (this.s = 0),
        (this.c = 0);
    },
    o = {
      p: { configurable: !0 },
      h: { configurable: !0 },
      l: { configurable: !0 },
      v: { configurable: !0 },
      d: { configurable: !0 },
      m: { configurable: !0 },
      g: { configurable: !0 },
      T: { configurable: !0 },
      L: { configurable: !0 },
      x: { configurable: !0 },
      b: { configurable: !0 },
      y: { configurable: !0 },
      S: { configurable: !0 },
      M: { configurable: !0 },
      _: { configurable: !0 },
      w: { configurable: !0 },
      W: { configurable: !0 },
      F: { configurable: !0 },
      $: { configurable: !0 }
    };
  (o.p.get = function() {
    return this.u;
  }),
    (o.p.set = function(t) {
      this.u !== t && ((this.u = t), (this.O = void 0));
    }),
    (o.h.get = function() {
      return $(this.o, this.i, this.r, this.p, this.e);
    }),
    (o.l.get = function() {
      return (
        void 0 === this.j &&
          (this.j = j(this.t, $(this.o, this.i, this.r, 1, this.e), this.a)),
        this.j
      );
    }),
    (o.v.get = function() {
      return this.l * this.p;
    }),
    (o.d.get = function() {
      return this.s * this.p;
    }),
    (o.d.set = function(t) {
      this.s = t / this.p;
    }),
    (o.m.get = function() {
      return this.c * this.p;
    }),
    (o.m.set = function(t) {
      this.c = t / this.p;
    }),
    (o.g.get = function() {
      return F(this.v, this.p, this.n);
    }),
    (o.T.get = function() {
      return E(this.v, this.p, this.n);
    }),
    (o.L.get = function() {
      return this.d - this.g / 2;
    }),
    (o.x.get = function() {
      return this.m - this.T / 2;
    }),
    (o.b.get = function() {
      return this.f;
    }),
    (o.b.set = function(t) {
      this.f !== t && ((this.f = t), (this.O = void 0));
    }),
    (o.y.get = function() {
      return (
        void 0 === this.O &&
          (this.O = (function(t, n, e, r, i, o, a, u, f) {
            var s = $(n, e, r, (i *= 4), o),
              c = j(t, s, f),
              p = a * i * 2,
              h = f(),
              l = h.getContext("2d"),
              v = O(F(p + 2 * i + c, p + 3 * i, u), 4),
              d = O(E(p + 2 * i + c, p + 3 * i, u), 4);
            (h.width = v),
              (h.height = d),
              l.translate(v / 2, d / 2),
              l.rotate(u),
              (l.font = s),
              (l.textAlign = "center"),
              (l.textBaseline = "middle"),
              l.fillText(t, 0, 0),
              0 < p &&
                ((l.miterLimit = 1), (l.lineWidth = p), l.strokeText(t, 0, 0));
            for (
              var m = l.getImageData(0, 0, v, d).data,
                g = [],
                T = 1 / 0,
                L = 0,
                x = 1 / 0,
                b = 0,
                y = v / 4,
                S = d / 4,
                M = 0;
              M < y;
              ++M
            )
              for (var _ = 0; _ < S; ++_)
                t: for (var w = 0; w < 4; ++w)
                  for (var W = 0; W < 4; ++W)
                    if (m[4 * (v * (4 * _ + W) + (4 * M + w)) + 3]) {
                      g.push([M, _]),
                        (T = Math.min(M, T)),
                        (L = Math.max(M + 1, L)),
                        (x = Math.min(_, x)),
                        (b = Math.max(_ + 1, b));
                      break t;
                    }
            return 0 < g.length
              ? [
                  g.map(function(t) {
                    var n = t[0],
                      e = t[1];
                    return [n - T, e - x];
                  }),
                  L - T,
                  b - x,
                  Math.ceil(y / 2) - T,
                  Math.ceil(S / 2) - x
                ]
              : [g, 0, 0, 0, 0];
          })(
            this.t,
            this.o,
            this.i,
            this.r,
            this.p,
            this.e,
            this.b,
            this.n,
            this.a
          )),
        this.O
      );
    }),
    (o.S.get = function() {
      return this.y[0];
    }),
    (o.M.get = function() {
      return this.y[1];
    }),
    (o._.get = function() {
      return this.y[2];
    }),
    (o.w.get = function() {
      return this.y[3];
    }),
    (o.W.get = function() {
      return this.y[4];
    }),
    (o.F.get = function() {
      return Math.ceil(this.d) - this.w;
    }),
    (o.F.set = function(t) {
      this.d = t + this.w;
    }),
    (o.$.get = function() {
      return Math.ceil(this.m) - this.W;
    }),
    (o.$.set = function(t) {
      this.m = t + this.W;
    }),
    Object.defineProperties(k.prototype, o);
  var M = "div";
  var a,
    u = {
      name: "VueWordCloud",
      mixins: [
        ((a = {
          cloudWords: {
            get: function(s) {
              var t,
                r = this,
                n = this,
                c = n.elementWidth,
                p = n.elementHeight,
                d = n.words,
                e = n.text,
                i = n.weight,
                o = n.rotation,
                a = n.rotationUnit,
                u = n.fontFamily,
                f = n.fontWeight,
                h = n.fontVariant,
                l = n.fontStyle,
                v = n.color,
                m = n.spacing,
                g = n.fontSizeRatio,
                T = n.createCanvas,
                L = n.loadFont,
                x = n.createWorker;
              (t = g), (g = 1 < (t = Math.abs(t)) ? 1 / t : t);
              var b,
                y,
                S,
                M =
                  ((y = (b = [c, p])[0]),
                  (S = b[1]) < y ? [1, S / y] : y < S ? [y / S, 1] : [1, 1]);
              if (0 < c && 0 < p) {
                var _ = D(e),
                  w = D(i),
                  W = D(o),
                  E = D(a),
                  F = D(u),
                  $ = D(f),
                  O = D(h),
                  j = D(l),
                  H = D(v);
                return (
                  (d = d.map(function(t, n) {
                    var e, r, i, o, a, u, f, s, c, p, h;
                    t &&
                      (B(t)
                        ? (i = t)
                        : Array.isArray(t)
                        ? ((i = (e = t)[0]), (o = e[1]))
                        : P(t) &&
                          ((i = (r = t).text),
                          (o = r.weight),
                          (a = r.rotation),
                          (u = r.rotationUnit),
                          (f = r.fontFamily),
                          (s = r.fontWeight),
                          (c = r.fontVariant),
                          (p = r.fontStyle),
                          (h = r.color))),
                      C(i) && (i = _(t, n, d)),
                      C(o) && (o = w(t, n, d)),
                      C(a) && (a = W(t, n, d)),
                      C(u) && (u = E(t, n, d)),
                      C(f) && (f = F(t, n, d)),
                      C(s) && (s = $(t, n, d)),
                      C(c) && (c = O(t, n, d)),
                      C(p) && (p = j(t, n, d)),
                      C(h) && (h = H(t, n, d));
                    var l = new k(
                      i,
                      (function() {
                        switch (u) {
                          case "turn":
                            return a * Math.PI * 2;
                          case "deg":
                            return (a * Math.PI) / 180;
                        }
                        return a;
                      })(),
                      f,
                      s,
                      c,
                      p,
                      T
                    );
                    return Object.assign(l, { H: t, A: o, P: h }), l;
                  })),
                  Promise.resolve()
                    .then(function() {
                      return Promise.all(
                        d.map(function(t) {
                          var n = t.e,
                            e = t.o,
                            r = t.r,
                            i = t.t;
                          return L(n, e, r, i);
                        })
                      );
                    })
                    .catch(A)
                    .then(function() {
                      if (
                        0 <
                        (d = d
                          .filter(function(t) {
                            return 0 < t.v;
                          })
                          .sort(function(t, n) {
                            return n.A - t.A;
                          })).length
                      ) {
                        var t = d[0],
                          n = (e = d)[e.length - 1],
                          i = t.A,
                          o = n.A;
                        if (o < i) {
                          var a =
                            0 < g
                              ? 1 / g
                              : 0 < o
                              ? i / o
                              : i < 0
                              ? o / i
                              : 1 + i - o;
                          d.forEach(function(t) {
                            var n, e, r;
                            t.p =
                              ((n = t.A),
                              (r = 1) + ((n - (e = o)) * (a - r)) / (i - e));
                          });
                        }
                        d.reduceRight(function(t, n) {
                          return (
                            n.p < 2 * t ? (n.p /= t) : ((t = n.p), (n.p = 1)),
                            (n.B = t)
                          );
                        }, 1),
                          d.forEach(function(t) {
                            t.p *= 4;
                          });
                        var u = x(
                            "(function () {\n\t'use strict';\n\n\tfunction findPixel(ref, ref$1, iteratee) {\n\t\tvar aspectWidth = ref[0];\n\t\tvar aspectHeight = ref[1];\n\t\tvar startLeft = ref$1[0];\n\t\tvar startTop = ref$1[1];\n\n\t\tvar stepLeft, stepTop;\n\t\tif (aspectWidth > aspectHeight) {\n\t\t\tstepLeft = 1;\n\t\t\tstepTop = aspectHeight / aspectWidth;\n\t\t} else\n\t\tif (aspectHeight > aspectWidth) {\n\t\t\tstepTop = 1;\n\t\t\tstepLeft = aspectWidth / aspectHeight;\n\t\t} else {\n\t\t\tstepLeft = stepTop = 1;\n\t\t}\n\n\t\tvar value = [startLeft, startTop];\n\t\tif (iteratee(value)) {\n\t\t\treturn value;\n\t\t}\n\n\t\tvar endLeft = startLeft;\n\t\tvar endTop = startTop;\n\n\t\tvar previousStartLeft = startLeft;\n\t\tvar previousStartTop = startTop;\n\t\tvar previousEndLeft = endLeft;\n\t\tvar previousEndTop = endTop;\n\n\t\tfor (;;) {\n\n\t\t\tstartLeft -= stepLeft;\n\t\t\tstartTop -= stepTop;\n\t\t\tendLeft += stepLeft;\n\t\t\tendTop += stepTop;\n\n\t\t\tvar currentStartLeft = Math.floor(startLeft);\n\t\t\tvar currentStartTop = Math.floor(startTop);\n\t\t\tvar currentEndLeft = Math.ceil(endLeft);\n\t\t\tvar currentEndTop = Math.ceil(endTop);\n\n\t\t\tif (currentEndLeft > previousEndLeft) {\n\t\t\t\tfor (var top = currentStartTop; top < currentEndTop; ++top) {\n\t\t\t\t\tvar value$1 = [currentEndLeft, top];\n\t\t\t\t\tif (iteratee(value$1)) {\n\t\t\t\t\t\treturn value$1;\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\n\t\t\tif (currentEndTop > previousEndTop) {\n\t\t\t\tfor (var left = currentEndLeft; left > currentStartLeft; --left) {\n\t\t\t\t\tvar value$2 = [left, currentEndTop];\n\t\t\t\t\tif (iteratee(value$2)) {\n\t\t\t\t\t\treturn value$2;\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\n\t\t\tif (currentStartLeft < previousStartLeft) {\n\t\t\t\tfor (var top$1 = currentEndTop; top$1 > currentStartTop; --top$1) {\n\t\t\t\t\tvar value$3 = [currentStartLeft, top$1];\n\t\t\t\t\tif (iteratee(value$3)) {\n\t\t\t\t\t\treturn value$3;\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\n\t\t\tif (currentStartTop < previousStartTop) {\n\t\t\t\tfor (var left$1 = currentStartLeft; left$1 < currentEndLeft; ++left$1) {\n\t\t\t\t\tvar value$4 = [left$1, currentStartTop];\n\t\t\t\t\tif (iteratee(value$4)) {\n\t\t\t\t\t\treturn value$4;\n\t\t\t\t\t}\n\t\t\t\t}\n\t\t\t}\n\n\t\t\tpreviousStartLeft = currentStartLeft;\n\t\t\tpreviousStartTop = currentStartTop;\n\t\t\tpreviousEndLeft = currentEndLeft;\n\t\t\tpreviousEndTop = currentEndTop;\n\t\t}\n\t}\n\n\tvar construct = function(event) {\n\t\tself.removeEventListener('message', construct);\n\n\t\tvar _aspect = event.data;\n\t\tvar _pixels;\n\t\tvar _minLeft;\n\t\tvar _maxLeftWidth;\n\t\tvar _minTop;\n\t\tvar _maxTopHeight;\n\n\t\tvar clear = function() {\n\t\t\t_pixels = {};\n\t\t\t_minLeft = 0;\n\t\t\t_maxLeftWidth = 0;\n\t\t\t_minTop = 0;\n\t\t\t_maxTopHeight = 0;\n\t\t};\n\n\t\tclear();\n\n\t\tvar getLeft = function() {\n\t\t\treturn Math.ceil((_minLeft + _maxLeftWidth) / 2);\n\t\t};\n\n\t\tvar getTop = function() {\n\t\t\treturn Math.ceil((_minTop + _maxTopHeight) / 2);\n\t\t};\n\n\t\tvar getWidth = function() {\n\t\t\treturn _maxLeftWidth - _minLeft;\n\t\t};\n\n\t\tvar getHeight = function() {\n\t\t\treturn _maxTopHeight - _minTop;\n\t\t};\n\n\t\tvar getBounds = function() {\n\t\t\treturn {\n\t\t\t\tleft: getLeft(),\n\t\t\t\ttop: getTop(),\n\t\t\t\twidth: getWidth(),\n\t\t\t\theight: getHeight(),\n\t\t\t};\n\t\t};\n\n\t\tvar put = function(pixels, pixelsLeft, pixelsTop) {\n\t\t\tpixels.forEach(function (ref) {\n\t\t\t\tvar pixelLeft = ref[0];\n\t\t\t\tvar pixelTop = ref[1];\n\n\t\t\t\tvar left = pixelsLeft + pixelLeft;\n\t\t\t\tvar top = pixelsTop + pixelTop;\n\t\t\t\t_pixels[(left + \"|\" + top)] = true;\n\t\t\t\t_minLeft = Math.min(left, _minLeft);\n\t\t\t\t_maxLeftWidth = Math.max(left + 1, _maxLeftWidth);\n\t\t\t\t_minTop = Math.min(top, _minTop);\n\t\t\t\t_maxTopHeight = Math.max(top + 1, _maxTopHeight);\n\t\t\t});\n\t\t};\n\n\t\tvar canFit = function(pixels, pixelsLeft, pixelsTop) {\n\t\t\treturn pixels.every(function (ref) {\n\t\t\t\tvar pixelLeft = ref[0];\n\t\t\t\tvar pixelTop = ref[1];\n\n\t\t\t\tvar left = pixelsLeft + pixelLeft;\n\t\t\t\tvar top = pixelsTop + pixelTop;\n\t\t\t\treturn !_pixels[(left + \"|\" + top)];\n\t\t\t});\n\t\t};\n\n\t\tvar findFit = function(pixels, pixelsLeft, pixelsTop) {\n\t\t\treturn findPixel(_aspect, [pixelsLeft + getLeft(), pixelsTop + getTop()], function (ref) {\n\t\t\t\tvar pixelsLeft = ref[0];\n\t\t\t\tvar pixelsTop = ref[1];\n\n\t\t\t\treturn canFit(pixels, pixelsLeft, pixelsTop);\n\t\t\t});\n\t\t};\n\n\t\tself.postMessage({});\n\t\tself.addEventListener('message', function(event) {\n\t\t\tself.postMessage({\n\t\t\t\tgetBounds: getBounds,\n\t\t\t\tput: put,\n\t\t\t\t//canFit,\n\t\t\t\tfindFit: findFit,\n\t\t\t\tclear: clear,\n\t\t\t}[event.data.name].apply(null, event.data.args));\n\t\t});\n\t};\n\tself.addEventListener('message', construct);\n\n}());\n"
                          ),
                          f = { completedWords: 0, totalWords: d.length };
                        return Promise.resolve()
                          .then(function() {
                            return (
                              s.throwIfInterrupted(), (r.progress = f), N(u, M)
                            );
                          })
                          .then(function() {
                            s.throwIfInterrupted(), f.completedWords++;
                            var n = Promise.resolve();
                            return (
                              d.reduce(function(t, r, i) {
                                return (
                                  (n = n
                                    .then(function() {
                                      return r.B < t.B
                                        ? Promise.resolve()
                                            .then(function() {
                                              return N(u, { name: "clear" });
                                            })
                                            .then(function() {
                                              var n = Promise.resolve(),
                                                e = t.B / r.B;
                                              return (
                                                d
                                                  .slice(0, i)
                                                  .forEach(function(t) {
                                                    n = n.then(function() {
                                                      return (
                                                        (t.p *= e),
                                                        N(u, {
                                                          name: "put",
                                                          args: [t.S, t.F, t.$]
                                                        })
                                                      );
                                                    });
                                                  }),
                                                n
                                              );
                                            })
                                        : N(u, {
                                            name: "put",
                                            args: [t.S, t.F, t.$]
                                          });
                                    })
                                    .then(function() {
                                      return (
                                        (r.b = m),
                                        N(u, {
                                          name: "findFit",
                                          args: [r.S, r.F, r.$]
                                        })
                                      );
                                    })
                                    .then(function(t) {
                                      var n = t[0],
                                        e = t[1];
                                      s.throwIfInterrupted(),
                                        f.completedWords++,
                                        (r.F = n),
                                        (r.$ = e),
                                        (r.b = 0);
                                    })),
                                  r
                                );
                              }),
                              n
                            );
                          })
                          .then(function() {
                            return N(u, { name: "put", args: [n.S, n.F, n.$] });
                          })
                          .then(function() {
                            return N(u, { name: "getBounds" });
                          })
                          .then(function(t) {
                            var n = t.left,
                              e = t.top,
                              r = t.width,
                              i = t.height;
                            if (0 < r && 0 < i) {
                              var o = Math.min(c / r, p / i);
                              d.forEach(function(t) {
                                (t.d -= n), (t.m -= e), (t.p *= o);
                              });
                            }
                            var v = new Set();
                            return d.map(function(t) {
                              for (
                                var n = t.H,
                                  e = t.t,
                                  r = t.A,
                                  i = t.n,
                                  o = t.e,
                                  a = t.r,
                                  u = t.i,
                                  f = t.o,
                                  s = t.h,
                                  c = t.d,
                                  p = t.m,
                                  h = t.P,
                                  l = JSON.stringify([e, o, a, u, f]);
                                v.has(l);

                              )
                                l += "!";
                              return (
                                v.add(l),
                                {
                                  key: l,
                                  word: n,
                                  text: e,
                                  weight: r,
                                  rotation: i,
                                  font: s,
                                  color: h,
                                  left: c,
                                  top: p
                                }
                              );
                            });
                          })
                          .finally(function() {
                            u.terminate();
                          })
                          .finally(function() {
                            s.throwIfInterrupted(), (r.progress = null);
                          });
                      }
                      var e;
                      return [];
                    })
                );
              }
              return [];
            },
            default: e
          }
        }),
        {
          data: function() {
            var n = {};
            return (
              Object.keys(a).forEach(function(t) {
                n[l + t] = {};
              }),
              n
            );
          },
          computed: {},
          beforeCreate: function() {
            var e = this,
              s = new Set();
            Object.entries(a).forEach(function(t) {
              var r = t[0],
                n = t[1],
                i = n.get,
                o = n.default,
                a = n.errorHandler;
              void 0 === a && (a = A);
              var u,
                f = !0;
              (e.$options.computed[r] = function() {
                return this[l + r], this[h + r], o;
              }),
                (e.$options.computed[h + r] = function() {
                  var n = this;
                  u && (u.interrupt(), s.delete(u)),
                    f && ((f = !1), c(o) && (o = o.call(this)));
                  var e = new p(o);
                  (u = e),
                    s.add(u),
                    new Promise(function(t) {
                      t(i.call(n, e));
                    })
                      .then(function(t) {
                        e.throwIfInterrupted(), (o = t), (n[l + r] = {});
                      })
                      .catch(a);
                });
            });
          }
        })
      ],
      props: r,
      data: function() {
        return { elementWidth: 0, elementHeight: 0, progress: null };
      },
      computed: i,
      watch: {
        cloudWords: function(t) {
          this.$emit("update:cloudWords", t);
        },
        progress: {
          handler: function(t) {
            this.$emit("update:progress", t);
          },
          deep: !0,
          immediate: !0
        }
      },
      mounted: function() {
        var n,
          e,
          r,
          i,
          t = this;
        (n = function() {
          if (t._isDestroyed) return !1;
          t.updateElementSize();
        }),
          (e = 1e3),
          (r = function(t) {
            requestAnimationFrame(function() {
              !1 !== n() && setTimeout(t, e);
            });
          }),
          (i = function() {
            for (var t = [], n = arguments.length; n--; ) t[n] = arguments[n];
            return r.call.apply(r, [this, i].concat(t));
          })();
      },
      methods: {
        updateElementSize: function() {
          (this.elementWidth = this.$el.offsetWidth),
            (this.elementHeight = this.$el.offsetHeight);
        }
      },
      render: function(T) {
        var t = this,
          L = t.$scopedSlots,
          x = t.animationEasing,
          b = t.animationOptions,
          n = t.cloudWords,
          y = t.separateAnimationDelay,
          S = t.separateAnimationDuration;
        L = Object.assign(
          {},
          {
            default: function(t) {
              return t.text;
            }
          },
          L
        );
        var e = n.map(function(t, n) {
            var e = t.word,
              r = t.key,
              i = t.text,
              o = t.weight,
              a = t.rotation,
              u = t.font,
              f = t.color,
              s = t.left,
              c = t.top,
              p = L.default({
                word: e,
                text: i,
                weight: o,
                font: u,
                color: f,
                left: s,
                top: c
              }),
              h = {
                position: "absolute",
                left: "50%",
                top: "50%",
                color: f,
                font: u,
                whiteSpace: "nowrap",
                transform: [
                  "translate(-50%,-50%)",
                  "rotate(" + a + "rad)"
                ].join(" ")
              },
              l = { position: "absolute", left: s + "px", top: c + "px" };
            if (0 < S) {
              var v = {
                  transitionProperty: "all",
                  transitionDuration: S + "ms",
                  transitionTimingFunction: x,
                  transitionDelay: y * n + "ms"
                },
                d = {
                  animationDuration: S + "ms",
                  animationTimingFunction: x,
                  animationDelay: y * n + "ms"
                };
              Object.assign(h, v), Object.assign(l, v, d);
            }
            var m = T(M, { style: h }, [p]),
              g = T(M, { key: r, style: l }, [m]);
            return T("transition", Object.assign({}, b), [g]);
          }),
          r = T(
            M,
            {
              style: {
                position: "absolute",
                top: "50%",
                left: "50%",
                transform: "translate(-50%,-50%)"
              }
            },
            e
          );
        return T(
          M,
          { style: { position: "relative", width: "100%", height: "100%" } },
          [r]
        );
      }
    };
  return (
    "undefined" != typeof window &&
      window.Vue &&
      window.Vue.component(u.name, u),
    u
  );
});
