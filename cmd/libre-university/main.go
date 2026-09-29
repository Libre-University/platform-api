// libre-university, LibreUniversity'nin varsayılan dağıtımıdır: çekirdek
// platform ile MVP modüllerini (obs, lms) tek binary olarak derler
// (ADR-0004, ADR-0011, ADR-0012). Farklı modül seti isteyen kurumlar bu
// dizini kopyalayıp içe aktarılan modülleri değiştirir.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Libre-University/platform-api/core"
	"github.com/Libre-University/platform-api/platform"

	// Modüller init() içinde kendilerini kaydeder.
	_ "github.com/Libre-University/platform-api/modules/lms"
	_ "github.com/Libre-University/platform-api/modules/obs"
)

// version derleme sırasında -ldflags "-X main.version=..." ile verilir.
var version = "dev"

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "hata:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	cmd := "serve"
	if len(args) > 0 {
		cmd = args[0]
		args = args[1:]
	}
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	cfg := platform.ConfigFromEnv()
	cfg.Version = version

	switch cmd {
	case "serve":
		return serve(cfg, log, args)
	case "migrate":
		return migrate(cfg, log)
	case "modules":
		for _, m := range core.Registered() {
			fmt.Printf("%s\t%s\t%d izin\n", m.Name(), m.Version(), len(m.Permissions()))
		}
		return nil
	case "version":
		fmt.Println(version)
		return nil
	case "help", "-h", "--help":
		usage()
		return nil
	default:
		usage()
		return fmt.Errorf("bilinmeyen komut %q", cmd)
	}
}

func usage() {
	fmt.Fprint(os.Stderr, `Kullanım: libre-university <komut>

Komutlar:
  serve     HTTP sunucusunu başlat (varsayılan)
  migrate   Çekirdek ve modül migration'larını uygula (LU_DATABASE_URL gerekli)
  modules   Derlenmiş modülleri listele
  version   Sürümü yazdır

Ortam değişkenleri:
  LU_HTTP_ADDR      Dinleme adresi (varsayılan :8080)
  LU_DATABASE_URL   PostgreSQL bağlantı adresi
`)
}

func serve(cfg platform.Config, log *slog.Logger, args []string) error {
	fs := flag.NewFlagSet("serve", flag.ContinueOnError)
	fs.StringVar(&cfg.HTTPAddr, "addr", cfg.HTTPAddr, "HTTP dinleme adresi")
	if err := fs.Parse(args); err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := platform.New(ctx, cfg, log, core.Registered())
	if err != nil {
		return err
	}
	defer app.Close()

	srv := &http.Server{
		Addr:              cfg.HTTPAddr,
		Handler:           app.Handler(),
		ReadHeaderTimeout: 10 * time.Second,
	}
	errc := make(chan error, 1)
	go func() {
		log.Info("sunucu başlıyor", "addr", cfg.HTTPAddr, "version", cfg.Version, "modules", len(app.Modules()))
		errc <- srv.ListenAndServe()
	}()

	select {
	case err := <-errc:
		if errors.Is(err, http.ErrServerClosed) {
			return nil
		}
		return err
	case <-ctx.Done():
		log.Info("kapatılıyor")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		return srv.Shutdown(shutdownCtx)
	}
}

func migrate(cfg platform.Config, log *slog.Logger) error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	app, err := platform.New(ctx, cfg, log, core.Registered())
	if err != nil {
		return err
	}
	defer app.Close()
	return app.Migrate(ctx)
}
