package ping

import (
	"fmt"
	"slices"
	"strings"
	"sync"
)

type metricElem struct {
	labels []string
	value  int
}

type PrometheusMetric struct {
	Name       string
	Help       string
	Type       string // gauge | counter
	LabelNames []string
	elems      []metricElem

	mu sync.Mutex
}

func (m *PrometheusMetric) SetValue(labels []string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, elem := range m.elems {
		if slices.Equal(elem.labels, labels) {
			m.elems[i].value = value
			return
		}
	}

	m.elems = append(m.elems, metricElem{
		labels: labels,
		value:  value,
	})
}

func (m *PrometheusMetric) AddValue(labels []string, value int) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, elem := range m.elems {
		if slices.Equal(elem.labels, labels) {
			m.elems[i].value = elem.value + value
			return
		}
	}

	m.elems = append(m.elems, metricElem{
		labels: labels,
		value:  value,
	})
}

func (m *PrometheusMetric) Delete(labels []string) {
	m.mu.Lock()
	defer m.mu.Unlock()

	for i, elem := range m.elems {
		if slices.Equal(elem.labels, labels) {
			m.elems = slices.Delete(m.elems, i, i+1)
			return
		}
	}
}

func (m *PrometheusMetric) GetString() string {
	m.mu.Lock()
	defer m.mu.Unlock()

	lines := []string{
		fmt.Sprintf("# HELP %v %v", m.Name, m.Help),
		fmt.Sprintf("# TYPE %v %v", m.Name, m.Type),
	}

	for _, elem := range m.elems {
		labelPairs := []string{}

		for i, labelName := range m.LabelNames {
			labelPairs = append(labelPairs,
				fmt.Sprintf(`%v="%v"`, labelName, elem.labels[i]),
			)
		}

		lines = append(lines,
			fmt.Sprintf("%v{%v} %v", m.Name, strings.Join(labelPairs, " "), elem.value),
		)
	}

	lines = append(lines, "")
	return strings.Join(lines, "\n")
}
