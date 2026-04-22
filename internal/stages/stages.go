package stages

import (
	"github.com/bootcraft-cn/tester-utils/tester_definition"
)

// GetDefinition returns the TesterDefinition for the tinydsa course.
func GetDefinition() tester_definition.TesterDefinition {
	return tester_definition.TesterDefinition{
		TestCases: []tester_definition.TestCase{
			// Phase 1: 线性结构
			s01DynamicArrayTestCase(),
			s02SinglyLinkedListTestCase(),
			s03DoublyLinkedListTestCase(),
			s04ArrayDequeTestCase(),
			s05StackAndQueueTestCase(),
			// Phase 2: 树与映射
			s06BSTSearchInsertTestCase(),
		},
	}
}

// --- Language rules (4 languages) ---

// allJavaSources lists all Java source files that need to be compiled.
// This is cumulative — every new stage that adds a .java file must be listed here.
var allJavaSources = []string{
	"src/main/java/cn/bootcraft/tinydsa/DynamicArray.java",
	"src/main/java/cn/bootcraft/tinydsa/SinglyLinkedList.java",
	"src/main/java/cn/bootcraft/tinydsa/DoublyLinkedList.java",
	"src/main/java/cn/bootcraft/tinydsa/ArrayDeque.java",
	"src/main/java/cn/bootcraft/tinydsa/Stack.java",
	"src/main/java/cn/bootcraft/tinydsa/Queue.java",
	"src/main/java/cn/bootcraft/tinydsa/BST.java",
}

// javaRule creates a LanguageRule for Java auto-detection.
func javaRule(testDriver string) tester_definition.LanguageRule {
	// Flags: [-encoding UTF-8] + extra source files (excluding first) + test driver
	// Source: first source file (appended last by compileJava)
	// Result: javac -d . -encoding UTF-8 <extra sources> tests/TestSXX.java <Source>
	flags := []string{"-encoding", "UTF-8"}
	// Add all sources except the first (which goes in Source)
	for _, s := range allJavaSources[1:] {
		flags = append(flags, s)
	}
	flags = append(flags, "tests/"+testDriver+".java")
	return tester_definition.LanguageRule{
		DetectFile: "src/main/java/cn/bootcraft/tinydsa/DynamicArray.java",
		Language:   "java",
		Source:     allJavaSources[0],
		Flags:      flags,
		RunCmd:     "java",
		RunArgs:    []string{"-cp", ".", testDriver},
	}
}

// pythonRule creates a LanguageRule for Python auto-detection.
func pythonRule(testDriver string) tester_definition.LanguageRule {
	return tester_definition.LanguageRule{
		DetectFile: "tinydsa/dynamic_array.py",
		Language:   "python",
		Source:     "tinydsa/dynamic_array.py",
		RunCmd:     "python3",
		RunArgs:    []string{"tests/" + testDriver + ".py"},
	}
}

// goRule creates a LanguageRule for Go auto-detection.
func goRule(testDriver string) tester_definition.LanguageRule {
	return tester_definition.LanguageRule{
		DetectFile: "pkg/tinydsa/dynamic_array.go",
		Language:   "go",
		Source:     "./pkg/tinydsa",
		RunCmd:     "go",
		RunArgs:    []string{"run", "tests/" + testDriver + ".go"},
	}
}

// tsRule creates a LanguageRule for TypeScript auto-detection.
// Uses `tsx` directly (pre-installed globally in the runtime image) instead
// of `npx tsx` to avoid a per-run npx download into HOME, which both wastes
// time and historically tripped a tsx@4.21 WASM OOM on Node 18.
func tsRule(testDriver string) tester_definition.LanguageRule {
	return tester_definition.LanguageRule{
		DetectFile: "src/dynamicArray.ts",
		Language:   "typescript",
		Source:     "src/dynamicArray.ts",
		RunCmd:     "tsx",
		RunArgs:    []string{"tests/" + testDriver + ".ts"},
	}
}

// autoCompileStep returns a CompileStep with auto-detection for all 4 languages.
func autoCompileStep(javaDriver, pythonDriver, goDriver, tsDriver string) *tester_definition.CompileStep {
	return &tester_definition.CompileStep{
		Language: "auto",
		AutoDetect: []tester_definition.LanguageRule{
			javaRule(javaDriver),
			pythonRule(pythonDriver),
			goRule(goDriver),
			tsRule(tsDriver),
		},
	}
}
