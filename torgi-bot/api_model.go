package main

type ApiResponse struct {
	Content          []ApiLot        `json:"content"`
	Pageable         Pageable        `json:"pageable"`
	CategoryFacet    []CategoryFacet `json:"categoryFacet"`
	TotalElements    int             `json:"totalElements"`
	Last             bool            `json:"last"`
	TotalPages       int             `json:"totalPages"`
	First            bool            `json:"first"`
	NumberOfElements int             `json:"numberOfElements"`
	Size             int             `json:"size"`
	Number           int             `json:"number"`
	Sort             Sort            `json:"sort"`
	Empty            bool            `json:"empty"`
}

type ApiLot struct {
	ID                                string           `json:"id"`
	NoticeNumber                      string           `json:"noticeNumber"`
	LotNumber                         int              `json:"lotNumber"`
	LotStatus                         string           `json:"lotStatus"`
	BiddType                          NamedCode        `json:"biddType"`
	BiddForm                          NamedCode        `json:"biddForm"`
	LotName                           string           `json:"lotName"`
	LotDescription                    string           `json:"lotDescription"`
	BiddEndTime                       string           `json:"biddEndTime"`
	LotImages                         []string         `json:"lotImages"`
	LotVideos                         []string         `json:"lotVideos"`
	Characteristics                   []Characteristic `json:"characteristics"`
	CurrencyCode                      string           `json:"currencyCode"`
	SubjectRFCode                     string           `json:"subjectRFCode"`
	Category                          NamedCode        `json:"category"`
	CreateDate                        string           `json:"createDate"`
	TimeZoneName                      string           `json:"timeZoneName"`
	TimezoneOffset                    string           `json:"timezoneOffset"`
	HasAppeals                        bool             `json:"hasAppeals"`
	IsStopped                         bool             `json:"isStopped"`
	Attributes                        []Attribute      `json:"attributes"`
	IsAnnulled                        bool             `json:"isAnnulled"`
	NoticeFirstVersionPublicationDate string           `json:"noticeFirstVersionPublicationDate"`
	NpaHintCode                       string           `json:"npaHintCode"`
	TypeTransaction                   string           `json:"typeTransaction"`
}

type NamedCode struct {
	Code string `json:"code"`
	Name string `json:"name"`
}

type Characteristic struct {
	CharacteristicValue interface{} `json:"characteristicValue"` // может быть строка, массив или объект
	Name                string      `json:"name"`
	Code                string      `json:"code"`
	Type                string      `json:"type"`
	Unit                *Unit       `json:"unit,omitempty"`
}

type Unit struct {
	Code   string `json:"code"`
	Name   string `json:"name"`
	Symbol string `json:"symbol"`
}

type Attribute struct {
	Code          string         `json:"code"`
	FullName      string         `json:"fullName"`
	Value         interface{}    `json:"value,omitempty"` // строка, bool или NamedCode
	AttributeType string         `json:"attributeType"`
	Group         AttributeGroup `json:"group"`
	SortOrder     int            `json:"sortOrder"`
}

type AttributeGroup struct {
	Code             string `json:"code"`
	Name             string `json:"name"`
	DisplayGroupType string `json:"displayGroupType"`
}

type CategoryFacet struct {
	ID    string `json:"_id"`
	Count int    `json:"count"`
}

type Pageable struct {
	PageNumber int  `json:"pageNumber"`
	PageSize   int  `json:"pageSize"`
	Sort       Sort `json:"sort"`
	Offset     int  `json:"offset"`
	Unpaged    bool `json:"unpaged"`
	Paged      bool `json:"paged"`
}

type Sort struct {
	Unsorted bool `json:"unsorted"`
	Sorted   bool `json:"sorted"`
	Empty    bool `json:"empty"`
}
