package cmd

import (
	"github.com/majd/ipatool/v2/pkg/appstore"
)

type SearchResult struct {
	Success bool           `json:"success"`
	Count   int            `json:"count"`
	Results []appstore.App `json:"results"`
}

func Search(term string, limit int64) (*SearchResult, error) {
	infoResult, err := dependencies.AppStore.AccountInfo()
	if err != nil {
		return nil, err
	}

	output, err := dependencies.AppStore.Search(appstore.SearchInput{
		Account: infoResult.Account,
		Term:    term,
		Limit:   limit,
	})
	if err != nil {
		return nil, err
	}

	return &SearchResult{
		Success: true,
		Count:   output.Count,
		Results: output.Results,
	}, nil
}
