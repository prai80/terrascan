package writer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/accurics/terrascan/pkg/iac-providers/output"
	"github.com/accurics/terrascan/pkg/policy"
	"github.com/accurics/terrascan/pkg/results"
)

const (
	configOnlyTestOutputYAML = `aws_s3_bucket:
    - id: aws_s3_bucket.bucket
      name: bucket
      source: modules/m1/main.tf
      line: 20
      type: aws_s3_bucket
      config:
        bucket: ${module.m3.fullbucketname}
        policy: ${module.m2.fullbucketpolicy}`

	scanTestOutputYAML = `results:
    violations:
        - rule_name: s3EnforceUserACL
          description: S3 bucket Access is allowed to all AWS Account Users.
          rule_id: AWS.S3Bucket.DS.High.1043
          severity: HIGH
          category: S3
          resource_name: bucket
          resource_type: aws_s3_bucket
          file: modules/m1/main.tf
          line: 20
    count:
        low: 0
        medium: 0
        high: 1
        total: 1`
)

func TestYAMLWriter(t *testing.T) {
	type funcInput interface{}
	tests := []struct {
		name           string
		input          funcInput
		expectedOutput string
	}{
		{
			name: "YAML Writer: ResourceConfig",
			input: output.AllResourceConfigs{
				"aws_s3_bucket": []output.ResourceConfig{
					{
						ID:     "aws_s3_bucket.bucket",
						Name:   "bucket",
						Source: "modules/m1/main.tf",
						Line:   20,
						Type:   "aws_s3_bucket",
						Config: map[string]string{
							"bucket": "${module.m3.fullbucketname}",
							"policy": "${module.m2.fullbucketpolicy}",
						},
					},
				},
			},
			expectedOutput: configOnlyTestOutputYAML,
		},
		{
			name: "YAML Writer: Violations",
			input: policy.EngineOutput{
				ViolationStore: &results.ViolationStore{
					Violations: []*results.Violation{
						{
							RuleName:     "s3EnforceUserACL",
							Description:  "S3 bucket Access is allowed to all AWS Account Users.",
							RuleID:       "AWS.S3Bucket.DS.High.1043",
							Severity:     "HIGH",
							Category:     "S3",
							ResourceName: "bucket",
							ResourceType: "aws_s3_bucket",
							File:         "modules/m1/main.tf",
							LineNumber:   20,
						},
					},
					Count: results.ViolationStats{
						LowCount:    0,
						MediumCount: 0,
						HighCount:   1,
						TotalCount:  1,
					},
				},
			},
			expectedOutput: scanTestOutputYAML,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			writer := &bytes.Buffer{}
			YAMLWriter(tt.input, writer)
			if gotOutput := writer.String(); !strings.EqualFold(strings.TrimSpace(gotOutput), strings.TrimSpace(tt.expectedOutput)) {
				t.Errorf("YAMLWriter() = got: %v, want: %v", gotOutput, tt.expectedOutput)
			}
		})
	}
}
