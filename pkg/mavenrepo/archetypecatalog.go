package mavenrepo

import (
	"fmt"
	"net/url"
	"strings"
)

type ArchetypeCatalog struct {
	Archetypes []Archetype `xml:"archetypes>archetype"`
}

// Archetype is the XML deserialized from the file archetype-catalog.xml which should be at the root of a Maven repo.
type Archetype struct {
	GroupID string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version string `xml:"version"`
	Repository string `xml:"repository"`
	Description string `xml:"description"`
}

// TODO: maven-metadata

// TODO: pom-file

// ToGroupURL generates the group repository URL from the archetype
func (a *Archetype) ToGroupURL(repoURL *url.URL) *url.URL {
	rootPath := strings.Split(a.GroupID, ".")

	newUrl, _ := url.Parse(repoURL.String())
	newUrl.Path = strings.Join(append(strings.Split(newUrl.Path, "/"), rootPath...), "/")
	return newUrl
}

// ToArtifactURL generates the URL to the artifact. It's maven-metadata.xml should be found here.
func (a *Archetype) ToArtifactURL(repoURL *url.URL) *url.URL {
	groupURL := a.ToGroupURL(repoURL)
	groupURL.Path = strings.Join(append(strings.Split(groupURL.Path, "/"), a.ArtifactID), "/")
	return groupURL
}

// ToMavenMetadataURL generates the URL to download the maven metadata of the artifact
func (a *Archetype) ToMavenMetadataURL(repoURL *url.URL) *url.URL {
	groupURL := a.ToGroupURL(repoURL)
	mavenMetadata, _ := url.Parse(groupURL.String())

	mavenMetadata.Path = strings.Join(append(strings.Split(groupURL.Path, "/"), "maven-metadata.xml"), "/")
	return groupURL
}

// ToVersionURL generates the URL to the version of the artifact. Here, a pom file should be found.
func (a *Archetype) ToVersionURL(repoURL *url.URL) *url.URL {
	artifactURL := a.ToArtifactURL(repoURL)
	versionURL, _ := url.Parse(artifactURL.String())

	versionURL.Path = strings.Join(append(strings.Split(artifactURL.Path, "/"), a.Version), "/")
	return versionURL
}

// ToPOMFile generates the URL which will download the POM file.
func (a *Archetype) ToPOMURL(repoURL *url.URL) *url.URL {
	versionURL := a.ToVersionURL(repoURL)
	pomURL, _ := url.Parse(versionURL.String())

	pomFilename := fmt.Sprintf("%v-%v.pom", a.ArtifactID, a.Version)

	pomURL.Path = strings.Join(append(strings.Split(versionURL.Path, "/"), pomFilename), "/")
	return pomURL
}