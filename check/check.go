package check

type Check interface {
	Check() CheckResult
	Description() string
}
