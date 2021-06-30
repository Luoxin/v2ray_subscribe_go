package main

type CheckResult struct {
	NodeName string
	Speed    float64
	Delay    float64
}

//go:generate pie CheckResultList.*
type CheckResultList []*CheckResult
