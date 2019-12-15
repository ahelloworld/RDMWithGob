package main

type TestObj struct {
	T int    `json:"T"`
	M string `json:"M"`
}

func init() {
	gobRegister("TestObj", &TestObj{})
}
