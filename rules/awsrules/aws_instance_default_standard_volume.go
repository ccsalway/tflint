package awsrules

import (
	"log"

	"github.com/hashicorp/hcl2/hcl"
	"github.com/wata727/tflint/issue"
	"github.com/wata727/tflint/tflint"
)

// AwsInstanceDefaultStandardVolumeRule checks whether the volume type is unspecified
type AwsInstanceDefaultStandardVolumeRule struct {
	resourceType string
}

// NewAwsInstanceDefaultStandardVolumeRule returns new rule with default attributes
func NewAwsInstanceDefaultStandardVolumeRule() *AwsInstanceDefaultStandardVolumeRule {
	return &AwsInstanceDefaultStandardVolumeRule{
		resourceType: "aws_instance",
	}
}

// Name returns the rule name
func (r *AwsInstanceDefaultStandardVolumeRule) Name() string {
	return "aws_instance_default_standard_volume"
}

// Enabled returns whether the rule is enabled by default
func (r *AwsInstanceDefaultStandardVolumeRule) Enabled() bool {
	return true
}

// Check checks whether `volume_type` is defined for `root_block_device` or `ebs_block_device`
func (r *AwsInstanceDefaultStandardVolumeRule) Check(runner *tflint.Runner) error {
	log.Printf("[INFO] Check `%s` rule for `%s` runner", r.Name(), runner.TFConfigPath())

	err := runner.WalkResourceAttributes(r.resourceType, "root_block_device", func(attribute *hcl.Attribute) error {
		return r.walker(runner, attribute)
	})
	if err != nil {
		return err
	}
	return runner.WalkResourceAttributes(r.resourceType, "ebs_block_device", func(attribute *hcl.Attribute) error {
		return r.walker(runner, attribute)
	})
}

func (r *AwsInstanceDefaultStandardVolumeRule) walker(runner *tflint.Runner, attribute *hcl.Attribute) error {
	var val map[string]string
	err := runner.EvaluateExpr(attribute.Expr, &val)

	return runner.EnsureNoError(err, func() error {
		if _, ok := val["volume_type"]; !ok {
			runner.Issues = append(runner.Issues, &issue.Issue{
				Detector: r.Name(),
				Type:     issue.WARNING,
				Message:  "\"volume_type\" is not specified. Default standard volume type is not recommended. You can use \"gp2\", \"io1\", etc instead.",
				Line:     attribute.Range.Start.Line,
				File:     runner.GetFileName(attribute.Range.Filename),
				Link:     "https://github.com/wata727/tflint/blob/master/docs/aws_instance_default_standard_volume.md",
			})
		}
		return nil
	})
}