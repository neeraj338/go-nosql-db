package command

import (
	"errors"
	"go-nosql-db/service"
	"strings"

	"github.com/alexflint/go-arg"
)

type Command interface {
	Run() interface{}
}

//--table students -k="name" -v="neeraj" -select "id,name,age --data="{}" -delete"
type Args struct {
	Table       string `arg:"required" help:"--table db table name"` // `arg:"required" help:"--table db table name"`
	FilterKey   string `arg:"-k,--filter-key" help:"optional: filter condition"`
	FilterValue string `arg:"-v,--filter-value" help:"optional: filter condition"`
	Select      string `arg:"-s,--select" help:"optional: json nodes to select"`
	Data        string `arg:"-d,--data" help:"optional json body"`
	Delete      bool   `arg:"--delete" help:"delete record for filter"`
}

func (args *Args) buildProjectionArr() []string {
	var projections []string = make([]string, 0)
	for _, value := range strings.Split(args.Select, ",") {
		if trm := strings.TrimSpace(value); len(trm) > 0 {
			projections = append(projections, trm)
		}
	}
	return projections
}

func BuildCommand() Command {
	var args Args
	arg.MustParse(&args)
	argValidator := CommandLineArgValidator{CommandLineArgs: args}

	err := argValidator.Validate()
	if err != nil {
		panic(err)
	}

	filter := service.Filter{FilterKey: args.FilterKey, FilterValue: args.FilterValue, Table: args.Table}

	if d := strings.TrimSpace(args.Data); len(d) > 0 {
		return &service.Persist{Table: filter.Table, Data: d}
	} else if args.Delete {
		return &service.Delete{Filter: filter}
	}

	return &service.Select{Filter: filter, Projection: args.buildProjectionArr()}
}

type CommandLineArgValidator struct {
	CommandLineArgs Args
}

func (oneCommand *CommandLineArgValidator) Validate() error {
	c := map[bool]int{true: 1, false: 0}[len(strings.TrimSpace(oneCommand.CommandLineArgs.Data)) > 0]
	d := map[bool]int{true: 1, false: 0}[oneCommand.CommandLineArgs.Delete]
	s := map[bool]int{true: 1, false: 0}[len(strings.TrimSpace(oneCommand.CommandLineArgs.Select)) > 0]

	fkv := map[bool]int{true: 1, false: 0}[len(strings.TrimSpace(oneCommand.CommandLineArgs.FilterKey)) > 0 ||
		len(strings.TrimSpace(oneCommand.CommandLineArgs.FilterValue)) > 0]

	fk := map[bool]int{true: 1, false: 0}[len(strings.TrimSpace(oneCommand.CommandLineArgs.FilterKey)) > 0]
	fv := map[bool]int{true: 1, false: 0}[len(strings.TrimSpace(oneCommand.CommandLineArgs.FilterValue)) > 0]

	if (c+d+s > 1) || c+d+s+fkv == 0 || (fk == 1 && fv == 0) || (d == 1 && (d+fv)%2 != 0) {
		return errors.New("**Invalid command line argument : run with --help for more")
	}
	return nil
}
