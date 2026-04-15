package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinydsa-tester/internal/helpers"
)

func s03DoublyLinkedListTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "doubly-linked-list",
		Timeout:     30 * time.Second,
		TestFunc:    testS03DoublyLinkedList,
		CompileStep: autoCompileStep("TestS03", "test_s03", "test_s03", "testS03"),
	}
}

func testS03DoublyLinkedList(harness *test_case_harness.TestCaseHarness) error {
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

	// --- addFirst ---

	addFirstTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"size_after_3_add_first", "3", "addFirst 3 次后 size == 3"},
		{"to_array_after_add_first", "10,20,30", "addFirst 后 toArray == [10,20,30]"},
		{"to_array_reverse_after_add_first", "30,20,10", "addFirst 后 toArrayReverse == [30,20,10]"},
	}

	for _, tc := range addFirstTests {
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
		{"to_array_reverse_after_add_last", "30,20,10", "addLast 后 toArrayReverse == [30,20,10]"},
	}

	for _, tc := range addLastTests {
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
		{"to_array_reverse_mixed", "40,30,20,10", "混合 addFirst/addLast toArrayReverse 正确"},
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
		{"to_array_reverse_after_remove_first", "30,20", "removeFirst 后 toArrayReverse == [30,20]"},
		{"remove_first_empty", "-1", "空链表 removeFirst 返回 -1"},
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
		{"to_array_reverse_after_remove_last", "20,10", "removeLast 后 toArrayReverse == [20,10]"},
		{"remove_last_empty", "-1", "空链表 removeLast 返回 -1"},
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
		{"to_array_reverse_after_alternate", "30,20", "交替删除后 toArrayReverse == [30,20]"},
		{"size_after_alternate", "2", "交替删除后 size == 2"},
	}

	for _, tc := range alternateTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- removeLast 连续至空 ---

	removeAllTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"remove_last_3", "30", "连续 removeLast 第 1 次返回 30"},
		{"remove_last_2", "20", "连续 removeLast 第 2 次返回 20"},
		{"remove_last_1", "10", "连续 removeLast 第 3 次返回 10"},
		{"size_after_remove_all", "0", "全部删除后 size == 0"},
		{"remove_last_when_empty", "-1", "删空后 removeLast 返回 -1"},
	}

	for _, tc := range removeAllTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 删空再添加（哨兵不被破坏） ---

	refillTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"size_after_clear", "0", "删空后 size == 0"},
		{"to_array_after_refill", "88,99", "重新添加后 toArray == [88,99]"},
		{"to_array_reverse_after_refill", "99,88", "重新添加后 toArrayReverse == [99,88]"},
		{"size_after_refill", "2", "重新添加后 size == 2"},
	}

	for _, tc := range refillTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- toArray / toArrayReverse 空链表 ---

	emptyTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_empty", "", "空链表 toArray 返回空"},
		{"to_array_reverse_empty", "", "空链表 toArrayReverse 返回空"},
	}

	for _, tc := range emptyTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	total := len(basicTests) + len(addFirstTests) + len(addLastTests) + len(mixedTests) +
		len(getTests) + len(removeFirstTests) + len(removeLastTests) + len(alternateTests) +
		len(removeAllTests) + len(refillTests) + len(emptyTests)
	logger.Successf("All %d S03 tests passed!", total)
	return nil
}
