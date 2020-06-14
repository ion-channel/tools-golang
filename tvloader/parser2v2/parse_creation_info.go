// SPDX-License-Identifier: Apache-2.0 OR GPL-2.0-or-later

package parser2v2

import (
	"fmt"

	"github.com/spdx/tools-golang/spdx"
)

func (parser *tvParser2_2) parsePairFromCreationInfo2_2(tag string, value string) error {
	// fail if not in Creation Info parser state
	if parser.st != psCreationInfo2_2 {
		return fmt.Errorf("Got invalid state %v in parsePairFromCreationInfo2_2", parser.st)
	}

	// create an SPDX Creation Info data struct if we don't have one already
	if parser.doc.CreationInfo == nil {
		parser.doc.CreationInfo = &spdx.CreationInfo2_2{}
	}

	ci := parser.doc.CreationInfo
	switch tag {
	case "SPDXVersion":
		ci.SPDXVersion = value
	case "DataLicense":
		ci.DataLicense = value
	case "SPDXID":
		ci.SPDXIdentifier = spdx.ElementID(value)
	case "DocumentName":
		ci.DocumentName = value
	case "DocumentNamespace":
		ci.DocumentNamespace = value
	case "ExternalDocumentRef":
		ci.ExternalDocumentReferences = append(ci.ExternalDocumentReferences, value)
	case "LicenseListVersion":
		ci.LicenseListVersion = value
	case "Creator":
		subkey, subvalue, err := extractSubs(value)
		if err != nil {
			return err
		}
		switch subkey {
		case "Person":
			ci.CreatorPersons = append(ci.CreatorPersons, subvalue)
		case "Organization":
			ci.CreatorOrganizations = append(ci.CreatorOrganizations, subvalue)
		case "Tool":
			ci.CreatorTools = append(ci.CreatorTools, subvalue)
		default:
			return fmt.Errorf("unrecognized Creator type %v", subkey)
		}
	case "Created":
		ci.Created = value
	case "CreatorComment":
		ci.CreatorComment = value
	case "DocumentComment":
		ci.DocumentComment = value

	// tag for going on to package section
	case "PackageName":
		parser.st = psPackage2_2
		parser.pkg = &spdx.Package2_2{
			FilesAnalyzed:             true,
			IsFilesAnalyzedTagPresent: false,
		}
		return parser.parsePairFromPackage2_2(tag, value)
	// tag for going on to _unpackaged_ file section
	case "FileName":
		// leave pkg as nil, so that packages will be placed in UnpackagedFiles
		parser.st = psFile2_2
		parser.pkg = nil
		return parser.parsePairFromFile2_2(tag, value)
	// tag for going on to other license section
	case "LicenseID":
		parser.st = psOtherLicense2_2
		return parser.parsePairFromOtherLicense2_2(tag, value)
	// tag for going on to review section (DEPRECATED)
	case "Reviewer":
		parser.st = psReview2_2
		return parser.parsePairFromReview2_2(tag, value)
	// for relationship tags, pass along but don't change state
	case "Relationship":
		parser.rln = &spdx.Relationship2_2{}
		parser.doc.Relationships = append(parser.doc.Relationships, parser.rln)
		return parser.parsePairForRelationship2_2(tag, value)
	case "RelationshipComment":
		return parser.parsePairForRelationship2_2(tag, value)
	// for annotation tags, pass along but don't change state
	case "Annotator":
		parser.ann = &spdx.Annotation2_2{}
		parser.doc.Annotations = append(parser.doc.Annotations, parser.ann)
		return parser.parsePairForAnnotation2_2(tag, value)
	case "AnnotationDate":
		return parser.parsePairForAnnotation2_2(tag, value)
	case "AnnotationType":
		return parser.parsePairForAnnotation2_2(tag, value)
	case "SPDXREF":
		return parser.parsePairForAnnotation2_2(tag, value)
	case "AnnotationComment":
		return parser.parsePairForAnnotation2_2(tag, value)
	default:
		return fmt.Errorf("received unknown tag %v in CreationInfo section", tag)
	}

	return nil
}