// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package views

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func MenuBar(city string) templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div id=\"menu\"><div id=\"menu-city\"><span>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(city)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `server/views/menu_bar.templ`, Line: 5, Col: 34}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("</span></div><div id=\"menu-date\"></div><div id=\"menu-time\"></div><div id=\"menu-settings\">Settings</div><script>\n            function updateDateTime() {\n                const now = new Date();\n                const dateOptions = { weekday: 'long', day: 'numeric', month: 'long' };\n                const frenchDate = now.toLocaleDateString('fr-FR', dateOptions);\n\n                const timeOptions = { hour: 'numeric', minute: 'numeric', second: 'numeric' };\n                const frenchTime = now.toLocaleTimeString('fr-FR', timeOptions);\n\n                document.getElementById(\"menu-date\").innerText = frenchDate;\n                document.getElementById(\"menu-time\").innerText = frenchTime;\n            }\n\n            function toggleQRCode() {\n                if (document.getElementById(\"qrcode\").innerHTML === \"\") {\n                    const url = new URL(window.location.href);\n                    const baseUrl = `${url.protocol}//${url.hostname}${url.port ? ':' + url.port : ''}`;\n                    const adminUrl = `${baseUrl}/admin`;\n                    console.log(adminUrl);\n                    new QRCode(document.getElementById(\"qrcode\"), {\n                        text: adminUrl,\n                        width: 256,\n                        height: 256\n                    });\n                } else {\n                    document.getElementById(\"qrcode\").innerHTML = \"\"; \n                }\n            }\n\n            // date update\n            updateDateTime();\n            setInterval(updateDateTime, 1000);\n\n            // QR code\n            document.getElementById(\"menu-settings\").addEventListener(\"click\", toggleQRCode);\n        </script></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
