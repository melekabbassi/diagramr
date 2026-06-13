package studio

import (
	"context"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	"github.com/melekabbassi/diagramr/internal/renderer"
	"github.com/melekabbassi/diagramr/internal/watcher"
)

//go:embed web
var webFS embed.FS

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Config holds everything the studio server needs to start.
type Config struct {
	InitialMermaid string
	OutputPath     string // empty → <RenderConfig.RootPath>/diagramr.mmd
	RenderConfig   RenderConfig
}

// Serve starts the studio web server, opens the browser, and blocks until
// the process receives an interrupt signal.
func Serve(cfg Config) error {
	if cfg.OutputPath == "" {
		cfg.OutputPath = filepath.Join(cfg.RenderConfig.RootPath, "diagramr.mmd")
	}

	hub := newHub()
	go hub.run()

	var mu sync.RWMutex
	current := cfg.InitialMermaid

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	go func() {
		_ = watcher.Watch(ctx, cfg.RenderConfig.RootPath, []string{".go"}, func() {
			out, err := renderMermaid(cfg.RenderConfig)
			if err != nil {
				return
			}
			mu.Lock()
			current = out
			mu.Unlock()
			hub.broadcast([]byte(out))
		})
	}()

	ln, err := net.Listen("tcp", ":0")
	if err != nil {
		return fmt.Errorf("studio: listen: %w", err)
	}
	port := ln.Addr().(*net.TCPAddr).Port

	staticFS, err := fs.Sub(webFS, "web")
	if err != nil {
		return err
	}

	r := mux.NewRouter()

	r.PathPrefix("/static/").Handler(
		http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))),
	)

	r.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		f, ferr := webFS.Open("web/index.html")
		if ferr != nil {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		defer f.Close()
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = io.Copy(w, f)
	})

	r.HandleFunc("/ws", func(w http.ResponseWriter, req *http.Request) {
		conn, uerr := upgrader.Upgrade(w, req, nil)
		if uerr != nil {
			return
		}
		c := &client{hub: hub, conn: conn, send: make(chan []byte, 8)}
		hub.registerC <- c

		mu.RLock()
		snap := current
		mu.RUnlock()
		c.send <- []byte(snap)

		go c.writePump()

		for {
			if _, _, rerr := conn.ReadMessage(); rerr != nil {
				hub.unregisterC <- c
				return
			}
		}
	})

	r.HandleFunc("/render", func(w http.ResponseWriter, req *http.Request) {
		var opts renderer.Options
		if derr := json.NewDecoder(req.Body).Decode(&opts); derr != nil {
			http.Error(w, derr.Error(), http.StatusBadRequest)
			return
		}
		rc := cfg.RenderConfig
		rc.RenderOpts = opts
		out, rerr := renderMermaid(rc)
		if rerr != nil {
			http.Error(w, rerr.Error(), http.StatusInternalServerError)
			return
		}
		mu.Lock()
		current = out
		mu.Unlock()
		hub.broadcast([]byte(out))
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodPost)

	r.HandleFunc("/save", func(w http.ResponseWriter, req *http.Request) {
		body, rerr := io.ReadAll(req.Body)
		if rerr != nil {
			http.Error(w, rerr.Error(), http.StatusBadRequest)
			return
		}
		if werr := os.WriteFile(cfg.OutputPath, body, 0o644); werr != nil {
			http.Error(w, werr.Error(), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	}).Methods(http.MethodPost)

	srv := &http.Server{Handler: r}

	fmt.Fprintf(os.Stderr, "\ndiagramr studio\n  http://localhost:%d\n\n  Ctrl+C to stop\n\n", port)
	openBrowser(fmt.Sprintf("http://localhost:%d", port))

	go func() {
		<-ctx.Done()
		shutCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutCtx)
	}()

	if serr := srv.Serve(ln); serr != nil && serr != http.ErrServerClosed {
		return serr
	}
	return nil
}

func openBrowser(url string) {
	var cmd string
	var args []string
	if isWSL() {
		cmd, args = "powershell.exe", []string{"-Command", "Start-Process", "'" + url + "'"}
	} else {
		switch runtime.GOOS {
		case "darwin":
			cmd, args = "open", []string{url}
		case "windows":
			cmd, args = "cmd", []string{"/c", "start", url}
		default:
			cmd, args = "xdg-open", []string{url}
		}
	}
	_ = exec.Command(cmd, args...).Start()
}

func isWSL() bool {
	_, ok := os.LookupEnv("WSL_DISTRO_NAME")
	return ok
}
