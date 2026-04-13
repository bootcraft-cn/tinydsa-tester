package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinydsa-tester/internal/helpers"
)

func s02SinglyLinkedListTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "singly-linked-list",
		Timeout:     30 * time.Second,
		TestFunc:    testS02SinglyLinkedList,
		CompileStep: autoCompileStep("TestS02", "test_s02", "test_s02", "testS02"),
	}
}

func testS02SinglyLinkedList(harness *test_case_harness.TestCaseHarness) error {
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
		{"to_array_mixed", "10,20,30", "混合 addFirst/addLast 顺序正确"},
		{"size_mixed", "3", "混合插入后 size == 3"},
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
		{"get_2", "30", "get(2) == 30"},
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
		{"remove_last_empty", "-1", "空链表 removeLast 返回 -1"},
	}

	for _, tc := range removeLastTests {
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

	// --- 删空再添加 ---

	refillTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"size_after_clear", "0", "删空后 size == 0"},
		{"get_after_refill", "99", "重新添加后 get(0) == 99"},
		{"size_after_refill", "1", "重新添加后 size == 1"},
	}

	for _, tc := range refillTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- toArray 空链表 ---

	emptyTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"to_array_empty", "", "空链表 toArray 返回空"},
	}

	for _, tc := range emptyTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	total := len(basicTests) + len(addFirstTests) + len(addLastTests) + len(mixedTests) +
		len(getTests) + len(removeFirstTests) + len(removeLastTests) + len(removeAllTests) +
		len(refillTests) + len(emptyTests)
	logger.Successf("All %d S02 tests passed!", total)
	return nil
}
