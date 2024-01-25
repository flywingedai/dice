package results

type Result struct {
	Chances  *Chances  `json:"chances"`
	Result   int64     `json:"result"`
	Children []*Result `json:"children"`
}
