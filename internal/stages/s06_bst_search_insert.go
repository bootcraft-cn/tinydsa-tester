package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinydsa-tester/internal/helpers"
)

func s06BSTSearchInsertTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "bst-search-insert",
		Timeout:     30 * time.Second,
		TestFunc:    testS06BSTSearchInsert,
		CompileStep: autoCompileStep("TestS06", "test_s06", "test_s06", "testS06"),
	}
}

func testS06BSTSearchInsert(harness *test_case_harness.TestCaseHarness) error {
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

	// --- Group 1: 空树 ---

	emptyTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_empty_size", "0", "空树 size == 0"},
		{"bst_empty_get", "-1", "空树 get 返回 -1"},
		{"bst_empty_contains", "false", "空树 contains 返回 false"},
		{"bst_empty_keys", "", "空树 keys 返回空"},
	}

	for _, tc := range emptyTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 2: 单节点插入 ---

	singleTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_single_size", "1", "单节点 size == 1"},
		{"bst_single_get", "42", "单节点 get 返回 42"},
		{"bst_single_contains", "true", "单节点 contains 返回 true"},
		{"bst_single_missing", "false", "不存在的 key contains 返回 false"},
	}

	for _, tc := range singleTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 3: 多节点插入 + 中序遍历 ---

	multiTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_multi_size", "5", "5 个节点 size == 5"},
		{"bst_multi_keys", "ant,bee,cat,dog,elk", "中序遍历返回字典序"},
		{"bst_multi_get_first", "2", "get 最小 key 返回 2"},
		{"bst_multi_get_last", "5", "get 最大 key 返回 5"},
	}

	for _, tc := range multiTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 4: 更新已有 key ---

	updateTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_update_get", "999", "更新后 get 返回新值 999"},
		{"bst_update_size", "5", "更新不改变 size"},
		{"bst_update_keys", "ant,bee,cat,dog,elk", "更新不改变 keys"},
	}

	for _, tc := range updateTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 5: Get/contains 边界 ---

	edgeTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_get_missing", "-1", "get 不存在的 key 返回 -1"},
		{"bst_contains_existing", "true", "contains 已有 key 返回 true"},
		{"bst_contains_missing", "false", "contains 不存在 key 返回 false"},
	}

	for _, tc := range edgeTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 6: 左偏树 ---

	leftTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_left_keys", "a,b,c,d,e", "左偏树中序遍历正确"},
		{"bst_left_size", "5", "左偏树 size == 5"},
		{"bst_left_get", "1", "左偏树 get 最深节点"},
	}

	for _, tc := range leftTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 7: 右偏树 ---

	rightTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_right_keys", "a,b,c,d,e", "右偏树中序遍历正确"},
		{"bst_right_size", "5", "右偏树 size == 5"},
		{"bst_right_get", "5", "右偏树 get 最深节点"},
	}

	for _, tc := range rightTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 8: 大树 ---

	largeTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_large_size", "10", "10 节点 size == 10"},
		{"bst_large_keys", "b,d,f,g,j,m,n,p,t,w", "10 节点中序遍历正确"},
		{"bst_large_get", "10", "大树 get 返回正确值"},
		{"bst_large_contains", "true", "大树 contains 返回 true"},
	}

	for _, tc := range largeTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Group 9: 混合操作 ---

	mixedTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"bst_mixed_overwrite", "300", "覆盖写后 get 返回最新值"},
		{"bst_mixed_final_size", "3", "覆盖写不增加 size"},
	}

	for _, tc := range mixedTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	logger.Successf("All 30 S06 BST tests passed!")
	return nil
}
