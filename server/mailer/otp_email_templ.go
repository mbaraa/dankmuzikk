// Code generated by templ - DO NOT EDIT.

// templ: version: v0.3.833
package mailer

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import templruntime "github.com/a-h/templ/runtime"

func otpEmail(name, code string) templ.Component {
	return templruntime.GeneratedTemplate(func(templ_7745c5c3_Input templruntime.GeneratedComponentInput) (templ_7745c5c3_Err error) {
		templ_7745c5c3_W, ctx := templ_7745c5c3_Input.Writer, templ_7745c5c3_Input.Context
		if templ_7745c5c3_CtxErr := ctx.Err(); templ_7745c5c3_CtxErr != nil {
			return templ_7745c5c3_CtxErr
		}
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templruntime.GetBuffer(templ_7745c5c3_W)
		if !templ_7745c5c3_IsBuffer {
			defer func() {
				templ_7745c5c3_BufErr := templruntime.ReleaseBuffer(templ_7745c5c3_Buffer)
				if templ_7745c5c3_Err == nil {
					templ_7745c5c3_Err = templ_7745c5c3_BufErr
				}
			}()
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 1, "<html lang=\"en\"><body><div style=\"\n                    font-family: Helvetica, Arial, sans-serif;\n                    min-width: 1000px;\n                    overflow: auto;\n                    line-height: 2;\n                    \"><div style=\"margin: 50px auto; width: 70%; padding: 20px 0\"><div style=\"border-bottom: 1px solid #eee\"><a href=\"https://dankmuzikk.com\" style=\"\n                               font-size: 1.4em;\n                               color: #00466a;\n                               text-decoration: none;\n                               font-weight: 600;\n                             \">DankMuzikk</a></div><p style=\"font-size: 1.1em\">Hi ")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var2 string
		templ_7745c5c3_Var2, templ_7745c5c3_Err = templ.JoinStringErrs(name)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `otp_email.templ`, Line: 26, Col: 42}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var2))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 2, ",</p><p>Thank you for using DankMuzikk, below lies your one-time-password, which will be valid for the next <b>30 minutes</b>, don't share this code with anyone in order to keep your account safe 😁</p><h2 style=\"\n                           background: #00466a;\n                           margin: 0 auto;\n                           width: max-content;\n                           padding: 0 10px;\n                           color: #fff;\n                           border-radius: 4px;\n                         \">")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		var templ_7745c5c3_Var3 string
		templ_7745c5c3_Var3, templ_7745c5c3_Err = templ.JoinStringErrs(code)
		if templ_7745c5c3_Err != nil {
			return templ.Error{Err: templ_7745c5c3_Err, FileName: `otp_email.templ`, Line: 42, Col: 12}
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString(templ.EscapeString(templ_7745c5c3_Var3))
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templruntime.WriteString(templ_7745c5c3_Buffer, 3, "</h2><p style=\"font-size: 0.9em\">Regards,<br>DankMuzikk Admin</p><hr style=\"border: none; border-top: 1px solid #eee\"><div style=\"\n                           float: right;\n                           padding: 8px 0;\n                           color: #aaa;\n                           font-size: 0.8em;\n                           line-height: 1;\n                           font-weight: 300;\n                         \"><p>DankMuzikk</p><p><a href=\"mailto:pub@mbaraa.com\">pub@mbaraa.com</a></p></div></div></div></body></html>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		return nil
	})
}

var _ = templruntime.GeneratedTemplate
