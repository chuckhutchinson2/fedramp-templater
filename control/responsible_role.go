package control

import (
	"github.com/jbowtie/gokogiri/xml"
	"errors"
	"regexp"
	"strings"
	"fmt"
)

// findResponsibleRole looks for the Responsible Role cell in the control table.
func findResponsibleRole(ct *Table) (*responsibleRole, error) {
	nodes, err := ct.searchSubtree(".//w:tc//w:t[contains(., 'Responsible Role')]")
	if err != nil {
		return nil, err
	}
	if len(nodes) != 1 {
		return nil, errors.New("could not find Responsible Role cell")
	}
	// Not sure why we have to get the parent's parent, but we need to.
	// If we only go up once, it won't find the other text nodes.
	parentNode := nodes[0].Parent().Parent()
	childNodes, err := parentNode.Search(".//w:t")
	if err != nil || len(childNodes) < 1 {
		return nil, errors.New("Should not happen, cannot find text nodes.")
	}
	return &responsibleRole{parentNode: parentNode, textNodes: &childNodes}, nil
}

// responsibleRole is the container for the responsible role cell.
type responsibleRole struct {
	parentNode xml.Node
	textNodes *[]xml.Node
}

// getContent returns the full string representation of the content of the cell itself.
func (r *responsibleRole) getContent() string {
	return r.parentNode.Content()
}

// setValue will set the value of the responsible role cell and do any needed formatting.
// In this case, it will just place the text after "Responsible Role: "
// If there are other nodes, we don't care about them, zero the content out.
func (r *responsibleRole) setValue(value string) {
	for idx, node := range *(r.textNodes) {
		if idx == 0 {
			node.SetContent(fmt.Sprintf("Responsible Role: %s", value))
		} else {
			node.SetContent("")
		}
	}
}

// isDefaultValue contains the logic to detect if the input is a default value. This is looking at the extracted
// value of responsible role and not the full string representation.
func (r *responsibleRole) isDefaultValue(value string) bool {
	return value == ""
}

// getValue extracts the unique value from the full string representation.
// It looks at all the text after "Responsible Role:".
func (r *responsibleRole) getValue() string {
	re := regexp.MustCompile("Responsible Role:?")
	result := ""
	// Get all the substrings
	subStrings := re.Split(r.parentNode.Content(), -1)
	// Go through the substrings and find the first one that is not empty.
	// (So far has always been the string at index 1)
	for _, subString := range subStrings {
		// If we find an non-empty string, return it.
		trimmedString := strings.TrimSpace(subString)
		if len(trimmedString) > 0 {
			result = trimmedString
			break
		}
	}
	return result
}