package api

import (
	"bytes"
	"encoding/csv"
	"github.com/sirupsen/logrus"
	"strings"
)

func (f *FaqResponse) GetCSVBytes() []byte {
	return GetCSVBytesForFaqs(f.Faqs)
}

func GetCSVBytesForFaqs(faqs *[]Faqs) []byte {
	b := &bytes.Buffer{}
	wr := csv.NewWriter(b)
	wr.Comma = ';'
	if err := wr.Write([]string{"Question", "Answer"}); err != nil {
		logrus.Warningf("Unable to write CSV Data")
	}
	for _, faq := range *faqs {
		if err := wr.Write([]string{strings.ReplaceAll(faq.Title, ";", ","), strings.ReplaceAll(faq.Answer, ";", ",")}); err != nil {
			logrus.Warningf("Unable to write CSV Data")
		}

	}
	wr.Flush()
	return b.Bytes()
}
