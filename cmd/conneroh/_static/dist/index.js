var M0 = (function () {
    let htmx = {
      onLoad: null,
      process: null,
      on: null,
      off: null,
      trigger: null,
      ajax: null,
      find: null,
      findAll: null,
      closest: null,
      values: function (J, Y) {
        return getInputValues(J, Y || "post").values;
      },
      remove: null,
      addClass: null,
      removeClass: null,
      toggleClass: null,
      takeClass: null,
      swap: null,
      defineExtension: null,
      removeExtension: null,
      logAll: null,
      logNone: null,
      logger: null,
      config: {
        historyEnabled: !0,
        historyCacheSize: 10,
        refreshOnHistoryMiss: !1,
        defaultSwapStyle: "innerHTML",
        defaultSwapDelay: 0,
        defaultSettleDelay: 20,
        includeIndicatorStyles: !0,
        indicatorClass: "htmx-indicator",
        requestClass: "htmx-request",
        addedClass: "htmx-added",
        settlingClass: "htmx-settling",
        swappingClass: "htmx-swapping",
        allowEval: !0,
        allowScriptTags: !0,
        inlineScriptNonce: "",
        inlineStyleNonce: "",
        attributesToSettle: ["class", "style", "width", "height"],
        withCredentials: !1,
        timeout: 0,
        wsReconnectDelay: "full-jitter",
        wsBinaryType: "blob",
        disableSelector: "[hx-disable], [data-hx-disable]",
        scrollBehavior: "instant",
        defaultFocusScroll: !1,
        getCacheBusterParam: !1,
        globalViewTransitions: !1,
        methodsThatUseUrlParams: ["get", "delete"],
        selfRequestsOnly: !0,
        ignoreTitle: !1,
        scrollIntoViewOnBoost: !0,
        triggerSpecsCache: null,
        disableInheritance: !1,
        responseHandling: [
          { code: "204", swap: !1 },
          { code: "[23]..", swap: !0 },
          { code: "[45]..", swap: !1, error: !0 },
        ],
        allowNestedOobSwaps: !0,
      },
      parseInterval: null,
      _: null,
      version: "2.0.4",
    };
    (htmx.onLoad = onLoadHelper),
      (htmx.process = processNode),
      (htmx.on = addEventListenerImpl),
      (htmx.off = removeEventListenerImpl),
      (htmx.trigger = triggerEvent),
      (htmx.ajax = ajaxHelper),
      (htmx.find = find),
      (htmx.findAll = findAll),
      (htmx.closest = closest),
      (htmx.remove = removeElement),
      (htmx.addClass = addClassToElement),
      (htmx.removeClass = removeClassFromElement),
      (htmx.toggleClass = toggleClassOnElement),
      (htmx.takeClass = takeClassForElement),
      (htmx.swap = swap),
      (htmx.defineExtension = defineExtension),
      (htmx.removeExtension = removeExtension),
      (htmx.logAll = logAll),
      (htmx.logNone = logNone),
      (htmx.parseInterval = parseInterval),
      (htmx._ = internalEval);
    let internalAPI = {
        addTriggerHandler,
        bodyContains,
        canAccessLocalStorage,
        findThisElement,
        filterValues,
        swap,
        hasAttribute,
        getAttributeValue,
        getClosestAttributeValue,
        getClosestMatch,
        getExpressionVars,
        getHeaders,
        getInputValues,
        getInternalData,
        getSwapSpecification,
        getTriggerSpecs,
        getTarget,
        makeFragment,
        mergeObjects,
        makeSettleInfo,
        oobSwap,
        querySelectorExt,
        settleImmediately,
        shouldCancel,
        triggerEvent,
        triggerErrorEvent,
        withExtensions,
      },
      VERBS = ["get", "post", "put", "delete", "patch"],
      VERB_SELECTOR = VERBS.map(function (J) {
        return "[hx-" + J + "], [data-hx-" + J + "]";
      }).join(", ");
    function parseInterval(J) {
      if (J == null) return;
      let Y = NaN;
      if (J.slice(-2) == "ms") Y = parseFloat(J.slice(0, -2));
      else if (J.slice(-1) == "s") Y = parseFloat(J.slice(0, -1)) * 1000;
      else if (J.slice(-1) == "m") Y = parseFloat(J.slice(0, -1)) * 1000 * 60;
      else Y = parseFloat(J);
      return isNaN(Y) ? void 0 : Y;
    }
    function getRawAttribute(J, Y) {
      return J instanceof Element && J.getAttribute(Y);
    }
    function hasAttribute(J, Y) {
      return (
        !!J.hasAttribute && (J.hasAttribute(Y) || J.hasAttribute("data-" + Y))
      );
    }
    function getAttributeValue(J, Y) {
      return getRawAttribute(J, Y) || getRawAttribute(J, "data-" + Y);
    }
    function parentElt(J) {
      let Y = J.parentElement;
      if (!Y && J.parentNode instanceof ShadowRoot) return J.parentNode;
      return Y;
    }
    function getDocument() {
      return document;
    }
    function getRootNode(J, Y) {
      return J.getRootNode ? J.getRootNode({ composed: Y }) : getDocument();
    }
    function getClosestMatch(J, Y) {
      while (J && !Y(J)) J = parentElt(J);
      return J || null;
    }
    function getAttributeValueWithDisinheritance(J, Y, Z) {
      let $ = getAttributeValue(Y, Z),
        W = getAttributeValue(Y, "hx-disinherit");
      var X = getAttributeValue(Y, "hx-inherit");
      if (J !== Y) {
        if (htmx.config.disableInheritance)
          if (X && (X === "*" || X.split(" ").indexOf(Z) >= 0)) return $;
          else return null;
        if (W && (W === "*" || W.split(" ").indexOf(Z) >= 0)) return "unset";
      }
      return $;
    }
    function getClosestAttributeValue(J, Y) {
      let Z = null;
      if (
        (getClosestMatch(J, function ($) {
          return !!(Z = getAttributeValueWithDisinheritance(
            J,
            asElement($),
            Y,
          ));
        }),
        Z !== "unset")
      )
        return Z;
    }
    function matches(J, Y) {
      let Z =
        J instanceof Element &&
        (J.matches ||
          J.matchesSelector ||
          J.msMatchesSelector ||
          J.mozMatchesSelector ||
          J.webkitMatchesSelector ||
          J.oMatchesSelector);
      return !!Z && Z.call(J, Y);
    }
    function getStartTag(J) {
      let Z = /<([a-z][^\/\0>\x20\t\r\n\f]*)/i.exec(J);
      if (Z) return Z[1].toLowerCase();
      else return "";
    }
    function parseHTML(J) {
      return new DOMParser().parseFromString(J, "text/html");
    }
    function takeChildrenFor(J, Y) {
      while (Y.childNodes.length > 0) J.append(Y.childNodes[0]);
    }
    function duplicateScript(J) {
      let Y = getDocument().createElement("script");
      if (
        (forEach(J.attributes, function (Z) {
          Y.setAttribute(Z.name, Z.value);
        }),
        (Y.textContent = J.textContent),
        (Y.async = !1),
        htmx.config.inlineScriptNonce)
      )
        Y.nonce = htmx.config.inlineScriptNonce;
      return Y;
    }
    function isJavaScriptScriptNode(J) {
      return (
        J.matches("script") &&
        (J.type === "text/javascript" || J.type === "module" || J.type === "")
      );
    }
    function normalizeScriptTags(J) {
      Array.from(J.querySelectorAll("script")).forEach((Y) => {
        if (isJavaScriptScriptNode(Y)) {
          let Z = duplicateScript(Y),
            $ = Y.parentNode;
          try {
            $.insertBefore(Z, Y);
          } catch (W) {
            logError(W);
          } finally {
            Y.remove();
          }
        }
      });
    }
    function makeFragment(J) {
      let Y = J.replace(/<head(\s[^>]*)?>[\s\S]*?<\/head>/i, ""),
        Z = getStartTag(Y),
        $;
      if (Z === "html") {
        $ = new DocumentFragment();
        let X = parseHTML(J);
        takeChildrenFor($, X.body), ($.title = X.title);
      } else if (Z === "body") {
        $ = new DocumentFragment();
        let X = parseHTML(Y);
        takeChildrenFor($, X.body), ($.title = X.title);
      } else {
        let X = parseHTML(
          '<body><template class="internal-htmx-wrapper">' +
            Y +
            "</template></body>",
        );
        ($ = X.querySelector("template").content), ($.title = X.title);
        var W = $.querySelector("title");
        if (W && W.parentNode === $) W.remove(), ($.title = W.innerText);
      }
      if ($)
        if (htmx.config.allowScriptTags) normalizeScriptTags($);
        else $.querySelectorAll("script").forEach((X) => X.remove());
      return $;
    }
    function maybeCall(J) {
      if (J) J();
    }
    function isType(J, Y) {
      return Object.prototype.toString.call(J) === "[object " + Y + "]";
    }
    function isFunction(J) {
      return typeof J === "function";
    }
    function isRawObject(J) {
      return isType(J, "Object");
    }
    function getInternalData(J) {
      let Z = J["htmx-internal-data"];
      if (!Z) Z = J["htmx-internal-data"] = {};
      return Z;
    }
    function toArray(J) {
      let Y = [];
      if (J) for (let Z = 0; Z < J.length; Z++) Y.push(J[Z]);
      return Y;
    }
    function forEach(J, Y) {
      if (J) for (let Z = 0; Z < J.length; Z++) Y(J[Z]);
    }
    function isScrolledIntoView(J) {
      let Y = J.getBoundingClientRect(),
        Z = Y.top,
        $ = Y.bottom;
      return Z < window.innerHeight && $ >= 0;
    }
    function bodyContains(J) {
      return J.getRootNode({ composed: !0 }) === document;
    }
    function splitOnWhitespace(J) {
      return J.trim().split(/\s+/);
    }
    function mergeObjects(J, Y) {
      for (let Z in Y) if (Y.hasOwnProperty(Z)) J[Z] = Y[Z];
      return J;
    }
    function parseJSON(J) {
      try {
        return JSON.parse(J);
      } catch (Y) {
        return logError(Y), null;
      }
    }
    function canAccessLocalStorage() {
      try {
        return (
          localStorage.setItem(
            "htmx:localStorageTest",
            "htmx:localStorageTest",
          ),
          localStorage.removeItem("htmx:localStorageTest"),
          !0
        );
      } catch (Y) {
        return !1;
      }
    }
    function normalizePath(J) {
      try {
        let Y = new URL(J);
        if (Y) J = Y.pathname + Y.search;
        if (!/^\/$/.test(J)) J = J.replace(/\/+$/, "");
        return J;
      } catch (Y) {
        return J;
      }
    }
    function internalEval(str) {
      return maybeEval(getDocument().body, function () {
        return eval(str);
      });
    }
    function onLoadHelper(J) {
      return htmx.on("htmx:load", function (Z) {
        J(Z.detail.elt);
      });
    }
    function logAll() {
      htmx.logger = function (J, Y, Z) {
        if (console) console.log(Y, J, Z);
      };
    }
    function logNone() {
      htmx.logger = null;
    }
    function find(J, Y) {
      if (typeof J !== "string") return J.querySelector(Y);
      else return find(getDocument(), J);
    }
    function findAll(J, Y) {
      if (typeof J !== "string") return J.querySelectorAll(Y);
      else return findAll(getDocument(), J);
    }
    function getWindow() {
      return window;
    }
    function removeElement(J, Y) {
      if (((J = resolveTarget(J)), Y))
        getWindow().setTimeout(function () {
          removeElement(J), (J = null);
        }, Y);
      else parentElt(J).removeChild(J);
    }
    function asElement(J) {
      return J instanceof Element ? J : null;
    }
    function asHtmlElement(J) {
      return J instanceof HTMLElement ? J : null;
    }
    function asString(J) {
      return typeof J === "string" ? J : null;
    }
    function asParentNode(J) {
      return J instanceof Element ||
        J instanceof Document ||
        J instanceof DocumentFragment
        ? J
        : null;
    }
    function addClassToElement(J, Y, Z) {
      if (((J = asElement(resolveTarget(J))), !J)) return;
      if (Z)
        getWindow().setTimeout(function () {
          addClassToElement(J, Y), (J = null);
        }, Z);
      else J.classList && J.classList.add(Y);
    }
    function removeClassFromElement(J, Y, Z) {
      let $ = asElement(resolveTarget(J));
      if (!$) return;
      if (Z)
        getWindow().setTimeout(function () {
          removeClassFromElement($, Y), ($ = null);
        }, Z);
      else if ($.classList) {
        if (($.classList.remove(Y), $.classList.length === 0))
          $.removeAttribute("class");
      }
    }
    function toggleClassOnElement(J, Y) {
      (J = resolveTarget(J)), J.classList.toggle(Y);
    }
    function takeClassForElement(J, Y) {
      (J = resolveTarget(J)),
        forEach(J.parentElement.children, function (Z) {
          removeClassFromElement(Z, Y);
        }),
        addClassToElement(asElement(J), Y);
    }
    function closest(J, Y) {
      if (((J = asElement(resolveTarget(J))), J && J.closest))
        return J.closest(Y);
      else {
        do if (J == null || matches(J, Y)) return J;
        while ((J = J && asElement(parentElt(J))));
        return null;
      }
    }
    function startsWith(J, Y) {
      return J.substring(0, Y.length) === Y;
    }
    function endsWith(J, Y) {
      return J.substring(J.length - Y.length) === Y;
    }
    function normalizeSelector(J) {
      let Y = J.trim();
      if (startsWith(Y, "<") && endsWith(Y, "/>"))
        return Y.substring(1, Y.length - 2);
      else return Y;
    }
    function querySelectorAllExt(J, Y, Z) {
      if (Y.indexOf("global ") === 0)
        return querySelectorAllExt(J, Y.slice(7), !0);
      J = resolveTarget(J);
      let $ = [];
      {
        let Q = 0,
          B = 0;
        for (let U = 0; U < Y.length; U++) {
          let _ = Y[U];
          if (_ === "," && Q === 0) {
            $.push(Y.substring(B, U)), (B = U + 1);
            continue;
          }
          if (_ === "<") Q++;
          else if (_ === "/" && U < Y.length - 1 && Y[U + 1] === ">") Q--;
        }
        if (B < Y.length) $.push(Y.substring(B));
      }
      let W = [],
        X = [];
      while ($.length > 0) {
        let Q = normalizeSelector($.shift()),
          B;
        if (Q.indexOf("closest ") === 0)
          B = closest(asElement(J), normalizeSelector(Q.substr(8)));
        else if (Q.indexOf("find ") === 0)
          B = find(asParentNode(J), normalizeSelector(Q.substr(5)));
        else if (Q === "next" || Q === "nextElementSibling")
          B = asElement(J).nextElementSibling;
        else if (Q.indexOf("next ") === 0)
          B = scanForwardQuery(J, normalizeSelector(Q.substr(5)), !!Z);
        else if (Q === "previous" || Q === "previousElementSibling")
          B = asElement(J).previousElementSibling;
        else if (Q.indexOf("previous ") === 0)
          B = scanBackwardsQuery(J, normalizeSelector(Q.substr(9)), !!Z);
        else if (Q === "document") B = document;
        else if (Q === "window") B = window;
        else if (Q === "body") B = document.body;
        else if (Q === "root") B = getRootNode(J, !!Z);
        else if (Q === "host") B = J.getRootNode().host;
        else X.push(Q);
        if (B) W.push(B);
      }
      if (X.length > 0) {
        let Q = X.join(","),
          B = asParentNode(getRootNode(J, !!Z));
        W.push(...toArray(B.querySelectorAll(Q)));
      }
      return W;
    }
    var scanForwardQuery = function (J, Y, Z) {
        let $ = asParentNode(getRootNode(J, Z)).querySelectorAll(Y);
        for (let W = 0; W < $.length; W++) {
          let X = $[W];
          if (X.compareDocumentPosition(J) === Node.DOCUMENT_POSITION_PRECEDING)
            return X;
        }
      },
      scanBackwardsQuery = function (J, Y, Z) {
        let $ = asParentNode(getRootNode(J, Z)).querySelectorAll(Y);
        for (let W = $.length - 1; W >= 0; W--) {
          let X = $[W];
          if (X.compareDocumentPosition(J) === Node.DOCUMENT_POSITION_FOLLOWING)
            return X;
        }
      };
    function querySelectorExt(J, Y) {
      if (typeof J !== "string") return querySelectorAllExt(J, Y)[0];
      else return querySelectorAllExt(getDocument().body, J)[0];
    }
    function resolveTarget(J, Y) {
      if (typeof J === "string") return find(asParentNode(Y) || document, J);
      else return J;
    }
    function processEventArgs(J, Y, Z, $) {
      if (isFunction(Y))
        return {
          target: getDocument().body,
          event: asString(J),
          listener: Y,
          options: Z,
        };
      else
        return {
          target: resolveTarget(J),
          event: asString(Y),
          listener: Z,
          options: $,
        };
    }
    function addEventListenerImpl(J, Y, Z, $) {
      return (
        ready(function () {
          let X = processEventArgs(J, Y, Z, $);
          X.target.addEventListener(X.event, X.listener, X.options);
        }),
        isFunction(Y) ? Y : Z
      );
    }
    function removeEventListenerImpl(J, Y, Z) {
      return (
        ready(function () {
          let $ = processEventArgs(J, Y, Z);
          $.target.removeEventListener($.event, $.listener);
        }),
        isFunction(Y) ? Y : Z
      );
    }
    let DUMMY_ELT = getDocument().createElement("output");
    function findAttributeTargets(J, Y) {
      let Z = getClosestAttributeValue(J, Y);
      if (Z)
        if (Z === "this") return [findThisElement(J, Y)];
        else {
          let $ = querySelectorAllExt(J, Z);
          if ($.length === 0)
            return (
              logError(
                'The selector "' + Z + '" on ' + Y + " returned no matches!",
              ),
              [DUMMY_ELT]
            );
          else return $;
        }
    }
    function findThisElement(J, Y) {
      return asElement(
        getClosestMatch(J, function (Z) {
          return getAttributeValue(asElement(Z), Y) != null;
        }),
      );
    }
    function getTarget(J) {
      let Y = getClosestAttributeValue(J, "hx-target");
      if (Y)
        if (Y === "this") return findThisElement(J, "hx-target");
        else return querySelectorExt(J, Y);
      else if (getInternalData(J).boosted) return getDocument().body;
      else return J;
    }
    function shouldSettleAttribute(J) {
      let Y = htmx.config.attributesToSettle;
      for (let Z = 0; Z < Y.length; Z++) if (J === Y[Z]) return !0;
      return !1;
    }
    function cloneAttributes(J, Y) {
      forEach(J.attributes, function (Z) {
        if (!Y.hasAttribute(Z.name) && shouldSettleAttribute(Z.name))
          J.removeAttribute(Z.name);
      }),
        forEach(Y.attributes, function (Z) {
          if (shouldSettleAttribute(Z.name)) J.setAttribute(Z.name, Z.value);
        });
    }
    function isInlineSwap(J, Y) {
      let Z = getExtensions(Y);
      for (let $ = 0; $ < Z.length; $++) {
        let W = Z[$];
        try {
          if (W.isInlineSwap(J)) return !0;
        } catch (X) {
          logError(X);
        }
      }
      return J === "outerHTML";
    }
    function oobSwap(J, Y, Z, $) {
      $ = $ || getDocument();
      let W = "#" + getRawAttribute(Y, "id"),
        X = "outerHTML";
      if (J === "true");
      else if (J.indexOf(":") > 0)
        (X = J.substring(0, J.indexOf(":"))),
          (W = J.substring(J.indexOf(":") + 1));
      else X = J;
      Y.removeAttribute("hx-swap-oob"), Y.removeAttribute("data-hx-swap-oob");
      let Q = querySelectorAllExt($, W, !1);
      if (Q)
        forEach(Q, function (B) {
          let U,
            _ = Y.cloneNode(!0);
          if (
            ((U = getDocument().createDocumentFragment()),
            U.appendChild(_),
            !isInlineSwap(X, B))
          )
            U = asParentNode(_);
          let G = { shouldSwap: !0, target: B, fragment: U };
          if (!triggerEvent(B, "htmx:oobBeforeSwap", G)) return;
          if (((B = G.target), G.shouldSwap))
            handlePreservedElements(U),
              swapWithStyle(X, B, B, U, Z),
              restorePreservedElements();
          forEach(Z.elts, function (z) {
            triggerEvent(z, "htmx:oobAfterSwap", G);
          });
        }),
          Y.parentNode.removeChild(Y);
      else
        Y.parentNode.removeChild(Y),
          triggerErrorEvent(getDocument().body, "htmx:oobErrorNoTarget", {
            content: Y,
          });
      return J;
    }
    function restorePreservedElements() {
      let J = find("#--htmx-preserve-pantry--");
      if (J) {
        for (let Y of [...J.children]) {
          let Z = find("#" + Y.id);
          Z.parentNode.moveBefore(Y, Z), Z.remove();
        }
        J.remove();
      }
    }
    function handlePreservedElements(J) {
      forEach(findAll(J, "[hx-preserve], [data-hx-preserve]"), function (Y) {
        let Z = getAttributeValue(Y, "id"),
          $ = getDocument().getElementById(Z);
        if ($ != null)
          if (Y.moveBefore) {
            let W = find("#--htmx-preserve-pantry--");
            if (W == null)
              getDocument().body.insertAdjacentHTML(
                "afterend",
                "<div id='--htmx-preserve-pantry--'></div>",
              ),
                (W = find("#--htmx-preserve-pantry--"));
            W.moveBefore($, null);
          } else Y.parentNode.replaceChild($, Y);
      });
    }
    function handleAttributes(J, Y, Z) {
      forEach(Y.querySelectorAll("[id]"), function ($) {
        let W = getRawAttribute($, "id");
        if (W && W.length > 0) {
          let X = W.replace("'", "\\'"),
            Q = $.tagName.replace(":", "\\:"),
            B = asParentNode(J),
            U = B && B.querySelector(Q + "[id='" + X + "']");
          if (U && U !== B) {
            let _ = $.cloneNode();
            cloneAttributes($, U),
              Z.tasks.push(function () {
                cloneAttributes($, _);
              });
          }
        }
      });
    }
    function makeAjaxLoadTask(J) {
      return function () {
        removeClassFromElement(J, htmx.config.addedClass),
          processNode(asElement(J)),
          processFocus(asParentNode(J)),
          triggerEvent(J, "htmx:load");
      };
    }
    function processFocus(J) {
      let Z = asHtmlElement(
        matches(J, "[autofocus]") ? J : J.querySelector("[autofocus]"),
      );
      if (Z != null) Z.focus();
    }
    function insertNodesBefore(J, Y, Z, $) {
      handleAttributes(J, Z, $);
      while (Z.childNodes.length > 0) {
        let W = Z.firstChild;
        if (
          (addClassToElement(asElement(W), htmx.config.addedClass),
          J.insertBefore(W, Y),
          W.nodeType !== Node.TEXT_NODE && W.nodeType !== Node.COMMENT_NODE)
        )
          $.tasks.push(makeAjaxLoadTask(W));
      }
    }
    function stringHash(J, Y) {
      let Z = 0;
      while (Z < J.length) Y = ((Y << 5) - Y + J.charCodeAt(Z++)) | 0;
      return Y;
    }
    function attributeHash(J) {
      let Y = 0;
      if (J.attributes)
        for (let Z = 0; Z < J.attributes.length; Z++) {
          let $ = J.attributes[Z];
          if ($.value)
            (Y = stringHash($.name, Y)), (Y = stringHash($.value, Y));
        }
      return Y;
    }
    function deInitOnHandlers(J) {
      let Y = getInternalData(J);
      if (Y.onHandlers) {
        for (let Z = 0; Z < Y.onHandlers.length; Z++) {
          let $ = Y.onHandlers[Z];
          removeEventListenerImpl(J, $.event, $.listener);
        }
        delete Y.onHandlers;
      }
    }
    function deInitNode(J) {
      let Y = getInternalData(J);
      if (Y.timeout) clearTimeout(Y.timeout);
      if (Y.listenerInfos)
        forEach(Y.listenerInfos, function (Z) {
          if (Z.on) removeEventListenerImpl(Z.on, Z.trigger, Z.listener);
        });
      deInitOnHandlers(J),
        forEach(Object.keys(Y), function (Z) {
          if (Z !== "firstInitCompleted") delete Y[Z];
        });
    }
    function cleanUpElement(J) {
      if (
        (triggerEvent(J, "htmx:beforeCleanupElement"),
        deInitNode(J),
        J.children)
      )
        forEach(J.children, function (Y) {
          cleanUpElement(Y);
        });
    }
    function swapOuterHTML(J, Y, Z) {
      if (J instanceof Element && J.tagName === "BODY")
        return swapInnerHTML(J, Y, Z);
      let $,
        W = J.previousSibling,
        X = parentElt(J);
      if (!X) return;
      if ((insertNodesBefore(X, J, Y, Z), W == null)) $ = X.firstChild;
      else $ = W.nextSibling;
      Z.elts = Z.elts.filter(function (Q) {
        return Q !== J;
      });
      while ($ && $ !== J) {
        if ($ instanceof Element) Z.elts.push($);
        $ = $.nextSibling;
      }
      if ((cleanUpElement(J), J instanceof Element)) J.remove();
      else J.parentNode.removeChild(J);
    }
    function swapAfterBegin(J, Y, Z) {
      return insertNodesBefore(J, J.firstChild, Y, Z);
    }
    function swapBeforeBegin(J, Y, Z) {
      return insertNodesBefore(parentElt(J), J, Y, Z);
    }
    function swapBeforeEnd(J, Y, Z) {
      return insertNodesBefore(J, null, Y, Z);
    }
    function swapAfterEnd(J, Y, Z) {
      return insertNodesBefore(parentElt(J), J.nextSibling, Y, Z);
    }
    function swapDelete(J) {
      cleanUpElement(J);
      let Y = parentElt(J);
      if (Y) return Y.removeChild(J);
    }
    function swapInnerHTML(J, Y, Z) {
      let $ = J.firstChild;
      if ((insertNodesBefore(J, $, Y, Z), $)) {
        while ($.nextSibling)
          cleanUpElement($.nextSibling), J.removeChild($.nextSibling);
        cleanUpElement($), J.removeChild($);
      }
    }
    function swapWithStyle(J, Y, Z, $, W) {
      switch (J) {
        case "none":
          return;
        case "outerHTML":
          swapOuterHTML(Z, $, W);
          return;
        case "afterbegin":
          swapAfterBegin(Z, $, W);
          return;
        case "beforebegin":
          swapBeforeBegin(Z, $, W);
          return;
        case "beforeend":
          swapBeforeEnd(Z, $, W);
          return;
        case "afterend":
          swapAfterEnd(Z, $, W);
          return;
        case "delete":
          swapDelete(Z);
          return;
        default:
          var X = getExtensions(Y);
          for (let Q = 0; Q < X.length; Q++) {
            let B = X[Q];
            try {
              let U = B.handleSwap(J, Z, $, W);
              if (U) {
                if (Array.isArray(U))
                  for (let _ = 0; _ < U.length; _++) {
                    let G = U[_];
                    if (
                      G.nodeType !== Node.TEXT_NODE &&
                      G.nodeType !== Node.COMMENT_NODE
                    )
                      W.tasks.push(makeAjaxLoadTask(G));
                  }
                return;
              }
            } catch (U) {
              logError(U);
            }
          }
          if (J === "innerHTML") swapInnerHTML(Z, $, W);
          else swapWithStyle(htmx.config.defaultSwapStyle, Y, Z, $, W);
      }
    }
    function findAndSwapOobElements(J, Y, Z) {
      var $ = findAll(J, "[hx-swap-oob], [data-hx-swap-oob]");
      return (
        forEach($, function (W) {
          if (htmx.config.allowNestedOobSwaps || W.parentElement === null) {
            let X = getAttributeValue(W, "hx-swap-oob");
            if (X != null) oobSwap(X, W, Y, Z);
          } else
            W.removeAttribute("hx-swap-oob"),
              W.removeAttribute("data-hx-swap-oob");
        }),
        $.length > 0
      );
    }
    function swap(J, Y, Z, $) {
      if (!$) $ = {};
      J = resolveTarget(J);
      let W = $.contextElement
          ? getRootNode($.contextElement, !1)
          : getDocument(),
        X = document.activeElement,
        Q = {};
      try {
        Q = {
          elt: X,
          start: X ? X.selectionStart : null,
          end: X ? X.selectionEnd : null,
        };
      } catch (_) {}
      let B = makeSettleInfo(J);
      if (Z.swapStyle === "textContent") J.textContent = Y;
      else {
        let _ = makeFragment(Y);
        if (((B.title = _.title), $.selectOOB)) {
          let G = $.selectOOB.split(",");
          for (let z = 0; z < G.length; z++) {
            let K = G[z].split(":", 2),
              A = K[0].trim();
            if (A.indexOf("#") === 0) A = A.substring(1);
            let L = K[1] || "true",
              H = _.querySelector("#" + A);
            if (H) oobSwap(L, H, B, W);
          }
        }
        if (
          (findAndSwapOobElements(_, B, W),
          forEach(findAll(_, "template"), function (G) {
            if (G.content && findAndSwapOobElements(G.content, B, W))
              G.remove();
          }),
          $.select)
        ) {
          let G = getDocument().createDocumentFragment();
          forEach(_.querySelectorAll($.select), function (z) {
            G.appendChild(z);
          }),
            (_ = G);
        }
        handlePreservedElements(_),
          swapWithStyle(Z.swapStyle, $.contextElement, J, _, B),
          restorePreservedElements();
      }
      if (Q.elt && !bodyContains(Q.elt) && getRawAttribute(Q.elt, "id")) {
        let _ = document.getElementById(getRawAttribute(Q.elt, "id")),
          G = {
            preventScroll:
              Z.focusScroll !== void 0
                ? !Z.focusScroll
                : !htmx.config.defaultFocusScroll,
          };
        if (_) {
          if (Q.start && _.setSelectionRange)
            try {
              _.setSelectionRange(Q.start, Q.end);
            } catch (z) {}
          _.focus(G);
        }
      }
      if (
        (J.classList.remove(htmx.config.swappingClass),
        forEach(B.elts, function (_) {
          if (_.classList) _.classList.add(htmx.config.settlingClass);
          triggerEvent(_, "htmx:afterSwap", $.eventInfo);
        }),
        $.afterSwapCallback)
      )
        $.afterSwapCallback();
      if (!Z.ignoreTitle) handleTitle(B.title);
      let U = function () {
        if (
          (forEach(B.tasks, function (_) {
            _.call();
          }),
          forEach(B.elts, function (_) {
            if (_.classList) _.classList.remove(htmx.config.settlingClass);
            triggerEvent(_, "htmx:afterSettle", $.eventInfo);
          }),
          $.anchor)
        ) {
          let _ = asElement(resolveTarget("#" + $.anchor));
          if (_) _.scrollIntoView({ block: "start", behavior: "auto" });
        }
        if ((updateScrollState(B.elts, Z), $.afterSettleCallback))
          $.afterSettleCallback();
      };
      if (Z.settleDelay > 0) getWindow().setTimeout(U, Z.settleDelay);
      else U();
    }
    function handleTriggerHeader(J, Y, Z) {
      let $ = J.getResponseHeader(Y);
      if ($.indexOf("{") === 0) {
        let W = parseJSON($);
        for (let X in W)
          if (W.hasOwnProperty(X)) {
            let Q = W[X];
            if (isRawObject(Q)) Z = Q.target !== void 0 ? Q.target : Z;
            else Q = { value: Q };
            triggerEvent(Z, X, Q);
          }
      } else {
        let W = $.split(",");
        for (let X = 0; X < W.length; X++) triggerEvent(Z, W[X].trim(), []);
      }
    }
    let WHITESPACE = /\s/,
      WHITESPACE_OR_COMMA = /[\s,]/,
      SYMBOL_START = /[_$a-zA-Z]/,
      SYMBOL_CONT = /[_$a-zA-Z0-9]/,
      STRINGISH_START = ['"', "'", "/"],
      NOT_WHITESPACE = /[^\s]/,
      COMBINED_SELECTOR_START = /[{(]/,
      COMBINED_SELECTOR_END = /[})]/;
    function tokenizeString(J) {
      let Y = [],
        Z = 0;
      while (Z < J.length) {
        if (SYMBOL_START.exec(J.charAt(Z))) {
          var $ = Z;
          while (SYMBOL_CONT.exec(J.charAt(Z + 1))) Z++;
          Y.push(J.substring($, Z + 1));
        } else if (STRINGISH_START.indexOf(J.charAt(Z)) !== -1) {
          let W = J.charAt(Z);
          var $ = Z;
          Z++;
          while (Z < J.length && J.charAt(Z) !== W) {
            if (J.charAt(Z) === "\\") Z++;
            Z++;
          }
          Y.push(J.substring($, Z + 1));
        } else {
          let W = J.charAt(Z);
          Y.push(W);
        }
        Z++;
      }
      return Y;
    }
    function isPossibleRelativeReference(J, Y, Z) {
      return (
        SYMBOL_START.exec(J.charAt(0)) &&
        J !== "true" &&
        J !== "false" &&
        J !== "this" &&
        J !== Z &&
        Y !== "."
      );
    }
    function maybeGenerateConditional(J, Y, Z) {
      if (Y[0] === "[") {
        Y.shift();
        let $ = 1,
          W = " return (function(" + Z + "){ return (",
          X = null;
        while (Y.length > 0) {
          let Q = Y[0];
          if (Q === "]") {
            if (($--, $ === 0)) {
              if (X === null) W = W + "true";
              Y.shift(), (W += ")})");
              try {
                let B = maybeEval(
                  J,
                  function () {
                    return Function(W)();
                  },
                  function () {
                    return !0;
                  },
                );
                return (B.source = W), B;
              } catch (B) {
                return (
                  triggerErrorEvent(getDocument().body, "htmx:syntax:error", {
                    error: B,
                    source: W,
                  }),
                  null
                );
              }
            }
          } else if (Q === "[") $++;
          if (isPossibleRelativeReference(Q, X, Z))
            W +=
              "((" +
              Z +
              "." +
              Q +
              ") ? (" +
              Z +
              "." +
              Q +
              ") : (window." +
              Q +
              "))";
          else W = W + Q;
          X = Y.shift();
        }
      }
    }
    function consumeUntil(J, Y) {
      let Z = "";
      while (J.length > 0 && !Y.test(J[0])) Z += J.shift();
      return Z;
    }
    function consumeCSSSelector(J) {
      let Y;
      if (J.length > 0 && COMBINED_SELECTOR_START.test(J[0]))
        J.shift(),
          (Y = consumeUntil(J, COMBINED_SELECTOR_END).trim()),
          J.shift();
      else Y = consumeUntil(J, WHITESPACE_OR_COMMA);
      return Y;
    }
    let INPUT_SELECTOR = "input, textarea, select";
    function parseAndCacheTrigger(J, Y, Z) {
      let $ = [],
        W = tokenizeString(Y);
      do {
        consumeUntil(W, NOT_WHITESPACE);
        let B = W.length,
          U = consumeUntil(W, /[,\[\s]/);
        if (U !== "")
          if (U === "every") {
            let _ = { trigger: "every" };
            consumeUntil(W, NOT_WHITESPACE),
              (_.pollInterval = parseInterval(consumeUntil(W, /[,\[\s]/))),
              consumeUntil(W, NOT_WHITESPACE);
            var X = maybeGenerateConditional(J, W, "event");
            if (X) _.eventFilter = X;
            $.push(_);
          } else {
            let _ = { trigger: U };
            var X = maybeGenerateConditional(J, W, "event");
            if (X) _.eventFilter = X;
            consumeUntil(W, NOT_WHITESPACE);
            while (W.length > 0 && W[0] !== ",") {
              let z = W.shift();
              if (z === "changed") _.changed = !0;
              else if (z === "once") _.once = !0;
              else if (z === "consume") _.consume = !0;
              else if (z === "delay" && W[0] === ":")
                W.shift(),
                  (_.delay = parseInterval(
                    consumeUntil(W, WHITESPACE_OR_COMMA),
                  ));
              else if (z === "from" && W[0] === ":") {
                if ((W.shift(), COMBINED_SELECTOR_START.test(W[0])))
                  var Q = consumeCSSSelector(W);
                else {
                  var Q = consumeUntil(W, WHITESPACE_OR_COMMA);
                  if (
                    Q === "closest" ||
                    Q === "find" ||
                    Q === "next" ||
                    Q === "previous"
                  ) {
                    W.shift();
                    let A = consumeCSSSelector(W);
                    if (A.length > 0) Q += " " + A;
                  }
                }
                _.from = Q;
              } else if (z === "target" && W[0] === ":")
                W.shift(), (_.target = consumeCSSSelector(W));
              else if (z === "throttle" && W[0] === ":")
                W.shift(),
                  (_.throttle = parseInterval(
                    consumeUntil(W, WHITESPACE_OR_COMMA),
                  ));
              else if (z === "queue" && W[0] === ":")
                W.shift(), (_.queue = consumeUntil(W, WHITESPACE_OR_COMMA));
              else if (z === "root" && W[0] === ":")
                W.shift(), (_[z] = consumeCSSSelector(W));
              else if (z === "threshold" && W[0] === ":")
                W.shift(), (_[z] = consumeUntil(W, WHITESPACE_OR_COMMA));
              else
                triggerErrorEvent(J, "htmx:syntax:error", { token: W.shift() });
              consumeUntil(W, NOT_WHITESPACE);
            }
            $.push(_);
          }
        if (W.length === B)
          triggerErrorEvent(J, "htmx:syntax:error", { token: W.shift() });
        consumeUntil(W, NOT_WHITESPACE);
      } while (W[0] === "," && W.shift());
      if (Z) Z[Y] = $;
      return $;
    }
    function getTriggerSpecs(J) {
      let Y = getAttributeValue(J, "hx-trigger"),
        Z = [];
      if (Y) {
        let $ = htmx.config.triggerSpecsCache;
        Z = ($ && $[Y]) || parseAndCacheTrigger(J, Y, $);
      }
      if (Z.length > 0) return Z;
      else if (matches(J, "form")) return [{ trigger: "submit" }];
      else if (matches(J, 'input[type="button"], input[type="submit"]'))
        return [{ trigger: "click" }];
      else if (matches(J, INPUT_SELECTOR)) return [{ trigger: "change" }];
      else return [{ trigger: "click" }];
    }
    function cancelPolling(J) {
      getInternalData(J).cancelled = !0;
    }
    function processPolling(J, Y, Z) {
      let $ = getInternalData(J);
      $.timeout = getWindow().setTimeout(function () {
        if (bodyContains(J) && $.cancelled !== !0) {
          if (
            !maybeFilterEvent(
              Z,
              J,
              makeEvent("hx:poll:trigger", { triggerSpec: Z, target: J }),
            )
          )
            Y(J);
          processPolling(J, Y, Z);
        }
      }, Z.pollInterval);
    }
    function isLocalLink(J) {
      return (
        location.hostname === J.hostname &&
        getRawAttribute(J, "href") &&
        getRawAttribute(J, "href").indexOf("#") !== 0
      );
    }
    function eltIsDisabled(J) {
      return closest(J, htmx.config.disableSelector);
    }
    function boostElement(J, Y, Z) {
      if (
        (J instanceof HTMLAnchorElement &&
          isLocalLink(J) &&
          (J.target === "" || J.target === "_self")) ||
        (J.tagName === "FORM" &&
          String(getRawAttribute(J, "method")).toLowerCase() !== "dialog")
      ) {
        Y.boosted = !0;
        let $, W;
        if (J.tagName === "A") ($ = "get"), (W = getRawAttribute(J, "href"));
        else {
          let X = getRawAttribute(J, "method");
          if (
            (($ = X ? X.toLowerCase() : "get"),
            (W = getRawAttribute(J, "action")),
            W == null || W === "")
          )
            W = getDocument().location.href;
          if ($ === "get" && W.includes("?")) W = W.replace(/\?[^#]+/, "");
        }
        Z.forEach(function (X) {
          addEventListener(
            J,
            function (Q, B) {
              let U = asElement(Q);
              if (eltIsDisabled(U)) {
                cleanUpElement(U);
                return;
              }
              issueAjaxRequest($, W, U, B);
            },
            Y,
            X,
            !0,
          );
        });
      }
    }
    function shouldCancel(J, Y) {
      let Z = asElement(Y);
      if (!Z) return !1;
      if (J.type === "submit" || J.type === "click") {
        if (Z.tagName === "FORM") return !0;
        if (
          matches(Z, 'input[type="submit"], button') &&
          (matches(Z, "[form]") || closest(Z, "form") !== null)
        )
          return !0;
        if (
          Z instanceof HTMLAnchorElement &&
          Z.href &&
          (Z.getAttribute("href") === "#" ||
            Z.getAttribute("href").indexOf("#") !== 0)
        )
          return !0;
      }
      return !1;
    }
    function ignoreBoostedAnchorCtrlClick(J, Y) {
      return (
        getInternalData(J).boosted &&
        J instanceof HTMLAnchorElement &&
        Y.type === "click" &&
        (Y.ctrlKey || Y.metaKey)
      );
    }
    function maybeFilterEvent(J, Y, Z) {
      let $ = J.eventFilter;
      if ($)
        try {
          return $.call(Y, Z) !== !0;
        } catch (W) {
          let X = $.source;
          return (
            triggerErrorEvent(getDocument().body, "htmx:eventFilter:error", {
              error: W,
              source: X,
            }),
            !0
          );
        }
      return !1;
    }
    function addEventListener(J, Y, Z, $, W) {
      let X = getInternalData(J),
        Q;
      if ($.from) Q = querySelectorAllExt(J, $.from);
      else Q = [J];
      if ($.changed) {
        if (!("lastValue" in X)) X.lastValue = new WeakMap();
        Q.forEach(function (B) {
          if (!X.lastValue.has($)) X.lastValue.set($, new WeakMap());
          X.lastValue.get($).set(B, B.value);
        });
      }
      forEach(Q, function (B) {
        let U = function (_) {
          if (!bodyContains(J)) {
            B.removeEventListener($.trigger, U);
            return;
          }
          if (ignoreBoostedAnchorCtrlClick(J, _)) return;
          if (W || shouldCancel(_, J)) _.preventDefault();
          if (maybeFilterEvent($, J, _)) return;
          let G = getInternalData(_);
          if (((G.triggerSpec = $), G.handledFor == null)) G.handledFor = [];
          if (G.handledFor.indexOf(J) < 0) {
            if ((G.handledFor.push(J), $.consume)) _.stopPropagation();
            if ($.target && _.target) {
              if (!matches(asElement(_.target), $.target)) return;
            }
            if ($.once)
              if (X.triggeredOnce) return;
              else X.triggeredOnce = !0;
            if ($.changed) {
              let z = event.target,
                K = z.value,
                A = X.lastValue.get($);
              if (A.has(z) && A.get(z) === K) return;
              A.set(z, K);
            }
            if (X.delayed) clearTimeout(X.delayed);
            if (X.throttle) return;
            if ($.throttle > 0) {
              if (!X.throttle)
                triggerEvent(J, "htmx:trigger"),
                  Y(J, _),
                  (X.throttle = getWindow().setTimeout(function () {
                    X.throttle = null;
                  }, $.throttle));
            } else if ($.delay > 0)
              X.delayed = getWindow().setTimeout(function () {
                triggerEvent(J, "htmx:trigger"), Y(J, _);
              }, $.delay);
            else triggerEvent(J, "htmx:trigger"), Y(J, _);
          }
        };
        if (Z.listenerInfos == null) Z.listenerInfos = [];
        Z.listenerInfos.push({ trigger: $.trigger, listener: U, on: B }),
          B.addEventListener($.trigger, U);
      });
    }
    let windowIsScrolling = !1,
      scrollHandler = null;
    function initScrollHandler() {
      if (!scrollHandler)
        (scrollHandler = function () {
          windowIsScrolling = !0;
        }),
          window.addEventListener("scroll", scrollHandler),
          window.addEventListener("resize", scrollHandler),
          setInterval(function () {
            if (windowIsScrolling)
              (windowIsScrolling = !1),
                forEach(
                  getDocument().querySelectorAll(
                    "[hx-trigger*='revealed'],[data-hx-trigger*='revealed']",
                  ),
                  function (J) {
                    maybeReveal(J);
                  },
                );
          }, 200);
    }
    function maybeReveal(J) {
      if (!hasAttribute(J, "data-hx-revealed") && isScrolledIntoView(J))
        if (
          (J.setAttribute("data-hx-revealed", "true"),
          getInternalData(J).initHash)
        )
          triggerEvent(J, "revealed");
        else
          J.addEventListener(
            "htmx:afterProcessNode",
            function () {
              triggerEvent(J, "revealed");
            },
            { once: !0 },
          );
    }
    function loadImmediately(J, Y, Z, $) {
      let W = function () {
        if (!Z.loaded) (Z.loaded = !0), triggerEvent(J, "htmx:trigger"), Y(J);
      };
      if ($ > 0) getWindow().setTimeout(W, $);
      else W();
    }
    function processVerbs(J, Y, Z) {
      let $ = !1;
      return (
        forEach(VERBS, function (W) {
          if (hasAttribute(J, "hx-" + W)) {
            let X = getAttributeValue(J, "hx-" + W);
            ($ = !0),
              (Y.path = X),
              (Y.verb = W),
              Z.forEach(function (Q) {
                addTriggerHandler(J, Q, Y, function (B, U) {
                  let _ = asElement(B);
                  if (closest(_, htmx.config.disableSelector)) {
                    cleanUpElement(_);
                    return;
                  }
                  issueAjaxRequest(W, X, _, U);
                });
              });
          }
        }),
        $
      );
    }
    function addTriggerHandler(J, Y, Z, $) {
      if (Y.trigger === "revealed")
        initScrollHandler(),
          addEventListener(J, $, Z, Y),
          maybeReveal(asElement(J));
      else if (Y.trigger === "intersect") {
        let W = {};
        if (Y.root) W.root = querySelectorExt(J, Y.root);
        if (Y.threshold) W.threshold = parseFloat(Y.threshold);
        new IntersectionObserver(function (Q) {
          for (let B = 0; B < Q.length; B++)
            if (Q[B].isIntersecting) {
              triggerEvent(J, "intersect");
              break;
            }
        }, W).observe(asElement(J)),
          addEventListener(asElement(J), $, Z, Y);
      } else if (!Z.firstInitCompleted && Y.trigger === "load") {
        if (!maybeFilterEvent(Y, J, makeEvent("load", { elt: J })))
          loadImmediately(asElement(J), $, Z, Y.delay);
      } else if (Y.pollInterval > 0)
        (Z.polling = !0), processPolling(asElement(J), $, Y);
      else addEventListener(J, $, Z, Y);
    }
    function shouldProcessHxOn(J) {
      let Y = asElement(J);
      if (!Y) return !1;
      let Z = Y.attributes;
      for (let $ = 0; $ < Z.length; $++) {
        let W = Z[$].name;
        if (
          startsWith(W, "hx-on:") ||
          startsWith(W, "data-hx-on:") ||
          startsWith(W, "hx-on-") ||
          startsWith(W, "data-hx-on-")
        )
          return !0;
      }
      return !1;
    }
    let HX_ON_QUERY = new XPathEvaluator().createExpression(
      './/*[@*[ starts-with(name(), "hx-on:") or starts-with(name(), "data-hx-on:") or starts-with(name(), "hx-on-") or starts-with(name(), "data-hx-on-") ]]',
    );
    function processHXOnRoot(J, Y) {
      if (shouldProcessHxOn(J)) Y.push(asElement(J));
      let Z = HX_ON_QUERY.evaluate(J),
        $ = null;
      while (($ = Z.iterateNext())) Y.push(asElement($));
    }
    function findHxOnWildcardElements(J) {
      let Y = [];
      if (J instanceof DocumentFragment)
        for (let Z of J.childNodes) processHXOnRoot(Z, Y);
      else processHXOnRoot(J, Y);
      return Y;
    }
    function findElementsToProcess(J) {
      if (J.querySelectorAll) {
        let $ = [];
        for (let X in extensions) {
          let Q = extensions[X];
          if (Q.getSelectors) {
            var Y = Q.getSelectors();
            if (Y) $.push(Y);
          }
        }
        return J.querySelectorAll(
          VERB_SELECTOR +
            ", [hx-boost] a, [data-hx-boost] a, a[hx-boost], a[data-hx-boost], form, [type='submit'], [hx-ext], [data-hx-ext], [hx-trigger], [data-hx-trigger]" +
            $.flat()
              .map((X) => ", " + X)
              .join(""),
        );
      } else return [];
    }
    function maybeSetLastButtonClicked(J) {
      let Y = closest(asElement(J.target), "button, input[type='submit']"),
        Z = getRelatedFormData(J);
      if (Z) Z.lastButtonClicked = Y;
    }
    function maybeUnsetLastButtonClicked(J) {
      let Y = getRelatedFormData(J);
      if (Y) Y.lastButtonClicked = null;
    }
    function getRelatedFormData(J) {
      let Y = closest(asElement(J.target), "button, input[type='submit']");
      if (!Y) return;
      let Z =
        resolveTarget("#" + getRawAttribute(Y, "form"), Y.getRootNode()) ||
        closest(Y, "form");
      if (!Z) return;
      return getInternalData(Z);
    }
    function initButtonTracking(J) {
      J.addEventListener("click", maybeSetLastButtonClicked),
        J.addEventListener("focusin", maybeSetLastButtonClicked),
        J.addEventListener("focusout", maybeUnsetLastButtonClicked);
    }
    function addHxOnEventHandler(J, Y, Z) {
      let $ = getInternalData(J);
      if (!Array.isArray($.onHandlers)) $.onHandlers = [];
      let W,
        X = function (Q) {
          maybeEval(J, function () {
            if (eltIsDisabled(J)) return;
            if (!W) W = new Function("event", Z);
            W.call(J, Q);
          });
        };
      J.addEventListener(Y, X), $.onHandlers.push({ event: Y, listener: X });
    }
    function processHxOnWildcard(J) {
      deInitOnHandlers(J);
      for (let Y = 0; Y < J.attributes.length; Y++) {
        let Z = J.attributes[Y].name,
          $ = J.attributes[Y].value;
        if (startsWith(Z, "hx-on") || startsWith(Z, "data-hx-on")) {
          let W = Z.indexOf("-on") + 3,
            X = Z.slice(W, W + 1);
          if (X === "-" || X === ":") {
            let Q = Z.slice(W + 1);
            if (startsWith(Q, ":")) Q = "htmx" + Q;
            else if (startsWith(Q, "-")) Q = "htmx:" + Q.slice(1);
            else if (startsWith(Q, "htmx-")) Q = "htmx:" + Q.slice(5);
            addHxOnEventHandler(J, Q, $);
          }
        }
      }
    }
    function initNode(J) {
      if (closest(J, htmx.config.disableSelector)) {
        cleanUpElement(J);
        return;
      }
      let Y = getInternalData(J),
        Z = attributeHash(J);
      if (Y.initHash !== Z) {
        deInitNode(J),
          (Y.initHash = Z),
          triggerEvent(J, "htmx:beforeProcessNode");
        let $ = getTriggerSpecs(J);
        if (!processVerbs(J, Y, $)) {
          if (getClosestAttributeValue(J, "hx-boost") === "true")
            boostElement(J, Y, $);
          else if (hasAttribute(J, "hx-trigger"))
            $.forEach(function (X) {
              addTriggerHandler(J, X, Y, function () {});
            });
        }
        if (
          J.tagName === "FORM" ||
          (getRawAttribute(J, "type") === "submit" && hasAttribute(J, "form"))
        )
          initButtonTracking(J);
        (Y.firstInitCompleted = !0), triggerEvent(J, "htmx:afterProcessNode");
      }
    }
    function processNode(J) {
      if (((J = resolveTarget(J)), closest(J, htmx.config.disableSelector))) {
        cleanUpElement(J);
        return;
      }
      initNode(J),
        forEach(findElementsToProcess(J), function (Y) {
          initNode(Y);
        }),
        forEach(findHxOnWildcardElements(J), processHxOnWildcard);
    }
    function kebabEventName(J) {
      return J.replace(/([a-z0-9])([A-Z])/g, "$1-$2").toLowerCase();
    }
    function makeEvent(J, Y) {
      let Z;
      if (window.CustomEvent && typeof window.CustomEvent === "function")
        Z = new CustomEvent(J, {
          bubbles: !0,
          cancelable: !0,
          composed: !0,
          detail: Y,
        });
      else
        (Z = getDocument().createEvent("CustomEvent")),
          Z.initCustomEvent(J, !0, !0, Y);
      return Z;
    }
    function triggerErrorEvent(J, Y, Z) {
      triggerEvent(J, Y, mergeObjects({ error: Y }, Z));
    }
    function ignoreEventForLogging(J) {
      return J === "htmx:afterProcessNode";
    }
    function withExtensions(J, Y) {
      forEach(getExtensions(J), function (Z) {
        try {
          Y(Z);
        } catch ($) {
          logError($);
        }
      });
    }
    function logError(J) {
      if (console.error) console.error(J);
      else if (console.log) console.log("ERROR: ", J);
    }
    function triggerEvent(J, Y, Z) {
      if (((J = resolveTarget(J)), Z == null)) Z = {};
      Z.elt = J;
      let $ = makeEvent(Y, Z);
      if (htmx.logger && !ignoreEventForLogging(Y)) htmx.logger(J, Y, Z);
      if (Z.error)
        logError(Z.error), triggerEvent(J, "htmx:error", { errorInfo: Z });
      let W = J.dispatchEvent($),
        X = kebabEventName(Y);
      if (W && X !== Y) {
        let Q = makeEvent(X, $.detail);
        W = W && J.dispatchEvent(Q);
      }
      return (
        withExtensions(asElement(J), function (Q) {
          W = W && Q.onEvent(Y, $) !== !1 && !$.defaultPrevented;
        }),
        W
      );
    }
    let currentPathForHistory = location.pathname + location.search;
    function getHistoryElement() {
      return (
        getDocument().querySelector("[hx-history-elt],[data-hx-history-elt]") ||
        getDocument().body
      );
    }
    function saveToHistoryCache(J, Y) {
      if (!canAccessLocalStorage()) return;
      let Z = cleanInnerHtmlForHistory(Y),
        $ = getDocument().title,
        W = window.scrollY;
      if (htmx.config.historyCacheSize <= 0) {
        localStorage.removeItem("htmx-history-cache");
        return;
      }
      J = normalizePath(J);
      let X = parseJSON(localStorage.getItem("htmx-history-cache")) || [];
      for (let B = 0; B < X.length; B++)
        if (X[B].url === J) {
          X.splice(B, 1);
          break;
        }
      let Q = { url: J, content: Z, title: $, scroll: W };
      triggerEvent(getDocument().body, "htmx:historyItemCreated", {
        item: Q,
        cache: X,
      }),
        X.push(Q);
      while (X.length > htmx.config.historyCacheSize) X.shift();
      while (X.length > 0)
        try {
          localStorage.setItem("htmx-history-cache", JSON.stringify(X));
          break;
        } catch (B) {
          triggerErrorEvent(getDocument().body, "htmx:historyCacheError", {
            cause: B,
            cache: X,
          }),
            X.shift();
        }
    }
    function getCachedHistory(J) {
      if (!canAccessLocalStorage()) return null;
      J = normalizePath(J);
      let Y = parseJSON(localStorage.getItem("htmx-history-cache")) || [];
      for (let Z = 0; Z < Y.length; Z++) if (Y[Z].url === J) return Y[Z];
      return null;
    }
    function cleanInnerHtmlForHistory(J) {
      let Y = htmx.config.requestClass,
        Z = J.cloneNode(!0);
      return (
        forEach(findAll(Z, "." + Y), function ($) {
          removeClassFromElement($, Y);
        }),
        forEach(findAll(Z, "[data-disabled-by-htmx]"), function ($) {
          $.removeAttribute("disabled");
        }),
        Z.innerHTML
      );
    }
    function saveCurrentPageToHistory() {
      let J = getHistoryElement(),
        Y = currentPathForHistory || location.pathname + location.search,
        Z;
      try {
        Z = getDocument().querySelector(
          '[hx-history="false" i],[data-hx-history="false" i]',
        );
      } catch ($) {
        Z = getDocument().querySelector(
          '[hx-history="false"],[data-hx-history="false"]',
        );
      }
      if (!Z)
        triggerEvent(getDocument().body, "htmx:beforeHistorySave", {
          path: Y,
          historyElt: J,
        }),
          saveToHistoryCache(Y, J);
      if (htmx.config.historyEnabled)
        history.replaceState(
          { htmx: !0 },
          getDocument().title,
          window.location.href,
        );
    }
    function pushUrlIntoHistory(J) {
      if (htmx.config.getCacheBusterParam) {
        if (
          ((J = J.replace(/org\.htmx\.cache-buster=[^&]*&?/, "")),
          endsWith(J, "&") || endsWith(J, "?"))
        )
          J = J.slice(0, -1);
      }
      if (htmx.config.historyEnabled) history.pushState({ htmx: !0 }, "", J);
      currentPathForHistory = J;
    }
    function replaceUrlInHistory(J) {
      if (htmx.config.historyEnabled) history.replaceState({ htmx: !0 }, "", J);
      currentPathForHistory = J;
    }
    function settleImmediately(J) {
      forEach(J, function (Y) {
        Y.call(void 0);
      });
    }
    function loadHistoryFromServer(J) {
      let Y = new XMLHttpRequest(),
        Z = { path: J, xhr: Y };
      triggerEvent(getDocument().body, "htmx:historyCacheMiss", Z),
        Y.open("GET", J, !0),
        Y.setRequestHeader("HX-Request", "true"),
        Y.setRequestHeader("HX-History-Restore-Request", "true"),
        Y.setRequestHeader("HX-Current-URL", getDocument().location.href),
        (Y.onload = function () {
          if (this.status >= 200 && this.status < 400) {
            triggerEvent(getDocument().body, "htmx:historyCacheMissLoad", Z);
            let $ = makeFragment(this.response),
              W =
                $.querySelector("[hx-history-elt],[data-hx-history-elt]") || $,
              X = getHistoryElement(),
              Q = makeSettleInfo(X);
            handleTitle($.title),
              handlePreservedElements($),
              swapInnerHTML(X, W, Q),
              restorePreservedElements(),
              settleImmediately(Q.tasks),
              (currentPathForHistory = J),
              triggerEvent(getDocument().body, "htmx:historyRestore", {
                path: J,
                cacheMiss: !0,
                serverResponse: this.response,
              });
          } else
            triggerErrorEvent(
              getDocument().body,
              "htmx:historyCacheMissLoadError",
              Z,
            );
        }),
        Y.send();
    }
    function restoreHistory(J) {
      saveCurrentPageToHistory(),
        (J = J || location.pathname + location.search);
      let Y = getCachedHistory(J);
      if (Y) {
        let Z = makeFragment(Y.content),
          $ = getHistoryElement(),
          W = makeSettleInfo($);
        handleTitle(Y.title),
          handlePreservedElements(Z),
          swapInnerHTML($, Z, W),
          restorePreservedElements(),
          settleImmediately(W.tasks),
          getWindow().setTimeout(function () {
            window.scrollTo(0, Y.scroll);
          }, 0),
          (currentPathForHistory = J),
          triggerEvent(getDocument().body, "htmx:historyRestore", {
            path: J,
            item: Y,
          });
      } else if (htmx.config.refreshOnHistoryMiss) window.location.reload(!0);
      else loadHistoryFromServer(J);
    }
    function addRequestIndicatorClasses(J) {
      let Y = findAttributeTargets(J, "hx-indicator");
      if (Y == null) Y = [J];
      return (
        forEach(Y, function (Z) {
          let $ = getInternalData(Z);
          ($.requestCount = ($.requestCount || 0) + 1),
            Z.classList.add.call(Z.classList, htmx.config.requestClass);
        }),
        Y
      );
    }
    function disableElements(J) {
      let Y = findAttributeTargets(J, "hx-disabled-elt");
      if (Y == null) Y = [];
      return (
        forEach(Y, function (Z) {
          let $ = getInternalData(Z);
          ($.requestCount = ($.requestCount || 0) + 1),
            Z.setAttribute("disabled", ""),
            Z.setAttribute("data-disabled-by-htmx", "");
        }),
        Y
      );
    }
    function removeRequestIndicators(J, Y) {
      forEach(J.concat(Y), function (Z) {
        let $ = getInternalData(Z);
        $.requestCount = ($.requestCount || 1) - 1;
      }),
        forEach(J, function (Z) {
          if (getInternalData(Z).requestCount === 0)
            Z.classList.remove.call(Z.classList, htmx.config.requestClass);
        }),
        forEach(Y, function (Z) {
          if (getInternalData(Z).requestCount === 0)
            Z.removeAttribute("disabled"),
              Z.removeAttribute("data-disabled-by-htmx");
        });
    }
    function haveSeenNode(J, Y) {
      for (let Z = 0; Z < J.length; Z++) if (J[Z].isSameNode(Y)) return !0;
      return !1;
    }
    function shouldInclude(J) {
      let Y = J;
      if (
        Y.name === "" ||
        Y.name == null ||
        Y.disabled ||
        closest(Y, "fieldset[disabled]")
      )
        return !1;
      if (
        Y.type === "button" ||
        Y.type === "submit" ||
        Y.tagName === "image" ||
        Y.tagName === "reset" ||
        Y.tagName === "file"
      )
        return !1;
      if (Y.type === "checkbox" || Y.type === "radio") return Y.checked;
      return !0;
    }
    function addValueToFormData(J, Y, Z) {
      if (J != null && Y != null)
        if (Array.isArray(Y))
          Y.forEach(function ($) {
            Z.append(J, $);
          });
        else Z.append(J, Y);
    }
    function removeValueFromFormData(J, Y, Z) {
      if (J != null && Y != null) {
        let $ = Z.getAll(J);
        if (Array.isArray(Y)) $ = $.filter((W) => Y.indexOf(W) < 0);
        else $ = $.filter((W) => W !== Y);
        Z.delete(J), forEach($, (W) => Z.append(J, W));
      }
    }
    function processInputValue(J, Y, Z, $, W) {
      if ($ == null || haveSeenNode(J, $)) return;
      else J.push($);
      if (shouldInclude($)) {
        let X = getRawAttribute($, "name"),
          Q = $.value;
        if ($ instanceof HTMLSelectElement && $.multiple)
          Q = toArray($.querySelectorAll("option:checked")).map(function (B) {
            return B.value;
          });
        if ($ instanceof HTMLInputElement && $.files) Q = toArray($.files);
        if ((addValueToFormData(X, Q, Y), W)) validateElement($, Z);
      }
      if ($ instanceof HTMLFormElement)
        forEach($.elements, function (X) {
          if (J.indexOf(X) >= 0) removeValueFromFormData(X.name, X.value, Y);
          else J.push(X);
          if (W) validateElement(X, Z);
        }),
          new FormData($).forEach(function (X, Q) {
            if (X instanceof File && X.name === "") return;
            addValueToFormData(Q, X, Y);
          });
    }
    function validateElement(J, Y) {
      let Z = J;
      if (Z.willValidate) {
        if ((triggerEvent(Z, "htmx:validation:validate"), !Z.checkValidity()))
          Y.push({
            elt: Z,
            message: Z.validationMessage,
            validity: Z.validity,
          }),
            triggerEvent(Z, "htmx:validation:failed", {
              message: Z.validationMessage,
              validity: Z.validity,
            });
      }
    }
    function overrideFormData(J, Y) {
      for (let Z of Y.keys()) J.delete(Z);
      return (
        Y.forEach(function (Z, $) {
          J.append($, Z);
        }),
        J
      );
    }
    function getInputValues(J, Y) {
      let Z = [],
        $ = new FormData(),
        W = new FormData(),
        X = [],
        Q = getInternalData(J);
      if (Q.lastButtonClicked && !bodyContains(Q.lastButtonClicked))
        Q.lastButtonClicked = null;
      let B =
        (J instanceof HTMLFormElement && J.noValidate !== !0) ||
        getAttributeValue(J, "hx-validate") === "true";
      if (Q.lastButtonClicked)
        B = B && Q.lastButtonClicked.formNoValidate !== !0;
      if (Y !== "get") processInputValue(Z, W, X, closest(J, "form"), B);
      if (
        (processInputValue(Z, $, X, J, B),
        Q.lastButtonClicked ||
          J.tagName === "BUTTON" ||
          (J.tagName === "INPUT" && getRawAttribute(J, "type") === "submit"))
      ) {
        let _ = Q.lastButtonClicked || J,
          G = getRawAttribute(_, "name");
        addValueToFormData(G, _.value, W);
      }
      let U = findAttributeTargets(J, "hx-include");
      return (
        forEach(U, function (_) {
          if (
            (processInputValue(Z, $, X, asElement(_), B), !matches(_, "form"))
          )
            forEach(
              asParentNode(_).querySelectorAll(INPUT_SELECTOR),
              function (G) {
                processInputValue(Z, $, X, G, B);
              },
            );
        }),
        overrideFormData($, W),
        { errors: X, formData: $, values: formDataProxy($) }
      );
    }
    function appendParam(J, Y, Z) {
      if (J !== "") J += "&";
      if (String(Z) === "[object Object]") Z = JSON.stringify(Z);
      let $ = encodeURIComponent(Z);
      return (J += encodeURIComponent(Y) + "=" + $), J;
    }
    function urlEncode(J) {
      J = formDataFromObject(J);
      let Y = "";
      return (
        J.forEach(function (Z, $) {
          Y = appendParam(Y, $, Z);
        }),
        Y
      );
    }
    function getHeaders(J, Y, Z) {
      let $ = {
        "HX-Request": "true",
        "HX-Trigger": getRawAttribute(J, "id"),
        "HX-Trigger-Name": getRawAttribute(J, "name"),
        "HX-Target": getAttributeValue(Y, "id"),
        "HX-Current-URL": getDocument().location.href,
      };
      if ((getValuesForElement(J, "hx-headers", !1, $), Z !== void 0))
        $["HX-Prompt"] = Z;
      if (getInternalData(J).boosted) $["HX-Boosted"] = "true";
      return $;
    }
    function filterValues(J, Y) {
      let Z = getClosestAttributeValue(Y, "hx-params");
      if (Z)
        if (Z === "none") return new FormData();
        else if (Z === "*") return J;
        else if (Z.indexOf("not ") === 0)
          return (
            forEach(Z.slice(4).split(","), function ($) {
              ($ = $.trim()), J.delete($);
            }),
            J
          );
        else {
          let $ = new FormData();
          return (
            forEach(Z.split(","), function (W) {
              if (((W = W.trim()), J.has(W)))
                J.getAll(W).forEach(function (X) {
                  $.append(W, X);
                });
            }),
            $
          );
        }
      else return J;
    }
    function isAnchorLink(J) {
      return (
        !!getRawAttribute(J, "href") &&
        getRawAttribute(J, "href").indexOf("#") >= 0
      );
    }
    function getSwapSpecification(J, Y) {
      let Z = Y || getClosestAttributeValue(J, "hx-swap"),
        $ = {
          swapStyle: getInternalData(J).boosted
            ? "innerHTML"
            : htmx.config.defaultSwapStyle,
          swapDelay: htmx.config.defaultSwapDelay,
          settleDelay: htmx.config.defaultSettleDelay,
        };
      if (
        htmx.config.scrollIntoViewOnBoost &&
        getInternalData(J).boosted &&
        !isAnchorLink(J)
      )
        $.show = "top";
      if (Z) {
        let Q = splitOnWhitespace(Z);
        if (Q.length > 0)
          for (let B = 0; B < Q.length; B++) {
            let U = Q[B];
            if (U.indexOf("swap:") === 0)
              $.swapDelay = parseInterval(U.slice(5));
            else if (U.indexOf("settle:") === 0)
              $.settleDelay = parseInterval(U.slice(7));
            else if (U.indexOf("transition:") === 0)
              $.transition = U.slice(11) === "true";
            else if (U.indexOf("ignoreTitle:") === 0)
              $.ignoreTitle = U.slice(12) === "true";
            else if (U.indexOf("scroll:") === 0) {
              var W = U.slice(7).split(":");
              let G = W.pop();
              var X = W.length > 0 ? W.join(":") : null;
              ($.scroll = G), ($.scrollTarget = X);
            } else if (U.indexOf("show:") === 0) {
              var W = U.slice(5).split(":");
              let z = W.pop();
              var X = W.length > 0 ? W.join(":") : null;
              ($.show = z), ($.showTarget = X);
            } else if (U.indexOf("focus-scroll:") === 0) {
              let _ = U.slice(13);
              $.focusScroll = _ == "true";
            } else if (B == 0) $.swapStyle = U;
            else logError("Unknown modifier in hx-swap: " + U);
          }
      }
      return $;
    }
    function usesFormData(J) {
      return (
        getClosestAttributeValue(J, "hx-encoding") === "multipart/form-data" ||
        (matches(J, "form") &&
          getRawAttribute(J, "enctype") === "multipart/form-data")
      );
    }
    function encodeParamsForBody(J, Y, Z) {
      let $ = null;
      if (
        (withExtensions(Y, function (W) {
          if ($ == null) $ = W.encodeParameters(J, Z, Y);
        }),
        $ != null)
      )
        return $;
      else if (usesFormData(Y))
        return overrideFormData(new FormData(), formDataFromObject(Z));
      else return urlEncode(Z);
    }
    function makeSettleInfo(J) {
      return { tasks: [], elts: [J] };
    }
    function updateScrollState(J, Y) {
      let Z = J[0],
        $ = J[J.length - 1];
      if (Y.scroll) {
        var W = null;
        if (Y.scrollTarget) W = asElement(querySelectorExt(Z, Y.scrollTarget));
        if (Y.scroll === "top" && (Z || W)) (W = W || Z), (W.scrollTop = 0);
        if (Y.scroll === "bottom" && ($ || W))
          (W = W || $), (W.scrollTop = W.scrollHeight);
      }
      if (Y.show) {
        var W = null;
        if (Y.showTarget) {
          let Q = Y.showTarget;
          if (Y.showTarget === "window") Q = "body";
          W = asElement(querySelectorExt(Z, Q));
        }
        if (Y.show === "top" && (Z || W))
          (W = W || Z),
            W.scrollIntoView({
              block: "start",
              behavior: htmx.config.scrollBehavior,
            });
        if (Y.show === "bottom" && ($ || W))
          (W = W || $),
            W.scrollIntoView({
              block: "end",
              behavior: htmx.config.scrollBehavior,
            });
      }
    }
    function getValuesForElement(J, Y, Z, $) {
      if ($ == null) $ = {};
      if (J == null) return $;
      let W = getAttributeValue(J, Y);
      if (W) {
        let X = W.trim(),
          Q = Z;
        if (X === "unset") return null;
        if (X.indexOf("javascript:") === 0) (X = X.slice(11)), (Q = !0);
        else if (X.indexOf("js:") === 0) (X = X.slice(3)), (Q = !0);
        if (X.indexOf("{") !== 0) X = "{" + X + "}";
        let B;
        if (Q)
          B = maybeEval(
            J,
            function () {
              return Function("return (" + X + ")")();
            },
            {},
          );
        else B = parseJSON(X);
        for (let U in B)
          if (B.hasOwnProperty(U)) {
            if ($[U] == null) $[U] = B[U];
          }
      }
      return getValuesForElement(asElement(parentElt(J)), Y, Z, $);
    }
    function maybeEval(J, Y, Z) {
      if (htmx.config.allowEval) return Y();
      else return triggerErrorEvent(J, "htmx:evalDisallowedError"), Z;
    }
    function getHXVarsForElement(J, Y) {
      return getValuesForElement(J, "hx-vars", !0, Y);
    }
    function getHXValsForElement(J, Y) {
      return getValuesForElement(J, "hx-vals", !1, Y);
    }
    function getExpressionVars(J) {
      return mergeObjects(getHXVarsForElement(J), getHXValsForElement(J));
    }
    function safelySetHeaderValue(J, Y, Z) {
      if (Z !== null)
        try {
          J.setRequestHeader(Y, Z);
        } catch ($) {
          J.setRequestHeader(Y, encodeURIComponent(Z)),
            J.setRequestHeader(Y + "-URI-AutoEncoded", "true");
        }
    }
    function getPathFromResponse(J) {
      if (J.responseURL && typeof URL !== "undefined")
        try {
          let Y = new URL(J.responseURL);
          return Y.pathname + Y.search;
        } catch (Y) {
          triggerErrorEvent(getDocument().body, "htmx:badResponseUrl", {
            url: J.responseURL,
          });
        }
    }
    function hasHeader(J, Y) {
      return Y.test(J.getAllResponseHeaders());
    }
    function ajaxHelper(J, Y, Z) {
      if (((J = J.toLowerCase()), Z))
        if (Z instanceof Element || typeof Z === "string")
          return issueAjaxRequest(J, Y, null, null, {
            targetOverride: resolveTarget(Z) || DUMMY_ELT,
            returnPromise: !0,
          });
        else {
          let $ = resolveTarget(Z.target);
          if ((Z.target && !$) || (Z.source && !$ && !resolveTarget(Z.source)))
            $ = DUMMY_ELT;
          return issueAjaxRequest(J, Y, resolveTarget(Z.source), Z.event, {
            handler: Z.handler,
            headers: Z.headers,
            values: Z.values,
            targetOverride: $,
            swapOverride: Z.swap,
            select: Z.select,
            returnPromise: !0,
          });
        }
      else return issueAjaxRequest(J, Y, null, null, { returnPromise: !0 });
    }
    function hierarchyForElt(J) {
      let Y = [];
      while (J) Y.push(J), (J = J.parentElement);
      return Y;
    }
    function verifyPath(J, Y, Z) {
      let $, W;
      if (typeof URL === "function")
        (W = new URL(Y, document.location.href)),
          ($ = document.location.origin === W.origin);
      else (W = Y), ($ = startsWith(Y, document.location.origin));
      if (htmx.config.selfRequestsOnly) {
        if (!$) return !1;
      }
      return triggerEvent(
        J,
        "htmx:validateUrl",
        mergeObjects({ url: W, sameHost: $ }, Z),
      );
    }
    function formDataFromObject(J) {
      if (J instanceof FormData) return J;
      let Y = new FormData();
      for (let Z in J)
        if (J.hasOwnProperty(Z))
          if (J[Z] && typeof J[Z].forEach === "function")
            J[Z].forEach(function ($) {
              Y.append(Z, $);
            });
          else if (typeof J[Z] === "object" && !(J[Z] instanceof Blob))
            Y.append(Z, JSON.stringify(J[Z]));
          else Y.append(Z, J[Z]);
      return Y;
    }
    function formDataArrayProxy(J, Y, Z) {
      return new Proxy(Z, {
        get: function ($, W) {
          if (typeof W === "number") return $[W];
          if (W === "length") return $.length;
          if (W === "push")
            return function (X) {
              $.push(X), J.append(Y, X);
            };
          if (typeof $[W] === "function")
            return function () {
              $[W].apply($, arguments),
                J.delete(Y),
                $.forEach(function (X) {
                  J.append(Y, X);
                });
            };
          if ($[W] && $[W].length === 1) return $[W][0];
          else return $[W];
        },
        set: function ($, W, X) {
          return (
            ($[W] = X),
            J.delete(Y),
            $.forEach(function (Q) {
              J.append(Y, Q);
            }),
            !0
          );
        },
      });
    }
    function formDataProxy(J) {
      return new Proxy(J, {
        get: function (Y, Z) {
          if (typeof Z === "symbol") {
            let W = Reflect.get(Y, Z);
            if (typeof W === "function")
              return function () {
                return W.apply(J, arguments);
              };
            else return W;
          }
          if (Z === "toJSON") return () => Object.fromEntries(J);
          if (Z in Y)
            if (typeof Y[Z] === "function")
              return function () {
                return J[Z].apply(J, arguments);
              };
            else return Y[Z];
          let $ = J.getAll(Z);
          if ($.length === 0) return;
          else if ($.length === 1) return $[0];
          else return formDataArrayProxy(Y, Z, $);
        },
        set: function (Y, Z, $) {
          if (typeof Z !== "string") return !1;
          if ((Y.delete(Z), $ && typeof $.forEach === "function"))
            $.forEach(function (W) {
              Y.append(Z, W);
            });
          else if (typeof $ === "object" && !($ instanceof Blob))
            Y.append(Z, JSON.stringify($));
          else Y.append(Z, $);
          return !0;
        },
        deleteProperty: function (Y, Z) {
          if (typeof Z === "string") Y.delete(Z);
          return !0;
        },
        ownKeys: function (Y) {
          return Reflect.ownKeys(Object.fromEntries(Y));
        },
        getOwnPropertyDescriptor: function (Y, Z) {
          return Reflect.getOwnPropertyDescriptor(Object.fromEntries(Y), Z);
        },
      });
    }
    function issueAjaxRequest(J, Y, Z, $, W, X) {
      let Q = null,
        B = null;
      if (
        ((W = W != null ? W : {}),
        W.returnPromise && typeof Promise !== "undefined")
      )
        var U = new Promise(function (F, V) {
          (Q = F), (B = V);
        });
      if (Z == null) Z = getDocument().body;
      let _ = W.handler || handleAjaxResponse,
        G = W.select || null;
      if (!bodyContains(Z)) return maybeCall(Q), U;
      let z = W.targetOverride || asElement(getTarget(Z));
      if (z == null || z == DUMMY_ELT)
        return (
          triggerErrorEvent(Z, "htmx:targetError", {
            target: getAttributeValue(Z, "hx-target"),
          }),
          maybeCall(B),
          U
        );
      let K = getInternalData(Z),
        A = K.lastButtonClicked;
      if (A) {
        let F = getRawAttribute(A, "formaction");
        if (F != null) Y = F;
        let V = getRawAttribute(A, "formmethod");
        if (V != null) {
          if (V.toLowerCase() !== "dialog") J = V;
        }
      }
      let L = getClosestAttributeValue(Z, "hx-confirm");
      if (X === void 0) {
        if (
          triggerEvent(Z, "htmx:confirm", {
            target: z,
            elt: Z,
            path: Y,
            verb: J,
            triggeringEvent: $,
            etc: W,
            issueRequest: function (v) {
              return issueAjaxRequest(J, Y, Z, $, W, !!v);
            },
            question: L,
          }) === !1
        )
          return maybeCall(Q), U;
      }
      let H = Z,
        M = getClosestAttributeValue(Z, "hx-sync"),
        P = null,
        C = !1;
      if (M) {
        let F = M.split(":"),
          V = F[0].trim();
        if (V === "this") H = findThisElement(Z, "hx-sync");
        else H = asElement(querySelectorExt(Z, V));
        if (
          ((M = (F[1] || "drop").trim()),
          (K = getInternalData(H)),
          M === "drop" && K.xhr && K.abortable !== !0)
        )
          return maybeCall(Q), U;
        else if (M === "abort")
          if (K.xhr) return maybeCall(Q), U;
          else C = !0;
        else if (M === "replace") triggerEvent(H, "htmx:abort");
        else if (M.indexOf("queue") === 0)
          P = (M.split(" ")[1] || "last").trim();
      }
      if (K.xhr)
        if (K.abortable) triggerEvent(H, "htmx:abort");
        else {
          if (P == null) {
            if ($) {
              let F = getInternalData($);
              if (F && F.triggerSpec && F.triggerSpec.queue)
                P = F.triggerSpec.queue;
            }
            if (P == null) P = "last";
          }
          if (K.queuedRequests == null) K.queuedRequests = [];
          if (P === "first" && K.queuedRequests.length === 0)
            K.queuedRequests.push(function () {
              issueAjaxRequest(J, Y, Z, $, W);
            });
          else if (P === "all")
            K.queuedRequests.push(function () {
              issueAjaxRequest(J, Y, Z, $, W);
            });
          else if (P === "last")
            (K.queuedRequests = []),
              K.queuedRequests.push(function () {
                issueAjaxRequest(J, Y, Z, $, W);
              });
          return maybeCall(Q), U;
        }
      let q = new XMLHttpRequest();
      (K.xhr = q), (K.abortable = C);
      let j = function () {
          if (
            ((K.xhr = null),
            (K.abortable = !1),
            K.queuedRequests != null && K.queuedRequests.length > 0)
          )
            K.queuedRequests.shift()();
        },
        N = getClosestAttributeValue(Z, "hx-prompt");
      if (N) {
        var E = prompt(N);
        if (
          E === null ||
          !triggerEvent(Z, "htmx:prompt", { prompt: E, target: z })
        )
          return maybeCall(Q), j(), U;
      }
      if (L && !X) {
        if (!confirm(L)) return maybeCall(Q), j(), U;
      }
      let k = getHeaders(Z, z, E);
      if (J !== "get" && !usesFormData(Z))
        k["Content-Type"] = "application/x-www-form-urlencoded";
      if (W.headers) k = mergeObjects(k, W.headers);
      let I = getInputValues(Z, J),
        g = I.errors,
        JJ = I.formData;
      if (W.values) overrideFormData(JJ, formDataFromObject(W.values));
      let dJ = formDataFromObject(getExpressionVars(Z)),
        mJ = overrideFormData(JJ, dJ),
        i = filterValues(mJ, Z);
      if (htmx.config.getCacheBusterParam && J === "get")
        i.set("org.htmx.cache-buster", getRawAttribute(z, "id") || "true");
      if (Y == null || Y === "") Y = getDocument().location.href;
      let pJ = getValuesForElement(Z, "hx-request"),
        yY = getInternalData(Z).boosted,
        CJ = htmx.config.methodsThatUseUrlParams.indexOf(J) >= 0,
        b = {
          boosted: yY,
          useUrlParams: CJ,
          formData: i,
          parameters: formDataProxy(i),
          unfilteredFormData: mJ,
          unfilteredParameters: formDataProxy(mJ),
          headers: k,
          target: z,
          verb: J,
          errors: g,
          withCredentials:
            W.credentials || pJ.credentials || htmx.config.withCredentials,
          timeout: W.timeout || pJ.timeout || htmx.config.timeout,
          path: Y,
          triggeringEvent: $,
        };
      if (!triggerEvent(Z, "htmx:configRequest", b))
        return maybeCall(Q), j(), U;
      if (
        ((Y = b.path),
        (J = b.verb),
        (k = b.headers),
        (i = formDataFromObject(b.parameters)),
        (g = b.errors),
        (CJ = b.useUrlParams),
        g && g.length > 0)
      )
        return (
          triggerEvent(Z, "htmx:validation:halted", b), maybeCall(Q), j(), U
        );
      let hY = Y.split("#"),
        z0 = hY[0],
        iJ = hY[1],
        f = Y;
      if (CJ) {
        if (((f = z0), !i.keys().next().done)) {
          if (f.indexOf("?") < 0) f += "?";
          else f += "&";
          if (((f += urlEncode(i)), iJ)) f += "#" + iJ;
        }
      }
      if (!verifyPath(Z, f, b))
        return triggerErrorEvent(Z, "htmx:invalidPath", b), maybeCall(B), U;
      if (
        (q.open(J.toUpperCase(), f, !0),
        q.overrideMimeType("text/html"),
        (q.withCredentials = b.withCredentials),
        (q.timeout = b.timeout),
        pJ.noHeaders)
      );
      else
        for (let F in k)
          if (k.hasOwnProperty(F)) {
            let V = k[F];
            safelySetHeaderValue(q, F, V);
          }
      let D = {
        xhr: q,
        target: z,
        requestConfig: b,
        etc: W,
        boosted: yY,
        select: G,
        pathInfo: {
          requestPath: Y,
          finalRequestPath: f,
          responsePath: null,
          anchor: iJ,
        },
      };
      if (
        ((q.onload = function () {
          try {
            let F = hierarchyForElt(Z);
            if (
              ((D.pathInfo.responsePath = getPathFromResponse(q)),
              _(Z, D),
              D.keepIndicators !== !0)
            )
              removeRequestIndicators(FJ, jJ);
            if (
              (triggerEvent(Z, "htmx:afterRequest", D),
              triggerEvent(Z, "htmx:afterOnLoad", D),
              !bodyContains(Z))
            ) {
              let V = null;
              while (F.length > 0 && V == null) {
                let v = F.shift();
                if (bodyContains(v)) V = v;
              }
              if (V)
                triggerEvent(V, "htmx:afterRequest", D),
                  triggerEvent(V, "htmx:afterOnLoad", D);
            }
            maybeCall(Q), j();
          } catch (F) {
            throw (
              (triggerErrorEvent(
                Z,
                "htmx:onLoadError",
                mergeObjects({ error: F }, D),
              ),
              F)
            );
          }
        }),
        (q.onerror = function () {
          removeRequestIndicators(FJ, jJ),
            triggerErrorEvent(Z, "htmx:afterRequest", D),
            triggerErrorEvent(Z, "htmx:sendError", D),
            maybeCall(B),
            j();
        }),
        (q.onabort = function () {
          removeRequestIndicators(FJ, jJ),
            triggerErrorEvent(Z, "htmx:afterRequest", D),
            triggerErrorEvent(Z, "htmx:sendAbort", D),
            maybeCall(B),
            j();
        }),
        (q.ontimeout = function () {
          removeRequestIndicators(FJ, jJ),
            triggerErrorEvent(Z, "htmx:afterRequest", D),
            triggerErrorEvent(Z, "htmx:timeout", D),
            maybeCall(B),
            j();
        }),
        !triggerEvent(Z, "htmx:beforeRequest", D))
      )
        return maybeCall(Q), j(), U;
      var FJ = addRequestIndicatorClasses(Z),
        jJ = disableElements(Z);
      forEach(["loadstart", "loadend", "progress", "abort"], function (F) {
        forEach([q, q.upload], function (V) {
          V.addEventListener(F, function (v) {
            triggerEvent(Z, "htmx:xhr:" + F, {
              lengthComputable: v.lengthComputable,
              loaded: v.loaded,
              total: v.total,
            });
          });
        });
      }),
        triggerEvent(Z, "htmx:beforeSend", D);
      let K0 = CJ ? null : encodeParamsForBody(q, Z, i);
      return q.send(K0), U;
    }
    function determineHistoryUpdates(J, Y) {
      let Z = Y.xhr,
        $ = null,
        W = null;
      if (hasHeader(Z, /HX-Push:/i))
        ($ = Z.getResponseHeader("HX-Push")), (W = "push");
      else if (hasHeader(Z, /HX-Push-Url:/i))
        ($ = Z.getResponseHeader("HX-Push-Url")), (W = "push");
      else if (hasHeader(Z, /HX-Replace-Url:/i))
        ($ = Z.getResponseHeader("HX-Replace-Url")), (W = "replace");
      if ($)
        if ($ === "false") return {};
        else return { type: W, path: $ };
      let X = Y.pathInfo.finalRequestPath,
        Q = Y.pathInfo.responsePath,
        B = getClosestAttributeValue(J, "hx-push-url"),
        U = getClosestAttributeValue(J, "hx-replace-url"),
        _ = getInternalData(J).boosted,
        G = null,
        z = null;
      if (B) (G = "push"), (z = B);
      else if (U) (G = "replace"), (z = U);
      else if (_) (G = "push"), (z = Q || X);
      if (z) {
        if (z === "false") return {};
        if (z === "true") z = Q || X;
        if (Y.pathInfo.anchor && z.indexOf("#") === -1)
          z = z + "#" + Y.pathInfo.anchor;
        return { type: G, path: z };
      } else return {};
    }
    function codeMatches(J, Y) {
      var Z = new RegExp(J.code);
      return Z.test(Y.toString(10));
    }
    function resolveResponseHandling(J) {
      for (var Y = 0; Y < htmx.config.responseHandling.length; Y++) {
        var Z = htmx.config.responseHandling[Y];
        if (codeMatches(Z, J.status)) return Z;
      }
      return { swap: !1 };
    }
    function handleTitle(J) {
      if (J) {
        let Y = find("title");
        if (Y) Y.innerHTML = J;
        else window.document.title = J;
      }
    }
    function handleAjaxResponse(J, Y) {
      let { xhr: Z, target: $, etc: W, select: X } = Y;
      if (!triggerEvent(J, "htmx:beforeOnLoad", Y)) return;
      if (hasHeader(Z, /HX-Trigger:/i)) handleTriggerHeader(Z, "HX-Trigger", J);
      if (hasHeader(Z, /HX-Location:/i)) {
        saveCurrentPageToHistory();
        let C = Z.getResponseHeader("HX-Location");
        var Q;
        if (C.indexOf("{") === 0)
          (Q = parseJSON(C)), (C = Q.path), delete Q.path;
        ajaxHelper("get", C, Q).then(function () {
          pushUrlIntoHistory(C);
        });
        return;
      }
      let B =
        hasHeader(Z, /HX-Refresh:/i) &&
        Z.getResponseHeader("HX-Refresh") === "true";
      if (hasHeader(Z, /HX-Redirect:/i)) {
        (Y.keepIndicators = !0),
          (location.href = Z.getResponseHeader("HX-Redirect")),
          B && location.reload();
        return;
      }
      if (B) {
        (Y.keepIndicators = !0), location.reload();
        return;
      }
      if (hasHeader(Z, /HX-Retarget:/i))
        if (Z.getResponseHeader("HX-Retarget") === "this") Y.target = J;
        else
          Y.target = asElement(
            querySelectorExt(J, Z.getResponseHeader("HX-Retarget")),
          );
      let U = determineHistoryUpdates(J, Y),
        _ = resolveResponseHandling(Z),
        G = _.swap,
        z = !!_.error,
        K = htmx.config.ignoreTitle || _.ignoreTitle,
        A = _.select;
      if (_.target) Y.target = asElement(querySelectorExt(J, _.target));
      var L = W.swapOverride;
      if (L == null && _.swapOverride) L = _.swapOverride;
      if (hasHeader(Z, /HX-Retarget:/i))
        if (Z.getResponseHeader("HX-Retarget") === "this") Y.target = J;
        else
          Y.target = asElement(
            querySelectorExt(J, Z.getResponseHeader("HX-Retarget")),
          );
      if (hasHeader(Z, /HX-Reswap:/i)) L = Z.getResponseHeader("HX-Reswap");
      var H = Z.response,
        M = mergeObjects(
          {
            shouldSwap: G,
            serverResponse: H,
            isError: z,
            ignoreTitle: K,
            selectOverride: A,
            swapOverride: L,
          },
          Y,
        );
      if (_.event && !triggerEvent($, _.event, M)) return;
      if (!triggerEvent($, "htmx:beforeSwap", M)) return;
      if (
        (($ = M.target),
        (H = M.serverResponse),
        (z = M.isError),
        (K = M.ignoreTitle),
        (A = M.selectOverride),
        (L = M.swapOverride),
        (Y.target = $),
        (Y.failed = z),
        (Y.successful = !z),
        M.shouldSwap)
      ) {
        if (Z.status === 286) cancelPolling(J);
        if (
          (withExtensions(J, function (I) {
            H = I.transformResponse(H, Z, J);
          }),
          U.type)
        )
          saveCurrentPageToHistory();
        var P = getSwapSpecification(J, L);
        if (!P.hasOwnProperty("ignoreTitle")) P.ignoreTitle = K;
        $.classList.add(htmx.config.swappingClass);
        let C = null,
          q = null;
        if (X) A = X;
        if (hasHeader(Z, /HX-Reselect:/i))
          A = Z.getResponseHeader("HX-Reselect");
        let j = getClosestAttributeValue(J, "hx-select-oob"),
          N = getClosestAttributeValue(J, "hx-select"),
          E = function () {
            try {
              if (U.type)
                if (
                  (triggerEvent(
                    getDocument().body,
                    "htmx:beforeHistoryUpdate",
                    mergeObjects({ history: U }, Y),
                  ),
                  U.type === "push")
                )
                  pushUrlIntoHistory(U.path),
                    triggerEvent(getDocument().body, "htmx:pushedIntoHistory", {
                      path: U.path,
                    });
                else
                  replaceUrlInHistory(U.path),
                    triggerEvent(getDocument().body, "htmx:replacedInHistory", {
                      path: U.path,
                    });
              swap($, H, P, {
                select: A || N,
                selectOOB: j,
                eventInfo: Y,
                anchor: Y.pathInfo.anchor,
                contextElement: J,
                afterSwapCallback: function () {
                  if (hasHeader(Z, /HX-Trigger-After-Swap:/i)) {
                    let I = J;
                    if (!bodyContains(J)) I = getDocument().body;
                    handleTriggerHeader(Z, "HX-Trigger-After-Swap", I);
                  }
                },
                afterSettleCallback: function () {
                  if (hasHeader(Z, /HX-Trigger-After-Settle:/i)) {
                    let I = J;
                    if (!bodyContains(J)) I = getDocument().body;
                    handleTriggerHeader(Z, "HX-Trigger-After-Settle", I);
                  }
                  maybeCall(C);
                },
              });
            } catch (I) {
              throw (
                (triggerErrorEvent(J, "htmx:swapError", Y), maybeCall(q), I)
              );
            }
          },
          k = htmx.config.globalViewTransitions;
        if (P.hasOwnProperty("transition")) k = P.transition;
        if (
          k &&
          triggerEvent(J, "htmx:beforeTransition", Y) &&
          typeof Promise !== "undefined" &&
          document.startViewTransition
        ) {
          let I = new Promise(function (JJ, dJ) {
              (C = JJ), (q = dJ);
            }),
            g = E;
          E = function () {
            document.startViewTransition(function () {
              return g(), I;
            });
          };
        }
        if (P.swapDelay > 0) getWindow().setTimeout(E, P.swapDelay);
        else E();
      }
      if (z)
        triggerErrorEvent(
          J,
          "htmx:responseError",
          mergeObjects(
            {
              error:
                "Response Status Error Code " +
                Z.status +
                " from " +
                Y.pathInfo.requestPath,
            },
            Y,
          ),
        );
    }
    let extensions = {};
    function extensionBase() {
      return {
        init: function (J) {
          return null;
        },
        getSelectors: function () {
          return null;
        },
        onEvent: function (J, Y) {
          return !0;
        },
        transformResponse: function (J, Y, Z) {
          return J;
        },
        isInlineSwap: function (J) {
          return !1;
        },
        handleSwap: function (J, Y, Z, $) {
          return !1;
        },
        encodeParameters: function (J, Y, Z) {
          return null;
        },
      };
    }
    function defineExtension(J, Y) {
      if (Y.init) Y.init(internalAPI);
      extensions[J] = mergeObjects(extensionBase(), Y);
    }
    function removeExtension(J) {
      delete extensions[J];
    }
    function getExtensions(J, Y, Z) {
      if (Y == null) Y = [];
      if (J == null) return Y;
      if (Z == null) Z = [];
      let $ = getAttributeValue(J, "hx-ext");
      if ($)
        forEach($.split(","), function (W) {
          if (((W = W.replace(/ /g, "")), W.slice(0, 7) == "ignore:")) {
            Z.push(W.slice(7));
            return;
          }
          if (Z.indexOf(W) < 0) {
            let X = extensions[W];
            if (X && Y.indexOf(X) < 0) Y.push(X);
          }
        });
      return getExtensions(asElement(parentElt(J)), Y, Z);
    }
    var isReady = !1;
    getDocument().addEventListener("DOMContentLoaded", function () {
      isReady = !0;
    });
    function ready(J) {
      if (isReady || getDocument().readyState === "complete") J();
      else getDocument().addEventListener("DOMContentLoaded", J);
    }
    function insertIndicatorStyles() {
      if (htmx.config.includeIndicatorStyles !== !1) {
        let J = htmx.config.inlineStyleNonce
          ? ` nonce="${htmx.config.inlineStyleNonce}"`
          : "";
        getDocument().head.insertAdjacentHTML(
          "beforeend",
          "<style" +
            J +
            ">      ." +
            htmx.config.indicatorClass +
            "{opacity:0}      ." +
            htmx.config.requestClass +
            " ." +
            htmx.config.indicatorClass +
            "{opacity:1; transition: opacity 200ms ease-in;}      ." +
            htmx.config.requestClass +
            "." +
            htmx.config.indicatorClass +
            "{opacity:1; transition: opacity 200ms ease-in;}      </style>",
        );
      }
    }
    function getMetaConfig() {
      let J = getDocument().querySelector('meta[name="htmx-config"]');
      if (J) return parseJSON(J.content);
      else return null;
    }
    function mergeMetaConfig() {
      let J = getMetaConfig();
      if (J) htmx.config = mergeObjects(htmx.config, J);
    }
    return (
      ready(function () {
        mergeMetaConfig(), insertIndicatorStyles();
        let J = getDocument().body;
        processNode(J);
        let Y = getDocument().querySelectorAll(
          "[hx-trigger='restored'],[data-hx-trigger='restored']",
        );
        J.addEventListener("htmx:abort", function ($) {
          let W = $.target,
            X = getInternalData(W);
          if (X && X.xhr) X.xhr.abort();
        });
        let Z = window.onpopstate ? window.onpopstate.bind(window) : null;
        (window.onpopstate = function ($) {
          if ($.state && $.state.htmx)
            restoreHistory(),
              forEach(Y, function (W) {
                triggerEvent(W, "htmx:restored", {
                  document: getDocument(),
                  triggerEvent,
                });
              });
          else if (Z) Z($);
        }),
          getWindow().setTimeout(function () {
            triggerEvent(J, "htmx:load", {}), (J = null);
          }, 0);
      }),
      htmx
    );
  })(),
  UJ = M0;
(function () {
  UJ.defineExtension("preload", {
    onEvent: function (G, z) {
      if (G === "htmx:afterProcessNode") {
        let K = z.target || z.detail.elt;
        [
          ...(K.hasAttribute("preload") ? [K] : []),
          ...K.querySelectorAll("[preload]"),
        ].forEach(function (L) {
          J(L), L.querySelectorAll("[href],[hx-get],[data-hx-get]").forEach(J);
        });
        return;
      }
      if (G === "htmx:beforeRequest") {
        let K = z.detail.requestConfig.headers;
        if (!("HX-Preloaded" in K && K["HX-Preloaded"] === "true")) return;
        z.preventDefault();
        let A = z.detail.xhr;
        (A.onload = function () {
          Q(z.detail.elt, A.responseText);
        }),
          (A.onerror = null),
          (A.onabort = null),
          (A.ontimeout = null),
          A.send();
      }
    },
  });
  function J(G) {
    if (G.preloadState !== void 0) return;
    if (!U(G)) return;
    if (G instanceof HTMLFormElement) {
      let L = G;
      if (
        !(
          (L.hasAttribute("method") && L.method === "get") ||
          L.hasAttribute("hx-get") ||
          L.hasAttribute("hx-data-get")
        )
      )
        return;
      for (let H = 0; H < L.elements.length; H++) {
        let M = L.elements.item(H);
        J(M), M.labels.forEach(J);
      }
      return;
    }
    let z = B(G, "preload");
    if (((G.preloadAlways = z && z.includes("always")), G.preloadAlways))
      z = z.replace("always", "").trim();
    let K = z || "mousedown",
      A = K === "mouseover";
    if (
      (G.addEventListener(K, Y(G, A)), K === "mousedown" || K === "mouseover")
    )
      G.addEventListener("touchstart", Y(G));
    if (K === "mouseover")
      G.addEventListener("mouseout", function (L) {
        if (L.target === G && G.preloadState === "TIMEOUT")
          G.preloadState = "READY";
      });
    (G.preloadState = "READY"), UJ.trigger(G, "preload:init");
  }
  function Y(G, z = !1) {
    return function () {
      if (G.preloadState !== "READY") return;
      if (z) {
        G.preloadState = "TIMEOUT";
        let K = 100;
        window.setTimeout(function () {
          if (G.preloadState === "TIMEOUT") (G.preloadState = "READY"), Z(G);
        }, K);
        return;
      }
      Z(G);
    };
  }
  function Z(G) {
    if (G.preloadState !== "READY") return;
    G.preloadState = "LOADING";
    let z = G.getAttribute("hx-get") || G.getAttribute("data-hx-get");
    if (z) {
      W(z, G);
      return;
    }
    let K = B(G, "hx-boost") === "true";
    if (G.hasAttribute("href")) {
      let A = G.getAttribute("href");
      if (K) W(A, G);
      else X(A, G);
      return;
    }
    if (_(G)) {
      let A =
          G.form.getAttribute("action") ||
          G.form.getAttribute("hx-get") ||
          G.form.getAttribute("data-hx-get"),
        L = UJ.values(G.form),
        M = !(
          G.form.getAttribute("hx-get") ||
          G.form.getAttribute("data-hx-get") ||
          K
        )
          ? X
          : W;
      if (G.type === "submit") {
        M(A, G.form, L);
        return;
      }
      let P = G.name || G.control.name;
      if (G.tagName === "SELECT") {
        Array.from(G.options).forEach((N) => {
          if (N.selected) return;
          L.set(P, N.value);
          let E = $(G.form, L);
          M(A, G.form, E);
        });
        return;
      }
      let C = G.getAttribute("type") || G.control.getAttribute("type"),
        q = G.value || G.control?.value;
      if (C === "radio") L.set(P, q);
      else if (C === "checkbox") {
        let N = L.getAll(P);
        if (N.includes(q)) L[P] = N.filter((E) => E !== q);
        else L.append(P, q);
      }
      let j = $(G.form, L);
      M(A, G.form, j);
      return;
    }
  }
  function $(G, z) {
    let K = G.elements,
      A = new FormData();
    for (let L = 0; L < K.length; L++) {
      let H = K.item(L);
      if (z.has(H.name) && H.tagName === "SELECT") {
        A.append(H.name, z.get(H.name));
        continue;
      }
      if (z.has(H.name) && z.getAll(H.name).includes(H.value))
        A.append(H.name, H.value);
    }
    return A;
  }
  function W(G, z, K = void 0) {
    UJ.ajax("GET", G, {
      source: z,
      values: K,
      headers: { "HX-Preloaded": "true" },
    });
  }
  function X(G, z, K = void 0) {
    let A = new XMLHttpRequest();
    if (K) G += "?" + new URLSearchParams(K.entries()).toString();
    A.open("GET", G),
      A.setRequestHeader("HX-Preloaded", "true"),
      (A.onload = function () {
        Q(z, A.responseText);
      }),
      A.send();
  }
  function Q(G, z) {
    if (
      ((G.preloadState = G.preloadAlways ? "READY" : "DONE"),
      B(G, "preload-images") === "true")
    )
      document.createElement("div").innerHTML = z;
  }
  function B(G, z) {
    if (G == null) return;
    return (
      G.getAttribute(z) || G.getAttribute("data-" + z) || B(G.parentElement, z)
    );
  }
  function U(G) {
    let z = ["href", "hx-get", "data-hx-get"],
      K = (L) => z.some((H) => L.hasAttribute(H)) || L.method === "get",
      A = G.form instanceof HTMLFormElement && K(G.form) && _(G);
    if (!K(G) && !A) return !1;
    if (G instanceof HTMLInputElement && G.closest("label")) return !1;
    return !0;
  }
  function _(G) {
    if (G instanceof HTMLInputElement || G instanceof HTMLButtonElement) {
      let z = G.getAttribute("type");
      return ["checkbox", "radio", "submit"].includes(z);
    }
    if (G instanceof HTMLLabelElement) return G.control && _(G.control);
    return G instanceof HTMLSelectElement;
  }
})();
var rJ = !1,
  tJ = !1,
  n = [],
  eJ = -1;
function A0(J) {
  L0(J);
}
function L0(J) {
  if (!n.includes(J)) n.push(J);
  q0();
}
function P0(J) {
  let Y = n.indexOf(J);
  if (Y !== -1 && Y > eJ) n.splice(Y, 1);
}
function q0() {
  if (!tJ && !rJ) (rJ = !0), queueMicrotask(H0);
}
function H0() {
  (rJ = !1), (tJ = !0);
  for (let J = 0; J < n.length; J++) n[J](), (eJ = J);
  (n.length = 0), (eJ = -1), (tJ = !1);
}
var $J,
  e,
  WJ,
  rY,
  JY = !0;
function C0(J) {
  (JY = !1), J(), (JY = !0);
}
function F0(J) {
  ($J = J.reactive),
    (WJ = J.release),
    (e = (Y) =>
      J.effect(Y, {
        scheduler: (Z) => {
          if (JY) A0(Z);
          else Z();
        },
      })),
    (rY = J.raw);
}
function fY(J) {
  e = J;
}
function j0(J) {
  let Y = () => {};
  return [
    ($) => {
      let W = e($);
      if (!J._x_effects)
        (J._x_effects = new Set()),
          (J._x_runEffects = () => {
            J._x_effects.forEach((X) => X());
          });
      return (
        J._x_effects.add(W),
        (Y = () => {
          if (W === void 0) return;
          J._x_effects.delete(W), WJ(W);
        }),
        W
      );
    },
    () => {
      Y();
    },
  ];
}
function tY(J, Y) {
  let Z = !0,
    $,
    W = e(() => {
      let X = J();
      if ((JSON.stringify(X), !Z))
        queueMicrotask(() => {
          Y(X, $), ($ = X);
        });
      else $ = X;
      Z = !1;
    });
  return () => WJ(W);
}
var eY = [],
  JZ = [],
  YZ = [];
function R0(J) {
  YZ.push(J);
}
function AY(J, Y) {
  if (typeof Y === "function") {
    if (!J._x_cleanups) J._x_cleanups = [];
    J._x_cleanups.push(Y);
  } else (Y = J), JZ.push(Y);
}
function ZZ(J) {
  eY.push(J);
}
function $Z(J, Y, Z) {
  if (!J._x_attributeCleanups) J._x_attributeCleanups = {};
  if (!J._x_attributeCleanups[Y]) J._x_attributeCleanups[Y] = [];
  J._x_attributeCleanups[Y].push(Z);
}
function WZ(J, Y) {
  if (!J._x_attributeCleanups) return;
  Object.entries(J._x_attributeCleanups).forEach(([Z, $]) => {
    if (Y === void 0 || Y.includes(Z))
      $.forEach((W) => W()), delete J._x_attributeCleanups[Z];
  });
}
function O0(J) {
  J._x_effects?.forEach(P0);
  while (J._x_cleanups?.length) J._x_cleanups.pop()();
}
var LY = new MutationObserver(CY),
  PY = !1;
function qY() {
  LY.observe(document, {
    subtree: !0,
    childList: !0,
    attributes: !0,
    attributeOldValue: !0,
  }),
    (PY = !0);
}
function XZ() {
  T0(), LY.disconnect(), (PY = !1);
}
var GJ = [];
function T0() {
  let J = LY.takeRecords();
  GJ.push(() => J.length > 0 && CY(J));
  let Y = GJ.length;
  queueMicrotask(() => {
    if (GJ.length === Y) while (GJ.length > 0) GJ.shift()();
  });
}
function O(J) {
  if (!PY) return J();
  XZ();
  let Y = J();
  return qY(), Y;
}
var HY = !1,
  DJ = [];
function V0() {
  HY = !0;
}
function E0() {
  (HY = !1), CY(DJ), (DJ = []);
}
function CY(J) {
  if (HY) {
    DJ = DJ.concat(J);
    return;
  }
  let Y = [],
    Z = new Set(),
    $ = new Map(),
    W = new Map();
  for (let X = 0; X < J.length; X++) {
    if (J[X].target._x_ignoreMutationObserver) continue;
    if (J[X].type === "childList")
      J[X].removedNodes.forEach((Q) => {
        if (Q.nodeType !== 1) return;
        if (!Q._x_marker) return;
        Z.add(Q);
      }),
        J[X].addedNodes.forEach((Q) => {
          if (Q.nodeType !== 1) return;
          if (Z.has(Q)) {
            Z.delete(Q);
            return;
          }
          if (Q._x_marker) return;
          Y.push(Q);
        });
    if (J[X].type === "attributes") {
      let Q = J[X].target,
        B = J[X].attributeName,
        U = J[X].oldValue,
        _ = () => {
          if (!$.has(Q)) $.set(Q, []);
          $.get(Q).push({ name: B, value: Q.getAttribute(B) });
        },
        G = () => {
          if (!W.has(Q)) W.set(Q, []);
          W.get(Q).push(B);
        };
      if (Q.hasAttribute(B) && U === null) _();
      else if (Q.hasAttribute(B)) G(), _();
      else G();
    }
  }
  W.forEach((X, Q) => {
    WZ(Q, X);
  }),
    $.forEach((X, Q) => {
      eY.forEach((B) => B(Q, X));
    });
  for (let X of Z) {
    if (Y.some((Q) => Q.contains(X))) continue;
    JZ.forEach((Q) => Q(X));
  }
  for (let X of Y) {
    if (!X.isConnected) continue;
    YZ.forEach((Q) => Q(X));
  }
  (Y = null), (Z = null), ($ = null), (W = null);
}
function QZ(J) {
  return qJ(YJ(J));
}
function PJ(J, Y, Z) {
  return (
    (J._x_dataStack = [Y, ...YJ(Z || J)]),
    () => {
      J._x_dataStack = J._x_dataStack.filter(($) => $ !== Y);
    }
  );
}
function YJ(J) {
  if (J._x_dataStack) return J._x_dataStack;
  if (typeof ShadowRoot === "function" && J instanceof ShadowRoot)
    return YJ(J.host);
  if (!J.parentNode) return [];
  return YJ(J.parentNode);
}
function qJ(J) {
  return new Proxy({ objects: J }, N0);
}
var N0 = {
  ownKeys({ objects: J }) {
    return Array.from(new Set(J.flatMap((Y) => Object.keys(Y))));
  },
  has({ objects: J }, Y) {
    if (Y == Symbol.unscopables) return !1;
    return J.some(
      (Z) => Object.prototype.hasOwnProperty.call(Z, Y) || Reflect.has(Z, Y),
    );
  },
  get({ objects: J }, Y, Z) {
    if (Y == "toJSON") return I0;
    return Reflect.get(J.find(($) => Reflect.has($, Y)) || {}, Y, Z);
  },
  set({ objects: J }, Y, Z, $) {
    let W =
        J.find((Q) => Object.prototype.hasOwnProperty.call(Q, Y)) ||
        J[J.length - 1],
      X = Object.getOwnPropertyDescriptor(W, Y);
    if (X?.set && X?.get) return X.set.call($, Z) || !0;
    return Reflect.set(W, Y, Z);
  },
};
function I0() {
  return Reflect.ownKeys(this).reduce((Y, Z) => {
    return (Y[Z] = Reflect.get(this, Z)), Y;
  }, {});
}
function BZ(J) {
  let Y = ($) => typeof $ === "object" && !Array.isArray($) && $ !== null,
    Z = ($, W = "") => {
      Object.entries(Object.getOwnPropertyDescriptors($)).forEach(
        ([X, { value: Q, enumerable: B }]) => {
          if (B === !1 || Q === void 0) return;
          if (typeof Q === "object" && Q !== null && Q.__v_skip) return;
          let U = W === "" ? X : `${W}.${X}`;
          if (typeof Q === "object" && Q !== null && Q._x_interceptor)
            $[X] = Q.initialize(J, U, X);
          else if (Y(Q) && Q !== $ && !(Q instanceof Element)) Z(Q, U);
        },
      );
    };
  return Z(J);
}
function UZ(J, Y = () => {}) {
  let Z = {
    initialValue: void 0,
    _x_interceptor: !0,
    initialize($, W, X) {
      return J(
        this.initialValue,
        () => D0($, W),
        (Q) => YY($, W, Q),
        W,
        X,
      );
    },
  };
  return (
    Y(Z),
    ($) => {
      if (typeof $ === "object" && $ !== null && $._x_interceptor) {
        let W = Z.initialize.bind(Z);
        Z.initialize = (X, Q, B) => {
          let U = $.initialize(X, Q, B);
          return (Z.initialValue = U), W(X, Q, B);
        };
      } else Z.initialValue = $;
      return Z;
    }
  );
}
function D0(J, Y) {
  return Y.split(".").reduce((Z, $) => Z[$], J);
}
function YY(J, Y, Z) {
  if (typeof Y === "string") Y = Y.split(".");
  if (Y.length === 1) J[Y[0]] = Z;
  else if (Y.length === 0) throw error;
  else if (J[Y[0]]) return YY(J[Y[0]], Y.slice(1), Z);
  else return (J[Y[0]] = {}), YY(J[Y[0]], Y.slice(1), Z);
}
var GZ = {};
function y(J, Y) {
  GZ[J] = Y;
}
function ZY(J, Y) {
  let Z = w0(Y);
  return (
    Object.entries(GZ).forEach(([$, W]) => {
      Object.defineProperty(J, `$${$}`, {
        get() {
          return W(Y, Z);
        },
        enumerable: !1,
      });
    }),
    J
  );
}
function w0(J) {
  let [Y, Z] = LZ(J),
    $ = { interceptor: UZ, ...Y };
  return AY(J, Z), $;
}
function k0(J, Y, Z, ...$) {
  try {
    return Z(...$);
  } catch (W) {
    LJ(W, J, Y);
  }
}
function LJ(J, Y, Z = void 0) {
  (J = Object.assign(J ?? { message: "No error message given." }, {
    el: Y,
    expression: Z,
  })),
    console.warn(
      `Alpine Expression Error: ${J.message}

${
  Z
    ? 'Expression: "' +
      Z +
      `"

`
    : ""
}`,
      Y,
    ),
    setTimeout(() => {
      throw J;
    }, 0);
}
var NJ = !0;
function _Z(J) {
  let Y = NJ;
  NJ = !1;
  let Z = J();
  return (NJ = Y), Z;
}
function l(J, Y, Z = {}) {
  let $;
  return w(J, Y)((W) => ($ = W), Z), $;
}
function w(...J) {
  return zZ(...J);
}
var zZ = KZ;
function b0(J) {
  zZ = J;
}
function KZ(J, Y) {
  let Z = {};
  ZY(Z, J);
  let $ = [Z, ...YJ(J)],
    W = typeof Y === "function" ? S0($, Y) : y0($, Y, J);
  return k0.bind(null, J, Y, W);
}
function S0(J, Y) {
  return (Z = () => {}, { scope: $ = {}, params: W = [] } = {}) => {
    let X = Y.apply(qJ([$, ...J]), W);
    wJ(Z, X);
  };
}
var sJ = {};
function x0(J, Y) {
  if (sJ[J]) return sJ[J];
  let Z = Object.getPrototypeOf(async function () {}).constructor,
    $ =
      /^[\n\s]*if.*\(.*\)/.test(J.trim()) || /^(let|const)\s/.test(J.trim())
        ? `(async()=>{ ${J} })()`
        : J,
    X = (() => {
      try {
        let Q = new Z(
          ["__self", "scope"],
          `with (scope) { __self.result = ${$} }; __self.finished = true; return __self.result;`,
        );
        return Object.defineProperty(Q, "name", { value: `[Alpine] ${J}` }), Q;
      } catch (Q) {
        return LJ(Q, Y, J), Promise.resolve();
      }
    })();
  return (sJ[J] = X), X;
}
function y0(J, Y, Z) {
  let $ = x0(Y, Z);
  return (W = () => {}, { scope: X = {}, params: Q = [] } = {}) => {
    ($.result = void 0), ($.finished = !1);
    let B = qJ([X, ...J]);
    if (typeof $ === "function") {
      let U = $($, B).catch((_) => LJ(_, Z, Y));
      if ($.finished) wJ(W, $.result, B, Q, Z), ($.result = void 0);
      else
        U.then((_) => {
          wJ(W, _, B, Q, Z);
        })
          .catch((_) => LJ(_, Z, Y))
          .finally(() => ($.result = void 0));
    }
  };
}
function wJ(J, Y, Z, $, W) {
  if (NJ && typeof Y === "function") {
    let X = Y.apply(Z, $);
    if (X instanceof Promise)
      X.then((Q) => wJ(J, Q, Z, $)).catch((Q) => LJ(Q, W, Y));
    else J(X);
  } else if (typeof Y === "object" && Y instanceof Promise) Y.then((X) => J(X));
  else J(Y);
}
var FY = "x-";
function XJ(J = "") {
  return FY + J;
}
function h0(J) {
  FY = J;
}
var kJ = {};
function T(J, Y) {
  return (
    (kJ[J] = Y),
    {
      before(Z) {
        if (!kJ[Z]) {
          console.warn(
            String.raw`Cannot find directive \`${Z}\`. \`${J}\` will use the default order of execution`,
          );
          return;
        }
        let $ = o.indexOf(Z);
        o.splice($ >= 0 ? $ : o.indexOf("DEFAULT"), 0, J);
      },
    }
  );
}
function f0(J) {
  return Object.keys(kJ).includes(J);
}
function jY(J, Y, Z) {
  if (((Y = Array.from(Y)), J._x_virtualDirectives)) {
    let X = Object.entries(J._x_virtualDirectives).map(([B, U]) => ({
        name: B,
        value: U,
      })),
      Q = MZ(X);
    (X = X.map((B) => {
      if (Q.find((U) => U.name === B.name))
        return { name: `x-bind:${B.name}`, value: `"${B.value}"` };
      return B;
    })),
      (Y = Y.concat(X));
  }
  let $ = {};
  return Y.map(HZ((X, Q) => ($[X] = Q)))
    .filter(FZ)
    .map(g0($, Z))
    .sort(c0)
    .map((X) => {
      return u0(J, X);
    });
}
function MZ(J) {
  return Array.from(J)
    .map(HZ())
    .filter((Y) => !FZ(Y));
}
var $Y = !1,
  KJ = new Map(),
  AZ = Symbol();
function v0(J) {
  $Y = !0;
  let Y = Symbol();
  (AZ = Y), KJ.set(Y, []);
  let Z = () => {
      while (KJ.get(Y).length) KJ.get(Y).shift()();
      KJ.delete(Y);
    },
    $ = () => {
      ($Y = !1), Z();
    };
  J(Z), $();
}
function LZ(J) {
  let Y = [],
    Z = (B) => Y.push(B),
    [$, W] = j0(J);
  return (
    Y.push(W),
    [
      {
        Alpine: HJ,
        effect: $,
        cleanup: Z,
        evaluateLater: w.bind(w, J),
        evaluate: l.bind(l, J),
      },
      () => Y.forEach((B) => B()),
    ]
  );
}
function u0(J, Y) {
  let Z = () => {},
    $ = kJ[Y.type] || Z,
    [W, X] = LZ(J);
  $Z(J, Y.original, X);
  let Q = () => {
    if (J._x_ignore || J._x_ignoreSelf) return;
    $.inline && $.inline(J, Y, W),
      ($ = $.bind($, J, Y, W)),
      $Y ? KJ.get(AZ).push($) : $();
  };
  return (Q.runCleanups = X), Q;
}
var PZ =
    (J, Y) =>
    ({ name: Z, value: $ }) => {
      if (Z.startsWith(J)) Z = Z.replace(J, Y);
      return { name: Z, value: $ };
    },
  qZ = (J) => J;
function HZ(J = () => {}) {
  return ({ name: Y, value: Z }) => {
    let { name: $, value: W } = CZ.reduce(
      (X, Q) => {
        return Q(X);
      },
      { name: Y, value: Z },
    );
    if ($ !== Y) J($, Y);
    return { name: $, value: W };
  };
}
var CZ = [];
function RY(J) {
  CZ.push(J);
}
function FZ({ name: J }) {
  return jZ().test(J);
}
var jZ = () => new RegExp(`^${FY}([^:^.]+)\\b`);
function g0(J, Y) {
  return ({ name: Z, value: $ }) => {
    let W = Z.match(jZ()),
      X = Z.match(/:([a-zA-Z0-9\-_:]+)/),
      Q = Z.match(/\.[^.\]]+(?=[^\]]*$)/g) || [],
      B = Y || J[Z] || Z;
    return {
      type: W ? W[1] : null,
      value: X ? X[1] : null,
      modifiers: Q.map((U) => U.replace(".", "")),
      expression: $,
      original: B,
    };
  };
}
var WY = "DEFAULT",
  o = [
    "ignore",
    "ref",
    "data",
    "id",
    "anchor",
    "bind",
    "init",
    "for",
    "model",
    "modelable",
    "transition",
    "show",
    "if",
    WY,
    "teleport",
  ];
function c0(J, Y) {
  let Z = o.indexOf(J.type) === -1 ? WY : J.type,
    $ = o.indexOf(Y.type) === -1 ? WY : Y.type;
  return o.indexOf(Z) - o.indexOf($);
}
function MJ(J, Y, Z = {}) {
  J.dispatchEvent(
    new CustomEvent(Y, {
      detail: Z,
      bubbles: !0,
      composed: !0,
      cancelable: !0,
    }),
  );
}
function t(J, Y) {
  if (typeof ShadowRoot === "function" && J instanceof ShadowRoot) {
    Array.from(J.children).forEach((W) => t(W, Y));
    return;
  }
  let Z = !1;
  if ((Y(J, () => (Z = !0)), Z)) return;
  let $ = J.firstElementChild;
  while ($) t($, Y, !1), ($ = $.nextElementSibling);
}
function S(J, ...Y) {
  console.warn(`Alpine Warning: ${J}`, ...Y);
}
var vY = !1;
function d0() {
  if (vY)
    S(
      "Alpine has already been initialized on this page. Calling Alpine.start() more than once can cause problems.",
    );
  if (((vY = !0), !document.body))
    S(
      "Unable to initialize. Trying to load Alpine before `<body>` is available. Did you forget to add `defer` in Alpine's `<script>` tag?",
    );
  MJ(document, "alpine:init"),
    MJ(document, "alpine:initializing"),
    qY(),
    R0((Y) => u(Y, t)),
    AY((Y) => BJ(Y)),
    ZZ((Y, Z) => {
      jY(Y, Z).forEach(($) => $());
    });
  let J = (Y) => !SJ(Y.parentElement, !0);
  Array.from(document.querySelectorAll(TZ().join(",")))
    .filter(J)
    .forEach((Y) => {
      u(Y);
    }),
    MJ(document, "alpine:initialized"),
    setTimeout(() => {
      s0();
    });
}
var OY = [],
  RZ = [];
function OZ() {
  return OY.map((J) => J());
}
function TZ() {
  return OY.concat(RZ).map((J) => J());
}
function VZ(J) {
  OY.push(J);
}
function EZ(J) {
  RZ.push(J);
}
function SJ(J, Y = !1) {
  return QJ(J, (Z) => {
    if ((Y ? TZ() : OZ()).some((W) => Z.matches(W))) return !0;
  });
}
function QJ(J, Y) {
  if (!J) return;
  if (Y(J)) return J;
  if (J._x_teleportBack) J = J._x_teleportBack;
  if (!J.parentElement) return;
  return QJ(J.parentElement, Y);
}
function m0(J) {
  return OZ().some((Y) => J.matches(Y));
}
var NZ = [];
function p0(J) {
  NZ.push(J);
}
var i0 = 1;
function u(J, Y = t, Z = () => {}) {
  if (QJ(J, ($) => $._x_ignore)) return;
  v0(() => {
    Y(J, ($, W) => {
      if ($._x_marker) return;
      if (
        (Z($, W),
        NZ.forEach((X) => X($, W)),
        jY($, $.attributes).forEach((X) => X()),
        !$._x_ignore)
      )
        $._x_marker = i0++;
      $._x_ignore && W();
    });
  });
}
function BJ(J, Y = t) {
  Y(J, (Z) => {
    O0(Z), WZ(Z), delete Z._x_marker;
  });
}
function s0() {
  [
    ["ui", "dialog", ["[x-dialog], [x-popover]"]],
    ["anchor", "anchor", ["[x-anchor]"]],
    ["sort", "sort", ["[x-sort]"]],
  ].forEach(([Y, Z, $]) => {
    if (f0(Z)) return;
    $.some((W) => {
      if (document.querySelector(W))
        return S(`found "${W}", but missing ${Y} plugin`), !0;
    });
  });
}
var XY = [],
  TY = !1;
function VY(J = () => {}) {
  return (
    queueMicrotask(() => {
      TY ||
        setTimeout(() => {
          QY();
        });
    }),
    new Promise((Y) => {
      XY.push(() => {
        J(), Y();
      });
    })
  );
}
function QY() {
  TY = !1;
  while (XY.length) XY.shift()();
}
function o0() {
  TY = !0;
}
function EY(J, Y) {
  if (Array.isArray(Y)) return uY(J, Y.join(" "));
  else if (typeof Y === "object" && Y !== null) return n0(J, Y);
  else if (typeof Y === "function") return EY(J, Y());
  return uY(J, Y);
}
function uY(J, Y) {
  let Z = (X) => X.split(" ").filter(Boolean),
    $ = (X) =>
      X.split(" ")
        .filter((Q) => !J.classList.contains(Q))
        .filter(Boolean),
    W = (X) => {
      return (
        J.classList.add(...X),
        () => {
          J.classList.remove(...X);
        }
      );
    };
  return (Y = Y === !0 ? (Y = "") : Y || ""), W($(Y));
}
function n0(J, Y) {
  let Z = (B) => B.split(" ").filter(Boolean),
    $ = Object.entries(Y)
      .flatMap(([B, U]) => (U ? Z(B) : !1))
      .filter(Boolean),
    W = Object.entries(Y)
      .flatMap(([B, U]) => (!U ? Z(B) : !1))
      .filter(Boolean),
    X = [],
    Q = [];
  return (
    W.forEach((B) => {
      if (J.classList.contains(B)) J.classList.remove(B), Q.push(B);
    }),
    $.forEach((B) => {
      if (!J.classList.contains(B)) J.classList.add(B), X.push(B);
    }),
    () => {
      Q.forEach((B) => J.classList.add(B)),
        X.forEach((B) => J.classList.remove(B));
    }
  );
}
function xJ(J, Y) {
  if (typeof Y === "object" && Y !== null) return l0(J, Y);
  return a0(J, Y);
}
function l0(J, Y) {
  let Z = {};
  return (
    Object.entries(Y).forEach(([$, W]) => {
      if (((Z[$] = J.style[$]), !$.startsWith("--"))) $ = r0($);
      J.style.setProperty($, W);
    }),
    setTimeout(() => {
      if (J.style.length === 0) J.removeAttribute("style");
    }),
    () => {
      xJ(J, Z);
    }
  );
}
function a0(J, Y) {
  let Z = J.getAttribute("style", Y);
  return (
    J.setAttribute("style", Y),
    () => {
      J.setAttribute("style", Z || "");
    }
  );
}
function r0(J) {
  return J.replace(/([a-z])([A-Z])/g, "$1-$2").toLowerCase();
}
function BY(J, Y = () => {}) {
  let Z = !1;
  return function () {
    if (!Z) (Z = !0), J.apply(this, arguments);
    else Y.apply(this, arguments);
  };
}
T(
  "transition",
  (J, { value: Y, modifiers: Z, expression: $ }, { evaluate: W }) => {
    if (typeof $ === "function") $ = W($);
    if ($ === !1) return;
    if (!$ || typeof $ === "boolean") e0(J, Z, Y);
    else t0(J, $, Y);
  },
);
function t0(J, Y, Z) {
  IZ(J, EY, ""),
    {
      enter: (W) => {
        J._x_transition.enter.during = W;
      },
      "enter-start": (W) => {
        J._x_transition.enter.start = W;
      },
      "enter-end": (W) => {
        J._x_transition.enter.end = W;
      },
      leave: (W) => {
        J._x_transition.leave.during = W;
      },
      "leave-start": (W) => {
        J._x_transition.leave.start = W;
      },
      "leave-end": (W) => {
        J._x_transition.leave.end = W;
      },
    }[Z](Y);
}
function e0(J, Y, Z) {
  IZ(J, xJ);
  let $ = !Y.includes("in") && !Y.includes("out") && !Z,
    W = $ || Y.includes("in") || ["enter"].includes(Z),
    X = $ || Y.includes("out") || ["leave"].includes(Z);
  if (Y.includes("in") && !$) Y = Y.filter((P, C) => C < Y.indexOf("out"));
  if (Y.includes("out") && !$) Y = Y.filter((P, C) => C > Y.indexOf("out"));
  let Q = !Y.includes("opacity") && !Y.includes("scale"),
    B = Q || Y.includes("opacity"),
    U = Q || Y.includes("scale"),
    _ = B ? 0 : 1,
    G = U ? _J(Y, "scale", 95) / 100 : 1,
    z = _J(Y, "delay", 0) / 1000,
    K = _J(Y, "origin", "center"),
    A = "opacity, transform",
    L = _J(Y, "duration", 150) / 1000,
    H = _J(Y, "duration", 75) / 1000,
    M = "cubic-bezier(0.4, 0.0, 0.2, 1)";
  if (W)
    (J._x_transition.enter.during = {
      transformOrigin: K,
      transitionDelay: `${z}s`,
      transitionProperty: A,
      transitionDuration: `${L}s`,
      transitionTimingFunction: M,
    }),
      (J._x_transition.enter.start = { opacity: _, transform: `scale(${G})` }),
      (J._x_transition.enter.end = { opacity: 1, transform: "scale(1)" });
  if (X)
    (J._x_transition.leave.during = {
      transformOrigin: K,
      transitionDelay: `${z}s`,
      transitionProperty: A,
      transitionDuration: `${H}s`,
      transitionTimingFunction: M,
    }),
      (J._x_transition.leave.start = { opacity: 1, transform: "scale(1)" }),
      (J._x_transition.leave.end = { opacity: _, transform: `scale(${G})` });
}
function IZ(J, Y, Z = {}) {
  if (!J._x_transition)
    J._x_transition = {
      enter: { during: Z, start: Z, end: Z },
      leave: { during: Z, start: Z, end: Z },
      in($ = () => {}, W = () => {}) {
        UY(
          J,
          Y,
          {
            during: this.enter.during,
            start: this.enter.start,
            end: this.enter.end,
          },
          $,
          W,
        );
      },
      out($ = () => {}, W = () => {}) {
        UY(
          J,
          Y,
          {
            during: this.leave.during,
            start: this.leave.start,
            end: this.leave.end,
          },
          $,
          W,
        );
      },
    };
}
window.Element.prototype._x_toggleAndCascadeWithTransitions = function (
  J,
  Y,
  Z,
  $,
) {
  let W =
      document.visibilityState === "visible"
        ? requestAnimationFrame
        : setTimeout,
    X = () => W(Z);
  if (Y) {
    if (J._x_transition && (J._x_transition.enter || J._x_transition.leave))
      J._x_transition.enter &&
      (Object.entries(J._x_transition.enter.during).length ||
        Object.entries(J._x_transition.enter.start).length ||
        Object.entries(J._x_transition.enter.end).length)
        ? J._x_transition.in(Z)
        : X();
    else J._x_transition ? J._x_transition.in(Z) : X();
    return;
  }
  (J._x_hidePromise = J._x_transition
    ? new Promise((Q, B) => {
        J._x_transition.out(
          () => {},
          () => Q($),
        ),
          J._x_transitioning &&
            J._x_transitioning.beforeCancel(() =>
              B({ isFromCancelledTransition: !0 }),
            );
      })
    : Promise.resolve($)),
    queueMicrotask(() => {
      let Q = DZ(J);
      if (Q) {
        if (!Q._x_hideChildren) Q._x_hideChildren = [];
        Q._x_hideChildren.push(J);
      } else
        W(() => {
          let B = (U) => {
            let _ = Promise.all([
              U._x_hidePromise,
              ...(U._x_hideChildren || []).map(B),
            ]).then(([G]) => G?.());
            return delete U._x_hidePromise, delete U._x_hideChildren, _;
          };
          B(J).catch((U) => {
            if (!U.isFromCancelledTransition) throw U;
          });
        });
    });
};
function DZ(J) {
  let Y = J.parentNode;
  if (!Y) return;
  return Y._x_hidePromise ? Y : DZ(Y);
}
function UY(
  J,
  Y,
  { during: Z, start: $, end: W } = {},
  X = () => {},
  Q = () => {},
) {
  if (J._x_transitioning) J._x_transitioning.cancel();
  if (
    Object.keys(Z).length === 0 &&
    Object.keys($).length === 0 &&
    Object.keys(W).length === 0
  ) {
    X(), Q();
    return;
  }
  let B, U, _;
  J1(J, {
    start() {
      B = Y(J, $);
    },
    during() {
      U = Y(J, Z);
    },
    before: X,
    end() {
      B(), (_ = Y(J, W));
    },
    after: Q,
    cleanup() {
      U(), _();
    },
  });
}
function J1(J, Y) {
  let Z,
    $,
    W,
    X = BY(() => {
      O(() => {
        if (((Z = !0), !$)) Y.before();
        if (!W) Y.end(), QY();
        if ((Y.after(), J.isConnected)) Y.cleanup();
        delete J._x_transitioning;
      });
    });
  (J._x_transitioning = {
    beforeCancels: [],
    beforeCancel(Q) {
      this.beforeCancels.push(Q);
    },
    cancel: BY(function () {
      while (this.beforeCancels.length) this.beforeCancels.shift()();
      X();
    }),
    finish: X,
  }),
    O(() => {
      Y.start(), Y.during();
    }),
    o0(),
    requestAnimationFrame(() => {
      if (Z) return;
      let Q =
          Number(
            getComputedStyle(J)
              .transitionDuration.replace(/,.*/, "")
              .replace("s", ""),
          ) * 1000,
        B =
          Number(
            getComputedStyle(J)
              .transitionDelay.replace(/,.*/, "")
              .replace("s", ""),
          ) * 1000;
      if (Q === 0)
        Q =
          Number(getComputedStyle(J).animationDuration.replace("s", "")) * 1000;
      O(() => {
        Y.before();
      }),
        ($ = !0),
        requestAnimationFrame(() => {
          if (Z) return;
          O(() => {
            Y.end();
          }),
            QY(),
            setTimeout(J._x_transitioning.finish, Q + B),
            (W = !0);
        });
    });
}
function _J(J, Y, Z) {
  if (J.indexOf(Y) === -1) return Z;
  let $ = J[J.indexOf(Y) + 1];
  if (!$) return Z;
  if (Y === "scale") {
    if (isNaN($)) return Z;
  }
  if (Y === "duration" || Y === "delay") {
    let W = $.match(/([0-9]+)ms/);
    if (W) return W[1];
  }
  if (Y === "origin") {
    if (
      ["top", "right", "left", "center", "bottom"].includes(J[J.indexOf(Y) + 2])
    )
      return [$, J[J.indexOf(Y) + 2]].join(" ");
  }
  return $;
}
var d = !1;
function p(J, Y = () => {}) {
  return (...Z) => (d ? Y(...Z) : J(...Z));
}
function Y1(J) {
  return (...Y) => d && J(...Y);
}
var wZ = [];
function yJ(J) {
  wZ.push(J);
}
function Z1(J, Y) {
  wZ.forEach((Z) => Z(J, Y)),
    (d = !0),
    kZ(() => {
      u(Y, (Z, $) => {
        $(Z, () => {});
      });
    }),
    (d = !1);
}
var GY = !1;
function $1(J, Y) {
  if (!Y._x_dataStack) Y._x_dataStack = J._x_dataStack;
  (d = !0),
    (GY = !0),
    kZ(() => {
      W1(Y);
    }),
    (d = !1),
    (GY = !1);
}
function W1(J) {
  let Y = !1;
  u(J, ($, W) => {
    t($, (X, Q) => {
      if (Y && m0(X)) return Q();
      (Y = !0), W(X, Q);
    });
  });
}
function kZ(J) {
  let Y = e;
  fY((Z, $) => {
    let W = Y(Z);
    return WJ(W), () => {};
  }),
    J(),
    fY(Y);
}
function bZ(J, Y, Z, $ = []) {
  if (!J._x_bindings) J._x_bindings = $J({});
  switch (((J._x_bindings[Y] = Z), (Y = $.includes("camel") ? K1(Y) : Y), Y)) {
    case "value":
      X1(J, Z);
      break;
    case "style":
      B1(J, Z);
      break;
    case "class":
      Q1(J, Z);
      break;
    case "selected":
    case "checked":
      U1(J, Y, Z);
      break;
    default:
      SZ(J, Y, Z);
      break;
  }
}
function X1(J, Y) {
  if (hZ(J)) {
    if (J.attributes.value === void 0) J.value = Y;
    if (window.fromModel)
      if (typeof Y === "boolean") J.checked = IJ(J.value) === Y;
      else J.checked = gY(J.value, Y);
  } else if (NY(J))
    if (Number.isInteger(Y)) J.value = Y;
    else if (
      !Array.isArray(Y) &&
      typeof Y !== "boolean" &&
      ![null, void 0].includes(Y)
    )
      J.value = String(Y);
    else if (Array.isArray(Y)) J.checked = Y.some((Z) => gY(Z, J.value));
    else J.checked = !!Y;
  else if (J.tagName === "SELECT") z1(J, Y);
  else {
    if (J.value === Y) return;
    J.value = Y === void 0 ? "" : Y;
  }
}
function Q1(J, Y) {
  if (J._x_undoAddedClasses) J._x_undoAddedClasses();
  J._x_undoAddedClasses = EY(J, Y);
}
function B1(J, Y) {
  if (J._x_undoAddedStyles) J._x_undoAddedStyles();
  J._x_undoAddedStyles = xJ(J, Y);
}
function U1(J, Y, Z) {
  SZ(J, Y, Z), _1(J, Y, Z);
}
function SZ(J, Y, Z) {
  if ([null, void 0, !1].includes(Z) && A1(Y)) J.removeAttribute(Y);
  else {
    if (xZ(Y)) Z = Y;
    G1(J, Y, Z);
  }
}
function G1(J, Y, Z) {
  if (J.getAttribute(Y) != Z) J.setAttribute(Y, Z);
}
function _1(J, Y, Z) {
  if (J[Y] !== Z) J[Y] = Z;
}
function z1(J, Y) {
  let Z = [].concat(Y).map(($) => {
    return $ + "";
  });
  Array.from(J.options).forEach(($) => {
    $.selected = Z.includes($.value);
  });
}
function K1(J) {
  return J.toLowerCase().replace(/-(\w)/g, (Y, Z) => Z.toUpperCase());
}
function gY(J, Y) {
  return J == Y;
}
function IJ(J) {
  if ([1, "1", "true", "on", "yes", !0].includes(J)) return !0;
  if ([0, "0", "false", "off", "no", !1].includes(J)) return !1;
  return J ? Boolean(J) : null;
}
var M1 = new Set([
  "allowfullscreen",
  "async",
  "autofocus",
  "autoplay",
  "checked",
  "controls",
  "default",
  "defer",
  "disabled",
  "formnovalidate",
  "inert",
  "ismap",
  "itemscope",
  "loop",
  "multiple",
  "muted",
  "nomodule",
  "novalidate",
  "open",
  "playsinline",
  "readonly",
  "required",
  "reversed",
  "selected",
  "shadowrootclonable",
  "shadowrootdelegatesfocus",
  "shadowrootserializable",
]);
function xZ(J) {
  return M1.has(J);
}
function A1(J) {
  return ![
    "aria-pressed",
    "aria-checked",
    "aria-expanded",
    "aria-selected",
  ].includes(J);
}
function L1(J, Y, Z) {
  if (J._x_bindings && J._x_bindings[Y] !== void 0) return J._x_bindings[Y];
  return yZ(J, Y, Z);
}
function P1(J, Y, Z, $ = !0) {
  if (J._x_bindings && J._x_bindings[Y] !== void 0) return J._x_bindings[Y];
  if (J._x_inlineBindings && J._x_inlineBindings[Y] !== void 0) {
    let W = J._x_inlineBindings[Y];
    return (
      (W.extract = $),
      _Z(() => {
        return l(J, W.expression);
      })
    );
  }
  return yZ(J, Y, Z);
}
function yZ(J, Y, Z) {
  let $ = J.getAttribute(Y);
  if ($ === null) return typeof Z === "function" ? Z() : Z;
  if ($ === "") return !0;
  if (xZ(Y)) return !![Y, "true"].includes($);
  return $;
}
function NY(J) {
  return (
    J.type === "checkbox" ||
    J.localName === "ui-checkbox" ||
    J.localName === "ui-switch"
  );
}
function hZ(J) {
  return J.type === "radio" || J.localName === "ui-radio";
}
function fZ(J, Y) {
  var Z;
  return function () {
    var $ = this,
      W = arguments,
      X = function () {
        (Z = null), J.apply($, W);
      };
    clearTimeout(Z), (Z = setTimeout(X, Y));
  };
}
function vZ(J, Y) {
  let Z;
  return function () {
    let $ = this,
      W = arguments;
    if (!Z) J.apply($, W), (Z = !0), setTimeout(() => (Z = !1), Y);
  };
}
function uZ({ get: J, set: Y }, { get: Z, set: $ }) {
  let W = !0,
    X,
    Q,
    B = e(() => {
      let U = J(),
        _ = Z();
      if (W) $(oJ(U)), (W = !1);
      else {
        let G = JSON.stringify(U),
          z = JSON.stringify(_);
        if (G !== X) $(oJ(U));
        else if (G !== z) Y(oJ(_));
      }
      (X = JSON.stringify(J())), (Q = JSON.stringify(Z()));
    });
  return () => {
    WJ(B);
  };
}
function oJ(J) {
  return typeof J === "object" ? JSON.parse(JSON.stringify(J)) : J;
}
function q1(J) {
  (Array.isArray(J) ? J : [J]).forEach((Z) => Z(HJ));
}
var s = {},
  cY = !1;
function H1(J, Y) {
  if (!cY) (s = $J(s)), (cY = !0);
  if (Y === void 0) return s[J];
  if (
    ((s[J] = Y),
    BZ(s[J]),
    typeof Y === "object" &&
      Y !== null &&
      Y.hasOwnProperty("init") &&
      typeof Y.init === "function")
  )
    s[J].init();
}
function C1() {
  return s;
}
var gZ = {};
function F1(J, Y) {
  let Z = typeof Y !== "function" ? () => Y : Y;
  if (J instanceof Element) return cZ(J, Z());
  else gZ[J] = Z;
  return () => {};
}
function j1(J) {
  return (
    Object.entries(gZ).forEach(([Y, Z]) => {
      Object.defineProperty(J, Y, {
        get() {
          return (...$) => {
            return Z(...$);
          };
        },
      });
    }),
    J
  );
}
function cZ(J, Y, Z) {
  let $ = [];
  while ($.length) $.pop()();
  let W = Object.entries(Y).map(([Q, B]) => ({ name: Q, value: B })),
    X = MZ(W);
  return (
    (W = W.map((Q) => {
      if (X.find((B) => B.name === Q.name))
        return { name: `x-bind:${Q.name}`, value: `"${Q.value}"` };
      return Q;
    })),
    jY(J, W, Z).map((Q) => {
      $.push(Q.runCleanups), Q();
    }),
    () => {
      while ($.length) $.pop()();
    }
  );
}
var dZ = {};
function R1(J, Y) {
  dZ[J] = Y;
}
function O1(J, Y) {
  return (
    Object.entries(dZ).forEach(([Z, $]) => {
      Object.defineProperty(J, Z, {
        get() {
          return (...W) => {
            return $.bind(Y)(...W);
          };
        },
        enumerable: !1,
      });
    }),
    J
  );
}
var T1 = {
    get reactive() {
      return $J;
    },
    get release() {
      return WJ;
    },
    get effect() {
      return e;
    },
    get raw() {
      return rY;
    },
    version: "3.14.9",
    flushAndStopDeferringMutations: E0,
    dontAutoEvaluateFunctions: _Z,
    disableEffectScheduling: C0,
    startObservingMutations: qY,
    stopObservingMutations: XZ,
    setReactivityEngine: F0,
    onAttributeRemoved: $Z,
    onAttributesAdded: ZZ,
    closestDataStack: YJ,
    skipDuringClone: p,
    onlyDuringClone: Y1,
    addRootSelector: VZ,
    addInitSelector: EZ,
    interceptClone: yJ,
    addScopeToNode: PJ,
    deferMutations: V0,
    mapAttributes: RY,
    evaluateLater: w,
    interceptInit: p0,
    setEvaluator: b0,
    mergeProxies: qJ,
    extractProp: P1,
    findClosest: QJ,
    onElRemoved: AY,
    closestRoot: SJ,
    destroyTree: BJ,
    interceptor: UZ,
    transition: UY,
    setStyles: xJ,
    mutateDom: O,
    directive: T,
    entangle: uZ,
    throttle: vZ,
    debounce: fZ,
    evaluate: l,
    initTree: u,
    nextTick: VY,
    prefixed: XJ,
    prefix: h0,
    plugin: q1,
    magic: y,
    store: H1,
    start: d0,
    clone: $1,
    cloneNode: Z1,
    bound: L1,
    $data: QZ,
    watch: tY,
    walk: t,
    data: R1,
    bind: F1,
  },
  HJ = T1;
function mZ(J, Y) {
  let Z = Object.create(null),
    $ = J.split(",");
  for (let W = 0; W < $.length; W++) Z[$[W]] = !0;
  return Y ? (W) => !!Z[W.toLowerCase()] : (W) => !!Z[W];
}
var V1 =
    "itemscope,allowfullscreen,formnovalidate,ismap,nomodule,novalidate,readonly",
  k6 = mZ(
    V1 +
      ",async,autofocus,autoplay,controls,default,defer,disabled,hidden,loop,open,required,reversed,scoped,seamless,checked,muted,multiple,selected",
  ),
  E1 = Object.freeze({}),
  b6 = Object.freeze([]),
  N1 = Object.prototype.hasOwnProperty,
  hJ = (J, Y) => N1.call(J, Y),
  a = Array.isArray,
  AJ = (J) => pZ(J) === "[object Map]",
  I1 = (J) => typeof J === "string",
  IY = (J) => typeof J === "symbol",
  fJ = (J) => J !== null && typeof J === "object",
  D1 = Object.prototype.toString,
  pZ = (J) => D1.call(J),
  iZ = (J) => {
    return pZ(J).slice(8, -1);
  },
  DY = (J) =>
    I1(J) && J !== "NaN" && J[0] !== "-" && "" + parseInt(J, 10) === J,
  vJ = (J) => {
    let Y = Object.create(null);
    return (Z) => {
      return Y[Z] || (Y[Z] = J(Z));
    };
  },
  w1 = /-(\w)/g,
  S6 = vJ((J) => {
    return J.replace(w1, (Y, Z) => (Z ? Z.toUpperCase() : ""));
  }),
  k1 = /\B([A-Z])/g,
  x6 = vJ((J) => J.replace(k1, "-$1").toLowerCase()),
  sZ = vJ((J) => J.charAt(0).toUpperCase() + J.slice(1)),
  y6 = vJ((J) => (J ? `on${sZ(J)}` : "")),
  oZ = (J, Y) => J !== Y && (J === J || Y === Y),
  _Y = new WeakMap(),
  zJ = [],
  h,
  r = Symbol("iterate"),
  zY = Symbol("Map key iterate");
function b1(J) {
  return J && J._isEffect === !0;
}
function S1(J, Y = E1) {
  if (b1(J)) J = J.raw;
  let Z = h1(J, Y);
  if (!Y.lazy) Z();
  return Z;
}
function x1(J) {
  if (J.active) {
    if ((nZ(J), J.options.onStop)) J.options.onStop();
    J.active = !1;
  }
}
var y1 = 0;
function h1(J, Y) {
  let Z = function $() {
    if (!Z.active) return J();
    if (!zJ.includes(Z)) {
      nZ(Z);
      try {
        return v1(), zJ.push(Z), (h = Z), J();
      } finally {
        zJ.pop(), lZ(), (h = zJ[zJ.length - 1]);
      }
    }
  };
  return (
    (Z.id = y1++),
    (Z.allowRecurse = !!Y.allowRecurse),
    (Z._isEffect = !0),
    (Z.active = !0),
    (Z.raw = J),
    (Z.deps = []),
    (Z.options = Y),
    Z
  );
}
function nZ(J) {
  let { deps: Y } = J;
  if (Y.length) {
    for (let Z = 0; Z < Y.length; Z++) Y[Z].delete(J);
    Y.length = 0;
  }
}
var ZJ = !0,
  wY = [];
function f1() {
  wY.push(ZJ), (ZJ = !1);
}
function v1() {
  wY.push(ZJ), (ZJ = !0);
}
function lZ() {
  let J = wY.pop();
  ZJ = J === void 0 ? !0 : J;
}
function x(J, Y, Z) {
  if (!ZJ || h === void 0) return;
  let $ = _Y.get(J);
  if (!$) _Y.set(J, ($ = new Map()));
  let W = $.get(Z);
  if (!W) $.set(Z, (W = new Set()));
  if (!W.has(h)) {
    if ((W.add(h), h.deps.push(W), h.options.onTrack))
      h.options.onTrack({ effect: h, target: J, type: Y, key: Z });
  }
}
function m(J, Y, Z, $, W, X) {
  let Q = _Y.get(J);
  if (!Q) return;
  let B = new Set(),
    U = (G) => {
      if (G)
        G.forEach((z) => {
          if (z !== h || z.allowRecurse) B.add(z);
        });
    };
  if (Y === "clear") Q.forEach(U);
  else if (Z === "length" && a(J))
    Q.forEach((G, z) => {
      if (z === "length" || z >= $) U(G);
    });
  else {
    if (Z !== void 0) U(Q.get(Z));
    switch (Y) {
      case "add":
        if (!a(J)) {
          if ((U(Q.get(r)), AJ(J))) U(Q.get(zY));
        } else if (DY(Z)) U(Q.get("length"));
        break;
      case "delete":
        if (!a(J)) {
          if ((U(Q.get(r)), AJ(J))) U(Q.get(zY));
        }
        break;
      case "set":
        if (AJ(J)) U(Q.get(r));
        break;
    }
  }
  let _ = (G) => {
    if (G.options.onTrigger)
      G.options.onTrigger({
        effect: G,
        target: J,
        key: Z,
        type: Y,
        newValue: $,
        oldValue: W,
        oldTarget: X,
      });
    if (G.options.scheduler) G.options.scheduler(G);
    else G();
  };
  B.forEach(_);
}
var u1 = mZ("__proto__,__v_isRef,__isVue"),
  aZ = new Set(
    Object.getOwnPropertyNames(Symbol)
      .map((J) => Symbol[J])
      .filter(IY),
  ),
  g1 = rZ(),
  c1 = rZ(!0),
  dY = d1();
function d1() {
  let J = {};
  return (
    ["includes", "indexOf", "lastIndexOf"].forEach((Y) => {
      J[Y] = function (...Z) {
        let $ = R(this);
        for (let X = 0, Q = this.length; X < Q; X++) x($, "get", X + "");
        let W = $[Y](...Z);
        if (W === -1 || W === !1) return $[Y](...Z.map(R));
        else return W;
      };
    }),
    ["push", "pop", "shift", "unshift", "splice"].forEach((Y) => {
      J[Y] = function (...Z) {
        f1();
        let $ = R(this)[Y].apply(this, Z);
        return lZ(), $;
      };
    }),
    J
  );
}
function rZ(J = !1, Y = !1) {
  return function Z($, W, X) {
    if (W === "__v_isReactive") return !J;
    else if (W === "__v_isReadonly") return J;
    else if (W === "__v_raw" && X === (J ? (Y ? W6 : Y0) : Y ? $6 : J0).get($))
      return $;
    let Q = a($);
    if (!J && Q && hJ(dY, W)) return Reflect.get(dY, W, X);
    let B = Reflect.get($, W, X);
    if (IY(W) ? aZ.has(W) : u1(W)) return B;
    if (!J) x($, "get", W);
    if (Y) return B;
    if (KY(B)) return !Q || !DY(W) ? B.value : B;
    if (fJ(B)) return J ? Z0(B) : xY(B);
    return B;
  };
}
var m1 = p1();
function p1(J = !1) {
  return function Y(Z, $, W, X) {
    let Q = Z[$];
    if (!J) {
      if (((W = R(W)), (Q = R(Q)), !a(Z) && KY(Q) && !KY(W)))
        return (Q.value = W), !0;
    }
    let B = a(Z) && DY($) ? Number($) < Z.length : hJ(Z, $),
      U = Reflect.set(Z, $, W, X);
    if (Z === R(X)) {
      if (!B) m(Z, "add", $, W);
      else if (oZ(W, Q)) m(Z, "set", $, W, Q);
    }
    return U;
  };
}
function i1(J, Y) {
  let Z = hJ(J, Y),
    $ = J[Y],
    W = Reflect.deleteProperty(J, Y);
  if (W && Z) m(J, "delete", Y, void 0, $);
  return W;
}
function s1(J, Y) {
  let Z = Reflect.has(J, Y);
  if (!IY(Y) || !aZ.has(Y)) x(J, "has", Y);
  return Z;
}
function o1(J) {
  return x(J, "iterate", a(J) ? "length" : r), Reflect.ownKeys(J);
}
var n1 = { get: g1, set: m1, deleteProperty: i1, has: s1, ownKeys: o1 },
  l1 = {
    get: c1,
    set(J, Y) {
      return (
        console.warn(
          `Set operation on key "${String(Y)}" failed: target is readonly.`,
          J,
        ),
        !0
      );
    },
    deleteProperty(J, Y) {
      return (
        console.warn(
          `Delete operation on key "${String(Y)}" failed: target is readonly.`,
          J,
        ),
        !0
      );
    },
  },
  kY = (J) => (fJ(J) ? xY(J) : J),
  bY = (J) => (fJ(J) ? Z0(J) : J),
  SY = (J) => J,
  uJ = (J) => Reflect.getPrototypeOf(J);
function RJ(J, Y, Z = !1, $ = !1) {
  J = J.__v_raw;
  let W = R(J),
    X = R(Y);
  if (Y !== X) !Z && x(W, "get", Y);
  !Z && x(W, "get", X);
  let { has: Q } = uJ(W),
    B = $ ? SY : Z ? bY : kY;
  if (Q.call(W, Y)) return B(J.get(Y));
  else if (Q.call(W, X)) return B(J.get(X));
  else if (J !== W) J.get(Y);
}
function OJ(J, Y = !1) {
  let Z = this.__v_raw,
    $ = R(Z),
    W = R(J);
  if (J !== W) !Y && x($, "has", J);
  return !Y && x($, "has", W), J === W ? Z.has(J) : Z.has(J) || Z.has(W);
}
function TJ(J, Y = !1) {
  return (
    (J = J.__v_raw), !Y && x(R(J), "iterate", r), Reflect.get(J, "size", J)
  );
}
function mY(J) {
  J = R(J);
  let Y = R(this);
  if (!uJ(Y).has.call(Y, J)) Y.add(J), m(Y, "add", J, J);
  return this;
}
function pY(J, Y) {
  Y = R(Y);
  let Z = R(this),
    { has: $, get: W } = uJ(Z),
    X = $.call(Z, J);
  if (!X) (J = R(J)), (X = $.call(Z, J));
  else eZ(Z, $, J);
  let Q = W.call(Z, J);
  if ((Z.set(J, Y), !X)) m(Z, "add", J, Y);
  else if (oZ(Y, Q)) m(Z, "set", J, Y, Q);
  return this;
}
function iY(J) {
  let Y = R(this),
    { has: Z, get: $ } = uJ(Y),
    W = Z.call(Y, J);
  if (!W) (J = R(J)), (W = Z.call(Y, J));
  else eZ(Y, Z, J);
  let X = $ ? $.call(Y, J) : void 0,
    Q = Y.delete(J);
  if (W) m(Y, "delete", J, void 0, X);
  return Q;
}
function sY() {
  let J = R(this),
    Y = J.size !== 0,
    Z = AJ(J) ? new Map(J) : new Set(J),
    $ = J.clear();
  if (Y) m(J, "clear", void 0, void 0, Z);
  return $;
}
function VJ(J, Y) {
  return function Z($, W) {
    let X = this,
      Q = X.__v_raw,
      B = R(Q),
      U = Y ? SY : J ? bY : kY;
    return (
      !J && x(B, "iterate", r),
      Q.forEach((_, G) => {
        return $.call(W, U(_), U(G), X);
      })
    );
  };
}
function EJ(J, Y, Z) {
  return function (...$) {
    let W = this.__v_raw,
      X = R(W),
      Q = AJ(X),
      B = J === "entries" || (J === Symbol.iterator && Q),
      U = J === "keys" && Q,
      _ = W[J](...$),
      G = Z ? SY : Y ? bY : kY;
    return (
      !Y && x(X, "iterate", U ? zY : r),
      {
        next() {
          let { value: z, done: K } = _.next();
          return K
            ? { value: z, done: K }
            : { value: B ? [G(z[0]), G(z[1])] : G(z), done: K };
        },
        [Symbol.iterator]() {
          return this;
        },
      }
    );
  };
}
function c(J) {
  return function (...Y) {
    {
      let Z = Y[0] ? `on key "${Y[0]}" ` : "";
      console.warn(
        `${sZ(J)} operation ${Z}failed: target is readonly.`,
        R(this),
      );
    }
    return J === "delete" ? !1 : this;
  };
}
function a1() {
  let J = {
      get(X) {
        return RJ(this, X);
      },
      get size() {
        return TJ(this);
      },
      has: OJ,
      add: mY,
      set: pY,
      delete: iY,
      clear: sY,
      forEach: VJ(!1, !1),
    },
    Y = {
      get(X) {
        return RJ(this, X, !1, !0);
      },
      get size() {
        return TJ(this);
      },
      has: OJ,
      add: mY,
      set: pY,
      delete: iY,
      clear: sY,
      forEach: VJ(!1, !0),
    },
    Z = {
      get(X) {
        return RJ(this, X, !0);
      },
      get size() {
        return TJ(this, !0);
      },
      has(X) {
        return OJ.call(this, X, !0);
      },
      add: c("add"),
      set: c("set"),
      delete: c("delete"),
      clear: c("clear"),
      forEach: VJ(!0, !1),
    },
    $ = {
      get(X) {
        return RJ(this, X, !0, !0);
      },
      get size() {
        return TJ(this, !0);
      },
      has(X) {
        return OJ.call(this, X, !0);
      },
      add: c("add"),
      set: c("set"),
      delete: c("delete"),
      clear: c("clear"),
      forEach: VJ(!0, !0),
    };
  return (
    ["keys", "values", "entries", Symbol.iterator].forEach((X) => {
      (J[X] = EJ(X, !1, !1)),
        (Z[X] = EJ(X, !0, !1)),
        (Y[X] = EJ(X, !1, !0)),
        ($[X] = EJ(X, !0, !0));
    }),
    [J, Z, Y, $]
  );
}
var [r1, t1, e1, J6] = a1();
function tZ(J, Y) {
  let Z = Y ? (J ? J6 : e1) : J ? t1 : r1;
  return ($, W, X) => {
    if (W === "__v_isReactive") return !J;
    else if (W === "__v_isReadonly") return J;
    else if (W === "__v_raw") return $;
    return Reflect.get(hJ(Z, W) && W in $ ? Z : $, W, X);
  };
}
var Y6 = { get: tZ(!1, !1) },
  Z6 = { get: tZ(!0, !1) };
function eZ(J, Y, Z) {
  let $ = R(Z);
  if ($ !== Z && Y.call(J, $)) {
    let W = iZ(J);
    console.warn(
      `Reactive ${W} contains both the raw and reactive versions of the same object${W === "Map" ? " as keys" : ""}, which can lead to inconsistencies. Avoid differentiating between the raw and reactive versions of an object and only use the reactive version if possible.`,
    );
  }
}
var J0 = new WeakMap(),
  $6 = new WeakMap(),
  Y0 = new WeakMap(),
  W6 = new WeakMap();
function X6(J) {
  switch (J) {
    case "Object":
    case "Array":
      return 1;
    case "Map":
    case "Set":
    case "WeakMap":
    case "WeakSet":
      return 2;
    default:
      return 0;
  }
}
function Q6(J) {
  return J.__v_skip || !Object.isExtensible(J) ? 0 : X6(iZ(J));
}
function xY(J) {
  if (J && J.__v_isReadonly) return J;
  return $0(J, !1, n1, Y6, J0);
}
function Z0(J) {
  return $0(J, !0, l1, Z6, Y0);
}
function $0(J, Y, Z, $, W) {
  if (!fJ(J))
    return console.warn(`value cannot be made reactive: ${String(J)}`), J;
  if (J.__v_raw && !(Y && J.__v_isReactive)) return J;
  let X = W.get(J);
  if (X) return X;
  let Q = Q6(J);
  if (Q === 0) return J;
  let B = new Proxy(J, Q === 2 ? $ : Z);
  return W.set(J, B), B;
}
function R(J) {
  return (J && R(J.__v_raw)) || J;
}
function KY(J) {
  return Boolean(J && J.__v_isRef === !0);
}
y("nextTick", () => VY);
y("dispatch", (J) => MJ.bind(MJ, J));
y("watch", (J, { evaluateLater: Y, cleanup: Z }) => ($, W) => {
  let X = Y($),
    B = tY(() => {
      let U;
      return X((_) => (U = _)), U;
    }, W);
  Z(B);
});
y("store", C1);
y("data", (J) => QZ(J));
y("root", (J) => SJ(J));
y("refs", (J) => {
  if (J._x_refs_proxy) return J._x_refs_proxy;
  return (J._x_refs_proxy = qJ(B6(J))), J._x_refs_proxy;
});
function B6(J) {
  let Y = [];
  return (
    QJ(J, (Z) => {
      if (Z._x_refs) Y.push(Z._x_refs);
    }),
    Y
  );
}
var nJ = {};
function W0(J) {
  if (!nJ[J]) nJ[J] = 0;
  return ++nJ[J];
}
function U6(J, Y) {
  return QJ(J, (Z) => {
    if (Z._x_ids && Z._x_ids[Y]) return !0;
  });
}
function G6(J, Y) {
  if (!J._x_ids) J._x_ids = {};
  if (!J._x_ids[Y]) J._x_ids[Y] = W0(Y);
}
y("id", (J, { cleanup: Y }) => (Z, $ = null) => {
  let W = `${Z}${$ ? `-${$}` : ""}`;
  return _6(J, W, Y, () => {
    let X = U6(J, Z),
      Q = X ? X._x_ids[Z] : W0(Z);
    return $ ? `${Z}-${Q}-${$}` : `${Z}-${Q}`;
  });
});
yJ((J, Y) => {
  if (J._x_id) Y._x_id = J._x_id;
});
function _6(J, Y, Z, $) {
  if (!J._x_id) J._x_id = {};
  if (J._x_id[Y]) return J._x_id[Y];
  let W = $();
  return (
    (J._x_id[Y] = W),
    Z(() => {
      delete J._x_id[Y];
    }),
    W
  );
}
y("el", (J) => J);
X0("Focus", "focus", "focus");
X0("Persist", "persist", "persist");
function X0(J, Y, Z) {
  y(Y, ($) =>
    S(
      `You can't use [$${Y}] without first installing the "${J}" plugin here: https://alpinejs.dev/plugins/${Z}`,
      $,
    ),
  );
}
T(
  "modelable",
  (J, { expression: Y }, { effect: Z, evaluateLater: $, cleanup: W }) => {
    let X = $(Y),
      Q = () => {
        let G;
        return X((z) => (G = z)), G;
      },
      B = $(`${Y} = __placeholder`),
      U = (G) => B(() => {}, { scope: { __placeholder: G } }),
      _ = Q();
    U(_),
      queueMicrotask(() => {
        if (!J._x_model) return;
        J._x_removeModelListeners.default();
        let G = J._x_model.get,
          z = J._x_model.set,
          K = uZ(
            {
              get() {
                return G();
              },
              set(A) {
                z(A);
              },
            },
            {
              get() {
                return Q();
              },
              set(A) {
                U(A);
              },
            },
          );
        W(K);
      });
  },
);
T("teleport", (J, { modifiers: Y, expression: Z }, { cleanup: $ }) => {
  if (J.tagName.toLowerCase() !== "template")
    S("x-teleport can only be used on a <template> tag", J);
  let W = oY(Z),
    X = J.content.cloneNode(!0).firstElementChild;
  if (
    ((J._x_teleport = X),
    (X._x_teleportBack = J),
    J.setAttribute("data-teleport-template", !0),
    X.setAttribute("data-teleport-target", !0),
    J._x_forwardEvents)
  )
    J._x_forwardEvents.forEach((B) => {
      X.addEventListener(B, (U) => {
        U.stopPropagation(), J.dispatchEvent(new U.constructor(U.type, U));
      });
    });
  PJ(X, {}, J);
  let Q = (B, U, _) => {
    if (_.includes("prepend")) U.parentNode.insertBefore(B, U);
    else if (_.includes("append")) U.parentNode.insertBefore(B, U.nextSibling);
    else U.appendChild(B);
  };
  O(() => {
    Q(X, W, Y),
      p(() => {
        u(X);
      })();
  }),
    (J._x_teleportPutBack = () => {
      let B = oY(Z);
      O(() => {
        Q(J._x_teleport, B, Y);
      });
    }),
    $(() =>
      O(() => {
        X.remove(), BJ(X);
      }),
    );
});
var z6 = document.createElement("div");
function oY(J) {
  let Y = p(
    () => {
      return document.querySelector(J);
    },
    () => {
      return z6;
    },
  )();
  if (!Y) S(`Cannot find x-teleport element for selector: "${J}"`);
  return Y;
}
var Q0 = () => {};
Q0.inline = (J, { modifiers: Y }, { cleanup: Z }) => {
  Y.includes("self") ? (J._x_ignoreSelf = !0) : (J._x_ignore = !0),
    Z(() => {
      Y.includes("self") ? delete J._x_ignoreSelf : delete J._x_ignore;
    });
};
T("ignore", Q0);
T(
  "effect",
  p((J, { expression: Y }, { effect: Z }) => {
    Z(w(J, Y));
  }),
);
function MY(J, Y, Z, $) {
  let W = J,
    X = (U) => $(U),
    Q = {},
    B = (U, _) => (G) => _(U, G);
  if (Z.includes("dot")) Y = K6(Y);
  if (Z.includes("camel")) Y = M6(Y);
  if (Z.includes("passive")) Q.passive = !0;
  if (Z.includes("capture")) Q.capture = !0;
  if (Z.includes("window")) W = window;
  if (Z.includes("document")) W = document;
  if (Z.includes("debounce")) {
    let U = Z[Z.indexOf("debounce") + 1] || "invalid-wait",
      _ = bJ(U.split("ms")[0]) ? Number(U.split("ms")[0]) : 250;
    X = fZ(X, _);
  }
  if (Z.includes("throttle")) {
    let U = Z[Z.indexOf("throttle") + 1] || "invalid-wait",
      _ = bJ(U.split("ms")[0]) ? Number(U.split("ms")[0]) : 250;
    X = vZ(X, _);
  }
  if (Z.includes("prevent"))
    X = B(X, (U, _) => {
      _.preventDefault(), U(_);
    });
  if (Z.includes("stop"))
    X = B(X, (U, _) => {
      _.stopPropagation(), U(_);
    });
  if (Z.includes("once"))
    X = B(X, (U, _) => {
      U(_), W.removeEventListener(Y, X, Q);
    });
  if (Z.includes("away") || Z.includes("outside"))
    (W = document),
      (X = B(X, (U, _) => {
        if (J.contains(_.target)) return;
        if (_.target.isConnected === !1) return;
        if (J.offsetWidth < 1 && J.offsetHeight < 1) return;
        if (J._x_isShown === !1) return;
        U(_);
      }));
  if (Z.includes("self"))
    X = B(X, (U, _) => {
      _.target === J && U(_);
    });
  if (L6(Y) || B0(Y))
    X = B(X, (U, _) => {
      if (P6(_, Z)) return;
      U(_);
    });
  return (
    W.addEventListener(Y, X, Q),
    () => {
      W.removeEventListener(Y, X, Q);
    }
  );
}
function K6(J) {
  return J.replace(/-/g, ".");
}
function M6(J) {
  return J.toLowerCase().replace(/-(\w)/g, (Y, Z) => Z.toUpperCase());
}
function bJ(J) {
  return !Array.isArray(J) && !isNaN(J);
}
function A6(J) {
  if ([" ", "_"].includes(J)) return J;
  return J.replace(/([a-z])([A-Z])/g, "$1-$2")
    .replace(/[_\s]/, "-")
    .toLowerCase();
}
function L6(J) {
  return ["keydown", "keyup"].includes(J);
}
function B0(J) {
  return ["contextmenu", "click", "mouse"].some((Y) => J.includes(Y));
}
function P6(J, Y) {
  let Z = Y.filter((X) => {
    return ![
      "window",
      "document",
      "prevent",
      "stop",
      "once",
      "capture",
      "self",
      "away",
      "outside",
      "passive",
    ].includes(X);
  });
  if (Z.includes("debounce")) {
    let X = Z.indexOf("debounce");
    Z.splice(X, bJ((Z[X + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
  }
  if (Z.includes("throttle")) {
    let X = Z.indexOf("throttle");
    Z.splice(X, bJ((Z[X + 1] || "invalid-wait").split("ms")[0]) ? 2 : 1);
  }
  if (Z.length === 0) return !1;
  if (Z.length === 1 && nY(J.key).includes(Z[0])) return !1;
  let W = ["ctrl", "shift", "alt", "meta", "cmd", "super"].filter((X) =>
    Z.includes(X),
  );
  if (((Z = Z.filter((X) => !W.includes(X))), W.length > 0)) {
    if (
      W.filter((Q) => {
        if (Q === "cmd" || Q === "super") Q = "meta";
        return J[`${Q}Key`];
      }).length === W.length
    ) {
      if (B0(J.type)) return !1;
      if (nY(J.key).includes(Z[0])) return !1;
    }
  }
  return !0;
}
function nY(J) {
  if (!J) return [];
  J = A6(J);
  let Y = {
    ctrl: "control",
    slash: "/",
    space: " ",
    spacebar: " ",
    cmd: "meta",
    esc: "escape",
    up: "arrow-up",
    down: "arrow-down",
    left: "arrow-left",
    right: "arrow-right",
    period: ".",
    comma: ",",
    equal: "=",
    minus: "-",
    underscore: "_",
  };
  return (
    (Y[J] = J),
    Object.keys(Y)
      .map((Z) => {
        if (Y[Z] === J) return Z;
      })
      .filter((Z) => Z)
  );
}
T("model", (J, { modifiers: Y, expression: Z }, { effect: $, cleanup: W }) => {
  let X = J;
  if (Y.includes("parent")) X = J.parentNode;
  let Q = w(X, Z),
    B;
  if (typeof Z === "string") B = w(X, `${Z} = __placeholder`);
  else if (typeof Z === "function" && typeof Z() === "string")
    B = w(X, `${Z()} = __placeholder`);
  else B = () => {};
  let U = () => {
      let K;
      return Q((A) => (K = A)), lY(K) ? K.get() : K;
    },
    _ = (K) => {
      let A;
      if ((Q((L) => (A = L)), lY(A))) A.set(K);
      else B(() => {}, { scope: { __placeholder: K } });
    };
  if (typeof Z === "string" && J.type === "radio")
    O(() => {
      if (!J.hasAttribute("name")) J.setAttribute("name", Z);
    });
  var G =
    J.tagName.toLowerCase() === "select" ||
    ["checkbox", "radio"].includes(J.type) ||
    Y.includes("lazy")
      ? "change"
      : "input";
  let z = d
    ? () => {}
    : MY(J, G, Y, (K) => {
        _(lJ(J, Y, K, U()));
      });
  if (Y.includes("fill")) {
    if (
      [void 0, null, ""].includes(U()) ||
      (NY(J) && Array.isArray(U())) ||
      (J.tagName.toLowerCase() === "select" && J.multiple)
    )
      _(lJ(J, Y, { target: J }, U()));
  }
  if (!J._x_removeModelListeners) J._x_removeModelListeners = {};
  if (
    ((J._x_removeModelListeners.default = z),
    W(() => J._x_removeModelListeners.default()),
    J.form)
  ) {
    let K = MY(J.form, "reset", [], (A) => {
      VY(() => J._x_model && J._x_model.set(lJ(J, Y, { target: J }, U())));
    });
    W(() => K());
  }
  (J._x_model = {
    get() {
      return U();
    },
    set(K) {
      _(K);
    },
  }),
    (J._x_forceModelUpdate = (K) => {
      if (K === void 0 && typeof Z === "string" && Z.match(/\./)) K = "";
      (window.fromModel = !0),
        O(() => bZ(J, "value", K)),
        delete window.fromModel;
    }),
    $(() => {
      let K = U();
      if (Y.includes("unintrusive") && document.activeElement.isSameNode(J))
        return;
      J._x_forceModelUpdate(K);
    });
});
function lJ(J, Y, Z, $) {
  return O(() => {
    if (Z instanceof CustomEvent && Z.detail !== void 0)
      return Z.detail !== null && Z.detail !== void 0
        ? Z.detail
        : Z.target.value;
    else if (NY(J))
      if (Array.isArray($)) {
        let W = null;
        if (Y.includes("number")) W = aJ(Z.target.value);
        else if (Y.includes("boolean")) W = IJ(Z.target.value);
        else W = Z.target.value;
        return Z.target.checked
          ? $.includes(W)
            ? $
            : $.concat([W])
          : $.filter((X) => !q6(X, W));
      } else return Z.target.checked;
    else if (J.tagName.toLowerCase() === "select" && J.multiple) {
      if (Y.includes("number"))
        return Array.from(Z.target.selectedOptions).map((W) => {
          let X = W.value || W.text;
          return aJ(X);
        });
      else if (Y.includes("boolean"))
        return Array.from(Z.target.selectedOptions).map((W) => {
          let X = W.value || W.text;
          return IJ(X);
        });
      return Array.from(Z.target.selectedOptions).map((W) => {
        return W.value || W.text;
      });
    } else {
      let W;
      if (hZ(J))
        if (Z.target.checked) W = Z.target.value;
        else W = $;
      else W = Z.target.value;
      if (Y.includes("number")) return aJ(W);
      else if (Y.includes("boolean")) return IJ(W);
      else if (Y.includes("trim")) return W.trim();
      else return W;
    }
  });
}
function aJ(J) {
  let Y = J ? parseFloat(J) : null;
  return H6(Y) ? Y : J;
}
function q6(J, Y) {
  return J == Y;
}
function H6(J) {
  return !Array.isArray(J) && !isNaN(J);
}
function lY(J) {
  return (
    J !== null &&
    typeof J === "object" &&
    typeof J.get === "function" &&
    typeof J.set === "function"
  );
}
T("cloak", (J) =>
  queueMicrotask(() => O(() => J.removeAttribute(XJ("cloak")))),
);
EZ(() => `[${XJ("init")}]`);
T(
  "init",
  p((J, { expression: Y }, { evaluate: Z }) => {
    if (typeof Y === "string") return !!Y.trim() && Z(Y, {}, !1);
    return Z(Y, {}, !1);
  }),
);
T("text", (J, { expression: Y }, { effect: Z, evaluateLater: $ }) => {
  let W = $(Y);
  Z(() => {
    W((X) => {
      O(() => {
        J.textContent = X;
      });
    });
  });
});
T("html", (J, { expression: Y }, { effect: Z, evaluateLater: $ }) => {
  let W = $(Y);
  Z(() => {
    W((X) => {
      O(() => {
        (J.innerHTML = X), (J._x_ignoreSelf = !0), u(J), delete J._x_ignoreSelf;
      });
    });
  });
});
RY(PZ(":", qZ(XJ("bind:"))));
var U0 = (
  J,
  { value: Y, modifiers: Z, expression: $, original: W },
  { effect: X, cleanup: Q },
) => {
  if (!Y) {
    let U = {};
    j1(U),
      w(J, $)(
        (G) => {
          cZ(J, G, W);
        },
        { scope: U },
      );
    return;
  }
  if (Y === "key") return C6(J, $);
  if (
    J._x_inlineBindings &&
    J._x_inlineBindings[Y] &&
    J._x_inlineBindings[Y].extract
  )
    return;
  let B = w(J, $);
  X(() =>
    B((U) => {
      if (U === void 0 && typeof $ === "string" && $.match(/\./)) U = "";
      O(() => bZ(J, Y, U, Z));
    }),
  ),
    Q(() => {
      J._x_undoAddedClasses && J._x_undoAddedClasses(),
        J._x_undoAddedStyles && J._x_undoAddedStyles();
    });
};
U0.inline = (J, { value: Y, modifiers: Z, expression: $ }) => {
  if (!Y) return;
  if (!J._x_inlineBindings) J._x_inlineBindings = {};
  J._x_inlineBindings[Y] = { expression: $, extract: !1 };
};
T("bind", U0);
function C6(J, Y) {
  J._x_keyExpression = Y;
}
VZ(() => `[${XJ("data")}]`);
T("data", (J, { expression: Y }, { cleanup: Z }) => {
  if (F6(J)) return;
  Y = Y === "" ? "{}" : Y;
  let $ = {};
  ZY($, J);
  let W = {};
  O1(W, $);
  let X = l(J, Y, { scope: W });
  if (X === void 0 || X === !0) X = {};
  ZY(X, J);
  let Q = $J(X);
  BZ(Q);
  let B = PJ(J, Q);
  Q.init && l(J, Q.init),
    Z(() => {
      Q.destroy && l(J, Q.destroy), B();
    });
});
yJ((J, Y) => {
  if (J._x_dataStack)
    (Y._x_dataStack = J._x_dataStack),
      Y.setAttribute("data-has-alpine-state", !0);
});
function F6(J) {
  if (!d) return !1;
  if (GY) return !0;
  return J.hasAttribute("data-has-alpine-state");
}
T("show", (J, { modifiers: Y, expression: Z }, { effect: $ }) => {
  let W = w(J, Z);
  if (!J._x_doHide)
    J._x_doHide = () => {
      O(() => {
        J.style.setProperty(
          "display",
          "none",
          Y.includes("important") ? "important" : void 0,
        );
      });
    };
  if (!J._x_doShow)
    J._x_doShow = () => {
      O(() => {
        if (J.style.length === 1 && J.style.display === "none")
          J.removeAttribute("style");
        else J.style.removeProperty("display");
      });
    };
  let X = () => {
      J._x_doHide(), (J._x_isShown = !1);
    },
    Q = () => {
      J._x_doShow(), (J._x_isShown = !0);
    },
    B = () => setTimeout(Q),
    U = BY(
      (z) => (z ? Q() : X()),
      (z) => {
        if (typeof J._x_toggleAndCascadeWithTransitions === "function")
          J._x_toggleAndCascadeWithTransitions(J, z, Q, X);
        else z ? B() : X();
      },
    ),
    _,
    G = !0;
  $(() =>
    W((z) => {
      if (!G && z === _) return;
      if (Y.includes("immediate")) z ? B() : X();
      U(z), (_ = z), (G = !1);
    }),
  );
});
T("for", (J, { expression: Y }, { effect: Z, cleanup: $ }) => {
  let W = R6(Y),
    X = w(J, W.items),
    Q = w(J, J._x_keyExpression || "index");
  (J._x_prevKeys = []),
    (J._x_lookup = {}),
    Z(() => j6(J, W, X, Q)),
    $(() => {
      Object.values(J._x_lookup).forEach((B) =>
        O(() => {
          BJ(B), B.remove();
        }),
      ),
        delete J._x_prevKeys,
        delete J._x_lookup;
    });
});
function j6(J, Y, Z, $) {
  let W = (Q) => typeof Q === "object" && !Array.isArray(Q),
    X = J;
  Z((Q) => {
    if (O6(Q) && Q >= 0) Q = Array.from(Array(Q).keys(), (M) => M + 1);
    if (Q === void 0) Q = [];
    let { _x_lookup: B, _x_prevKeys: U } = J,
      _ = [],
      G = [];
    if (W(Q))
      Q = Object.entries(Q).map(([M, P]) => {
        let C = aY(Y, P, M, Q);
        $(
          (q) => {
            if (G.includes(q)) S("Duplicate key on x-for", J);
            G.push(q);
          },
          { scope: { index: M, ...C } },
        ),
          _.push(C);
      });
    else
      for (let M = 0; M < Q.length; M++) {
        let P = aY(Y, Q[M], M, Q);
        $(
          (C) => {
            if (G.includes(C)) S("Duplicate key on x-for", J);
            G.push(C);
          },
          { scope: { index: M, ...P } },
        ),
          _.push(P);
      }
    let z = [],
      K = [],
      A = [],
      L = [];
    for (let M = 0; M < U.length; M++) {
      let P = U[M];
      if (G.indexOf(P) === -1) A.push(P);
    }
    U = U.filter((M) => !A.includes(M));
    let H = "template";
    for (let M = 0; M < G.length; M++) {
      let P = G[M],
        C = U.indexOf(P);
      if (C === -1) U.splice(M, 0, P), z.push([H, M]);
      else if (C !== M) {
        let q = U.splice(M, 1)[0],
          j = U.splice(C - 1, 1)[0];
        U.splice(M, 0, j), U.splice(C, 0, q), K.push([q, j]);
      } else L.push(P);
      H = P;
    }
    for (let M = 0; M < A.length; M++) {
      let P = A[M];
      if (!(P in B)) continue;
      O(() => {
        BJ(B[P]), B[P].remove();
      }),
        delete B[P];
    }
    for (let M = 0; M < K.length; M++) {
      let [P, C] = K[M],
        q = B[P],
        j = B[C],
        N = document.createElement("div");
      O(() => {
        if (!j) S('x-for ":key" is undefined or invalid', X, C, B);
        j.after(N),
          q.after(j),
          j._x_currentIfEl && j.after(j._x_currentIfEl),
          N.before(q),
          q._x_currentIfEl && q.after(q._x_currentIfEl),
          N.remove();
      }),
        j._x_refreshXForScope(_[G.indexOf(C)]);
    }
    for (let M = 0; M < z.length; M++) {
      let [P, C] = z[M],
        q = P === "template" ? X : B[P];
      if (q._x_currentIfEl) q = q._x_currentIfEl;
      let j = _[C],
        N = G[C],
        E = document.importNode(X.content, !0).firstElementChild,
        k = $J(j);
      if (
        (PJ(E, k, X),
        (E._x_refreshXForScope = (I) => {
          Object.entries(I).forEach(([g, JJ]) => {
            k[g] = JJ;
          });
        }),
        O(() => {
          q.after(E), p(() => u(E))();
        }),
        typeof N === "object")
      )
        S(
          "x-for key cannot be an object, it must be a string or an integer",
          X,
        );
      B[N] = E;
    }
    for (let M = 0; M < L.length; M++)
      B[L[M]]._x_refreshXForScope(_[G.indexOf(L[M])]);
    X._x_prevKeys = G;
  });
}
function R6(J) {
  let Y = /,([^,\}\]]*)(?:,([^,\}\]]*))?$/,
    Z = /^\s*\(|\)\s*$/g,
    $ = /([\s\S]*?)\s+(?:in|of)\s+([\s\S]*)/,
    W = J.match($);
  if (!W) return;
  let X = {};
  X.items = W[2].trim();
  let Q = W[1].replace(Z, "").trim(),
    B = Q.match(Y);
  if (B) {
    if (((X.item = Q.replace(Y, "").trim()), (X.index = B[1].trim()), B[2]))
      X.collection = B[2].trim();
  } else X.item = Q;
  return X;
}
function aY(J, Y, Z, $) {
  let W = {};
  if (/^\[.*\]$/.test(J.item) && Array.isArray(Y))
    J.item
      .replace("[", "")
      .replace("]", "")
      .split(",")
      .map((Q) => Q.trim())
      .forEach((Q, B) => {
        W[Q] = Y[B];
      });
  else if (
    /^\{.*\}$/.test(J.item) &&
    !Array.isArray(Y) &&
    typeof Y === "object"
  )
    J.item
      .replace("{", "")
      .replace("}", "")
      .split(",")
      .map((Q) => Q.trim())
      .forEach((Q) => {
        W[Q] = Y[Q];
      });
  else W[J.item] = Y;
  if (J.index) W[J.index] = Z;
  if (J.collection) W[J.collection] = $;
  return W;
}
function O6(J) {
  return !Array.isArray(J) && !isNaN(J);
}
function G0() {}
G0.inline = (J, { expression: Y }, { cleanup: Z }) => {
  let $ = SJ(J);
  if (!$._x_refs) $._x_refs = {};
  ($._x_refs[Y] = J), Z(() => delete $._x_refs[Y]);
};
T("ref", G0);
T("if", (J, { expression: Y }, { effect: Z, cleanup: $ }) => {
  if (J.tagName.toLowerCase() !== "template")
    S("x-if can only be used on a <template> tag", J);
  let W = w(J, Y),
    X = () => {
      if (J._x_currentIfEl) return J._x_currentIfEl;
      let B = J.content.cloneNode(!0).firstElementChild;
      return (
        PJ(B, {}, J),
        O(() => {
          J.after(B), p(() => u(B))();
        }),
        (J._x_currentIfEl = B),
        (J._x_undoIf = () => {
          O(() => {
            BJ(B), B.remove();
          }),
            delete J._x_currentIfEl;
        }),
        B
      );
    },
    Q = () => {
      if (!J._x_undoIf) return;
      J._x_undoIf(), delete J._x_undoIf;
    };
  Z(() =>
    W((B) => {
      B ? X() : Q();
    }),
  ),
    $(() => J._x_undoIf && J._x_undoIf());
});
T("id", (J, { expression: Y }, { evaluate: Z }) => {
  Z(Y).forEach((W) => G6(J, W));
});
yJ((J, Y) => {
  if (J._x_ids) Y._x_ids = J._x_ids;
});
RY(PZ("@", qZ(XJ("on:"))));
T(
  "on",
  p((J, { value: Y, modifiers: Z, expression: $ }, { cleanup: W }) => {
    let X = $ ? w(J, $) : () => {};
    if (J.tagName.toLowerCase() === "template") {
      if (!J._x_forwardEvents) J._x_forwardEvents = [];
      if (!J._x_forwardEvents.includes(Y)) J._x_forwardEvents.push(Y);
    }
    let Q = MY(J, Y, Z, (B) => {
      X(() => {}, { scope: { $event: B }, params: [B] });
    });
    W(() => Q());
  }),
);
gJ("Collapse", "collapse", "collapse");
gJ("Intersect", "intersect", "intersect");
gJ("Focus", "trap", "focus");
gJ("Mask", "mask", "mask");
function gJ(J, Y, Z) {
  T(Y, ($) =>
    S(
      `You can't use [x-${Y}] without first installing the "${J}" plugin here: https://alpinejs.dev/plugins/${Z}`,
      $,
    ),
  );
}
HJ.setEvaluator(KZ);
HJ.setReactivityEngine({ reactive: xY, effect: S1, release: x1, raw: R });
var T6 = HJ,
  cJ = T6;
function V6(J) {
  J.directive(
    "intersect",
    J.skipDuringClone(
      (
        Y,
        { value: Z, expression: $, modifiers: W },
        { evaluateLater: X, cleanup: Q },
      ) => {
        let B = X($),
          U = { rootMargin: I6(W), threshold: E6(W) },
          _ = new IntersectionObserver((G) => {
            G.forEach((z) => {
              if (z.isIntersecting === (Z === "leave")) return;
              B(), W.includes("once") && _.disconnect();
            });
          }, U);
        _.observe(Y),
          Q(() => {
            _.disconnect();
          });
      },
    ),
  );
}
function E6(J) {
  if (J.includes("full")) return 0.99;
  if (J.includes("half")) return 0.5;
  if (!J.includes("threshold")) return 0;
  let Y = J[J.indexOf("threshold") + 1];
  if (Y === "100") return 1;
  if (Y === "0") return 0;
  return Number(`.${Y}`);
}
function N6(J) {
  let Y = J.match(/^(-?[0-9]+)(px|%)?$/);
  return Y ? Y[1] + (Y[2] || "px") : void 0;
}
function I6(J) {
  let $ = J.indexOf("margin");
  if ($ === -1) return "0px 0px 0px 0px";
  let W = [];
  for (let X = 1; X < 5; X++) W.push(N6(J[$ + X] || ""));
  return (
    (W = W.filter((X) => X !== void 0)),
    W.length ? W.join(" ").trim() : "0px 0px 0px 0px"
  );
}
var _0 = V6;
window.Alpine = cJ;
cJ.plugin(_0);
cJ.start();
