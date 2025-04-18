package twerge

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"
	"os"
	"regexp"
	"slices"
	"strings"

	"github.com/a-h/templ"
	"github.com/conneroisu/twerge/internal/files"
	"github.com/dave/jennifer/jen"
)

var (
	// DefaultMerger is the default template merger
	DefaultMerger = NewMerger(defaultConfig)
	// DefaultGenerator is the default implementation of the Generater interface.
	DefaultGenerator = NewGenerator(defaultConfig)
	// DefaultCache is the default cache object.
	DefaultCache = NewCache()
	// Merge is the default template merger
	//
	// It takes a space-delimited string of TailwindCSS classes and returns
	// a merged string
	//
	// It also adds the merged class to the ClassMapStr when used
	//
	// It will quickly return the generated class name from ClassMapStr if
	// available
	Merge = DefaultMerger.Merge
	// It returns a short unique CSS class name from the merged classes
	// given a raw class string.
	//
	// It will quickly return the existing class name if the class name already
	// exists in the cache.
	//
	// Example:
	//
	//	templ ExampleComp () {
	//		<div class={twerge.It("bg-red-500 text-white")}>
	//			Hello World
	//		</div>
	//	}
	It = DefaultGenerator.It
	// If returns a short unique CSS class name from the merged classes for a conditional.
	If = DefaultGenerator.If
	// GenCSS generates the CSS classes from the provided components.
	GenCSS = DefaultGenerator.GenCSS

	arbitraryPropertyRegex = regexp.MustCompile(`^\[(.+)\]$`)

	_ Generater = (*Generator)(nil)
	_ Merger    = (*MergerImpl)(nil)
)

const (
	// twergeBeginMarker is the beginning of the section where the generated
	// CSS will be placed
	twergeBeginMarker = "/* twerge:begin */"
	// twergeEndMarker is the end of the section where the generated CSS
	// will be placed
	twergeEndMarker = "/* twerge:end */"
)

// getClassGroupIDFn returns the class group id for a given class
type (
	// cacher is the interface for a LRU cacher
	cacher interface {
		Get(string) (generated string, merged string)
		Set(string, string) string
		All() func(yield func(string, CacheValue) bool)
	}
	// TODO: remove this type
	getClassGroupIDFn func(string) (isTwClass bool, groupId string)
	// Merger is the minimal interface for merging Tailwind CSS classes
	Merger interface {
		Merge(classes string) string
	}
	// MergerImpl is the default template merger
	MergerImpl struct {
		Config *config
		Cache  cacher
		Logger *slog.Logger
	}
	// Generater is the minimal interface for generating Tailwind CSS
	// classes.
	Generater interface {
		It(classes string) string
		If(cond bool, trueClass, falseClass string) string
		Cached() cacher
		GenCSS(ctx context.Context,
			classPath, twPath, templPath string,
			components ...templ.Component,
		) error
	}
	// Generator is the default implementation of the Generator interface.
	Generator struct {
		Config *config
		Merger Merger
		Cache  cacher
		Logger *slog.Logger
	}
	// AmnesicMerger is the default template merger for generating Tailwind CSS classes.
	//
	// It does not modify it's cache.
	AmnesicMerger struct {
		config *config
	}
)

// NewMerger creates a new template merger
func NewMerger(
	config *config,
) *MergerImpl {
	if config == nil {
		config = defaultConfig
	}
	m := &MergerImpl{
		Config: config,
		Cache:  DefaultCache,
	}
	return m
}

// Merge is the default template merger's Merge function.
//
// This functions definition makes the template merger implement the Merger
// interface.
func (m *MergerImpl) Merge(
	classes string,
) string {
	if strings.Contains(classes, "  ") {
		panic("two consecutive spaces are not allowed in class names: " + classes)
	}

	classList := strings.TrimSpace(classes)
	if classList == "" {
		return ""
	}

	// Check if we've seen this class list before in the cache
	_, merged := m.Cache.Get(classList)
	if merged != "" {
		return merged
	}

	mergeClassList := makeMergeClassList(
		m.Config,
		makeGetClassGroupID(m.Config),
	)

	// Merge the classes
	merged = mergeClassList(classList)

	_ = m.Cache.Set(classList, merged)

	return merged
}

// NewGenerator creates a new Generator.
func NewGenerator(config *config) *Generator {
	if config == nil {
		config = defaultConfig
	}
	return &Generator{
		Config: config,
		Merger: &AmnesicMerger{
			config: config,
		},
		Cache:  DefaultCache,
		Logger: slog.Default(),
	}
}

// makeMergeClassList creates a function that merges a class list
func makeMergeClassList(
	conf *config,
	getClassGroupID getClassGroupIDFn,
) func(classList string) string {
	if conf == nil {
		panic("conf is nil")
	}
	separator := conf.ModifierSeparator

	return func(classList string) string {
		classes := strings.Split(strings.TrimSpace(classList), " ")
		unqClasses := make(map[string]string, len(classes))
		resultClassList := ""

		for _, class := range classes {
			modifiers := []string{}
			modifierStart := 0
			bracketDepth := 0
			// used for bg-red-500/50 (50% opacity)
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
						modifiers = append(modifiers, class[modifierStart:i])
						modifierStart = i + 1
						continue
					}

					if char == conf.PostfixModifier {
						postFixMod = i
					}
				}
			}

			// TODO: Add panic here if len(className) == 0 (Meaning that there may be instances of two spaces in a row in the class string)
			baseClassWithImportant := class[modifierStart:]
			if len(baseClassWithImportant) == 0 {
				return ""
			}
			hasImportant := baseClassWithImportant[0] == byte(conf.ImportantModifier)

			var baseClass string
			if hasImportant {
				baseClass = baseClassWithImportant[1:]
			} else {
				baseClass = baseClassWithImportant
			}

			// fix case where there is modifier & maybePostfix which causes maybePostfix to be beyond size of baseClass!
			if postFixMod != -1 && postFixMod > modifierStart {
				postFixMod -= modifierStart
			} else {
				postFixMod = -1
			}

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

// makeGetClassGroupID returns a getClassGroupIdfn
func makeGetClassGroupID(
	conf *config,
) getClassGroupIDFn {
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
		isTwClass, groupID := getClassGroupIDRecursive(conf, classParts, 0, &conf.ClassGroups)
		if isTwClass {
			return isTwClass, groupID
		}

		return getGroupIDForArbitraryProperty(baseClass)
	}

}

func getClassGroupIDRecursive(
	conf *config,
	classParts []string,
	i int,
	classMap *classPart,
) (bool, string) {
	if i >= len(classParts) {
		if classMap.ClassGroupID != "" {
			return true, classMap.ClassGroupID
		}

		return false, ""
	}

	if classMap.NextPart != nil {
		nextClassMap := classMap.NextPart[classParts[i]]
		isTw, id := getClassGroupIDRecursive(conf, classParts, i+1, &nextClassMap)
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

// It returns a short unique CSS class name from the merged classes.
//
// If the class name already exists, it will return the existing class name.
//
// If the class name does not exist, it will generate a new class name and return it.
func (g *Generator) It(classes string) string {
	generated, _ := g.Cache.Get(classes)
	if generated != "" {
		return generated
	}
	classList := strings.TrimSpace(classes)
	if classList == "" {
		return ""
	}
	merged := g.Merger.Merge(classList)

	generated = g.Cache.Set(classList, merged)
	g.Logger.Debug("Generated class", slog.String("class", generated))

	return generated
}

// Merge is the default template merger's Merge function.
//
// This functions definition makes the template merger implement the Merger
// interface.
func (m *AmnesicMerger) Merge(classes string) string {
	if strings.Contains(classes, "  ") {
		panic("two spaces are not allowed in class names: " + classes)
	}

	mergeClassList := makeMergeClassList(
		m.config,
		makeGetClassGroupID(m.config),
	)
	// Merge the classes
	merged := mergeClassList(classes)

	return merged
}

// If returns a short unique CSS class name from the merged classes.
//
// If the class name already exists, it will return the existing class name.
//
// If the class name does not exist, it will generate a new class name and return it.
func (g *Generator) If(cond bool, trueClass, falseClass string) string {
	trueEval := g.It(trueClass)
	falseEval := g.It(falseClass)
	if cond {
		return trueEval
	}
	return falseEval
}

// GenCSS generates the CSS classes from the provided components.
func (g *Generator) GenCSS(
	ctx context.Context,
	classPath, inputCSSPath, templPath string,
	components ...templ.Component,
) error {
	for _, c := range components {
		err := c.Render(ctx, io.Discard)
		if err != nil {
			return err
		}
	}
	err := generateClassMapCode(g, classPath)
	if err != nil {
		return err
	}
	err = generateTailwind(g, inputCSSPath)
	if err != nil {
		return err
	}
	err = generateHTML(g, templPath)
	if err != nil {
		return err
	}
	return nil
}

// GenerateTailwind creates an input CSS file for the Tailwind CLI
// that includes all the @apply directives from the provided class map.
//
// This is useful for building a production CSS file with Tailwind's CLI.
//
// The marker is used to identify the start and end of the @apply directives generated
// by Twerge.
func generateTailwind(
	gen Generater,
	cssPath string,
) error {
	var builder strings.Builder
	baseContent, err := os.ReadFile(cssPath)
	if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error reading input file: %w", err)
	}
	if os.IsNotExist(err) {
		baseContent = []byte(`@tailwind base;
@tailwind components;
@tailwind utilities;

/* twerge:start */
/* twerge:end */
`)
	}
	var gendClasses []string
	for _, key := range gen.Cached().All() {
		if slices.Contains(gendClasses, key.Generated) {
			continue
		}
		builder.WriteString(".")
		builder.WriteString(key.Generated)
		builder.WriteString(" { \n\t@apply ")
		builder.WriteString(key.Merged)
		builder.WriteString("; \n}\n")
	}
	cssContent := builder.String()

	// Add to file content
	newContent, err := files.ReplaceBetweenMarkers(
		baseContent,
		[]byte(cssContent),
		twergeBeginMarker,
		twergeEndMarker,
	)
	if err != nil {
		return fmt.Errorf("error replacing twerge content between markers: %w", err)
	}

	// Write to output path
	err = os.WriteFile(cssPath, newContent, 0644)
	if err != nil {
		return fmt.Errorf("error writing output file: %w", err)
	}

	return nil
}

// generateHTML creates a .templ file that can be used to generate a CSS file
// with the provided class map.
func generateHTML(
	gen Generater,
	htmlPath string,
) error {
	var buf bytes.Buffer

	for _, v := range gen.Cached().All() {
		buf.WriteString("<div class=\"")
		buf.WriteString(v.Generated)
		buf.WriteString("\"></div>\n")
	}
	buf.WriteString("}")

	err := os.WriteFile(htmlPath, buf.Bytes(), 0644)
	if err != nil {
		return fmt.Errorf("error writing .templ file: %w", err)
	}

	return nil
}

// GenerateClassMapCode generates Go code for a variable containing the class mapping
func generateClassMapCode(
	gen Generater,
	classPath string,
) error {
	packageName := files.GetPackageName(classPath)

	// Create a new file
	f := jen.NewFile(packageName)

	// Add a package comment
	f.PackageComment("Code generated by twerge. DO NOT EDIT.")

	// Create the ClassMapStr variable
	f.Var().Id("ClassMapStr").Op("=").Map(
		jen.String(),
	).String().Values(jen.DictFunc(func(d jen.Dict) {
		var (
			values []string
			key    string
		)

		for key = range gen.Cached().All() {
			values = append(values, key)
		}

		for _, key := range values {
			generated, _ := gen.Cached().Get(key)
			d[jen.Lit(key)] = jen.Lit(generated)
		}
	}))

	f.Var().Id("MergedMapStr").Op("=").Map(
		jen.String(),
	).String().Values(jen.DictFunc(func(d jen.Dict) {
		var (
			values []string
			value  CacheValue
		)

		for _, value = range gen.Cached().All() {
			values = append(values, value.Merged)
		}

		for _, k := range values {
			generated, _ := gen.Cached().Get(k)
			d[jen.Lit(k)] = jen.Lit(generated)
		}
	}))

	return files.WriteJen(f, classPath)
}

// Cached returns the cacher for the Generator.
func (g *Generator) Cached() cacher {
	return g.Cache
}
