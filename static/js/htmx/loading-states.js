!(function () {
  let t = [];
  function n(t, n) {
    document.body.contains(t) && n();
  }
  function a(a, o, i, c) {
    const e = htmx.closest(a, "[data-loading-delay]");
    if (e) {
      const a = e.getAttribute("data-loading-delay") || 200,
        f = setTimeout(function () {
          i(),
            t.push(function () {
              n(o, c);
            });
        }, a);
      t.push(function () {
        n(o, function () {
          clearTimeout(f);
        });
      });
    } else
      i(),
        t.push(function () {
          n(o, c);
        });
  }
  function o(t, n, a) {
    return Array.from(htmx.findAll(t, "[" + n + "]")).filter(function (t) {
      return (function (t, n) {
        const a = htmx.closest(t, "[data-loading-path]");
        return !a || a.getAttribute("data-loading-path") === n;
      })(t, a);
    });
  }
  function i(t) {
    return t.getAttribute("data-loading-target")
      ? Array.from(htmx.findAll(t.getAttribute("data-loading-target")))
      : [t];
  }
  htmx.defineExtension("loading-states", {
    onEvent: function (n, c) {
      if ("htmx:beforeRequest" === n) {
        const t =
          ((e = c.target),
          htmx.closest(e, "[data-loading-states]") || document.body);
        let n = {};
        [
          "data-loading",
          "data-loading-class",
          "data-loading-class-remove",
          "data-loading-disable",
          "data-loading-aria-busy",
        ].forEach(function (a) {
          n[a] = o(t, a, c.detail.pathInfo.requestPath);
        }),
          n["data-loading"].forEach(function (t) {
            i(t).forEach(function (n) {
              a(
                t,
                n,
                function () {
                  n.style.display =
                    t.getAttribute("data-loading") || "inline-block";
                },
                function () {
                  n.style.display = "none";
                },
              );
            });
          }),
          n["data-loading-class"].forEach(function (t) {
            const n = t.getAttribute("data-loading-class").split(" ");
            i(t).forEach(function (o) {
              a(
                t,
                o,
                function () {
                  n.forEach(function (t) {
                    o.classList.add(t);
                  });
                },
                function () {
                  n.forEach(function (t) {
                    o.classList.remove(t);
                  });
                },
              );
            });
          }),
          n["data-loading-class-remove"].forEach(function (t) {
            const n = t.getAttribute("data-loading-class-remove").split(" ");
            i(t).forEach(function (o) {
              a(
                t,
                o,
                function () {
                  n.forEach(function (t) {
                    o.classList.remove(t);
                  });
                },
                function () {
                  n.forEach(function (t) {
                    o.classList.add(t);
                  });
                },
              );
            });
          }),
          n["data-loading-disable"].forEach(function (t) {
            i(t).forEach(function (n) {
              a(
                t,
                n,
                function () {
                  n.disabled = !0;
                },
                function () {
                  n.disabled = !1;
                },
              );
            });
          }),
          n["data-loading-aria-busy"].forEach(function (t) {
            i(t).forEach(function (n) {
              a(
                t,
                n,
                function () {
                  n.setAttribute("aria-busy", "true");
                },
                function () {
                  n.removeAttribute("aria-busy");
                },
              );
            });
          });
      }
      var e;
      if ("htmx:beforeOnLoad" === n) for (; t.length > 0; ) t.shift()();
    },
  });
})();
