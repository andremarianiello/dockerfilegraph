package dockerfile2dot

import (
	"github.com/awalterschulze/gographviz"
)

// BuildDotFile builds a GraphViz .dot file from a Google Cloud Build configuration
func BuildDotFile(simplifiedDockerfile SimplifiedDockerfile) string {
	graph := gographviz.NewEscape()
	graph.SetName("G")
	graph.SetDir(true)
	graph.AddAttr("G", "splines", "ortho")
	graph.AddAttr("G", "rankdir", "LR")
	graph.AddAttr("G", "nodesep", "1")

	for index, stage := range simplifiedDockerfile.Stages {
		attrs := map[string]string{
			"label": "\"" + getStageLabel(stage) + "\"",
			"shape": "Mrecord",
			"width": "2",
		}

		// Color the last stage, because it is the default build target
		if index == len(simplifiedDockerfile.Stages)-1 {
			attrs["style"] = "filled"
			attrs["fillcolor"] = "grey90"
		}

		graph.AddNode("G", stage.ID, attrs)

		for _, waitForStageID := range stage.WaitFor {
			if waitForStageID == "" {
				continue
			}

			graph.AddEdge(
				getRealStageID(simplifiedDockerfile, waitForStageID),
				stage.ID,
				true,
				nil,
			)
		}
	}

	return graph.String()
}

func getStageLabel(stage Stage) string {
	if stage.Name != "" {
		return stage.Name
	}
	return stage.ID
}

func getRealStageID(simplifiedDockerfile SimplifiedDockerfile, stageID string) string {
	for _, stage := range simplifiedDockerfile.Stages {
		if stageID == stage.ID || stageID == stage.Name {
			return stage.ID
		}
	}
	// This should never happen
	return ""
}
