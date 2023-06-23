// Copyright 2023 The Wa Authors. All rights reserved.

function HTTPTransport() {
  'use strict';

  function playback(output, events) {
    var timeout;
    output({Kind: 'start'});
    function next() {
      if (!events || events.length === 0) {
        output({Kind: 'end'});
        return;
      }
      var e = events.shift();
      if (e.Delay === 0) {
        output({Kind: 'stdout', Body: e.Message});
        next();
        return;
      }
      timeout = setTimeout(function() {
        output({Kind: 'stdout', Body: e.Message});
        next();
      }, e.Delay / 1000000);
    }
    next();
    return {
      Stop: function() {
        clearTimeout(timeout);
      }
    }
  }

  function error(output, msg) {
    output({Kind: 'start'});
    output({Kind: 'stderr', Body: msg});
    output({Kind: 'end'});
  }

  var seq = 0;
  return {
    Run: function(body, output, options) {
      seq++;
      var cur = seq;
      var playing;
      $.ajax(playgroundOptions.compileURL, {
        type: 'POST',
        data: {'version': 2, 'body': body},
        dataType: 'json',
        success: function(data) {
          if (seq != cur) return;
          if (!data) return;
          if (playing != null) playing.Stop();
          if (data.Errors) {
            error(output, data.Errors);
            return;
          }
          playing = playback(output, data.Events);
        },
        error: function() {
          error(output, 'Error communicating with remote server.');
        }
      });
      return {
        Kill: function() {
          if (playing != null) playing.Stop();
          output({Kind: 'end', Body: 'killed'});
        }
      };
    }
  };
}

function SocketTransport() {
  'use strict';

  var id = 0;
  var outputs = {};
  var started = {};
  var websocket = new WebSocket('ws://' + window.location.host + '/socket');

  websocket.onclose = function() {
    console.log('websocket connection closed');
  }

  websocket.onmessage = function(e) {
    var m = JSON.parse(e.data);
    var output = outputs[m.Id];
    if (output === null)
      return;
    if (!started[m.Id]) {
      output({Kind: 'start'});
      started[m.Id] = true;
    }
    output({Kind: m.Kind, Body: m.Body});
  }

  function send(m) {
    websocket.send(JSON.stringify(m));
  }

  return {
    Run: function(body, output, options) {
      var thisID = id+'';
      id++;
      outputs[thisID] = output;
      send({Id: thisID, Kind: 'run', Body: body, Options: options});
      return {
        Kill: function() {
          send({Id: thisID, Kind: 'kill'});
        }
      };
    }
  };
}

function PlaygroundOutput(el) {
  'use strict';

  return function(write) {
    if (write.Kind == 'start') {
      el.innerHTML = '';
      return;
    }

    var cl = 'system';
    if (write.Kind == 'stdout' || write.Kind == 'stderr')
      cl = write.Kind;

    var m = write.Body;
    if (write.Kind == 'end')
      m = '\nProgram exited' + (m?(': '+m):'.');

    if (m.indexOf('IMAGE:') === 0) {
      // TODO(adg): buffer all writes before creating image
      var url = 'data:image/png;base64,' + m.substr(6);
      var img = document.createElement('img');
      img.src = url;
      el.appendChild(img);
      return;
    }

    // ^L clears the screen.
    var s = m.split('\x0c');
    if (s.length > 1) {
      el.innerHTML = '';
      m = s.pop();
    }

    m = m.replace(/&/g, '&amp;');
    m = m.replace(/</g, '&lt;');
    m = m.replace(/>/g, '&gt;');

    var needScroll = (el.scrollTop + el.offsetHeight) == el.scrollHeight;

    var span = document.createElement('span');
    span.className = cl;
    span.innerHTML = m;
    el.appendChild(span);
    $(el).fadeIn()

    if (needScroll)
      el.scrollTop = el.scrollHeight - el.offsetHeight;
  }
}

var playgroundOptions = {}

var defaultOptions = {
  'compileURL': '/-/play/compile',
  'fmtURL': '/-/play/fmt',
};

function kclPlaygroundOptions(opts) {
  playgroundOptions = $.extend(defaultOptions, playgroundOptions, opts);
}

kclPlaygroundOptions({});

(function() {
  function lineHighlight(error) {
    var regex = /prog.wa:([0-9]+)/g;
    var r = regex.exec(error);
    while (r) {
      $(".lines div").eq(r[1]-1).addClass("lineerror");
      r = regex.exec(error);
    }
  }
  function highlightOutput(wrappedOutput) {
    return function(write) {
      if (write.Body) lineHighlight(write.Body);
      wrappedOutput(write);
    }
  }
  function lineClear() {
    $(".lineerror").removeClass("lineerror");
  }

  // opts is an object with these keys
  //  codeEl - code editor element
  //  outputEl - program output element
  //  runEl - run button element
  //  fmtEl - fmt button element (optional)
  //  enableHistory - enable using HTML5 history API (optional)
  //  transport - playground transport to use (default is HTTPTransport)
  function playground(opts) {
    var opts = $.extend(opts, playgroundOptions);
    var code = $(opts.codeEl);
    var transport = opts['transport'] || new HTTPTransport();
    var running;

    console.log(code);
    var editorProps = {
      lineNumbers: true,
      indentWithTabs: false,
      mode: 'wa',
      smartIndent: true,
      tabSize: 4,
      indentUnit: 4,
    };

    if (typeof opts['theme'] !== 'undefined') {
      editorProps.theme = opts['theme'];
    }

    var editor = CodeMirror.fromTextArea(code[0], editorProps);

    var outdiv = $(opts.outputEl).empty().hide();
    var output = $('<pre/>').appendTo(outdiv);

    function body() {
      return editor.getValue();
    }
    function setBody(text) {
      editor.setValue(text);
    }
    function origin(href) {
      return (""+href).split("/").slice(0, 3).join("/");
    }

    var pushedEmpty = (window.location.pathname == "/");
    function inputChanged() {
      if (pushedEmpty) {
        return;
      }
      pushedEmpty = true;
      $(opts.shareURLEl).hide();
      window.history.pushState(null, "", "/");
    }
    function popState(e) {
      if (e === null) {
        return;
      }
      if (e && e.state && e.state.code) {
        setBody(e.state.code);
      }
    }
    var rewriteHistory = false;
    if (window.history && window.history.pushState && window.addEventListener && opts.enableHistory) {
      rewriteHistory = true;
      code[0].addEventListener('input', inputChanged);
      window.addEventListener('popstate', popState);
    }

    function setError(error) {
      if (running) running.Kill();
      lineClear();
      lineHighlight(error);
      output.empty().addClass("error").text(error);
      if (error === "") {
        outdiv.hide();
      } else {
        outdiv.fadeIn();
      }
    }
    function loading() {
      lineClear();
      if (running) running.Kill();
      output.removeClass("error").fadeIn().text('Waiting for remote server...');
    }
    function run() {
      $(opts.outputEl).fadeIn();
      loading();
      running = transport.Run(body(), highlightOutput(PlaygroundOutput(output[0])));
    }

    function fmt() {
      loading();
      var data = {"body": body()};
      data["imports"] = "true";
      $.ajax(playgroundOptions.fmtURL, {
        data: data,
        type: "POST",
        dataType: "json",
        success: function(data) {
          if (data.Error) {
            setError(data.Error);
          } else {
            setBody(data.Body);
            setError("");
          }
        }
      });
    }

    $(opts.runEl).click(run);
    $(opts.fmtEl).click(fmt);
  }

  window.playground = playground;
})();
