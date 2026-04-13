package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinydsa-tester/internal/helpers"
)

func s01DynamicArrayTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "dynamic-array",
		Timeout:     30 * time.Second,
		TestFunc:    testS01DynamicArray,
		CompileStep: autoCompileStep("TestS01", "test_s01", "test_s01", "testS01"),
	}
}

func testS01DynamicArray(harness *test_case_harness.TestCaseHarness) error {
	logger := harness.Logger
	workDir := harness.SubmissionDir
	lang := harness.DetectedLang

	// Run test driver
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
		{"initial_capacity", "8", "初始 capacity == 8"},
	}

	for _, tc := range basicTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- add + get ---

	addGetTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"size_after_3_adds", "3", "add 3 次后 size == 3"},
		{"get_0", "10", "get(0) == 10"},
		{"get_1", "20", "get(1) == 20"},
		{"get_2", "30", "get(2) == 30"},
	}

	for _, tc := range addGetTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- get 越界 ---

	oobTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"get_out_of_bounds", "-1", "get 越界返回 -1"},
		{"get_negative", "-1", "get 负索引返回 -1"},
	}

	for _, tc := range oobTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- set ---

	setTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"get_after_set", "99", "set(1, 99) 后 get(1) == 99"},
		{"size_after_oob_set", "3", "越界 set 不改变 size"},
	}

	for _, tc := range setTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 扩容 ---

	expandTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"cap_before_expand", "8", "8 个元素时 capacity == 8"},
		{"cap_after_expand", "16", "第 9 个元素触发扩容 → capacity == 16"},
		{"size_after_expand", "9", "扩容后 size == 9"},
		{"get_8_after_expand", "8", "扩容后 get(8) == 8"},
		{"data_intact_after_expand", "true", "扩容后数据完整"},
	}

	for _, tc := range expandTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- removeAt ---

	removeTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"removeAt_return", "20", "removeAt(1) 返回 20"},
		{"size_after_remove", "3", "删除后 size == 3"},
		{"get_1_after_remove", "30", "删除后元素左移 get(1) == 30"},
		{"get_2_after_remove", "40", "删除后 get(2) == 40"},
		{"removeAt_oob", "-1", "removeAt 越界返回 -1"},
	}

	for _, tc := range removeTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- 缩容 ---

	shrinkTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"cap_16_elements", "16", "16 个元素 capacity == 16"},
		{"size_after_removes", "3", "删除 13 个后 size == 3"},
		{"cap_after_shrink", "8", "缩容后 capacity == 8"},
		{"get_0_after_shrink", "0", "缩容后 get(0) == 0"},
		{"get_1_after_shrink", "1", "缩容后 get(1) == 1"},
		{"get_2_after_shrink", "2", "缩容后 get(2) == 2"},
	}

	for _, tc := range shrinkTests {
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
		{"size_empty", "0", "全部删除后 size == 0"},
		{"get_after_refill", "100", "重新添加后 get(0) == 100"},
		{"size_after_refill", "1", "重新添加后 size == 1"},
	}

	for _, tc := range refillTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	total := len(basicTests) + len(addGetTests) + len(oobTests) + len(setTests) +
		len(expandTests) + len(removeTests) + len(shrinkTests) + len(refillTests)
	logger.Successf("All %d S01 tests passed!", total)
	return nil
}
