package commands

import (
	"fmt"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/section/sectionctl/api"
)

// DomainsCmd manages domains on Section
type DomainsCmd struct {
	List DomainsListCmd `cmd help:"List domains on Section." default:"1"`
}

// DomainsListCmd handles listing domains on Section
type DomainsListCmd struct {
	AccountID int `short:"a" help:"ID of account to list domains under"`
}

// Run executes the command
func (c *DomainsListCmd) Run(cli *CLI, logWriters *LogWriters) (err error) {
	var aids []int
	if c.AccountID == 0 {
		s := NewSpinner("Looking up accounts",logWriters)
		s.Start()

		as, err := api.Accounts()
		if err != nil {
			return fmt.Errorf("unable to look up accounts: %w", err)
		}
		for _, a := range as {
			aids = append(aids, a.ID)
		}

		s.Stop()
	} else {
		aids = append(aids, c.AccountID)
	}

	s := NewSpinner("Looking up domains",logWriters)
	s.Start()
	domains := make(map[int][]api.DomainsResponse)
	for _, id := range aids {
		ds, err := api.Domains(id)
		if err != nil {
			return fmt.Errorf("unable to look up domains: %w", err)
		}
		domains[id] = ds
	}
	s.Stop()

	table := NewTable(cli, os.Stdout)
	table.SetHeader([]string{"Account ID", "Domain", "Engaged"})
	table.SetHeaderColor(tablewriter.Colors{tablewriter.Normal,tablewriter.FgWhiteColor},tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor},tablewriter.Colors{tablewriter.Normal, tablewriter.FgWhiteColor})
	table.SetColumnColor(tablewriter.Colors{tablewriter.Normal,tablewriter.FgWhiteColor},tablewriter.Colors{tablewriter.Bold, tablewriter.FgHiCyanColor},tablewriter.Colors{tablewriter.Normal,tablewriter.FgWhiteColor})
	for id, ds := range domains {
		for _, d := range ds {
			r := []string{strconv.Itoa(id), d.DomainName, fmt.Sprintf("%t", d.Engaged)}
			table.Append(r)
		}
	}

	table.Render()
	return err
}
