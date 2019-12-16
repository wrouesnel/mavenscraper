package mavenrepo

type ArchetypeCatalogRoot struct {
	ArchetypeCatalog ArchetypeCatalog `xml:"archetype-catalog"`
}

type ArchetypeCatalog struct {
	Archetypes []Archetype `xml:"archetypes"`
}

// Archetype is the XML deserialized from the file archetype-catalog.xml which should be at the root of a Maven repo.
type Archetype struct {
	GroupID string `xml:"groupId"`
	ArtifactID string `xml:"artifactId"`
	Version string `xml:"version"`
	Repository string `xml:"repository"`
	Description string `xml:"description"`
}