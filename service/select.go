package service

import (
	"fmt"
	"strings"

	"github.com/itchyny/gojq"
)

type Select struct {
	Filter
	Projection []string
}

func (s *Select) buildSelectProjection(id string) (selectorFormat string, isIdSelected bool) {
	var projection []string = make([]string, 0)

	for _, value := range s.Projection {
		node := strings.TrimSpace(value)
		isId := strings.EqualFold(node, "id")
		if isId {
			isIdSelected = isId
		}

		if len(node) > 0 && !isId {
			projection = append(projection, fmt.Sprintf("%v:.%v", strings.ReplaceAll(node, ".", "_"), node))
		}
	}
	return fmt.Sprintf(".\"%v\" | {%v}", id, strings.Join(projection, ", ")), isIdSelected
}

func buildProjectionResult(isSelectId bool, id string, projectionResult interface{}) interface{} {
	if isSelectId {
		var idMap map[string]interface{} = make(map[string]interface{})
		idMap[id] = projectionResult
		return idMap
	}
	return projectionResult
}

func (s *Select) selectByProjection(searchResults []map[string]interface{}) interface{} {
	var projectionResult []interface{}
	for _, result := range searchResults {
		for key := range result {
			format, isSelectId := s.buildSelectProjection(key)
			query, err := gojq.Parse(format)

			if err != nil {
				fmt.Println(err)
			}

			iter := query.Run(result)

			r, isSelected := iter.Next()
			if isSelected {
				projectionResult = append(projectionResult, buildProjectionResult(isSelectId, key, r))
			}
		}

	}

	return projectionResult
}

func (s *Select) Run() interface{} {

	result, _ := s.Search()
	if s.Projection != nil && len(s.Projection) > 0 {

		return s.selectByProjection(result)
	}

	return result
}
