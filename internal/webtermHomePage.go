package internal

import (
	"fmt"
	"net/http"
)

func home(tabs []*WebTermTab) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		wTabsContainer := func() {
			fmt.Fprint(w, `<div class="tabs">
		<div class="card">
			<div class="header">
				<h3>Khayyam <i class="fas fa-angle-down iconM"></i></h3>
			</div>
			<div class="body">
				<ul>`)
			for _, tab := range tabs {
				onClick := `openframe('` + tab.path + `')`
				fmt.Fprintln(w, `<li class="tab" onClick="`+onClick+`">`)
				fmt.Fprintln(w, `<i class="fas fa-home icon"></i>`)
				fmt.Fprintln(w, `<span>`+tab.title+`</span>`)
				fmt.Fprintln(w, `</li>`)
			}
			fmt.Fprintln(w, ` </ul>
			</div>
		</div>
	</div>`)
			fmt.Fprintln(w)
		}

		wHTML := func() {
			fmt.Fprint(w, `<html>
			<head>
			  <link rel="stylesheet" href="https://use.fontawesome.com/releases/v5.8.1/css/all.css" integrity="sha384-50oBUHEmvpQ+1lW4y57PTFmhCaXp0ML5d60M1M7uH2+nqUivzIebhndOJK28anvf" crossorigin="anonymous">
				<link rel="stylesheet" href="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.css" integrity="sha512-iLYuqv+v/P4u9erpk+KM83Ioe/l7SEmr7wB6g+Kg1qmEit8EShDKnKtLHlv2QXUp7GGJhmqDI+1PhJYLTsfb8w==" crossorigin="anonymous" referrerpolicy="no-referrer" />
				<script src="https://cdnjs.cloudflare.com/ajax/libs/xterm/3.14.5/xterm.min.js" integrity="sha512-2PRgAav8Os8vLcOAh1gSaDoNLe1fAyq8/G3QSdyjFFD+OqNjLeHE/8q4+S4MEZgPsuo+itHopj+hJvqS8XUQ8A==" crossorigin="anonymous" referrerpolicy="no-referrer"></script>
				<style>
					@import url('https://fonts.googleapis.com/css?family=Roboto');
					document, body {
						width: 100%;
						height: 100%;
						padding: 0px;
						margin: 0px;
						border: 0px;
					}
					body {
						display: flex;			
						flex-direction: columns;
						overflow: hidden;
					}
					#app {
						display: flex;			
						flex-direction: rows;
						overflow: hidden;
						width: 100%;
						height: 100%;
						background-color: black;
					}
			
					.tabs {
						padding: 50px;
					}
					
					.card {
						width: 300px;
						height: 420px;
						background-color: #1E2B32;
						border-radius: 10px 10px;
					}
					
					.header{
						border-radius: 10px 10px 0px 0px;
						padding: 5px;
						background-color: #2A3942;
					}
					
					h3 {
						color: #FFFFFF;
						font-family: 'Roboto', sans-serif;
						margin-left: 1rem;
					}
					
					.iconM{
						font-size: 18px;
						margin-left: 170px;
						color: #2f89fc;
					}
					
					.icon{
						margin-right: 8px;
					}
					
					.body li{
						transition: 1s all;
						font-family: 'Roboto', sans-serif;
						font-size: 18px;
						padding: 15px;
						margin-left: -40px;
						margin-top: 0px;
						color: #fff;
						list-style: none;
						display: block;
						border-top-right-radius: 10px 10px;
						border-bottom-right-radius: 10px 10px;
					}
					
					li:hover{
						transition: 1s all;
						color: #2f89fc;
						background-color: rgba(42, 56, 65, 0.82);
						border-top-right-radius: 10px 10px;
						border-bottom-right-radius: 10px 10px;
						cursor: pointer;
					}
					
					.body > li {
						float: left;
					}
					
					.body li ul{
						background: #1E2B32;
						margin-left: 280px;
						margin-top: -38px;
						display: none;
						position: absolute;
						border-top-right-radius: 15px 15px;
						border-bottom-right-radius: 15px 15px;
					}
					
					.body li:hover > ul{
						display: block;
						cursor: pointer;
					}
			
					.frame {
						flex-grow: 1;
					}

					.frame iframe {
						position: relative;
						top: 0px;
						left: 0px;
						width: 100%;
						height: 100%;
					}
				</style>
				<script>
					function openframe(id) {
						const elements = document.querySelectorAll('iframe');
						elements.forEach(element => {
							element.setAttribute('z-index', element.id === id ? '1' : '-1');
						})
					}
				</script>
			</head>
			<body>
			<div id="app">`)
			wTabsContainer()
			fmt.Fprint(w, `<div class="frame">`)
			for _, tab := range tabs {
				fmt.Fprint(w, `<iframe id="`+tab.path+`" src="/_tab?q=`+tab.path+`" />`)
			}
			fmt.Fprint(w, `</div>`)
			fmt.Fprint(w, `</div></body></html>`)
		}

		wHTML()
	}
}
