package grok

import (
	"fmt"
	"strings"
)

type patternMap map[string]*grokPattern

// resolve references inside a pattern so that all substitutions are added
// in the correct order.
func (knownPatterns *patternMap) resolve(key, pattern string, newPatterns map[string]string) error {
	for _, keys := range namedReference.FindAllStringSubmatch(pattern, -1) {
		names := strings.Split(keys[1], ":")
		refKey := names[0]

		if _, refKeyCompiled := (*knownPatterns)[refKey]; !refKeyCompiled {
			refPattern, refKeyFound := newPatterns[refKey]
			if !refKeyFound {
				return fmt.Errorf("no pattern found for %%{%s}", refKey)
			}
			knownPatterns.resolve(refKey, refPattern, newPatterns)
		}
	}
	return knownPatterns.add(key, pattern)
}

// add a list of patterns to the map
func (knownPatterns *patternMap) addList(newPatterns map[string]string) error {
	for key, pattern := range newPatterns {
		if _, alreadyCompiled := (*knownPatterns)[key]; alreadyCompiled {
			continue
		}
		if err := knownPatterns.resolve(key, pattern, newPatterns); err != nil {
			return err
		}
	}

	return nil
}

// add a single pattern to the map
func (knownPatterns *patternMap) add(name, pattern string) error {
	p, err := newPattern(pattern, *knownPatterns)
	if err != nil {
		return err
	}

	(*knownPatterns)[name] = p
	return nil
}
