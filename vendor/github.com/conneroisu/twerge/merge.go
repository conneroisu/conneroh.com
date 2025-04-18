package twerge

import (
	"fmt"
	"maps"
	"regexp"
	"slices"
	"strings"
	"sync"
)

// ClassMapStr is a map of class strings to their generated class names
// This variable can be populated by code generation or manually
// It is protected by mapMutex for concurrent access
var (
	// Merge is the default template merger
	// It takes a space-delimited string of TailwindCSS classes and returns a merged string
	// It also adds the merged class to the ClassMapStr when used
	// It will quickly return the generated class name from ClassMapStr if available
	Merge = createTwMerge(nil, nil)

	ClassMapStr = make(map[string]string)

	// GenClassMergeStr is a map of merged class strings to their generated class names
	// This variable can be populated by code generation or manually
	// It is protected by mapMutex for concurrent access
	GenClassMergeStr = make(map[string]string)

	// mapMutex protects ClassMapStr for concurrent access
	mapMutex sync.RWMutex

	classID                int
	arbitraryPropertyRegex = regexp.MustCompile(`^\[(.+)\]$`)
)

// classMap is a mapping of original class strings to generated class names
type (
	// twMergeFn is the type of the template merger.
	twMergeFn func(classes string) string
	// splitModifiersFn is the type of the function used to split modifiers
	splitModifiersFn = func(string) (
		baseClass string,
		modifiers []string,
		hasImportant bool,
		maybePostfixModPosition int,
	)
	// getClassGroupIDFn returns the class group id for a given class
	getClassGroupIDFn func(string) (isTwClass bool, groupId string)
)

// It returns a short unique CSS class name from the merged classes.
//
// If the class name already exists, it will return the existing class name.
//
// If the class name does not exist, it will generate a new class name and return it.
func It(classes string) string {
	mapMutex.RLock()
	if className, exists := ClassMapStr[classes]; exists {
		mapMutex.RUnlock()
		return className
	}
	mapMutex.RUnlock()

	// First, merge the classes
	merged := Merge(classes)

	// Store the mapping
	mapMutex.Lock()
	classname := fmt.Sprintf("tw-%d", classID)
	ClassMapStr[classes] = classname
	GenClassMergeStr[classname] = merged
	classID++
	mapMutex.Unlock()

	return classname
}

// If returns a short unique CSS class name from the merged classes taking an additional boolean parameter.
//
// If the class name already exists, it will return the existing class name.
//
// If the class name does not exist, it will generate a new class name and return it.
func If(cond bool, trueClass, falseClass string) string {
	if className, exists := ClassMapStr[trueClass]; exists && cond {
		return className
	}
	if className, exists := ClassMapStr[falseClass]; exists && !cond {
		return className
	}

	trueEval := It(trueClass)
	falseEval := It(falseClass)
	if cond {
		return trueEval
	}
	return falseEval
}

// createTwMerge creates a new template merger
func createTwMerge(
	config *config,
	cache icache,
) twMergeFn {
	var (
		fnToCall        twMergeFn
		splitModifiers  splitModifiersFn
		getClassGroupID getClassGroupIDFn
		mergeClassList  func(classList string) string
	)

	merger := func(classes string) string {
		classList := strings.TrimSpace(classes)
		if classList == "" {
			return ""
		}

		// Check if we've seen this class list before in the cache
		cached := cache.Get(classList)
		if cached != "" {
			return cached
		}

		// Merge the classes
		merged := mergeClassList(classList)
		cache.Set(classList, merged)

		// Add to ClassMapStr for lookup by other functions
		if classList != merged {
			mapMutex.Lock()
			className := fmt.Sprintf("tw-%d", classID)
			ClassMapStr[classList] = className
			GenClassMergeStr[className] = merged
			classID++
			mapMutex.Unlock()
		}

		return merged
	}

	init := func(classes string) string {
		if config == nil {
			config = defaultConfig
		}
		if cache == nil {
			cache = newCache(config.MaxCacheSize)
		}

		splitModifiers = makeSplitModifiers(config)

		getClassGroupID = makeGetClassGroupID(config)

		mergeClassList = makeMergeClassList(config, splitModifiers, getClassGroupID)

		fnToCall = merger
		return fnToCall(classes)
	}

	fnToCall = init
	return func(classes string) string {
		return fnToCall(classes)
	}
}

// makeMergeClassList creates a function that merges a class list
func makeMergeClassList(
	conf *config,
	splitModifiers splitModifiersFn,
	getClassGroupID getClassGroupIDFn,
) func(classList string) string {
	return func(classList string) string {
		classes := strings.Split(strings.TrimSpace(classList), " ")
		unqClasses := make(map[string]string, len(classes))
		resultClassList := ""

		for _, class := range classes {
			baseClass, modifiers, hasImportant, postFixMod := splitModifiers(class)

			// there is a postfix modifier -> text-lg/8
			if postFixMod != -1 {
				baseClass = baseClass[:postFixMod]
			}
			isTwClass, groupID := getClassGroupID(baseClass)
			if !isTwClass {
				resultClassList += class + " "
				continue
			}
			// we have to sort the modifiers bc hover:focus:bg-red-500 == focus:hover:bg-red-500
			modifiers = sortModifiers(modifiers)
			if hasImportant {
				modifiers = append(modifiers, "!")
			}
			unqClasses[groupID+strings.Join(modifiers, string(conf.ModifierSeparator))] = class

			conflicts := conf.ConflictingClassGroups[groupID]
			if conflicts == nil {
				continue
			}
			for _, conflict := range conflicts {
				// erase the conflicts with the same modifiers
				unqClasses[conflict+strings.Join(modifiers, string(conf.ModifierSeparator))] = ""
			}
		}

		for _, class := range unqClasses {
			if class == "" {
				continue
			}
			resultClassList += class + " "
		}
		return strings.TrimSpace(resultClassList)
	}

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

// makeSplitModifiers creates a function that splits modifiers
func makeSplitModifiers(config *config) splitModifiersFn {

	return func(className string) (string, []string, bool, int) {
		separator := config.ModifierSeparator
		modifiers := []string{}
		modifierStart := 0
		bracketDepth := 0
		// used for bg-red-500/50 (50% opacity)
		maybePostfixModPosition := -1

		for i := range len(className) {
			char := rune(className[i])

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
					modifiers = append(modifiers, className[modifierStart:i])
					modifierStart = i + 1
					continue
				}

				if char == config.PostfixModifier {
					maybePostfixModPosition = i
				}
			}
		}

		baseClassWithImportant := className[modifierStart:]
		if len(baseClassWithImportant) == 0 {
			return "", nil, false, -1
		}
		hasImportant := baseClassWithImportant[0] == byte(config.ImportantModifier)

		var baseClass string
		if hasImportant {
			baseClass = baseClassWithImportant[1:]
		} else {
			baseClass = baseClassWithImportant
		}

		// fix case where there is modifier & maybePostfix which causes maybePostfix to be beyond size of baseClass!
		if maybePostfixModPosition != -1 && maybePostfixModPosition > modifierStart {
			maybePostfixModPosition -= modifierStart
		} else {
			maybePostfixModPosition = -1
		}

		return baseClass, modifiers, hasImportant, maybePostfixModPosition

	}
}

func getMapping() map[string]string {
	mapMutex.RLock()
	defer mapMutex.RUnlock()

	// Create a copy to avoid concurrent map access issues
	mapping := make(map[string]string, len(ClassMapStr))
	maps.Copy(mapping, ClassMapStr)

	return mapping
}

// makeGetClassGroupID returns a getClassGroupIdfn
func makeGetClassGroupID(conf *config) getClassGroupIDFn {
	var getClassGroupIDRecursive func(
		classParts []string,
		i int,
		classMap *classPart,
	) (isTwClass bool, groupId string)

	getClassGroupIDRecursive = func(
		classParts []string,
		i int,
		classMap *classPart,
	) (isTwClass bool, groupId string) {
		if i >= len(classParts) {
			if classMap.ClassGroupID != "" {
				return true, classMap.ClassGroupID
			}

			return false, ""
		}

		if classMap.NextPart != nil {
			nextClassMap := classMap.NextPart[classParts[i]]
			isTw, id := getClassGroupIDRecursive(classParts, i+1, &nextClassMap)
			if isTw {
				return isTw, id
			}
		}

		if len(classMap.Validators) > 0 {
			remainingClass := strings.Join(classParts[i:], string(conf.ClassSeparator))

			for _, validator := range classMap.Validators {
				if validator.Fn(remainingClass) {
					return true, validator.ClassGroupID
				}
			}

		}
		return false, ""
	}

	getGroupIDForArbitraryProperty := func(class string) (bool, string) {
		if arbitraryPropertyRegex.MatchString(class) {
			arbitraryPropertyClassName := arbitraryPropertyRegex.FindStringSubmatch(class)[1]
			property := arbitraryPropertyClassName[:strings.Index(arbitraryPropertyClassName, ":")]

			if property != "" {
				// two dots here because one dot is used as prefix for class groups in plugins
				return true, "arbitrary.." + property
			}
		}

		return false, ""
	}

	return func(baseClass string) (isTwClass bool, groupdId string) {
		classParts := strings.Split(baseClass, string(conf.ClassSeparator))
		// remove first element if empty for things like -px-4
		if len(classParts) > 0 && classParts[0] == "" {
			classParts = classParts[1:]
		}
		isTwClass, groupID := getClassGroupIDRecursive(classParts, 0, &conf.ClassGroups)
		if isTwClass {
			return isTwClass, groupID
		}

		return getGroupIDForArbitraryProperty(baseClass)
	}

}
