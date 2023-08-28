package main

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"
)

type BitwardenFolder struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

type BitwardenLoginURI struct {
	Match string `json:"match"`
	Uri   string `json:"uri"`
}

type BitwardenItemsLogin struct {
	Uris     []BitwardenLoginURI `json:"uris"`
	Username string              `json:"username"`
	Password string              `json:"password"`
	Totp     []string            `json:"totp"`
}

type BitwardenItems struct {
	Id             string              `json:"id"`
	OrganizationId string              `json:"organizationId,omitempty"`
	FolderId       string              `json:"folderId,omitempty"`
	Type           int                 `json:"type"`
	Reprompt       int                 `json:"reprompt"`
	Name           string              `json:"name"`
	Notes          []string            `json:"notes"`
	Favorite       bool                `json:"favorite"`
	Login          BitwardenItemsLogin `json:"login"`
	CollectionIds  []string            `json:"collectionIds"`
}

type BitwardenExport struct {
	Encrypted bool              `json:"encrypted"`
	Folders   []BitwardenFolder `json:"folders"`
	Items     []BitwardenItems  `json:"items"`
}

func bitwardenExport(base *[]Base) {
	knuvi_folder := BitwardenFolder{
		Name: "knuvi",
		Id:   uuid.New().String(),
	}
	record := BitwardenExport{
		Folders: []BitwardenFolder{knuvi_folder},
	}
	for _, item := range *base {
		bitem := BitwardenItems{
			Id:       uuid.New().String(),
			FolderId: knuvi_folder.Id,
			Type:     1,
			Reprompt: 0,
			Name:     item.Host,
			Favorite: false,
			Login: BitwardenItemsLogin{
				Username: item.Login,
				Password: item.Password,
				Uris: []BitwardenLoginURI{
					BitwardenLoginURI{
						Uri: item.Host,
					},
				},
			},
		}
		record.Items = append(record.Items, bitem)
	}
	f, _ := os.Create("output.json")
	defer f.Close()
	as_json, _ := json.MarshalIndent(record, "", "\t")
	f.Write(as_json)
}
