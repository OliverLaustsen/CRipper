package model

type VisitList []Visit

type Visit struct {
	Chapter int64
	Comic   string
	Url     string
}
