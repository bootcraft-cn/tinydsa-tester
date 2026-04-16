package stages

import (
	"fmt"
	"time"

	"github.com/bootcraft-cn/tester-utils/runner"
	"github.com/bootcraft-cn/tester-utils/test_case_harness"
	"github.com/bootcraft-cn/tester-utils/tester_definition"
	"github.com/bootcraft-cn/tinydsa-tester/internal/helpers"
)

func s05StackAndQueueTestCase() tester_definition.TestCase {
	return tester_definition.TestCase{
		Slug:        "stack-and-queue",
		Timeout:     30 * time.Second,
		TestFunc:    testS05StackAndQueue,
		CompileStep: autoCompileStep("TestS05", "test_s05", "test_s05", "testS05"),
	}
}

func testS05StackAndQueue(harness *test_case_harness.TestCaseHarness) error {
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

	// --- Stack: 基本操作 ---

	stackBasicTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"stack_initial_size", "0", "栈初始 size == 0"},
	}

	for _, tc := range stackBasicTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Stack: push + pop (LIFO) ---

	stackLIFOTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"stack_pop_1", "3", "push 1,2,3 → pop 返回 3"},
		{"stack_pop_2", "2", "pop 返回 2"},
		{"stack_pop_3", "1", "pop 返回 1"},
		{"stack_pop_empty", "-1", "空栈 pop 返回 -1"},
	}

	for _, tc := range stackLIFOTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Stack: peek ---

	stackPeekTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"stack_peek", "20", "peek 返回栈顶 20"},
		{"stack_size_after_peek", "2", "peek 不改变 size"},
		{"stack_peek_empty", "-1", "空栈 peek 返回 -1"},
	}

	for _, tc := range stackPeekTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Stack: size ---

	stackSizeTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"stack_size_3", "3", "push 3 次后 size == 3"},
		{"stack_size_after_pop", "2", "pop 后 size == 2"},
	}

	for _, tc := range stackSizeTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Stack: 交替 push/pop ---

	stackAlternateTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"stack_alternate_pop_1", "2", "交替 pop 返回 2"},
		{"stack_alternate_peek", "3", "交替 peek 返回 3"},
		{"stack_alternate_pop_2", "3", "交替 pop 返回 3"},
		{"stack_alternate_pop_3", "1", "交替 pop 返回 1"},
		{"stack_alternate_size", "0", "交替操作后 size == 0"},
	}

	for _, tc := range stackAlternateTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Queue: 基本操作 ---

	queueBasicTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"queue_initial_size", "0", "队列初始 size == 0"},
	}

	for _, tc := range queueBasicTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Queue: enqueue + dequeue (FIFO) ---

	queueFIFOTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"queue_dequeue_1", "1", "enqueue 1,2,3 → dequeue 返回 1"},
		{"queue_dequeue_2", "2", "dequeue 返回 2"},
		{"queue_dequeue_3", "3", "dequeue 返回 3"},
		{"queue_dequeue_empty", "-1", "空队列 dequeue 返回 -1"},
	}

	for _, tc := range queueFIFOTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Queue: front ---

	queueFrontTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"queue_front", "10", "front 返回队首 10"},
		{"queue_size_after_front", "2", "front 不改变 size"},
		{"queue_front_empty", "-1", "空队列 front 返回 -1"},
	}

	for _, tc := range queueFrontTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Queue: size ---

	queueSizeTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"queue_size_3", "3", "enqueue 3 次后 size == 3"},
		{"queue_size_after_dequeue", "2", "dequeue 后 size == 2"},
	}

	for _, tc := range queueSizeTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	// --- Queue: 交替 enqueue/dequeue ---

	queueAlternateTests := []struct {
		name     string
		expected string
		label    string
	}{
		{"queue_alternate_dequeue_1", "1", "交替 dequeue 返回 1"},
		{"queue_alternate_front", "2", "交替 front 返回 2"},
		{"queue_alternate_dequeue_2", "2", "交替 dequeue 返回 2"},
		{"queue_alternate_dequeue_3", "3", "交替 dequeue 返回 3"},
		{"queue_alternate_size", "0", "交替操作后 size == 0"},
	}

	for _, tc := range queueAlternateTests {
		if err := helpers.AssertEqual(results, tc.name, tc.expected); err != nil {
			return err
		}
		logger.Successf("✓ %s", tc.label)
	}

	total := len(stackBasicTests) + len(stackLIFOTests) + len(stackPeekTests) + len(stackSizeTests) +
		len(stackAlternateTests) + len(queueBasicTests) + len(queueFIFOTests) + len(queueFrontTests) +
		len(queueSizeTests) + len(queueAlternateTests)
	logger.Successf("All %d S05 tests passed!", total)
	return nil
}
