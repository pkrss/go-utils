package beans

type Page struct {
	Content          interface{}   `json:"content"`          // content
	Last             bool          `json:"last"`             // is last page
	TotalPages       int           `json:"totalPages"`       // total pages
	TotalElements    int           `json:"totalElements"`    // total elements
	Sort             []interface{} `json:"sort"`             // sort params
	First            bool          `json:"first"`            // is first page
	NumberOfElements int           `json:"numberOfElements"` // content list size
	Size             int           `json:"size"`             // query limit size
	Number           int           `json:"number"`           // query offset, from 0
}
