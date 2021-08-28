package internal

import (
	"fmt"
	"net/http"
)

func tab() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		q := r.FormValue("q")
		fmt.Fprintln(w, `<html>
<head>
  <script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.js" integrity="sha512-2PRgAav8Os8vLcOAh1gSaDoNLe1fAyq8/G3QSdyjFFD+OqNjLeHE/8q4+S4MEZgPsuo+itHopj+hJvqS8XUQ8A==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/addons/fit/fit.min.js" integrity="sha512-+wh8VA1djpWk3Dj9/IJDu6Ufi4vVQ0zxLv9Vmfo70AbmYFJm0z3NLnV98vdRKBdPDV4Kwpi7EZdr8mDY9L8JIA==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
	<script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/addons/webLinks/webLinks.min.js" integrity="sha512-obszFFlN3K8h7hpqVwXAODf9IOnd1P4PuYRFAwZKTaykxzyMmizo9+eStvrFobjmFs6r6QVsXHMa7ksl34jecg==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
  <link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.css" integrity="sha512-iLYuqv+v/P4u9erpk+KM83Ioe/l7SEmr7wB6g+Kg1qmEit8EShDKnKtLHlv2QXUp7GGJhmqDI+1PhJYLTsfb8w==" crossorigin="anonymous" referrerpolicy="no-referrer" />
	<style>
	  document, body {
      width: 100%;
      height: 100%;
			padding: 0px;
			margin: 0px;
			border: 0px;
      overflow: hidden;
  	}
		#terminal {
	    background: #000000;
		  color: #ffffff;
		  display: inline-block;
      width: 100%;
      height: 100%;
		}
	</style>
</head>
<body>

	<pre id="terminal"></pre>

	<script>
	
		var elem = document.getElementById("terminal");
		elem.tabindex = 0;

		webLinks.apply(Terminal);
		fit.apply(Terminal);

		var terminal = new Terminal();
		terminal.open(elem);

		terminal.webLinksInit();
		var socket = new WebSocket('ws://'+document.location.host+"`+q+`", 'echo');

		socket.addEventListener("open", function () {
			terminal.on('data', function (evt) {
						socket.send(evt);
		 		});
		 });

		 setInterval(function () {
			terminal.fit();
		 }, 1000)

		socket.addEventListener("message", function (evt) {
				terminal.write(event.data);
		});
	</script>
</body>
</html>
`)
	}
}
