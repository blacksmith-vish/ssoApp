package main

import (
	"encoding/xml"
	"fmt"
)

type MarksSet struct {
	XMLName xml.Name `xml:"MarksSet"`
	//Text       string   `xml:",chardata"`
	MarkRecord []struct {
		//Text         string `xml:",chardata"`
		MarkingName  string `xml:"MarkingName"`
		MarkingXmlId string `xml:"MarkingXmlId"`
		Description  string `xml:"Description"`
	} `xml:"MarkRecord"`
}

func main() {
	menu := new(MarksSet)
	err := xml.Unmarshal([]byte(data), menu)
	if err != nil {
		fmt.Printf("error: %v", err)
		return
	}
	fmt.Printf("--- Unmarshal ---\n\n")
	for _, foodNode := range menu.MarkRecord {
		fmt.Printf("MarkingName: %s\n", foodNode.MarkingName)
		fmt.Printf("MarkingXmlId: %s\n", foodNode.MarkingXmlId)
		fmt.Printf("Description: %s\n", foodNode.Description)
		fmt.Printf("---\n")
	}

	xmlText, err := xml.MarshalIndent(menu, " ", " ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
		return
	}
	fmt.Printf("\n--- Marshal ---\n\n")
	fmt.Printf("xml: %s\n", string(xmlText))
}

var data = `
<MarksSet>
	<MarkRecord>
		<MarkingName>Некорректная ссылка</MarkingName>
		<MarkingXmlId>a1op74r8u3fubqijykrrsdbh1e</MarkingXmlId>
		<Description>тут нерабочая ссылка</Description>
	</MarkRecord>
</MarksSet>
`
