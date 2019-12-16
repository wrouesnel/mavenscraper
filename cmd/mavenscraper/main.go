// mavenscraper is a prototype for selectively mirroring a maven repository

package main

import (
	"encoding/xml"
	"fmt"
	"github.com/dustin/go-humanize"
	"github.com/wrouesnel/mavenscraper/pkg/mavenrepo"
	"go.uber.org/zap"
	"gopkg.in/alecthomas/kingpin.v2"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func main() {
	appConfig := struct{
		mavenURL *url.URL

		logLevel string
		logTarget string
		logFormat string
	}{}

	app := kingpin.New("mavenscraper", "Prototype tool for selective mirroring of Maven repos in Golang")

	app.Flag("log.level", "Log Level").Default("info").StringVar(&appConfig.logLevel)
	app.Flag("log.target", "Log Target").Default("stderr").EnumVar(&appConfig.logTarget, "stderr", "stdout" )
	app.Flag("log.format", "Log level").Default("console").EnumVar(&appConfig.logFormat, "console", "json")

	app.Flag("maven.repo", "Maven repository to analyze").Default("https://repo.maven.apache.org/maven2").URLVar(&appConfig.mavenURL)

	_ = kingpin.MustParse(app.Parse(os.Args[1:]))

	config := zap.NewProductionConfig()
	logLevel := zap.AtomicLevel{}
	if err := logLevel.UnmarshalText([]byte(appConfig.logLevel)); err != nil {
		app.Fatalf("Log level is not valid: %v", err.Error())
	}
	config.Level = logLevel
	config.Encoding = appConfig.logFormat
	config.OutputPaths = []string{appConfig.logTarget}
	config.ErrorOutputPaths = []string{appConfig.logTarget}

	log, err := config.Build()
	if err != nil {
		app.Fatalf("Could not configure logging: %v", err.Error())
	}

	log.Debug("Application logging ready")

	archetypeURL, _ := url.Parse(appConfig.mavenURL.String())
	archetypeURL.Path = strings.Join(append(strings.Split(archetypeURL.Path, "/"),"archetype-catalog.xml"),"/")

	log.Info("Downloading archetype-catalog.xml", zap.String("url", archetypeURL.String()))

	resp, err := http.Head(archetypeURL.String())
	if err != nil {
		log.Fatal("HTTP request failed", zap.Error(err))
	}

	log.Info(fmt.Sprintf("archetype-catalog.xml is %v", humanize.Bytes(uint64(resp.ContentLength))))
	resp.Body.Close()

	resp, err = http.Get(archetypeURL.String())
	if err != nil {
		log.Fatal("HTTP request failed", zap.Error(err))
	}
	defer resp.Body.Close()
	// Decode the body as an archetype
	decoder := xml.NewDecoder(resp.Body)

	catalog := &mavenrepo.ArchetypeCatalogRoot{}

	if err := decoder.Decode(catalog); err != nil {
		log.Fatal("Failed to decode archetype catalog", zap.Error(err))
	}

	log.Info(fmt.Sprintf("Got %v archetypes", len(catalog.ArchetypeCatalog.Archetypes)))

	for

	os.Exit(0)
}