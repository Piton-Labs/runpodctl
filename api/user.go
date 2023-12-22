package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
)

type UserOut struct {
	Myself *MySelf         `json:"data"`
	Errors []*GraphQLError `json:"errors"`
}

type MySelf struct {
	Id                     string `json:"id"`
	ContainerRegistryCreds *ContainerRegistryCreds
}

type ContainerRegistryCreds struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func GetCreds() (creds *ContainerRegistryCreds, err error) {
	input := Input{
		Query: `
			query getMyself {
				myself {
				id
				email
				containerRegistryCreds {
					id
					name
				}
			}
		`,
	}
	res, err := Query(input)
	if err != nil {
		return
	}
	if res.StatusCode != 200 {
		err = fmt.Errorf("statuscode %d", res.StatusCode)
	}
	defer res.Body.Close()

	rawData, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	data := &UserOut{}
	if err = json.Unmarshal(rawData, data); err != nil {
		return
	}
	if len(data.Errors) > 0 {
		err = errors.New(data.Errors[0].Message)
		return
	}
	if data == nil || data.Myself == nil || data.Myself.ContainerRegistryCreds == nil {
		err = fmt.Errorf("data is nil: %s", string(rawData))
		return
	}

	creds = data.Myself.ContainerRegistryCreds
	return
}
