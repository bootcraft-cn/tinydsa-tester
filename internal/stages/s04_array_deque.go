package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinydsa-tester/internal/helpers"
)

func s04ArrayDequeTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "array-deque",
		Timeout:     30 * time.Second,
		TestFunc:    testS04ArrayDeque,
		CompileStep: autoCompileStep("TestS04", "test_s04", "test_s04", "testS04"),
	}
}

func testS04ArrayDeque(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir
	lang := harness.DetectedLang

	r := runner.Run(workDir, lang.RunCmd, lang.RunArgs...).
		WithTimeout(10 * time.Second).
		WithLogger(logger).
		Execute().
		Exit(0)

	if err := r.Error(); err != nil {
		return fmt.Errorf("test driver failed: %v", err)
	}

	results := helpers.ParseStructuredOutput(string(r.Result().Stdout))

	// --- 基本操作 ---

	basicTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"initial_size", "0", "初始 size == 0"},
	}

	for _, tc := range basicTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- addLast ---

	addLastTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_after_add_last", "10,20,30", "addLast 后 toArray == [10,20,30]"},
		{"size_after_add_last", "3", "addLast 后 size == 3"},
	}

	for _, tc := range addLastTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- addFirst ---

	addFirstTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_after_add_first", "10,20,30", "addFirst 后 toArray == [10,20,30]"},
		{"size_after_add_first", "3", "addFirst 后 size == 3"},
	}

	for _, tc := range addFirstTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 混合插入 ---

	mixedTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_mixed", "10,20,30,40", "混合 addFirst/addLast toArray 正确"},
		{"size_mixed", "4", "混合插入后 size == 4"},
	}

	for _, tc := range mixedTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- get ---

	getTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"get_0", "10", "get(0) == 10"},
		{"get_1", "20", "get(1) == 20"},
		{"get_3", "40", "get(3) == 40"},
		{"get_out_of_bounds", "-1", "get 越界返回 -1"},
		{"get_negative", "-1", "get 负索引返回 -1"},
	}

	for _, tc := range getTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- removeFirst ---

	removeFirstTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"remove_first_val", "10", "removeFirst 返回 10"},
		{"size_after_remove_first", "2", "removeFirst 后 size == 2"},
		{"to_array_after_remove_first", "20,30", "removeFirst 后 toArray == [20,30]"},
		{"remove_first_empty", "-1", "空队列 removeFirst 返回 -1"},
	}

	for _, tc := range removeFirstTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- removeLast ---

	removeLastTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"remove_last_val", "30", "removeLast 返回 30"},
		{"size_after_remove_last", "2", "removeLast 后 size == 2"},
		{"to_array_after_remove_last", "10,20", "removeLast 后 toArray == [10,20]"},
		{"remove_last_empty", "-1", "空队列 removeLast 返回 -1"},
	}

	for _, tc := range removeLastTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 交替 removeFirst/removeLast ---

	alternateTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"alternate_remove_first", "10", "交替删除 removeFirst 返回 10"},
		{"alternate_remove_last", "40", "交替删除 removeLast 返回 40"},
		{"to_array_after_alternate", "20,30", "交替删除后 toArray == [20,30]"},
		{"size_after_alternate", "2", "交替删除后 size == 2"},
	}

	for _, tc := range alternateTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 环绕测试 ---

	wraparoundTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_wraparound", "1,5,10,20,30", "环绕后 toArray == [1,5,10,20,30]"},
		{"size_wraparound", "5", "环绕后 size == 5"},
		{"get_0_wraparound", "1", "环绕 get(0) == 1"},
		{"get_4_wraparound", "30", "环绕 get(4) == 30"},
	}

	for _, tc := range wraparoundTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 扩容测试（addLast） ---

	resizeTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_after_resize", "1,2,3,4,5,6,7,8,9", "addLast 扩容后 toArray 正确"},
		{"size_after_resize", "9", "扩容后 size == 9"},
		{"get_8_after_resize", "9", "扩容后 get(8) == 9"},
	}

	for _, tc := range resizeTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 扩容测试（addFirst） ---

	resizeFirstTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_add_first_resize", "9,8,7,6,5,4,3,2,1", "addFirst 扩容后 toArray 正确"},
		{"size_add_first_resize", "9", "addFirst 扩容后 size == 9"},
	}

	for _, tc := range resizeFirstTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 删空再添加 ---

	refillTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"size_after_clear", "0", "删空后 size == 0"},
		{"to_array_after_refill", "88,99", "重新添加后 toArray == [88,99]"},
		{"size_after_refill", "2", "重新添加后 size == 2"},
	}

	for _, tc := range refillTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- toArray 空队列 ---

	emptyTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_empty", "", "空队列 toArray 返回空"},
	}

	for _, tc := range emptyTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	total := len(basicTests) + len(addLastTests) + len(addFirstTests) + len(mixedTests) +
		len(getTests) + len(removeFirstTests) + len(removeLastTests) + len(alternateTests) +
		len(wraparoundTests) + len(resizeTests) + len(resizeFirstTests) + len(refillTests) + len(emptyTests)
	logger.Successf("All %d S04 tests passed!", total)
	return nil
}
