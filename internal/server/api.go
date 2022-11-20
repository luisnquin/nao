package server

/*
import (
	"net/http"
	"os"

	"github.com/andreaskoch/go-fswatch"
	"github.com/fatih/color"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luisnquin/nao/v2/internal/config"
	"github.com/luisnquin/nao/v2/internal/store"
)

func New(port string, quiet, verbose bool) *Server {
	s := &Server{
		itWasMyFault: make(chan bool),
		router:       echo.New(),
		box:          store.New(),
		quiet:        quiet,
		port:         port,
	}

	if verbose {
		s.router.Use(middleware.Logger())
	}

	s.router.Use(middleware.CORS())

	return s
}

func (a *Server) Start() error {
	go a.watchAndRefreshData()
	a.mountHandlers()

	return http.ListenAndServe(a.port, a.router)
}

func (a *Server) watchAndRefreshData() {
	w := fswatch.NewFileWatcher(config.App.Paths.DataFile, 1)
	w.Start()

	if !a.quiet {
		color.New(color.FgHiCyan).Fprintln(os.Stdout, "👀  Watching "+config.App.Paths.DataFile+"\n")
	}

	var timesMod int

	for w.IsRunning() {
		select {
		case <-w.Modified():
			select {
			case <-a.itWasMyFault:
				continue

			default:
				timesMod++
				a.box.ReplaceBox(store.JustLoadBox())

				if !a.quiet {
					color.New(color.FgHiBlue).Fprintf(os.Stdout, "\rData refreshed(x%d)", timesMod)
				}
			}

		case <-w.Moved():
			color.New(color.FgHiRed).Fprintln(os.Stderr, "Error: Unable to find data file, apparently moved")
			os.Exit(1)
		}
	}
}

func (a *Server) mountHandlers() {
	notes := a.router.Group("/notes")
	notes.GET("", a.GetNotesHandler())
	notes.POST("", a.NewNoteHandler())
	notes.GET("/:id", a.GetNoteHandler())
	notes.PUT("/:id", a.ModifyNoteHandler())
	notes.DELETE("/:id", a.DeleteNoteHandler())
	notes.PATCH("/:id", a.ModifyNoteContentHandler())
}
*/
