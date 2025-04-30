package twerge

import (
	"slices"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

// CacheValue contains the value of a cache entry.
//
// It is used to store the generated and merged classes.
//
// As twerge is meant to be used statically, aka at build/compile time,
// it is trying to maxmimize performance at runtime.
type CacheValue struct {
	// Generated is the generated class. It is a short unique CSS class name.
	//
	// Example: tw-123
	Generated string
	// Merged is the merged class. It is the final class name that is used in the CSS.
	//
	// Example: min-h-screen bg-gray-50 text-gray-900 flex flex-col
	Merged string
}

var defaultGenerator atomic.Pointer[Generator]

func init() {
	defaultGenerator.Store(New(newDefaultHandler()))
}

// Default returns the default [Generator].
func Default() *Generator { return defaultGenerator.Load() }

// SetDefault sets the default [Generator].
func SetDefault(g *Generator) { defaultGenerator.Store(g) }

// If returns a short unique CSS class name from the merged classes taking an additional boolean parameter.
func If(
	ok bool,
	trueClass string,
	falseClass string,
) string {
	trueClass = It(trueClass)
	falseClass = It(falseClass)
	if ok {
		return trueClass
	}
	return falseClass
}

// It returns a short unique CSS class name from the merged classes.
func It(raw string) string {
	return Default().It(raw)
}

// Generator generates all the code needed to use Twerge statically.
//
// At runtime, it uses the statically defined code, if configured, to
// map the class names to the generated class names.
type Generator struct{ Handler Handler }

// New creates a new Generator with the given non-nil Handler.
func New(h Handler) *Generator {
	return &Generator{Handler: h}
}

// Handler is the interface that needs to be implemented to customize the
// behavior of the [Generator].
type Handler interface {
	It(string) string
	Cache() map[string]CacheValue
	SetCache(map[string]CacheValue)
}

// Cache returns the cache of the [Generator].
func (Generator) Cache() map[string]CacheValue {
	return defaultGenerator.Load().Handler.Cache()
}

// Cache returns the cache of the [Generator].
func (g *defaultHandler) Cache() map[string]CacheValue { return g.entries }

// SetCache sets the cache of the [Generator].
func (g *defaultHandler) SetCache(entries map[string]CacheValue) {
	g.entries = entries
}

// It returns a short unique CSS class name from the merged classes.
//
// If the class name already exists, it will return the existing class name.
//
// If the class name does not exist, it will generate a new class name and
// return it.
func (g *Generator) It(classes string) string {
	return g.Handler.It(classes)
}

func newDefaultHandler() *defaultHandler {
	return &defaultHandler{
		entries: make(map[string]CacheValue),
		config:  defaultConfig,
	}

}

type defaultHandler struct {
	config  *config
	entries map[string]CacheValue
	mu      sync.RWMutex
}

func (g *defaultHandler) It(classes string) string {
	// Read Safe Lock
	g.mu.RLock()
	if className, exists := g.entries[classes]; exists {
		g.mu.RUnlock()
		return className.Generated
	}
	g.mu.RUnlock()

	// Write Safe Lock
	g.mu.Lock()
	className := "tw-" + strconv.Itoa(len(g.entries))
	g.entries[classes] = CacheValue{
		Generated: className,
		Merged:    g.merge(classes),
	}
	g.mu.Unlock()

	return className
}

func (g *defaultHandler) merge(classes string) string {
	var (
		uniques = make(map[string]string)
		merged  string
	)
	classesSeq := strings.SplitSeq(strings.TrimSpace(classes), " ")

	for class := range classesSeq {
		var (
			modifiers     []string
			modifierStart int
			bracketDepth  int

			isTwClass bool
			groupID   string
		)
		separator := g.config.ModifierSeparator
		// used for examples like 'bg-red-500/50' (50% opacity)
		postFixMod := -1

		for i := range len(class) {
			char := rune(class[i])

			if char == '[' {
				bracketDepth++
				continue
			}
			if char == ']' {
				bracketDepth--
				continue
			}

			if bracketDepth == 0 {
				if char == separator {
					modifiers = append(
						modifiers,
						class[modifierStart:i],
					)
					modifierStart = i + 1
					continue
				}

				if char == g.config.PostfixModifier {
					postFixMod = i
				}
			}
		}

		base := class[modifierStart:]
		hasImportant := base[0] == byte(g.config.ImportantModifier)

		if hasImportant {
			base = base[1:]
		}

		// if there is modifier & maybePostfix which causes
		// postfixModPos to be beyond size of baseClass
		if postFixMod != -1 && postFixMod > modifierStart {
			postFixMod -= modifierStart
		} else {
			postFixMod = -1
		}

		// there is a postfix modifier -> text-lg/8
		if postFixMod != -1 {
			base = base[:postFixMod]
		}
		isTwClass, groupID = g.getClassGroupID(base)
		if !isTwClass {
			merged += class + " "
			continue
		}
		// sort as hover:focus:bg-red-500 == focus:hover:bg-red-500
		modifiers = sortModifiers(modifiers)
		if hasImportant {
			modifiers = append(modifiers, "!")
		}
		uniques[groupID+strings.Join(
			modifiers,
			string(g.config.ModifierSeparator),
		)] = class

		conflicts := g.config.ConflictingClassGroups[groupID]
		if conflicts == nil {
			continue
		}
		for _, conflict := range conflicts {
			// erase the conflicts with the same modifiers
			uniques[conflict+strings.Join(
				modifiers,
				string(g.config.ModifierSeparator),
			)] = ""
		}
	}

	for _, unique := range uniques {
		if unique == "" {
			continue
		}
		merged += unique + " "
	}
	return strings.TrimSpace(merged)
}

func (g *defaultHandler) getClassGroupIDRecursive(
	classParts []string,
	i int,
	configClassGroups *classPart,
) (isTwClass bool, groupID string) {
	if i >= len(classParts) {
		if configClassGroups.ClassGroupID != "" {
			return true, configClassGroups.ClassGroupID
		}

		return false, ""
	}

	if configClassGroups.NextPart != nil {
		nextClassMap := configClassGroups.NextPart[classParts[i]]
		isTw, id := g.getClassGroupIDRecursive(
			classParts,
			i+1,
			&nextClassMap,
		)
		if isTw {
			return isTw, id
		}
	}

	if len(configClassGroups.Validators) > 0 {
		remainingClass := strings.Join(
			classParts[i:],
			string(g.config.ClassSeparator),
		)

		for _, validator := range configClassGroups.Validators {
			if validator.Fn(remainingClass) {
				return true, validator.ClassGroupID
			}
		}

	}
	return false, ""
}

func (g *defaultHandler) getGroupIDForArbitraryProperty(class string) (bool, string) {
	if arbitraryPropertyRegex.MatchString(class) {
		name := arbitraryPropertyRegex.FindStringSubmatch(class)[1]
		property := name[:strings.Index(name, ":")]

		if property != "" {
			// two dots here because one dot is used as prefix for class groups in plugins
			return true, "arbitrary.." + property
		}
	}

	return false, ""
}

// getClassGroupID returns a boolean and a string
func (g *defaultHandler) getClassGroupID(baseClass string) (bool, string) {
	classParts := strings.Split(baseClass, string(g.config.ClassSeparator))
	// remove first element if empty for things like -px-4
	if len(classParts) > 0 && classParts[0] == "" {
		classParts = classParts[1:]
	}
	isTwClass, groupID := g.getClassGroupIDRecursive(
		classParts,
		0,
		&g.config.ClassGroups,
	)
	if isTwClass {
		return isTwClass, groupID
	}

	return g.getGroupIDForArbitraryProperty(baseClass)
}

// sortModifiers Sorts modifiers according to following schema:
// - Predefined modifiers are sorted alphabetically
// - When an arbitrary variant appears, it must be preserved which modifiers are before and after it
func sortModifiers(modifiers []string) []string {
	if len(modifiers) < 2 {
		return modifiers
	}

	unsortedModifiers := []string{}
	sorted := make([]string, len(modifiers))

	for _, modifier := range modifiers {
		isArbitraryVariant := modifier[0] == '['
		if isArbitraryVariant {
			slices.Sort(unsortedModifiers)
			sorted = append(sorted, unsortedModifiers...)
			sorted = append(sorted, modifier)
			unsortedModifiers = []string{}
			continue
		}
		unsortedModifiers = append(unsortedModifiers, modifier)
	}

	slices.Sort(unsortedModifiers)
	sorted = append(sorted, unsortedModifiers...)

	return sorted
}
