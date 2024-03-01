package model

type Product struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Size     string `json:"size"`
	Code     string `json:"code"`
	Quantity int    `json:"quantity"`
}
