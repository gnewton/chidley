package main

import (
	"encoding/xml"
)

type PubmedArticleSet_Type struct {
	PubmedArticle     []PubmedArticle_Type     `xml:"PubmedArticle,omitempty" json:",omitempty"`
	PubmedBookArticle []PubmedBookArticle_Type `xml:"PubmedBookArticle,omitempty" json:",omitempty"`
	XMLName           xml.Name                 `xml:"PubmedArticleSet,omitempty" json:",omitempty"`
}

type PubmedArticle_Type struct {
	MedlineCitation MedlineCitation_Type `xml:"MedlineCitation,omitempty" json:",omitempty"`
	PubmedData      PubmedData_Type      `xml:"PubmedData,omitempty" json:",omitempty"`
	XMLName         xml.Name             `xml:"PubmedArticle,omitempty" json:",omitempty"`
}

type MedlineCitation_Type struct {
	Article                 Article_Type                 `xml:"Article,omitempty" json:",omitempty"`
	ChemicalList            ChemicalList_Type            `xml:"ChemicalList,omitempty" json:",omitempty"`
	CitationSubset          []string                     `xml:"CitationSubset,omitempty" json:",omitempty"`
	CommentsCorrectionsList CommentsCorrectionsList_Type `xml:"CommentsCorrectionsList,omitempty" json:",omitempty"`
	DateCompleted           DateCompleted_Type           `xml:"DateCompleted,omitempty" json:",omitempty"`
	DateCreated             DateCreated_Type             `xml:"DateCreated,omitempty" json:",omitempty"`
	DateRevised             DateRevised_Type             `xml:"DateRevised,omitempty" json:",omitempty"`
	GeneSymbolList          GeneSymbolList_Type          `xml:"GeneSymbolList,omitempty" json:",omitempty"`
	GeneralNote             []GeneralNote_Type           `xml:"GeneralNote,omitempty" json:",omitempty"`
	InvestigatorList        InvestigatorList_Type        `xml:"InvestigatorList,omitempty" json:",omitempty"`
	KeywordList             KeywordList_Type             `xml:"KeywordList,omitempty" json:",omitempty"`
	MedlineJournalInfo      MedlineJournalInfo_Type      `xml:"MedlineJournalInfo,omitempty" json:",omitempty"`
	MeshHeadingList         MeshHeadingList_Type         `xml:"MeshHeadingList,omitempty" json:",omitempty"`
	NumberOfReferences      string                       `xml:"NumberOfReferences,omitempty" json:",omitempty"`
	OtherAbstract           OtherAbstract_Type           `xml:"OtherAbstract,omitempty" json:",omitempty"`
	OtherID                 []OtherID_Type               `xml:"OtherID,omitempty" json:",omitempty"`
	Owner                   string                       `xml:"Owner,attr"`
	PMID                    PMID_Type                    `xml:"PMID,omitempty" json:",omitempty"`
	PersonalNameSubjectList PersonalNameSubjectList_Type `xml:"PersonalNameSubjectList,omitempty" json:",omitempty"`
	SpaceFlightMission      []string                     `xml:"SpaceFlightMission,omitempty" json:",omitempty"`
	Status                  string                       `xml:"Status,attr"`
	SupplMeshList           SupplMeshList_Type           `xml:"SupplMeshList,omitempty" json:",omitempty"`
	VersionDate             string                       `xml:"VersionDate,attr"`
	VersionID               string                       `xml:"VersionID,attr"`
	XMLName                 xml.Name                     `xml:"MedlineCitation,omitempty" json:",omitempty"`
}

type KeywordList_Type struct {
	Keyword []Keyword_Type `xml:"Keyword,omitempty" json:",omitempty"`
	Owner   string         `xml:"Owner,attr"`
	XMLName xml.Name       `xml:"KeywordList,omitempty" json:",omitempty"`
}

type Keyword_Type struct {
	MajorTopicYN string   `xml:"MajorTopicYN,attr"`
	Text         string   `xml:",chardata" json:",omitempty"`
	XMLName      xml.Name `xml:"Keyword,omitempty" json:",omitempty"`
}

type OtherAbstract_Type struct {
	AbstractText         AbstractText_Type `xml:"AbstractText,omitempty" json:",omitempty"`
	CopyrightInformation string            `xml:"CopyrightInformation,omitempty" json:",omitempty"`
	Language             string            `xml:"Language,attr"`
	Type                 string            `xml:"Type,attr"`
	XMLName              xml.Name          `xml:"OtherAbstract,omitempty" json:",omitempty"`
}

type AbstractText_Type struct {
	Label       string   `xml:"Label,attr"`
	NlmCategory string   `xml:"NlmCategory,attr"`
	Text        string   `xml:",chardata" json:",omitempty"`
	XMLName     xml.Name `xml:"AbstractText,omitempty" json:",omitempty"`
}

type PersonalNameSubjectList_Type struct {
	PersonalNameSubject []PersonalNameSubject_Type `xml:"PersonalNameSubject,omitempty" json:",omitempty"`
	XMLName             xml.Name                   `xml:"PersonalNameSubjectList,omitempty" json:",omitempty"`
}

type PersonalNameSubject_Type struct {
	ForeName string   `xml:"ForeName,omitempty" json:",omitempty"`
	Initials string   `xml:"Initials,omitempty" json:",omitempty"`
	LastName string   `xml:"LastName,omitempty" json:",omitempty"`
	Suffix   string   `xml:"Suffix,omitempty" json:",omitempty"`
	XMLName  xml.Name `xml:"PersonalNameSubject,omitempty" json:",omitempty"`
}

type PMID_Type struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	Version string   `xml:"Version,attr"`
	XMLName xml.Name `xml:"PMID,omitempty" json:",omitempty"`
}

type DateCompleted_Type struct {
	Day     string   `xml:"Day,omitempty" json:",omitempty"`
	Month   string   `xml:"Month,omitempty" json:",omitempty"`
	XMLName xml.Name `xml:"DateCompleted,omitempty" json:",omitempty"`
	Year    string   `xml:"Year,omitempty" json:",omitempty"`
}

type DateRevised_Type struct {
	Day     string   `xml:"Day,omitempty" json:",omitempty"`
	Month   string   `xml:"Month,omitempty" json:",omitempty"`
	XMLName xml.Name `xml:"DateRevised,omitempty" json:",omitempty"`
	Year    string   `xml:"Year,omitempty" json:",omitempty"`
}

type Article_Type struct {
	Abstract            Abstract_Type            `xml:"Abstract,omitempty" json:",omitempty"`
	ArticleDate         ArticleDate_Type         `xml:"ArticleDate,omitempty" json:",omitempty"`
	ArticleTitle        ArticleTitle_Type        `xml:"ArticleTitle,omitempty" json:",omitempty"`
	AuthorList          AuthorList_Type          `xml:"AuthorList,omitempty" json:",omitempty"`
	DataBankList        DataBankList_Type        `xml:"DataBankList,omitempty" json:",omitempty"`
	ELocationID         []ELocationID_Type       `xml:"ELocationID,omitempty" json:",omitempty"`
	GrantList           GrantList_Type           `xml:"GrantList,omitempty" json:",omitempty"`
	Journal             Journal_Type             `xml:"Journal,omitempty" json:",omitempty"`
	Language            []string                 `xml:"Language,omitempty" json:",omitempty"`
	Pagination          Pagination_Type          `xml:"Pagination,omitempty" json:",omitempty"`
	PubModel            string                   `xml:"PubModel,attr"`
	PublicationTypeList PublicationTypeList_Type `xml:"PublicationTypeList,omitempty" json:",omitempty"`
	VernacularTitle     string                   `xml:"VernacularTitle,omitempty" json:",omitempty"`
	XMLName             xml.Name                 `xml:"Article,omitempty" json:",omitempty"`
}

type ELocationID_Type struct {
	EIdType string   `xml:"EIdType,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	ValidYN string   `xml:"ValidYN,attr"`
	XMLName xml.Name `xml:"ELocationID,omitempty" json:",omitempty"`
}

type ArticleDate_Type struct {
	DateType string   `xml:"DateType,attr"`
	Day      string   `xml:"Day,omitempty" json:",omitempty"`
	Month    string   `xml:"Month,omitempty" json:",omitempty"`
	XMLName  xml.Name `xml:"ArticleDate,omitempty" json:",omitempty"`
	Year     string   `xml:"Year,omitempty" json:",omitempty"`
}

type GrantList_Type struct {
	CompleteYN string       `xml:"CompleteYN,attr"`
	Grant      []Grant_Type `xml:"Grant,omitempty" json:",omitempty"`
	XMLName    xml.Name     `xml:"GrantList,omitempty" json:",omitempty"`
}

type Grant_Type struct {
	Acronym string   `xml:"Acronym,omitempty" json:",omitempty"`
	Agency  string   `xml:"Agency,omitempty" json:",omitempty"`
	Country string   `xml:"Country,omitempty" json:",omitempty"`
	GrantID string   `xml:"GrantID,omitempty" json:",omitempty"`
	XMLName xml.Name `xml:"Grant,omitempty" json:",omitempty"`
}

type Journal_Type struct {
	ISOAbbreviation string            `xml:"ISOAbbreviation,omitempty" json:",omitempty"`
	ISSN            ISSN_Type         `xml:"ISSN,omitempty" json:",omitempty"`
	JournalIssue    JournalIssue_Type `xml:"JournalIssue,omitempty" json:",omitempty"`
	Title           string            `xml:"Title,omitempty" json:",omitempty"`
	XMLName         xml.Name          `xml:"Journal,omitempty" json:",omitempty"`
}

type JournalIssue_Type struct {
	CitedMedium string       `xml:"CitedMedium,attr"`
	Issue       string       `xml:"Issue,omitempty" json:",omitempty"`
	PubDate     PubDate_Type `xml:"PubDate,omitempty" json:",omitempty"`
	Volume      string       `xml:"Volume,omitempty" json:",omitempty"`
	XMLName     xml.Name     `xml:"JournalIssue,omitempty" json:",omitempty"`
}

type PubDate_Type struct {
	Day         string   `xml:"Day,omitempty" json:",omitempty"`
	MedlineDate string   `xml:"MedlineDate,omitempty" json:",omitempty"`
	Month       string   `xml:"Month,omitempty" json:",omitempty"`
	Season      string   `xml:"Season,omitempty" json:",omitempty"`
	XMLName     xml.Name `xml:"PubDate,omitempty" json:",omitempty"`
	Year        string   `xml:"Year,omitempty" json:",omitempty"`
}

type ISSN_Type struct {
	IssnType string   `xml:"IssnType,attr"`
	Text     string   `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"ISSN,omitempty" json:",omitempty"`
}

type Pagination_Type struct {
	MedlinePgn string   `xml:"MedlinePgn,omitempty" json:",omitempty"`
	XMLName    xml.Name `xml:"Pagination,omitempty" json:",omitempty"`
}

type Abstract_Type struct {
	AbstractText         []AbstractText_Type `xml:"AbstractText,omitempty" json:",omitempty"`
	CopyrightInformation string              `xml:"CopyrightInformation,omitempty" json:",omitempty"`
	XMLName              xml.Name            `xml:"Abstract,omitempty" json:",omitempty"`
}

type DataBankList_Type struct {
	CompleteYN string          `xml:"CompleteYN,attr"`
	DataBank   []DataBank_Type `xml:"DataBank,omitempty" json:",omitempty"`
	XMLName    xml.Name        `xml:"DataBankList,omitempty" json:",omitempty"`
}

type DataBank_Type struct {
	AccessionNumberList AccessionNumberList_Type `xml:"AccessionNumberList,omitempty" json:",omitempty"`
	DataBankName        string                   `xml:"DataBankName,omitempty" json:",omitempty"`
	XMLName             xml.Name                 `xml:"DataBank,omitempty" json:",omitempty"`
}

type AccessionNumberList_Type struct {
	AccessionNumber []string `xml:"AccessionNumber,omitempty" json:",omitempty"`
	XMLName         xml.Name `xml:"AccessionNumberList,omitempty" json:",omitempty"`
}

type ArticleTitle_Type struct {
	Book    string   `xml:"book,attr"`
	Part    string   `xml:"part,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"ArticleTitle,omitempty" json:",omitempty"`
}

type AuthorList_Type struct {
	Author     []Author_Type `xml:"Author,omitempty" json:",omitempty"`
	CompleteYN string        `xml:"CompleteYN,attr"`
	Type       string        `xml:"Type,attr"`
	XMLName    xml.Name      `xml:"AuthorList,omitempty" json:",omitempty"`
}

type Author_Type struct {
	Affiliation    string          `xml:"Affiliation,omitempty" json:",omitempty"`
	CollectiveName string          `xml:"CollectiveName,omitempty" json:",omitempty"`
	ForeName       string          `xml:"ForeName,omitempty" json:",omitempty"`
	Identifier     Identifier_Type `xml:"Identifier,omitempty" json:",omitempty"`
	Initials       string          `xml:"Initials,omitempty" json:",omitempty"`
	LastName       string          `xml:"LastName,omitempty" json:",omitempty"`
	Suffix         string          `xml:"Suffix,omitempty" json:",omitempty"`
	ValidYN        string          `xml:"ValidYN,attr"`
	XMLName        xml.Name        `xml:"Author,omitempty" json:",omitempty"`
}

type Identifier_Type struct {
	Source  string   `xml:"Source,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"Identifier,omitempty" json:",omitempty"`
}

type PublicationTypeList_Type struct {
	PublicationType []string `xml:"PublicationType,omitempty" json:",omitempty"`
	XMLName         xml.Name `xml:"PublicationTypeList,omitempty" json:",omitempty"`
}

type ChemicalList_Type struct {
	Chemical []Chemical_Type `xml:"Chemical,omitempty" json:",omitempty"`
	XMLName  xml.Name        `xml:"ChemicalList,omitempty" json:",omitempty"`
}

type Chemical_Type struct {
	NameOfSubstance string   `xml:"NameOfSubstance,omitempty" json:",omitempty"`
	RegistryNumber  string   `xml:"RegistryNumber,omitempty" json:",omitempty"`
	XMLName         xml.Name `xml:"Chemical,omitempty" json:",omitempty"`
}

type CommentsCorrectionsList_Type struct {
	CommentsCorrections []CommentsCorrections_Type `xml:"CommentsCorrections,omitempty" json:",omitempty"`
	XMLName             xml.Name                   `xml:"CommentsCorrectionsList,omitempty" json:",omitempty"`
}

type CommentsCorrections_Type struct {
	Note      string    `xml:"Note,omitempty" json:",omitempty"`
	PMID      PMID_Type `xml:"PMID,omitempty" json:",omitempty"`
	RefSource string    `xml:"RefSource,omitempty" json:",omitempty"`
	RefType   string    `xml:"RefType,attr"`
	XMLName   xml.Name  `xml:"CommentsCorrections,omitempty" json:",omitempty"`
}

type OtherID_Type struct {
	Source  string   `xml:"Source,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"OtherID,omitempty" json:",omitempty"`
}

type GeneSymbolList_Type struct {
	GeneSymbol []string `xml:"GeneSymbol,omitempty" json:",omitempty"`
	XMLName    xml.Name `xml:"GeneSymbolList,omitempty" json:",omitempty"`
}

type GeneralNote_Type struct {
	Owner   string   `xml:"Owner,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"GeneralNote,omitempty" json:",omitempty"`
}

type InvestigatorList_Type struct {
	Investigator []Investigator_Type `xml:"Investigator,omitempty" json:",omitempty"`
	XMLName      xml.Name            `xml:"InvestigatorList,omitempty" json:",omitempty"`
}

type Investigator_Type struct {
	Affiliation string   `xml:"Affiliation,omitempty" json:",omitempty"`
	ForeName    string   `xml:"ForeName,omitempty" json:",omitempty"`
	Initials    string   `xml:"Initials,omitempty" json:",omitempty"`
	LastName    string   `xml:"LastName,omitempty" json:",omitempty"`
	Suffix      string   `xml:"Suffix,omitempty" json:",omitempty"`
	ValidYN     string   `xml:"ValidYN,attr"`
	XMLName     xml.Name `xml:"Investigator,omitempty" json:",omitempty"`
}

type SupplMeshList_Type struct {
	SupplMeshName []SupplMeshName_Type `xml:"SupplMeshName,omitempty" json:",omitempty"`
	XMLName       xml.Name             `xml:"SupplMeshList,omitempty" json:",omitempty"`
}

type SupplMeshName_Type struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	Type    string   `xml:"Type,attr"`
	XMLName xml.Name `xml:"SupplMeshName,omitempty" json:",omitempty"`
}

type DateCreated_Type struct {
	Day     string   `xml:"Day,omitempty" json:",omitempty"`
	Month   string   `xml:"Month,omitempty" json:",omitempty"`
	XMLName xml.Name `xml:"DateCreated,omitempty" json:",omitempty"`
	Year    string   `xml:"Year,omitempty" json:",omitempty"`
}

type MedlineJournalInfo_Type struct {
	Country     string   `xml:"Country,omitempty" json:",omitempty"`
	ISSNLinking string   `xml:"ISSNLinking,omitempty" json:",omitempty"`
	MedlineTA   string   `xml:"MedlineTA,omitempty" json:",omitempty"`
	NlmUniqueID string   `xml:"NlmUniqueID,omitempty" json:",omitempty"`
	XMLName     xml.Name `xml:"MedlineJournalInfo,omitempty" json:",omitempty"`
}

type MeshHeadingList_Type struct {
	MeshHeading []MeshHeading_Type `xml:"MeshHeading,omitempty" json:",omitempty"`
	XMLName     xml.Name           `xml:"MeshHeadingList,omitempty" json:",omitempty"`
}

type MeshHeading_Type struct {
	DescriptorName DescriptorName_Type  `xml:"DescriptorName,omitempty" json:",omitempty"`
	QualifierName  []QualifierName_Type `xml:"QualifierName,omitempty" json:",omitempty"`
	XMLName        xml.Name             `xml:"MeshHeading,omitempty" json:",omitempty"`
}

type DescriptorName_Type struct {
	MajorTopicYN string   `xml:"MajorTopicYN,attr"`
	Text         string   `xml:",chardata" json:",omitempty"`
	Type         string   `xml:"Type,attr"`
	XMLName      xml.Name `xml:"DescriptorName,omitempty" json:",omitempty"`
}

type QualifierName_Type struct {
	MajorTopicYN string   `xml:"MajorTopicYN,attr"`
	Text         string   `xml:",chardata" json:",omitempty"`
	XMLName      xml.Name `xml:"QualifierName,omitempty" json:",omitempty"`
}

type PubmedData_Type struct {
	ArticleIdList     ArticleIdList_Type `xml:"ArticleIdList,omitempty" json:",omitempty"`
	History           History_Type       `xml:"History,omitempty" json:",omitempty"`
	PublicationStatus string             `xml:"PublicationStatus,omitempty" json:",omitempty"`
	XMLName           xml.Name           `xml:"PubmedData,omitempty" json:",omitempty"`
}

type History_Type struct {
	PubMedPubDate []PubMedPubDate_Type `xml:"PubMedPubDate,omitempty" json:",omitempty"`
	XMLName       xml.Name             `xml:"History,omitempty" json:",omitempty"`
}

type PubMedPubDate_Type struct {
	Day       string   `xml:"Day,omitempty" json:",omitempty"`
	Hour      string   `xml:"Hour,omitempty" json:",omitempty"`
	Minute    string   `xml:"Minute,omitempty" json:",omitempty"`
	Month     string   `xml:"Month,omitempty" json:",omitempty"`
	PubStatus string   `xml:"PubStatus,attr"`
	XMLName   xml.Name `xml:"PubMedPubDate,omitempty" json:",omitempty"`
	Year      string   `xml:"Year,omitempty" json:",omitempty"`
}

type ArticleIdList_Type struct {
	ArticleId []ArticleId_Type `xml:"ArticleId,omitempty" json:",omitempty"`
	XMLName   xml.Name         `xml:"ArticleIdList,omitempty" json:",omitempty"`
}

type ArticleId_Type struct {
	IdType  string   `xml:"IdType,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"ArticleId,omitempty" json:",omitempty"`
}

type PubmedBookArticle_Type struct {
	BookDocument   BookDocument_Type   `xml:"BookDocument,omitempty" json:",omitempty"`
	PubmedBookData PubmedBookData_Type `xml:"PubmedBookData,omitempty" json:",omitempty"`
	XMLName        xml.Name            `xml:"PubmedBookArticle,omitempty" json:",omitempty"`
}

type BookDocument_Type struct {
	Abstract         Abstract_Type         `xml:"Abstract,omitempty" json:",omitempty"`
	ArticleIdList    ArticleIdList_Type    `xml:"ArticleIdList,omitempty" json:",omitempty"`
	ArticleTitle     ArticleTitle_Type     `xml:"ArticleTitle,omitempty" json:",omitempty"`
	AuthorList       AuthorList_Type       `xml:"AuthorList,omitempty" json:",omitempty"`
	Book             Book_Type             `xml:"Book,omitempty" json:",omitempty"`
	ContributionDate ContributionDate_Type `xml:"ContributionDate,omitempty" json:",omitempty"`
	DateRevised      DateRevised_Type      `xml:"DateRevised,omitempty" json:",omitempty"`
	ItemList         []ItemList_Type       `xml:"ItemList,omitempty" json:",omitempty"`
	KeywordList      KeywordList_Type      `xml:"KeywordList,omitempty" json:",omitempty"`
	Language         string                `xml:"Language,omitempty" json:",omitempty"`
	LocationLabel    LocationLabel_Type    `xml:"LocationLabel,omitempty" json:",omitempty"`
	PMID             PMID_Type             `xml:"PMID,omitempty" json:",omitempty"`
	Sections         Sections_Type         `xml:"Sections,omitempty" json:",omitempty"`
	XMLName          xml.Name              `xml:"BookDocument,omitempty" json:",omitempty"`
}

type Book_Type struct {
	AuthorList      []AuthorList_Type    `xml:"AuthorList,omitempty" json:",omitempty"`
	BeginningDate   BeginningDate_Type   `xml:"BeginningDate,omitempty" json:",omitempty"`
	BookTitle       BookTitle_Type       `xml:"BookTitle,omitempty" json:",omitempty"`
	CollectionTitle CollectionTitle_Type `xml:"CollectionTitle,omitempty" json:",omitempty"`
	Edition         string               `xml:"Edition,omitempty" json:",omitempty"`
	EndingDate      EndingDate_Type      `xml:"EndingDate,omitempty" json:",omitempty"`
	Isbn            []string             `xml:"Isbn,omitempty" json:",omitempty"`
	Medium          string               `xml:"Medium,omitempty" json:",omitempty"`
	PubDate         PubDate_Type         `xml:"PubDate,omitempty" json:",omitempty"`
	Publisher       Publisher_Type       `xml:"Publisher,omitempty" json:",omitempty"`
	XMLName         xml.Name             `xml:"Book,omitempty" json:",omitempty"`
}

type BeginningDate_Type struct {
	XMLName xml.Name `xml:"BeginningDate,omitempty" json:",omitempty"`
	Year    string   `xml:"Year,omitempty" json:",omitempty"`
}

type EndingDate_Type struct {
	XMLName xml.Name `xml:"EndingDate,omitempty" json:",omitempty"`
	Year    string   `xml:"Year,omitempty" json:",omitempty"`
}

type Publisher_Type struct {
	PublisherLocation string   `xml:"PublisherLocation,omitempty" json:",omitempty"`
	PublisherName     string   `xml:"PublisherName,omitempty" json:",omitempty"`
	XMLName           xml.Name `xml:"Publisher,omitempty" json:",omitempty"`
}

type BookTitle_Type struct {
	Book    string   `xml:"book,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"BookTitle,omitempty" json:",omitempty"`
}

type CollectionTitle_Type struct {
	Book    string   `xml:"book,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"CollectionTitle,omitempty" json:",omitempty"`
}

type ItemList_Type struct {
	Item     string   `xml:"Item,omitempty" json:",omitempty"`
	ListType string   `xml:"ListType,attr"`
	XMLName  xml.Name `xml:"ItemList,omitempty" json:",omitempty"`
}

type Sections_Type struct {
	Section []Section_Type `xml:"Section,omitempty" json:",omitempty"`
	XMLName xml.Name       `xml:"Sections,omitempty" json:",omitempty"`
}

type Section_Type struct {
	LocationLabel LocationLabel_Type `xml:"LocationLabel,omitempty" json:",omitempty"`
	Section       []Section_Type     `xml:"Section,omitempty" json:",omitempty"`
	SectionTitle  SectionTitle_Type  `xml:"SectionTitle,omitempty" json:",omitempty"`
	XMLName       xml.Name           `xml:"Section,omitempty" json:",omitempty"`
}

type SectionTitle_Type struct {
	Book    string   `xml:"book,attr"`
	Part    string   `xml:"part,attr"`
	Sec     string   `xml:"sec,attr"`
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"SectionTitle,omitempty" json:",omitempty"`
}

type LocationLabel_Type struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	Type    string   `xml:"Type,attr"`
	XMLName xml.Name `xml:"LocationLabel,omitempty" json:",omitempty"`
}

type ContributionDate_Type struct {
	Day     string   `xml:"Day,omitempty" json:",omitempty"`
	Month   string   `xml:"Month,omitempty" json:",omitempty"`
	XMLName xml.Name `xml:"ContributionDate,omitempty" json:",omitempty"`
	Year    string   `xml:"Year,omitempty" json:",omitempty"`
}

type PubmedBookData_Type struct {
	ArticleIdList     ArticleIdList_Type `xml:"ArticleIdList,omitempty" json:",omitempty"`
	History           History_Type       `xml:"History,omitempty" json:",omitempty"`
	PublicationStatus string             `xml:"PublicationStatus,omitempty" json:",omitempty"`
	XMLName           xml.Name           `xml:"PubmedBookData,omitempty" json:",omitempty"`
}
