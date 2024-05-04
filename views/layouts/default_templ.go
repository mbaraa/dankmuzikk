// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.663
package layouts

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"dankmuzikk/views/components/header"
	"dankmuzikk/views/components/loading"
	"dankmuzikk/views/components/player"
)

func Default(isMobile bool, children ...templ.Component) templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!doctype html><html lang=\"en\"><head><meta charset=\"UTF-8\"><title>DankMuzikk</title><link rel=\"apple-touch-icon\" sizes=\"180x180\" href=\"/static/apple-touch-icon.png\"><link rel=\"icon\" type=\"image/png\" sizes=\"32x32\" href=\"/static/favicon-32x32.png\"><link rel=\"icon\" type=\"image/png\" sizes=\"16x16\" href=\"/static/favicon-16x16.png\"><link rel=\"manifest\" href=\"/static/site.webmanifest\"><link rel=\"mask-icon\" href=\"/static/safari-pinned-tab.svg\" color=\"#4C8C36\"><meta name=\"msapplication-TileColor\" content=\"#4C8C36\"><meta name=\"theme-color\" content=\"#4C8C36\"><meta name=\"viewport\" content=\"width=device-width, initial-scale=1\"><meta name=\"description\" content=\"Create, Share and Play Music Playlists.\"><meta name=\"keywords\" content=\"dankmuzikk,dank,dank music,music,music playlist,share playlist,group playlist\"><link async defer rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/player-seeker.css\"><link async defer rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/myrad-pro-font.css\"><link async defer rel=\"stylesheet\" type=\"text/css\" href=\"/static/css/audio-nugget-font.css\"><link href=\"/static/css/tailwind.css\" rel=\"stylesheet\"><script src=\"/static/js/htmx.min.js\"></script><script async defer src=\"/static/js/json-enc.js\"></script></head><body style=\"\n                background-image: url(&#34;/static/images/dankground.svg&#34;), linear-gradient(to right, #4C8C36 , #9EE07E);\n                background-repeat: repeat;\n                background-size: cover;\n                padding: 0px;\n                margin: 0px;\n                background-position: center;\n                background-attachment: fixed;\n                min-height: 100dvh;\n            \">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = header.Header(isMobile).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		for _, child := range children {
			templ_7745c5c3_Err = child.Render(ctx, templ_7745c5c3_Buffer)
			if templ_7745c5c3_Err != nil {
				return templ_7745c5c3_Err
			}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"loading\" class=\"hidden\">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = loading.Loading().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = player.PlayerSticky(isMobile).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script src=\"/static/js/utils.js\"></script><script type=\"module\">\n\t\t\t    function registerServiceWorkers() {\n\t\t\t    \tif (!(\"serviceWorker\" in navigator)) {\n\t\t\t    \t\tconsole.error(\"Browser doesn't support service workers\");\n\t\t\t    \t\treturn;\n\t\t\t    \t}\n\t\t\t    \twindow.addEventListener(\"load\", () => {\n\t\t\t    \t\tnavigator.serviceWorker\n\t\t\t    \t\t\t.register(\"/static/js/service-worker.js\")\n\t\t\t    \t\t\t.then((reg) => {\n\t\t\t    \t\t\t\tconsole.log(\"Service Worker Registered\", reg);\n\t\t\t    \t\t\t})\n\t\t\t    \t\t\t.catch((err) => {\n\t\t\t    \t\t\t\tconsole.log(\"Service Worker Registration failed:\", err);\n\t\t\t    \t\t\t});\n\t\t                });\n\t\t\t    }\n\t\t\t    registerServiceWorkers();\n\t\t    </script></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
