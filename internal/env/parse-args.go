package env

import "flag"

func parseArgs() map[string]string {
	names := []string{"c"}
	values := make([]string, len(names))

	for i := range names {
		flag.StringVar(&values[i], names[i], "", "flag "+names[i])
	}

	flag.Parse()

	args := make(map[string]string)
	for i := range names {
		if len(values[i]) > 0 {
			args[names[i]] = values[i]
		}
	}

	return args
}
