package front

import (
	"fmt"

	"github.com/maxence-charriere/go-app/v9/pkg/app"
)

type Home struct {
	app.Compo
}

func NewHome() *Home {
	return &Home{}
}

func GetAPIRoute(ctx app.Context) string {
	return fmt.Sprintf("%s://%s", ctx.Page().URL().Scheme, ctx.Page().URL().Host)
}

func (h *Home) OnMount(ctx app.Context) {
	fmt.Println("component mounted version")
}

func (h *Home) Render() app.UI {
	return app.Main().Body(

		app.Div().Body(
			app.H1().Text("fldarkstat"),
		).Class("header"),

		app.Div().Body().Class("footer"),
	)
}
