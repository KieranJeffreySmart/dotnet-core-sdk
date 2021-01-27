package dotnetcoresdk

import (
	"regexp"
	"sort"

	"github.com/paketo-buildpacks/packit"
)

type PlanEntryResolver struct {
	logger LogEmitter
}

func NewPlanEntryResolver(logger LogEmitter) PlanEntryResolver {
	return PlanEntryResolver{
		logger: logger,
	}
}

func (r PlanEntryResolver) Resolve(entries []packit.BuildpackPlanEntry) packit.BuildpackPlanEntry {
	sort.Slice(entries, func(i, j int) bool {
		leftSource := entries[i].Metadata["version-source"]
		left, _ := leftSource.(string)

		rightSource := entries[j].Metadata["version-source"]
		right, _ := rightSource.(string)

		return getPriority(left) > getPriority(right)
	})

	chosenEntry := entries[0]

	if chosenEntry.Metadata == nil {
		chosenEntry.Metadata = map[string]interface{}{}
	}

	for _, entry := range entries {
		if entry.Metadata["build"] == true {
			chosenEntry.Metadata["build"] = true
		}
		if entry.Metadata["launch"] == true {
			chosenEntry.Metadata["launch"] = true
		}
	}

	r.logger.Candidates(entries)

	return chosenEntry
}

func getPriority(source string) int {
	var (
		priorities = map[string]int{
			"RUNTIME_VERSION":    4,
			"buildpack.yml":      3,
			"global.json":        2,
			"runtimeconfig.json": 1,
			"":                   -1,
		}
	)

	if match, _ := regexp.MatchString(`.*\.(cs)|(fs)|(vb)proj`, source); match {
		return 1
	}

	return priorities[source]
}
