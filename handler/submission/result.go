package submission

import "exercise/handler/submission/model"

type Result struct {
	Success   bool             `json:"success"`
	TestSuite *model.TestSuite `json:"test_suite"`
}
