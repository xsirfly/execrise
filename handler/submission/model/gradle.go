package model

import (
	"encoding/xml"
	"io/ioutil"
	"os"
)

type TestSuite struct {
	XMLName   xml.Name    `xml:"testsuite"`
	Name      string      `xml:"name,attr"`
	Tests     int64       `xml:"tests,attr"`
	Failures  int64       `xml:"failures,attr"`
	Errors    int64       `xml:"errors,attr"`
	Time      float64     `xml:"time,attr"`
	TestCases []*TestCase `xml:"testcase"`
}

type TestCase struct {
	XMLName   xml.Name         `xml:"testcase"`
	Name      string           `xml:"name,attr"`
	Classname string           `xml:"classname,attr"`
	Time      float64          `xml:"time,attr"`
	Failure   *TestCaseFailure `xml:"failure"`
}

type TestCaseFailure struct {
	XMLName xml.Name `xml:"failure"`
	Message string   `xml:"message,attr"`
	Type    string   `xml:"type,attr"`
	Stack   string   `xml:",innerxml"`
}

func (t *TestSuite) UnmarshalFromFile(filename string) error {
	file, err := os.Open(filename) // For read access.
	if err != nil {
		return err
	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return err
	}

	if err = xml.Unmarshal(data, t); err != nil {
		return err
	}
	return nil
}
